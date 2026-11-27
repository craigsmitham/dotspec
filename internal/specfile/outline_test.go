package specfile

import (
	"errors"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
)

func TestSharedOutlinePreservesProseAndLocations(t *testing.T) {
	source := "---\ntype: system-spec\n---\n# **Café** `parser`\n\n## Purpose\nAuthored *prose*.\n\n### Detail\nKeep [this](#detail).\n\n```md\n# Example only\n```\n\n> ## Quoted section\n\n## Scope\nOne project.\n"
	for _, data := range []string{source, "\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")} {
		outline, err := NewMarkdown([]byte(data)).TypedOutline("system", &Metadata{})
		if err != nil {
			t.Fatal(err)
		}
		if outline.Name != "Café parser" || outline.NameLocation.Line != 4 {
			t.Fatalf("heading extraction: %+v", outline)
		}
		if len(outline.Sections) != 2 || len(outline.Sections[0].Children) != 1 {
			t.Fatalf("example headings became structure: %+v", outline.Sections)
		}
		purpose := outline.Body(outline.Sections[0])
		if !strings.Contains(purpose, "Authored *prose*.") || !strings.Contains(purpose, "> ## Quoted section") {
			t.Fatalf("lost authored prose: %q", purpose)
		}
		if location := outline.Location(outline.Sections[1].Start, "Scope"); location.Line != 18 {
			t.Fatalf("wrong source location: %+v", location)
		}
	}
}

// parseStructure follows the parse path of a built-in kind: typed outline,
// explicit identity, captured sections and same-document references.
func parseStructure(data []byte) (string, error) {
	md := NewMarkdown(data)
	var meta struct {
		Metadata `yaml:",inline"`
		ID       string `yaml:"id"`
	}
	doc, err := md.TypedOutline("system", &meta)
	if err != nil {
		return "", err
	}
	var ids Identifiers
	if err := ids.Register(meta.ID, "system", md.MetadataLocation("id")); err != nil {
		return "", err
	}
	ids.AddAnchors(doc.Anchors)
	var purpose, scope string
	if err := doc.Capture(map[string]*string{"Purpose": &purpose, "Scope": &scope}, map[string]diagnostic.Location{}); err != nil {
		return "", err
	}
	return purpose, doc.ValidateReferences()
}

func TestSharedStructureDiagnosticsAndSourceLocations(t *testing.T) {
	const draft = "---\ntype: system-spec\n---\n# Job Map\n"
	for _, tt := range []struct {
		name, source, rule    string
		line, column, related int
	}{
		{"frontmatter", "# Draft", "metadata.frontmatter", 1, 1, 0},
		{"type", strings.Replace(draft, "type: system-spec", "type: job-spec", 1), "metadata.type", 2, 7, 0},
		{"unknown field", strings.Replace(draft, "type: system-spec", "type: system-spec\ndefined_by: other", 1), "metadata.invalid", 3, 1, 0},
		{"invalid ID", strings.Replace(draft, "type: system-spec", "type: system-spec\nid: system.task", 1), "identity.invalid", 3, 5, 0},
		{"wrong ID type", strings.Replace(draft, "type: system-spec", "type: system-spec\nid: [task]", 1), "metadata.invalid", 3, 1, 0},
		{"duplicate section", draft + "## Purpose\n\n## Purpose\n", "structure.duplicate-section", 7, 1, 5},
		{"wrong level", draft + "### Purpose\n", "structure.section-level", 5, 1, 0},
		{"second title", draft + "# Another\n", "structure.title", 5, 1, 4},
		{"empty title", strings.Replace(draft, "# Job Map", "#", 1), "structure.title", 4, 1, 0},
		{"before title", "---\ntype: system-spec\n---\n## Purpose\n# Title\n", "structure.section-order", 4, 1, 0},
		{"unresolved reference", draft + "## Evidence\n[Unknown](#missing)\n", "reference.unresolved", 6, 1, 0},
		{"unicode reference", draft + "## Evidence\né [Unknown](#missing)\n", "reference.unresolved", 6, 3, 0},
		{"malformed fragment", draft + "## Evidence\n[Unknown](#bad%XX)\n", "reference.invalid", 6, 1, 0},
		{"fenced example", draft + "## Evidence\n```markdown\n# Example\n## Purpose\n[not checked](#missing)\n```\n\n> ## Purpose\n> quoted content\n", "", 0, 0, 0},
		{"explicit ID after rename", strings.Replace(strings.Replace(draft, "# Job Map", "# Renamed map", 1), "type: system-spec", "type: system-spec\nid: job-map", 1) + "## Evidence\n[This system](#job%2Dmap)\n", "", 0, 0, 0},
		{"renamed without ID", strings.Replace(draft, "# Job Map", "# Renamed map", 1) + "## Evidence\n[This system](#job%2Dmap)\n", "reference.unresolved", 6, 1, 0},
		{"other file", draft + "[elsewhere](missing.spec.md#unknown)\n", "", 0, 0, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{tt.source, "\ufeff" + strings.ReplaceAll(tt.source, "\n", "\r\n")} {
				purpose, err := parseStructure([]byte(source))
				if tt.rule == "" {
					if err != nil || purpose != "" {
						t.Fatalf("example content became structure: %q, %v", purpose, err)
					}
					continue
				}
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
