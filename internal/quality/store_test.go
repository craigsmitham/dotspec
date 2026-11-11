package quality

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const model = `---
type: system
---
# Root
## Purpose
Help contributors.
## Scope
This repository.
## Quality attributes
### Reliability
**ID:** reliability
Consistent behavior.
### Repeatability
**ID:** repeatability
**Sub-attribute of:** reliability
Stable results.
### Diagnosability
**ID:** diagnosability
**Sub-attribute of:** reliability
Understandable findings.
`

func TestResolutionGroupingAndNoAttributeInheritance(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	write(t, filepath.Join(dir, "specs/system.spec.md"), model)
	write(t, Layout.Path(dir, "work"), example)
	child := filepath.Join(dir, "child")
	write(t, Layout.Path(child, "work"), example)
	spec, err := Show(child, "work")
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Classification.Resolved) != 2 || !reflect.DeepEqual(spec.Classification.Ancestors, []string{"reliability"}) || spec.Classification.SystemPath != filepath.Join(dir, "specs/system.spec.md") {
		t.Fatalf("resolution: %+v", spec)
	}
	groups, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	filtered := Filter(groups, "reliability")
	if len(filtered) != 1 || len(filtered[0].Requirements) != 2 {
		t.Fatalf("multiple paths double counted requirements: %+v", filtered)
	}
	if len(groups[0].Requirements) != 2 || len(Filter(groups, "unknown")) != 0 {
		t.Fatal("filter modified original groups or matched unknown attribute")
	}
	// The child system deliberately has no attribute declarations. It must not acquire the
	// parent's attributes simply because the attribute labels happen to match.
	write(t, filepath.Join(child, "specs/system.spec.md"), strings.Split(model, "## Quality attributes")[0])
	spec, err = Show(child, "work")
	if err != nil || spec.Classification.SystemPath != "" || len(spec.Classification.Resolved) != 0 || len(spec.Classification.Undefined) != 0 {
		t.Fatalf("inherited attributes: %+v %v", spec, err)
	}
	groups, err = List(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(Filter(groups, "reliability")) != 1 || len(Filter(groups, "repeatability")) != 2 {
		t.Fatalf("attribute grouping crossed system boundary: %+v", groups)
	}
	// Declaring an empty local attribute list changes unbound labels to undefined ones.
	write(t, filepath.Join(child, "specs/system.spec.md"), strings.Split(model, "## Quality attributes")[0]+"## Quality attributes\n")
	spec, err = Show(child, "work")
	if err != nil || len(spec.Classification.Undefined) != 2 {
		t.Fatalf("undefined labels: %+v %v", spec, err)
	}
	advice := CheckAttributes(spec, spec.System)
	if len(advice) != 2 || advice[0].Location.Path != spec.Path || advice[0].Location.Line != 4 {
		t.Fatalf("advice: %+v", advice)
	}
	write(t, filepath.Join(child, "specs/system.spec.md"), "broken")
	if _, err := Show(child, "work"); err == nil {
		t.Fatal("ignored malformed local model")
	}
}

func TestCreateAndLocalIdentity(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, slug := range []string{"", "../escape", "UPPER", "two--parts", "a/b"} {
		if _, err := Init(dir, slug); err == nil {
			t.Errorf("accepted %q", slug)
		}
	}
	if files, _ := os.ReadDir(dir); len(files) != 0 {
		t.Fatal("invalid invocation wrote files")
	}
	path, err := Init(dir, "work")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Init(dir, "work"); err == nil {
		t.Fatal("overwrote existing file")
	}
	if data, _ := os.ReadFile(path); string(data) != Template {
		t.Fatal("template changed")
	}
	write(t, path, example)
	spec, err := Show(dir, "work")
	if err != nil || spec.ID != "repeatable-inspection" || spec.Slug != "work" {
		t.Fatalf("identity: %+v %v", spec, err)
	}
	child := filepath.Join(dir, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Show(child, "work"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("slug lookup searched parent: %v", err)
	}
	write(t, path, "broken")
	if _, err := List(dir); err == nil {
		t.Fatal("list omitted malformed requirement")
	}
}
