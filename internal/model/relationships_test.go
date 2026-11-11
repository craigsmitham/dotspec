package model

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func relationshipFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	conformanceWrite(t, root, "specs/system.spec.md", "---\ntype: system\n---\n# Authoring context\n## Purpose\nProvide context.\n## Scope\nThis repository.\n")
	conformanceWrite(t, root, "external/specs/system.spec.md", "---\ntype: system\n---\n# Delivery system\n## Purpose\nDeliver links.\n## Scope\nExternal delivery.\n")
	conformanceWrite(t, root, "specs/jobs/share.job.spec.md", "---\ntype: job\n---\n# Share information\n## Job\nWhen collaborating, I want to share information so others can act.\n## Steps\n### Send a link\nSend the address to a recipient.\n## Outcomes\n- Minimize time to share a usable address.\n")
	conformanceWrite(t, root, "specs/use-cases/share.use-case.spec.md", `---
type: use-case
---
# Share a link
## Goal
Deliver an address to the recipient.
## System boundary
**Subject system:** [Delivery](../../external/specs/system.spec.md)

The actor submits the address through the delivery service.
## Supports
[Sharing](../jobs/share.job.spec.md) and [sending](../jobs/share.job.spec.md#send-a-link).
## Related jobs
[Sharing](../jobs/share.job.spec.md) provides research context.
## Main success scenario
1. The actor submits the address.
2. The system confirms delivery.
## Success criterion
The recipient can read the address.
`)
	conformanceWrite(t, root, "specs/business-rules/consent.business-rule.spec.md", `---
type: business-rule
---
# Recipient consent
## Statement
The recipient must have consented to delivery.
## Applies to
**Targets:** [Sharing](../use-cases/share.use-case.spec.md)

This rule applies to sharing links.
`)
	conformanceWrite(t, root, "specs/requirements/delivery.requirement.spec.md", `---
type: requirement
---
# Confirm delivery
## Statement
The service must confirm delivery.
## Applies to
**Targets:** [Sharing](../use-cases/share.use-case.spec.md)
### Subjects
[Delivery context](../../external/specs/system.spec.md) is contextual background.
`)
	return root
}

func TestExplicitSemanticRelationships(t *testing.T) {
	root := relationshipFixture(t)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatalf("unexpected authoring error: %+v", result.Snapshot.Diagnostics)
	}
	s := result.Snapshot
	uc := elementKind(t, s, "UseCase")
	context := designRecord(t, s, "System", "Authoring context")
	boundary := designRecord(t, s, "System", "Delivery system")
	if uc["system_context"] != context["id"] || uc["subject_system"] != boundary["id"] {
		t.Fatalf("directory context conflated with explicit boundary: %+v", uc)
	}
	if len(Strings(uc["supports"])) != 2 {
		t.Fatal("missing supported job/step", uc)
	}
	for _, kind := range []string{"BusinessRule", "GeneralRequirement"} {
		r := elementKind(t, s, kind)
		if !slices.Equal(Strings(r["applies_to"]), []string{String(uc, "id")}) {
			t.Fatal(r)
		}
	}
	projections := map[string]int{}
	for _, a := range s.Assertions {
		projections[String(a, "projection_path")]++
	}
	for path, count := range map[string]int{"UseCase.subject_system": 1, "TraceableElement.supports": 2, "BusinessRule.applies_to": 1, "Obligation.applies_to": 1} {
		if projections[path] != count {
			t.Fatalf("projection %s: %d, want %d", path, projections[path], count)
		}
	}
	forged := cloneSnapshot(t, s)
	elementKind(t, forged, "UseCase")["subject_system"] = context["id"]
	requireFinding(t, forged, "concept.projection")
	conformancePayloads(t, s)
}

func TestContextualLinksDoNotAssertSemanticRelations(t *testing.T) {
	root := relationshipFixture(t)
	for _, file := range []string{"specs/use-cases/share.use-case.spec.md", "specs/business-rules/consent.business-rule.spec.md", "specs/requirements/delivery.requirement.spec.md"} {
		p := filepath.Join(root, file)
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		text := strings.NewReplacer("**Subject system:** ", "", "**Targets:** ", "", "## Supports", "## Notes").Replace(string(data))
		if err := os.WriteFile(p, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	s := conformanceExport(t, root).Snapshot
	for _, a := range s.Assertions {
		if slices.Contains([]string{"HasSubjectSystem", "Supports", "AppliesTo"}, String(a, "kind")) {
			t.Fatal("citation inferred a domain claim", a)
		}
	}
}

func TestUnknownConceptScopeDefersDependentRelationships(t *testing.T) {
	root := compositionFixture(t)
	p := filepath.Join(root, "specs/system.spec.md")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(data), "[Shorthand](concepts/shorthand.concept.spec.md).", "[Shorthand](https://example.com/shorthand).", 1)
	if err := os.WriteFile(p, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatalf("unknown scope must not be an incompatible scope: %+v", result.Snapshot.Diagnostics)
	}
	deferred := false
	for _, d := range result.Snapshot.Diagnostics {
		if String(d, "severity") == "error" {
			t.Fatal(d)
		}
		deferred = deferred || strings.Contains(String(d, "message"), "not yet verifiable")
	}
	if !deferred {
		t.Fatal("missing deferred-scope advice")
	}
	requireValid(t, result.Snapshot)
}

func TestRelationshipContractsRejectCorruptSnapshots(t *testing.T) {
	s := conformanceExport(t, compositionFixture(t)).Snapshot
	t.Run("assertions retain declaration role", func(t *testing.T) {
		paths := map[string]bool{}
		for _, a := range s.Assertions {
			if String(a, "kind") == "HasConceptType" {
				paths[String(a, "projection_path")] = true
			}
		}
		if !paths["Concept.type_parameters"] || !paths["Concept.introduced_types"] {
			t.Fatal(paths)
		}
		forged := cloneSnapshot(t, s)
		for _, a := range forged.Assertions {
			if String(a, "kind") == "HasConceptType" {
				a["projection_path"] = "Concept.actions.actions"
				break
			}
		}
		requireFinding(t, forged, "assertion.projection")
	})
	for _, kind := range []string{"UsesConcept", "HasConceptUse", "HasSynchronization", "ConceptDependency", "HasConceptType", "HasTypeBinding"} {
		t.Run(kind+" requires endpoints", func(t *testing.T) {
			forged := cloneSnapshot(t, s)
			for _, a := range forged.Assertions {
				if String(a, "kind") == kind {
					delete(a, "subject")
					delete(a, "object")
					break
				}
			}
			requireFinding(t, forged, "shape.invalid")
		})
	}
	t.Run("systems have only one hierarchy", func(t *testing.T) {
		forged := cloneSnapshot(t, s)
		sys := elementKind(t, forged, "System")
		sys["system_context"] = sys["id"]
		requireFinding(t, forged, "shape.invalid")
	})
	t.Run("bindings retain section content", func(t *testing.T) {
		forged := cloneSnapshot(t, s)
		use := designRecord(t, forged, "ConceptUse", "PublicLinks")
		delete(use, "bindings_content")
		forged.ContentBindings = slices.DeleteFunc(forged.ContentBindings, func(b Record) bool { return b["owner"] == use["id"] && b["field_path"] == "bindings_content" })
		requireFinding(t, forged, "declaration.pair")
	})
}

func TestOptionalDecompositionRetainsSectionContent(t *testing.T) {
	s := conformanceExport(t, compositionFixture(t)).Snapshot
	for _, tc := range []struct{ kind, field, content string }{
		{"ConceptAction", "input_parameters", "inputs"}, {"ConceptAction", "output_parameters", "outputs"},
		{"Synchronization", "when_actions", "when"}, {"Synchronization", "then_actions", "then"},
	} {
		t.Run(tc.field, func(t *testing.T) {
			forged := cloneSnapshot(t, s)
			found := false
			for _, r := range forged.Elements {
				if String(r, "kind") == tc.kind && len(Strings(r[tc.field])) > 0 {
					delete(r, tc.content)
					found = true
					break
				}
			}
			if !found {
				t.Fatal("fixture missing decomposition")
			}
			requireFinding(t, forged, "declaration.content")
		})
	}
}
