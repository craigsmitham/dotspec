package system

import (
	"errors"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func TestSystemDiagnosticLocations(t *testing.T) {
	for _, tt := range []struct {
		source, rule string
		line         int
	}{
		{"---\ntype: job\n---\n# Example\n", "metadata.type", 2},
		{"---\ntype: system\nunknown: field\n---\n# Example\n", "metadata.invalid", 3},
		{"---\ntype: system\n---\n# Example\n## Scope\nThe product.\n", "structure.required-section", 4},
		{"---\ntype: system\n---\n# Example\n## Purpose\nHelp.\n## Scope\n[Missing](#absent)\n", "reference.unresolved", 8},
	} {
		for _, source := range []string{tt.source, "\ufeff" + strings.ReplaceAll(tt.source, "\n", "\r\n")} {
			_, err := Parse([]byte(source))
			var d *diagnostic.Diagnostic
			if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line != tt.line || d.Location.Column < 1 {
				t.Errorf("%s: expected source line %d, got %+v (%v)", tt.rule, tt.line, d, err)
			}
		}
	}
}

func TestSystemHeadingReferences(t *testing.T) {
	source := "---\ntype: system\n---\n# Example\n## Purpose\nSee [scope](#scope) and [details](#details).\n## Scope\nThe product.\n### Details\n[External](https://example.com) and [file](other.md#other) stay unchecked.\n```markdown\n[Example](#not-an-actual-link)\n```\n"
	if _, err := Parse([]byte(source)); err != nil {
		t.Fatal(err)
	}
}
