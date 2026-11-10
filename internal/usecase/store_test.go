package usecase

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

func writeUseCase(t *testing.T, dir, slug string) string {
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
	path, err := Init(dir, "check-specifications")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Show(dir, "check-specifications")
	if err != nil || got.Slug != "check-specifications" || got.System != nil || got.Name != "New use case" {
		t.Fatalf("initialized use case: %+v, %v", got, err)
	}
	if _, err := Init(dir, "check-specifications"); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("overwrite: %v", err)
	}
	if content, err := os.ReadFile(path); err != nil || string(content) != Template {
		t.Fatalf("existing file changed: %s, %v", content, err)
	}
	write(t, path, strings.Replace(example, "# Check specifications", "# Review structural findings", 1))
	got, err = Show(dir, "check-specifications")
	if err != nil || got.Slug != "check-specifications" || got.ID != "check-specifications" || got.Name != "Review structural findings" {
		t.Fatalf("identity changed with title: %+v, %v", got, err)
	}
	child := filepath.Join(dir, "child")
	writeSystem(t, child, "Child")
	if _, err := Show(child, "check-specifications"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lookup must not fall back to an ancestor use case: %v", err)
	}
	writeUseCase(t, child, "check-specifications")
	got, err = Show(child, "check-specifications")
	if err != nil || got.System == nil || got.System.Name != "Child" || got.SystemBoundary != "The dotspec CLI." {
		t.Fatalf("association must not rewrite the stated boundary: %+v, %v", got, err)
	}
}

func TestValidationBeforeCreatingFiles(t *testing.T) {
	dir := t.TempDir()
	for _, slug := range []string{"", "../escape", "two words", "two--hyphens", "Uppercase", "a/b", `a\b`, "work.use-case.spec.md", "café"} {
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
	for _, slug := range []string{"check-specifications", "2-steps", "check2"} {
		if err := ValidateSlug(slug); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGroupingAndAncestorContext(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	writeUseCase(t, dir, "unassociated")
	for _, child := range []string{"a", "b"} {
		owner := filepath.Join(dir, child)
		writeSystem(t, owner, "Same name")
		writeUseCase(t, owner, "work")
	}
	deep := filepath.Join(dir, "a", "deep")
	writeUseCase(t, deep, "deep-work")
	groups, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 3 || len(groups[0].UseCases) != 2 || groups[0].System.Path == groups[1].System.Path || groups[2].System != nil {
		t.Fatalf("grouping merged distinct systems or lost unassociated cases: %+v", groups)
	}
	if groups[0].UseCases[0].Slug != "deep-work" || groups[0].UseCases[1].Slug != "work" {
		t.Fatalf("use cases not in source order: %+v", groups[0].UseCases)
	}
	groups, err = List(deep)
	if err != nil || len(groups) != 1 || len(groups[0].UseCases) != 1 || groups[0].System.Name != "Same name" {
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
	writeUseCase(t, dir, "tracked")
	git(t, dir, "add", "specs/use-cases/tracked.use-case.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "ignored/\ntracked.use-case.spec.md\n")
	write(t, filepath.Join(dir, "ignored/broken.use-case.spec.md"), "ignored")
	nested := filepath.Join(dir, "nested")
	writeUseCase(t, nested, "nested")
	git(t, nested, "init", "-q")
	writeSystem(t, dir, "Outer")
	got, err := Show(nested, "nested")
	if err != nil || got.System != nil {
		t.Fatalf("association crossed a repository boundary: %+v, %v", got, err)
	}
	outside := t.TempDir()
	outsideFile := writeUseCase(t, outside, "outside")
	if err := os.Symlink(outsideFile, filepath.Join(dir, "linked.use-case.spec.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "linked-directory")); err != nil {
		t.Fatal(err)
	}
	groups, err := List(dir)
	if err != nil || len(groups) != 1 || len(groups[0].UseCases) != 1 || groups[0].UseCases[0].Slug != "tracked" {
		t.Fatalf("Git visibility: %+v, %v", groups, err)
	}
	plain := t.TempDir()
	write(t, filepath.Join(plain, ".gitignore"), "ignored/\n")
	writeUseCase(t, filepath.Join(plain, "ignored"), "visible")
	groups, err = List(plain)
	if err != nil || len(groups) != 1 || groups[0].UseCases[0].Slug != "visible" {
		t.Fatalf("ignore rules applied outside Git: %+v, %v", groups, err)
	}
}

func TestReadDiagnosticsAndPreservation(t *testing.T) {
	dir := t.TempDir()
	path := writeUseCase(t, dir, "work")
	source := draft + "## Goal\n\n## Goal\n"
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
	for _, name := range []string{"Wrong.use-case.spec.md", "work.md"} {
		_, err := Read(filepath.Join(dir, "specs", Collection, name))
		if !errors.As(err, &d) || d.Rule != "layout.use-case-name" {
			t.Fatalf("filename diagnostic: %v", err)
		}
	}
}
