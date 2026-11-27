package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/concept"
	"github.com/dotspec-md/dotspec/internal/usecase"
)

// Each specification collection kind checks end to end through check.Focused
// and check.Project.
var collectionKinds = []struct {
	kind, dir, section, template string
	typ                          func() check.Type
	advice                       int
}{
	{kind: "concept", dir: "concepts", section: "Purpose", template: concept.Template, typ: concept.CheckType, advice: 3},
	{kind: "use-case", dir: "use-cases", section: "Goal", template: usecase.ExpandedTemplate, typ: usecase.CheckType, advice: 4},
}

func TestCollectionProjectAndFocusedReports(t *testing.T) {
	for _, k := range collectionKinds {
		t.Run(k.kind, func(t *testing.T) {
			dir := t.TempDir()
			git(t, dir, "init", "-q")
			child := filepath.Join(dir, "component")
			specs := filepath.Join(child, "specs", k.dir)
			write(t, filepath.Join(dir, "specs/product.system.spec.md"), validSystem)
			write(t, filepath.Join(specs, "work."+k.kind+".spec.md"), "---\ntype: "+k.kind+"-spec\n---\n# Work\n## "+k.section+"\nInspect findings.\n")
			focused, project := check.Focused(child, "work", k.typ()), checkProject(child)
			if !reflect.DeepEqual(focused.Diagnostics, project.Diagnostics) || !reflect.DeepEqual(focused.ContextFiles, project.ContextFiles) {
				t.Fatalf("focused and project differ:\n%+v\n%+v", focused, project)
			}
			if focused.ExitCode() != 0 || !focused.Complete || focused.Summary.FilesChecked != 1 || focused.Summary.ContextFilesChecked != 1 || focused.Summary.Advice != k.advice {
				t.Fatalf("focused accounting: %+v", focused)
			}
			if !reflect.DeepEqual(focused.Scope.Types, []string{k.kind}) || focused.Scope.Discovery != "direct" || !reflect.DeepEqual(project.Scope.Types, []string{"system", "job", "use-case", "concept", "requirement"}) {
				t.Fatalf("scope: %+v; %+v", focused.Scope, project.Scope)
			}
			// A job and this kind sharing an ancestor report that context only once.
			write(t, filepath.Join(child, "specs/jobs/work.job.spec.md"), validJob)
			write(t, filepath.Join(dir, "specs/product.system.spec.md"), "broken shared context")
			project = checkProject(child)
			if project.ExitCode() != 1 || !project.Complete || project.Summary.Errors != 1 || project.Summary.FilesChecked != 2 || project.Summary.ContextFilesChecked != 1 {
				t.Fatalf("shared malformed context: %+v", project)
			}
			// The focused check does not inspect malformed sibling specifications.
			write(t, filepath.Join(dir, "specs/product.system.spec.md"), validSystem)
			write(t, filepath.Join(specs, "broken."+k.kind+".spec.md"), "broken sibling")
			if focused := check.Focused(child, "work", k.typ()); focused.ExitCode() != 0 || focused.Summary.FilesChecked != 1 {
				t.Fatalf("focused check scanned a sibling: %+v", focused)
			}
		})
	}
}

func TestCollectionFailuresAndMissingInput(t *testing.T) {
	for _, k := range collectionKinds {
		t.Run(k.kind, func(t *testing.T) {
			dir := t.TempDir()
			misplaced := "misplaced." + k.kind + ".spec.md"
			for name, source := range map[string]string{
				"specs/" + k.dir + "/broken." + k.kind + ".spec.md": "broken",
				"specs/" + k.dir + "/valid." + k.kind + ".spec.md":  k.template,
				misplaced:                      k.template,
				"specs/jobs/valid.job.spec.md": validJob,
			} {
				write(t, filepath.Join(dir, name), source)
			}
			r := checkProject(dir)
			if !r.Complete || r.ExitCode() != 1 || r.Summary.FilesChecked != 4 || r.Summary.Errors != 1 || r.Summary.Advice != 0 {
				t.Fatalf("independent checks did not continue: %+v", r)
			}
			if r.Diagnostics[0].Rule != "metadata.frontmatter" {
				t.Fatalf("unexpected findings: %+v", r.Diagnostics)
			}
			if err := os.Remove(filepath.Join(dir, misplaced)); err != nil {
				t.Fatal(err)
			}
			missing := check.Focused(dir, "absent", k.typ())
			if missing.Complete || missing.ExitCode() != 2 || missing.Summary.FilesChecked != 1 || missing.Diagnostics[0].Rule != "scan.read" {
				t.Fatalf("missing input: %+v", missing)
			}
			broken := check.Focused(dir, "broken", k.typ())
			if !broken.Complete || broken.ExitCode() != 1 || broken.Summary.Advice != 0 {
				t.Fatalf("malformed file must skip dependent checks: %+v", broken)
			}
		})
	}
}

func TestCollectionCheckVisibilityAndIgnoredContext(t *testing.T) {
	for _, k := range collectionKinds {
		t.Run(k.kind, func(t *testing.T) {
			dir := t.TempDir()
			git(t, dir, "init", "-q")
			tracked := "specs/" + k.dir + "/tracked." + k.kind + ".spec.md"
			write(t, filepath.Join(dir, tracked), k.template)
			git(t, dir, "add", tracked)
			write(t, filepath.Join(dir, ".gitignore"), "tracked."+k.kind+".spec.md\nignored/\n*.system.spec.md\n")
			write(t, filepath.Join(dir, "ignored/bad."+k.kind+".spec.md"), "ignored")
			write(t, filepath.Join(dir, "specs/product.system.spec.md"), validSystem)
			r := checkProject(dir)
			if r.ExitCode() != 0 || r.Scope.Discovery != "git" || r.Summary.FilesChecked != 1 || r.Summary.ContextFilesChecked != 1 {
				t.Fatalf("tracked file or ignored ancestor context lost: %+v", r)
			}
		})
	}
}
