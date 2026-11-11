package requirement

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/craigsmitham/dotspec/internal/quality"
)

func put(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestInventoryIncludesQualityAndLocalModel(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("initialize repository: %v: %s", err, output)
	}
	put(t, filepath.Join(dir, "specs/system.spec.md"), `---
type: system
---
# System
## Purpose
Support work.
## Scope
This project.
## Quality attributes
### Reliability
**ID:** reliability
Consistent operation.
### Preservation
**ID:** source-preservation
**Sub-attribute of:** reliability
Unchanged sources.
`)
	path, err := Init(dir, "work")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Init(dir, "work"); err == nil {
		t.Fatal("overwrote target")
	}
	put(t, path, example)
	put(t, quality.Layout.Path(dir, "work"), "---\ntype: quality\nattributes: [source-preservation]\n---\n# Preservation\n## Criterion\nUnchanged sources.\n## Measurement\nCompare snapshots.\n")
	groups, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	filtered := Filter(groups, "reliability")
	if len(filtered) != 1 || len(filtered[0].Requirements) != 2 {
		t.Fatalf("inventory: %+v", filtered)
	}
	q := filtered[0].Requirements[0]
	if q.Type != "quality" || q.Statement != "Unchanged sources." || q.Verification != "Compare snapshots." {
		t.Fatalf("quality view: %+v", q)
	}
	r, err := Show(dir, "work")
	if err != nil || r.Type != "requirement" || len(r.Classification.Resolved) != 2 {
		t.Fatalf("show: %+v %v", r, err)
	}
	put(t, Layout.Path(filepath.Join(dir, "child"), "work"), example)
	put(t, filepath.Join(dir, "child/specs/system.spec.md"), "---\ntype: system\n---\n# Child\n## Purpose\nWork.\n## Scope\nChild.\n")
	child, err := Show(filepath.Join(dir, "child"), "work")
	if err != nil || child.Classification.SystemPath != "" || len(child.Classification.Undefined) != 0 {
		t.Fatalf("model inherited: %+v %v", child, err)
	}
}
