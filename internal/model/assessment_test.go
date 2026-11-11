package model

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/obligation"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
)

func TestAssessmentPlansPreservePurposeSubjectBasisAndCases(t *testing.T) {
	root := reviewFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	d := newDesignResolver(s)
	for _, tc := range []struct{ kind, name string }{
		{"System", "Document workspace"}, {"Concept", "Trash"}, {"UseCase", "Restore a document"},
		{"GeneralRequirement", "Preserve document metadata on restoration"}, {"Assumption", "Trusted time"},
	} {
		r := designRecord(t, s, tc.kind, tc.name)
		plan := Object(r["validation_plan"])
		for _, f := range []string{"subject", "participants", "method", "success_criteria", "conditions"} {
			if strings.TrimSpace(d.text(plan[f])) == "" {
				t.Errorf("%s missing %s", tc.kind, f)
			}
		}
		if len(Strings(Object(plan["basis"])["references"])) == 0 {
			t.Errorf("%s missing contextual basis", tc.kind)
		}
	}
	sys := designRecord(t, s, "System", "Document workspace")
	cases := reviewObjects(Object(sys["validation_plan"])["cases"])
	if len(cases) != 1 || !strings.Contains(d.text(cases[0]["expected_result"]), "resumes the task") || d.text(cases[0]["evidence_to_retain"]) == "" {
		t.Fatal("lost assessment case", cases)
	}
	if len(Strings(sys["supports"])) != 0 || sys["verified"] != nil || sys["validated"] != nil {
		t.Fatal("plan inferred an association or result")
	}
	review, err := ReviewMarkdown(s, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"## Requirements quality", "## Needs and intent", "## Assumptions: evidence and justified reliance", "## Validation planning", "## System verification planning", "No system verification or measurement guidance recorded", "no separately identified criterion, conditions or applicability", "../specs/jobs/recover.job.spec.md", "not a guarantee"} {
		if !strings.Contains(review, want) {
			t.Errorf("missing review distinction %q", want)
		}
	}
	forged := cloneSnapshot(t, s)
	plan := Object(designRecord(t, forged, "System", "Document workspace")["validation_plan"])
	Object(plan["basis"])["references"] = []string{String(sys, "id")}
	requireFinding(t, forged, "reference.kind")
}

func TestNestedAndEmptyValidationPlansAreExplicitAndIdempotent(t *testing.T) {
	root := traceFixture(t)
	// Named actions inherit the same plan role. Repeated supplementary sections
	// preserve content and append cases without duplicating the nested mapping pass.
	traceEdit(t, root, "specs/concepts/trash.concept.spec.md", "### Restore", `#### Validation plan

##### Subject

The proposed Delete action.

##### Cases

###### Accidental deletion

Discuss whether recovery remains possible.

#### Validation plan

##### Basis

The operational principle of reversible deletion.

##### Cases

###### Intentional deletion

Discuss whether the recovery policy is appropriate.

### Restore`)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Rationale", "## Validation plan\n\n## Rationale")
	result := conformanceExport(t, root)
	conformancePayloads(t, result.Snapshot)
	if result.ExitCode != 0 {
		t.Fatal(result.Snapshot.Diagnostics)
	}
	action := designRecord(t, result.Snapshot, "ConceptAction", "Delete")
	plan := Object(action["validation_plan"])
	cases := reviewObjects(plan["cases"])
	if len(cases) != 2 {
		t.Fatal("cases overwritten or duplicated", cases)
	}
	d := newDesignResolver(result.Snapshot)
	if strings.Count(d.text(plan["subject"]), "proposed Delete action") != 1 {
		t.Fatal("subject duplicated")
	}
	req := designRecord(t, result.Snapshot, "GeneralRequirement", "Preserve document metadata on restoration")
	if String(Object(Object(req["validation_plan"])["content"]), "presence") != "empty" {
		t.Fatal("empty plan fabricated", req)
	}
	uc := designRecord(t, result.Snapshot, "UseCase", "Restore a document")
	if uc["validation_plan"] != nil {
		t.Fatal("validation inferred from supports or narrative")
	}
}

func TestMeasurementReviewAgreesAcrossCommandsAndModel(t *testing.T) {
	cases := []struct {
		name, body string
		missing    []string
	}{
		{"measurement only", "## Measurement plan\nCollect timings.\n", []string{"criterion", "conditions or applicability"}},
		{"measure without comparison", "## Measurement plan\nCollect timings.\n## Applies to\nManaged production requests.\n## Criterion\n### Measure\nElapsed time.\n", []string{"comparison rule and target"}},
		{"complete narrative criterion", "## Measurement plan\nCollect timings.\n## Applies to\nManaged production requests.\n## Criterion\nThe 95th percentile is at most 100 ms.\n", nil},
		{"structured criterion", "## Measurement plan\nCollect timings.\n### Conditions\nProduction requests.\n## Criterion\n### Comparison\nAt most.\n### Target\n100 ms at the 95th percentile.\n", nil},
		{"measurement success criteria", "## Measurement plan\n### Procedure\nCollect timings.\n### Conditions\nProduction requests.\n### Success criteria\nThe 95th percentile is at most 100 ms.\n", nil},
		{"explicit verification guidance", "## Measurement plan\nCollect timings.\n## Verification plan\nCompare timings to the externally defined criterion for the stated workload.\n", nil},
		{"validation is distinct", "## Validation plan\nObserve whether the proposed latency serves user needs.\n## Measurement plan\nCollect timings.\n", []string{"criterion", "conditions or applicability"}},
	}
	for _, kind := range []string{"requirement", "quality"} {
		for _, tc := range cases {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				root := t.TempDir()
				dir := filepath.Join(root, "specs", "requirements")
				if kind == "quality" {
					dir = filepath.Join(root, "specs", "qualities")
				}
				// Quality uses its established collection name.
				if kind == "quality" {
					dir = filepath.Join(root, "specs", quality.Collection)
				}
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				source := "---\ntype: " + kind + "\n---\n# Latency\n## Statement\nThe service must respond within the agreed bound.\n" + tc.body
				if err := os.WriteFile(filepath.Join(dir, "latency."+kind+".spec.md"), []byte(source), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "specs", "system.spec.md"), []byte("---\ntype: system\n---\n# Service\n## Purpose\nServe requests.\n## Scope\nRequest handling.\n"), 0644); err != nil {
					t.Fatal(err)
				}
				native, err := obligation.Parse([]byte(source), kind)
				if err != nil {
					t.Fatal(err)
				}
				if !slices.Equal(native.MeasurementGuidance.Missing(), tc.missing) {
					t.Fatalf("native: %v", native.MeasurementGuidance.Missing())
				}
				hasAdvice := false
				if kind == "quality" {
					q, err := quality.Parse([]byte(source))
					if err != nil {
						t.Fatal(err)
					}
					for _, f := range quality.Check(q) {
						hasAdvice = hasAdvice || f.Rule == "quality.measurement-basis"
					}
				} else {
					q, err := requirement.Parse([]byte(source))
					if err != nil {
						t.Fatal(err)
					}
					for _, f := range requirement.Check(q) {
						hasAdvice = hasAdvice || f.Rule == "requirement.measurement-basis"
					}
				}
				if hasAdvice != (len(tc.missing) > 0) {
					t.Fatal("native advice disagrees")
				}
				result := conformanceExport(t, root)
				conformancePayloads(t, result.Snapshot)
				if result.ExitCode != 0 {
					t.Fatal(result.Snapshot.Diagnostics)
				}
				modelKind := "GeneralRequirement"
				if kind == "quality" {
					modelKind = "QualityRequirement"
				}
				r := designRecord(t, result.Snapshot, modelKind, "Latency")
				d := newDesignResolver(result.Snapshot)
				if got := d.measurementGuidance(r).Missing(); !slices.Equal(got, tc.missing) {
					t.Fatalf("model %v, want %v", got, tc.missing)
				}
			})
		}
	}
}
