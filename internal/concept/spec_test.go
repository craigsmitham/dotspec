package concept

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
)

const draft = "---\ntype: concept-spec\n---\n# Job Map\n"

const example = `---
type: concept-spec
id: job-map
---
# Job Map

## Purpose
Understand an accomplishment and its success criteria.

## State
A job, ordered steps, and outcomes associated with the job or one step.

| Element | Meaning |
| --- | --- |
| Step | A part of the accomplishment |

## Actions
### Add a step
Describe a step in the job document.

#### Effects
This subsection describes the action's state; it is not another concept section.

### View a map
Display steps in authored order. This does not require that execution order.

## Operational principles
Describe the steps, attach outcomes, and view the map to see how they fit together.

## Related specifications
JobMap<StepId> binds identities without interpreting another concept's state.

## Evidence
See [this concept](#job%2Dmap), [the action](#add-a-step), and
[an external definition](missing.spec.md#definition).
`

func TestParsePreservesConceptAndDistinctFields(t *testing.T) {
	spec, err := Parse([]byte(example))
	if err != nil {
		t.Fatal(err)
	}
	body := example[strings.Index(example, "# Job Map"):]
	if spec.Markdown != strings.TrimSpace(body) || spec.ID != "job-map" || spec.Name != "Job Map" {
		t.Fatalf("authored content or identity changed: %+v", spec)
	}
	if spec.Purpose != "Understand an accomplishment and its success criteria." || !strings.Contains(spec.State, "| Step |") || !strings.Contains(spec.Actions, "#### Effects") || !strings.HasPrefix(spec.OperationalPrinciples, "Describe the steps") {
		t.Fatalf("lost section content: %+v", spec)
	}
	if strings.Contains(spec.State, "This subsection") || strings.Contains(spec.OperationalPrinciples, "JobMap<StepId>") || len(Check(spec)) != 0 {
		t.Fatalf("section boundaries or advice: %+v; %+v", spec, Check(spec))
	}
}

func TestDraftsAndAdvice(t *testing.T) {
	for _, source := range []string{draft, draft + "\n## State\n", draft + "\n## Evidence\nAn unfinished idea.\n"} {
		spec, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		spec.Path = "specs/concepts/job-map.concept.spec.md"
		advice := Check(spec)
		var rules []string
		for _, d := range advice {
			rules = append(rules, d.Rule)
			if d.Severity != "advice" || d.Location.Path != spec.Path || d.Location.Column != 1 || d.Location.Line < 4 {
				t.Fatalf("invalid advice location/severity: %+v", d)
			}
			if d.Rule == "concept.state" && strings.Contains(source, "## State") && d.Location.Line != 6 {
				t.Fatalf("empty section should point at its heading: %+v", d)
			}
		}
		if !reflect.DeepEqual(rules, []string{"concept.purpose", "concept.state", "concept.actions", "concept.operational-principle"}) {
			t.Fatalf("unexpected advice: %+v", advice)
		}
	}
}

// Shared structure diagnostics are tested in specfile; these rows check that
// the concept parser reaches them with its own type and catalog class.
func TestStructureDiagnostics(t *testing.T) {
	for _, tt := range []struct {
		name, source, rule string
		line, column       int
	}{
		{"type", strings.Replace(draft, "type: concept-spec", "type: job-spec", 1), "metadata.type", 2, 7},
		{"wrong level", draft + "### Purpose\n", "structure.section-level", 5, 1},
		{"unresolved reference", draft + "## Evidence\n[Unknown](#missing)\n", "reference.unresolved", 6, 1},
	} {
		_, err := Parse([]byte(tt.source))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line != tt.line || d.Location.Column != tt.column {
			t.Fatalf("%s: %+v, %v; want %s %d:%d", tt.name, d, err, tt.rule, tt.line, tt.column)
		}
	}
}

func TestEveryConventionalSectionRejectsDuplicates(t *testing.T) {
	for _, name := range []string{"Purpose", "State", "Actions", "Operational principles"} {
		_, err := Parse([]byte(draft + "## " + name + "\n\n## " + name + "\n"))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != "structure.duplicate-section" {
			t.Fatalf("duplicate %s: %v", name, err)
		}
	}
}

func TestExampleHeadingsAreNotCaptured(t *testing.T) {
	spec, err := Parse([]byte(example + "\n```markdown\n## Purpose\nAn example.\n```\n\n> ## Purpose\n> quoted content\n"))
	if err != nil {
		t.Fatal(err)
	}
	if spec.Purpose != "Understand an accomplishment and its success criteria." {
		t.Fatalf("example heading replaced the purpose: %q", spec.Purpose)
	}
}

func FuzzParse(f *testing.F) {
	for _, source := range []string{draft, example, Template, draft + "## Purpose\n\n## Purpose\n", draft + "[bad](#missing)\n"} {
		f.Add([]byte(source))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if spec, err := Parse(data); err == nil {
			_ = Check(spec)
		}
	})
}
