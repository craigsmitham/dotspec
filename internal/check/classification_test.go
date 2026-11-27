package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
)

func TestClassificationContextAndAdvice(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	modelPath := filepath.Join(dir, "specs/product.system.spec.md")
	attrPath := filepath.Join(dir, "specs/quality/known/known.quality-attribute.spec.md")
	attr := "---\ntype: quality-attribute-spec\n---\n# Known\n"
	write(t, attrPath, attr)
	model := validSystem
	write(t, modelPath, model)
	source := strings.Replace(requirement.QualityTemplate, "type: requirement-spec", "type: requirement-spec\nquality_attributes: [known, unknown]", 1)
	path := requirement.Layout.Path(child, "work")
	write(t, path, source)
	focused := check.Focused(child, "work", requirement.CheckType())
	checkQuality := func(path string) check.Report {
		return check.Project(path, system.CheckType(), system.AttributeCheckType(), requirement.CheckType())
	}
	project := checkQuality(child)
	if !reflect.DeepEqual(focused.Diagnostics, project.Diagnostics) || focused.ExitCode() != 0 || focused.Summary.Advice != 2 || focused.Summary.ContextFilesChecked != 2 {
		t.Fatalf("context advice differs: %+v %+v", focused, project)
	}
	// A system discovered directly and also used by two requirement owners
	// contributes its missing-definition advice only once.
	write(t, requirement.Layout.Path(dir, "root"), source)
	project = checkQuality(dir)
	if project.Summary.Advice != 3 || project.Summary.ContextFilesChecked != 0 {
		t.Fatalf("duplicate model advice: %+v", project)
	}
	write(t, attrPath, strings.Replace(attr, "type: quality-attribute-spec", "type: quality-attribute-spec\nsubattribute_of: [missing]", 1))
	project = checkQuality(dir)
	if !project.Complete || project.Summary.Errors != 1 || project.Summary.Advice != 0 || project.Diagnostics[0].Rule != "quality.subattribute" {
		t.Fatalf("cascading model errors: %+v", project)
	}
	// Failure to locate a valid context does not fabricate undefined labels.
	write(t, modelPath, "broken")
	focused = check.Focused(child, "work", requirement.CheckType())
	if !focused.Complete || focused.ExitCode() != 1 || focused.Summary.Errors != 1 || focused.Summary.Advice != 0 {
		t.Fatalf("broken context: %+v", focused)
	}
	if err := os.Remove(modelPath); err != nil {
		t.Fatal(err)
	}
	focused = check.Focused(child, "work", requirement.CheckType())
	if focused.ExitCode() != 0 || focused.Summary.Advice != 0 || focused.Summary.ContextFilesChecked != 0 {
		t.Fatalf("absent model: %+v", focused)
	}
	missing := check.Focused(child, "missing", requirement.CheckType())
	if missing.Complete || missing.ExitCode() != 2 || missing.Summary.FilesChecked != 1 {
		t.Fatalf("missing input: %+v", missing)
	}
}

func TestClassifiedRequirementDiscoveryAndSourceLocations(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	path := requirement.Layout.Path(dir, "tracked")
	write(t, path, requirement.QualityTemplate)
	git(t, dir, "add", "specs/requirements/tracked.requirement.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "tracked.requirement.spec.md\nignored/\n*.system.spec.md\n")
	write(t, filepath.Join(dir, "ignored/broken.requirement.spec.md"), "ignored")
	write(t, filepath.Join(dir, "specs/product.system.spec.md"), validSystem)
	r := checkProject(dir)
	if r.ExitCode() != 0 || r.Summary.FilesChecked != 1 || r.Summary.ContextFilesChecked != 1 {
		t.Fatalf("discovery or ignored context: %+v", r)
	}
	write(t, path, "---\ntype: requirement-spec\nquality_attributes:\n  - known\n  - missing\n---\n# Draft\n")
	r = checkProject(dir)
	var lines []int
	for _, d := range r.Diagnostics {
		if d.Rule == "requirement.undefined-attribute" {
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
