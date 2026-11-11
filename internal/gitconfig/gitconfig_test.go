package gitconfig

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestInstallListRemoveInIsolatedGlobalConfig(t *testing.T) {
	g := Git{Env: append(os.Environ(), "GIT_CONFIG_GLOBAL="+filepath.Join(t.TempDir(), "gitconfig"))}

	got, err := g.ListGlobalAliases()
	if err != nil || len(got) != 0 {
		t.Fatalf("initial list = %v, %v; want empty", got, err)
	}

	if err := g.InstallAliases([]Alias{{"cp", "cherry-pick"}, {"sp", "stash push"}}); err != nil {
		t.Fatal(err)
	}
	got, err = g.ListGlobalAliases()
	want := []Alias{{"cp", "cherry-pick"}, {"sp", "stash push"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("after install = %v, %v; want %v", got, err, want)
	}

	if err := g.RemoveGlobalAlias("cp"); err != nil {
		t.Fatal(err)
	}
	got, err = g.ListGlobalAliases()
	want = []Alias{{"sp", "stash push"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("after remove = %v, %v; want %v", got, err, want)
	}
}

func TestValidation(t *testing.T) {
	for _, bad := range []string{"bad alias", "-bad"} {
		if err := ValidateAliasName(bad); err == nil || !strings.Contains(err.Error(), "Alias names must start") {
			t.Errorf("ValidateAliasName(%q) = %v", bad, err)
		}
	}
	if err := ValidateAliasCommand("  "); err == nil || !strings.Contains(err.Error(), "commands cannot be empty") {
		t.Errorf("ValidateAliasCommand(blank) = %v", err)
	}
	if err := ValidateAliasName("log-graph2"); err != nil {
		t.Error(err)
	}
	if err := ValidateAliasCommand("log --oneline"); err != nil {
		t.Error(err)
	}
}

func TestFindAliasReplacements(t *testing.T) {
	got := FindAliasReplacements(
		[]Alias{{"co", "checkout"}, {"old", "first"}, {"old", "second"}},
		[]Alias{{"co", "switch"}, {"old", "first"}, {"st", "status --short"}},
	)
	want := []Replacement{
		{Name: "co", Current: []string{"checkout"}, Next: "switch"},
		{Name: "old", Current: []string{"first", "second"}, Next: "first"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCommonCommandsAreValidAndUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range CommonCommands {
		if err := ValidateAliasCommand(c.Command); err != nil || c.Description == "" {
			t.Errorf("bad entry %+v: %v", c, err)
		}
		if seen[c.Command] {
			t.Errorf("duplicate command %q", c.Command)
		}
		seen[c.Command] = true
	}
}
