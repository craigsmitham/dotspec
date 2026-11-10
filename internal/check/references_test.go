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
)

func referenceCheck(dir string) check.Report {
	return check.Focused(dir, "rule", requirement.CheckType(), system.CheckType(), concept.CheckType(), job.CheckType(), quality.CheckType())
}

func TestCrossFileReferencesAndContext(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	path := requirement.Layout.Path(dir, "rule")
	source := requirement.Template + `
## References
[Heading](../concepts/inspection.concept.spec.md#actions)
[ID](../concepts/inspection.concept.spec.md#stable-inspection)
[Job step](../jobs/work.job.spec.md#prepare)
[Attribute](../system.spec.md#reliability)
[Native Markdown](../../Guide%20to%20rules.md#example-rule)
[Native contract](../../schema.json#/properties/name)
[External](https://example.invalid/missing#unknown)
[External malformed escape](https://example.invalid/missing%zz)
[Opaque URI](urn:example:missing%zz)
[Absolute path](/missing%zz)
[Protocol relative](//example.invalid/missing%zz)
[Directory](../../docs/)
`
	write(t, path, source)
	write(t, concept.Layout.Path(dir, "inspection"), strings.Replace(concept.Template, "type: concept", "type: concept\nid: stable-inspection", 1)+"\n[Not recursively checked](missing.md)\n")
	write(t, job.Layout.Path(dir, "work"), "---\ntype: job\n---\n# Work\n## Steps\n### Prepare\n**ID:** prepare\n")
	write(t, system.Layout.Path(dir, ""), validSystem+"\n## Quality attributes\n### Reliability\n**ID:** reliability\nConsistent operation.\n")
	write(t, filepath.Join(dir, "Guide to rules.md"), "---\nlayout: article\n---\n## Example rule\nA rule.\n[Not followed](missing.md)\n")
	write(t, filepath.Join(dir, "schema.json"), "{}")
	if err := os.Mkdir(filepath.Join(dir, "docs"), 0755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, ".gitignore"), "Guide to rules.md\nconcepts/\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	r := referenceCheck(dir)
	if r.ExitCode() != 0 || r.Summary.FilesChecked != 1 || r.Summary.ContextFilesChecked != 5 {
		t.Fatalf("reference context: %+v", r)
	}
	if after, _ := os.ReadFile(path); !reflect.DeepEqual(before, after) {
		t.Fatal("check changed source")
	}
	write(t, path, source+"\n[Missing](../concepts/inspection.concept.spec.md#missing)\n[Missing file](gone.md)\n[Missing file again](gone.md)\n")
	r = referenceCheck(dir)
	if !r.Complete || r.ExitCode() != 1 || r.Summary.Errors != 3 || r.Summary.ContextFilesChecked != 6 {
		t.Fatalf("invalid links: %+v", r)
	}
	for _, d := range r.Diagnostics {
		if d.Rule != "reference.unresolved" || d.Location.Path != path || d.Location.Line < 20 || len(d.Related) != 1 {
			t.Fatalf("reference location: %+v", d)
		}
	}
}

func TestMalformedAndNonDirectoryReferences(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	write(t, filepath.Join(dir, "file.txt"), "An ordinary file.\n")
	write(t, requirement.Layout.Path(dir, "rule"), requirement.Template+"\n[Malformed](bad%zz.md)\n[Not a directory](../../file.txt/missing.md)\n")
	r := referenceCheck(dir)
	if !r.Complete || r.Summary.Errors != 2 || r.Summary.ContextFilesChecked != 1 {
		t.Fatalf("reference errors: %+v", r)
	}
	if r.Diagnostics[0].Rule != "reference.invalid" || r.Diagnostics[1].Rule != "reference.unresolved" {
		t.Fatalf("reference error kinds: %+v", r.Diagnostics)
	}
}

func TestReferenceTargetErrorsAndCycles(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	a := requirement.Layout.Path(dir, "rule")
	b := requirement.Layout.Path(dir, "other")
	write(t, a, requirement.Template+"\n[Other](other.requirement.spec.md#statement)\n")
	write(t, b, requirement.Template+"\n[Back](rule.requirement.spec.md#statement)\n")
	r := checkProject(dir)
	if r.ExitCode() != 0 || r.Summary.FilesChecked != 2 || r.Summary.ContextFilesChecked != 0 {
		t.Fatalf("cycle: %+v", r)
	}
	write(t, b, "broken")
	r = checkProject(dir)
	if !r.Complete || r.Summary.Errors != 1 || r.Diagnostics[0].Rule != "metadata.frontmatter" {
		t.Fatalf("cascading target errors: %+v", r)
	}
	r = referenceCheck(dir)
	if !r.Complete || r.Summary.Errors != 1 || r.Summary.ContextFilesChecked != 1 {
		t.Fatalf("focused target error: %+v", r)
	}
}

func TestReferenceBoundaries(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	write(t, filepath.Join(dir, "nested/README.md"), "# Nested\n")
	git(t, filepath.Join(dir, "nested"), "init", "-q")
	write(t, requirement.Layout.Path(dir, "rule"), requirement.Template+"\n[Outside](../../../outside.md)\n[Nested](../../nested/README.md)\n[Metadata](../../.git/config)\n")
	r := referenceCheck(dir)
	if r.Summary.Errors != 3 || !r.Complete || r.Summary.ContextFilesChecked != 0 {
		t.Fatalf("boundaries: %+v", r)
	}
	for _, d := range r.Diagnostics {
		if d.Rule != "reference.boundary" {
			t.Fatal(d)
		}
	}
	if err := os.Symlink(filepath.Join(dir, "nested"), filepath.Join(dir, "shortcut")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	write(t, requirement.Layout.Path(dir, "rule"), requirement.Template+"\n[Symlink](../../shortcut/README.md)\n")
	r = referenceCheck(dir)
	if r.Summary.Errors != 1 || r.Diagnostics[0].Rule != "reference.boundary" {
		t.Fatalf("symlink: %+v", r)
	}
}

func TestRequirementsShareSystemAttributeContext(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	write(t, system.Layout.Path(dir, ""), validSystem+"\n## Quality attributes\n")
	write(t, requirement.Layout.Path(dir, "rule"), strings.Replace(requirement.Template, "type: requirement", "type: requirement\nattributes: [missing]", 1))
	r := referenceCheck(dir)
	if r.ExitCode() != 0 || r.Summary.Advice != 1 || r.Diagnostics[0].Rule != "requirement.undefined-attribute" || r.Diagnostics[0].Location.Line != 3 {
		t.Fatalf("attribute advice: %+v", r)
	}
	write(t, system.Layout.Path(dir, ""), "broken")
	r = referenceCheck(dir)
	if r.Summary.Errors != 1 || r.Summary.Advice != 0 {
		t.Fatalf("invalid context cascaded: %+v", r)
	}
}

func TestReferencesDoNotRepeatFailedContextReads(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	write(t, requirement.Layout.Path(dir, "rule"), requirement.Template)
	if err := os.Mkdir(system.Layout.Path(dir, ""), 0755); err != nil {
		t.Fatal(err)
	}
	r := referenceCheck(dir)
	if r.Complete || r.Summary.Errors != 1 || r.Summary.ContextFilesChecked != 1 || r.Diagnostics[0].Rule != "scan.context" {
		t.Fatalf("repeated unavailable context: %+v", r)
	}
}
