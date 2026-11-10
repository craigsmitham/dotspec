package specfile

import (
	"strings"
	"testing"
)

func TestSharedOutlinePreservesProseAndLocations(t *testing.T) {
	source := "---\ntype: system\n---\n# **Café** `parser`\n\n## Purpose\nAuthored *prose*.\n\n### Detail\nKeep [this](#detail).\n\n```md\n# Example only\n```\n\n> ## Quoted section\n\n## Scope\nOne project.\n"
	for _, data := range []string{source, "\ufeff" + strings.ReplaceAll(source, "\n", "\r\n")} {
		outline, err := ParseOutline([]byte(data), "system")
		if err != nil {
			t.Fatal(err)
		}
		document, err := Parse([]byte(data), "system", "Purpose", "Scope")
		if err != nil {
			t.Fatal(err)
		}
		if outline.Name != "Café parser" || outline.Name != document.Name || outline.NameLocation.Line != 4 {
			t.Fatalf("heading extraction: %+v", outline)
		}
		if len(outline.Sections) != 2 || len(outline.Sections[0].Children) != 1 {
			t.Fatalf("example headings became structure: %+v", outline.Sections)
		}
		purpose := outline.Body(outline.Sections[0])
		if purpose != document.Sections["Purpose"] || !strings.Contains(purpose, "Authored *prose*.") || !strings.Contains(purpose, "> ## Quoted section") {
			t.Fatalf("lost authored prose: %q", purpose)
		}
		if location := outline.Location(outline.Sections[1].Start, "Scope"); location.Line != 18 {
			t.Fatalf("wrong source location: %+v", location)
		}
	}
}
