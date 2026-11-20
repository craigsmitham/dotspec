package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/testproject"
)

func reviewFixture(t *testing.T) string {
	t.Helper()
	return testproject.Open(t, "requirements-review")
}

func TestReviewModelSourceFidelityAndSemanticBoundaries(t *testing.T) {
	root := reviewFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
	d := newDesignResolver(s)
	sys := designRecord(t, s, "System", "Document workspace")
	partner := designRecord(t, s, "System", "Platform clock")
	if ids := Strings(sys["interacts_with"]); len(ids) != 1 || ids[0] != String(partner, "id") || sys["part_of"] != nil {
		t.Fatal("interaction confused with structure", sys)
	}
	if d.text(sys["out_of_scope"]) == "" || d.text(Object(sys["project_references"])["content"]) == "" {
		t.Fatal("unmapped system context")
	}
	required := Strings(sys["required_interfaces"])
	provided := Strings(sys["provided_interfaces"])
	if len(required) != 1 || len(provided) != 1 || len(Strings(Object(d.index[required[0]]["contracts"])["references"])) != 1 || d.index[provided[0]]["contracts"] != nil {
		t.Fatal("interface direction/contract information lost")
	}
	assumption := designRecord(t, s, "Assumption", "Trusted time")
	for _, field := range []string{"basis", "applicability"} {
		if d.text(assumption[field]) == "" {
			t.Fatal("unmapped assumption", field)
		}
	}
	if strings.Contains(d.text(assumption["statement"]), "Assumption validation plan") {
		t.Fatal("assumption guidance became its premise")
	}
	questions := reviewObjects(sys["open_questions"])
	if len(questions) != 2 {
		t.Fatal(questions)
	}
	if len(questions[0]) != 1 || !strings.Contains(d.text(questions[0]["question"]), "Before acceptance") {
		t.Fatal("question prose or linked context lost")
	}
	count := 0
	for _, r := range s.Elements {
		if String(r, "kind") == "ScenarioStep" && len(Strings(r["illustrates"])) > 0 {
			count++
			if String(r, "authored_label") != "2" || len(Strings(r["illustrates"])) != 2 {
				t.Fatal(r)
			}
		}
	}
	if count != 1 {
		t.Fatal("ordinary references incorrectly promoted or explicit illustration lost", count)
	}
	concept := designRecord(t, s, "Concept", "Trash")
	if !strings.Contains(d.text(concept["effects_table"]), "| Expire | unchanged | false | |") {
		t.Fatal("effects table lost")
	}
	assertions := map[string]int{}
	for _, a := range s.Relationships() {
		assertions[String(a, "projection_path")]++
	}
	if assertions["ScenarioStep.illustrates"] != 2 || assertions["System.interacts_with"] != 1 {
		t.Fatal(assertions)
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "System", "Document workspace")["interacts_with"] = []string{String(concept, "id")}
	requireFinding(t, forged, "shape.invalid")
	forged = cloneSnapshot(t, s)
	for _, r := range forged.Elements {
		if String(r, "kind") == "ScenarioStep" && String(r, "authored_label") == "1" {
			r["illustrates"] = []string{String(concept, "id")}
			break
		}
	}
	requireFinding(t, forged, "shape.invalid")
}

func TestReviewViewsAreAdvisoryAndRebaseSources(t *testing.T) {
	root := reviewFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	text, err := ReviewMarkdown(result.Snapshot, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"No contribution path recorded", "Contribution path recorded", "ScenarioStep 2", "no rule or action relationship recorded", "no native contract reference recorded", "Before acceptance", "Review required for **recovery copy**", "../specs/document-workspace.system.spec.md#reduce-retention-cost", "../project-plan.md#retention-boundary-decision"} {
		if !strings.Contains(text, want) {
			t.Errorf("review missing %q", want)
		}
	}
	pegs, err := PEGSMarkdown(result.Snapshot, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Project", "## Environment and boundaries", "## Goals", "## System", "../project-plan.md", "../contracts/time.openapi.yaml", "## Open questions"} {
		if !strings.Contains(pegs, want) {
			t.Errorf("PEGS missing %q", want)
		}
	}
	design, err := DesignMarkdown(result.Snapshot, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Effects table", "| Expire | unchanged | false | |", "Blank cells are unspecified"} {
		if !strings.Contains(design, want) {
			t.Errorf("design missing %q", want)
		}
	}
	overview, err := OverviewMarkdown(result.Snapshot, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(overview, "Out of scope") || !strings.Contains(overview, "Time API") {
		t.Fatal("overview omits boundary information")
	}
	saved := cloneSnapshot(t, result.Snapshot)
	for _, view := range []struct {
		name   string
		render func(*Snapshot, string) (string, error)
		want   string
	}{{"review", ReviewMarkdown, text}, {"overview", OverviewMarkdown, overview}, {"design", DesignMarkdown, design}, {"pegs", PEGSMarkdown, pegs}} {
		again, err := view.render(saved, filepath.Join(root, "docs"))
		if err != nil || again != view.want {
			t.Errorf("%s differs after JSON round trip: %v", view.name, err)
		}
	}
}

func TestIllustrationIsExplicitAndRejectsWrongEndpoint(t *testing.T) {
	for _, tc := range []struct {
		name, replacement string
		invalid           bool
	}{
		{"citation", "Related: [Restore retained content](../requirements/restore-content.requirement.spec.md).", false},
		{"wrong endpoint", "**Illustrates:** [Workspace](../document-workspace.system.spec.md).", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := reviewFixture(t)
			traceEdit(t, root, "specs/use-cases/restore.use-case.spec.md", "**Illustrates:** [Restore retained content](../requirements/restore-content.requirement.spec.md), [Restore](../concepts/trash.concept.spec.md#restore).", tc.replacement)
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.invalid {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			checkSourceText(t, result.Snapshot)
			for _, r := range result.Snapshot.Elements {
				if String(r, "kind") == "ScenarioStep" && len(Strings(r["illustrates"])) > 0 {
					t.Fatal("invalid illustration projected", r)
				}
			}
		})
	}
}

func TestRepeatedComponentInterfaceSectionsRetainAllEntries(t *testing.T) {
	root := reviewFixture(t)
	conformanceWrite(t, root, "worker/specs/retention-worker.component.spec.md", `---
type: component-spec
---
# Retention worker
## Part of
[Workspace](../../specs/document-workspace.system.spec.md)

## Required interfaces

### Time readings

**ID:** component-time

Read the current time.

#### Contracts

[Clock contract](../../contracts/time.openapi.yaml).

#### Contracts

[Project references](../../project-plan.md).

## Required interfaces

### Storage reads

Read retained items.

## Implementation

A batch worker can use [Time readings](#component-time).
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	checkSourceText(t, result.Snapshot)
	component := designRecord(t, result.Snapshot, "Component", "Retention worker")
	d := newDesignResolver(result.Snapshot)
	items := Strings(component["required_interfaces"])
	if len(items) != 2 || len(Strings(Object(d.index[items[0]]["contracts"])["references"])) != 2 {
		t.Fatal("repeated interface content overwritten", items)
	}
	if strings.Count(d.text(Object(component["implementation"])["content"]), "A batch worker") != 1 {
		t.Fatal("nested explanation duplicated")
	}
}
