package concept

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/speccollection"
)

// members lists concepts below dir as dotspec list concept does.
func members(dir string) ([]*speccollection.Member, error) {
	return speccollection.Members(dir, Layout, "concepts", func(path string) (string, error) {
		spec, err := Read(path)
		if err != nil {
			return "", err
		}
		return spec.Name, nil
	})
}

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

func TestInitializationAndAssociation(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	path, err := Init(dir, "job-map")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Read(path)
	if err != nil || got.Slug != "job-map" || got.Name != "Job map" {
		t.Fatalf("initialized concept: %+v, %v", got, err)
	}
	if listed, err := members(dir); err != nil || len(listed) != 1 || listed[0].System != nil {
		t.Fatalf("unassociated concept: %+v, %v", listed, err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Init(dir, "job-map"); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("overwrite: %v", err)
	}
	if content, err := os.ReadFile(path); err != nil || string(content) != string(before) {
		t.Fatalf("existing file changed: %s, %v", content, err)
	}
	write(t, path, strings.Replace(example, "# Job Map", "# Map an accomplishment", 1))
	got, err = Read(path)
	if err != nil || got.Slug != "job-map" || got.ID != "job-map" || got.Name != "Map an accomplishment" {
		t.Fatalf("identity changed with title: %+v, %v", got, err)
	}
	child := filepath.Join(dir, "child")
	writeSystem(t, child, "Child")
	local := writeConcept(t, child, "job-map")
	listed, err := members(child)
	if err != nil || len(listed) != 1 || listed[0].Path != local || listed[0].System == nil || listed[0].System.Name != "Child" {
		t.Fatalf("a subtree lists only its own concepts with their system: %+v, %v", listed, err)
	}
	if got, err := Read(local); err != nil || got.Purpose != "Understand an accomplishment and its success criteria." {
		t.Fatalf("association must not rewrite the purpose: %+v, %v", got, err)
	}
}

func TestValidationBeforeCreatingFiles(t *testing.T) {
	dir := t.TempDir()
	slug := "work.concept.spec.md"
	if _, err := Init(dir, slug); err == nil {
		t.Errorf("init accepted %q", slug)
	}
	if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
		t.Fatalf("invalid invocation wrote files: %v, %v", entries, err)
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
	if _, err := members(dir); err == nil {
		t.Fatal("list silently omitted malformed input")
	}
	if got, err := os.ReadFile(path); err != nil || !reflect.DeepEqual(got, []byte(source)) {
		t.Fatalf("read changed source: %q, %v", got, err)
	}
	for _, name := range []string{"Wrong.concept.spec.md", "work.md"} {
		_, err := Read(filepath.Join(dir, "specs", Layout.Collection, name))
		if !errors.As(err, &d) || d.Rule != "layout.concept-name" {
			t.Fatalf("filename diagnostic: %v", err)
		}
	}
}
