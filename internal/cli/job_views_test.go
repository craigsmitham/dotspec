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
	t.Chdir(dir)
	if err := os.Mkdir(".git", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("change.job.spec.md", []byte(viewFixture), 0o644); err != nil {
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
		{[]string{"job", "check", "change"}, "writing suggestions (advisory)", "dotspec:"},
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
	if err := os.WriteFile("change.job.spec.md", []byte(viewFixture+"\n[Broken](#missing)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := invoke("job", "check", "change"); code == 0 || out != "" || !strings.Contains(errOut, "unresolved local reference") {
		t.Fatalf("structural error: %d %q %q", code, out, errOut)
	}
}

func TestEmptyMapsAndNewCommandHelp(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "draft.job.spec.md"), []byte("---\ntype: job\n---\n# Draft\n"), 0o644); err != nil {
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
