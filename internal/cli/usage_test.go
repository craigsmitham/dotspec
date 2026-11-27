package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestUsageErrorsSuggestAndNameOperands(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"new", "requirment", "fast-save"}, `dotspec: unknown command "requirment" for "dotspec new"; did you mean "requirement"?` + "\n"},
		{[]string{"lst", "requirement"}, `dotspec: unknown command "lst" for "dotspec"; did you mean "list"?` + "\n"},
		{[]string{"help", "chek"}, `dotspec: unknown command "chek" for "dotspec"; did you mean "check"?` + "\n"},
		{[]string{"bogus"}, `dotspec: unknown command "bogus" for "dotspec"; see dotspec --help` + "\n"},
		{[]string{"new", "requirement"}, "dotspec: new requirement: missing <slug>; see dotspec new requirement --help\n"},
		{[]string{"new", "system", "Billing"}, "dotspec: new system: missing --part-of <path>; see dotspec new system --help\n"},
		{[]string{"requirment", "--help"}, `dotspec: unknown command "requirment" for "dotspec"; see dotspec --help` + "\n"},
		{[]string{"new", "requirment", "--help"}, `dotspec: unknown command "requirment" for "dotspec new"; did you mean "requirement"?` + "\n"},
		{[]string{"-C", ".", "snapshot", "exprt", "-h"}, `dotspec: unknown command "exprt" for "dotspec snapshot"; did you mean "export"?` + "\n"},
		{[]string{"snapshot", "validate"}, "dotspec: snapshot validate: missing <file>; see dotspec snapshot validate --help\n"},
		{[]string{"check", "a", "b"}, `dotspec: check: unexpected argument "b"; see dotspec check --help` + "\n"},
		{[]string{"sources", "a", "b"}, `dotspec: sources: unexpected argument "b"; see dotspec sources --help` + "\n"},
		{[]string{"check", "--bogus"}, "dotspec: unknown flag: --bogus; see dotspec check --help\n"},
	} {
		code, out, errOut := invoke(tc.args...)
		if code != 2 || out != "" || errOut != tc.want {
			t.Errorf("%v: %d %q %q, want %q", tc.args, code, out, errOut, tc.want)
		}
	}
}

// Every command group rejects a misspelled subcommand rather than printing
// its help with status 0.
// new or list without a kind names each kind with its domain description and
// its handbook route, as their help does.
func TestMissingKindListsKinds(t *testing.T) {
	t.Chdir(t.TempDir())
	for _, verb := range []string{"new", "list"} {
		code, out, errOut := invoke(verb)
		if code != 2 || out != "" || !strings.HasPrefix(errOut, "dotspec: "+verb+": missing <kind>; see dotspec "+verb+" --help\n\nKinds:\n") {
			t.Fatalf("%s: %d %q %q", verb, code, out, errOut)
		}
		for _, k := range kinds() {
			if !strings.Contains(errOut, "  "+k.layout.Kind+" ") || !strings.Contains(errOut, k.short+"\n") || !strings.Contains(errOut, "dotspec docs "+k.layout.Kind+"\n") {
				t.Errorf("%s: kind %s is not listed with its description and route: %q", verb, k.layout.Kind, errOut)
			}
		}
	}
}

func TestEveryGroupRejectsUnknownSubcommands(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		if cmd.HasAvailableSubCommands() {
			if err := cmd.ValidateArgs([]string{"zzz"}); err == nil || !strings.Contains(err.Error(), `unknown command "zzz"`) {
				t.Errorf("%s accepts an unknown subcommand: %v", cmd.CommandPath(), err)
			}
		}
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
}

func TestFilesystemErrorsOmitSystemCallsAndAbsolutePaths(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	if err := os.Mkdir("specs", 0o755); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "absent")
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"-C", missing, "list", "requirement"}, "dotspec: --chdir absent: directory does not exist\n"},
		{[]string{"snapshot", "validate", "absent.json"}, "dotspec: absent.json: path does not exist\n"},
		{[]string{"docs", "guides/absent"}, `dotspec: cannot open "guides/absent": no such page; run dotspec docs to browse topics` + "\n"},
	} {
		code, out, errOut := invoke(tc.args...)
		if code != 2 || out != "" || errOut != tc.want {
			t.Errorf("%v: %d %q %q, want %q", tc.args, code, out, errOut, tc.want)
		}
	}
	if code, _, errOut := invoke("new", "requirement", "taken"); code != 0 || errOut != "" {
		t.Fatalf("init: %d %q", code, errOut)
	}
	_, _, errOut := invoke("new", "requirement", "taken")
	if want := "dotspec: " + filepath.Join("specs", "requirements", "taken.requirement.spec.md") + " already exists; refusing to overwrite it\n"; errOut != want {
		t.Errorf("overwrite: %q, want %q", errOut, want)
	}
}

func TestRelativePaths(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	cwd := filepath.Join(dir, "a", "b")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(cwd)
	sep := string(filepath.Separator)
	for message, want := range map[string]string{
		cwd + sep + "specs" + sep + "x.md exists":    "specs" + sep + "x.md exists",
		"source must be inside " + cwd + ".":         "source must be inside ." + sep + ".",
		"source must be inside " + dir + sep + "a":   "source must be inside ..",
		dir + sep + "a" + sep + "c" + sep + "y.md":   ".." + sep + "c" + sep + "y.md",
		cwd + "-sibling" + sep + "z.md":              ".." + sep + "b-sibling" + sep + "z.md",
		"see " + sep + "elsewhere" + sep + "file.md": "see " + sep + "elsewhere" + sep + "file.md",
	} {
		if got := relativePaths(message); got != want {
			t.Errorf("%q: got %q, want %q", message, got, want)
		}
	}
}
