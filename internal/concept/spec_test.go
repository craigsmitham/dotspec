package concept

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const draft = "---\ntype: concept\n---\n# Job Map\n"

const example = `---
type: concept
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

#### State
This subsection describes the action's state; it is not another concept section.

### View a map
Display steps in authored order. This does not require that execution order.

## Operational principle
Describe the steps, attach outcomes, and view the map to see how they fit together.

## Instantiation
JobMap<StepId> binds identities without interpreting another concept's state.

## Notes
See [this concept](#job%2Dmap), [the action](#add-a-step), and
[an external definition](missing.spec.md#definition).
`

func TestParsePreservesConceptAndDistinctFields(t *testing.T) {
	spec, err := Parse([]byte(example))
	if err != nil {
		t.Fatal(err)
	}
	body := example[strings.Index(example, "# Job Map"):]
	if spec.Markdown != strings.TrimSpace(body) || spec.ID != "job-map" || spec.Name != "Job Map" || spec.System != nil {
		t.Fatalf("authored content or identity changed: %+v", spec)
	}
	if spec.Purpose != "Understand an accomplishment and its success criteria." || !strings.Contains(spec.State, "| Step |") || !strings.Contains(spec.Actions, "#### State") || !strings.HasPrefix(spec.OperationalPrinciple, "Describe the steps") {
		t.Fatalf("lost section content: %+v", spec)
	}
	if strings.Contains(spec.State, "This subsection") || strings.Contains(spec.OperationalPrinciple, "JobMap<StepId>") || len(Check(spec)) != 0 {
		t.Fatalf("section boundaries or advice: %+v; %+v", spec, Check(spec))
	}
}

func TestDraftsAndAdvice(t *testing.T) {
	for _, source := range []string{draft, draft + "\n## State\n", draft + "\n## Notes\nAn unfinished idea.\n"} {
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

func TestStructureAndSourceLocations(t *testing.T) {
	for _, tt := range []struct {
		name, source, rule    string
		line, column, related int
	}{
		{"frontmatter", "# Draft", "metadata.frontmatter", 1, 1, 0},
		{"type", strings.Replace(draft, "type: concept", "type: job", 1), "metadata.type", 2, 7, 0},
		{"unknown field", strings.Replace(draft, "type: concept", "type: concept\ndefined_by: other", 1), "metadata.invalid", 3, 1, 0},
		{"invalid ID", strings.Replace(draft, "type: concept", "type: concept\nid: concept.task", 1), "identity.invalid", 3, 5, 0},
		{"wrong ID type", strings.Replace(draft, "type: concept", "type: concept\nid: [task]", 1), "metadata.invalid", 3, 1, 0},
		{"duplicate section", draft + "## Purpose\n\n## Purpose\n", "structure.duplicate-section", 7, 1, 5},
		{"wrong level", draft + "### Purpose\n", "structure.section-level", 5, 1, 0},
		{"second title", draft + "# Another\n", "structure.title", 5, 1, 4},
		{"empty title", strings.Replace(draft, "# Job Map", "#", 1), "structure.title", 4, 1, 0},
		{"before title", "---\ntype: concept\n---\n## Purpose\n# Title\n", "structure.section-order", 4, 1, 0},
		{"unresolved reference", draft + "## Notes\n[Unknown](#missing)\n", "reference.unresolved", 6, 1, 0},
		{"unicode reference", draft + "## Notes\né [Unknown](#missing)\n", "reference.unresolved", 6, 3, 0},
		{"malformed fragment", draft + "## Notes\n[Unknown](#bad%XX)\n", "reference.invalid", 6, 1, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{tt.source, "\ufeff" + strings.ReplaceAll(tt.source, "\n", "\r\n")} {
				_, err := Parse([]byte(source))
				var d *diagnostic.Diagnostic
				if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line != tt.line || d.Location.Column != tt.column {
					t.Fatalf("diagnostic: %+v, %v; want %s %d:%d", d, err, tt.rule, tt.line, tt.column)
				}
				if tt.related > 0 && (len(d.Related) != 1 || d.Related[0].Line != tt.related) {
					t.Fatalf("missing previous definition: %+v", d)
				}
			}
		})
	}
}

func TestEveryConventionalSectionRejectsDuplicates(t *testing.T) {
	for _, name := range []string{"Purpose", "State", "Actions", "Operational principle"} {
		_, err := Parse([]byte(draft + "## " + name + "\n\n## " + name + "\n"))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != "structure.duplicate-section" {
			t.Fatalf("duplicate %s: %v", name, err)
		}
	}
}

func TestExamplesAndUnknownProseRemainMarkdown(t *testing.T) {
	extra := "\n## More notes\n\n```markdown\n# Example\n## Purpose\n[not checked](#missing)\n```\n\n> ## Purpose\n> quoted content\n\n## purpose\nLowercase headings are ordinary prose.\n"
	spec, err := Parse([]byte(example + extra))
	if err != nil {
		t.Fatal(err)
	}
	if spec.Purpose != "Understand an accomplishment and its success criteria." || !strings.Contains(spec.Markdown, strings.TrimSpace(extra)) {
		t.Fatalf("examples changed parsed content: %+v", spec)
	}
	// Explicit identity is independent of the title's generated heading anchor.
	if _, err := Parse([]byte(strings.Replace(example, "# Job Map", "# Renamed map", 1))); err != nil {
		t.Fatalf("explicit ID stopped resolving after rename: %v", err)
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

func TestOperationalPrinciplesAliasAndDuplicate(t *testing.T) {
	source := "---\ntype: concept\n---\n# Example\n## Type parameters\n### Item\nAn identity.\n## Introduced types\n### Token\nA token.\n## Operational principles\n### First story\nOne.\n### Second story\nTwo.\n"
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if spec.TypeParameters == "" || spec.IntroducedTypes == "" || !strings.Contains(spec.OperationalPrinciple, "Second story") {
		t.Fatal("plural stories or declarations lost")
	}
	if _, err = Parse([]byte(source + "\n## Operational principle\nDuplicate.\n")); err == nil {
		t.Fatal("singular/plural duplicate accepted")
	}
}
