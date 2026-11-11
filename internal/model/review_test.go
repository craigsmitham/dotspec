package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func reviewFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/requirements-review")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestReviewModelSourceFidelityAndSemanticBoundaries(t *testing.T) {
	root := reviewFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	d := newDesignResolver(s)
	sys := designRecord(t, s, "System", "Document workspace")
	partner := designRecord(t, s, "System", "Platform clock")
	if ids := Strings(sys["interacts_with"]); len(ids) != 1 || ids[0] != String(partner, "id") || sys["part_of"] != nil {
		t.Fatal("interaction confused with structure", sys)
	}
	if d.text(sys["out_of_scope"]) == "" || d.text(Object(sys["project_context"])["content"]) == "" {
		t.Fatal("unmapped system context")
	}
	required := reviewObjects(Object(sys["required_interfaces"])["interfaces"])
	provided := reviewObjects(Object(sys["provided_interfaces"])["interfaces"])
	if len(required) != 1 || len(provided) != 1 || len(Strings(Object(required[0]["contracts"])["references"])) != 1 || provided[0]["contracts"] != nil {
		t.Fatal("interface direction/contract information lost")
	}
	assumption := designRecord(t, s, "Assumption", "Trusted time")
	for _, field := range []string{"basis", "applicability"} {
		if d.text(assumption[field]) == "" {
			t.Fatal("unmapped assumption", field)
		}
	}
	if !strings.Contains(d.text(Object(assumption["validation_plan"])["method"]), "independent reference") {
		t.Fatal("missing validation method")
	}
	if strings.Contains(d.text(assumption["statement"]), "Validation plan") {
		t.Fatal("assumption guidance became its premise")
	}
	questions := reviewObjects(Object(sys["open_questions"])["questions"])
	if len(questions) != 2 {
		t.Fatal(questions)
	}
	for _, field := range []string{"raised_by", "raised_on", "resolve_by", "importance", "evidence_needed"} {
		if d.text(questions[0][field]) == "" {
			t.Fatal("missing question field", field)
		}
	}
	if len(Strings(Object(questions[0]["tracking_reference"])["references"])) != 1 || questions[1]["resolve_by"] != nil {
		t.Fatal("question metadata inferred or dropped")
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
	rows := reviewObjects(Object(concept["effects_table"])["rows"])
	if len(rows) != 4 {
		t.Fatal("lost effects table rows", rows)
	}
	last := reviewObjects(rows[3]["cells"])
	if len(last) != 4 || strings.TrimSpace(d.text(last[3])) != "" || !strings.Contains(d.text(last[1]), "unchanged") {
		t.Fatal("blank effects inferred as unchanged", last)
	}
	if strings.Count(d.text(concept["implementation_suggestions"]), "An index could") != 1 {
		t.Fatal("explanation duplicated or lost")
	}
	assertions := map[string]int{}
	for _, a := range s.Assertions {
		assertions[String(a, "kind")]++
	}
	if assertions["Illustrates"] != 2 || assertions["InteractsWith"] != 1 {
		t.Fatal(assertions)
	}
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "System", "Document workspace")["interacts_with"] = []string{String(concept, "id")}
	requireFinding(t, forged, "reference.kind")
	forged = cloneSnapshot(t, s)
	for _, r := range forged.Elements {
		if String(r, "kind") == "ScenarioStep" && String(r, "authored_label") == "1" {
			r["illustrates"] = []string{String(concept, "id")}
			break
		}
	}
	requireFinding(t, forged, "reference.kind")
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
	for _, want := range []string{"No contribution path recorded", "Contribution path recorded", "ScenarioStep 2", "no rule or action relationship recorded", "no native contract reference recorded", "No system verification or measurement guidance recorded", "Before acceptance", "no Raised by, Raised on, Resolve by, Importance, Evidence needed recorded", "Review required for **recovery copy**", "../specs/system.spec.md#reduce-retention-cost", "../project-plan.md#retention-boundary-decision", "nonbinding"} {
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
	for _, want := range []string{"Effects table", "| Expire | unchanged | false | |", "Blank effects cells are unspecified", "Implementation suggestions (nonbinding)"} {
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
	// Review a serialized snapshot too; embedded values are then []any, not []Record.
	again, err := ReviewMarkdown(cloneSnapshot(t, result.Snapshot), filepath.Join(root, "docs"))
	if err != nil || again != text {
		t.Fatal("review differs after JSON round trip", err)
	}
}

func TestIllustrationIsExplicitAndRejectsWrongEndpoint(t *testing.T) {
	for _, tc := range []struct {
		name, replacement string
		invalid           bool
	}{
		{"citation", "Related: [Restore retained content](#restore-content).", false},
		{"wrong endpoint", "**Illustrates:** [Workspace](../system.spec.md).", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := reviewFixture(t)
			traceEdit(t, root, "specs/use-cases/restore.use-case.spec.md", "**Illustrates:** [Restore retained content](#restore-content), [Restore](../concepts/trash.concept.spec.md#restore).", tc.replacement)
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.invalid {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			conformancePayloads(t, result.Snapshot)
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
	traceEdit(t, root, "specs/system.spec.md", "## Concept uses", `## Components

### Retention worker

#### Required interfaces

##### Time readings

Read the current time.

###### Contracts

[Clock contract](../contracts/time.openapi.yaml).

###### Contracts

[Project context](../project-plan.md).

#### Required interfaces

##### Storage reads

Read retained items.

#### Implementation hints

A batch worker is one implementation option.

## Concept uses`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	conformancePayloads(t, result.Snapshot)
	component := designRecord(t, result.Snapshot, "Component", "Retention worker")
	items := reviewObjects(Object(component["required_interfaces"])["interfaces"])
	if len(items) != 2 || len(Strings(Object(items[0]["contracts"])["references"])) != 2 {
		t.Fatal("repeated interface content overwritten", items)
	}
	d := newDesignResolver(result.Snapshot)
	if strings.Count(d.text(component["implementation_suggestions"]), "A batch worker") != 1 {
		t.Fatal("nested explanation duplicated")
	}
}
