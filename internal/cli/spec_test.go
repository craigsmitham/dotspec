package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/model"
)

func TestGeneratedSpecificationWorkspaceAndScopes(t *testing.T) {
	root := projectFixture(t)
	demo := filepath.Join(root, "apps/demo")
	writeFixture(t, filepath.Join(root, "dotspec.yaml"), "name: tool\nreport:\n  enabled: false\n")
	readReport := func(wantCode int, state string, args ...string) check.Report {
		t.Helper()
		code, out, errOut := invoke(append(args, "--json")...)
		var report check.Report
		if code != wantCode || errOut != "" || json.Unmarshal([]byte(out), &report) != nil || report.GeneratedSpec != state {
			t.Fatalf("unexpected check for %v: code %d, stdout %s, stderr %s", args, code, out, errOut)
		}
		return report
	}

	// A name-only manifest requires reports from the first whole-project check.
	readReport(1, "checked", "check", "--project", "demo")
	readReport(0, "", "check", "--project", "tool")
	readReport(0, "out-of-scope", "check", "--project", "demo", "--path", "internal")
	readReport(0, "out-of-scope", "job", "check", "work", "--project", "demo")
	if code, _, errOut := invoke("spec", "build", "--project", "demo", "--path", "internal"); code != 2 || !strings.Contains(errOut, "whole project") {
		t.Fatal(code, errOut)
	}
	if code, _, errOut := invoke("spec", "build", "--project", "tool"); code != 2 || !strings.Contains(errOut, "report.enabled is false") {
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
	writeFixture(t, filepath.Join(root, "specs/tool.system.spec.md"), "broken\n")
	readReport(1, "", "check", "--project", "tool")
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
	writeFixture(t, filepath.Join(root, "other/specs/other.system.spec.md"), "---\ntype: system-spec\n---\n# Other\n\n## Purpose\nOther work.\n\n## Scope\nOther only.\n")
	if code, out, errs := invoke("spec", "build"); code != 0 {
		t.Fatal(out, errs)
	}
	if code, out, errs := invoke("check"); code != 0 {
		t.Fatal(out, errs)
	}
	for _, name := range []string{"tool", "other"} {
		if code, out, errs := invoke("model", "product", name); code != 0 || !strings.Contains(out, name+".system.spec.md") {
			t.Fatal(out, errs)
		}
	}

}

func TestGeneratedSpecScopedSystemsAndComponents(t *testing.T) {
	root := projectFixture(t)
	writeFixture(t, filepath.Join(root, "README.md"), "# Tool\n\n[Read the spec](SPEC.md).\n")
	writeFixture(t, filepath.Join(root, "specs/tool.system.spec.md"), "---\ntype: system-spec\n---\n# Tool\n\n## Purpose\nProcess work.\n\n## Scope\nQueued work.\n\n## Evidence\n[Introduction](../README.md).\n")
	writeFixture(t, filepath.Join(root, "worker/specs/worker.system.spec.md"), "---\ntype: system-spec\n---\n# Worker\n\n## Description\nA **local** job executor.\n\n## Framing\nCould workers recover interrupted jobs?\n\n## Purpose\nClaim jobs.\n\n## Scope\nExecutors.\n\n## Part of\n[Tool](../../specs/tool.system.spec.md).\n")
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

## Interacts with
The [Tool](../../../specs/tool.system.spec.md) schedules retries.

`)
	writeFixture(t, filepath.Join(root, grandchild, "quality/availability/availability.quality-attribute.spec.md"), `---
type: quality-attribute-spec
---
# Availability
## Definition
Accept work while either A | B is available.
### Limit
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

## Responsibilities

Receive submissions through the [storage contracts](#required-interfaces).

## Architecture

Journal keeps the durable record.

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

### Conditions
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
	reader("specification/tool.system.spec-report.md", "### Subsystems", "Subsystem: Worker", "[Open questions](<open-questions.md>)", "[Systems](<systems.md>)")
	reader("specification/systems/worker/worker.system.spec-report.md", "## Framing", "Could workers recover interrupted jobs?")
	reader("specification/systems.md", "- [Worker](<systems/worker/worker.system.spec-report.md>) — A local job executor.")
	reader("specification/systems/worker/worker.system.spec-report.md", "[Tool](<../../tool.system.spec-report.md>)", "[Systems](<systems.md>)")
	reader("specification/systems/worker/systems.md", "- [Queue](<systems/queue/queue.system.spec-report.md>)")
	reader(section+"queue.system.spec-report.md", "### System structure", "### System context", "### Interactions", "owner ---|interacts with|", "schedules retries.", "[Worker](<../../worker.system.spec-report.md>)", "[Requirements](<requirements.md>)", "[Components](<components.md>)")
	reader(section+"components.md", "- [Intake](<components/intake/intake.component.spec-report.md>)")
	reader(section+"requirements.md", "- [Retain jobs](<requirements/retain-jobs.requirement.spec-report.md>) — Each submitted job must remain in the queue until accepted.", "[Intake requirements](<components/intake/requirements.md>)")
	reader(section+"quality-attributes.md", "**[Availability](<quality/availability/availability.quality-attribute.spec-report.md>)** — Accept work while either A &#124; B is available.")
	reader(section+"quality/availability/availability.quality-attribute.spec-report.md", "# Availability", "### Limit", "Maintenance windows are excluded.", "[Retain jobs](<../../requirements/retain-jobs.requirement.spec-report.md>)")
	reader("specification/open-questions.md",
		"[Which outages count?](<systems/worker/systems/queue/quality/availability/availability.quality-attribute.spec-report.md#open-questions>) — Availability",
		"[How is uptime measured?](<systems/worker/systems/queue/quality/availability/availability.quality-attribute.spec-report.md#open-questions>) — Availability")
	intake := section + "components/intake/"
	reader(intake+"components/journal/journal.component.spec-report.md", "../../intake.component.spec-report.md", "[Requirements](<requirements.md>)")
	reader(intake+"components/journal/requirements.md", "## Applicable requirements\n\n- [Retain jobs](<../../../../requirements/retain-jobs.requirement.spec-report.md>)")
	reader(intake+"intake.component.spec-report.md", "## Architecture\n\nJournal keeps the durable record.", "- **Provided** · [Submission]", "- **Required** · [Storage]", "[Interfaces](<interfaces.md>)", "[Required interfaces](<interfaces.md#required-interfaces>)", "[storage contracts](<interfaces.md#required-interfaces>)", "[Components](<components.md>)")
	componentPage, err := os.ReadFile(filepath.Join(root, intake, "intake.component.spec-report.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(componentPage), "## Required interfaces") || strings.Contains(string(componentPage), "## Provided interfaces") {
		t.Fatal("collection-only interface sections in component report", string(componentPage))
	}
	reader(intake+"interfaces.md", "## Required interfaces\n\n- [Storage](<interfaces/storage.interface.spec-report.md>)", "## Provided interfaces\n\n- [Submission](<interfaces/submission.interface.spec-report.md>)")
	reader(intake+"components.md", "- [Journal](<components/journal/journal.component.spec-report.md>)")
	reader(intake+"requirements.md", "## Applicable requirements\n\n- [Retain jobs](<../../requirements/retain-jobs.requirement.spec-report.md>)", "[Journal requirements](<components/journal/requirements.md>)")
	reader(section+"requirements/retain-jobs.requirement.spec-report.md", "../components/intake/intake.component.spec-report.md", "../components/intake/components/journal/journal.component.spec-report.md", "apps/demo/specs/demo.system.spec.md", "../quality/availability/availability.quality-attribute.spec-report.md")
	if _, err := os.Stat(filepath.Join(root, "apps/demo/SPEC.md")); !os.IsNotExist(err) {
		t.Fatalf("referenced project must not acquire generated files: %v", err)
	}
}

func TestGeneratedOpenQuestionsLifecycle(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
	const system = "---\ntype: system-spec\n---\n# Notes\n\n## Purpose\nKeep notes.\n\n## Scope\nOne notebook.\n"
	const source = "specs/notes.system.spec.md"
	const index = "specification/open-questions.md"
	const overview = "specification/notes.system.spec-report.md"
	const retired = "specification/requirements-review.report.md"
	writeFixture(t, retired, "---\ngenerated_by: dotspec\ntype: requirements-review-report\nsource: null\n---\n# Requirements review\n")
	build := func() {
		t.Helper()
		for _, args := range [][]string{{"spec", "build"}, {"check"}} {
			if code, out, errOut := invoke(args...); code != 0 {
				t.Fatal(code, out, errOut)
			}
		}
	}
	read := func(file string) string {
		t.Helper()
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	writeFixture(t, source, system+"\n## Open questions\n\n### Who can recover a note?\nShould `_owner_` be required for [recovery][policy]?\n\n[policy]: #scope\n")
	build()
	if page := read(index); !strings.Contains(page, "[Who can recover a note? Should &#95;owner&#95; be required for recovery?](<notes.system.spec-report.md#who-can-recover-a-note>) — Notes") ||
		strings.Contains(page, "Inspection findings") || !strings.Contains(read(overview), "[Open questions](<open-questions.md>)") {
		t.Fatal(page)
	}
	if _, err := os.Stat(retired); !os.IsNotExist(err) {
		t.Fatalf("obsolete review report after build: %v", err)
	}

	writeFixture(t, source, system+"\n## Open questions\n")
	build()
	if _, err := os.Stat(index); !os.IsNotExist(err) || strings.Contains(read(overview), "open-questions.md") {
		t.Fatalf("question index or link remains after resolving the last question: %v", err)
	}
}

func TestOverviewKeepsNavigationOnlyCollectionsInTheMap(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
	const system = "---\ntype: system-spec\n---\n# Notes\n\n## Purpose\nKeep notes.\n\n## Scope\nOne notebook.\n\n"
	for _, name := range []string{"Draft", "Review", "Publish", "Archive"} {
		slug := strings.ToLower(name)
		writeFixture(t, "specs/business-rules/"+slug+".business-rule.spec.md",
			"---\ntype: business-rule-spec\n---\n# "+name+" rule\n\n## Statement\nKeep the "+slug+" record.\n\n## Applies to\n**Targets:** [Notes](../notes.system.spec.md).\n")
		writeFixture(t, "specs/features/"+slug+".feature.spec.md", "---\ntype: feature-spec\n---\n# "+name+" notes\n")
	}
	for _, narrative := range []string{"", "The archive keeps all published notes."} {
		writeFixture(t, "specs/notes.system.spec.md", system+narrative+"\n")
		if code, out, errOut := invoke("spec", "build"); code != 0 {
			t.Fatal(code, out, errOut)
		}
		data, err := os.ReadFile("specification/notes.system.spec-report.md")
		if err != nil {
			t.Fatal(err)
		}
		page := string(data)
		for _, route := range []string{"[Business rules](<business-rules.md>)", "[Features](<features.md>)"} {
			if strings.Count(page, route) != 1 {
				t.Fatalf("expected one navigation route %s in %s", route, page)
			}
		}
		if strings.Contains(page, "## Requirements\n") || strings.Contains(page, "## Features\n") ||
			strings.Contains(page, "## Business rules\n") || !strings.Contains(page, narrative) {
			t.Fatal(page)
		}
	}
}

func TestGeneratedSystemReadingLayout(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
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
	writeFixture(t, "specs/quality-models/shared.quality-model.md", "---\ntype: quality-model\n---\n# Shared model\n\n## Quality attributes\n\n### Precision\n\n**ID:** precision\n\nCloseness of a measurement.\n")
	page := build()
	if strings.Contains(page, "## Quality Attributes") || strings.Contains(page, "[Precision]") ||
		!strings.HasSuffix(strings.TrimSpace(page), " · [Quality models](<quality-models.md>)") ||
		strings.Count(page, "(<quality-models.md>)") != 1 {
		t.Fatalf("model-only system overview:\n%s", page)
	}
	if !strings.Contains(page, "<summary>On this page</summary>") ||
		!strings.Contains(page, "## Purpose\n\nKeep notes.") || strings.Contains(page, "## Goals\n") ||
		!strings.Contains(page, "## Scope\n\nOne notebook.") ||
		strings.Index(page, "## Scope") < strings.Index(page, "Individual authors.") ||
		strings.Contains(page, "open-questions.md") || strings.Contains(page, "## Architecture") {
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
	writeFixture(t, "plan.md", "# Release plan\n")
	writeFixture(t, "specs/notes.system.spec.md", opening+"\n## Project references\n\nThe [release plan](../plan.md) names the reviewer.\n\n## Quality attributes\n")
	for name, labels := range map[string]string{"Save": "[durability, reliability]", "Restore": "[durability]"} {
		writeFixture(t, filepath.Join("specs/requirements", strings.ToLower(name)+".requirement.spec.md"),
			"---\ntype: requirement-spec\nquality_attributes: "+labels+"\n---\n# "+name+"\n\n## Statement\nKeep the note.\n")
	}
	page = build()
	if !strings.Contains(page, "[Quality Attributes](<#quality-attributes>)") || strings.Contains(page, "(<quality-attributes.md>)") || strings.Count(page, "## Quality Attributes\n") != 1 {
		t.Fatalf("unexpected system reading order:\n%s", page)
	}
	if !strings.Contains(page, "## Project references\n\nThe [release plan](<../plan.md>) names the reviewer.") || strings.Index(page, "## Project references") < strings.Index(page, "## Quality Attributes") {
		t.Fatal(page)
	}
	durability := "- [Durability](<quality/reliability/durability/durability.quality-attribute.spec-report.md>) — The ID: below is permanent."
	index, err := os.ReadFile("specification/quality-attributes.md")
	if err != nil {
		t.Fatal(err)
	}
	_, overview, _ := strings.Cut(string(index), "# Quality attributes\n")
	wantOverview := strings.Join([]string{
		"**[Reliability](<quality/reliability/reliability.quality-attribute.spec-report.md>)**", "",
		durability,
		"  - [Recovery](<quality/reliability/durability/recovery/recovery.quality-attribute.spec-report.md>)", "",
		"**[Safety](<quality/safety/safety.quality-attribute.spec-report.md>)**", "",
		durability,
	}, "\n")
	_, hierarchy, _ := strings.Cut(page, "## Quality Attributes\n")
	hierarchy, _, _ = strings.Cut(hierarchy, "\n## ")
	if strings.TrimSpace(hierarchy) != wantOverview || !strings.HasPrefix(strings.TrimSpace(overview), wantOverview) {
		t.Fatalf("system hierarchy:\n%s\nindex:\n%s\nwant:\n%s", hierarchy, overview, wantOverview)
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

func TestSystemContextSectionsStandAlone(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
	const system = "---\ntype: system-spec\n---\n# Notes\n\n## Purpose\nKeep drafts.\n"
	const source = "specs/notes.system.spec.md"
	const report = "specification/notes.system.spec-report.md"
	const promise = "specs/requirements/read.requirement.spec.md"
	writeFixture(t, promise, "---\ntype: requirement-spec\n---\n# Read notes\n\n## Statement\nSee [winter](../notes.system.spec.md#winter).\n")
	writeFixture(t, source, system+`
## Constraints

## Scope

## Assumptions

## Environment

### Winter

- Readers wear gloves.

## Out of scope
`)
	build := func() string {
		t.Helper()
		for _, args := range [][]string{{"spec", "build"}, {"check"}} {
			if code, out, errOut := invoke(args...); code != 0 {
				t.Fatal(code, out, errOut)
			}
		}
		data, err := os.ReadFile(report)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	page := build()
	for _, title := range []string{"Scope", "Out of scope", "Assumptions", "Constraints"} {
		if strings.Contains(page, "## "+title+"\n") || strings.Contains(page, "["+title+"](") {
			t.Fatalf("empty %s visible in %s", title, page)
		}
	}
	if !strings.Contains(page, "[Environment](<#environment>)") || !strings.Contains(page, "## Environment\n\n### Winter\n\n- Readers wear gloves.") {
		t.Fatal(page)
	}
	data, err := os.ReadFile("specification/requirements/read-notes.requirement.spec-report.md")
	if err != nil || !strings.Contains(string(data), "[winter](<../notes.system.spec-report.md#winter>)") {
		t.Fatalf("nested section citation: %s, %v", data, err)
	}

	writeFixture(t, promise, "---\ntype: requirement-spec\n---\n# Read notes\n\n## Statement\nRead saved notes.\n")
	writeFixture(t, source, system+"\n## Scope\n\n## Constraints\nThe archive agreement requires local storage.\n")
	page = build()
	if strings.Contains(page, "## Scope\n") || !strings.Contains(page, "[Constraints](<#constraints>)") ||
		!strings.Contains(page, "## Constraints\n\nThe archive agreement requires local storage.") {
		t.Fatal(page)
	}

	writeFixture(t, source, system+"\n## Out of scope\nSharing notes.\n\n### Exports\nPrinting is excluded.\n")
	writeFixture(t, promise, "---\ntype: requirement-spec\n---\n# Read notes\n\n## Statement\nSee [exclusions](../notes.system.spec.md#out-of-scope).\n")
	page = build()
	if !strings.Contains(page, "## Scope\n\n### Out of scope\n\nSharing notes.\n\n#### Exports\n\nPrinting is excluded.") ||
		strings.Count(page, "[Scope](<#scope>)") != 1 || strings.Contains(page, "[Out of scope](<#") {
		t.Fatal(page)
	}
	data, err = os.ReadFile("specification/requirements/read-notes.requirement.spec-report.md")
	if err != nil || !strings.Contains(string(data), "[exclusions](<../notes.system.spec-report.md#out-of-scope>)") {
		t.Fatalf("exclusions citation: %s, %v", data, err)
	}
}

func TestNamedOverviewEntriesKeepScopedDetailAndLinks(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
	writeFixture(t, "specs/notes.system.spec.md", `---
type: system-spec
---
# Notes

## Purpose
Keep drafts.

## Scope
Local edits.

## Assumptions
### The author is present

- The author is present.
- A second authorized person can attend when the author asks for help identifying the entries that may be archived after review.

### Local storage
A draft stays on disk.

#### Basis
The [store constraint](../store/specs/store.component.spec.md#storage) qualifies this assumption.
`)
	source := "store/specs/store.component.spec.md"
	component := `---
type: component-spec
---
# Store

## Part of
[Notes](../../specs/notes.system.spec.md).

## Responsibilities
Apply the [retention constraint](#storage) and its [exception](#exception).

## Evidence
Storage layout.[^disk]

[^disk]: A file on the author's device.
`
	writeFixture(t, source, component+`
## Boundary
Writes occur locally.

### Storage
Retain edits.

#### Exception
A reset clears the draft.[^disk] See [the basis](../../specs/notes.system.spec.md#basis).

### Storage
A separate archive retains old drafts.
`)
	for _, args := range [][]string{{"spec", "build"}, {"check"}} {
		if code, out, errOut := invoke(args...); code != 0 {
			t.Fatal(code, out, errOut)
		}
	}
	read := func(file string) string {
		t.Helper()
		data, err := os.ReadFile("specification/" + file)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	premises := read("notes.system.spec-report.md")
	for _, want := range []string{
		"## Assumptions\n\n- [The author is present](<assumptions.md#the-author-is-present>) — The author is present. …",
		"- [Local storage](<assumptions.md#local-storage>) — A draft stays on disk.",
	} {
		if !strings.Contains(premises, want) {
			t.Fatalf("missing assumption preview %q in %s", want, premises)
		}
	}
	overview := read("components/store/store.component.spec-report.md")
	for _, want := range []string{"## Scope\n\n### Boundary\n\nWrites occur locally.", "#### Storage\n\nRetain edits.", "A reset clears the draft.", "A separate archive retains old drafts.", "[the basis](<../../assumptions.md#basis>)", "[^disk]: A file"} {
		if !strings.Contains(overview, want) {
			t.Fatalf("missing %q in %s", want, overview)
		}
	}
	if !strings.Contains(read("assumptions.md"), "[store constraint](<components/store/store.component.spec-report.md#storage>)") {
		t.Fatal(read("assumptions.md"))
	}
	// A rebuild retires both the component collection and its system navigation.
	writeFixture(t, source, component[:strings.Index(component, "## Responsibilities")])
	writeFixture(t, "specs/notes.system.spec.md", "---\ntype: system-spec\n---\n# Notes\n\n## Purpose\nKeep drafts.\n\n## Scope\nLocal edits.\n")
	if code, out, errOut := invoke("spec", "build"); code != 0 {
		t.Fatal(code, out, errOut)
	}
	for _, file := range []string{"constraints.md", "components/store/constraints.md", "assumptions.md"} {
		if _, err := os.Stat("specification/" + file); !os.IsNotExist(err) {
			t.Fatalf("%s still present: %v", file, err)
		}
	}
}

// The lead loses its authored heading and the map introduces new headings.
// Rebased Markdown in a title and repeated section names must keep distinct
// destinations through both changes.
func TestSystemOverviewRichTitleAndCollidingHeadings(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
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

### Open questions
An authored boundary detail.

## Open questions
See [the boundary detail](#open-questions).
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
	for _, want := range []string{"[the description](<#overview.lead>)", "[the boundary detail](<#passage.open-questions>)", "[Open questions](<#overview.open-questions-1>)"} {
		if !strings.Contains(string(page), want) {
			t.Fatalf("missing %s in %s", want, page)
		}
	}
}

func TestGeneratedSpecificationMissingSystem(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	writeFixture(t, "dotspec.yaml", "name: empty\n")
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

func TestReportPreviewsAsCollectionsGrow(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFixture(t, "dotspec.yaml", "name: notes\n")
	const systemFile = "specs/notes.system.spec.md"
	const system = "---\ntype: system-spec\n---\n# Notes\n\n## Purpose\nRecord notes.\n\n## Scope\nLocal work.\n"
	const goals = "\n## Goals\n\n### One\nFirst goal.\n\n### Two\nSecond goal.\n\n### Three\nThird goal.\n"
	const fourthGoal = "\n### Four\nFourth goal.\n"
	const uses = "\n## Concept uses\n\n### My notebook\n\n#### Description\nThe current author's notes.\n\n#### Definition\n[Notebook](concepts/notebook.concept.spec.md).\n"
	writeFixture(t, systemFile, system+goals+"\n## Concept uses\n\n## Synchronizations\n")
	writeFixture(t, "specs/concepts/notebook.concept.spec.md", "---\ntype: concept-spec\n---\n# Notebook\n\n## Purpose\nKeep written ideas.\n")
	addUseCase := func(name, extra string) {
		writeFixture(t, "specs/use-cases/"+strings.ToLower(name)+".use-case.spec.md",
			"---\ntype: use-case-spec\n---\n# "+name+"\n\n## Goal\n"+name+" a note.\n"+extra)
	}
	const sharedFile = "specs/requirements/keep-notes.requirement.spec.md"
	const sharedSource = "---\ntype: requirement-spec\n---\n# Keep notes\n\n## Statement\nSaved notes remain available.\n"
	writeFixture(t, sharedFile, sharedSource)
	const illustrates = "\n## Main success scenario\n1. The author opens a note.\n   - **Illustrates:** [Keep notes](../requirements/keep-notes.requirement.spec.md).\n"
	addUseCase("Archive", illustrates)
	addUseCase("Edit", illustrates)
	addUseCase("Review", "")
	addRequirement := func(name string) {
		writeFixture(t, "specs/requirements/"+strings.ToLower(name)+".requirement.spec.md",
			"---\ntype: requirement-spec\n---\n# "+name+"\n\n## Statement\nKeep a "+name+" copy.\n\n## Refines\n[Keep notes](keep-notes.requirement.spec.md).\n")
	}
	addRequirement("Local")
	const featureFile = "specs/features/editing.feature.spec.md"
	const feature = "---\ntype: feature-spec\n---\n# Editing\n\n## Groups\n[Edit](../use-cases/edit.use-case.spec.md).\n[Keep notes](../requirements/keep-notes.requirement.spec.md).\n"
	writeFixture(t, featureFile, feature)
	read := func(file string) string {
		t.Helper()
		b, err := os.ReadFile("specification/" + file)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	build := func() string {
		t.Helper()
		for _, args := range [][]string{{"spec", "build"}, {"check"}} {
			if code, out, msg := invoke(args...); code != 0 {
				t.Fatal(args, out, msg)
			}
		}
		return read("notes.system.spec-report.md")
	}
	section := func(page, heading string) string {
		t.Helper()
		_, body, found := strings.Cut(page, "\n## "+heading+"\n")
		if !found {
			t.Fatalf("missing %s section:\n%s", heading, page)
		}
		body, _, _ = strings.Cut(body, "\n## ")
		body, _, _ = strings.Cut(body, "\n---\n\n[View source](")
		return strings.TrimSpace(body)
	}
	before := build()
	const featureReport = "features/editing.feature.spec-report.md"
	const editGoal = "- [Edit](<../use-cases/edit.use-case.spec-report.md#goal>) — Edit a note."
	if !strings.Contains(read(featureReport), editGoal) {
		t.Fatal(read(featureReport))
	}
	for _, heading := range []string{"Concepts", "Use cases", "Requirements"} {
		if strings.Contains(before, "\n## "+heading+"\n") {
			t.Fatal(before)
		}
	}
	if !strings.Contains(read("concepts.md"), "\n# Concept definitions\n") || !strings.Contains(read("use-cases.md"), "Review a note.") {
		t.Fatal("missing index entries")
	}
	refines := "\n  - Refines [Keep notes](<requirements/keep-notes.requirement.spec-report.md>)"
	if strings.Count(read("requirements.md"), refines) != 1 {
		t.Fatal(read("requirements.md"))
	}
	writeFixture(t, sharedFile, strings.Replace(sharedSource, "Saved notes remain available.", "Reopened notes retain their text.", 1))

	addUseCase("Share", "")
	writeFixture(t, featureFile, feature+"[Share](../use-cases/share.use-case.spec.md).\n")
	addRequirement("Remote")
	writeFixture(t, systemFile, system+goals+fourthGoal+uses+"\n## Synchronizations\n\n### Archive updates\n\n#### When\nThe note is saved.\n\n#### Then\nArchive the current version.\n")
	after := build()
	if !strings.Contains(read(featureReport), "- [Share](<../use-cases/share.use-case.spec-report.md#goal>) — Share a note.") {
		t.Fatal(read(featureReport))
	}
	if strings.Contains(after, "\n## Use cases\n") || !strings.Contains(read("use-cases.md"), "Share a note.") {
		t.Fatal(after)
	}
	if section(after, "Goals") != section(before, "Goals")+"\n\n- [Four](<goals.md#four>) — Fourth goal." || !strings.Contains(read("goals.md"), "Fourth goal.") {
		t.Fatal(after)
	}
	use := "- [My notebook](<concept-uses/my-notebook.concept-use.spec-report.md>) — The current author's notes.\n  - Uses [Notebook](<concepts/notebook.concept.spec-report.md>)"
	sync := "### Synchronizations\n\n- [Archive updates](<synchronizations/archive-updates.synchronization.spec-report.md>) — When — The note is saved. · Then — Archive the current version."
	if section(after, "Concepts") != "### Concept uses\n\n"+use+"\n\n"+sync || !strings.Contains(read("concept-uses.md"), use) {
		t.Fatal(after)
	}
	if !strings.Contains(read("concept-uses/my-notebook.concept-use.spec-report.md"), "\n## Uses\n\n[Notebook](<../concepts/notebook.concept.spec-report.md>)") {
		t.Fatal("concept use report lost its definition")
	}
	for _, name := range []string{"Local", "Remote"} {
		want := "- [" + name + "](<requirements/" + strings.ToLower(name) + ".requirement.spec-report.md>) — Keep a " + name + " copy." + refines
		if !strings.Contains(read("requirements.md"), want) {
			t.Fatalf("missing refinement for %s", name)
		}
	}
	for _, file := range []string{"requirements.md", "requirements/keep-notes.requirement.spec-report.md"} {
		page := read(file)
		if !strings.Contains(page, "Reopened notes retain their text.") || strings.Contains(page, "Saved notes remain available.") {
			t.Fatalf("stale shared requirement in %s:\n%s", file, page)
		}
	}
	if !strings.Contains(read(featureReport), "[Keep notes](<../requirements/keep-notes.requirement.spec-report.md>)") {
		t.Fatal("feature lost its link to the shared requirement")
	}
	code, output, msg := invoke("model", "export")
	var snapshot model.Snapshot
	if code != 0 || json.Unmarshal([]byte(output), &snapshot) != nil {
		t.Fatal("model export failed", msg)
	}
	var sharedIDs []string
	for _, element := range snapshot.Elements {
		if element["kind"] == "Requirement" && element["slug"] == "keep-notes" {
			sharedIDs = append(sharedIDs, model.String(element, "id"))
		}
	}
	if len(sharedIDs) != 1 {
		t.Fatalf("expected one Keep notes requirement, got %v", sharedIDs)
	}
	roles := map[string]int{}
	for _, ref := range snapshot.References {
		if model.Object(ref["resolution"])["target"] == sharedIDs[0] {
			roles[model.String(ref, "role")]++
		}
	}
	if roles["illustrates"] != 2 || roles["groups"] != 1 || roles["refines"] != 2 {
		t.Fatalf("unexpected references to Keep notes: %v", roles)
	}
	for _, name := range []string{"archive", "edit"} {
		if !strings.Contains(read("use-cases/"+name+".use-case.spec-report.md"), "requirements/keep-notes.requirement.spec-report.md") {
			t.Fatalf("scenario lost shared requirement: %s", name)
		}
	}

	// A second parent can live in another project, where the authored file
	// remains the destination. Both parents belong beside the same entry.
	const rulePath = "policy/specs/business-rules/retention.business-rule.spec.md"
	writeFixture(t, "policy/dotspec.yaml", "name: policy\n")
	writeFixture(t, "policy/specs/policy.system.spec.md", "---\ntype: system-spec\n---\n# Policy\n\n## Purpose\nRetain notes.\n\n## Scope\nShared policy.\n")
	writeFixture(t, rulePath, "---\ntype: business-rule-spec\n---\n# Retention\n\n## Statement\nKeep archived notes.\n")
	localFile := "specs/requirements/local.requirement.spec.md"
	local, err := os.ReadFile(localFile)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, localFile, string(local)+"[Retention](../../"+rulePath+").\n")
	build()
	if !strings.Contains(read("requirements.md"), refines+"\n  - Refines [Retention](<../"+rulePath+">)") {
		t.Fatalf("missing refinement parents:\n%s", read("requirements.md"))
	}
}
