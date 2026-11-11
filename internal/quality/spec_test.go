package quality

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const draft = "---\ntype: quality\n---\n# Draft\n"
const example = `---
type: quality
id: repeatable-inspection
quality_attributes: [repeatability, diagnosability]
---
# Repeatable inspection

## Applies to
Completed checks using the same build, arguments, paths, environment and source tree.

## Statement
Repeated inspections produce equivalent structured reports, including diagnostic order.

## Measurement plan
Compare parsed JSON from repeated fixture runs.

### Criterion
This is a measurement subsection, not another requirement criterion.

## Rationale
Stable findings make changes easier to identify.

## Related specifications
See [this requirement](#repeatable-inspection) and [the statement](#statement).
`

func TestRequirementPreservesContentAndLabels(t *testing.T) {
	spec, err := Parse([]byte(example))
	if err != nil {
		t.Fatal(err)
	}
	if spec.ID != "repeatable-inspection" || spec.System != nil || !reflect.DeepEqual(spec.QualityAttributes, []string{"repeatability", "diagnosability"}) || spec.Applicability == "" || spec.Statement == "" || spec.MeasurementPlan == "" || spec.Rationale == "" {
		t.Fatalf("fields: %+v", spec)
	}
	if !strings.Contains(spec.MeasurementPlan, "### Criterion") || strings.Contains(spec.Statement, "measurement subsection") || !strings.Contains(spec.Markdown, "## Related specifications") || len(Check(spec)) != 0 {
		t.Fatalf("Markdown and advice: %+v", spec)
	}
	if _, err := Parse([]byte(strings.Replace(example, "# Repeatable inspection", "# A different title", 1))); err != nil {
		t.Fatal(err)
	}
}

func TestRequirementDraftAdvice(t *testing.T) {
	spec, err := Parse([]byte(draft + "## Verification plan\n"))
	if err != nil {
		t.Fatal(err)
	}
	spec.Path = "source.spec.md"
	if spec.QualityAttributes == nil {
		t.Fatal("empty labels must be an array")
	}
	advice := Check(spec)
	if len(advice) != 4 {
		t.Fatalf("advice: %+v", advice)
	}
	for _, d := range advice {
		if d.Severity != "advice" || d.Location.Path != spec.Path || d.Location.Line < 4 {
			t.Fatalf("location: %+v", d)
		}
		if d.Rule == "quality.verification-plan" && d.Location.Line != 5 {
			t.Fatalf("missing empty heading location: %+v", d)
		}
	}
}

func TestRequirementRejectsMalformedLabels(t *testing.T) {
	for _, value := range []string{"null", "", "performance", "{}", "[null]", "[true]", "[12]", "['']", "[' ']", "[[]]", "[performance, performance]", "[quality.performance]", "[' leading']"} {
		source := strings.Replace(draft, "type: quality", "type: quality\nquality_attributes: "+value, 1)
		for _, source := range []string{source, "\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")} {
			_, err := Parse([]byte(source))
			var d *diagnostic.Diagnostic
			if !errors.As(err, &d) || d.Location.Line != 3 || d.Location.Column < 1 {
				t.Errorf("accepted or mislocated %q: %+v %v", value, d, err)
			}
		}
	}
	for _, value := range []string{"[]", "[performance]", "\n  - performance\n  - usability"} {
		if _, err := Parse([]byte(strings.Replace(draft, "type: quality", "type: quality\nquality_attributes: "+value, 1))); err != nil {
			t.Fatalf("valid labels %q: %v", value, err)
		}
	}
}

func TestRequirementStructure(t *testing.T) {
	for _, tt := range []struct{ source, rule string }{
		{"# Missing frontmatter", "metadata.frontmatter"},
		{strings.Replace(draft, "type: quality", "type: concept", 1), "metadata.type"},
		{strings.Replace(draft, "type: quality", "type: quality\nmodel: ISO", 1), "metadata.invalid"},
		{strings.Replace(draft, "type: quality", "type: quality\nid: quality.test", 1), "identity.invalid"},
		{draft + "# Second title", "structure.title"},
		{draft + "[broken](#missing)", "reference.unresolved"},
		{draft + "### Criterion\n", "structure.section-level"},
	} {
		_, err := Parse([]byte(tt.source))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != tt.rule {
			t.Errorf("wanted %s, got %+v %v", tt.rule, d, err)
		}
	}
	for _, section := range []string{"Applies to", "Criterion", "Measurement plan", "Rationale"} {
		_, err := Parse([]byte(draft + "## " + section + "\n\n## " + section + "\n"))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != "structure.duplicate-section" || len(d.Related) != 1 {
			t.Fatalf("duplicate section: %+v %v", d, err)
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{draft, example, Template, strings.Replace(draft, "type: quality", "type: quality\nquality_attributes: [null]", 1)} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if spec, err := Parse(data); err == nil {
			_ = Check(spec)
		}
	})
}
