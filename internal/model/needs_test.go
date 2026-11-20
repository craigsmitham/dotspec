package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/testproject"
)

func needsFixture(t *testing.T) string {
	t.Helper()
	return testproject.Open(t, "user-needs")
}

func TestNeedsRelationshipsAndEvidence(t *testing.T) {
	root := needsFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	checkSourceText(t, s)
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
	text, err := NeedsMarkdown(s, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"DesiredOutcome: Timely awareness", "-->|applies to|", "-->|depends on|", "-->|provided by|", "Minimize the time between", "../specs/jobs/reach-work.job.spec.md#timely-awareness"} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q", want)
		}
	}
	// Ownership is not applicability, and provider claims must have source evidence.
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "DesiredOutcome", "Timely awareness")["applies_to"] = []string{String(j, "id")}
	requireFinding(t, forged, "shape.invalid")
	forged = cloneSnapshot(t, s)
	designRecord(t, forged, "Capability", "Service-change detection")["provided_by"] = []string{String(j, "id")}
	requireFinding(t, forged, "shape.invalid")
}

func TestNeedsReferencesRejectWrongKindsAndPreserveUnresolvedNarrative(t *testing.T) {
	for _, tc := range []struct {
		name, old, replacement string
		invalid                bool
	}{
		{"wrong capability", "[Service-change detection](service-change-detection.capability.spec.md)", "[Feature](../features/disruption-alerts.feature.spec.md)", true},
		{"wrong provider", "[Notification service](../../components/notification-service/specs/notification-service.component.spec.md)", "[Feature](../features/disruption-alerts.feature.spec.md)", true},
		{"provider prose", "[Notification service](../../components/notification-service/specs/notification-service.component.spec.md)", "Provider not yet selected.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := needsFixture(t)
			traceEdit(t, root, "specs/capabilities/disruption-communication.capability.spec.md", tc.old, tc.replacement)
			result := conformanceExport(t, root)
			if (result.ExitCode != 0) != tc.invalid {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			checkSourceText(t, result.Snapshot)
			if tc.name == "provider prose" {
				before, err := NeedsMarkdown(result.Snapshot, root)
				if err != nil || !strings.Contains(before, tc.replacement) {
					t.Fatalf("provider qualification missing: %v", err)
				}
				after, err := NeedsMarkdown(cloneSnapshot(t, result.Snapshot), root)
				if err != nil || before != after {
					t.Fatalf("needs view changed after saving the snapshot: %v", err)
				}
			}
		})
	}
}
