// Command git-aliases interactively manages global Git aliases.
package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"git-aliases/internal/gitconfig"
)

var (
	// errCancelled means the user pressed Ctrl+C; the CLI exits quietly.
	errCancelled = errors.New("cancelled")
	// errBack means the user pressed Esc; the current step is abandoned.
	errBack = errors.New("back")
)

var (
	git = gitconfig.Git{}

	dim     = lipgloss.NewStyle().Faint(true)
	title   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")).Padding(0, 1)
	box     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).Padding(0, 1)
	infoSym = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render("●")
	okSym   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("◆")
	errSym  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("■")
)

var (
	accent    = lipgloss.Color("#A78BFA") // matches the banner's purple
	promptBox = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accent).Padding(1, 2)
	hintKey   = lipgloss.NewStyle().Foreground(lipgloss.Color("#C4B5FD")).Bold(true)
	hintRule  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// theme is huh's Charm theme with bolder questions and a gap below each one.
func theme() *huh.Theme {
	t := huh.ThemeCharm()
	t.Focused.Title = t.Focused.Title.Foreground(accent).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.MarginBottom(1)
	t.Blurred.Title = t.Blurred.Title.MarginBottom(1)
	t.Blurred.Description = t.Blurred.Description.MarginBottom(1)
	return t
}

// keys renders a footer such as "↑/↓ move · Enter select" from alternating
// key/action pairs.
func keys(pairs ...string) string {
	var parts []string
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, hintKey.Render(pairs[i])+" "+dim.Render(pairs[i+1]))
	}
	return strings.Join(parts, dim.Render("  ·  "))
}

func logInfo(msg string)    { fmt.Println(infoSym, msg) }
func logSuccess(msg string) { fmt.Println(okSym, msg) }
func logError(msg string)   { fmt.Fprintln(os.Stderr, errSym, msg) }

func note(body, heading string) {
	fmt.Println(box.Render(lipgloss.NewStyle().Bold(true).Render(heading) + "\n\n" + body))
}

// escForm wraps a huh form so Esc leaves the prompt instead of being ignored.
type escForm struct {
	form     *huh.Form
	hint     string // key hints shown at the bottom of the box
	width    int    // inner width of the box
	back     bool
	onUpdate func() // optional; runs after every update to sync fields
}

// searching reports whether the focused select list is in search mode, where
// Esc should clear the search instead of going back.
func (m *escForm) searching() bool {
	f, ok := m.form.GetFocusedField().(interface{ GetFiltering() bool })
	return ok && f.GetFiltering()
}

func (m *escForm) Init() tea.Cmd { return m.form.Init() }

func (m *escForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && k.Type == tea.KeyEsc && !m.searching() {
		m.back = true
		return m, tea.Quit
	}
	_, cmd := m.form.Update(msg)
	if m.onUpdate != nil {
		m.onUpdate()
	}
	return m, cmd
}

// View draws the form inside a rounded box with the key hints as a footer.
// Once the form is finished it draws nothing, so the box disappears.
func (m *escForm) View() string {
	if m.back || m.form.State != huh.StateNormal {
		return ""
	}
	body := strings.TrimRight(m.form.View(), "\n")
	if m.hint != "" {
		rule := hintRule.Render(strings.Repeat("─", m.width))
		hint := lipgloss.NewStyle().Width(m.width).Render(m.hint)
		body += "\n\n" + rule + "\n" + hint
	}
	return promptBox.Render(body)
}

// boxWidth is the usable width inside a prompt box for the current terminal.
func boxWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		w = 80
	}
	// Leave room for the border (2) and padding (4), capped for readability.
	return max(min(w-6, 100), 20)
}

// prompt runs a single huh field. Ctrl+C returns errCancelled and Esc
// returns errBack. hint is shown at the bottom of the box.
func prompt(field huh.Field, hint string) error {
	return runForm(huh.NewForm(huh.NewGroup(field)), hint, nil)
}

// runForm runs a form with the same Esc/Ctrl+C handling as prompt. onUpdate,
// if set, runs after every key press or event.
func runForm(form *huh.Form, hint string, onUpdate func()) error {
	width := boxWidth()
	form.WithShowHelp(false).WithTheme(theme()).WithWidth(width)
	form.SubmitCmd = tea.Quit
	form.CancelCmd = tea.Interrupt

	m := &escForm{form: form, hint: hint, width: width, onUpdate: onUpdate}
	_, err := tea.NewProgram(m).Run()
	switch {
	case m.back:
		return errBack
	case form.State == huh.StateAborted || errors.Is(err, tea.ErrInterrupted):
		fmt.Println(dim.Render("Cancelled."))
		return errCancelled
	case err != nil:
		return err
	}
	return nil
}

func confirm(message string) (bool, error) {
	var ok bool
	err := prompt(huh.NewConfirm().Title(message).Affirmative("Yes").Negative("No").Value(&ok),
		keys("←/→", "choose", "y/n", "answer", "Enter", "confirm", "Esc", "back"))
	return ok, err
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func confirmReplacements(aliases []gitconfig.Alias) (bool, error) {
	current, err := git.ListGlobalAliases()
	if err != nil {
		return false, err
	}
	replacements := gitconfig.FindAliasReplacements(current, aliases)
	if len(replacements) == 0 {
		return true, nil
	}
	lines := make([]string, len(replacements))
	for i, r := range replacements {
		lines[i] = fmt.Sprintf("git %s: %s -> %s", r.Name, strings.Join(r.Current, " | "), r.Next)
	}
	note(strings.Join(lines, "\n"), "Existing aliases that will change")
	return confirm("Replace these aliases in your global Git config?")
}

func saveAliases(aliases []gitconfig.Alias, successMessage string) error {
	ok, err := confirmReplacements(aliases)
	if errors.Is(err, errCancelled) || errors.Is(err, errBack) {
		return err
	}
	if err == nil && !ok {
		logInfo("No changes made.")
		return nil
	}
	if err == nil {
		err = git.InstallAliases(aliases)
	}
	if err != nil {
		logError("Could not update global Git config: " + err.Error())
		return nil
	}
	logSuccess(successMessage)
	return nil
}

var multiSelectKeys = keys("↑/↓", "move", "Space", "toggle", "Ctrl+A", "all/none", "Enter", "confirm", "Esc", "back")

func installBuiltIns() error {
	options := make([]huh.Option[string], len(gitconfig.DefaultAliases))
	for i, a := range gitconfig.DefaultAliases {
		options[i] = huh.NewOption(fmt.Sprintf("git %-8s %s", a.Name, dim.Render("git "+a.Command)), a.Name)
	}
	var selected []string
	if err := prompt(huh.NewMultiSelect[string]().
		Title("Choose built-in aliases to install").
		Options(options...).
		Value(&selected), multiSelectKeys); err != nil {
		return err
	}
	if len(selected) == 0 {
		logInfo("No aliases selected.")
		return nil
	}

	var aliases []gitconfig.Alias
	for _, a := range gitconfig.DefaultAliases {
		if slices.Contains(selected, a.Name) {
			aliases = append(aliases, a)
		}
	}
	return saveAliases(aliases, fmt.Sprintf("Installed %d %s in global Git config.",
		len(aliases), plural(len(aliases), "alias", "aliases")))
}

// ownCommand is the select value for "type the command yourself".
const ownCommand = "\x00own"

func commandOptions() []huh.Option[string] {
	width := 0
	for _, c := range gitconfig.CommonCommands {
		width = max(width, len(c.Command))
	}
	options := []huh.Option[string]{huh.NewOption("✎ Type my own command", ownCommand)}
	for _, c := range gitconfig.CommonCommands {
		label := fmt.Sprintf("%-*s  %s", width, c.Command, dim.Render(c.Description))
		options = append(options, huh.NewOption(label, c.Command))
	}
	return options
}

// addCustomAlias shows the alias name, a command list, and the editable
// command on one screen. Moving through the list fills in the command field.
// Enter/Tab moves to the next field, Shift+Tab to the previous one, and Esc
// returns to the main menu.
func addCustomAlias() error {
	var name, picked, command string

	commandInput := huh.NewInput().
		TitleFunc(func() string {
			if n := strings.TrimSpace(name); n != "" {
				return "Command for git " + n
			}
			return "Command"
		}, &name).
		Description("Add options if you like, then Enter to save").
		Placeholder("log --oneline").
		Validate(gitconfig.ValidateAliasCommand).
		Value(&command)

	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Alias name").
			Placeholder("e.g. lg").
			Validate(func(v string) error { return gitconfig.ValidateAliasName(strings.TrimSpace(v)) }).
			Value(&name),
		huh.NewSelect[string]().
			Title("Git command").
			Description("Pick one to start from, or type your own").
			Options(commandOptions()...).
			Height(14).
			Value(&picked),
		commandInput,
	))

	// Whenever the highlighted list item changes, copy it into the command field.
	lastPicked := picked
	syncCommand := func() {
		if picked == lastPicked {
			return
		}
		lastPicked = picked
		command = picked
		if picked == ownCommand {
			command = ""
		}
		commandInput.Value(&command) // reloads the field's text from command
	}
	syncCommand()

	hint := keys("Enter/Tab", "next", "Shift+Tab", "previous", "↑/↓", "move in list", "/", "search list", "Esc", "back")
	if err := runForm(form, hint, syncCommand); err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)
	note("git "+name+" "+command, "Command preview")
	return saveAliases([]gitconfig.Alias{{Name: name, Command: command}}, "Saved git "+name+".")
}

func manageAliases() error {
	all, err := git.ListGlobalAliases()
	if err != nil {
		logError("Could not read global Git aliases: " + err.Error())
		return nil
	}
	names, commands := gitconfig.GroupAliases(all)
	if len(names) == 0 {
		logInfo("No global aliases found.")
		return nil
	}

	var filter string
	if err := prompt(huh.NewInput().
		Title("Filter aliases by name or command").
		Placeholder("Leave blank to show all").
		Value(&filter), keys("Enter", "continue", "Esc", "back")); err != nil {
		return err
	}
	query := strings.ToLower(strings.TrimSpace(filter))
	var visible []string
	for _, n := range names {
		if query == "" || strings.Contains(strings.ToLower(n), query) ||
			slices.ContainsFunc(commands[n], func(c string) bool { return strings.Contains(strings.ToLower(c), query) }) {
			visible = append(visible, n)
		}
	}
	if len(visible) == 0 {
		logInfo(fmt.Sprintf("No aliases match %q.", filter))
		return nil
	}

	lines := make([]string, len(visible))
	options := make([]huh.Option[string], len(visible))
	for i, n := range visible {
		joined := strings.Join(commands[n], " | ")
		lines[i] = fmt.Sprintf("git %-12s git %s", n, joined)
		options[i] = huh.NewOption(fmt.Sprintf("git %-8s %s", n, dim.Render("git "+joined)), n)
	}
	note(strings.Join(lines, "\n"), fmt.Sprintf("Global aliases (%d)", len(visible)))

	var selected []string
	if err := prompt(huh.NewMultiSelect[string]().
		Title("Select aliases to remove, or press Enter to keep them").
		Options(options...).
		Value(&selected), multiSelectKeys); err != nil {
		return err
	}
	if len(selected) == 0 {
		logInfo("No aliases selected.")
		return nil
	}

	word := plural(len(selected), "alias", "aliases")
	ok, err := confirm(fmt.Sprintf("Remove %d selected %s?", len(selected), word))
	if err != nil {
		return err
	}
	if !ok {
		logInfo("No changes made.")
		return nil
	}
	for _, n := range selected {
		if err := git.RemoveGlobalAlias(n); err != nil {
			logError("Could not remove aliases: " + err.Error())
			return nil
		}
	}
	logSuccess(fmt.Sprintf("Removed %d %s from global Git config.", len(selected), word))
	return nil
}

// menuOption lines up main menu descriptions in a column.
func menuOption(label, description, value string) huh.Option[string] {
	return huh.NewOption(fmt.Sprintf("%-32s %s", label, dim.Render(description)), value)
}

func run() int {
	gitVersion, err := git.Check()
	if err != nil {
		logError("Git is unavailable: " + err.Error())
		return 1
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		logError("Run this command in an interactive terminal.")
		return 1
	}

	printBanner(gitVersion)
	for {
		var action string
		err := prompt(huh.NewSelect[string]().
			Title("What would you like to do?").
			Options(
				menuOption("Install built-in aliases", "Choose shortcuts for common Git commands", "install"),
				menuOption("Add or update a custom alias", "Pick a command and preview it before saving", "custom"),
				menuOption("View, filter, or remove aliases", "Manage your global Git config", "manage"),
				huh.NewOption("Quit", "quit"),
			).
			Value(&action), keys("↑/↓", "move", "Enter", "select", "Esc", "quit"))
		if errors.Is(err, errBack) {
			// Esc on the main menu quits.
			fmt.Println(dim.Render("Done."))
			return 0
		}
		if err == nil {
			switch action {
			case "quit":
				fmt.Println(dim.Render("Done."))
				return 0
			case "install":
				err = installBuiltIns()
			case "custom":
				err = addCustomAlias()
			case "manage":
				err = manageAliases()
			}
		}
		if errors.Is(err, errCancelled) {
			return 0
		}
		if errors.Is(err, errBack) {
			continue // Esc inside a sub-menu returns to the main menu.
		}
		if err != nil {
			logError(err.Error())
			return 1
		}
	}
}

func main() {
	// "git aliases --help" is intercepted by Git (it opens a man page), so
	// -h is the help flag that works when run as a Git subcommand.
	args := os.Args[1:]
	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		fmt.Println("Usage: git aliases [-h | --version]\n\nInteractively manage global Git aliases.")
		return
	}
	if slices.Contains(args, "--version") {
		fmt.Println("git-aliases", version)
		return
	}
	os.Exit(run())
}
