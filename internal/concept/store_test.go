package concept

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeConcept(t *testing.T, dir, slug string) string {
	t.Helper()
	path := filepath.Join(dir, "specs", Collection, slug+Suffix)
	write(t, path, example)
	return path
}

func writeSystem(t *testing.T, dir, name string) {
	t.Helper()
	write(t, filepath.Join(dir, "specs/system.spec.md"), "---\ntype: system\n---\n# "+name+"\n## Purpose\nHelp people.\n## Scope\nThe project.\n")
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func TestInitializationAndLocalLookup(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	path, err := Init(dir, "job-map")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Show(dir, "job-map")
	if err != nil || got.Slug != "job-map" || got.System != nil || got.Name != "New concept" {
		t.Fatalf("initialized concept: %+v, %v", got, err)
	}
	if _, err := Init(dir, "job-map"); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("overwrite: %v", err)
	}
	if content, err := os.ReadFile(path); err != nil || string(content) != Template {
		t.Fatalf("existing file changed: %s, %v", content, err)
	}
	write(t, path, strings.Replace(example, "# Job Map", "# Map an accomplishment", 1))
	got, err = Show(dir, "job-map")
	if err != nil || got.Slug != "job-map" || got.ID != "job-map" || got.Name != "Map an accomplishment" {
		t.Fatalf("identity changed with title: %+v, %v", got, err)
	}
	child := filepath.Join(dir, "child")
	writeSystem(t, child, "Child")
	if _, err := Show(child, "job-map"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lookup must not fall back to an ancestor concept: %v", err)
	}
	writeConcept(t, child, "job-map")
	got, err = Show(child, "job-map")
	if err != nil || got.System == nil || got.System.Name != "Child" || got.Purpose != "Understand an accomplishment and its success criteria." {
		t.Fatalf("association must not rewrite the purpose: %+v, %v", got, err)
	}
}

func TestValidationBeforeCreatingFiles(t *testing.T) {
	dir := t.TempDir()
	for _, slug := range []string{"", "../escape", "two words", "two--hyphens", "Uppercase", "a/b", `a\b`, "work.concept.spec.md", "café"} {
		if _, err := Init(dir, slug); err == nil {
			t.Errorf("init accepted %q", slug)
		}
		if _, err := Show(dir, slug); err == nil {
			t.Errorf("show accepted %q", slug)
		}
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid invocation wrote files: %v, %v", entries, err)
	}
	if _, err := Init(filepath.Join(dir, "missing"), "valid"); err == nil {
		t.Fatal("init created the project directory")
	}
	for _, slug := range []string{"job-map", "2-steps", "check2"} {
		if err := ValidateSlug(slug); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGroupingAndAncestorContext(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	writeConcept(t, dir, "unassociated")
	for _, child := range []string{"a", "b"} {
		owner := filepath.Join(dir, child)
		writeSystem(t, owner, "Same name")
		writeConcept(t, owner, "work")
	}
	deep := filepath.Join(dir, "a", "deep")
	writeConcept(t, deep, "deep-work")
	groups, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 3 || len(groups[0].Concepts) != 2 || groups[0].System.Path == groups[1].System.Path || groups[2].System != nil {
		t.Fatalf("grouping merged distinct systems or lost unassociated cases: %+v", groups)
	}
	if groups[0].Concepts[0].Slug != "deep-work" || groups[0].Concepts[1].Slug != "work" {
		t.Fatalf("concepts not in source order: %+v", groups[0].Concepts)
	}
	groups, err = List(deep)
	if err != nil || len(groups) != 1 || len(groups[0].Concepts) != 1 || groups[0].System.Name != "Same name" {
		t.Fatalf("subtree context: %+v, %v", groups, err)
	}
	write(t, filepath.Join(dir, "a/specs/system.spec.md"), "broken")
	if _, err := List(deep); err == nil {
		t.Fatal("malformed enclosing system was ignored")
	}
}

func TestGitDiscoveryAndRepositoryBoundaries(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	writeConcept(t, dir, "tracked")
	git(t, dir, "add", "specs/concepts/tracked.concept.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "ignored/\ntracked.concept.spec.md\n")
	write(t, filepath.Join(dir, "ignored/broken.concept.spec.md"), "ignored")
	nested := filepath.Join(dir, "nested")
	writeConcept(t, nested, "nested")
	git(t, nested, "init", "-q")
	writeSystem(t, dir, "Outer")
	got, err := Show(nested, "nested")
	if err != nil || got.System != nil {
		t.Fatalf("association crossed a repository boundary: %+v, %v", got, err)
	}
	outside := t.TempDir()
	outsideFile := writeConcept(t, outside, "outside")
	if err := os.Symlink(outsideFile, filepath.Join(dir, "linked.concept.spec.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "linked-directory")); err != nil {
		t.Fatal(err)
	}
	groups, err := List(dir)
	if err != nil || len(groups) != 1 || len(groups[0].Concepts) != 1 || groups[0].Concepts[0].Slug != "tracked" {
		t.Fatalf("Git visibility: %+v, %v", groups, err)
	}
	plain := t.TempDir()
	write(t, filepath.Join(plain, ".gitignore"), "ignored/\n")
	writeConcept(t, filepath.Join(plain, "ignored"), "visible")
	groups, err = List(plain)
	if err != nil || len(groups) != 1 || groups[0].Concepts[0].Slug != "visible" {
		t.Fatalf("ignore rules applied outside Git: %+v, %v", groups, err)
	}
}

func TestReadDiagnosticsAndPreservation(t *testing.T) {
	dir := t.TempDir()
	path := writeConcept(t, dir, "work")
	source := draft + "## Purpose\n\n## Purpose\n"
	write(t, path, source)
	_, err := Read(path)
	var d *diagnostic.Diagnostic
	if !errors.As(err, &d) || d.Rule != "structure.duplicate-section" || d.Location.Path != path || d.Related[0].Path != path {
		t.Fatalf("source path missing from diagnostic: %+v", d)
	}
	if _, err := List(dir); err == nil {
		t.Fatal("list silently omitted malformed input")
	}
	if got, err := os.ReadFile(path); err != nil || !reflect.DeepEqual(got, []byte(source)) {
		t.Fatalf("read changed source: %q, %v", got, err)
	}
	for _, name := range []string{"Wrong.concept.spec.md", "work.md"} {
		_, err := Read(filepath.Join(dir, "specs", Collection, name))
		if !errors.As(err, &d) || d.Rule != "layout.concept-name" {
			t.Fatalf("filename diagnostic: %v", err)
		}
	}
}
