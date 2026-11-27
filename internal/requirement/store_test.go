package requirement

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/system"
)

// classify resolves a requirement's labels against its associated system.
func classify(t *testing.T, path string) Classification {
	t.Helper()
	spec, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := system.Associated(Layout.ContextDir(path), path)
	if err != nil {
		t.Fatal(err)
	}
	return Classify(spec.QualityAttributes, owner)
}

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
	put(t, filepath.Join(dir, "specs/system.system.spec.md"), "---\ntype: system-spec\n---\n# System\n## Purpose\nSupport work.\n## Scope\nThis project.\n")
	put(t, filepath.Join(dir, "specs/quality/reliability/reliability.quality-attribute.spec.md"), "---\ntype: quality-attribute-spec\nid: reliability\n---\n# Reliability\n\n## Definition\n\nConsistent operation.\n")
	put(t, filepath.Join(dir, "specs/quality/reliability/preservation/preservation.quality-attribute.spec.md"), "---\ntype: quality-attribute-spec\nid: source-preservation\nsubattribute_of: [reliability]\n---\n# Preservation\n\n## Definition\n\nUnchanged sources.\n")

	path, err := InitWithAttribute(dir, "work", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := InitWithAttribute(dir, "work", ""); err == nil {
		t.Fatal("overwrote target")
	}
	put(t, path, example)
	put(t, Layout.Path(dir, "preservation"), "---\ntype: requirement-spec\nquality_attributes: [source-preservation]\n---\n# Preservation\n## Statement\nUnchanged sources.\n## Evidence\nCompare snapshots.\n")
	listed, err := speccollection.Members(dir, Layout, "requirements", func(path string) (string, error) {
		spec, err := Read(path)
		if err != nil {
			return "", err
		}
		return spec.Name, nil
	})
	if err != nil || len(listed) != 2 || listed[0].System == nil {
		t.Fatalf("inventory: %+v %v", listed, err)
	}
	if q := classify(t, Layout.Path(dir, "preservation")); !slices.Contains(q.Ancestors, "reliability") {
		t.Fatalf("classified requirement: %+v", q)
	}
	if r := classify(t, path); len(r.Resolved) != 2 {
		t.Fatalf("classification: %+v", r)
	}
	put(t, Layout.Path(filepath.Join(dir, "child"), "work"), example)
	put(t, filepath.Join(dir, "child/specs/child.system.spec.md"), "---\ntype: system-spec\n---\n# Child\n## Purpose\nWork.\n## Scope\nChild.\n")

	if child := classify(t, Layout.Path(filepath.Join(dir, "child"), "work")); child.SystemPath != "" || len(child.Undefined) != 0 {
		t.Fatalf("model inherited: %+v", child)
	}
}

// Classification, not an authoring template, selects quality-specific advice.
func TestClassificationSelectsQualityAdvice(t *testing.T) {
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
	if general, quality := messages(draft), messages(classified); !strings.Contains(general, "required result, rule, property, or restriction") || !strings.Contains(quality, "assessable quality requirement") || !strings.Contains(quality, "requirement.statement") {
		t.Fatalf("advice wording should follow classification:\n%s\n%s", general, quality)
	}
}
