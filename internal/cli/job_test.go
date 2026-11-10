package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobWorkflow(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if code, _, errOut := invoke("system", "init"); code != 0 {
		t.Fatal(errOut)
	}
	for _, slug := range []string{"develop-requirements", "manage-requirements"} {
		code, out, errOut := invoke("job", "init", slug)
		if code != 0 || errOut != "" || !strings.Contains(out, "Created specs/jobs/"+slug+".job.spec.md") {
			t.Fatalf("init: %d %q %q", code, out, errOut)
		}
	}
	code, out, errOut := invoke("job", "show", "develop-requirements")
	for _, want := range []string{"New job", "Situation", "Motivation", "Intended progress", "System: New system (specs/system.spec.md)", "Source: specs/jobs/develop-requirements.job.spec.md"} {
		if code != 0 || errOut != "" || !strings.Contains(out, want) {
			t.Fatalf("show missing %q: %d %q %q", want, code, out, errOut)
		}
	}
	if strings.Contains(out, "## Desired outcome") {
		t.Fatalf("new template uses the legacy heading: %s", out)
	}
	code, out, errOut = invoke("job", "list")
	want := "Jobs\n\nSystem: New system (specs/system.spec.md)\n  New job  (specs/jobs/develop-requirements.job.spec.md)\n  New job  (specs/jobs/manage-requirements.job.spec.md)\n"
	if code != 0 || out != want || errOut != "" {
		t.Fatalf("list: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("job", "init", "develop-requirements")
	if code == 0 || out != "" || !strings.Contains(errOut, "refusing to overwrite") {
		t.Fatalf("overwrite: %d %q %q", code, out, errOut)
	}
}

func TestJobFlagsBeforeAndAfterSlug(t *testing.T) {
	dir := t.TempDir()
	if code, _, errOut := invoke("job", "init", "one", "--path", dir); code != 0 {
		t.Fatal(errOut)
	}
	if code, _, errOut := invoke("job", "init", "--path", dir, "two"); code != 0 {
		t.Fatal(errOut)
	}
	for _, args := range [][]string{
		{"job", "show", "one", "--path", dir},
		{"job", "show", "--path", dir, "two"},
		{"job", "show", "one", "--path=" + dir},
		{"job", "show", "--path=" + dir, "--", "two"},
	} {
		if code, out, errOut := invoke(args...); code != 0 || !strings.Contains(out, "New job") || errOut != "" {
			t.Fatalf("flags %v: %d %q %q", args, code, out, errOut)
		}
	}
}

func TestJobHelpAndUsageErrors(t *testing.T) {
	for _, args := range [][]string{{"job"}, {"job", "--help"}, {"job", "init", "--help"}, {"job", "show", "slug", "--help"}, {"job", "list", "--help"}} {
		code, out, errOut := invoke(args...)
		if code != 0 || !strings.Contains(out, "Usage: dotspec job") || errOut != "" {
			t.Errorf("help %v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{
		{"job", "unknown"}, {"job", "init"}, {"job", "show"}, {"job", "list", "slug"},
		{"job", "init", "one", "two"}, {"job", "show", "../escape"}, {"job", "init", "Upper"},
		{"job", "show", "slug", "--unknown"}, {"job", "init", "slug", "--path"},
		{"job", "init", "slug", "--path", ""}, {"job", "init", "--", "slug", "--path", "."},
	} {
		code, out, errOut := invoke(args...)
		if code == 0 || out != "" || !strings.HasPrefix(errOut, "dotspec:") {
			t.Errorf("usage %v: %d %q %q", args, code, out, errOut)
		}
	}
}

func TestJobEmptyMissingUnassociatedAndMalformed(t *testing.T) {
	dir := t.TempDir()
	code, out, errOut := invoke("job", "list", "--path", dir)
	if code != 0 || !strings.Contains(out, "No jobs found") || errOut != "" {
		t.Fatalf("empty: %d %q %q", code, out, errOut)
	}
	code, out, errOut = invoke("job", "show", "missing", "--path", dir)
	if code == 0 || out != "" || !strings.Contains(errOut, "missing.job.spec.md") {
		t.Fatalf("missing: %d %q %q", code, out, errOut)
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if code, _, errOut := invoke("job", "init", "work", "--path", dir); code != 0 {
		t.Fatal(errOut)
	}
	code, out, errOut = invoke("job", "show", "work", "--path", dir)
	if code != 0 || !strings.Contains(out, "System: None") || errOut != "" {
		t.Fatalf("unassociated: %d %q %q", code, out, errOut)
	}
	if err := os.Remove(filepath.Join(dir, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "specs", "jobs", "work.job.spec.md"), []byte("invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"job", "show", "work", "--path", dir}, {"job", "list", "--path", dir}} {
		code, out, errOut = invoke(args...)
		if code == 0 || out != "" || !strings.Contains(errOut, "work.job.spec.md") {
			t.Fatalf("malformed: %d %q %q", code, out, errOut)
		}
	}
}
