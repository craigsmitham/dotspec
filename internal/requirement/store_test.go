package requirement

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

func TestInventoryClassifiesAgainstLocalModel(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("initialize repository: %v: %s", err, output)
	}
	put(t, filepath.Join(dir, "specs/system.system.spec.md"), "---\ntype: system-spec\n---\n# System\n## Purpose\nSupport work.\n## Scope\nThis project.\n## Quality attributes\n")
	put(t, filepath.Join(dir, "specs/quality/reliability/reliability.quality-attribute.spec.md"), "---\ntype: quality-attribute-spec\nid: reliability\n---\n# Reliability\n\n## Definition\n\nConsistent operation.\n")
	put(t, filepath.Join(dir, "specs/quality/reliability/preservation/preservation.quality-attribute.spec.md"), "---\ntype: quality-attribute-spec\nid: source-preservation\nsubattribute_of: [reliability]\n---\n# Preservation\n\n## Definition\n\nUnchanged sources.\n")

	path, err := Init(dir, "work")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Init(dir, "work"); err == nil {
		t.Fatal("overwrote target")
	}
	put(t, path, example)
	put(t, Layout.Path(dir, "preservation"), "---\ntype: requirement-spec\nquality_attributes: [source-preservation]\n---\n# Preservation\n## Statement\nUnchanged sources.\n## Evidence\nCompare snapshots.\n")
	groups, err := List(dir)
	if err != nil {
		t.Fatal(err)
	}
	filtered := Filter(groups, "reliability")
	if len(filtered) != 1 || len(filtered[0].Requirements) != 2 {
		t.Fatalf("inventory: %+v", filtered)
	}
	q := filtered[0].Requirements[0]
	if q.Statement != "Unchanged sources." {
		t.Fatalf("classified requirement: %+v", q)
	}
	r, err := Show(dir, "work")
	if err != nil || len(r.Classification.Resolved) != 2 {
		t.Fatalf("show: %+v %v", r, err)
	}
	put(t, Layout.Path(filepath.Join(dir, "child"), "work"), example)
	put(t, filepath.Join(dir, "child/specs/child.system.spec.md"), "---\ntype: system-spec\n---\n# Child\n## Purpose\nWork.\n## Scope\nChild.\n")

	child, err := Show(filepath.Join(dir, "child"), "work")
	if err != nil || child.Classification.SystemPath != "" || len(child.Classification.Undefined) != 0 {
		t.Fatalf("model inherited: %+v %v", child, err)
	}
}

// The quality preset is an authoring aid. It creates an ordinary requirement
// and classification, not the template, selects quality-specific advice.
func TestQualityPresetCreatesAnOrdinaryRequirement(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("initialize repository: %v: %s", err, output)
	}
	path, err := InitWithPreset(dir, "bounded", "quality")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !Layout.Match(path) || !strings.Contains(string(data), "type: requirement-spec") || !strings.Contains(string(data), "## Criterion") {
		t.Fatalf("quality preset: %s\n%s", path, data)
	}
	if _, err := InitWithPreset(dir, "other", "performance"); err == nil {
		t.Fatal("accepted an unknown preset")
	}
	draft := "---\ntype: requirement-spec\n---\n# Draft\n"
	classified := strings.Replace(draft, "type: requirement-spec", "type: requirement-spec\nquality_attributes: [reliability]", 1)
	messages := func(source string) string {
		spec, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		text := ""
		for _, d := range Check(spec) {
			text += d.Rule + ": " + d.Message + "\n"
		}
		return text
	}
	if general, quality := messages(draft), messages(classified); !strings.Contains(general, "required result, rule, property, or restriction") || !strings.Contains(quality, "assessable quality obligation") || !strings.Contains(quality, "requirement.statement") {
		t.Fatalf("advice wording should follow classification:\n%s\n%s", general, quality)
	}
}
