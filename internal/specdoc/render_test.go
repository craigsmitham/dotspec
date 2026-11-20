package specdoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
)

// Moving authored sections to another file changes both the base URL and the
// availability of reference definitions and footnotes.
func TestLiftedSectionPreservesAuthoredMarkdown(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"dotspec.yaml": "name: reader\ngenerated_spec: true\n",
		"module#1/specs/sample-system.system.spec.md":          "---\ntype: system-spec\n---\n# Sample system\n\n## Scope\n\nLocal.\n\n## Purpose\n\nKeep notes.\n\n## Concept uses\n\n### Display a note\n\nSee [![diagram](icon.svg)](guide%23a.md?mode=reader#setup) and [details][SUBJECT-SOURCE].[^note]\n\n## Evidence\n\n[SUBJECT-SOURCE]: guide%23a.md#setup\n\n[^note]: See [setup][SUBJECT-SOURCE].\n",
		"module#1/specs/guide#a.md":                            "# Setup\n",
		"module#1/specs/icon.svg":                              "<svg/>",
		"module#1/specs/requirements/save.requirement.spec.md": "---\ntype: requirement-spec\n---\n# Save a note\n\n## Rationale\n\nThe reader needs the note.\n\n## Statement\n\nPreserve the note.\n",
	}
	for file, body := range files {
		loc := filepath.Join(root, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(loc), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(loc, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	result, err := model.Export(root, system.CheckType(), requirement.CheckType())
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("source check: %v", result.Snapshot.Diagnostics)
	}
	plan, err := Generate(root, result.Snapshot)
	if err != nil {
		t.Fatal(err)
	}
	file := "specification/concept-uses/display-a-note.concept-use.spec-report.md"
	out := plan.Pages[file].Markdown
	for _, want := range []string{
		"[![diagram](<../../module%231/specs/icon.svg>)](<../../module%231/specs/guide%23a.md?mode=reader#setup>)",
		"[details](<../../module%231/specs/guide%23a.md#setup>)",
		"[^note]: See [setup](<../../module%231/specs/guide%23a.md#setup>).",
		`source: "../../module%231/specs/sample-system.system.spec.md#display-a-note"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("%s missing %q\n%s", file, want, out)
		}
	}
	detail := plan.Pages["specification/requirements/save-a-note.requirement.spec-report.md"].Markdown
	if strings.Index(detail, "## Rationale") < strings.Index(detail, "## Statement") {
		t.Fatal(detail)
	}
	if _, err := Write(root, plan); err != nil {
		t.Fatal(err)
	}
}
