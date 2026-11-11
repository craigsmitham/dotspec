package model

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	modelschema "github.com/craigsmitham/dotspec/model"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func mapped(t *testing.T, kind, body string) (*ParsedDocument, Mapping, Record) {
	t.Helper()
	d := ParseContent("specs/example."+kind+".spec.md", []byte("---\ntype: "+kind+"\n---\n\n# Example\n\n"+body))
	m := MapDomain(d, kind)
	if len(m.Diagnostics) > 0 {
		t.Fatalf("unexpected mapping diagnostics: %#v", m.Diagnostics)
	}
	for _, r := range m.Elements {
		if String(r, "id") == m.RootID {
			return d, m, r
		}
	}
	t.Fatal("no mapped root")
	return nil, Mapping{}, nil
}

func recordsOf(m Mapping, kind string) []Record {
	var out []Record
	for _, r := range m.Elements {
		if String(r, "kind") == kind {
			out = append(out, r)
		}
	}
	return out
}

func richText(doc *ParsedDocument, value any) string {
	var out []string
	for _, id := range Strings(Object(value)["content_nodes"]) {
		out = append(out, doc.Text(id))
	}
	return strings.Join(out, "")
}

func TestMappingAllCheckedInSpecimens(t *testing.T) {
	schemaDoc, err := jsonschema.UnmarshalJSON(bytes.NewReader(modelschema.SchemaJSON))
	if err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	const schemaURL = "https://dotspec.invalid/mapping-test.json"
	if err := compiler.AddResource(schemaURL, schemaDoc); err != nil {
		t.Fatal(err)
	}
	schemas := map[string]*jsonschema.Schema{}
	counts := map[string]int{}
	for _, root := range []string{"../../specs", "../../examples/starter/specs", "../../internal/integration/specs"} {
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".spec.md") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			d := ParseContent(path, data)
			kind := String(d.Document, "declared_type")
			m := MapDomain(d, kind)
			if len(m.Diagnostics) > 0 || m.RootID == "" {
				t.Errorf("%s: root=%q diagnostics=%#v", path, m.RootID, m.Diagnostics)
			}
			expectedClass := map[string]string{"business-rule": "BusinessRule", "system": "System", "job": "Job", "use-case": "UseCase", "concept": "Concept", "requirement": "GeneralRequirement", "quality": "QualityRequirement"}[kind]
			rootMatches := false
			for _, element := range m.Elements {
				if String(element, "id") == m.RootID && String(element, "kind") == expectedClass {
					rootMatches = true
				}
			}
			if !rootMatches {
				t.Errorf("%s: root does not match declared %s family", path, kind)
			}
			if got := recordsOf(m, "ExtensionSection"); len(got) > 0 {
				t.Errorf("%s: known corpus unexpectedly has unmapped sections: %#v", path, got)
			}
			for _, binding := range m.Bindings {
				if d.NodeByID[String(binding, "content_node")] == nil {
					t.Errorf("%s: binding has no syntax node: %#v", path, binding)
				}
			}
			all := append(append(append([]Record{}, m.Elements...), m.Bindings...), m.Labels...)
			for _, record := range all {
				class := String(record, "kind")
				compiled := schemas[class]
				if compiled == nil {
					var err error
					compiled, err = compiler.Compile(schemaURL + "#/$defs/" + class)
					if err != nil {
						t.Fatal(err)
					}
					schemas[class] = compiled
				}
				encoded, err := json.Marshal(record)
				if err != nil {
					t.Fatal(err)
				}
				value, err := jsonschema.UnmarshalJSON(bytes.NewReader(encoded))
				if err != nil {
					t.Fatal(err)
				}
				if err := compiled.Validate(value); err != nil {
					t.Errorf("%s %s schema: %v", path, class, err)
				}
			}
			counts[kind]++
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	total := 0
	for _, kind := range []string{"system", "job", "use-case", "concept", "requirement", "quality", "business-rule"} {
		if counts[kind] == 0 {
			t.Errorf("%s family absent from inventory", kind)
		}
		total += counts[kind]
	}
	if total < 23 {
		t.Errorf("specimen inventory unexpectedly shrank to %d (baseline 23)", total)
	}
}

func TestMappingProgressCaseAndJobIdentity(t *testing.T) {
	d, m, r := mapped(t, "job", `## Desired outcome

Agreed work.

## Steps

An explicitly partial map.

### Locate sources

**ID:** locate
**Stage:** custom-stage

Find evidence.

#### Outcomes

- **out-1:** Minimize time spent searching.
  - A supporting bullet, not a second outcome.

## Progress cases

### Adopt the tool

This case is a hypothesis.

#### Current approach

Manual review.

#### Desired progress

Faster review.

#### Forces

- **Push:** Repeated effort.
- **Pull:** Faster feedback.
- **Anxiety:** Another artifact.
- **Habit:** Existing conventions.

#### Timeline

To be established through interviews.

## Evidence

Unvalidated hypothesis.
`)
	if _, ok := r["intended_progress"]; !ok {
		t.Error("alias lost")
	}
	steps := recordsOf(m, "JobStep")
	if len(steps) != 1 || String(steps[0], "authored_id") != "locate" {
		t.Fatalf("steps: %#v", steps)
	}
	if got := strings.TrimSpace(richText(d, steps[0]["stage"])); got != "custom-stage" {
		t.Errorf("stage value contains unrelated metadata: %q", got)
	}
	if got := richText(d, steps[0]["description"]); strings.Contains(got, "**ID:") || strings.Contains(got, "Outcomes") {
		t.Errorf("step description includes metadata or outcomes: %q", got)
	}
	if len(recordsOf(m, "Outcome")) != 1 {
		t.Fatal("nested supporting bullets became outcomes")
	}
	cases := recordsOf(m, "ProgressCase")
	if len(cases) != 1 {
		t.Fatal("progress case missing")
	}
	forces := Object(cases[0]["forces"])["items"].([]Record)
	if len(forces) != 4 || String(forces[0], "kind") != "push" {
		t.Errorf("forces: %#v", forces)
	}
	if Object(Object(cases[0]["timeline"])["content"])["presence"] != "authored" {
		t.Error("unknown timeline was erased")
	}
	for _, x := range recordsOf(m, "Outcome") {
		if String(x, "owner") != String(steps[0], "id") {
			t.Error("step outcome assigned to job")
		}
	}
}

func TestMappingSystemQualityAttributesAndLabels(t *testing.T) {
	d, m, r := mapped(t, "system", `## Purpose

Describe a system.

## Scope

Its bounded behavior.

## Goals

- Improve clarity.
- Reduce effort.

## Glossary

### Draft

Unfinished material.

## Quality attributes

Local quality concerns.

### Usability

**ID:** usability

Can be used effectively.

### Reliability

**ID:** reliability

Works consistently.

### Diagnosability

**ID:** diagnostic
**Broader:** usability, reliability

Problems can be explained.
`)
	if len(recordsOf(m, "SystemGoal")) != 2 || len(recordsOf(m, "GlossaryTerm")) != 1 {
		t.Error("named system items missing")
	}
	if len(Strings(r["quality_attributes"])) != 3 || len(recordsOf(m, "QualityAttribute")) != 3 {
		t.Error("quality attributes missing")
	}
	if len(m.Labels) != 2 {
		t.Fatalf("broader labels: %#v", m.Labels)
	}
	for _, a := range recordsOf(m, "QualityAttribute") {
		if strings.Contains(richText(d, a["definition"]), "**ID:") {
			t.Error("attribute metadata included in definition")
		}
	}
	if got := richText(d, r["quality_attributes_content"]); !strings.Contains(got, "Local quality concerns.") || !strings.Contains(got, "Problems can be explained.") {
		t.Errorf("quality section content lost: %s", got)
	}
	for _, a := range recordsOf(m, "QualityAttribute") {
		if String(a, "owner") != String(r, "id") {
			t.Error("attribute must belong directly to the system")
		}
	}
	_, _, absent := mapped(t, "system", "## Purpose\n\nP.\n\n## Scope\n\nS.\n")
	_, _, empty := mapped(t, "system", "## Purpose\n\nP.\n\n## Scope\n\nS.\n\n## Quality attributes\n")
	if absent["quality_attributes"] != nil || absent["quality_attributes_content"] != nil || empty["quality_attributes"] == nil || len(Strings(empty["quality_attributes"])) != 0 || empty["quality_attributes_content"] == nil {
		t.Error("absent and empty attribute declarations conflated")
	}
}

func TestMappingUseCaseBranchesAndNarratives(t *testing.T) {
	_, m, r := mapped(t, "use-case", `## Goal

Achieve a result.

## Primary actor

A contributor.

## Supporting actors

- A service.

## Preconditions

- A working directory exists.

## Minimum guarantees

Source is retained.

## Main success scenario

1. Contributor requests inspection.
2. System reports findings.

## Extensions

An introductory narrative.

### Malformed source

The report identifies the error; handling follows the linked requirement.

#### Condition

The selected source is malformed.

#### Response

Report the error.

#### Flow

1. Preserve the source.
2. Continue inspection.

#### Continuation

Resume reporting.

#### A nested alternative

An additional explicitly headed branch.

## Sequence

**Mode:** representative

## Related use cases

See [another case](another.use-case.spec.md).
`)
	if len(recordsOf(m, "Scenario")) != 1 || len(recordsOf(m, "ScenarioStep")) != 4 {
		t.Errorf("scenarios/steps: %#v", m.Elements)
	}
	branches := recordsOf(m, "Extension")
	if len(branches) != 2 {
		t.Fatalf("branches: %#v", branches)
	}
	if branches[0]["condition"] == nil || branches[0]["continuation"] == nil {
		t.Error("explicit branch detail absent")
	}
	if branches[0]["handling_state"] != nil {
		t.Error("inferred handling status from prose")
	}
	if Object(r["sequence"])["mode"] != "representative" {
		t.Error("explicit sequence mode absent")
	}
	if r["related_use_cases"] == nil {
		t.Error("related use cases were treated as an extension")
	}
}

func TestMappingConceptActionDetailsAndQualityPlans(t *testing.T) {
	_, m, _ := mapped(t, "concept", `## Actions

Actions may be described without signatures.

### Complete

Complete the task.

#### Inputs

A task identifier.

#### Caller assumptions

- Caller has chosen a task.

#### Guards

- Task exists.

#### Effects

The task becomes complete.

#### Rejections

Unknown task is rejected.
`)
	actions := recordsOf(m, "ConceptAction")
	if len(actions) != 1 {
		t.Fatal("named action missing")
	}
	if actions[0]["caller_assumptions"] == nil || actions[0]["guards"] == nil {
		t.Error("assumptions and guards not represented separately")
	}
	_, _, r := mapped(t, "quality", `## Applies to

The selected source files.

### Exclusions

Concurrent edits.

## Criterion

Files remain unchanged.

## Measurement

Compare snapshots.

### Procedure

Read before and after inspection.

### Evidence to retain

Retain changed paths.

### Exclusions

Access times do not count.
`)
	if Object(r["applies_to"])["exclusions"] == nil || Object(r["measurement"])["evidence_to_retain"] == nil {
		t.Error("plan subfields missing")
	}
	if Object(r["criterion"])["target"] != nil {
		t.Error("invented quantitative target")
	}
}

func TestMappingPresentEmptyUnknownAndInvalid(t *testing.T) {
	_, m, r := mapped(t, "job", "## Job\n\n## Bespoke view\n\nKeep this custom material.\n")
	if Object(r["statement"])["presence"] != "empty" {
		t.Error("present empty section lost")
	}
	if r["performer"] != nil {
		t.Error("absent field manufactured")
	}
	if len(recordsOf(m, "ExtensionSection")) != 1 {
		t.Error("unknown domain gap not explicit")
	}
	bad := ParseContent("bad.job.spec.md", []byte("---\ntype: job\n---\n# Bad\n\n## Intended progress\n\nA.\n\n## Desired outcome\n\nB.\n"))
	invalid := MapDomain(bad, "job")
	if invalid.RootID != "" || len(invalid.Elements) != 0 || len(invalid.Diagnostics) != 1 {
		t.Fatalf("invalid source fabricated a domain object: %#v", invalid)
	}
	if String(invalid.Diagnostics[0], "code") != "structure.duplicate-section" {
		t.Error("existing authoring diagnostic changed")
	}
	if related, ok := invalid.Diagnostics[0]["related"].([]Record); !ok || len(related) != 1 || Number(related[0]["line"]) == 0 || String(related[0], "source") != String(bad.Source, "id") {
		t.Errorf("lost prior definition location: %#v", invalid.Diagnostics[0])
	}
	if bad.Source["payload"] == nil || len(bad.Nodes) == 0 {
		t.Error("invalid document lost source/content")
	}
}

func TestMappingMetadataLabelsHaveSourceOccurrences(t *testing.T) {
	data := []byte("---\ntype: quality\nid: sample\nattributes: [usability, reliability]\n---\n# Sample\n")
	d := ParseContent("sample.quality.spec.md", data)
	m := MapDomain(d, "quality")
	if len(m.Labels) != 2 {
		t.Fatalf("classification labels missing: %#v", m.Labels)
	}
	for _, label := range m.Labels {
		n := d.NodeByID[String(label, "content_node")]
		if String(n, "scalar_string") != String(label, "label") {
			t.Errorf("label does not point to its authored scalar: %#v", label)
		}
	}
}

func TestMappingOrdinaryMarkdownDoesNotBecomeASpecification(t *testing.T) {
	d := ParseContent("README.md", []byte("# Read me\n\nIntro.\n\n## Usage\n\nRun the tool.\n"))
	m := MapDomain(d, "")
	if m.RootID != "" || len(m.Diagnostics) != 0 {
		t.Fatalf("ordinary Markdown fabricated validation errors: %#v", m)
	}
	if len(recordsOf(m, "ExtensionSection")) != 1 {
		t.Error("ordinary section lost")
	}
	if d.Document["extension_sections"] != nil {
		t.Error("Document has invalid extension_sections property")
	}
}

func TestMappingRepeatedSupplementaryRolesDoNotOverwrite(t *testing.T) {
	d, m, r := mapped(t, "concept", "## Notes\n\nFirst note.\n\n## Notes\n\nSecond note.\n\n## Actions\n\n### Act\n\n#### Inputs\n\nFirst input.\n\n#### Inputs\n\nSecond input.\n")
	if text := richText(d, r["notes"]); !strings.Contains(text, "First note") || !strings.Contains(text, "Second note") {
		t.Errorf("lost repeated notes: %q", text)
	}
	a := recordsOf(m, "ConceptAction")[0]
	if text := richText(d, a["inputs"]); !strings.Contains(text, "First input") || !strings.Contains(text, "Second input") {
		t.Errorf("lost repeated input descriptions: %q", text)
	}
}
