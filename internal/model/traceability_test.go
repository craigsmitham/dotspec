package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/testproject"
)

func traceFixture(t *testing.T) string {
	t.Helper()
	return testproject.Open(t, "requirements-traceability")
}
func traceEdit(t *testing.T, root, path, old, replacement string) {
	t.Helper()
	p := filepath.Join(root, path)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), old) {
		t.Fatalf("missing edit target %q", old)
	}
	if err = os.WriteFile(p, []byte(strings.Replace(string(data), old, replacement, 1)), 0644); err != nil {
		t.Fatal(err)
	}
}
func TestLayeredTraceability(t *testing.T) {
	root := traceFixture(t)
	result := conformanceExport(t, root)
	s := result.Snapshot
	if result.ExitCode != 0 || len(s.Diagnostics) != 0 {
		t.Fatalf("unexpected findings: %d %v", result.ExitCode, s.Diagnostics)
	}
	parent := designRecord(t, s, "Requirement", "Restore retained content")
	uc := elementKind(t, s, "UseCase")
	req := designRecord(t, s, "Requirement", "Preserve document metadata on restoration")
	use := elementKind(t, s, "ConceptUse")
	if parent["authored_id"] != "restore-content" || parent["source_document"] == uc["source_document"] {
		t.Fatal(parent)
	}
	for _, tc := range []struct {
		r      Record
		field  string
		target Record
	}{{req, "refines", parent}, {uc, "involves", use}, {parent, "derived_from", elementKind(t, s, "BusinessRule")}, {req, "relies_on", elementKind(t, s, "Assumption")}, {elementKind(t, s, "Job"), "supports", elementKind(t, s, "Goal")}} {
		if !slices.Equal(Strings(tc.r[tc.field]), []string{String(tc.target, "id")}) {
			t.Fatalf("%s: %v", tc.field, tc.r)
		}
		if len(newDesignResolver(s).refs(tc.r, tc.field)) != 1 {
			t.Fatal("missing reference provenance", tc.r)
		}
	}
	if !strings.Contains(newDesignResolver(s).text(req["refines_content"]), "without weakening") {
		t.Fatal("rationale lost")
	}
	kinds := map[string]bool{}
	for _, a := range s.Relationships() {
		kinds[String(a, "projection_path")] = true
	}
	for _, kind := range []string{"TraceableElement.supports", "TraceableElement.derived_from", "Requirement.refines", "TraceableElement.relies_on", "UseCase.involves", "Requirement.applies_to"} {
		if !kinds[kind] {
			t.Fatal("missing projection", kind)
		}
	}
	checkSourceText(t, s)
	forged := cloneSnapshot(t, s)
	designRecord(t, forged, "Requirement", "Preserve document metadata on restoration")["refines"] = []string{String(elementKind(t, forged, "BusinessRule"), "id")}
	requireFinding(t, forged, "shape.invalid")
	borrowed := cloneSnapshot(t, s)
	for _, ref := range borrowed.References {
		if ref["role"] == "refines" {
			ref["context_owner"] = elementKind(t, borrowed, "Job")["id"]
			break
		}
	}
	requireFinding(t, borrowed, "reference.context")
}
func TestTraceTargetErrorsAndDeferredReferences(t *testing.T) {
	for _, tc := range []struct {
		name, path, old, replacement, field string
		code                                int
	}{
		{"definition is not use", "specs/use-cases/restore.use-case.spec.md", "[DocumentTrash](../document-workspace.system.spec.md#documenttrash)", "[Trash](../concepts/trash.concept.spec.md)", "involves", 1},
		{"goal cannot be refined", "specs/requirements/restoration.requirement.spec.md", "[Restore retained content](restore-content.requirement.spec.md)", "[Goal](../document-workspace.system.spec.md#reduce-lost-work)", "refines", 1},
		{"goal is not assumption", "specs/requirements/restoration.requirement.spec.md", "[Trusted time](../document-workspace.system.spec.md#trusted-time)", "[Goal](../document-workspace.system.spec.md#reduce-lost-work)", "relies_on", 1},
		{"concept is not motivation", "specs/jobs/recover.job.spec.md", "[Reduce lost work](../document-workspace.system.spec.md#reduce-lost-work)", "[Trash](../concepts/trash.concept.spec.md)", "supports", 1},
		{"external basis deferred", "specs/requirements/restore-content.requirement.spec.md", "[Recovery window](../business-rules/retention.business-rule.spec.md)", "[Policy](https://example.com/policy)", "derived_from", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := traceFixture(t)
			traceEdit(t, root, tc.path, tc.old, tc.replacement)
			r := conformanceExport(t, root)
			if r.ExitCode != tc.code {
				t.Fatalf("exit %d: %v", r.ExitCode, r.Snapshot.Diagnostics)
			}
			found := false
			for _, f := range r.Snapshot.Diagnostics {
				if String(f, "code") == "relationship."+tc.field {
					found = true
				}
			}
			if !found {
				t.Fatal("missing relationship diagnostic", r.Snapshot.Diagnostics)
			}
			checkSourceText(t, r.Snapshot)
		})
	}
}
func TestRefinementCycleRetainsSourceWithoutCyclicClaim(t *testing.T) {
	for _, self := range []bool{true, false} {
		root := traceFixture(t)
		if self {
			traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "[Restore retained content](restore-content.requirement.spec.md)", "[This requirement](restoration.requirement.spec.md)")
		} else {
			traceEdit(t, root, "specs/requirements/restore-content.requirement.spec.md", "## Derived from", "## Refines\n\n[Metadata](restoration.requirement.spec.md)\n\n## Derived from")
		}
		r := conformanceExport(t, root)
		if r.ExitCode != 1 {
			t.Fatalf("cycle accepted: %v", r.Snapshot.Diagnostics)
		}
		for _, e := range r.Snapshot.Elements {
			if len(Strings(e["refines"])) != 0 {
				t.Fatal("cyclic edge exported", e)
			}
		}
		for _, a := range r.Snapshot.Relationships() {
			if String(a, "projection_path") == "Requirement.refines" {
				t.Fatal("cyclic assertion exported")
			}
		}
		forged := cloneSnapshot(t, r.Snapshot)
		req := designRecord(t, forged, "Requirement", "Preserve document metadata on restoration")
		for _, ref := range forged.References {
			if ref["role"] == "refines" && ref["context_owner"] == req["id"] {
				Object(ref["resolution"])["target"] = req["id"]
				Object(ref["resolution"])["status"] = "resolved"
				break
			}
		}
		requireFinding(t, forged, "relationship.cycle")
		checkSourceText(t, r.Snapshot)
	}
}
func TestTraceCitationsAndDrafts(t *testing.T) {
	root := traceFixture(t)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Refines", "## Related specifications")
	r := conformanceExport(t, root)
	if r.ExitCode != 0 {
		t.Fatal(r.Snapshot.Diagnostics)
	}
	if len(Strings(designRecord(t, r.Snapshot, "Requirement", "Preserve document metadata on restoration")["refines"])) != 0 {
		t.Fatal("citation inferred refinement")
	}
	codes := map[string]bool{}
	for _, f := range r.Snapshot.Diagnostics {
		if String(f, "severity") == "error" {
			t.Fatal(f)
		}
		codes[String(f, "code")] = true
	}
	if !codes["relationship.motivation"] || codes["relationship.system_verification"] {
		t.Fatal("missing draft advice", codes)
	}
	// A selected subtree does not traverse references from its direct targets.
	data, err := os.ReadFile(filepath.Join(root, "specs/requirements/restoration.requirement.spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	selected := strings.NewReplacer("](restore-content.requirement.spec.md)", "](../../../specs/requirements/restore-content.requirement.spec.md)", "../use-cases/", "../../../specs/use-cases/", "../document-workspace.system.spec.md", "../../../specs/document-workspace.system.spec.md").Replace(string(data))
	conformanceWrite(t, root, "selected/specs/requirements/restoration.requirement.spec.md", selected)
	bounded := conformanceExport(t, filepath.Join(root, "selected"))
	found := false
	for _, f := range bounded.Snapshot.Diagnostics {
		found = found || strings.Contains(String(f, "message"), "not_attempted")
	}
	if !found {
		t.Fatal("bounded scope was presented as complete")
	}
}
func TestTraceViewIsPortableAndDoesNotClaimProof(t *testing.T) {
	root := traceFixture(t)
	s := conformanceExport(t, root).Snapshot
	before, _ := json.Marshal(s)
	view, err := TraceMarkdown(s, filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"## Why:", "## Behavior:", "## Change review", "restore-content.requirement.spec.md", "refines", "supports", "Further review candidates"} {
		if !strings.Contains(view, part) {
			t.Fatal("missing", part)
		}
	}
	if strings.Contains(view, root) {
		t.Fatal("absolute source paths leaked into generated view")
	}
	after, _ := json.Marshal(s)
	if string(before) != string(after) {
		t.Fatal("view mutated snapshot")
	}
	second := traceFixture(t)
	other, err := TraceMarkdown(conformanceExport(t, second).Snapshot, filepath.Join(second, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	if view != other {
		t.Fatal("generated output depends on checkout path or snapshot IDs")
	}
	again, err := TraceMarkdown(s, filepath.Join(root, "docs"))
	if err != nil || again != view {
		t.Fatal("nondeterministic view", err)
	}
}
func TestSeveralRefinementParents(t *testing.T) {
	root := traceFixture(t)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Refines", "## Refines\n\n[Policy](../business-rules/retention.business-rule.spec.md)\n")
	r := conformanceExport(t, root)
	if r.ExitCode != 0 {
		t.Fatal(r.Snapshot.Diagnostics)
	}
	if len(Strings(designRecord(t, r.Snapshot, "Requirement", "Preserve document metadata on restoration")["refines"])) != 2 {
		t.Fatal("multiple refinement parents lost")
	}
	requireValid(t, r.Snapshot)
}
