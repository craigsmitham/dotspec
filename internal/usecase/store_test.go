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
	path := filepath.Join(dir, "specs", Layout.Collection, slug+Layout.Suffix)
	write(t, path, example)
	return path
}

func writeSystem(t *testing.T, dir, name string) {
	t.Helper()
	write(t, filepath.Join(dir, "specs/"+strings.ToLower(strings.ReplaceAll(name, " ", "-"))+".system.spec.md"), "---\ntype: system-spec\n---\n# "+name+"\n## Purpose\nHelp people.\n## Scope\nThe project.\n")
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
	path, err := InitTemplate(dir, "check-specifications", "brief")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Show(dir, "check-specifications")
	if err != nil || got.Slug != "check-specifications" || got.System != nil || got.Name != "New use case" {
		t.Fatalf("initialized use case: %+v, %v", got, err)
	}
	if _, err := InitTemplate(dir, "check-specifications", "brief"); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
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
	if err != nil || got.System == nil || got.System.Name != "Child" || got.BoundaryDescription != "The dotspec CLI." {
		t.Fatalf("association must not rewrite the stated boundary: %+v, %v", got, err)
	}
}

func TestValidationBeforeCreatingFiles(t *testing.T) {
	dir := t.TempDir()
	slug := "work.use-case.spec.md"
	if _, err := InitTemplate(dir, slug, "brief"); err == nil {
		t.Errorf("init accepted %q", slug)
	}
	if _, err := Show(dir, slug); err == nil {
		t.Errorf("show accepted %q", slug)
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid invocation wrote files: %v, %v", entries, err)
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
		_, err := Read(filepath.Join(dir, "specs", Layout.Collection, name))
		if !errors.As(err, &d) || d.Rule != "layout.use-case-name" {
			t.Fatalf("filename diagnostic: %v", err)
		}
	}
}
