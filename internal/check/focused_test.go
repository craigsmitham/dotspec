package check_test

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestFocusedCheckSharesReportContract(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "specs/system.spec.md"), "---\ntype: system\n---\n# Product\n## Purpose\nHelp.\n## Scope\nThe project.\n")
	write(t, filepath.Join(dir, "specs/jobs/work.job.spec.md"), "---\ntype: job\n---\n# Work\n## Job\nImprove something.\n")
	focused, project := checkJob(dir, "work"), checkProject(dir)
	if !reflect.DeepEqual(focused.Diagnostics, project.Diagnostics) {
		t.Fatalf("focused/project findings differ:\n%+v\n%+v", focused.Diagnostics, project.Diagnostics)
	}
	if focused.ExitCode() != 0 || !focused.Complete || focused.Summary.FilesChecked != 1 || focused.Summary.ContextFilesChecked != 1 {
		t.Fatalf("wrong accounting: %+v", focused)
	}
	if len(focused.Diagnostics) != 2 || focused.Diagnostics[0].Rule != "job.performer" {
		t.Fatalf("findings must be ordered by location: %+v", focused.Diagnostics)
	}
	missing := checkJob(dir, "absent")
	if missing.ExitCode() != 2 || missing.Complete || missing.Diagnostics[0].Rule != "scan.read" {
		t.Fatalf("missing input: %+v", missing)
	}
}
