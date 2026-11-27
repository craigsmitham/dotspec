package qualitymodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"strings"
	"testing"
)

const modelPrefix = "---\ntype: quality-model\n---\n# Example\n## Description\nQuality.\n## Sources\nOur team.\n## Edition\nFirst edition.\n"

func TestQualityAttributesStructure(t *testing.T) {
	for _, tt := range []struct{ name, body, rule string }{
		{"duplicate section", "## Quality attributes\n## Quality attributes\n", "structure.duplicate-section"},
		{"wrong attribute level", "## Quality attributes\n#### Usability\n**ID:** usability", "structure.member-heading"},
		{"empty attribute name", "## Quality attributes\n###\n**ID:** empty", "structure.member-heading"},
		{"missing ID", "## Quality attributes\n### Usability\nAn idea.", "identity.required"},
		{"empty ID", "## Quality attributes\n### Usability\n**ID:**", "identity.required"},
		{"invalid ID", "## Quality attributes\n### Usability\n**ID:** quality.usability", "identity.invalid"},
		{"duplicate ID", "## Quality attributes\n### A\n**ID:** same\n### B\n**ID:** same", "identity.duplicate"},
		{"duplicate field", "## Quality attributes\n### A\n**ID:** a\n**ID:** a", "structure.duplicate-field"},
		{"duplicate hierarchy field", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** b\n**Sub-attribute of:** c", "structure.duplicate-field"},
		{"unknown parent", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** missing", "quality.subattribute"},
		{"empty parent", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:**", "quality.subattribute"},
		{"trailing comma", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** a,", "quality.subattribute"},
		{"duplicate parent", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** a, a", "identity.duplicate"},
		{"self cycle", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** a", "quality.cycle"},
		{"cycle", "## Quality attributes\n### A\n**ID:** a\n**Sub-attribute of:** b\n### B\n**ID:** b\n**Sub-attribute of:** a", "quality.cycle"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{modelPrefix + tt.body, "\ufeff" + strings.ReplaceAll(modelPrefix+tt.body, "\n", "\r\n")} {
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
	body := "\n## Quality attributes\n### A\n**ID:** a\n\n```markdown\n**ID:** code\n**Sub-attribute of:** absent\n### Fake\n```\n\n> **ID:** quote\n\n- **Sub-attribute of:** absent\n\n#### Evidence\n**ID:** documentation\n"
	spec, err := Parse([]byte(modelPrefix + body))
	if err != nil {
		t.Fatal(err)
	}
	a := spec.QualityAttributes[0]
	if a.ID != "a" || len(a.SubattributeOf) != 0 || !strings.Contains(a.Definition, "**ID:** code") || !strings.Contains(a.Definition, "**ID:** documentation") {
		t.Fatalf("metadata swallowed prose: %+v", a)
	}
	blank, err := Parse([]byte(modelPrefix + "\n## Quality attributes\n### A\n**ID:** a\n"))
	if err != nil {
		t.Fatal(err)
	}
	if blank.QualityAttributes[0].Definition != "" {
		t.Fatal("unexpected definition")
	}

}

func TestQualityAttributeDeepHierarchy(t *testing.T) {
	var body strings.Builder
	body.WriteString(modelPrefix + "\n## Quality attributes\n")
	for i := 0; i < 64; i++ {
		fmt.Fprintf(&body, "### Attribute %d\n**ID:** a%d\n", i, i)
		if i < 63 {
			fmt.Fprintf(&body, "**Sub-attribute of:** a%d\n", i+1)
		}
		body.WriteString("Definition.\n")
	}
	s, err := Parse([]byte(body.String()))
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(s.QualityAttributes[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"subattribute_of":["a1"]`) {
		t.Fatalf("references must be an array: %s", data)
	}
	data, err = json.Marshal(s.QualityAttributes[63])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"subattribute_of":[]`) {
		t.Fatalf("root must have an empty reference array: %s", data)
	}
}
