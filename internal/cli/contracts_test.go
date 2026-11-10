package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"github.com/spf13/cobra"
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
	if command == "dotspec" || strings.HasSuffix(command, " init") {
		return
	}
	if _, err := system.Init("."); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, "specs/jobs/develop-requirements.job.spec.md", viewFixture)
	writeFixture(t, "specs/use-cases/check-specifications.use-case.spec.md", usecase.Template)
}

func TestDeclaredCommandExamples(t *testing.T) {
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		if cmd.Example != "" {
			t.Run(cmd.CommandPath(), func(t *testing.T) {
				prepareExample(t, cmd.CommandPath())
				for _, example := range strings.Split(cmd.Example, "\n") {
					args := strings.Fields(example)
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

func TestReadmeQuickStart(t *testing.T) {
	data, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatal(err)
	}
	_, run, found := strings.Cut(string(data), "## Run\n")
	if !found {
		t.Fatal("missing Run section")
	}
	_, run, found = strings.Cut(run, "```sh\n")
	if !found {
		t.Fatal("missing quick start")
	}
	commands, _, found := strings.Cut(run, "```")
	if !found {
		t.Fatal("unclosed quick start")
	}
	prepareExample(t, "readme")
	for _, line := range strings.Split(strings.TrimSpace(commands), "\n") {
		args := strings.Fields(line)
		if len(args) < 4 || strings.Join(args[:3], " ") != "go run ." {
			t.Fatalf("unexpected quick start: %s", line)
		}
		if code, out, errOut := invoke(args[3:]...); code != 0 || errOut != "" {
			t.Fatalf("%s: %d %s %s", line, code, out, errOut)
		}
	}
}

func TestStructuredReadContracts(t *testing.T) {
	prepareExample(t, "read")
	for _, args := range [][]string{
		{"system", "show"}, {"system", "list"}, {"job", "show", "develop-requirements"},
		{"job", "list"}, {"job", "step", "list", "develop-requirements"},
		{"use-case", "show", "check-specifications"}, {"use-case", "list"},
		{"job", "outcome", "list", "develop-requirements"}, {"job", "map", "develop-requirements", "--outcomes"},
	} {
		code, out, errOut := invoke(append(args, "--json")...)
		if code != 0 || errOut != "" {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
		var result map[string]any
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatal(err)
		}
		if result["schema_version"] != float64(1) {
			t.Fatalf("missing version: %s", out)
		}
		if strings.Contains(out, "\x1b") {
			t.Fatalf("styled JSON: %s", out)
		}
	}
	_, output, _ := invoke("job", "show", "develop-requirements", "--json")
	var result jobResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatal(err)
	}
	parsed, err := job.Parse([]byte(viewFixture))
	if err != nil {
		t.Fatal(err)
	}
	if result.Job.Markdown != parsed.Markdown || len(result.Job.Steps) != 2 || !filepath.IsAbs(result.Job.Path) || result.Job.System == nil {
		t.Fatalf("lost authored content or association: %+v", result.Job)
	}
	_, output, _ = invoke("job", "outcome", "list", "develop-requirements", "--step", "gather", "--json")
	var outcomes outcomesResult
	if err := json.Unmarshal([]byte(output), &outcomes); err != nil {
		t.Fatal(err)
	}
	if len(outcomes.Overall) != 0 || len(outcomes.Steps) != 1 || outcomes.Steps[0].ID != "gather" || outcomes.Steps[0].Outcomes[0].ID != "OUT-002" {
		t.Fatalf("wrong selected scope: %+v", outcomes)
	}
	_, output, _ = invoke("job", "map", "develop-requirements", "--json")
	if err := json.Unmarshal([]byte(output), &outcomes); err != nil {
		t.Fatal(err)
	}
	if len(outcomes.Overall) != 0 || len(outcomes.Steps[0].Outcomes) != 0 {
		t.Fatalf("unrequested outcomes: %s", output)
	}
}

func TestOutputFailuresAcrossCommands(t *testing.T) {
	prepareExample(t, "output")
	for _, args := range [][]string{
		{"--help"}, {"--version"}, {"completion", "bash"}, {"man"},
		{"system", "show"}, {"system", "list", "--json"}, {"job", "show", "develop-requirements"},
		{"job", "check", "develop-requirements", "--json"},
		{"use-case", "show", "check-specifications"}, {"use-case", "list", "--json"},
		{"use-case", "check", "check-specifications", "--json"},
	} {
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
		{"system", "show", "--path", "specs"}, {"system", "list", "--path", "specs"},
		{"job", "init", "new", "--path", "specs"}, {"job", "show", "absent"},
	} {
		if code, out, errOut := invoke(args...); code != 2 || out != "" || errOut == "" {
			t.Fatalf("invocation/operation %v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{
		{"check", "--path", "specs", "--json"}, {"job", "check", "absent", "--json"},
	} {
		code, out, errOut := invoke(args...)
		if code != 2 || !strings.Contains(out, `"complete": false`) || errOut != "" {
			t.Fatalf("incomplete %v: %d %q %q", args, code, out, errOut)
		}
	}
	writeFixture(t, "specs/jobs/broken.job.spec.md", "broken")
	for _, args := range [][]string{{"job", "show", "broken"}, {"job", "list"}} {
		if code, out, errOut := invoke(args...); code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
			t.Fatalf("specification %v: %d %q %q", args, code, out, errOut)
		}
	}
}
