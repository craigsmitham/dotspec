package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func validationSnapshot(t *testing.T) *Snapshot {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"specs/system.spec.md": `---
type: system
---
# Example
## Purpose
Support explicit source-backed specifications.
## Scope
These examples.
## Goals
- Preserve author intent.
## Glossary
### Source
The original authored bytes.
## Quality attributes
### Reliability
**ID:** reliable
Works as specified.
### Consistency
**ID:** consistent
**Sub-attribute of:** reliable
Produces the same result from the same source.
`,
		"specs/use-cases/check.use-case.spec.md": `---
type: use-case
id: check
---
# Check specifications
## Goal
Inspect specifications.
## System boundary
The specification CLI.
## Primary actor
Developer.
## Main success scenario
1. The developer requests inspection.
2. The system returns findings.
## Extensions
### Invalid source
#### Condition
An input is invalid.
#### Handling
Unresolved
## Related specifications
[The system](../system.spec.md) and [requirements](../requirements/).
`,
		"specs/jobs/inspect.job.spec.md": `---
type: job
id: inspect
---
# Inspect intent
## Job
When making a change, I want to inspect intent, so I can preserve it.
## Performer
Developer.
## Situation
Preparing a change.
## Motivation
Understand constraints.
## Intended progress
An informed decision.
## Steps
### Gather intent
**ID:** gather
**Stage:** custom
Read the specifications.
#### Outcomes
- **clear:** Minimize overlooked constraints.
`,
		"specs/concepts/inspection.concept.spec.md": `---
type: concept
---
# Inspection
## Purpose
Expose authored content.
## State
A set of source files.
## Actions
### Inspect
Return the authored content and findings.
## Operational principle
Inspection preserves original sources.
## Additional
Keep this authored content even though its domain role is unknown.
`,
		"specs/requirements/preserve.requirement.spec.md": `---
type: requirement
attributes:
- reliable
---
# Preserve source
## Statement
Inspection preserves the original source bytes.
## Applies to
Source inspection.
## Verification
Compare the bytes before and after inspection.
`,
	}
	for name, text := range files {
		data := []byte(text)
		if strings.Contains(name, "use-case") {
			data = append([]byte{0xef, 0xbb, 0xbf}, bytes.ReplaceAll(data, []byte("\n"), []byte("\r\n"))...)
		}
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Export(dir, system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), quality.CheckType(), requirement.CheckType())
	if err != nil {
		t.Fatal(err)
	}
	requireValid(t, result.Snapshot)
	return result.Snapshot
}

func requireValid(t *testing.T, s *Snapshot) {
	t.Helper()
	findings, err := Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		t.Errorf("%s %s: %s", f.Code, f.Path, f.Message)
	}
	if len(findings) > 0 {
		t.FailNow()
	}
}
func requireFinding(t *testing.T, s *Snapshot, code string) {
	t.Helper()
	findings, err := Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		if f.Code == code {
			return
		}
	}
	t.Fatalf("expected %s, got %+v", code, findings)
}
func cloneSnapshot(t *testing.T, s *Snapshot) *Snapshot {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var clone Snapshot
	if err := json.Unmarshal(data, &clone); err != nil {
		t.Fatal(err)
	}
	return &clone
}
func elementKind(t *testing.T, s *Snapshot, kind string) Record {
	t.Helper()
	for _, r := range s.Elements {
		if String(r, "kind") == kind {
			return r
		}
	}
	t.Fatalf("missing %s", kind)
	return nil
}
func snapshotRecord(t *testing.T, s *Snapshot, id string) Record {
	t.Helper()
	for _, collection := range [][]Record{s.Sources, s.Documents, s.ContentNodes, s.Elements, s.Fragments, s.References, s.LabelOccurrences, s.Resolutions, s.Assertions, s.ContentBindings, s.Diagnostics} {
		for _, r := range collection {
			if String(r, "id") == id {
				return r
			}
		}
	}
	t.Fatalf("missing record %q", id)
	return nil
}

func TestValidationAcceptsCompleteSourceBackedExport(t *testing.T) {
	snapshot := validationSnapshot(t)
	if len(snapshot.Assertions) == 0 || len(snapshot.Resolutions) == 0 {
		t.Fatal("fixture did not exercise materialized assertions and resolution")
	}
	system := elementKind(t, snapshot, "System")
	attributeIDs := Strings(system["quality_attributes"])
	if len(attributeIDs) != 2 {
		t.Fatalf("system attributes: %v", attributeIDs)
	}
	for _, id := range attributeIDs {
		if String(snapshotRecord(t, snapshot, id), "owner") != String(system, "id") {
			t.Fatal("attribute is not owned directly by its system")
		}
	}
	declarations := 0
	for _, assertion := range snapshot.Assertions {
		if String(assertion, "kind") == "DeclaresQualityAttribute" {
			declarations++
			if String(assertion, "subject") != String(system, "id") || !slices.Contains(attributeIDs, String(assertion, "object")) {
				t.Fatalf("incorrect direct declaration: %+v", assertion)
			}
		}
	}
	if declarations != len(attributeIDs) {
		t.Fatalf("declaration count: %d", declarations)
	}
	for _, label := range snapshot.LabelOccurrences {
		if String(label, "system_context") != String(system, "id") {
			t.Fatalf("label lacks declaring system context: %+v", label)
		}
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	findings, err := ValidateJSON(data)
	if err != nil || len(findings) != 0 {
		t.Fatalf("round-trip validation: %v, %+v", err, findings)
	}
}

func TestValidationRejectsCorruptGraph(t *testing.T) {
	original := validationSnapshot(t)
	cases := []struct {
		name, code string
		mutate     func(*Snapshot)
	}{
		{"attribute declaration without section content", "quality.declaration", func(s *Snapshot) { delete(elementKind(t, s, "System"), "quality_attributes_content") }},
		{"section content without attribute declaration", "quality.declaration", func(s *Snapshot) { delete(elementKind(t, s, "System"), "quality_attributes") }},
		{"duplicate identity", "identity.duplicate", func(s *Snapshot) { s.Elements = append(s.Elements, s.Elements[0]) }},
		{"missing typed reference", "reference.missing", func(s *Snapshot) { elementKind(t, s, "UseCase")["main_success_scenario"] = "missing" }},
		{"wrong reference subtype", "reference.kind", func(s *Snapshot) {
			elementKind(t, s, "UseCase")["main_success_scenario"] = elementKind(t, s, "Concept")["id"]
		}},
		{"any-of owner narrowing", "reference.kind", func(s *Snapshot) { elementKind(t, s, "ScenarioStep")["owner"] = elementKind(t, s, "Job")["id"] }},
		{"missing inverse membership", "membership.missing", func(s *Snapshot) { delete(elementKind(t, s, "UseCase"), "main_success_scenario") }},
		{"duplicate membership", "membership.duplicate", func(s *Snapshot) {
			r := elementKind(t, s, "Scenario")
			r["steps"] = append(Strings(r["steps"]), Strings(r["steps"])[0])
		}},
		{"authored identity scope", "identity.authored", func(s *Snapshot) {
			elementKind(t, s, "UseCase")["authored_id"] = "duplicate"
			elementKind(t, s, "Scenario")["authored_id"] = "duplicate"
		}},
		{"unjustified classification", "label.projection", func(s *Snapshot) { elementKind(t, s, "GeneralRequirement")["attributes"] = []string{} }},
		{"unbound quality attributes", "label.unbound", func(s *Snapshot) {
			delete(elementKind(t, s, "GeneralRequirement"), "owning_system")
			for _, label := range s.LabelOccurrences {
				if String(label, "role") == "quality_classification" {
					delete(label, "system_context")
				}
			}
		}},
		{"mismatched resolution backpointer", "resolution.reciprocal", func(s *Snapshot) { s.Resolutions[0]["occurrence"] = s.Resolutions[1]["occurrence"] }},
		{"target on unresolved reference", "resolution.state", func(s *Snapshot) {
			s.Resolutions[0]["status"] = "unresolved"
			s.Resolutions[0]["target"] = s.Sources[0]["id"]
		}},
		{"assertion field disagreement", "assertion.projection", func(s *Snapshot) {
			for _, r := range s.Assertions {
				if String(r, "kind") == "HasScenario" {
					r["object"] = elementKind(t, s, "UseCase")["id"]
					return
				}
			}
		}},
		{"missing derivation input", "assertion.derivation", func(s *Snapshot) { s.Assertions[0]["inputs"] = []string{} }},
		{"missing assertion provenance", "assertion.provenance", func(s *Snapshot) {
			s.Assertions[0]["inputs"] = []string{}
			delete(s.Assertions[0], "source_nodes")
			delete(s.Assertions[0], "source_spans")
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) { s := cloneSnapshot(t, original); test.mutate(s); requireFinding(t, s, test.code) })
	}
}

func TestValidationRejectsCorruptSourceAndCoverage(t *testing.T) {
	original := validationSnapshot(t)
	usecaseSource := String(elementKind(t, original, "UseCase"), "source_document")
	documentID := ""
	for _, d := range original.Documents {
		if String(d, "source") == usecaseSource {
			documentID = String(d, "id")
		}
	}
	cases := []struct {
		name, code string
		mutate     func(*Snapshot)
	}{
		{"source digest", "payload.digest", func(s *Snapshot) {
			Object(snapshotRecord(t, s, usecaseSource)["payload"])["sha256"] = strings.Repeat("0", 64)
		}},
		{"source byte count", "payload.length", func(s *Snapshot) {
			p := Object(snapshotRecord(t, s, usecaseSource)["payload"])
			p["byte_length"] = Number(p["byte_length"]) + 1
		}},
		{"invalid base64", "payload.encoding", func(s *Snapshot) {
			Object(snapshotRecord(t, s, usecaseSource)["payload"])["bytes_base64"] = "not base64!"
		}},
		{"failed read with bytes", "payload.read-state", func(s *Snapshot) { snapshotRecord(t, s, usecaseSource)["read_status"] = "failed" }},
		{"missing offset map", "source.offset-map", func(s *Snapshot) { delete(snapshotRecord(t, s, documentID), "normalization") }},
		{"wrong normalization transform", "source.offset-rule", func(s *Snapshot) {
			maps := records(snapshotRecord(t, s, documentID)["normalization"])
			maps[0]["transformation"] = "CRLF to LF"
		}},
		{"root span omits bytes", "content.root-span", func(s *Snapshot) {
			d := snapshotRecord(t, s, documentID)
			r := snapshotRecord(t, s, String(d, "root"))
			Object(r["span"])["start"] = 1
		}},
		{"duplicate syntax parent", "content.parent", func(s *Snapshot) {
			d := snapshotRecord(t, s, documentID)
			r := snapshotRecord(t, s, String(d, "root"))
			r["children"] = append(Strings(r["children"]), Strings(r["children"])[0])
		}},
		{"syntax cycle", "content.cycle", func(s *Snapshot) {
			d := snapshotRecord(t, s, documentID)
			r := snapshotRecord(t, s, String(d, "root"))
			r["children"] = append(Strings(r["children"]), String(r, "id"))
		}},
		{"orphan syntax", "content.orphan", func(s *Snapshot) {
			d := snapshotRecord(t, s, documentID)
			r := snapshotRecord(t, s, String(d, "root"))
			r["children"] = Strings(r["children"])[1:]
		}},
		{"invalid binding field", "binding.field", func(s *Snapshot) { s.ContentBindings[0]["field_path"] = "does_not_exist" }},
		{"cross-source binding", "binding.source", func(s *Snapshot) {
			s.ContentBindings[0]["owner"] = elementKind(t, s, "UseCase")["id"]
			s.ContentBindings[0]["field_path"] = "goal"
			if v := String(snapshotRecord(t, s, String(s.ContentBindings[0], "content_node")), "document"); v == documentID {
				s.ContentBindings[0]["owner"] = elementKind(t, s, "Concept")["id"]
				s.ContentBindings[0]["field_path"] = "purpose"
			}
		}},
		{"rich presence mismatch", "content.presence", func(s *Snapshot) { Object(elementKind(t, s, "UseCase")["goal"])["presence"] = "empty" }},
		{"rich content source mismatch", "content.ownership", func(s *Snapshot) {
			Object(elementKind(t, s, "UseCase")["goal"])["content_nodes"] = Object(elementKind(t, s, "Concept")["purpose"])["content_nodes"]
		}},
		{"coverage profile mismatch", "coverage.profile", func(s *Snapshot) { s.Coverage["mapping_profile"] = "different/1" }},
		{"coverage false inventory", "coverage.accounting", func(s *Snapshot) {
			Object(snapshotRecord(t, s, documentID)["coverage"])["unbound_content_nodes"] = []string{String(snapshotRecord(t, s, documentID), "root")}
		}},
		{"aggregate inventory mismatch", "coverage.aggregate", func(s *Snapshot) {
			s.Coverage["unbound_content_nodes"] = []string{String(snapshotRecord(t, s, documentID), "root")}
		}},
		{"root structural subtree laundering", "coverage.incomplete", func(s *Snapshot) {
			root := String(snapshotRecord(t, s, documentID), "root")
			kept := []Record{}
			for _, b := range s.ContentBindings {
				node := snapshotRecord(t, s, String(b, "content_node"))
				if String(node, "document") != documentID {
					kept = append(kept, b)
				}
			}
			kept = append(kept, Record{"id": "laundered", "kind": "ContentBinding", "content_node": root, "scope": "subtree", "owner": documentID, "field_path": "root", "role": "structural_wrapper", "mapping_rule": "bad/1"})
			s.ContentBindings = kept
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) { s := cloneSnapshot(t, original); test.mutate(s); requireFinding(t, s, test.code) })
	}
}

func TestValidationAllowsFaithfulFailuresAndUnboundLabels(t *testing.T) {
	original := validationSnapshot(t)
	t.Run("source diagnostic does not invalidate snapshot", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		s.Diagnostics = append(s.Diagnostics, Record{"id": "source-error", "kind": "Diagnostic", "code": "source.invalid", "severity": "error", "message": "Authored source is invalid.", "source": elementKind(t, s, "UseCase")["source_document"], "phase": "parse", "element": "Example", "related": []Record{{"locator": "not-in-this-snapshot.md", "line": 2, "element": "Related"}}})
		requireValid(t, s)
	})
	t.Run("unread source remains representable", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		s.Sources = append(s.Sources, Record{"id": "unread", "kind": "SourceDocument", "locator": "unread.md", "read_status": "failed"})
		s.Coverage["status"] = "partial"
		s.Scope["reads_complete"] = false
		s.Scope["inspection_complete"] = false
		requireValid(t, s)
	})
	t.Run("external payload can be checked offline", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		p := Object(snapshotRecord(t, s, String(elementKind(t, s, "UseCase"), "source_document"))["payload"])
		p["storage"] = "external"
		p["immutable_locator"] = "urn:sha256:" + String(p, "sha256")
		delete(p, "bytes_base64")
		requireValid(t, s)
	})
	t.Run("unbound classification is not an invalid model", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		r := elementKind(t, s, "GeneralRequirement")
		delete(r, "owning_system")
		delete(r, "attributes")
		for _, label := range s.LabelOccurrences {
			if String(label, "owner") == String(r, "id") {
				delete(label, "system_context")
				result := snapshotRecord(t, s, String(label, "resolution"))
				result["status"] = "unbound_label"
				delete(result, "target")
			}
		}
		s.Assertions = slices.DeleteFunc(s.Assertions, func(a Record) bool { return String(a, "subject") == String(r, "id") })
		requireValid(t, s)
	})
}

func TestValidationClosedShapesAndMalformedJSON(t *testing.T) {
	snapshot := validationSnapshot(t)
	elementKind(t, snapshot, "UseCase")["unexpected"] = true
	requireFinding(t, snapshot, "shape.invalid")
	if findings, err := ValidateJSON([]byte(`{"model_version":`)); err != nil || len(findings) != 1 || findings[0].Code != "shape.json" {
		t.Fatalf("invalid JSON must be a validation finding, got %v %+v", err, findings)
	}
}

func TestValidationPayloadRetainsExactBytes(t *testing.T) {
	snapshot := validationSnapshot(t)
	source := snapshotRecord(t, snapshot, String(elementKind(t, snapshot, "UseCase"), "source_document"))
	data, err := base64.StdEncoding.DecodeString(String(Object(source["payload"]), "bytes_base64"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) || !bytes.Contains(data, []byte("\r\n")) {
		t.Fatal("fixture lost original BOM or CRLF")
	}
	requireValid(t, snapshot)
}

func TestValidationRejectsDomainCyclesAndAttributeScopeEscape(t *testing.T) {
	original := validationSnapshot(t)
	t.Run("quality hierarchy cycle", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		reliable, consistent := Record(nil), Record(nil)
		for _, r := range s.Elements {
			if String(r, "kind") == "QualityAttribute" {
				if String(r, "authored_id") == "reliable" {
					reliable = r
				} else {
					consistent = r
				}
			}
		}
		reliable["subattribute_of"] = []string{String(consistent, "id")}
		reliable["subattribute_of_labels"] = []string{"cycle-label"}
		s.LabelOccurrences = append(s.LabelOccurrences, Record{"id": "cycle-label", "kind": "LabelOccurrence", "owner": reliable["id"], "role": "subattribute_of", "label": consistent["authored_id"], "content_node": Strings(Object(consistent["definition"])["content_nodes"])[0], "resolution": "cycle-resolution"})
		s.Resolutions = append(s.Resolutions, Record{"id": "cycle-resolution", "kind": "Resolution", "occurrence": "cycle-label", "status": "resolved", "target": consistent["id"], "rule": "explicit-test/1"})
		requireFinding(t, s, "quality.subattribute-cycle")
	})
	t.Run("cross-system parent", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		var broad Record
		for _, r := range s.Elements {
			if String(r, "authored_id") == "reliable" {
				broad = r
			}
		}
		old := snapshotRecord(t, s, String(broad, "owner"))
		old["quality_attributes"] = slices.DeleteFunc(Strings(old["quality_attributes"]), func(id string) bool { return id == String(broad, "id") })
		s.Elements = append(s.Elements, Record{"id": "other-system", "kind": "System", "source_document": broad["source_document"], "decomposition": "unsegmented", "quality_attributes_content": old["quality_attributes_content"], "quality_attributes": []string{String(broad, "id")}})
		broad["owner"] = "other-system"
		requireFinding(t, s, "quality.scope")
	})
	t.Run("nested branch cycle", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		branch := elementKind(t, s, "Extension")
		second := Record{"id": "second-branch", "kind": "Extension", "source_document": branch["source_document"], "decomposition": "unsegmented", "owner": branch["id"], "content": branch["content"], "branches": []string{String(branch, "id")}}
		branch["owner"] = "second-branch"
		branch["branches"] = []string{"second-branch"}
		s.Elements = append(s.Elements, second)
		Object(elementKind(t, s, "UseCase")["extensions"])["branches"] = []string{}
		requireFinding(t, s, "ownership.cycle")
	})
	t.Run("system containment cycle", func(t *testing.T) {
		s := cloneSnapshot(t, original)
		sys := elementKind(t, s, "System")
		sys["parent_system"] = "second-system"
		s.Elements = append(s.Elements, Record{"id": "second-system", "kind": "System", "source_document": sys["source_document"], "decomposition": "unsegmented", "parent_system": sys["id"]})
		requireFinding(t, s, "system.parent-cycle")
	})
}

func TestValidationAllowsUndefinedQualityLabel(t *testing.T) {
	s := validationSnapshot(t)
	requirement := elementKind(t, s, "GeneralRequirement")
	requirement["attributes"] = []string{}
	for _, label := range s.LabelOccurrences {
		if String(label, "owner") == String(requirement, "id") {
			label["label"] = "unavailable"
			result := snapshotRecord(t, s, String(label, "resolution"))
			result["status"] = "undefined_label"
			delete(result, "target")
		}
	}
	s.Assertions = slices.DeleteFunc(s.Assertions, func(r Record) bool {
		return String(r, "kind") == "QualityClassification" && String(r, "subject") == String(requirement, "id")
	})
	requireValid(t, s)
}

func TestValidationAllowsInvalidSourceWithHonestFallback(t *testing.T) {
	s := validationSnapshot(t)
	doc := ParseContent("invalid.use-case.spec.md", []byte("---\ntype: [\n---\n# Unfinished\n## Goal\nUnknown.\n"))
	mapped := MapDomain(doc, "use-case")
	s.Sources = append(s.Sources, doc.Source)
	s.Documents = append(s.Documents, doc.Document)
	s.ContentNodes = append(s.ContentNodes, doc.Nodes...)
	s.Fragments = append(s.Fragments, doc.Fragments...)
	s.ContentBindings = append(s.ContentBindings, mapped.Bindings...)
	s.Diagnostics = append(s.Diagnostics, doc.Diagnostics...)
	s.Diagnostics = append(s.Diagnostics, mapped.Diagnostics...)
	if len(mapped.Elements) > 0 {
		t.Fatal("invalid source unexpectedly had a domain root")
	}
	refreshTestCoverage(s)
	requireValid(t, s)
	if String(Object(doc.Document["coverage"]), "status") != "partial" {
		t.Fatal("invalid source should expose partial domain coverage")
	}
}

func refreshTestCoverage(s *Snapshot) {
	b := snapshotBuilder{snapshot: s, byPath: map[string]*ParsedDocument{}}
	for _, d := range s.Documents {
		p := &ParsedDocument{Document: d, NodeByID: map[string]Record{}}
		for _, source := range s.Sources {
			if String(source, "id") == String(d, "source") {
				p.Source = source
			}
		}
		for _, node := range s.ContentNodes {
			if String(node, "document") == String(d, "id") {
				p.Nodes = append(p.Nodes, node)
				p.NodeByID[String(node, "id")] = node
			}
		}
		b.byPath[String(p.Source, "locator")] = p
	}
	for _, source := range s.Sources {
		if String(source, "kind") == "SourceDocument" && String(source, "read_status") != "read" {
			b.readFailed = true
		}
	}
	b.finishCoverage()
}

func TestValidationRejectsAmbiguousJSONTransport(t *testing.T) {
	for _, data := range [][]byte{
		[]byte(`{"model_version":"first","model_version":"second"}`),
		[]byte(`{"value":{"a":1,"\u0061":2}}`),
		[]byte(`{} {}`), []byte(`{} trailing`),
		[]byte(`{"value":"\ud800"}`), []byte(`{"value":"\udfff"}`),
		append([]byte(`{"value":"`), append([]byte{0xff}, []byte(`"}`)...)...),
	} {
		findings, err := ValidateJSON(data)
		if err != nil || len(findings) != 1 || findings[0].Code != "shape.json" {
			t.Errorf("ambiguous transport %q: %v %+v", data, err, findings)
		}
	}
}

func TestValidationAcceptsEquivalentIntegerJSONNotation(t *testing.T) {
	snapshot := validationSnapshot(t)
	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	numbers := regexp.MustCompile(`("(?:byte_length|start|end)":)([0-9]+)`)
	data = numbers.ReplaceAll(data, []byte("${1}${2}.0"))
	findings, err := ValidateJSON(data)
	if err != nil || len(findings) > 0 {
		t.Fatalf("equivalent integral numbers must validate: %v %+v", err, findings)
	}
	if _, err := decodeSnapshotJSON([]byte(`{"emoji":"\ud83d\ude00","literal":"\\ud800"}`)); err != nil {
		t.Fatalf("valid Unicode escapes: %v", err)
	}
}

func TestValidationKnownRoleFallbackCannotClaimCompleteMapping(t *testing.T) {
	snapshot := validationSnapshot(t)
	for _, binding := range snapshot.ContentBindings {
		if String(binding, "owner") == String(elementKind(t, snapshot, "UseCase"), "id") && String(binding, "field_path") == "goal" && String(binding, "role") == "named_field" {
			binding["role"] = "fallback"
			binding["known_role"] = "UseCase.goal"
			binding["reason"] = "Mapping was not completed."
		}
	}
	requireFinding(t, snapshot, "coverage.known-role")
	requireFinding(t, snapshot, "coverage.incomplete")
	refreshTestCoverage(snapshot)
	requireValid(t, snapshot)
}

func TestValidationShapeFailureUsesConcreteRecordAndSource(t *testing.T) {
	original := validationSnapshot(t)
	cases := []struct {
		name, path string
		mutate     func(Record)
	}{
		{"wrong domain field type", "/goal", func(r Record) { r["goal"] = 42 }},
		{"unknown domain field", "/unexpected", func(r Record) { r["unexpected"] = true }},
		{"missing required identity", "/source_document", func(r Record) { delete(r, "source_document") }},
		{"unknown record kind", "/kind", func(r Record) { r["kind"] = "UnknownUseCase" }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			snapshot := cloneSnapshot(t, original)
			record := elementKind(t, snapshot, "UseCase")
			test.mutate(record)
			findings, err := Validate(snapshot)
			if err != nil {
				t.Fatal(err)
			}
			if len(findings) != 1 {
				t.Fatalf("expected one targeted shape failure, got %+v", findings)
			}
			finding := findings[0]
			if finding.Code != "shape.invalid" || !strings.HasSuffix(finding.Path, test.path) || finding.Path == "/" || len(finding.Message) > 1100 {
				t.Fatalf("not an actionable field failure: %+v", finding)
			}
			if strings.Contains(finding.Message, "GeneralRequirement") || strings.Contains(finding.Message, "QualityAttribute") {
				t.Fatalf("unrelated subtype errors leaked: %s", finding.Message)
			}
			if test.name != "missing required identity" && (!strings.HasSuffix(finding.Source, "check.use-case.spec.md") || finding.Line < 1) {
				t.Fatalf("expected original source location: %+v", finding)
			}
		})
	}
}

func TestValidationRejectsDeletedSourceSubtreeDespiteRetainedPayload(t *testing.T) {
	snapshot := validationSnapshot(t)
	extension := elementKind(t, snapshot, "ExtensionSection")
	owner := snapshotRecord(t, snapshot, String(extension, "owner"))
	removed := map[string]bool{}
	var remove func(string)
	remove = func(id string) {
		if removed[id] {
			return
		}
		removed[id] = true
		for _, child := range Strings(snapshotRecord(t, snapshot, id)["children"]) {
			remove(child)
		}
	}
	for _, field := range []string{"heading", "content"} {
		for _, id := range Strings(Object(extension[field])["content_nodes"]) {
			remove(id)
		}
	}
	extensionID := String(extension, "id")
	owner["extension_sections"] = slices.DeleteFunc(Strings(owner["extension_sections"]), func(id string) bool { return id == extensionID })
	for _, document := range snapshot.Documents {
		document["elements"] = slices.DeleteFunc(Strings(document["elements"]), func(id string) bool { return id == extensionID })
	}
	for _, node := range snapshot.ContentNodes {
		if node["children"] != nil {
			node["children"] = slices.DeleteFunc(Strings(node["children"]), func(id string) bool { return removed[id] })
		}
	}
	snapshot.ContentNodes = slices.DeleteFunc(snapshot.ContentNodes, func(r Record) bool { return removed[String(r, "id")] })
	snapshot.Elements = slices.DeleteFunc(snapshot.Elements, func(r Record) bool { return String(r, "id") == extensionID })
	snapshot.Fragments = slices.DeleteFunc(snapshot.Fragments, func(r Record) bool { return removed[String(r, "content_node")] })
	snapshot.ContentBindings = slices.DeleteFunc(snapshot.ContentBindings, func(r Record) bool { return String(r, "owner") == extensionID || removed[String(r, "content_node")] })
	refreshTestCoverage(snapshot)
	// The original payload, digest and root span still look intact, and all
	// remaining records are internally bound. Byte partitioning reveals the loss.
	findings, err := Validate(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) == 0 {
		t.Fatal("unrepresented source subtree was accepted")
	}
	for _, finding := range findings {
		if finding.Code != "content.gap" {
			t.Errorf("unexpected collateral invalidity: %+v", finding)
		}
	}
}
