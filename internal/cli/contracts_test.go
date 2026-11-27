package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func prepareExample(t *testing.T, command string) {
	t.Helper()
	t.Chdir(t.TempDir())
	if output, err := exec.Command("git", "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	if err := os.Mkdir("internal", 0o755); err != nil {
		t.Fatal(err)
	}
	if command == "dotspec" || command == "dotspec init" {
		return
	}
	if _, err := system.InitNamed(".", "New system"); err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(command, "dotspec new") {
		return
	}
	if command == "dotspec build" {
		writeFixture(t, "dotspec.yaml", "name: demo\n")
		return
	}
	writeFixture(t, "specs/features/resume-work.feature.spec.md", "---\ntype: feature-spec\n---\n# Resume work\n")
	writeFixture(t, "specs/capabilities/local-storage.capability.spec.md", "---\ntype: capability-spec\n---\n# Local storage\n")
	writeFixture(t, "specs/jobs/develop-requirements.job.spec.md", jobFixture)
	writeFixture(t, "specs/use-cases/check-specifications.use-case.spec.md", usecase.Template)
	writeFixture(t, "specs/concepts/job-map.concept.spec.md", concept.Template)
	writeFixture(t, "specs/business-rules/purchase-approval.business-rule.spec.md", businessrule.Template)
	writeFixture(t, "specs/requirements/diagnostic-severity.requirement.spec.md", requirement.Template)
	writeFixture(t, "specs/requirements/reject-existing-target.requirement.spec.md", requirement.Template)
	if _, err := qualitymodel.Init(".", "organization-quality", "iso-25010-2023"); err != nil {
		t.Fatal(err)
	}
}

// exampleArgs splits an example invocation into arguments. A double-quoted
// operand, such as a system name, stays one argument.
func exampleArgs(example string) []string {
	var args []string
	for i, part := range strings.Split(example, `"`) {
		if i%2 == 1 {
			args = append(args, part)
			continue
		}
		args = append(args, strings.Fields(part)...)
	}
	return args
}

func TestDeclaredCommandExamples(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		if cmd.Example != "" {
			t.Run(cmd.CommandPath(), func(t *testing.T) {
				prepareExample(t, cmd.CommandPath())
				for _, example := range strings.Split(cmd.Example, "\n") {
					args := exampleArgs(example)
					if len(args) == 0 || args[0] != "dotspec" {
						t.Fatalf("example must be an executable dotspec invocation: %s", example)
					}
					if code, out, errOut := invoke(args[1:]...); code != 0 || errOut != "" {
						t.Fatalf("%s: %d %s %s", example, code, out, errOut)
					}
				}
			})
		}
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
}

func TestStructuredReadContracts(t *testing.T) {
	prepareExample(t, "read")
	for _, args := range [][]string{
		{"sources"}, {"check"}, {"check", "specs/jobs/develop-requirements.job.spec.md"},
		{"new", "requirement", "fresh"}, {"new", "quality-attribute", "recoverability", "--from", "iso-25010-2023", "--dry-run"},
	} {
		code, out, errOut := invoke(append(args, "--json")...)
		if code != 0 || errOut != "" {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
		assertResultContract(t, args, out)
	}
	for _, k := range kinds() {
		args := []string{"list", k.layout.Kind}
		code, out, errOut := invoke(append(args, "--json")...)
		if code != 0 || errOut != "" {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
		assertResultContract(t, args, out)
		var result listResult
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatal(err)
		}
		for _, e := range result.Entries {
			if e.Kind != k.layout.Kind || !filepath.IsAbs(e.Path) || e.Name == "" || e.System != nil && !filepath.IsAbs(e.System.Path) {
				t.Errorf("%v: entry %+v", args, e)
			}
		}
		if k.layout.Kind != "component" && k.layout.Kind != "quality-attribute" && len(result.Entries) == 0 {
			t.Errorf("%v: no entries", args)
		}
	}
}

// assertResultContract requires a command result to begin with schema version 1, to
// carry no terminal styling and to give every location an absolute source path.
func assertResultContract(t *testing.T, args []string, out string) {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(out))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		t.Fatalf("%v: result is not an object: %s", args, out)
	}
	if key, err := decoder.Token(); err != nil || key != "schema_version" {
		t.Fatalf("%v: schema_version is not the first key: %s", args, out)
	}
	if version, err := decoder.Token(); err != nil || version != float64(1) {
		t.Fatalf("%v: schema_version %v, want 1", args, version)
	}
	if strings.Contains(out, "\x1b") {
		t.Fatalf("%v: styled JSON: %s", args, out)
	}
	var result any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	var visit func(string, any)
	visit = func(at string, value any) {
		switch value := value.(type) {
		case map[string]any:
			for key, child := range value {
				if key == "location" || key == "related" || key == "referenced_by" {
					locations, ok := child.([]any)
					if !ok {
						locations = []any{child}
					}
					for _, location := range locations {
						fields, _ := location.(map[string]any)
						path, _ := fields["path"].(string)
						if !filepath.IsAbs(path) && !strings.HasPrefix(path, "builtin:") {
							t.Errorf("%v: %s.%s has no absolute path: %v", args, at, key, location)
						}
					}
				}
				visit(at+"."+key, child)
			}
		case []any:
			for _, child := range value {
				visit(at+"[]", child)
			}
		}
	}
	visit("", result)
}

// --json is the only output switch. It is local to each command that returns a
// result, so help and documentation commands never accept a flag that does nothing.
func TestJSONFlagPlacement(t *testing.T) {
	noResult := map[string]bool{"dotspec docs": true, "dotspec help": true, "dotspec completion": true, "dotspec man": true,
		"dotspec snapshot export": true, "dotspec snapshot schema": true, "dotspec build": true}
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		path := cmd.CommandPath()
		if cmd.PersistentFlags().Lookup("json") != nil {
			t.Errorf("%s: --json must be a local flag", path)
		}
		if cmd.Flags().Lookup("format") != nil {
			t.Errorf("%s: --json is the only output switch; remove --format", path)
		}
		exempt := noResult[path] || cmd.HasParent() && noResult[cmd.Parent().CommandPath()]
		hasJSON := cmd.LocalFlags().Lookup("json") != nil
		switch {
		case exempt && hasJSON:
			t.Errorf("%s returns no result and must not accept --json", path)
		case !exempt && cmd.Runnable() && !cmd.HasSubCommands() && !hasJSON:
			t.Errorf("%s returns a result and needs a local --json flag", path)
		}
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
}

type failedOutput struct{}

func (failedOutput) Write([]byte) (int, error) { return 0, errors.New("output unavailable") }

func TestOutputFailuresAcrossCommands(t *testing.T) {
	prepareExample(t, "output")
	commands := [][]string{
		{"--help"}, {"--version"}, {"completion", "bash"}, {"man"}, {"sources"}, {"sources", "--json"},
		{"new", "quality-attribute", "recoverability", "--from", "iso-25010-2023", "--dry-run"},
		{"check", "specs/jobs/develop-requirements.job.spec.md", "--json"},
		{"snapshot", "export"}, {"snapshot", "schema"},
		{"check"}, {"check", "--json"},
	}
	for _, k := range kinds() {
		commands = append(commands, []string{"list", k.layout.Kind}, []string{"list", k.layout.Kind, "--json"})
	}
	for _, args := range commands {
		var stderr bytes.Buffer
		if code := Run(args, failedOutput{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "output unavailable") {
			t.Errorf("%v: %d %s", args, code, &stderr)
		}
	}
	if code := Run([]string{"unknown"}, &bytes.Buffer{}, failedOutput{}); code != 2 {
		t.Fatalf("stderr failure returned %d", code)
	}
}

func TestVersionMetadata(t *testing.T) {
	oldVersion, oldCommit := version, commit
	t.Cleanup(func() { version, commit = oldVersion, oldCommit })
	version, commit = "1.2.3", "abcdef1234567890"
	code, out, errOut := invoke("--version")
	if code != 0 || out != "dotspec 1.2.3 (abcdef1)\n" || errOut != "" {
		t.Fatalf("version: %d %q %q", code, out, errOut)
	}
}

func TestExitPolicyForSelectorsAndSpecifications(t *testing.T) {
	prepareExample(t, "selectors")
	for _, args := range [][]string{
		{"-C", "missing-scope", "list", "system"}, {"-C", "missing-scope", "new", "job", "new"},
	} {
		if code, out, errOut := invoke(args...); code != 2 || out != "" || errOut == "" {
			t.Fatalf("invocation/operation %v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{
		{"check", "missing-scope", "--json"}, {"check", "specs/jobs/absent.job.spec.md", "--json"},
	} {
		code, out, errOut := invoke(args...)
		if code != 2 || !strings.Contains(out, `"complete": false`) || errOut != "" {
			t.Fatalf("incomplete %v: %d %q %q", args, code, out, errOut)
		}
	}
	writeFixture(t, "specs/jobs/broken.job.spec.md", "broken")
	if code, out, errOut := invoke("list", "job"); code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
		t.Fatalf("specification: %d %q %q", code, out, errOut)
	}
}

// Usage lines write <x> for a required operand, [x] or [<x>] for an optional
// one and ... for repetition, always in lowercase. Flag help reads as a
// lowercase phrase without a closing period.
func TestHelpConventions(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		cmd.InitDefaultHelpFlag()
		if _, operands, _ := strings.Cut(cmd.Use, " "); operands != "" {
			if err := checkOperands(operands); err != nil {
				t.Errorf("%s: usage %q: %v", cmd.CommandPath(), cmd.Use, err)
			}
		}
		check := func(f *pflag.Flag) {
			first, _ := utf8.DecodeRuneInString(f.Usage)
			if !unicode.IsLower(first) || strings.HasSuffix(f.Usage, ".") {
				t.Errorf("%s --%s: help %q must start lowercase and have no closing period", cmd.CommandPath(), f.Name, f.Usage)
			}
			// A value flag names its value, as in --part-of <path> or --template brief|expanded.
			if value := f.Value.Type(); value != "bool" && !flagValueName.MatchString(value) {
				t.Errorf("%s --%s: value shows as %q; add it with stringFlag and a <name> or a|b placeholder", cmd.CommandPath(), f.Name, value)
			}
		}
		cmd.Flags().VisitAll(check)
		cmd.PersistentFlags().VisitAll(check)
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(NewRoot())
}

func TestOperandConventions(t *testing.T) {
	for _, tc := range []struct {
		operands string
		valid    bool
	}{
		{"<slug>", true}, {"<id>...", true}, {"<file>|-", true}, {"<file>|x", false}, {"[source-path]", true}, {"[<path>]", true},
		{"[<topic> | <page>[#<section>] | <rule> | <words>...]", true},
		{"SNAPSHOT", false}, {"[SYSTEM]", false}, {"snapshot", false}, {"<Slug>", false},
		{"<slug", false}, {"[<path>", false}, {"<a <b>>", false},
	} {
		if err := checkOperands(tc.operands); (err == nil) != tc.valid {
			t.Errorf("%q: valid %v, got %v", tc.operands, tc.valid, err)
		}
	}
}

var operandName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

var flagValueName = regexp.MustCompile(`^(<[a-z][a-z0-9-]*>|[a-z][a-z0-9-]*(\|[a-z][a-z0-9-]*)+)$`)

func checkOperands(operands string) error {
	if strings.ContainsFunc(operands, unicode.IsUpper) {
		return errors.New("placeholders are lowercase")
	}
	depth := 0
	for i := 0; i < len(operands); i++ {
		switch c := operands[i]; {
		case c == '<':
			end := strings.IndexByte(operands[i:], '>')
			if end < 0 || !operandName.MatchString(operands[i+1:i+end]) {
				return errors.New("write a placeholder as <name>")
			}
			i += end
		case c == '[':
			depth++
		case c == ']':
			if depth--; depth < 0 {
				return errors.New("unbalanced ]")
			}
		case c == '-' && i > 0 && operands[i-1] == '|' && (i+1 == len(operands) || operands[i+1] == ' '):
			// A lone - names standard input, as in <file>|-.
		case c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-':
			if depth == 0 {
				return errors.New("write a required operand as <name>")
			}
		case c == '>':
			return errors.New("unbalanced >")
		case !strings.ContainsRune(" |.#", rune(c)):
			return errors.New("unexpected character " + string(c))
		}
	}
	if depth != 0 {
		return errors.New("unbalanced [")
	}
	return nil
}
