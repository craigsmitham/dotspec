package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/craigsmitham/dotspec/internal/usecase"
)

func TestUseCaseProjectAndFocusedReports(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	write(t, filepath.Join(dir, "specs/system.spec.md"), validSystem)
	write(t, filepath.Join(child, "specs/use-cases/work.use-case.spec.md"), "---\ntype: use-case\n---\n# Work\n## Goal\nInspect findings.\n")
	focused, project := checkUseCase(child, "work"), checkProject(child)
	if !reflect.DeepEqual(focused.Diagnostics, project.Diagnostics) || !reflect.DeepEqual(focused.ContextFiles, project.ContextFiles) {
		t.Fatalf("focused and project differ:\n%+v\n%+v", focused, project)
	}
	if focused.ExitCode() != 0 || !focused.Complete || focused.Summary.FilesChecked != 1 || focused.Summary.ContextFilesChecked != 1 || focused.Summary.Advice != 4 {
		t.Fatalf("focused accounting: %+v", focused)
	}
	if !reflect.DeepEqual(focused.Scope.Types, []string{"use-case"}) || focused.Scope.Discovery != "direct" || !reflect.DeepEqual(project.Scope.Types, []string{"system", "job", "use-case", "concept", "quality", "requirement"}) {
		t.Fatalf("scope: %+v; %+v", focused.Scope, project.Scope)
	}
	// A job and use case sharing an ancestor report that context only once.
	write(t, filepath.Join(child, "specs/jobs/work.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs/system.spec.md"), "broken shared context")
	project = checkProject(child)
	if project.ExitCode() != 1 || !project.Complete || project.Summary.Errors != 1 || project.Summary.FilesChecked != 2 || project.Summary.ContextFilesChecked != 1 {
		t.Fatalf("shared malformed context: %+v", project)
	}
	// The focused check does not inspect malformed sibling specifications.
	write(t, filepath.Join(dir, "specs/system.spec.md"), validSystem)
	write(t, filepath.Join(child, "specs/use-cases/broken.use-case.spec.md"), "broken sibling")
	if focused := checkUseCase(child, "work"); focused.ExitCode() != 0 || focused.Summary.FilesChecked != 1 {
		t.Fatalf("focused check scanned a sibling: %+v", focused)
	}
}

func TestUseCaseFailuresAndMissingInput(t *testing.T) {
	dir := t.TempDir()
	for name, source := range map[string]string{
		"specs/use-cases/broken.use-case.spec.md": "broken",
		"specs/use-cases/valid.use-case.spec.md":  usecase.ExpandedTemplate,
		"misplaced.use-case.spec.md":              usecase.ExpandedTemplate,
		"specs/jobs/valid.job.spec.md":            validJob,
	} {
		write(t, filepath.Join(dir, name), source)
	}
	r := checkProject(dir)
	if !r.Complete || r.ExitCode() != 1 || r.Summary.FilesChecked != 4 || r.Summary.Errors != 2 || r.Summary.Advice != 0 {
		t.Fatalf("independent checks did not continue: %+v", r)
	}
	if r.Diagnostics[0].Rule != "layout.use-case-location" || r.Diagnostics[1].Rule != "metadata.frontmatter" {
		t.Fatalf("unexpected findings: %+v", r.Diagnostics)
	}
	if err := os.Remove(filepath.Join(dir, "misplaced.use-case.spec.md")); err != nil {
		t.Fatal(err)
	}
	missing := checkUseCase(dir, "absent")
	if missing.Complete || missing.ExitCode() != 2 || missing.Summary.FilesChecked != 1 || missing.Diagnostics[0].Rule != "scan.read" {
		t.Fatalf("missing input: %+v", missing)
	}
	broken := checkUseCase(dir, "broken")
	if !broken.Complete || broken.ExitCode() != 1 || broken.Summary.Advice != 0 {
		t.Fatalf("malformed file must skip dependent checks: %+v", broken)
	}
}

func TestUseCaseCheckVisibilityAndIgnoredContext(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	write(t, filepath.Join(dir, "specs/use-cases/tracked.use-case.spec.md"), usecase.ExpandedTemplate)
	git(t, dir, "add", "specs/use-cases/tracked.use-case.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "tracked.use-case.spec.md\nignored/\nsystem.spec.md\n")
	write(t, filepath.Join(dir, "ignored/bad.use-case.spec.md"), "ignored")
	write(t, filepath.Join(dir, "specs/system.spec.md"), validSystem)
	r := checkProject(dir)
	if r.ExitCode() != 0 || r.Scope.Discovery != "git" || r.Summary.FilesChecked != 1 || r.Summary.ContextFilesChecked != 1 {
		t.Fatalf("tracked file or ignored ancestor context lost: %+v", r)
	}
}
