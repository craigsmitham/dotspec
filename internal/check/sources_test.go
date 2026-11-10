package check_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
)

func sourceReport(dir string) check.SourceReport {
	return check.Sources(dir, system.CheckType(), job.CheckType(), usecase.CheckType(), concept.CheckType(), quality.CheckType(), requirement.CheckType())
}

func sourceMap(r check.SourceReport) map[string]check.SourceFile {
	files := make(map[string]check.SourceFile)
	for _, file := range r.Files {
		files[file.Path] = file
	}
	return files
}

func TestSourcesCollectLinksFromEveryTypeWithoutInterpretingNativeFiles(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	native := filepath.Join(dir, "contract.tla")
	guide := filepath.Join(dir, "Guide to tests.md")
	write(t, native, "native\x00content\n[Not Markdown](missing.md)\n")
	write(t, guide, "# Guide\n[No recursive discovery](unlisted.txt)\n")
	write(t, filepath.Join(dir, "unlisted.txt"), "Must stay outside the inventory.\n")
	write(t, filepath.Join(dir, ".gitignore"), "contract.tla\nGuide to tests.md\n")
	specs := map[string]string{
		system.Layout.Path(dir, ""):          system.Template,
		job.Layout.Path(dir, "work"):         job.Template,
		usecase.Layout.Path(dir, "work"):     usecase.Template,
		concept.Layout.Path(dir, "work"):     concept.Template,
		quality.Layout.Path(dir, "work"):     quality.Template,
		requirement.Layout.Path(dir, "work"): requirement.Template,
	}
	for path, content := range specs {
		rel, err := filepath.Rel(filepath.Dir(path), native)
		if err != nil {
			t.Fatal(err)
		}
		content += "\n## Related material\n[Native source](" + filepath.ToSlash(rel) + "#native-fragment)\n"
		write(t, path, content)
		specs[path] = content
	}
	path := concept.Layout.Path(dir, "work")
	specs[path] += "[Another reference](../../contract.tla)\n[Guide](../../Guide%20to%20tests.md#guide)\n[Directory](../../)\n[Remote](https://example.invalid/file)\n`[Code](absent.txt)`\n![Image](absent.png)\n"
	write(t, path, specs[path])
	r := sourceReport(dir)
	if r.ExitCode() != 0 || len(r.Files) != 8 || !r.Complete {
		t.Fatalf("sources: %+v", r)
	}
	files := sourceMap(r)
	if len(files[native].ReferencedBy) != 7 || len(files[guide].ReferencedBy) != 1 {
		t.Fatalf("lost incoming references: %+v", files)
	}
	for _, location := range files[native].ReferencedBy {
		if !filepath.IsAbs(location.Path) || location.Line == 0 || location.Column == 0 || location.Element == "" {
			t.Fatalf("location: %+v", location)
		}
	}
	if _, found := files[filepath.Join(dir, "unlisted.txt")]; found {
		t.Fatal("followed referenced Markdown recursively")
	}
	if next := sourceReport(dir); !reflect.DeepEqual(r, next) {
		t.Fatal("inventory not deterministic")
	}
	if report := checkProject(dir); !reflect.DeepEqual(r.Diagnostics, report.Diagnostics) || r.ExitCode() != report.ExitCode() {
		t.Fatal("sources and check disagree")
	}
	for path, content := range specs {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != content {
			t.Fatalf("source changed: %s: %v", path, err)
		}
	}
	data, err := os.ReadFile(native)
	if err != nil || string(data) != "native\x00content\n[Not Markdown](missing.md)\n" {
		t.Fatal("native file changed", err)
	}
}

func TestSourcesIncludeAncestorContextAndKeepBoundaryFailuresVisible(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "module")
	write(t, system.Layout.Path(dir, ""), system.Template+"\n[Contract](../contract.txt)\n")
	write(t, filepath.Join(dir, "contract.txt"), "A native source.\n")
	path := concept.Layout.Path(child, "work")
	write(t, path, concept.Template)
	r := sourceReport(child)
	if r.ExitCode() != 0 || len(r.Files) != 3 || len(sourceMap(r)[filepath.Join(dir, "contract.txt")].ReferencedBy) != 1 {
		t.Fatalf("ancestor context: %+v", r)
	}
	write(t, path, concept.Template+"\n[Missing](missing.tla)\n[Boundary](../../../../outside.txt)\n")
	r = sourceReport(child)
	if r.ExitCode() != 1 || !r.Complete || len(r.Files) != 3 {
		t.Fatalf("failures: %+v", r)
	}
	if report := checkProject(child); !reflect.DeepEqual(r.Diagnostics, report.Diagnostics) {
		t.Fatal("different resolution policy")
	}
	if err := os.Symlink(filepath.Join(dir, "contract.txt"), filepath.Join(child, "shortcut.txt")); err == nil {
		write(t, path, concept.Template+"\n[Symlink](../../shortcut.txt)\n")
		r = sourceReport(child)
		if r.ExitCode() != 1 || len(r.Files) != 3 || r.Diagnostics[0].Rule != "reference.boundary" {
			t.Fatalf("symlink: %+v", r)
		}
	}
}

func TestSourcesEmptyInvalidAndIncomplete(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	if r := sourceReport(dir); r.ExitCode() != 0 || r.Files == nil || len(r.Files) != 0 {
		t.Fatalf("empty inventory: %+v", r)
	}
	path := requirement.Layout.Path(dir, "broken")
	write(t, path, "broken")
	if r := sourceReport(dir); r.ExitCode() != 1 || !r.Complete || len(r.Files) != 1 {
		t.Fatalf("invalid source: %+v", r)
	}
	write(t, path, requirement.Template)
	if err := os.Mkdir(system.Layout.Path(dir, ""), 0755); err != nil {
		t.Fatal(err)
	}
	if r := sourceReport(dir); r.ExitCode() != 2 || r.Complete || len(r.Files) != 1 {
		t.Fatalf("unreadable context directory: %+v", r)
	}
	if r := sourceReport(filepath.Join(dir, "absent")); r.ExitCode() != 2 || r.Complete || len(r.Files) != 0 || !strings.Contains(r.Diagnostics[0].Message, "absent") {
		t.Fatalf("missing scope: %+v", r)
	}
}
