package specfile

import (
	"os"
	"strings"
	"testing"
)

func TestLayoutEnforcesSlugGrammarForCreationAndLookup(t *testing.T) {
	notes := Layout{Kind: "note", Collection: "notes", Suffix: ".note.spec.md", SlugExample: "work"}
	dir, err := Directory(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	want := `invalid note slug "Bad--Slug": use lowercase letters or digits separated by single hyphens (for example, work)`
	if _, err := notes.Create(dir, "Bad--Slug", "# Note\n"); err == nil || err.Error() != want {
		t.Fatalf("create error = %v, want %q", err, want)
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid slug wrote files: %v, %v", entries, err)
	}
	if _, err := notes.LocalPath(dir, "Bad--Slug"); err == nil || err.Error() != want {
		t.Fatalf("lookup error = %v, want %q", err, want)
	}
	named := Layout{Kind: "system", Suffix: ".system.spec.md", NamedSubject: true}
	if _, err := named.LocalPath(dir, "Café"); err == nil || strings.Contains(err.Error(), "invalid system slug") {
		t.Fatalf("named subject lookup applied the collection slug grammar: %v", err)
	}
}
