package qualitymodel

import (
	"reflect"
	"strings"
	"testing"
)

func TestBundledStarters(t *testing.T) {
	models, err := Builtins()
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 3 {
		t.Fatalf("starters: %d", len(models))
	}
	counts := map[string]int{"iso-25010-2023": 49, "furps": 36, "q42": 200}
	for _, m := range models {
		if len(m.QualityAttributes) != counts[m.Slug] || len(Check(m)) != 0 || m.Sources == "" || m.Edition == "" {
			t.Fatalf("incomplete starter: %+v", m)
		}
		if !strings.Contains(m.Description, "starter") {
			t.Fatalf("starter provenance missing: %s", m.Slug)
		}
		roots, children := 0, map[string]int{}
		for _, a := range m.QualityAttributes {
			if len(a.SubattributeOf) == 0 {
				roots++
			} else {
				for _, parent := range a.SubattributeOf {
					children[parent]++
				}
			}
		}
		wantRoots := map[string]int{"iso-25010-2023": 9, "furps": 5, "q42": 9}[m.Slug]
		if roots != wantRoots {
			t.Fatalf("%s roots: %d, want %d", m.Slug, roots, wantRoots)
		}
		expected := map[string]map[string]int{
			"iso-25010-2023": {"functional-suitability": 3, "performance-efficiency": 3, "compatibility": 2, "interaction-capability": 8, "reliability": 4, "security": 6, "maintainability": 5, "flexibility": 4, "safety": 5},
			"furps":          {"functionality": 3, "usability": 7, "reliability": 5, "performance": 8, "supportability": 9},
			"q42":            {"secure": 28, "usable": 62, "reliable": 61, "flexible": 28, "suitable": 35, "efficient": 35, "maintainable": 19, "operable": 41, "safe": 16},
		}[m.Slug]
		if !reflect.DeepEqual(children, expected) {
			t.Fatalf("%s hierarchy: %v, want %v", m.Slug, children, expected)
		}
	}
	if Default != "iso-25010-2023" {
		t.Fatal("unexpected default")
	}
}

func TestParseContentAndValidation(t *testing.T) {
	source := "---\ntype: quality-model\nid: local\n---\n# Local\n## Description\nA local model.\n## Edition\nTeam revision.\n## Sources\n[Origin](https://example.test).\n## Quality attributes\nA narrative.\n### Specific\n**ID:** specific\n**Sub-attribute of:** broad\nA definition.\n#### Detail\nPreserve this.\n### Broad\n**ID:** broad\nA broad concern.\n## Custom\nKeep unknown content.\n"
	m, err := Parse([]byte("\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")))
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "local" || len(m.QualityAttributes) != 2 || !strings.Contains(m.QualityAttributes[0].Definition, "#### Detail") || !strings.Contains(m.Markdown, "Keep unknown content.") {
		t.Fatalf("lost source: %+v", m)
	}
	for _, bad := range []string{
		strings.Replace(source, "**Sub-attribute of:** broad", "**Sub-attribute of:** missing", 1),
		strings.Replace(source, "**ID:** broad", "**ID:** specific", 1),
		strings.Replace(source, "A broad concern.", "**Sub-attribute of:** specific", 1),
		source + "## Edition\nDuplicate\n",
		strings.Replace(source, "## Description", "### Description", 1),
		strings.Replace(source, "id: local", "unknown: true", 1),
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Fatalf("accepted invalid source: %s", bad)
		}
	}
	absent, err := Parse([]byte("---\ntype: quality-model\n---\n# Draft\n"))
	if err != nil {
		t.Fatal(err)
	}
	empty, err := Parse([]byte("---\ntype: quality-model\n---\n# Draft\n## Quality attributes\n"))
	if err != nil {
		t.Fatal(err)
	}
	if absent.QualityAttributes != nil || empty.QualityAttributes == nil || empty.QualityAttributesContent == nil {
		t.Fatal("absent vs empty")
	}
}

func FuzzParse(f *testing.F) {
	f.Add([]byte(Template))
	models, err := Builtins()
	if err != nil {
		f.Fatal(err)
	}
	for _, m := range models {
		f.Add(m.Data)
	}
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = Parse(data) })
}
