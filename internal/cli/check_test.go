package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dotspec-md/dotspec/internal/check"
)

func TestProjectCheckTextAndJSON(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "specs/jobs/duplicate.job.spec.md"), "---\ntype: job-spec\nid: repeated\n---\n# Example\n## Steps\n### Gather\n**ID:** repeated\n")
	writeFixture(t, filepath.Join(dir, "specs/jobs/draft.job.spec.md"), "---\ntype: job-spec\n---\n# Draft\n")
	code, out, errOut := invoke("check", dir, "--json")
	if code != 1 || errOut != "" {
		t.Fatalf("JSON check: %d %q %q", code, out, errOut)
	}
	var report check.Report
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != 1 || !report.Complete || report.Summary.Errors != 1 || report.Summary.Advice != 2 || report.Summary.FilesChecked != 2 {
		t.Fatalf("JSON report: %+v", report)
	}
	code, out, errOut = invoke("check", dir)
	if code != 1 || errOut != "" || !strings.Contains(out, "(scan complete): 1 error, 2 advisories; 0 context files.") || !strings.HasPrefix(out, diagnosticPath(report.Diagnostics[0].Location)) {
		t.Fatalf("text check: %d %q %q", code, out, errOut)
	}
	for _, d := range report.Diagnostics {
		if !strings.Contains(out, "["+d.Rule+"]") || !strings.Contains(out, d.Message) || !strings.Contains(out, fmt.Sprintf(":%d:%d", d.Location.Line, d.Location.Column)) {
			t.Errorf("text and JSON disagree on %+v: %s", d, out)
		}
		for _, related := range d.Related {
			if !strings.Contains(out, "Related: "+diagnosticPath(related)) {
				t.Errorf("missing related source: %s", out)
			}
		}
	}
	// The focused command still reports the same rule and both definitions.
	code, out, errOut = invoke("-C", dir, "check", "specs/jobs/duplicate.job.spec.md")
	if code != 1 || errOut != "" || !strings.Contains(out, "[identity.duplicate]") || !strings.Contains(out, "Related:") || !strings.Contains(out, ":3:5") {
		t.Fatalf("focused check lost shared diagnostics: %d %q %q", code, out, errOut)
	}
}

func TestCheckTextLeadsWithFindingsAndListsExistingContext(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "dotspec.yaml"), "name: desk\nbuild:\n  enabled: false\n")
	writeFixture(t, filepath.Join(dir, "notes.md"), "# Notes\n")
	writeFixture(t, filepath.Join(dir, "specs/jobs/work.job.spec.md"), "---\ntype: job-spec\n---\n# Work\n\nSee [notes](../../notes.md) and [gone](../../gone.md).\n")
	t.Chdir(dir)
	code, out, errOut := invoke("check")
	if code != 1 || errOut != "" {
		t.Fatalf("check: %d %q %q", code, out, errOut)
	}
	summary := "Checked 1 file in project desk (scan complete): 1 error, 2 advisories; 2 context files.\n"
	context := "\nContext files:\n  dotspec.yaml\n  notes.md\n"
	if !strings.HasPrefix(out, "specs/jobs/work.job.spec.md:") || !strings.Contains(out, summary+context) || !strings.HasSuffix(out, context) {
		t.Fatalf("findings, summary, then context:\n%s", out)
	}
	// The unresolved target is a finding, not context; JSON still accounts for it.
	if strings.Contains(out, "  gone.md") || !strings.Contains(out, "Related: gone.md") {
		t.Fatalf("unresolved target listed as context:\n%s", out)
	}
	code, out, _ = invoke("check", "--json")
	var report check.Report
	if code != 1 || json.Unmarshal([]byte(out), &report) != nil || report.Summary.ContextFilesChecked != 3 || len(report.ContextFiles) != 3 {
		t.Fatalf("JSON context accounting: %s", out)
	}
	code, out, _ = invoke("check", "specs/jobs/work.job.spec.md")
	if code != 1 || !strings.Contains(out, "Checked 1 file in specs/jobs (scan complete)") || strings.Contains(out, "Generated specification:") {
		t.Fatalf("focused check of a project without a generated specification:\n%s", out)
	}
}

func TestCheckTextStatesGeneratedSpecificationScope(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "dotspec.yaml"), "name: desk\n")
	writeFixture(t, filepath.Join(dir, "specs/jobs/work.job.spec.md"), "---\ntype: job-spec\n---\n# Work\n")
	t.Chdir(dir)
	for args, want := range map[string]string{
		"specs/jobs/work.job.spec.md": "Generated specification: not checked (focused check)\n",
		"specs":                       "Generated specification: not checked (subtree check)\n",
	} {
		code, out, errOut := invoke("check", args)
		if code != 0 || errOut != "" || !strings.Contains(out, want) || strings.Contains(out, "out-of-scope") {
			t.Errorf("check %s: %d %q %q", args, code, out, errOut)
		}
	}
}

func TestProjectCheckEmptyAdvisoryAndIncomplete(t *testing.T) {
	dir := t.TempDir()
	code, out, errOut := invoke("check", dir)
	if code != 0 || errOut != "" || !strings.Contains(out, "No supported source files found.") {
		t.Fatalf("empty: %d %q %q", code, out, errOut)
	}
	writeFixture(t, filepath.Join(dir, "specs/jobs/draft.job.spec.md"), "---\ntype: job-spec\n---\n# Draft\n")
	code, out, errOut = invoke("check", "--json", "--chdir="+dir)
	if code != 0 || errOut != "" || !strings.Contains(out, `"severity": "advice"`) {
		t.Fatalf("advice should not fail: %d %q %q", code, out, errOut)
	}
	for _, asJSON := range []bool{false, true} {
		args := []string{"check", filepath.Join(dir, "missing")}
		if asJSON {
			args = append(args, "--json")
		}
		code, out, errOut = invoke(args...)
		if code != 2 || errOut != "" {
			t.Fatalf("incomplete %v: %d %q %q", args, code, out, errOut)
		}
		if asJSON {
			var r check.Report
			if err := json.Unmarshal([]byte(out), &r); err != nil || r.Complete || r.Diagnostics[0].Rule != "scan.scope" {
				t.Fatalf("incomplete JSON: %v, %+v", err, r)
			}
		} else if !strings.Contains(out, "(scan incomplete)") || strings.Contains(out, "No supported source files found.") {
			t.Fatalf("incomplete scan looks successful: %s", out)
		}
	}
}

func TestProjectCheckUsage(t *testing.T) {
	for _, args := range [][]string{
		{"check", "--bad"}, {"check", "one", "two"}, {"check", "-C"},
		{"check", "-C", ""}, {"check", "--json=yaml"},
	} {
		if code, out, errOut := invoke(args...); code != 2 || out != "" || errOut == "" {
			t.Errorf("invalid invocation: %v: %d %q %q", args, code, out, errOut)
		}
	}
	t.Chdir(t.TempDir())
	writeFixture(t, "notes.txt", "Not a specification.\n")
	if code, out, errOut := invoke("check", "notes.txt"); code != 2 || out != "" || !strings.Contains(errOut, "unknown source document type: notes.txt") {
		t.Errorf("non-specification file: %d %q %q", code, out, errOut)
	}
	if code, out, errOut := invoke("check", "--help"); code != 0 || errOut != "" || !strings.Contains(out, "Exits 1 when error findings are reported.") || !strings.Contains(out, "--json") {
		t.Fatalf("help: %d %q %q", code, out, errOut)
	}
}

func TestSourcesTextNamesEachProjectOnce(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "dotspec.yaml"), "name: desk\nbuild:\n  enabled: false\n")
	writeFixture(t, filepath.Join(dir, "notes.md"), "# Notes\n")
	writeFixture(t, filepath.Join(dir, "specs/jobs/work.job.spec.md"), "---\ntype: job-spec\n---\n# Work\n\nSee [notes](../../notes.md).\n")
	t.Chdir(filepath.Join(dir, "specs"))
	code, out, errOut := invoke("sources")
	if code != 0 || errOut != "" || strings.Count(out, "Project: desk") != 1 || !strings.HasPrefix(out, "Project: desk\n\n../dotspec.yaml\n") {
		t.Fatalf("sources: %d %q %q", code, out, errOut)
	}
	// Each listed specification path can be pasted into check.
	if !strings.Contains(out, "\njobs/work.job.spec.md\n") || !strings.HasSuffix(out, "Listed 3 source files in project desk (scan complete).\n") {
		t.Fatalf("sources paths: %s", out)
	}
	if code, out, errOut := invoke("check", "jobs/work.job.spec.md"); code != 0 || errOut != "" {
		t.Fatalf("listed path does not check: %d %q %q", code, out, errOut)
	}
}
