package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const validJob = "---\ntype: job\n---\n# Prepare a change\n\n## Job\nPrepare a change for review.\n\n## Performer\nA reviewer.\n"
const validSystem = "---\ntype: system\n---\n# Product\n\n## Purpose\nHelp people.\n\n## Scope\nThe product.\n"

func write(t *testing.T, path, source string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	args = append([]string{"-C", dir}, args...)
	if output, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %s: %v", args, output, err)
	}
}

func TestProjectContinuesAcrossMalformedFiles(t *testing.T) {
	dir := t.TempDir()
	sources := map[string]string{
		"adjacent.job.spec.md":                 validJob,
		"specs/system.spec.md":                 "---\ntype: [\n---\n# Broken\n",
		"specs/jobs/duplicate.job.spec.md":     "---\ntype: job\nid: shared\n---\n# Duplicate\n\n## Steps\n\n### Gather\n**ID:** shared\n",
		"specs/jobs/advice.job.spec.md":        "---\ntype: job\n---\n# Draft\n",
		"specs/jobs/valid.job.spec.md":         validJob,
		"specs/terms/unsupported.term.spec.md": "ignored unsupported type",
	}
	for path, source := range sources {
		write(t, filepath.Join(dir, path), source)
	}
	r := Project(dir)
	if !r.Complete || r.ExitCode() != 1 || r.Summary.FilesDiscovered != 5 || r.Summary.FilesChecked != 5 || r.Summary.Errors != 3 || r.Summary.Advice != 2 {
		t.Fatalf("report: %+v", r)
	}
	rules := map[string]int{}
	for _, d := range r.Diagnostics {
		rules[d.Rule]++
		if d.Rule == "identity.duplicate" {
			if d.Location.Line != 10 || d.Location.Column != 1 || len(d.Related) != 1 || d.Related[0].Line != 3 || d.Related[0].Column != 5 || d.Related[0].Path != d.Location.Path {
				t.Fatalf("duplicate must locate both definitions: %+v", d)
			}
		}
	}
	want := map[string]int{"layout.job-location": 1, "metadata.invalid": 1, "identity.duplicate": 1, "job.statement": 1, "job.performer": 1}
	if !reflect.DeepEqual(rules, want) {
		t.Fatalf("unexpected or cascading findings: %v", rules)
	}
	if next := Project(dir); !reflect.DeepEqual(r, next) {
		t.Fatal("unchanged inputs produced a different report")
	}
	for path, source := range sources {
		data, err := os.ReadFile(filepath.Join(dir, path))
		if err != nil || string(data) != source {
			t.Fatalf("check changed %s: %v", path, err)
		}
	}
}

func TestSubtreeReportsAncestorContext(t *testing.T) {
	dir := t.TempDir()
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	git(t, dir, "init", "-q")
	child := filepath.Join(dir, "component")
	write(t, filepath.Join(child, "specs/jobs/one.job.spec.md"), validJob)
	write(t, filepath.Join(child, "specs/jobs/two.job.spec.md"), validJob)
	ancestor := filepath.Join(dir, "specs/system.spec.md")
	write(t, ancestor, validSystem)
	for _, broken := range []bool{false, true} {
		if broken {
			write(t, ancestor, "not a spec")
		}
		r := Project(child)
		if !r.Complete || r.Summary.FilesChecked != 2 || r.Summary.ContextFilesChecked != 1 || !reflect.DeepEqual(r.ContextFiles, []string{ancestor}) {
			t.Fatalf("subtree context: %+v", r)
		}
		if broken && (r.ExitCode() != 1 || len(r.Diagnostics) != 1 || r.Diagnostics[0].Location.Path != ancestor) {
			t.Fatalf("malformed shared context should appear once: %+v", r)
		}
		if !broken && r.ExitCode() != 0 {
			t.Fatalf("valid context: %+v", r)
		}
	}
}

func TestDiscoveryUsesGitVisibility(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	write(t, filepath.Join(dir, "specs/jobs/valid.job.spec.md"), validJob)
	write(t, filepath.Join(dir, "specs/jobs/tracked.job.spec.md"), "broken tracked file")
	git(t, dir, "add", "specs/jobs/tracked.job.spec.md")
	write(t, filepath.Join(dir, ".gitignore"), "ignored/\ntracked.job.spec.md\n")
	write(t, filepath.Join(dir, "ignored/broken.job.spec.md"), "ignored")
	nested := filepath.Join(dir, "nested")
	write(t, filepath.Join(nested, "bad.job.spec.md"), "nested repository")
	git(t, nested, "init", "-q")
	outside := t.TempDir()
	write(t, filepath.Join(outside, "bad.job.spec.md"), "outside the project")
	if err := os.Symlink(outside, filepath.Join(dir, "linked")); err != nil {
		t.Fatal(err)
	}
	r := Project(dir)
	if r.Scope.Discovery != "git" || !r.Complete || r.Summary.FilesChecked != 2 || r.Summary.Errors != 1 || r.Diagnostics[0].Rule != "metadata.frontmatter" {
		t.Fatalf("Git visibility changed: %+v", r)
	}
	plain := t.TempDir()
	write(t, filepath.Join(plain, ".gitignore"), "ignored/\n")
	write(t, filepath.Join(plain, "ignored/bad.job.spec.md"), "not ignored outside Git")
	if r := Project(plain); r.Scope.Discovery != "filesystem" || r.Summary.FilesChecked != 1 || r.ExitCode() != 1 {
		t.Fatalf("filesystem visibility changed: %+v", r)
	}
}

func TestEmptyAndAdvisoryOnlyProjectsSucceed(t *testing.T) {
	dir := t.TempDir()
	r := Project(dir)
	if !r.Complete || r.ExitCode() != 0 || r.Summary.FilesChecked != 0 || r.Diagnostics == nil {
		t.Fatalf("empty project: %+v", r)
	}
	write(t, filepath.Join(dir, "specs/jobs/draft.job.spec.md"), "---\ntype: job\n---\n# Draft\n")
	r = Project(dir)
	if !r.Complete || r.ExitCode() != 0 || r.Summary.Errors != 0 || r.Summary.Advice != 2 {
		t.Fatalf("advisory-only project: %+v", r)
	}
}

func TestIncompleteScanIsNotAValidationFailure(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	t.Setenv("PATH", t.TempDir()) // Git discovery cannot run.
	r := Project(dir)
	if r.Complete || r.ExitCode() != 2 || len(r.Diagnostics) != 1 || r.Diagnostics[0].Rule != "scan.discovery" {
		t.Fatalf("failed discovery must not look like an empty successful scan: %+v", r)
	}
	for _, path := range []string{filepath.Join(dir, "missing"), filepath.Join(dir, "specs")} {
		if strings.HasSuffix(path, "specs") {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		r = Project(path)
		if r.Complete || r.ExitCode() != 2 || r.Diagnostics[0].Rule != "scan.scope" {
			t.Fatalf("invalid scope: %+v", r)
		}
	}
}

func TestUnreadableSpecDoesNotPreventOtherChecks(t *testing.T) {
	dir := t.TempDir()
	unreadable := filepath.Join(dir, "specs/jobs/a-private.job.spec.md")
	write(t, unreadable, validJob)
	write(t, filepath.Join(dir, "specs/jobs/z-draft.job.spec.md"), "---\ntype: job\n---\n# Draft\n")
	if err := os.Chmod(unreadable, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(unreadable, 0o600) })
	if _, err := os.ReadFile(unreadable); err == nil {
		t.Skip("this user can read files despite permissions")
	}
	r := Project(dir)
	if r.Complete || r.ExitCode() != 2 || r.Summary.FilesChecked != 2 || r.Summary.Advice != 2 || r.Summary.Errors != 1 {
		t.Fatalf("file access failure prevented independent checks: %+v", r)
	}
	if r.Diagnostics[0].Rule != "scan.read" || !strings.HasSuffix(r.Diagnostics[0].Location.Path, "a-private.job.spec.md") {
		t.Fatalf("file access failure lost its source: %+v", r)
	}
}
