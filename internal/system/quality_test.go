package system

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const qualityExample = `
## Quality attributes
This is a local vocabulary, not an acceptance decision.

### Responsiveness
**ID:** response
**Broader:** performance, usability

Time before the user receives a result.

#### Examples
Keep this subsection in the definition.

### Performance
**ID:** performance
Resources and time used under a defined workload.

### Usability
**ID:** usability
The user's ability to complete tasks.

### First result
**ID:** first-result
**Broader:** response
Time before the first useful result.
`

func TestQualityModelAndStableIDs(t *testing.T) {
	source := validSpec + qualityExample + "\n## Notes\nSee [response](#response).\n"
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	m := spec.QualityModel
	if m == nil || len(m.Attributes) != 4 || !strings.Contains(m.Markdown, "local vocabulary") {
		t.Fatalf("model: %+v", m)
	}
	a := m.Find("response")
	if a.Name != "Responsiveness" || !reflect.DeepEqual(a.Broader, []string{"performance", "usability"}) || !strings.Contains(a.Definition, "#### Examples") || strings.Contains(a.Definition, "**ID:") {
		t.Fatalf("attribute: %+v", a)
	}
	if got := m.Ancestors([]string{"first-result", "response", "unbound"}); !reflect.DeepEqual(got, []string{"response", "performance", "usability"}) {
		t.Fatalf("ancestors: %v", got)
	}
	if _, err := Parse([]byte(strings.Replace(source, "### Responsiveness", "### Renamed quality", 1))); err != nil {
		t.Fatalf("renaming a heading broke identity: %v", err)
	}
	plain, _ := Parse([]byte(validSpec))
	empty, err := Parse([]byte(validSpec + "\n## Quality attributes\n"))
	if err != nil || plain.QualityModel != nil || empty.QualityModel == nil || empty.QualityModel.Attributes == nil {
		t.Fatalf("absent vs empty: %+v %+v %v", plain, empty, err)
	}
}

func TestQualityModelStructure(t *testing.T) {
	for _, tt := range []struct{ name, body, rule string }{
		{"duplicate section", "## Quality attributes\n## Quality attributes\n", "structure.duplicate-section"},
		{"wrong attribute level", "## Quality attributes\n#### Usability\n**ID:** usability", "structure.attribute-heading"},
		{"empty attribute name", "## Quality attributes\n###\n**ID:** empty", "structure.attribute-heading"},
		{"missing ID", "## Quality attributes\n### Usability\nAn idea.", "identity.required"},
		{"empty ID", "## Quality attributes\n### Usability\n**ID:**", "identity.required"},
		{"invalid ID", "## Quality attributes\n### Usability\n**ID:** quality.usability", "identity.invalid"},
		{"duplicate ID", "## Quality attributes\n### A\n**ID:** same\n### B\n**ID:** same", "identity.duplicate"},
		{"duplicate field", "## Quality attributes\n### A\n**ID:** a\n**ID:** a", "structure.duplicate-field"},
		{"unknown parent", "## Quality attributes\n### A\n**ID:** a\n**Broader:** missing", "quality.broader"},
		{"empty parent", "## Quality attributes\n### A\n**ID:** a\n**Broader:**", "quality.broader"},
		{"trailing comma", "## Quality attributes\n### A\n**ID:** a\n**Broader:** a,", "quality.broader"},
		{"duplicate parent", "## Quality attributes\n### A\n**ID:** a\n**Broader:** a, a", "identity.duplicate"},
		{"self cycle", "## Quality attributes\n### A\n**ID:** a\n**Broader:** a", "quality.cycle"},
		{"cycle", "## Quality attributes\n### A\n**ID:** a\n**Broader:** b\n### B\n**ID:** b\n**Broader:** a", "quality.cycle"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{validSpec + tt.body, "\ufeff" + strings.ReplaceAll(validSpec+tt.body, "\n", "\r\n")} {
				_, err := Parse([]byte(source))
				var d *diagnostic.Diagnostic
				if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line <= 10 || d.Location.Column < 1 {
					t.Fatalf("got %+v (%v), want %s", d, err, tt.rule)
				}
				if tt.name == "duplicate ID" && len(d.Related) != 1 {
					t.Fatalf("missing first declaration: %+v", d)
				}
			}
		})
	}
}

func TestQualityMetadataOnlyComesFromIntroParagraphs(t *testing.T) {
	body := "\n## Quality attributes\n### A\n**ID:** a\n\n```markdown\n**ID:** code\n**Broader:** absent\n### Fake\n```\n\n> **ID:** quote\n\n- **Broader:** absent\n\n#### Notes\n**ID:** documentation\n"
	spec, err := Parse([]byte(validSpec + body))
	if err != nil {
		t.Fatal(err)
	}
	a := spec.QualityModel.Attributes[0]
	if a.ID != "a" || len(a.Broader) != 0 || !strings.Contains(a.Definition, "**ID:** code") || !strings.Contains(a.Definition, "**ID:** documentation") {
		t.Fatalf("metadata swallowed prose: %+v", a)
	}
	blank, err := Parse([]byte(validSpec + "\n## Quality attributes\n### A\n**ID:** a\n"))
	if err != nil {
		t.Fatal(err)
	}
	blank.Path = "specs/system.spec.md"
	advice := Check(blank)
	if len(advice) != 1 || advice[0].Rule != "quality.definition" || advice[0].Location.Path != blank.Path {
		t.Fatalf("definition advice: %+v", advice)
	}
}
