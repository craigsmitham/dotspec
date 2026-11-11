package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/job"
)

const viewFixture = `---
type: job
---
# Prepare a change
## Outcomes
- **OUT-001:** Minimize the likelihood of missed needs.
## Steps
### Gather needs
**ID:** gather
**Stage:** locate

Talk to stakeholders.
#### Outcomes
- **OUT-002:** Minimize the time needed to identify stakeholders.
### Reach agreement
**ID:** agree

## Progress cases
### Considering a new approach
An ordinary research note.
## Extra notes
Keep this text and its **formatting**.
`

func TestStructuredJobCommands(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "specs", "jobs"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	if err := os.Mkdir(".git", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("specs/jobs/change.job.spec.md", []byte(viewFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		args         []string
		want, absent string
	}{
		{[]string{"job", "show", "change"}, "Keep this text and its **formatting**.", "advisory"},
		{[]string{"job", "step", "list", "change"}, "1. Gather needs [gather] (locate)", "OUT-001"},
		{[]string{"job", "outcome", "list", "change"}, "OUT-001", "advisory"},
		{[]string{"job", "outcome", "list", "change", "--step", "gather"}, "OUT-002", "OUT-001"},
		{[]string{"job", "outcome", "list", "change", "--step", "agree"}, "No outcomes defined", "OUT-001"},
		{[]string{"job", "map", "change"}, "1. Gather needs [gather] (locate)\n   |\n   v\n2. Reach agreement [agree]", "OUT-001"},
		{[]string{"job", "map", "change", "--outcomes"}, "OUT-001", "advisory"},
		{[]string{"job", "map", "change", "--format", "mermaid"}, "s1 --> s2", "OUT-001"},
		{[]string{"job", "map", "change", "--format=mermaid", "--outcomes"}, "job -.-> jo1", "advisory"},
		{[]string{"job", "check", "change"}, "advisory findings: 2", "dotspec:"},
	}
	for _, tt := range tests {
		code, out, errOut := invoke(tt.args...)
		if code != 0 || errOut != "" || !strings.Contains(out, tt.want) || strings.Contains(out, tt.absent) {
			t.Errorf("%v: %d %q %q", tt.args, code, out, errOut)
		}
	}
	for _, args := range [][]string{
		{"job", "outcome", "list", "change", "--step", "missing"},
		{"job", "step", "show", "change"},
		{"job", "map", "change", "--format", "unknown"},
		{"job", "show", "change", "--outcomes"},
		{"job", "step", "list"},
	} {
		if code, out, errOut := invoke(args...); code == 0 || out != "" || errOut == "" {
			t.Errorf("expected error: %v: %d %q %q", args, code, out, errOut)
		}
	}
	if err := os.WriteFile("specs/jobs/change.job.spec.md", []byte(viewFixture+"\n[Broken](#missing)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := invoke("job", "check", "change"); code != 1 || errOut != "" || !strings.Contains(out, "unresolved local reference") {
		t.Fatalf("structural error: %d %q %q", code, out, errOut)
	}
}

func TestEmptyMapsAndNewCommandHelp(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "specs", "jobs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "specs", "jobs", "draft.job.spec.md"), []byte("---\ntype: job\n---\n# Draft\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"job", "map"}, {"job", "step", "list"}, {"job", "outcome", "list"}} {
		full := append(append([]string{}, args...), "draft", "--path", dir)
		if code, out, errOut := invoke(full...); code != 0 || !strings.Contains(out, "No ") || errOut != "" {
			t.Errorf("empty %v: %d %q %q", full, code, out, errOut)
		}
	}
	for _, args := range [][]string{{"job", "map", "--help"}, {"job", "check", "--help"}, {"job", "step", "--help"}, {"job", "step", "list", "--help"}, {"job", "outcome", "list", "--help"}} {
		if code, out, errOut := invoke(args...); code != 0 || !strings.Contains(out, "Usage:") || errOut != "" {
			t.Errorf("help %v: %d %q %q", args, code, out, errOut)
		}
	}
}

func TestOutcomeScopes(t *testing.T) {
	spec, err := job.Parse([]byte(viewFixture))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := listOutcomes(spec, "", &out); err != nil {
		t.Fatal(err)
	}
	want := "Outcomes: Prepare a change\n\nOutcomes declared by job\n  OUT-001: Minimize the likelihood of missed needs.\n\nOutcomes declared by step: Gather needs [gather] (locate)\n  OUT-002: Minimize the time needed to identify stakeholders.\n"
	if out.String() != want {
		t.Fatalf("outcome scopes:\n%s\nwant:\n%s", out.String(), want)
	}
	for _, format := range []string{"text", "mermaid"} {
		out.Reset()
		if err := writeJobMap(spec, jobOptions{format: format, outcomes: true}, &out); err != nil {
			t.Fatal(err)
		}
		text := out.String()
		for _, label := range []string{"Outcomes declared by job", "Outcomes declared by step: Gather needs", "Authored step order"} {
			if !strings.Contains(text, label) {
				t.Errorf("%s map missing %q: %s", format, label, text)
			}
		}
		for _, id := range []string{"OUT-001", "OUT-002"} {
			if strings.Count(text, id) != 1 {
				t.Errorf("%s map should display %s once in its own scope: %s", format, id, text)
			}
		}
		if format == "mermaid" {
			if !strings.Contains(text, "job -.-> jo1") || !strings.Contains(text, "s1 -.-> o1_1") || strings.Contains(text, "s2 -.->") {
				t.Errorf("incorrect outcome associations: %s", text)
			}
			if !strings.Contains(text, "not a required sequence; completeness not assessed") {
				t.Errorf("map must explain the scope of ordering: %s", text)
			}
		} else if !strings.Contains(text, "arrows do not require this execution sequence.\nCompleteness is not assessed.") {
			t.Errorf("map must explain the scope of ordering: %s", text)
		}
	}
}

func TestJobMapWithoutSteps(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "specs", "jobs"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, outcome := range []string{"", "\n## Outcomes\n\n- **O1:** Minimize the likelihood of missed needs.\n"} {
		if err := os.WriteFile(filepath.Join(dir, "specs", "jobs", "draft.job.spec.md"), []byte("---\ntype: job\n---\n# Draft\n"+outcome), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, format := range []string{"text", "mermaid"} {
			for _, includeOutcomes := range []bool{false, true} {
				args := []string{"job", "map", "draft", "--path", dir, "--format", format}
				if includeOutcomes {
					args = append(args, "--outcomes")
				}
				code, out, errOut := invoke(args...)
				if code != 0 || errOut != "" || !strings.Contains(out, "Draft") || !strings.Contains(out, "No steps defined.") {
					t.Fatalf("map without steps: %v: %d %q %q", args, code, out, errOut)
				}
				showOutcome := includeOutcomes && outcome != ""
				if strings.Contains(out, "O1:") != showOutcome || strings.Contains(out, "Outcomes declared by job") != showOutcome {
					t.Errorf("outcomes with no steps: %v: %s", args, out)
				}
				if format == "mermaid" && (!strings.HasPrefix(out, "flowchart TD\n") || strings.Contains(out, "s1")) {
					t.Errorf("expected Mermaid without invented step nodes: %s", out)
				}
			}
		}
	}
}

func TestStepSelectionExplainsMissingIdentity(t *testing.T) {
	spec := &job.Job{Name: "Prepare a change", Steps: []job.Step{{Name: "Locate commitments"}}}
	var out bytes.Buffer
	err := listOutcomes(spec, "locate", &out)
	if err == nil || !strings.Contains(err.Error(), "add **ID:** to that step") || out.Len() != 0 {
		t.Fatalf("selection should explain how to address an unidentified step: %v, %q", err, out.String())
	}
}

func TestMermaidEscapesAuthorContent(t *testing.T) {
	spec := &job.Job{Name: "A job", Steps: []job.Step{
		{Name: `Close "] --> injected["node <script>`},
		{Name: "end\nclick s1 \"https://example.com\"", ID: "last"},
	}}
	var out bytes.Buffer
	if err := writeJobMap(spec, jobOptions{format: "mermaid"}, &out); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if strings.Contains(s, "<script>") || strings.Contains(s, `"] --> injected`) || strings.Contains(s, "\nclick") {
		t.Fatalf("author content escaped its label: %s", s)
	}
	if !strings.Contains(s, "#34;") || !strings.Contains(s, "s1 --> s2") {
		t.Fatalf("missing escapes or sequence: %s", s)
	}
}
