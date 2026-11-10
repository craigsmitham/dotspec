package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
)

func TestProjectCheckTextAndJSON(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "specs/jobs/duplicate.job.spec.md"), "---\ntype: job\nid: repeated\n---\n# Example\n## Steps\n### Gather\n**ID:** repeated\n")
	writeFixture(t, filepath.Join(dir, "specs/jobs/draft.job.spec.md"), "---\ntype: job\n---\n# Draft\n")
	code, out, errOut := invoke("check", "--path", dir, "--format", "json")
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
	code, out, errOut = invoke("check", "--path", dir)
	if code != 1 || errOut != "" || !strings.Contains(out, "Scan complete.") || !strings.Contains(out, "Errors: 1; advisory findings: 2.") {
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
	code, out, errOut = invoke("job", "check", "duplicate", "--path", dir)
	if code != 1 || out != "" || !strings.Contains(errOut, "[identity.duplicate]") || !strings.Contains(errOut, "Related:") || !strings.Contains(errOut, ":3:5") {
		t.Fatalf("focused check lost shared diagnostics: %d %q %q", code, out, errOut)
	}
}

func TestProjectCheckEmptyAdvisoryAndIncomplete(t *testing.T) {
	dir := t.TempDir()
	code, out, errOut := invoke("check", "--path", dir)
	if code != 0 || errOut != "" || !strings.Contains(out, "No supported specifications found.") {
		t.Fatalf("empty: %d %q %q", code, out, errOut)
	}
	writeFixture(t, filepath.Join(dir, "specs/jobs/draft.job.spec.md"), "---\ntype: job\n---\n# Draft\n")
	code, out, errOut = invoke("check", "--format=json", "--path="+dir)
	if code != 0 || errOut != "" || !strings.Contains(out, `"severity": "advice"`) {
		t.Fatalf("advice should not fail: %d %q %q", code, out, errOut)
	}
	for _, format := range []string{"text", "json"} {
		code, out, errOut = invoke("check", "--path", filepath.Join(dir, "missing"), "--format", format)
		if code != 2 || errOut != "" {
			t.Fatalf("incomplete %s: %d %q %q", format, code, out, errOut)
		}
		if format == "json" {
			var r check.Report
			if err := json.Unmarshal([]byte(out), &r); err != nil || r.Complete || r.Diagnostics[0].Rule != "scan.scope" {
				t.Fatalf("incomplete JSON: %v, %+v", err, r)
			}
		} else if !strings.Contains(out, "Scan incomplete.") || strings.Contains(out, "No supported specifications found.") {
			t.Fatalf("incomplete scan looks successful: %s", out)
		}
	}
}

func TestProjectCheckUsage(t *testing.T) {
	for _, args := range [][]string{
		{"check", "--bad"}, {"check", "extra"}, {"check", "--path"},
		{"check", "--path", ""}, {"check", "--format", "yaml"},
	} {
		if code, out, errOut := invoke(args...); code != 2 || out != "" || errOut == "" {
			t.Errorf("invalid invocation: %v: %d %q %q", args, code, out, errOut)
		}
	}
	if code, out, errOut := invoke("check", "--help"); code != 0 || errOut != "" || !strings.Contains(out, "Exit codes:") || !strings.Contains(out, "--format") {
		t.Fatalf("help: %d %q %q", code, out, errOut)
	}
}

type failedOutput struct{}

func (failedOutput) Write([]byte) (int, error) { return 0, fmt.Errorf("output unavailable") }

func TestProjectCheckOutputFailure(t *testing.T) {
	for _, format := range []string{"text", "json"} {
		var stderr bytes.Buffer
		code := Run([]string{"check", "--path", t.TempDir(), "--format", format}, failedOutput{}, &stderr)
		if code != 2 || !strings.Contains(stderr.String(), "output unavailable") {
			t.Fatalf("output failure: %d %s", code, &stderr)
		}
	}
}
