package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/quality"
)

func TestQualityContextAndAdvice(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	modelPath := filepath.Join(dir, "specs/system.spec.md")
	model := validSystem + "\n## Quality attributes\n### Known\n**ID:** known\n"
	write(t, modelPath, model)
	source := strings.Replace(quality.Template, "type: quality", "type: quality\nquality_attributes: [known, unknown]", 1)
	path := quality.Layout.Path(child, "work")
	write(t, path, source)
	focused := check.Focused(child, "work", quality.CheckType())
	project := checkProject(child)
	if !reflect.DeepEqual(focused.Diagnostics, project.Diagnostics) || focused.ExitCode() != 0 || focused.Summary.Advice != 2 || focused.Summary.ContextFilesChecked != 1 {
		t.Fatalf("context advice differs: %+v %+v", focused, project)
	}
	// A system discovered directly and also used by two requirement owners
	// contributes its missing-definition advice only once.
	write(t, quality.Layout.Path(dir, "root"), source)
	project = checkProject(dir)
	if project.Summary.Advice != 3 || project.Summary.ContextFilesChecked != 0 {
		t.Fatalf("duplicate model advice: %+v", project)
	}
	write(t, modelPath, strings.Replace(model, "**ID:** known", "**ID:** known\n**Sub-attribute of:** missing", 1))
	project = checkProject(dir)
	if !project.Complete || project.Summary.Errors != 1 || project.Summary.Advice != 0 || project.Diagnostics[0].Rule != "quality.subattribute" {
		t.Fatalf("cascading model errors: %+v", project)
	}
	// Failure to locate a valid context does not fabricate undefined labels.
	write(t, modelPath, "broken")
	focused = check.Focused(child, "work", quality.CheckType())
	if !focused.Complete || focused.ExitCode() != 1 || focused.Summary.Errors != 1 || focused.Summary.Advice != 0 {
		t.Fatalf("broken context: %+v", focused)
	}
	if err := os.Remove(modelPath); err != nil {
		t.Fatal(err)
	}
	focused = check.Focused(child, "work", quality.CheckType())
	if focused.ExitCode() != 0 || focused.Summary.Advice != 0 || focused.Summary.ContextFilesChecked != 0 {
		t.Fatalf("absent model: %+v", focused)
	}
	missing := check.Focused(child, "missing", quality.CheckType())
	if missing.Complete || missing.ExitCode() != 2 || missing.Summary.FilesChecked != 1 {
		t.Fatalf("missing input: %+v", missing)
	}
}

func TestQualityDiscoveryAndSourceLocations(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	path := quality.Layout.Path(dir, "tracked")
	write(t, path, quality.Template)
	git(t, dir, "add", "specs/quality/tracked.quality.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "tracked.quality.spec.md\nignored/\nsystem.spec.md\n")
	write(t, filepath.Join(dir, "ignored/broken.quality.spec.md"), "ignored")
	write(t, filepath.Join(dir, "specs/system.spec.md"), validSystem+"\n## Quality attributes\n")
	r := checkProject(dir)
	if r.ExitCode() != 0 || r.Summary.FilesChecked != 1 || r.Summary.ContextFilesChecked != 1 {
		t.Fatalf("discovery or ignored context: %+v", r)
	}
	write(t, path, "---\ntype: quality\nquality_attributes:\n  - known\n  - missing\n---\n# Draft\n")
	r = checkProject(dir)
	var lines []int
	for _, d := range r.Diagnostics {
		if d.Rule == "quality.undefined-attribute" {
			lines = append(lines, d.Location.Line)
		}
	}
	if !reflect.DeepEqual(lines, []int{4, 5}) {
		t.Fatalf("attribute locations: %+v", r.Diagnostics)
	}
	if after, _ := os.ReadFile(path); !strings.Contains(string(after), "  - missing") {
		t.Fatal("check rewrote labels")
	}
}
