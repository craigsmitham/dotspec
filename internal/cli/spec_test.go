package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
)

func TestGeneratedSpecificationWorkspaceAndScopes(t *testing.T) {
	root := projectFixture(t)
	demo := filepath.Join(root, "apps/demo")
	writeFixture(t, filepath.Join(demo, "dotspec.yaml"), "name: demo\ngenerated_spec: true\n")
	readReport := func(wantCode int, state string, args ...string) check.Report {
		t.Helper()
		code, out, errOut := invoke(append(args, "--json")...)
		var report check.Report
		if code != wantCode || errOut != "" || json.Unmarshal([]byte(out), &report) != nil || report.GeneratedSpec != state {
			t.Fatalf("unexpected check for %v: code %d, stdout %s, stderr %s", args, code, out, errOut)
		}
		return report
	}

	// The manifest opts in even before there are any generated files.
	readReport(1, "checked", "check", "--project", "demo")
	readReport(0, "", "check", "--project", "tool")
	readReport(0, "out-of-scope", "check", "--project", "demo", "--path", "internal")
	readReport(0, "out-of-scope", "job", "check", "work", "--project", "demo")
	if code, _, errOut := invoke("spec", "build", "--project", "demo", "--path", "internal"); code != 2 || !strings.Contains(errOut, "whole project") {
		t.Fatal(code, errOut)
	}
	if code, _, errOut := invoke("spec", "build", "--project", "tool"); code != 2 || errOut == "" {
		t.Fatal(code, errOut)
	}
	if code, out, errOut := invoke("spec", "build", "--project", "demo"); code != 0 || !strings.Contains(out, "SPEC.md") || errOut != "" {
		t.Fatal(code, out, errOut)
	}
	if code, out, errOut := invoke("spec", "build", "--project", "demo"); code != 0 || !strings.Contains(out, "0 written, 0 removed") || errOut != "" {
		t.Fatal(code, out, errOut)
	}
	readReport(0, "checked", "check", "--project", "demo")
	if err := os.Remove(filepath.Join(demo, "SPEC.md")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := invoke("check", "--all", "--json")
	var workspace workspaceResult
	if code != 1 || json.Unmarshal([]byte(out), &workspace) != nil || len(workspace.Projects) != 2 || workspace.Projects[0].ExitCode != 1 || workspace.Projects[1].ExitCode != 0 || errOut != "" {
		t.Fatal(code, out, errOut)
	}
	writeFixture(t, filepath.Join(demo, "specs/system.spec.md"), "broken\n")
	report := readReport(1, "blocked-by-source", "check", "--project", "demo")
	for _, finding := range report.Diagnostics {
		if strings.HasPrefix(finding.Rule, "spec.generated") {
			t.Fatalf("stale finding after invalid sources: %+v", finding)
		}
	}
	t.Chdir(filepath.Join(demo, "internal"))
	readReport(1, "blocked-by-source", "check")
}

func TestGeneratedSpecificationMultipleRoots(t *testing.T) {
	root := projectFixture(t)
	writeFixture(t, filepath.Join(root, "dotspec.yaml"), "name: tool\ngenerated_spec: true\n")
	writeFixture(t, filepath.Join(root, "other/specs/system.spec.md"), "---\ntype: system\n---\n# Other\n\n## Purpose\nOther work.\n\n## Scope\nOther only.\n")
	code, out, errOut := invoke("check")
	if code != 1 || errOut != "" || !strings.Contains(out, "root") {
		t.Fatalf("expected ambiguous root finding: %d %s %s", code, out, errOut)
	}
}

func TestGeneratedSpecScopedSystemsAndComponents(t *testing.T) {
	root := projectFixture(t)
	writeFixture(t, filepath.Join(root, "dotspec.yaml"), "name: tool\ngenerated_spec: true\n")
	writeFixture(t, filepath.Join(root, "README.md"), "# Tool\n\n[Read the spec](SPEC.md).\n")
	writeFixture(t, filepath.Join(root, "specs/system.spec.md"), "---\ntype: system\n---\n# Tool\n\n## Purpose\nProcess work.\n\n## Scope\nQueued work.\n\n## Notes\n[Introduction](../README.md).\n")
	writeFixture(t, filepath.Join(root, "worker/specs/system.spec.md"), "---\ntype: system\n---\n# Worker\n\n## Purpose\nClaim jobs.\n\n## Scope\nExecutors.\n\n## Part of\n[Tool](../../specs/system.spec.md).\n")
	const grandchild = "worker/queues/specs"
	writeFixture(t, filepath.Join(root, grandchild, "system.spec.md"), `---
type: system
---
# Queue

## Purpose
Retain queued work.

## Scope
Persistence.

## Part of
[Worker](../../specs/system.spec.md).

## Components

### Intake

#### Part of
[Queue](#queue).

### Journal

#### Part of
[Intake](#intake).
`)
	writeFixture(t, filepath.Join(root, grandchild, "requirements/retain.requirement.spec.md"), `---
type: requirement
---
# Retain jobs

## Statement
Each submitted job must remain in the queue until accepted.

## Applies to
**Targets:** [Intake](../system.spec.md#intake) and [Journal](../system.spec.md#journal).

## Other applications
The [Demo system](../../../../apps/demo/specs/system.spec.md) supplies examples.
`)
	if code, out, errOut := invoke("spec", "build"); code != 0 {
		t.Fatal(code, out, errOut)
	}
	if code, out, errOut := invoke("check"); code != 0 {
		t.Fatal(code, out, errOut)
	}
	section := "specification/systems/worker/systems/queue/"
	reader := func(path string, contains ...string) {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		for _, expected := range contains {
			if !strings.Contains(string(data), expected) {
				t.Errorf("%s missing %q:\n%s", path, expected, data)
			}
		}
	}
	reader("specification/systems/index.md", "worker/SPEC.md")
	reader("specification/systems/worker/systems/index.md", "queue/SPEC.md")
	reader("specification/systems/worker/SPEC.md", "[Tool](<../../../SPEC.md>) · [Full contents]")
	reader(section+"SPEC.md", "[Worker](<../../SPEC.md>) · [Requirements](<requirements/index.md>) · [Full contents]")
	reader(section+"components/journal.md", "intake.md", "../requirements/retain-jobs.md")
	reader(section+"components/intake.md", "journal.md", "../requirements/retain-jobs.md")
	reader(section+"requirements/retain-jobs.md", "../components/intake.md", "../components/journal.md", "apps/demo/specs/system.spec.md")
	reader("specification/index.md", "systems/worker/systems/queue/requirements/index.md")
	if _, err := os.Stat(filepath.Join(root, "apps/demo/SPEC.md")); !os.IsNotExist(err) {
		t.Fatalf("referenced project must not acquire generated files: %v", err)
	}
}

func TestGeneratedSystemReadingLayout(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\ngenerated_spec: true\n")
	opening := `---
type: system
---
# Notes

## Scope
One notebook.

### Storage
Local files.

## Purpose
Keep notes.

### Audience
Individual authors.
`
	writeFixture(t, "specs/system.spec.md", opening)
	build := func() string {
		t.Helper()
		if code, out, errOut := invoke("spec", "build"); code != 0 {
			t.Fatal(code, out, errOut)
		}
		data, err := os.ReadFile("SPEC.md")
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	page := build()
	nav := "[Full contents](<specification/index.md>) · [Requirements review](<specification/review.md>)"
	if !strings.Contains(page, "Local files.\n\n---\n\n"+nav+"\n\n## Explore this system") || strings.Count(page, nav) != 1 {
		t.Fatal(page)
	}

	attributes := `
## Quality attributes
### Reliability
**ID:** reliability

### Durability
The **ID:** below is permanent.

**Sub-attribute of:** reliability
**ID:** durability

An example:

~~~md
**ID:** example
**Sub-attribute of:** example-parent
~~~
`
	writeFixture(t, "specs/system.spec.md", strings.ReplaceAll(opening+attributes, "\n", "\r\n"))
	for name, labels := range map[string]string{"Save": "[durability, reliability]", "Restore": "[durability]"} {
		writeFixture(t, filepath.Join("specs/requirements", strings.ToLower(name)+".requirement.spec.md"),
			"---\ntype: requirement\nquality_attributes: "+labels+"\n---\n# "+name+"\n\n## Statement\nKeep the note.\n")
	}
	page = build()
	for _, want := range []string{
		"Local files.\n\n---\n\n[Requirements](<specification/requirements/index.md>) · " + nav + "\n\n## Quality attributes",
		"**Sub-attribute of:** reliability\\\n**ID:** durability",
		"~~~md\n**ID:** example\n**Sub-attribute of:** example-parent\n~~~",
		"**Classifies** · [Save](<specification/requirements/save.md>)",
		"**Classifies**\n\n- [Restore](<specification/requirements/restore.md>)\n- [Save](<specification/requirements/save.md>)",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("missing %q in:\n%s", want, page)
		}
	}
}

func TestGeneratedSpecificationMissingSystem(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	writeFixture(t, "dotspec.yaml", "name: empty\ngenerated_spec: true\n")
	code, out, errOut := invoke("check")
	if code != 1 || errOut != "" || !strings.Contains(out, "root") {
		t.Fatalf("expected missing root finding: %d %s %s", code, out, errOut)
	}
}
