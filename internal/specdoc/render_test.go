package specdoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
)

// Moving authored sections to another file changes both the base URL and the
// availability of reference definitions and footnotes.
func TestLiftedSectionPreservesAuthoredMarkdown(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"dotspec.yaml": "name: reader\n",
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

// Only the destination and title change; the authored label, including a code
// span that contains a bracket, is kept.
func TestRewriteDestinationKeepsLabels(t *testing.T) {
	for _, c := range []struct {
		raw  string
		n    model.Record
		want string
	}{
		{"[`]`](x.md)", model.Record{"node_kind": "link"}, "[`]`](<to.md>)"},
		{"![a [b]](x.png)", model.Record{"node_kind": "image"}, "![a [b]](<to.md>)"},
		{"[a \\]][ref]", model.Record{"node_kind": "link", "title": `say \"hi\"`}, `[a \]](<to.md> "say \"hi\"")`},
		{"[ref]: x.md 'it\\'s \"so\"'\n", model.Record{"node_kind": "link_definition", "title": `it\'s "so"`}, "[ref]: <to.md> \"it\\'s \\\"so\\\"\"\n"},
	} {
		if got := rewriteDestination(c.raw, c.n, "to.md"); got != c.want {
			t.Errorf("rewriteDestination(%q) = %q, want %q", c.raw, got, c.want)
		}
	}
}

// A wording change may give two relationship directions the same caption.
// Their targets still belong in the directory or the detail according to role.
func TestRelationshipDirectorySurvivesCaptionChanges(t *testing.T) {
	system := &subject{r: model.Record{"kind": "System"}, file: "system.md", source: &sourceFile{}}
	system.name = &heading{node: "system", owner: system}
	system.index("Components").add(entry{title: "Part", target: "part.md"})
	system.index("Requirements").add(entry{title: "Rule", target: "rule.md"})

	part, _ := reportRelation("Component.part_of")
	supports, _ := reportRelation("Requirement.supports")
	part.reverse.label, supports.reverse.label = "Related", "Related"
	p := planner{routes: map[string][]route{"system": {
		{relationCaption: part.reverse, title: "Part", target: "part.md"},
		{relationCaption: supports.reverse, title: "Rule", target: "rule.md"},
	}}}
	w := writer{page: &Page{}, file: system.file}
	p.relationships(&w, system.name, true)
	if got := w.out.String(); strings.Contains(got, "part.md") || !strings.Contains(got, "[Rule](<rule.md>)") {
		t.Fatal(got)
	}
}

// A job's page is the view of its map, so a job without steps says so.
func TestJobPageReportsAbsentSteps(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"dotspec.yaml":                  "name: jobs\n",
		"specs/sample.system.spec.md":   "---\ntype: system-spec\n---\n# Sample\n\n## Purpose\n\nKeep notes.\n",
		"specs/jobs/draft.job.spec.md":  "---\ntype: job-spec\n---\n# Draft\n\n## Job\n\nSettle an errand.\n\n## Related specifications\n\n[Mapped](mapped.job.spec.md)\n",
		"specs/jobs/empty.job.spec.md":  "---\ntype: job-spec\n---\n# Empty\n\n## Job\n\nSettle an errand.\n\n## Steps\n\n## Evidence\n\nAn interview.\n",
		"specs/jobs/mapped.job.spec.md": "---\ntype: job-spec\n---\n# Mapped\n\n## Job\n\nSettle an errand.\n\n## Steps\n\n### Recall outstanding work\n\nRecover the commitments.\n",
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
	result, err := model.Export(root, system.CheckType(), job.CheckType())
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
	draft := plan.Pages["specification/jobs/draft.job.spec-report.md"].Markdown
	steps, related := strings.Index(draft, "## Steps\n\nNo steps recorded."), strings.Index(draft, "## Related specifications")
	if steps < 0 || related < steps {
		t.Fatalf("absent steps should be reported before cross-references:\n%s", draft)
	}
	if empty := plan.Pages["specification/jobs/empty.job.spec-report.md"].Markdown; strings.Count(empty, "## Steps") != 1 || !strings.Contains(empty, "No steps recorded.") {
		t.Fatalf("an empty Steps heading records no steps:\n%s", empty)
	}
	if mapped := plan.Pages["specification/jobs/mapped.job.spec-report.md"].Markdown; strings.Contains(mapped, "No steps recorded") || !strings.Contains(mapped, "### Recall outstanding work") {
		t.Fatalf("authored steps:\n%s", mapped)
	}
}
