package model

import (
	"strings"
	"testing"
)

func TestObligationFormsShareAssessmentAndClassification(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/service.system.spec.md", "---\ntype: system-spec\n---\n# Service\n\n## Purpose\n\nServe requests.\n\n## Scope\n\nThis service.\n\n## Quality attributes\n\n")
	conformanceWrite(t, root, "specs/quality/speed/speed.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nid: speed\n---\n# Speed\n\n## Definition\n\nResponse speed.\n")

	body := "## Statement\n\nRespond within one second.\n\n## Applies to\n\n**Targets:** [Service](../service.system.spec.md)\n\n## Criterion\n\n### Target\n\n1 second.\n\n## Rationale\n\n### Method\n\nCheck the response to every request.\n"
	conformanceWrite(t, root, "specs/requirements/speed.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [speed]\n---\n# Speed\n\n"+body)
	clause := "---\ntype: use-case-spec\n---\n# Respond\n\n## Obligations\n\n### Response time\n\n**Quality attributes:** speed\n\nRespond within one second.\n\n" + strings.ReplaceAll(strings.TrimPrefix(body, "## Statement\n\nRespond within one second.\n\n"), "##", "####")
	conformanceWrite(t, root, "specs/use-cases/respond.use-case.spec.md", clause)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatalf("%v", result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	for _, kind := range []string{"Requirement", "ObligationClause"} {
		r := elementKind(t, s, kind)
		for _, field := range []string{"statement", "applicability", "criterion"} {
			if r[field] == nil {
				t.Fatalf("%s lacks %s", kind, field)
			}
		}
		if len(Strings(r["quality_attributes"])) != 1 || len(Strings(r["applies_to"])) != 1 {
			t.Fatalf("%s lost semantic links: %v", kind, r)
		}
	}
	checkSourceText(t, s)
	traceEdit(t, root, "specs/use-cases/respond.use-case.spec.md", "**Quality attributes:** speed", "**Quality attributes:** unknown")
	result = conformanceExport(t, root)
	found := false
	for _, d := range result.Snapshot.Diagnostics {
		found = found || String(d, "code") == "obligation.undefined-attribute"
	}
	if !found {
		t.Fatal("missing undefined clause classification advice")
	}
	checkSourceText(t, result.Snapshot)
}

func TestTraceabilityDoesNotConferObligationContainment(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/service.system.spec.md", "---\ntype: system-spec\n---\n# Service\n\n## Purpose\n\nServe requests.\n\n## Scope\n\nThis service.\n\n## Goals\n\n### Faster service\n\nRespond promptly.\n\n#### Obligations\n\n##### Response time\n\nRespond within one second.\n")
	result := conformanceExport(t, root)
	found := false
	for _, d := range result.Snapshot.Diagnostics {
		found = found || String(d, "code") == "obligation.container"
	}
	if result.ExitCode != 1 || !found {
		t.Fatalf("expected unsupported container error: %d %v", result.ExitCode, result.Snapshot.Diagnostics)
	}
	for _, r := range result.Snapshot.Elements {
		if String(r, "kind") == "ObligationClause" {
			t.Fatal("goal acquired a clause", r)
		}
	}
	checkSourceText(t, result.Snapshot)
}

func TestClauseQualityLabelsRejectAmbiguousAuthoring(t *testing.T) {
	for _, tc := range []struct{ name, field, code string }{
		{"empty", "**Quality attributes:**", "identity.invalid"},
		{"invalid", "**Quality attributes:** speed.fast", "identity.invalid"},
		{"duplicate label", "**Quality attributes:** speed, speed", "identity.duplicate"},
		{"duplicate field", "**Quality attributes:** speed\n**Quality attributes:** other", "structure.duplicate-field"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := ParseContent("respond.use-case.spec.md", []byte("---\ntype: use-case-spec\n---\n# Respond\n\n## Obligations\n\n### Speed\n\n"+tc.field+"\n\nRespond promptly.\n"))
			mapped := MapDomain(doc, "use-case")
			found := false
			for _, d := range mapped.Diagnostics {
				found = found || String(d, "code") == tc.code
			}
			if !found || mapped.RootID != "" {
				t.Fatalf("ambiguous clause was accepted: %+v", mapped)
			}
		})
	}
}
