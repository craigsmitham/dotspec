package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func needsFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/user-needs")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestNeedsRelationshipsAndEvidence(t *testing.T) {
	root := needsFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	outcome := designRecord(t, s, "DesiredOutcome", "Timely awareness")
	j := designRecord(t, s, "Job", "Reach work predictably")
	if outcome["owner"] != j["id"] || len(Strings(outcome["applies_to"])) != 2 {
		t.Fatal(outcome)
	}
	c := designRecord(t, s, "Capability", "Disruption communication")
	if len(Strings(c["depends_on"])) != 1 || len(Strings(c["provided_by"])) != 1 || len(Strings(c["supports"])) != 2 {
		t.Fatal(c)
	}
	providerless := designRecord(t, s, "Capability", "Service-change detection")
	if len(Strings(providerless["provided_by"])) != 0 {
		t.Fatal("inferred provider from declaration", providerless)
	}
	pc := designRecord(t, s, "ProgressCase", "Repeated late arrivals")
	d := newDesignResolver(s)
	for _, f := range []string{"circumstances", "struggling_moment", "participants", "alternatives_considered", "tradeoffs_and_selection_criteria"} {
		if d.text(pc[f]) == "" {
			t.Fatal("unmapped case field", f)
		}
	}
	if !strings.Contains(d.text(Object(pc["evidence"])["research_basis"]), "Illustrative hypothesis") {
		t.Fatal("lost evidence basis")
	}
	text, err := NeedsMarkdown(s, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"DesiredOutcome: Timely awareness", "-->|applies to|", "-->|depends on|", "-->|provided by|", "Minimize the time between", "Illustrative hypothesis", "../specs/jobs/reach-work.job.spec.md#timely-awareness"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q", want)
		}
	}
	// Ownership is not applicability, and provider claims must have source evidence.
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "DesiredOutcome", "Timely awareness")["applies_to"] = []string{String(j, "id")}
	requireFinding(t, forged, "concept.projection")
	forged = cloneSnapshot(t, s)
	designRecord(t, forged, "Capability", "Service-change detection")["provided_by"] = []string{String(j, "id")}
	requireFinding(t, forged, "reference.kind")
}

func TestNeedsReferencesRejectWrongKindsAndPreserveUnresolvedNarrative(t *testing.T) {
	for _, tc := range []struct {
		name, old, replacement string
		invalid                bool
	}{
		{"wrong capability", "[Service-change detection](#change-detection)", "[Feature](#disruption-alerts)", true},
		{"wrong provider", "[Notification service](#notification-service)", "[Feature](#disruption-alerts)", true},
		{"provider prose", "[Notification service](#notification-service)", "Provider not yet selected.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := needsFixture(t)
			traceEdit(t, root, "specs/system.spec.md", tc.old, tc.replacement)
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.invalid {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			conformancePayloads(t, result.Snapshot)
		})
	}
}

func TestNeedsDiagramLabelsExplicitCaseApplicability(t *testing.T) {
	root := needsFixture(t)
	traceEdit(t, root, "specs/requirements/notify-commuters.requirement.spec.md", "**Targets:** [Notification service](../system.spec.md#notification-service)", "**Targets:** [Case](../jobs/reach-work.job.spec.md#repeated-late-arrivals)")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	text, err := NeedsMarkdown(result.Snapshot, root)
	if err != nil {
		t.Fatal(err)
	}
	diagram := strings.Split(strings.Split(text, "```mermaid\n")[1], "```")[0]
	if !strings.Contains(diagram, "ProgressCase: Repeated late arrivals") || strings.Contains(diagram, "|declares|") {
		t.Fatal(diagram)
	}
}
