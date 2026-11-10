package system

import (
	"errors"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestOptionalSystemContext(t *testing.T) {
	source := Template + "\n## Goals\nImprove reviews.\n## Environment\nGit repositories.\n## Assumptions\nSources do not change during checks.\n## Constraints\nPortable Markdown.\n## Glossary\n### Complete scan\nAll selected files were examined.\n"
	s, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if s.Goals != "Improve reviews." || s.Environment == "" || s.Assumptions == "" || s.Constraints == "" || !strings.Contains(s.Glossary, "### Complete scan") || !s.References.Anchors["complete-scan"] {
		t.Fatalf("context: %+v", s)
	}
	for _, name := range []string{"Goals", "Environment", "Assumptions", "Constraints", "Glossary"} {
		for _, bad := range []string{source + "\n## " + name + "\n", strings.Replace(Template, "## Purpose", "### "+name+"\n\n## Purpose", 1)} {
			_, err := Parse([]byte(bad))
			var d *diagnostic.Diagnostic
			if !errors.As(err, &d) || !strings.HasPrefix(d.Rule, "structure.") {
				t.Errorf("accepted bad context %s: %v", name, err)
			}
		}
	}
}
