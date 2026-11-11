package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func traceFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../../examples/requirements-traceability")); err != nil {
		t.Fatal(err)
	}
	return root
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
	clause := elementKind(t, s, "ObligationClause")
	uc := elementKind(t, s, "UseCase")
	req := elementKind(t, s, "GeneralRequirement")
	use := elementKind(t, s, "ConceptUse")
	if clause["authored_id"] != "restore-content" || clause["owner"] != uc["id"] {
		t.Fatal(clause)
	}
	for _, tc := range []struct {
		r      Record
		field  string
		target Record
	}{{req, "refines", clause}, {uc, "involves", use}, {clause, "derived_from", elementKind(t, s, "BusinessRule")}, {req, "relies_on", elementKind(t, s, "Assumption")}, {elementKind(t, s, "Job"), "supports", elementKind(t, s, "SystemGoal")}} {
		if !slices.Equal(Strings(tc.r[tc.field]), []string{String(tc.target, "id")}) {
			t.Fatalf("%s: %v", tc.field, tc.r)
		}
		if len(Strings(tc.r[tc.field+"_references"])) != 1 {
			t.Fatal("missing reference provenance", tc.r)
		}
	}
	if String(s.Coverage, "status") != "complete_for_profile" || len(Strings(s.Coverage["extension_sections"])) != 0 {
		t.Fatal(s.Coverage)
	}
	if !strings.Contains(newDesignResolver(s).text(req["refines_content"]), "without weakening") {
		t.Fatal("rationale lost")
	}
	kinds := map[string]bool{}
	for _, a := range s.Assertions {
		kinds[String(a, "kind")] = true
	}
	for _, kind := range []string{"Supports", "DerivedFrom", "Refines", "ReliesOn", "Involves", "HasObligation", "AppliesTo"} {
		if !kinds[kind] {
			t.Fatal("missing projection", kind)
		}
	}
	conformancePayloads(t, s)
	forged := cloneSnapshot(t, s)
	elementKind(t, forged, "GeneralRequirement")["refines"] = []string{String(elementKind(t, forged, "BusinessRule"), "id")}
	requireFinding(t, forged, "concept.projection")
	borrowed := cloneSnapshot(t, s)
	br := elementKind(t, borrowed, "GeneralRequirement")
	br["refines_references"] = br["relies_on_references"]
	requireFinding(t, borrowed, "relationship.provenance")
}
func TestTraceTargetErrorsAndDeferredReferences(t *testing.T) {
	for _, tc := range []struct {
		name, path, old, replacement, field string
		code                                int
	}{
		{"definition is not use", "specs/use-cases/restore.use-case.spec.md", "[DocumentTrash](../system.spec.md#documenttrash)", "[Trash](../concepts/trash.concept.spec.md)", "involves", 1},
		{"goal is not obligation", "specs/requirements/restoration.requirement.spec.md", "[Restore retained content](../use-cases/restore.use-case.spec.md#restore-content)", "[Goal](../system.spec.md#reduce-lost-work)", "refines", 1},
		{"goal is not assumption", "specs/requirements/restoration.requirement.spec.md", "[Trusted time](../system.spec.md#trusted-time)", "[Goal](../system.spec.md#reduce-lost-work)", "relies_on", 1},
		{"concept is not motivation", "specs/jobs/recover.job.spec.md", "[Reduce lost work](../system.spec.md#reduce-lost-work)", "[Trash](../concepts/trash.concept.spec.md)", "supports", 1},
		{"external basis deferred", "specs/use-cases/restore.use-case.spec.md", "[Recovery window](../business-rules/retention.business-rule.spec.md)", "[Policy](https://example.com/policy)", "derived_from", 0},
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
			conformancePayloads(t, r.Snapshot)
		})
	}
}
func TestRefinementCycleRetainsSourceWithoutCyclicClaim(t *testing.T) {
	for _, self := range []bool{true, false} {
		root := traceFixture(t)
		if self {
			traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "[Restore retained content](../use-cases/restore.use-case.spec.md#restore-content)", "[This requirement](restoration.requirement.spec.md)")
		} else {
			traceEdit(t, root, "specs/use-cases/restore.use-case.spec.md", "#### Derived from", "#### Refines\n\n[Metadata](../requirements/restoration.requirement.spec.md)\n\n#### Derived from")
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
		for _, a := range r.Snapshot.Assertions {
			if String(a, "kind") == "Refines" {
				t.Fatal("cyclic assertion exported")
			}
		}
		forged := cloneSnapshot(t, r.Snapshot)
		req := elementKind(t, forged, "GeneralRequirement")
		req["refines"] = []string{String(req, "id")}
		requireFinding(t, forged, "relationship.cycle")
		conformancePayloads(t, r.Snapshot)
	}
}
func TestTraceCitationsAndDrafts(t *testing.T) {
	root := traceFixture(t)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Refines", "## Related specifications")
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Verification plan", "## Notes")
	r := conformanceExport(t, root)
	if r.ExitCode != 0 {
		t.Fatal(r.Snapshot.Diagnostics)
	}
	if len(Strings(elementKind(t, r.Snapshot, "GeneralRequirement")["refines"])) != 0 {
		t.Fatal("citation inferred refinement")
	}
	codes := map[string]bool{}
	for _, f := range r.Snapshot.Diagnostics {
		if String(f, "severity") == "error" {
			t.Fatal(f)
		}
		codes[String(f, "code")] = true
	}
	if !codes["relationship.motivation"] || !codes["relationship.verification_plan"] {
		t.Fatal("missing draft advice", codes)
	}
	// A selected subtree does not traverse references from its direct targets.
	data, err := os.ReadFile(filepath.Join(root, "specs/requirements/restoration.requirement.spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	selected := strings.NewReplacer("../use-cases/", "../../../specs/use-cases/", "../system.spec.md", "../../../specs/system.spec.md").Replace(string(data))
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
	for _, part := range []string{"## Why:", "## Behavior:", "## Change review", "## Verification guidance", "Not evaluated", "#restore-retained-content", "refines", "supports", "Further review candidates"} {
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
func TestSeveralRefinementParentsAndNestedConceptObligation(t *testing.T) {
	root := traceFixture(t)
	traceEdit(t, root, "specs/requirements/restoration.requirement.spec.md", "## Refines", "## Refines\n\n[Policy](../business-rules/retention.business-rule.spec.md)\n")
	traceEdit(t, root, "specs/concepts/trash.concept.spec.md", "### Expire", `#### Obligations
##### Preserve identity
Restoration must preserve the item identity.
###### Derived from
[Recovery](../use-cases/restore.use-case.spec.md#restore-content)
###### Verification plan
Compare the item identity before deletion and after restoration.

### Expire`)
	r := conformanceExport(t, root)
	if r.ExitCode != 0 {
		t.Fatal(r.Snapshot.Diagnostics)
	}
	if len(Strings(elementKind(t, r.Snapshot, "GeneralRequirement")["refines"])) != 2 {
		t.Fatal("multiple refinement parents lost")
	}
	clause := designRecord(t, r.Snapshot, "ObligationClause", "Preserve identity")
	owner := conformanceIndex(r.Snapshot)[String(clause, "owner")]
	if String(owner, "kind") != "ConceptAction" {
		t.Fatal("nested clause misowned", clause)
	}
	requireValid(t, r.Snapshot)
}
