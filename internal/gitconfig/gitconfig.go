// Package gitconfig reads and writes aliases in the user's global Git config.
package gitconfig

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Alias is a single alias entry. Git allows the same name to appear more than
// once, so lists of Alias may contain duplicate names.
type Alias struct {
	Name    string
	Command string
}

// Replacement describes an alias whose current value differs from the one
// about to be written.
type Replacement struct {
	Name    string
	Current []string
	Next    string
}

// DefaultAliases are the built-in aliases offered by the CLI, in display order.
var DefaultAliases = []Alias{
	{"st", "status --short --branch"},
	{"co", "checkout"},
	{"br", "branch"},
	{"cp", "cherry-pick"},
	{"ci", "commit"},
	{"last", "log -1 HEAD"},
	{"lg", "log --oneline --graph --decorate --all"},
	{"sp", "stash push"},
	{"spp", "stash pop"},
	{"sl", "stash list"},
	{"theirs", "checkout --theirs -- ."},
	{"unstage", "reset HEAD --"},
}

// GitCommand is a common Git command offered when creating a custom alias.
type GitCommand struct {
	Command     string
	Description string
}

// CommonCommands are suggested starting points for custom aliases.
var CommonCommands = []GitCommand{
	{"status", "Show the working tree status"},
	{"status --short --branch", "Compact status with branch info"},
	{"add", "Stage file contents"},
	{"add -p", "Stage changes hunk by hunk"},
	{"commit", "Record staged changes"},
	{"commit -m", "Commit with an inline message"},
	{"commit --amend", "Rewrite the last commit"},
	{"commit --amend --no-edit", "Add staged changes to the last commit"},
	{"switch", "Switch branches"},
	{"switch -c", "Create and switch to a new branch"},
	{"checkout", "Switch branches or restore files"},
	{"restore", "Restore working tree files"},
	{"restore --staged", "Unstage files"},
	{"branch", "List, create, or delete branches"},
	{"branch -vv", "List branches with upstream info"},
	{"merge", "Join branches together"},
	{"rebase", "Reapply commits on another base"},
	{"rebase -i", "Interactively edit commits"},
	{"cherry-pick", "Apply changes from existing commits"},
	{"revert", "Undo commits with new commits"},
	{"reset", "Move HEAD or unstage changes"},
	{"reset --soft HEAD~1", "Undo the last commit, keep changes staged"},
	{"stash push", "Save local changes away"},
	{"stash pop", "Reapply the latest stash"},
	{"stash list", "List stashes"},
	{"log", "Show commit history"},
	{"log --oneline --graph --decorate --all", "Compact graph of all branches"},
	{"log -1 HEAD", "Show the last commit"},
	{"diff", "Show unstaged changes"},
	{"diff --staged", "Show staged changes"},
	{"show", "Show a commit"},
	{"fetch --prune", "Download objects and prune deleted branches"},
	{"pull", "Fetch and integrate changes"},
	{"pull --rebase", "Fetch and rebase local commits"},
	{"push", "Upload local commits"},
	{"push --force-with-lease", "Force-push only if the remote is unchanged"},
	{"remote -v", "List remotes with URLs"},
	{"tag", "List or create tags"},
	{"blame", "Show who changed each line"},
	{"reflog", "Show where HEAD has been"},
	{"clean -fd", "Delete untracked files and folders"},
	{"shortlog -sn", "Commit counts per author"},
	{"worktree list", "List linked worktrees"},
}

var aliasNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]*$`)

// ValidateAliasName reports whether name is a usable alias name.
func ValidateAliasName(name string) error {
	if !aliasNamePattern.MatchString(name) {
		return errors.New("Alias names must start with a letter and contain only letters, numbers, or hyphens.")
	}
	return nil
}

// ValidateAliasCommand reports whether command is a usable alias command.
func ValidateAliasCommand(command string) error {
	if strings.TrimSpace(command) == "" {
		return errors.New("Alias commands cannot be empty.")
	}
	return nil
}

// Git runs git commands. The zero value runs "git" from PATH with the
// current process environment.
type Git struct {
	Executable string   // defaults to "git"
	Env        []string // nil means inherit the current environment
}

type result struct {
	stdout, stderr string
	code           int // -1 when the process was terminated by a signal
}

func (g Git) run(args ...string) (result, error) {
	exe := g.Executable
	if exe == "" {
		exe = "git"
	}
	cmd := exec.Command(exe, args...)
	cmd.Env = g.Env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	res := result{stdout: stdout.String(), stderr: stderr.String()}
	var exitErr *exec.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exitErr):
		res.code = exitErr.ExitCode()
	default:
		return res, err
	}
	return res, nil
}

func gitError(res result, args []string) error {
	if detail := strings.TrimSpace(res.stderr); detail != "" {
		return errors.New(detail)
	}
	status := fmt.Sprintf("exited with code %d", res.code)
	if res.code < 0 {
		status = "terminated"
	}
	return fmt.Errorf("git %s %s.", strings.Join(args, " "), status)
}

func (g Git) mustSucceed(args ...string) (result, error) {
	res, err := g.run(args...)
	if err != nil {
		return res, err
	}
	if res.code != 0 {
		return res, gitError(res, args)
	}
	return res, nil
}

// Check returns the output of "git --version", or an error if Git is unusable.
func (g Git) Check() (string, error) {
	res, err := g.mustSucceed("--version")
	return strings.TrimSpace(res.stdout), err
}

// InstallAliases writes each alias to the global config, replacing all
// existing values for that name.
func (g Git) InstallAliases(aliases []Alias) error {
	for _, a := range aliases {
		if err := ValidateAliasName(a.Name); err != nil {
			return err
		}
		if err := ValidateAliasCommand(a.Command); err != nil {
			return err
		}
		if _, err := g.mustSucceed("config", "--global", "--replace-all", "alias."+a.Name, a.Command); err != nil {
			return err
		}
	}
	return nil
}

// ListGlobalAliases returns every alias in the global config, in file order.
func (g Git) ListGlobalAliases() ([]Alias, error) {
	args := []string{"config", "--global", "--get-regexp", `^alias\.`}
	res, err := g.run(args...)
	if err != nil {
		return nil, err
	}
	if res.code == 1 && strings.TrimSpace(res.stdout) == "" {
		return []Alias{}, nil
	}
	if res.code != 0 {
		return nil, gitError(res, args)
	}

	aliases := []Alias{}
	for _, line := range strings.Split(res.stdout, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if line == "" {
			continue
		}
		key, command, _ := strings.Cut(line, " ")
		if strings.HasPrefix(strings.ToLower(key), "alias.") {
			key = key[len("alias."):]
		}
		aliases = append(aliases, Alias{Name: key, Command: command})
	}
	return aliases, nil
}

// RemoveGlobalAlias deletes every value of the named alias.
func (g Git) RemoveGlobalAlias(name string) error {
	if err := ValidateAliasName(name); err != nil {
		return err
	}
	_, err := g.mustSucceed("config", "--global", "--unset-all", "alias."+name)
	return err
}

// GroupAliases merges duplicate names, keeping first-seen order.
func GroupAliases(aliases []Alias) (names []string, commands map[string][]string) {
	commands = map[string][]string{}
	for _, a := range aliases {
		if _, seen := commands[a.Name]; !seen {
			names = append(names, a.Name)
		}
		commands[a.Name] = append(commands[a.Name], a.Command)
	}
	return names, commands
}

// FindAliasReplacements lists desired aliases that already exist with a
// different value, or with several values.
func FindAliasReplacements(current, desired []Alias) []Replacement {
	_, byName := GroupAliases(current)
	replacements := []Replacement{}
	for _, d := range desired {
		cur, ok := byName[d.Name]
		if !ok {
			continue
		}
		if !allEqual(cur) || cur[0] != d.Command {
			replacements = append(replacements, Replacement{Name: d.Name, Current: cur, Next: d.Command})
		}
	}
	return replacements
}

func allEqual(values []string) bool {
	for _, v := range values[1:] {
		if v != values[0] {
			return false
		}
	}
	return true
}
