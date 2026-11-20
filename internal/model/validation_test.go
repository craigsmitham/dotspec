package model

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func validationSnapshot(t *testing.T) *Snapshot {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"specs/example.system.spec.md":                                                "---\ntype: system-spec\n---\n# Example\n## Purpose\nSupport explicit source-backed specifications.\n## Scope\nThese examples.\n## Goals\n### Preserve author intent\n\nPreserve author intent.\n## Glossary\n### Source\nThe original authored bytes.\n## Quality attributes\n",
		"specs/quality/reliability/reliability.quality-attribute.spec.md":             "---\ntype: quality-attribute-spec\nid: reliable\n---\n# Reliability\n\n## Definition\n\nWorks as specified.\n",
		"specs/quality/reliability/consistency/consistency.quality-attribute.spec.md": "---\ntype: quality-attribute-spec\nid: consistent\nsubattribute_of: [reliable]\n---\n# Consistency\n\n## Definition\n\nProduces the same result from the same source.\n",

		"specs/use-cases/check.use-case.spec.md": `---
type: use-case-spec
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
[The system](../example.system.spec.md) and [requirements](../requirements/).
## Supports
[Inspect intent](../jobs/inspect.job.spec.md).
`,
		"specs/jobs/inspect.job.spec.md": `---
type: job-spec
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
##### Minimize overlooked constraints

**ID:** clear

Minimize overlooked constraints.
`,
		"specs/concepts/inspection.concept.spec.md": `---
type: concept-spec
---
# Inspection
## Purpose
Expose authored content.
## State
A set of source files.
## Actions
### Inspect
Return the authored content and findings.
## Operational principles
Inspection preserves original sources.
### Example
Retain this characteristic example in the operational principles.
`,
		"specs/requirements/preserve.requirement.spec.md": `---
type: requirement-spec
quality_attributes:
- reliable
---
# Preserve source
## Statement
Inspection preserves the original source bytes.
## Applies to
Source inspection.
## Evidence
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
	result, err := Export(dir, system.CheckType(), system.AttributeCheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), requirement.CheckType())
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
	for _, collection := range [][]Record{s.Sources, s.Documents, s.ContentNodes, s.Elements, s.Fragments, s.References, s.LabelOccurrences, s.Diagnostics} {
		for _, r := range collection {
			if String(r, "id") == id {
				return r
			}
		}
	}
	t.Fatalf("missing record %q", id)
	return nil
}

func TestValidationChecksCanonicalSourceAndGraph(t *testing.T) {
	original := validationSnapshot(t)
	requireValid(t, cloneSnapshot(t, original))
	cases := []struct {
		name, code string
		mutate     func(*Snapshot)
	}{
		{"missing owned member", "reference.missing", func(s *Snapshot) { elementKind(t, s, "Job")["steps"] = []string{"absent"} }},
		{"duplicate identity", "identity.duplicate", func(s *Snapshot) { s.Elements = append(s.Elements, s.Elements[0]) }},
		{"source outside field", "source.span", func(s *Snapshot) { records(Object(elementKind(t, s, "UseCase")["goal"])["spans"])[0]["end"] = 1000000 }},
		{"overlapping ranges", "content.order", func(s *Snapshot) {
			r := Object(elementKind(t, s, "Requirement")["statement"])
			span := records(r["spans"])[0]
			r["spans"] = []Record{span, span}
		}},
		{"unresolved endpoint", "resolution.state", func(s *Snapshot) { Object(s.LabelOccurrences[0]["resolution"])["status"] = "unresolved" }},
		{"wrong owner source", "reference.context", func(s *Snapshot) { s.References[0]["context_owner"] = elementKind(t, s, "Job")["id"] }},
		{"changed authored destination", "relationship.destination", func(s *Snapshot) {
			for _, ref := range s.References {
				if ref["role"] != nil && String(Object(ref["resolution"]), "status") == "resolved" {
					ref["authored_target"] = "different.spec.md#different"
					return
				}
			}
		}},
		{"unknown relationship", "relationship.role", func(s *Snapshot) { s.References[0]["role"] = "missing" }},
		{"enclosing cycle", "system.enclosing-cycle", func(s *Snapshot) { r := elementKind(t, s, "System"); r["enclosing_system"] = r["id"] }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { s := cloneSnapshot(t, original); c.mutate(s); requireFinding(t, s, c.code) })
	}
}
func TestValidationClosedShapesAndMalformedJSON(t *testing.T) {
	snapshot := validationSnapshot(t)
	elementKind(t, snapshot, "UseCase")["unexpected"] = true
	requireFinding(t, snapshot, "shape.invalid")
	if findings, err := ValidateJSON([]byte(`{"model_version":`)); err != nil || len(findings) != 1 || findings[0].Code != "shape.json" {
		t.Fatalf("invalid JSON must be a validation finding, got %v %+v", err, findings)
	}
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
			if strings.Contains(finding.Message, "Requirement") || strings.Contains(finding.Message, "QualityAttribute") {
				t.Fatalf("unrelated subtype errors leaked: %s", finding.Message)
			}
			if test.name != "missing required identity" && (!strings.HasSuffix(finding.Source, "check.use-case.spec.md") || finding.Line < 1) {
				t.Fatalf("expected original source location: %+v", finding)
			}
		})
	}
}
