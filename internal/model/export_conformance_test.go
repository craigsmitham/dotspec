package model

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func conformanceTypes() []check.Type {
	return []check.Type{qualitymodel.CheckType(), businessrule.CheckType(), system.CheckType(), system.ComponentCheckType(), system.ProductCheckType(system.FeatureLayout), system.ProductCheckType(system.CapabilityLayout), system.AttributeCheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType()}
}

func conformanceExport(t *testing.T, path string) ExportResult {
	t.Helper()
	result, err := Export(path, conformanceTypes()...)
	if err != nil {
		t.Fatal(err)
	}
	findings, err := Validate(result.Snapshot)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range findings {
		t.Errorf("graph %s %s: %s", finding.Code, finding.Path, finding.Message)
	}
	return result
}

func conformanceWrite(t *testing.T, root, path, body string) string {
	t.Helper()
	path = filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func conformanceIndex(snapshot *Snapshot) map[string]Record {
	index := map[string]Record{}
	for _, collection := range [][]Record{snapshot.Sources, snapshot.Documents, snapshot.ContentNodes, snapshot.Elements, snapshot.References, snapshot.LabelOccurrences, snapshot.Fragments, snapshot.Diagnostics} {
		for _, record := range collection {
			index[String(record, "id")] = record
		}
	}
	return index
}

func checkSourceText(t *testing.T, snapshot *Snapshot) {
	t.Helper()
	for _, source := range snapshot.Sources {
		if String(source, "kind") != "SourceDocument" || String(source, "read_status") != "read" {
			continue
		}
		actual, present := source["markdown"].(string)
		if String(source, "media_type") != "text/markdown" {
			if present || source["payload"] != nil {
				t.Errorf("native reference embedded: %s", source["locator"])
			}
			continue
		}
		original, err := os.ReadFile(filepath.FromSlash(String(source, "locator")))
		if err != nil {
			t.Fatal(err)
		}
		want := strings.ToValidUTF8(strings.ReplaceAll(strings.TrimPrefix(string(original), "\ufeff"), "\r\n", "\n"), "�")
		if !present || actual != want {
			t.Errorf("Markdown mismatch: %s", source["locator"])
		}

	}
}

func TestExportConformanceCheckedInProjects(t *testing.T) {
	for _, dir := range []string{"../..", "../../examples/todo"} {
		result := conformanceExport(t, dir)
		if result.ExitCode != 0 {
			t.Fatal(result.Snapshot.Diagnostics)
		}
		s := result.Snapshot
		checkSourceText(t, s)
		requireValid(t, cloneSnapshot(t, s))
		if len(s.References) == 0 || len(s.Elements) == 0 {
			t.Fatal("missing records", dir)
		}
		data, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		var doc map[string]json.RawMessage
		if err := json.Unmarshal(data, &doc); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"content_nodes", "content_bindings", "coverage", "resolutions", "assertions"} {
			if doc[key] != nil {
				t.Fatal("redundant interchange field", key)
			}
		}
	}
}
func TestExportConformanceClassificationAndDocumentLocalIDs(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/root.system.spec.md", "---\ntype: system-spec\n---\n# Root\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n\n## Quality attributes\n\n")
	conformanceWrite(t, root, "specs/quality/broad/broad.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nid: broad\n---\n# Broad\n\n## Definition\n\nA concern.\n")
	conformanceWrite(t, root, "specs/quality/broad/narrow/narrow.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nid: narrow\n---\n# Narrow\n\n## Definition\n\nA narrower concern.\n")

	conformanceWrite(t, root, "specs/requirements/resolved.requirement.spec.md", "---\ntype: requirement-spec\nid: reused\nquality_attributes: [narrow]\n---\n# Resolved\n")
	conformanceWrite(t, root, "specs/requirements/undefined.requirement.spec.md", "---\ntype: requirement-spec\nid: reused\nquality_attributes: [unknown]\n---\n# Undefined\n")
	conformanceWrite(t, root, "sub/specs/subsystem.system.spec.md", "---\ntype: system-spec\n---\n# Subsystem\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n")
	conformanceWrite(t, root, "sub/specs/requirements/unbound.requirement.spec.md", "---\ntype: requirement-spec\nid: reused\nquality_attributes: [narrow]\n---\n# Unbound\n")
	conformanceWrite(t, root, "empty/specs/empty-attributes.system.spec.md", "---\ntype: system-spec\n---\n# Empty attributes\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n\n## Quality attributes\n")
	conformanceWrite(t, root, "empty/specs/requirements/undefined-empty.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [narrow]\n---\n# Undefined in empty attributes\n")
	conformanceWrite(t, root, "sibling/specs/sibling.system.spec.md", "---\ntype: system-spec\n---\n# Sibling\n## Purpose\nP.\n## Scope\nS.\n## Quality attributes\n")
	conformanceWrite(t, root, "sibling/specs/quality/local-narrow/local-narrow.quality-attribute.spec.md", "---\ntype: quality-attribute-spec\nid: narrow\n---\n# Local narrow\n\n## Definition\n\nA different concern in this system.\n")

	conformanceWrite(t, root, "sibling/specs/requirements/local.requirement.spec.md", "---\ntype: requirement-spec\nquality_attributes: [narrow]\n---\n# Local classification\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Errorf("advice-only fixtures exit %d", result.ExitCode)
	}
	index := conformanceIndex(result.Snapshot)
	statuses := map[string]int{}
	reused := map[string]bool{}
	for _, label := range result.Snapshot.LabelOccurrences {
		if String(label, "role") != "quality_classification" {
			continue
		}
		resolution := Object(label["resolution"])
		statuses[String(resolution, "status")]++
		owner := index[String(label, "owner")]
		if String(resolution, "status") == "resolved" {
			target := index[String(resolution, "target")]
			if String(target, "specification_system") != String(label, "governing_system") || String(target, "specification_system") != String(owner, "specification_system") {
				t.Fatalf("classification escaped its declaring system: %+v", label)
			}
		}
		if String(resolution, "status") == "resolved" && len(Strings(owner["quality_attributes"])) != 1 {
			t.Error("parent category was fabricated as a directly authored classification")
		}
	}
	for _, element := range result.Snapshot.Elements {
		if String(element, "authored_id") == "reused" {
			reused[String(element, "id")] = true
		}
	}
	if len(reused) != 3 || statuses["resolved"] != 2 || statuses["undefined_label"] != 2 || statuses["unbound_label"] != 1 {
		t.Errorf("IDs/classification states: ids=%d states=%v", len(reused), statuses)
	}
	undefinedAdvice := 0
	for _, d := range result.Snapshot.Diagnostics {
		if strings.HasSuffix(String(d, "code"), "undefined-attribute") {
			undefinedAdvice++
			if String(d, "severity") != "advice" {
				t.Error("undefined label became a structural error")
			}
		}
	}
	if undefinedAdvice != 3 {
		t.Errorf("undefined-label advice count %d", undefinedAdvice)
	}
	checkSourceText(t, result.Snapshot)
}

func TestExportConformanceInvalidAndNativeSources(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/jobs/valid.job.spec.md", "---\ntype: job-spec\n---\n# Valid\n\n## Job\n\nInspect from request to report.\n\n## Performer\n\nContributor.\n\n## Related specifications\n\n[Native](../../reference.md#native) and [Broken](../../missing.md).\n")
	badPath := conformanceWrite(t, root, "specs/jobs/bad.job.spec.md", "---\ntype: job-spec\n---\n# Bad\n\n## Intended progress\n\nFirst.\n\n## Intended progress\n\nDuplicate.\n")
	conformanceWrite(t, root, "reference.md", "\xef\xbb\xbf# Native\r\n\r\n<!-- retain comment -->\r\n\r\n## Custom material\r\n\r\nA [second hop](second.md) must remain uninspected.  \r\n")
	second := conformanceWrite(t, root, "second.md", "# Second hop\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 1 {
		t.Errorf("malformed source exit %d", result.ExitCode)
	}
	badID := SourceID(filepath.ToSlash(badPath))
	validRoots := 0
	badDomain := false
	foundNative := false
	for _, element := range result.Snapshot.Elements {
		if String(element, "kind") == "Job" {
			validRoots++
		}
		if String(element, "source_document") == badID {
			badDomain = true
		}
	}
	for _, source := range result.Snapshot.Sources {
		if strings.HasSuffix(String(source, "locator"), "/reference.md") {
			foundNative = true
		}
		if String(source, "locator") == filepath.ToSlash(second) {
			t.Error("reference traversal expanded beyond one edge")
		}
	}
	if validRoots != 1 || badDomain || !foundNative {
		t.Errorf("valid/invalid/native handling roots=%d badDomain=%v native=%v", validRoots, badDomain, foundNative)
	}
	duplicate := false
	for _, d := range result.Snapshot.Diagnostics {
		if String(d, "code") == "structure.duplicate-section" {
			duplicate = true
			if len(records(d["related"])) != 1 {
				t.Errorf("duplicate finding lost earlier source location: %#v", d)
			}
		}
	}
	if !duplicate {
		t.Error("invalid source diagnostic missing")
	}
	checkSourceText(t, result.Snapshot)
}

func TestExportConformanceDeterminismAndContext(t *testing.T) {
	root := t.TempDir()
	context := conformanceWrite(t, root, "specs/context.system.spec.md", "---\ntype: system-spec\n---\n# Context\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n")
	conformanceWrite(t, root, "module/specs/use-cases/example.use-case.spec.md", "---\ntype: use-case-spec\n---\n# Example\n\n## Main success scenario\n\n1. **Actor:** Contributor\n   **Action:** Inspect\n   **Response:** Report\n   **DesiredOutcome:** Informed\n\n## Extensions\n\n### Failure\n\n**Condition:** Unreadable source\n**At step:** 1\n**Response:** Report a failure\n**Continuation:** End\n")
	selected := filepath.Join(root, "module")
	first := conformanceExport(t, selected)
	encoded, err := json.Marshal(first.Snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Snapshot.ContextSources) == 0 {
		t.Error("enclosing context outside selected subtree omitted")
	}
	contextID := SourceID(filepath.ToSlash(context))
	found := false
	for _, source := range first.Snapshot.Sources {
		found = found || String(source, "id") == contextID
	}
	if !found {
		t.Error("required context source not captured")
	}
	for i := 0; i < 8; i++ {
		next := conformanceExport(t, selected)
		again, err := json.Marshal(next.Snapshot)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(encoded, again) {
			t.Fatalf("same-source export changed on iteration %d", i)
		}
	}
}

func TestExportUsesTheSourceAcquiredDuringInspection(t *testing.T) {
	root := t.TempDir()
	const path = "specs/requirements/save.requirement.spec.md"
	const initial = "---\ntype: requirement-spec\n---\n# Save\n## Statement\nSave the record.\n"
	conformanceWrite(t, root, path, initial)
	inspection := check.InspectSources(root, conformanceTypes()...)
	conformanceWrite(t, root, path, strings.ReplaceAll(initial, "Save the record.", "Save its revision."))
	first := ExportInspection(inspection, conformanceTypes()...).Snapshot
	next := conformanceExport(t, root).Snapshot
	if first.Sources[0]["markdown"] != initial || next.Sources[0]["markdown"] == initial {
		t.Fatal("an inspection must use one source revision throughout check and export")
	}
}
