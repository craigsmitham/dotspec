package specfile

import (
	"strings"
	"testing"
)

func TestObligationIdentityDeclarations(t *testing.T) {
	for _, tc := range []struct{ name, body, want string }{
		{"valid", "## Obligations\n### Preserve content\n**ID:** preserve\nKeep contents.\n", ""},
		{"empty", "## Obligations\n### Preserve content\n**ID:**\n", "nonempty"},
		{"invalid", "## Obligations\n### Preserve content\n**ID:** bad id\n", "invalid ID"},
		{"duplicate field", "## Obligations\n### Preserve content\n**ID:** one\n**ID:** two\n", "duplicate ID field"},
		{"duplicate identity", "## Obligations\n### First\n**ID:** one\n### Second\n**ID:** one\n", "duplicate ID"},
		{"fenced sample", "## Obligations\n### First\n```markdown\n**ID:** fake\n```\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := ParseOutline([]byte("---\ntype: system\n---\n# System\n"+tc.body), "system", &Metadata{})
			if err != nil {
				t.Fatal(err)
			}
			var ids Identifiers
			err = ids.RegisterObligations(doc)
			if tc.want != "" {
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("%v, want %s", err, tc.want)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			ids.AddAnchors(doc.Anchors)
			if tc.name == "valid" && !doc.Anchors["preserve"] {
				t.Fatal("identity not addressable")
			}
			if doc.Anchors["fake"] {
				t.Fatal("fenced example declared an identity")
			}
		})
	}
}
