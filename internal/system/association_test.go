package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssociatedSystemHonorsReferenceBoundaries(t *testing.T) {
	root := t.TempDir()
	git(t, root, "init", "-q")
	writeSpec(t, root, "Default")
	target := writeSpec(t, filepath.Join(root, "selected"), "Selected")
	writeFile(t, filepath.Join(root, "nested", ".git", "marker"), "")
	writeSpec(t, filepath.Join(root, "nested"), "Nested")
	outside := writeSpec(t, t.TempDir(), "Outside")
	link := filepath.Join(root, "linked", "specs", Filename)
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "specs", "requirements", "selection.requirement.spec.md")
	outsideRel, err := filepath.Rel(filepath.Dir(source), outside)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		destination string
		valid       bool
	}{
		{"../../selected/specs/system.spec.md", true},
		{"../../selected/specs/system.spec.md#selected", true},
		{"../../selected/specs/system.spec.md#purpose", false},
		{"../../selected/specs/system.spec.md?alternate=true", false},
		{"../../nested/specs/system.spec.md", false},
		{"../../linked/specs/system.spec.md", false},
		{filepath.ToSlash(outsideRel), false},
		{"https://example.test/system.spec.md", false},
	} {
		t.Run(tc.destination, func(t *testing.T) {
			writeFile(t, source, fmt.Sprintf("---\ntype: requirement\n---\n# Selection\n## Specification system\n[Selected](%s)\n", tc.destination))
			got, err := Associated(root, source)
			if (err == nil) != tc.valid {
				t.Fatal(got, err)
			}
			if tc.valid && got.Path != target {
				t.Fatal(got)
			}
		})
	}
	writeFile(t, source, "---\ntype: requirement\n---\n# Selection\n## Specification system\n[Selected][system?]\n\n[system?]: ../../selected/specs/system.spec.md\n")
	got, err := Associated(root, source)
	if err != nil || got.Name != "Selected" {
		t.Fatal(got, err)
	}
	writeFile(t, source, strings.Replace(validSpec, "type: system", "type: requirement", 1))
	got, err = Associated(root, source)
	if err != nil || got.Name != "Default" {
		t.Fatal(got, err)
	}
}
