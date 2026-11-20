package model

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/craigsmitham/dotspec/internal/obligation"
	"github.com/craigsmitham/dotspec/internal/requirement"
)

func TestMeasurementReviewAgreesAcrossCommandsAndModel(t *testing.T) {
	cases := []struct {
		name, body string
		missing    []string
	}{
		{"measure only", "## Criterion\n### Measure\nElapsed time.\n", []string{"target", "conditions or applicability"}},
		{"measure without target", "## Applies to\nManaged production requests.\n## Criterion\n### Measure\nElapsed time.\n", []string{"target"}},
		{"complete narrative criterion", "## Applies to\nManaged production requests.\n## Criterion\nThe 95th percentile is at most 100 ms.\n", nil},
		{"structured criterion", "## Applies to\n### Conditions\nProduction requests.\n## Criterion\n### Measure\nElapsed time.\n### Target\nAt most 100 ms at the 95th percentile.\n", nil},
		{"plan success criteria are assessment guidance", "## Criterion\n### Measure\nElapsed time.\n## System verification plan\n### Method\nCollect timings.\n### Conditions\nProduction requests.\n### Success criteria\nThe 95th percentile is at most 100 ms.\n", []string{"target", "conditions or applicability"}},
		{"empty verification plan", "## Criterion\n### Measure\nElapsed time.\n## System verification plan\n", []string{"target", "conditions or applicability"}},
		{"assessment conditions do not define applicability", "## Criterion\n### Measure\nElapsed time.\n### Target\nAt most 100 ms at the 95th percentile.\n## System verification plan\n### Conditions\nTest with ten synthetic requests.\n", []string{"conditions or applicability"}},
		{"criterion references authoritative statement", "## Applies to\nManaged production requests.\n## Criterion\nUse the bound defined in [Statement](#statement).\n## System verification plan\nCollect timings and compare them to [Criterion](#criterion).\n", nil},
		{"no explicitly recorded measure", "## Applies to\nManaged production requests.\n## System verification plan\nInspect compliance with the specified response behavior.\n", nil},
		{"validation is distinct", "## System validation plan\nObserve whether the proposed latency serves user needs.\n## Criterion\n### Measure\nElapsed time.\n", []string{"target", "conditions or applicability"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "specs", "requirements")
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			source := "---\ntype: requirement-spec\n---\n# Latency\n## Statement\nThe service must respond within the agreed bound.\n" + tc.body
			if err := os.WriteFile(filepath.Join(dir, "latency.requirement.spec.md"), []byte(source), 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "specs", "service.system.spec.md"), []byte("---\ntype: system-spec\n---\n# Service\n## Purpose\nServe requests.\n## Scope\nRequest handling.\n"), 0644); err != nil {
				t.Fatal(err)
			}
			native, err := obligation.Parse([]byte(source), "requirement")
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(native.MeasurementGuidance.Missing(), tc.missing) {
				t.Fatalf("native: %v", native.MeasurementGuidance.Missing())
			}
			hasAdvice := false
			q, err := requirement.Parse([]byte(source))
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range requirement.Check(q) {
				hasAdvice = hasAdvice || f.Rule == "requirement.measurement-basis"
				if f.Rule == "requirement.measurement-basis" && (f.Location.Line == 0 || f.Location.Element != "Criterion") {
					t.Fatalf("measurement advice has wrong source location: %+v", f.Location)
				}
			}
			if hasAdvice != (len(tc.missing) > 0) {
				t.Fatal("native advice disagrees")
			}
			result := conformanceExport(t, root)
			checkSourceText(t, result.Snapshot)
			if result.ExitCode != 0 {
				t.Fatal(result.Snapshot.Diagnostics)
			}
			r := designRecord(t, result.Snapshot, "Requirement", "Latency")
			d := newDesignResolver(result.Snapshot)
			if got := d.measurementGuidance(r).Missing(); !slices.Equal(got, tc.missing) {
				t.Fatalf("model %v, want %v", got, tc.missing)
			}
		})
	}
}
