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
	writeFixture(t, filepath.Join(demo, "specs/demo.system.spec.md"), "broken\n")
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
	writeFixture(t, filepath.Join(root, "other/specs/other.system.spec.md"), "---\ntype: system-spec\n---\n# Other\n\n## Purpose\nOther work.\n\n## Scope\nOther only.\n")
	code, out, errOut := invoke("check")
	if code != 1 || errOut != "" || !strings.Contains(out, "root") {
		t.Fatalf("expected ambiguous root finding: %d %s %s", code, out, errOut)
	}
}

func TestGeneratedSpecScopedSystemsAndComponents(t *testing.T) {
	root := projectFixture(t)
	writeFixture(t, filepath.Join(root, "dotspec.yaml"), "name: tool\ngenerated_spec: true\n")
	writeFixture(t, filepath.Join(root, "README.md"), "# Tool\n\n[Read the spec](SPEC.md).\n")
	writeFixture(t, filepath.Join(root, "specs/tool.system.spec.md"), "---\ntype: system-spec\n---\n# Tool\n\n## Purpose\nProcess work.\n\n## Scope\nQueued work.\n\n## Notes\n[Introduction](../README.md).\n")
	writeFixture(t, filepath.Join(root, "worker/specs/worker.system.spec.md"), "---\ntype: system-spec\n---\n# Worker\n\n## Description\nA **local** job executor.\n\n## Purpose\nClaim jobs.\n\n## Scope\nExecutors.\n\n## Part of\n[Tool](../../specs/tool.system.spec.md).\n")
	const grandchild = "worker/queues/specs"
	writeFixture(t, filepath.Join(root, grandchild, "queue.system.spec.md"), `---
type: system-spec
id: queue-boundary
---
# Queue

## Purpose
Retain queued work.

## Scope
Persistence.

## Part of
[Worker](../../specs/worker.system.spec.md).

`)
	writeFixture(t, filepath.Join(root, grandchild, "quality/availability/availability.quality-attribute.spec.md"), `---
type: quality-attribute-spec
---
# Availability
## Definition
Accept work while either A | B is available.
## Limit
Maintenance windows are excluded.
## Open questions
- Which outages count?
- How is uptime measured?
`)
	writeFixture(t, filepath.Join(root, "worker/queues/components/intake/specs/intake.component.spec.md"), `---
type: component-spec
id: incoming-jobs
---

# Intake

## Part of
[Queue](../../../specs/queue.system.spec.md#queue).

## Required interfaces

### Storage

Read persisted jobs.

## Provided interfaces

### Submission

Accept jobs.
`)
	writeFixture(t, filepath.Join(root, "worker/queues/components/journal/specs/journal.component.spec.md"), `---
type: component-spec
---

# Journal

## Part of
[Intake](../../intake/specs/intake.component.spec.md).
`)
	writeFixture(t, filepath.Join(root, grandchild, "requirements/retain.requirement.spec.md"), `---
type: requirement-spec
quality_attributes: [availability]
---
# Retain jobs

## Statement
Each submitted job must remain in the queue until accepted.

## Applies to
**Targets:** [Intake](../../components/intake/specs/intake.component.spec.md) and [Journal](../../components/journal/specs/journal.component.spec.md).

## Other applications
The [Demo system](../../../../apps/demo/specs/demo.system.spec.md) supplies examples.
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
	reader("SPEC.md", "[Read the system specification](<specification/tool.system.spec-report.md>)")
	reader("specification/tool.system.spec-report.md", "[Systems](<systems.md>)")
	reader("specification/systems.md", "- [Worker](<systems/worker/worker.system.spec-report.md>) — A local job executor.")
	reader("specification/systems/worker/worker.system.spec-report.md", "[Tool](<../../tool.system.spec-report.md>)", "[Systems](<systems.md>)")
	reader("specification/systems/worker/systems.md", "- [Queue](<systems/queue/queue.system.spec-report.md>)")
	reader(section+"queue.system.spec-report.md", "[Worker](<../../worker.system.spec-report.md>)", "[All requirements](<requirements.md>)", "[Components](<components.md>)")
	reader(section+"components.md", "- [Intake](<components/intake/intake.component.spec-report.md>)")
	reader(section+"requirements.md", "- [Retain jobs](<requirements/retain-jobs.requirement.spec-report.md>) — Each submitted job must remain in the queue until accepted.", "[Intake](<components/intake/requirements.md>)")
	reader(section+"quality-attributes.md", "**[Availability](<quality/availability/availability.quality-attribute.spec-report.md>)** — Accept work while either A &#124; B is available.")
	reader(section+"quality/availability/availability.quality-attribute.spec-report.md", "# Availability", "## Limit", "Maintenance windows are excluded.", "[Retain jobs](<../../requirements/retain-jobs.requirement.spec-report.md>)")
	reader("specification/requirements-review.report.md",
		"[Which outages count?](<systems/worker/systems/queue/quality/availability/availability.quality-attribute.spec-report.md#open-questions>)",
		"[How is uptime measured?](<systems/worker/systems/queue/quality/availability/availability.quality-attribute.spec-report.md#open-questions>)")
	intake := section + "components/intake/"
	reader(intake+"components/journal/journal.component.spec-report.md", "../../intake.component.spec-report.md", "[Requirements](<requirements.md>)")
	reader(intake+"components/journal/requirements.md", "## Applicable requirements\n\n- [Retain jobs](<../../../../requirements/retain-jobs.requirement.spec-report.md>)")
	reader(intake+"intake.component.spec-report.md", "[Interfaces](<interfaces.md>)", "[Browse required interfaces](<interfaces.md#required-interfaces>)", "[Components](<components.md>)")
	reader(intake+"interfaces.md", "## Required interfaces\n\n- [Storage](<interfaces/storage.interface.spec-report.md>)", "## Provided interfaces\n\n- [Submission](<interfaces/submission.interface.spec-report.md>)")
	reader(intake+"components.md", "- [Journal](<components/journal/journal.component.spec-report.md>)")
	reader(intake+"requirements.md", "## Applicable requirements\n\n- [Retain jobs](<../../requirements/retain-jobs.requirement.spec-report.md>)", "[Journal](<components/journal/requirements.md>)")
	reader(section+"requirements/retain-jobs.requirement.spec-report.md", "../components/intake/intake.component.spec-report.md", "../components/intake/components/journal/journal.component.spec-report.md", "apps/demo/specs/demo.system.spec.md", "../quality/availability/availability.quality-attribute.spec-report.md")
	if _, err := os.Stat(filepath.Join(root, "apps/demo/SPEC.md")); !os.IsNotExist(err) {
		t.Fatalf("referenced project must not acquire generated files: %v", err)
	}
}

func TestGeneratedSystemReadingLayout(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\ngenerated_spec: true\n")
	opening := `---
type: system-spec
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
	writeFixture(t, "specs/notes.system.spec.md", opening)
	build := func() string {
		t.Helper()
		if code, out, errOut := invoke("spec", "build"); code != 0 {
			t.Fatal(code, out, errOut)
		}
		data, err := os.ReadFile("specification/notes.system.spec-report.md")
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	page := build()
	nav := "[Requirements review](<requirements-review.report.md>)"
	if !strings.Contains(page, nav+"\n\n| Section | Explore |") ||
		!strings.Contains(page, "## Purpose\n\nKeep notes.") || strings.Contains(page, "## Goals\n") ||
		!strings.Contains(page, "## Scope\n\nOne notebook.") ||
		strings.Index(page, "## Scope") < strings.Index(page, "Individual authors.") ||
		strings.Count(page, nav) != 1 || strings.Contains(page, "## System topology") {
		t.Fatal(page)
	}

	writeFixture(t, "specs/quality/reliability/reliability.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\n---\n# Reliability\n")
	writeFixture(t, "specs/quality/reliability/durability/durability.quality-attribute.spec.md", strings.ReplaceAll(`---
type: quality-attribute-spec
subattribute_of: [reliability, safety]
---
# Durability
## Definition
The **ID:** below is permanent.

An example:

~~~md
**ID:** example
**Sub-attribute of:** example-parent
~~~
`, "\n", "\r\n"))
	writeFixture(t, "specs/quality/safety/safety.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\n---\n# Safety\n")
	writeFixture(t, "specs/quality/reliability/durability/recovery/recovery.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nsubattribute_of: [durability]\n---\n# Recovery\n")
	writeFixture(t, "specs/notes.system.spec.md", opening+"\n## Quality attributes\n\nProtect recorded work.\n")
	for name, labels := range map[string]string{"Save": "[durability, reliability]", "Restore": "[durability]"} {
		writeFixture(t, filepath.Join("specs/requirements", strings.ToLower(name)+".requirement.spec.md"),
			"---\ntype: requirement-spec\nquality_attributes: "+labels+"\n---\n# "+name+"\n\n## Statement\nKeep the note.\n")
	}
	page = build()
	if !strings.Contains(page, "| [Quality attributes](<#quality-attributes>) | [Attribute index](<quality-attributes.md>)") || !strings.Contains(page, "## Quality attributes\n\nProtect recorded work.") {
		t.Fatalf("unexpected system reading order:\n%s", page)
	}
	durability := "- [Durability](<quality/reliability/durability/durability.quality-attribute.spec-report.md>) — The ID: below is permanent."
	index, err := os.ReadFile("specification/quality-attributes.md")
	if err != nil {
		t.Fatal(err)
	}
	_, overview, _ := strings.Cut(string(index), "[Notes](<notes.system.spec-report.md>)")
	wantOverview := strings.Join([]string{
		"**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)**", "",
		durability,
		"  - [Recovery](<quality/reliability/durability/recovery/recovery.quality-attribute.spec-report.md>)", "",
		"**[Safety](<quality/safety/safety.quality-attribute.spec-report.md>)**", "",
		durability,
	}, "\n")
	if strings.TrimSpace(overview) != wantOverview || !strings.Contains(page, wantOverview) {
		t.Fatalf("quality attributes overview:\n%s\nwant:\n%s", overview, wantOverview)
	}
	attributesPage, err := os.ReadFile("specification/quality/reliability/durability/durability.quality-attribute.spec-report.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"**Parent attribute**\n\n- [Reliability](<../reliability.quality-attribute.spec-report.md>)\n- [Safety](<../../safety/safety.quality-attribute.spec-report.md>)",
		"~~~md\n**ID:** example\n**Sub-attribute of:** example-parent\n~~~",
		"**Classifies**\n\n- [Restore](<../../../requirements/restore.requirement.spec-report.md>)\n- [Save](<../../../requirements/save.requirement.spec-report.md>)",
	} {
		if !strings.Contains(string(attributesPage), want) {
			t.Errorf("detail missing %q:\n%s", want, attributesPage)
		}
	}
}

// The lead loses its authored heading and the map introduces new headings.
// Rebased Markdown in a title and repeated section names must keep distinct
// destinations through both changes.
func TestSystemOverviewRichTitleAndCollidingHeadings(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\ngenerated_spec: true\n")
	writeFixture(t, "specs/guide.md", "# Guide\n")
	writeFixture(t, "specs/scope.system.spec.md", `---
type: system-spec
---
# [Scope](guide.md)

## Description
A local notebook.

## Purpose
Remember ideas.

## Scope
Read [the description](#description).

### Validation
An authored boundary detail.

## System validation plan
See [the boundary detail](#validation).
`)
	for _, args := range [][]string{{"spec", "build"}, {"check"}} {
		if code, out, errOut := invoke(args...); code != 0 {
			t.Fatal(code, out, errOut)
		}
	}
	page, err := os.ReadFile("specification/scope.system.spec-report.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[the description](<#overview.lead>)", "[the boundary detail](<#passage.validation>)", "[Validation](<#overview.validation-1>)"} {
		if !strings.Contains(string(page), want) {
			t.Fatalf("missing %s in %s", want, page)
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

func TestNamedSubjectFilesAgreeWithTheirHeadings(t *testing.T) {
	for _, tc := range []struct{ kind, name, slug, dir string }{
		{"system", "Dispatch & Routing", "dispatch-routing", "external"},
		{"component", "HTTP Task Service", "http-task-service", "adapter"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			root := projectFixture(t)
			if err := os.Mkdir(filepath.Join(root, tc.dir), 0755); err != nil {
				t.Fatal(err)
			}
			args := []string{tc.kind, "init", tc.name, "--path", tc.dir}
			if tc.kind == "component" {
				args = append(args, "--part-of", "specs/tool.system.spec.md")
			}
			if code, _, msg := invoke(args...); code != 0 {
				t.Fatal(msg)
			}
			suffix := tc.kind + ".spec.md"
			expected := tc.slug + "." + suffix
			path := filepath.Join(root, tc.dir, "specs", expected)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "# "+tc.name+"\n") {
				t.Fatal(string(data))
			}
			if code, out, msg := invoke("check", "--path", tc.dir); code != 0 {
				t.Fatal(out, msg)
			}
			for _, bad := range []string{suffix, "different." + suffix} {
				renamed := filepath.Join(filepath.Dir(path), bad)
				if err := os.Rename(path, renamed); err != nil {
					t.Fatal(err)
				}
				if code, out, msg := invoke("check", "--path", tc.dir); code != 1 || !strings.Contains(out+msg, expected) {
					t.Fatalf("expected diagnostic naming %s, got %d %s %s", expected, code, out, msg)
				}
				if err := os.Rename(renamed, path); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
