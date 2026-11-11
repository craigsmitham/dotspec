package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func conformanceTypes() []check.Type {
	return []check.Type{qualitymodel.CheckType(), businessrule.CheckType(), system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType(), quality.CheckType()}
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
	for _, collection := range [][]Record{snapshot.Sources, snapshot.Documents, snapshot.ContentNodes, snapshot.Elements, snapshot.References, snapshot.LabelOccurrences, snapshot.Resolutions, snapshot.Fragments, snapshot.Assertions, snapshot.ContentBindings, snapshot.Diagnostics} {
		for _, record := range collection {
			index[String(record, "id")] = record
		}
	}
	return index
}

func conformancePayloads(t *testing.T, snapshot *Snapshot) {
	t.Helper()
	for _, source := range snapshot.Sources {
		if String(source, "kind") != "SourceDocument" || String(source, "read_status") != "read" {
			continue
		}
		payload := Object(source["payload"])
		original, err := os.ReadFile(filepath.FromSlash(String(source, "locator")))
		if err != nil {
			t.Fatal(err)
		}
		recovered, err := base64.StdEncoding.DecodeString(String(payload, "bytes_base64"))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(original, recovered) {
			t.Errorf("source bytes differ: %s", source["locator"])
		}
		digest := sha256.Sum256(recovered)
		if hex.EncodeToString(digest[:]) != String(payload, "sha256") || len(recovered) != Number(payload["byte_length"]) {
			t.Errorf("payload checksum/size mismatch: %s", source["locator"])
		}
	}
}

func TestExportConformanceCheckedInDomainCoverage(t *testing.T) {
	result := conformanceExport(t, "../..")
	if result.ExitCode != 0 {
		t.Errorf("checked-in corpus exit %d", result.ExitCode)
	}
	conformancePayloads(t, result.Snapshot)
	index := conformanceIndex(result.Snapshot)
	docs := map[string]Record{}
	for _, document := range result.Snapshot.Documents {
		docs[String(document, "source")] = document
	}
	bindings := map[string][]Record{}
	for _, binding := range result.Snapshot.ContentBindings {
		bindings[String(binding, "content_node")] = append(bindings[String(binding, "content_node")], binding)
	}
	selectedSources := map[string]bool{}
	expectedRoots := map[string]int{}
	documents, sections := 0, 0
	for _, root := range []string{"../../specs", "../../examples/starter/specs", "../../examples/concept-design/specs", "../../examples/requirements-traceability/specs", "../../internal/integration/specs", "../../internal/qualitymodel/specs"} {
		if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".spec.md") {
				return nil
			}
			absolute, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			absolute = filepath.ToSlash(absolute)
			sourceID := SourceID(absolute)
			selectedSources[sourceID] = true
			document := docs[sourceID]
			if document == nil {
				t.Errorf("source missing from actual Export: %s", path)
				return nil
			}
			coverage := Object(document["coverage"])
			if String(coverage, "status") != "complete_for_profile" || len(Strings(coverage["unbound_content_nodes"])) != 0 || len(Strings(coverage["extension_sections"])) != 0 {
				t.Errorf("%s domain coverage: %#v", path, coverage)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			parsed := ParseContent(absolute, data)
			class := map[string]string{"quality-model": "QualityModel", "business-rule": "BusinessRule", "system": "System", "job": "Job", "use-case": "UseCase", "concept": "Concept", "requirement": "GeneralRequirement", "quality": "QualityRequirement"}[String(parsed.Document, "declared_type")]
			if class == "" {
				t.Errorf("%s: unknown specimen family", path)
			}
			expectedRoots[class]++
			rootCount := 0
			for _, id := range Strings(document["elements"]) {
				if String(index[id], "kind") == class {
					rootCount++
				}
			}
			if rootCount != 1 {
				t.Errorf("%s: exported %d %s roots, want one", path, rootCount, class)
			}
			for _, section := range parsed.Sections {
				if section.Level != 2 {
					continue
				}
				sections++
				mapped := false
				for _, binding := range bindings[section.HeadingID] {
					owner := index[String(binding, "owner")]
					if String(owner, "kind") != "Document" && String(owner, "kind") != "ExtensionSection" && String(binding, "role") != "fallback" && String(binding, "role") != "extension" {
						mapped = true
					}
				}
				if !mapped {
					t.Errorf("%s section %q has no named domain-field binding", path, section.Title)
				}
			}
			documents++
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	if documents < 23 || sections < 150 {
		t.Errorf("inventory unexpectedly shrank: %d documents, %d sections (baseline 23/150)", documents, sections)
	}
	counts := map[string]int{}
	for _, element := range result.Snapshot.Elements {
		if selectedSources[String(element, "source_document")] {
			counts[String(element, "kind")]++
		}
	}
	for kind, want := range expectedRoots {
		if counts[kind] != want {
			t.Errorf("%s exported roots: got %d want %d inventoried sources", kind, counts[kind], want)
		}
	}
	for kind, want := range map[string]int{"JobStep": 16, "DesiredOutcome": 21, "ProgressCase": 1, "ConceptAction": 6, "GlossaryTerm": 3, "QualityAttribute": 293} {
		if counts[kind] < want {
			t.Errorf("%s domain records: got %d want at least baseline %d", kind, counts[kind], want)
		}
	}
	if len(result.Snapshot.ContextSources) == 0 || len(result.Snapshot.Resolutions) == 0 {
		t.Error("actual export did not inspect context and direct references")
	}
	encoded, err := json.Marshal(result.Snapshot)
	if err != nil {
		t.Fatal(err)
	}
	var roundtrip Snapshot
	if err := json.Unmarshal(encoded, &roundtrip); err != nil {
		t.Fatal(err)
	}
	again, err := json.Marshal(&roundtrip)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, again) {
		t.Error("snapshot JSON roundtrip loses subtype or rich-content fields")
	}
	findings, err := Validate(&roundtrip)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) > 0 {
		t.Errorf("roundtripped snapshot invalid: %#v", findings)
	}
}

func TestExportConformanceClassificationAndDocumentLocalIDs(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/system.spec.md", "---\ntype: system\n---\n# Root\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n\n## Quality attributes\n\n### Broad\n\n**ID:** broad\n\nA concern.\n\n### Narrow\n\n**ID:** narrow\n**Sub-attribute of:** broad\n\nA narrower concern.\n")
	conformanceWrite(t, root, "specs/quality/resolved.quality.spec.md", "---\ntype: quality\nid: reused\nquality_attributes: [narrow]\n---\n# Resolved\n")
	conformanceWrite(t, root, "specs/requirements/undefined.requirement.spec.md", "---\ntype: requirement\nid: reused\nquality_attributes: [unknown]\n---\n# Undefined\n")
	conformanceWrite(t, root, "sub/specs/system.spec.md", "---\ntype: system\n---\n# Subsystem\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n")
	conformanceWrite(t, root, "sub/specs/quality/unbound.quality.spec.md", "---\ntype: quality\nid: reused\nquality_attributes: [narrow]\n---\n# Unbound\n")
	conformanceWrite(t, root, "empty/specs/system.spec.md", "---\ntype: system\n---\n# Empty attributes\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n\n## Quality attributes\n")
	conformanceWrite(t, root, "empty/specs/quality/undefined-empty.quality.spec.md", "---\ntype: quality\nquality_attributes: [narrow]\n---\n# Undefined in empty attributes\n")
	conformanceWrite(t, root, "sibling/specs/system.spec.md", "---\ntype: system\n---\n# Sibling\n## Purpose\nP.\n## Scope\nS.\n## Quality attributes\n### Local narrow\n**ID:** narrow\nA different concern in this system.\n")
	conformanceWrite(t, root, "sibling/specs/quality/local.quality.spec.md", "---\ntype: quality\nquality_attributes: [narrow]\n---\n# Local classification\n")
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
		resolution := index[String(label, "resolution")]
		statuses[String(resolution, "status")]++
		owner := index[String(label, "owner")]
		if String(resolution, "status") == "resolved" {
			target := index[String(resolution, "target")]
			if String(target, "owner") != String(label, "system_context") || String(target, "owner") != String(owner, "system_context") {
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
	if undefinedAdvice != 2 {
		t.Errorf("undefined-label advice count %d", undefinedAdvice)
	}
	conformancePayloads(t, result.Snapshot)
}

func TestExportConformanceInvalidAndNativeSources(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/jobs/valid.job.spec.md", "---\ntype: job\n---\n# Valid\n\n## Job\n\nInspect from request to report.\n\n## Performer\n\nContributor.\n\n## Related specifications\n\n[Native](../../reference.md#native) and [Broken](../../missing.md).\n")
	badPath := conformanceWrite(t, root, "specs/jobs/bad.job.spec.md", "---\ntype: job\n---\n# Bad\n\n## Intended progress\n\nFirst.\n\n## Desired outcome\n\nDuplicate.\n")
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
	if String(result.Snapshot.Coverage, "status") != "partial" {
		t.Error("invalid domain extraction reported complete coverage")
	}
	conformancePayloads(t, result.Snapshot)
}

func TestExportConformanceDeterminismAndContext(t *testing.T) {
	root := t.TempDir()
	context := conformanceWrite(t, root, "specs/system.spec.md", "---\ntype: system\n---\n# Context\n\n## Purpose\n\nP.\n\n## Scope\n\nS.\n")
	conformanceWrite(t, root, "module/specs/use-cases/example.use-case.spec.md", "---\ntype: use-case\n---\n# Example\n\n## Main success scenario\n\n1. **Actor:** Contributor\n   **Action:** Inspect\n   **Response:** Report\n   **DesiredOutcome:** Informed\n\n## Extensions\n\n### Failure\n\n**Condition:** Unreadable source\n**At step:** 1\n**Response:** Report a failure\n**Continuation:** End\n")
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

func conformanceRichText(snapshot *Snapshot, value any) string {
	index := conformanceIndex(snapshot)
	var out strings.Builder
	for _, id := range Strings(Object(value)["content_nodes"]) {
		span := Object(index[id]["span"])
		payload := Object(index[String(span, "source")]["payload"])
		data, _ := base64.StdEncoding.DecodeString(String(payload, "bytes_base64"))
		a, b := Number(span["start"]), Number(span["end"])
		if a >= 0 && b >= a && b <= len(data) {
			out.Write(data[a:b])
		}
	}
	return out.String()
}

func TestExportConformanceAdversarialDomainCompletion(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/system.spec.md", `---
type: system
---
# System
## Purpose
P.
## Scope
S.
## Glossary
### Term
A definition.
#### Aliases
- First alias
#### Alternative names
- Second alias
## Quality attributes
### Concern
**ID:** concern
A definition.
#### Definition example
Ordinary nested definition prose remains part of the definition.
#### Open question
Which threshold?
##### Evidence needed
Measure the baseline.
##### Decision caveat
An unmapped structured decision detail.
`)
	conformanceWrite(t, root, "specs/jobs/inspect.job.spec.md", `---
type: job
---
# Inspect
## Outcomes
- An actual outcome.
### Supporting notes
- A note, not another outcome.
## Steps
### Review
**Stage:**
#### Surprises
Unmapped step detail.
## Progress cases
### Adoption
#### Timeline
First timeline.
#### Timeline
Second timeline.
#### Forces
- **Push:** First pressure.
#### Forces
- **Pull:** Second attraction.
#### Custom case detail
Unmapped case detail.
## Open question
Which option?
### Evidence needed
Collect evidence.
`)
	conformanceWrite(t, root, "specs/concepts/action.concept.spec.md", `---
type: &family concept
id: *family
---
# Action
## Actions
### Act
#### Inputs
#### Implementation details
Unmapped action detail.
`)
	conformanceWrite(t, root, "specs/use-cases/inspect.use-case.spec.md", `---
type: use-case
---
# Inspect
## Main success scenario
1. **Actor:**
   **Action:**
   **Response:**
## Extensions
### Failure
**Condition:**
**Response:**
#### Flow
1. Report a failure.
##### Flow caveat
1. This explanatory item is not a flow step.
`)
	conformanceWrite(t, root, "specs/quality/measure.quality.spec.md", `---
type: quality
---
# Measure
## Criterion
A qualitative criterion.
## Measurement plan
Measure the criterion.
### Custom protocol
Unmapped protocol detail.
`)
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Errorf("accepted domain fixtures exit %d", result.ExitCode)
	}
	s := result.Snapshot
	extensions := map[string]string{}
	counts := map[string]int{}
	index := conformanceIndex(s)
	for _, element := range s.Elements {
		kind := String(element, "kind")
		counts[kind]++
		switch kind {
		case "ExtensionSection":
			extensions[strings.TrimSpace(conformanceRichText(s, element["heading"]))] = String(index[String(element, "owner")], "kind")
		case "Job":
			decisions := records(Object(element["open_questions"])["questions"])
			if len(decisions) != 1 || !strings.Contains(conformanceRichText(s, decisions[0]["question"]), "Which option?") || !strings.Contains(conformanceRichText(s, decisions[0]["evidence_needed"]), "Collect evidence") {
				t.Errorf("direct decision/evidence misbound: %#v", element["open_questions"])
			}
		case "JobStep":
			if Object(element["stage"])["presence"] != "empty" {
				t.Error("empty stage confused with absent/authored value")
			}
		case "ProgressCase":
			text := conformanceRichText(s, Object(element["timeline"])["content"])
			if !strings.Contains(text, "First timeline") || !strings.Contains(text, "Second timeline") {
				t.Errorf("repeated timeline overwritten: %q", text)
			}
			if len(records(Object(element["forces"])["items"])) != 2 {
				t.Error("repeated force collection overwritten")
			}
		case "Concept":
			if String(element, "authored_id") != "concept" || String(element, "source_fragment") == "" {
				t.Error("YAML alias identity lacks its authored fragment")
			}
			bound := false
			for _, binding := range s.ContentBindings {
				if String(binding, "owner") == String(element, "id") && String(binding, "field_path") == "authored_id" {
					bound = true
				}
			}
			if !bound {
				t.Error("YAML alias identity lacks a named binding")
			}
		case "ConceptAction":
			if Object(element["inputs"])["presence"] != "empty" {
				t.Error("empty action inputs lost")
			}
		case "UseCaseExtension":
			if Object(element["condition"])["presence"] != "empty" {
				t.Error("empty extension condition lost")
			}
			if !strings.Contains(conformanceRichText(s, element["handling"]), "Report a failure.") {
				t.Error("flow narrative missing from handling")
			}
			if !strings.Contains(conformanceRichText(s, element["content"]), "**Response:**") {
				t.Error("empty response field lost from source narrative")
			}
		case "ScenarioStep":
			if String(index[String(element, "owner")], "kind") == "Scenario" {
				text := conformanceRichText(s, element["content"])
				for _, label := range []string{"**Actor:**", "**Action:**", "**Response:**"} {
					if !strings.Contains(text, label) {
						t.Errorf("empty step field %s lost from narrative", label)
					}
				}
			}

		case "GlossaryTerm":
			if len(records(element["alternative_names"])) != 2 {
				t.Error("alias sections overwritten")
			}
		case "QualityAttribute":
			if !strings.Contains(conformanceRichText(s, element["definition"]), "Ordinary nested definition prose") {
				t.Error("nested quality definition prose lost")
			}
		}
	}
	for heading, owner := range map[string]string{"#### Surprises": "JobStep", "#### Custom case detail": "ProgressCase", "#### Implementation details": "ConceptAction", "##### Flow caveat": "UseCaseExtension", "### Supporting notes": "Job", "### Custom protocol": "QualityRequirement", "##### Decision caveat": "QualityAttribute"} {
		if extensions[heading] != owner {
			t.Errorf("unknown section %q owner=%q want %q; extensions=%v", heading, extensions[heading], owner, extensions)
		}
	}
	if counts["DesiredOutcome"] != 1 || counts["ScenarioStep"] != 2 {
		t.Errorf("supporting lists fabricated domain items: outcomes=%d steps=%d", counts["DesiredOutcome"], counts["ScenarioStep"])
	}
	if len(Strings(s.Coverage["extension_sections"])) != len(extensions) {
		t.Error("unknown nested sections missing from coverage inventory")
	}
}

func TestExportUseCaseNarrativesAndExplicitNesting(t *testing.T) {
	root := t.TempDir()
	conformanceWrite(t, root, "specs/use-cases/reserve.use-case.spec.md", `---
type: use-case
---
# Reserve equipment
## Main success scenario
A contributor requests equipment and the system confirms the reservation.
## Extensions
- 1a. Equipment unavailable:
  - 1a1. Offer an alternative.
    - 1a1a. Alternative unavailable:
      - .1. Ask the contributor to return later.
      - .1a. Contributor cannot return:
        Handling unresolved.
- *a. Contributor cancels.

### Service unavailable
**Continuation:** End without a reservation.
**Response:** Explain the service outage.
**At step:** At any time.
#### Handling
1. Preserve the request.
2. Consult [the guide](../../guide.md).
##### 2a. Guide unavailable
Ask a colleague.
#### Flow
1. Report the result.
#### Unresolved handling
Can another service accept the request?
#### Notes
This explanatory heading is not an alternative path.
#### Extensions
##### Colleague unavailable
Ask the contributor to return later.
`)
	conformanceWrite(t, root, "guide.md", "# Guide\n\nTry another service.\n")
	result := conformanceExport(t, root)
	if result.ExitCode != 0 {
		t.Fatalf("export exit %d", result.ExitCode)
	}
	s := result.Snapshot
	conformancePayloads(t, s)
	index := conformanceIndex(s)
	labels := map[string]Record{}
	var useCase, service, notes Record
	counts := map[string]int{}
	for _, e := range s.Elements {
		kind := String(e, "kind")
		counts[kind]++
		switch kind {
		case "UseCase":
			useCase = e
		case "Scenario":
			if len(Strings(e["steps"])) != 0 {
				t.Error("paragraph scenario acquired invented steps")
			}
		case "UseCaseExtension":
			labels[String(e, "authored_label")] = e
			if strings.Contains(conformanceRichText(s, e["name"]), "Service unavailable") {
				service = e
			}
			for _, removed := range []string{"at_step", "response", "flow", "continuation", "handling_state", "governing_references"} {
				if e[removed] != nil {
					t.Errorf("obsolete extension field %s", removed)
				}
			}
		case "ScenarioStep":
			for _, removed := range []string{"actor", "action", "response", "outcome", "governing_references"} {
				if e[removed] != nil {
					t.Errorf("obsolete step field %s", removed)
				}
			}
		case "ExtensionSection":
			if strings.Contains(conformanceRichText(s, e["heading"]), "Notes") {
				notes = e
			}
		}
	}
	if counts["UseCaseExtension"] != 7 || counts["ScenarioStep"] != 5 || counts["ExtensionSection"] != 1 {
		t.Fatalf("unexpected decomposition: %v", counts)
	}
	if service == nil || notes == nil || String(notes, "owner") != String(service, "id") {
		t.Fatal("Notes should belong to its enclosing branch as an extension section")
	}
	rootBranches := Strings(Object(useCase["extensions"])["branches"])
	if len(rootBranches) != 3 || rootBranches[0] != String(labels["1a"], "id") || rootBranches[1] != String(labels["*a"], "id") || rootBranches[2] != String(service, "id") {
		t.Fatalf("branch source order lost: %v", rootBranches)
	}
	if String(labels["1a1a"], "owner") != String(labels["1a"], "id") || String(labels[".1a"], "owner") != String(labels["1a1a"], "id") {
		t.Error("nested labels should follow Markdown containment, including branches under steps")
	}
	for label, want := range map[string]string{"1a": "1a1", "1a1a": ".1"} {
		steps := Strings(labels[label]["steps"])
		if len(steps) != 1 || String(index[steps[0]], "authored_label") != want {
			t.Errorf("branch %s handling step label lost: %v", label, steps)
		}
	}
	if labels["*a"]["handling"] != nil || labels["*a"]["condition"] != nil || labels["*a"]["location"] != nil {
		t.Error("condition-only prose acquired inferred semantic fields")
	}
	if len(Strings(service["branches"])) != 2 || len(Strings(service["steps"])) != 3 {
		t.Error("explicit nested heading/grouping or handling steps missing")
	}
	handling := conformanceRichText(s, service["handling"])
	previous := -1
	for _, part := range []string{"End without a reservation.", "Explain the service outage.", "Preserve the request.", "Report the result.", "Unresolved handling", "Can another service accept the request?"} {
		position := strings.Index(handling, part)
		if position <= previous {
			t.Errorf("handling absent or out of source order: %q in %q", part, handling)
		}
		previous = position
	}
	if !strings.Contains(conformanceRichText(s, service["location"]), "At any time.") {
		t.Error("free-prose location lost")
	}
	if len(Strings(service["references"])) != 1 {
		t.Error("ordinary inline reference missing")
	}
}
