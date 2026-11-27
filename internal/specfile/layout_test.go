package specfile

import (
	"os"
	"path/filepath"
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

func TestCreateTitlesCollectionDraftFromSlug(t *testing.T) {
	dir := t.TempDir()
	layout := Layout{Kind: "requirement", Collection: "requirements", Suffix: ".requirement.spec.md"}
	path, err := layout.Create(dir, "fast-save", "---\ntype: requirement-spec\n---\n\n# New requirement\n\n## Statement\n")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "---\ntype: requirement-spec\n---\n\n# Fast save\n\n## Statement\n"; string(data) != want {
		t.Fatalf("draft = %q, want %q", data, want)
	}
}

func TestManifestUnknownKeyNamesTheKey(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ProjectFilename), []byte("name: x\nreport:\n  enabled: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readProject(dir)
	if err == nil || !strings.Contains(err.Error(), `line 2: unknown key "report"; dotspec.yaml accepts name and build.enabled`) || strings.Contains(err.Error(), "specfile.Project") {
		t.Fatalf("err = %v", err)
	}
}
