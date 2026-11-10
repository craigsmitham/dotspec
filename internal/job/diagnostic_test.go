package job

import (
	"errors"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestDiagnosticLocationsSurviveLineEndings(t *testing.T) {
	source := "---\ntype: job\nid: shared\n---\n# Example\n\n## Steps\n\n### Gather\n**ID:** shared\n"
	for _, data := range []string{source, "\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")} {
		_, err := Parse([]byte(data))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != "identity.duplicate" || d.Location.Line != 10 || d.Location.Column != 1 || len(d.Related) != 1 || d.Related[0].Line != 3 || d.Related[0].Column != 5 {
			t.Fatalf("duplicate source locations: %#v, %v", d, err)
		}
	}
}

func TestStructuralDiagnosticsHaveSourceLocations(t *testing.T) {
	for _, tt := range []struct {
		body, rule string
		line       int
		related    int
	}{
		{"## Intended progress\nOne\n## Desired outcome\nTwo\n", "structure.duplicate-section", 7, 5},
		{"## Outcomes\n- **same:** First\n- **same:** Second\n", "identity.duplicate", 7, 6},
		{"## Steps\n### Gather\n**Stage:** locate\n**Stage:** execute\n", "structure.duplicate-field", 8, 7},
		{"## Notes\n[Missing](#unknown)\n", "reference.unresolved", 6, 0},
	} {
		_, err := Parse([]byte("---\ntype: job\n---\n# Example\n" + tt.body))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line != tt.line || d.Location.Column < 1 {
			t.Fatalf("%s: %#v, %v", tt.rule, d, err)
		}
		if tt.related != 0 && (len(d.Related) != 1 || d.Related[0].Line != tt.related) {
			t.Fatalf("missing first definition: %+v", d)
		}
	}
}

func TestAdviceRetainsOutcomeScopeAndLocation(t *testing.T) {
	source := "---\ntype: job\n---\n# Example\n## Job\nPrepare a change for review.\n## Performer\nA reviewer.\n## Outcomes\n- **O1:** Better results.\n## Steps\n### Gather\n#### Outcomes\n- **O2:** Better results.\n"
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	spec.Path = "example.job.spec.md"
	advice := Check(spec)
	if len(advice) != 2 || advice[0].Location.Line != 10 || advice[1].Location.Line != 14 {
		t.Fatalf("outcome locations: %+v", advice)
	}
	if advice[0].Location.Element != "Overall job outcomes / O1" || advice[1].Location.Element != "Outcomes for step: Gather / O2" {
		t.Fatalf("outcome scopes: %+v", advice)
	}
	for _, d := range advice {
		if d.Rule != "job.outcome-measure" || d.Severity != "advice" || d.Location.Path != spec.Path {
			t.Fatalf("advice lost its rule or path: %+v", d)
		}
	}
}
