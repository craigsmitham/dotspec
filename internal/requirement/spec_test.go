package requirement

import (
	"errors"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const example = `---
type: requirement
id: preserve
quality_attributes: [reliability, source-preservation]
---
# Preserve files
## Applies to
Inspection, including failed scans.
## Statement
The operation must preserve the input state.

### Postcondition
` + "```text\nstate_after = state_before\n```" + `

## Rationale
Inspection must not edit the work being inspected.
## Verification plan
Compare snapshots of source-file bytes and paths.
## Assumptions
No concurrent edits by other processes.
## References
[Rule](#preserve), [postcondition](#postcondition).
`

func TestGeneralRequirement(t *testing.T) {
	r, err := Parse([]byte(example))
	if err != nil {
		t.Fatal(err)
	}
	if r.Type != "requirement" || r.ID != "preserve" || len(r.QualityAttributes) != 2 || !strings.Contains(r.Statement, "state_after = state_before") || !strings.Contains(r.Markdown, "## Assumptions") || r.Rationale == "" || r.VerificationPlan == "" || len(Check(r)) != 0 {
		t.Fatalf("lost authored content: %+v", r)
	}
	if !r.References.Anchors["preserve"] || !r.References.Anchors["postcondition"] {
		t.Fatal("missing identity targets")
	}
}

func TestDraftAndStructure(t *testing.T) {
	draft := "---\ntype: requirement\n---\n# Draft\n"
	r, err := Parse([]byte(draft + "## Verification plan\n"))
	if err != nil {
		t.Fatal(err)
	}
	advice := Check(r)
	if len(advice) != 4 {
		t.Fatalf("draft advice: %+v", advice)
	}
	for _, d := range advice {
		if d.Severity != "advice" || (d.Rule == "requirement.verification" && d.Location.Line != 5) {
			t.Fatalf("advice location: %+v", d)
		}
	}
	for _, tt := range []struct{ text, rule string }{
		{strings.Replace(example, "type: requirement", "type: quality", 1), "metadata.type"},
		{strings.Replace(example, "id: preserve", "kind: functional", 1), "metadata.invalid"},
		{draft + "### Statement\n", "structure.section-level"},
		{draft + "## Statement\nOne.\n## Statement\nTwo.\n", "structure.duplicate-section"},
		{strings.Replace(example, "[reliability, source-preservation]", "[reliability, reliability]", 1), "identity.duplicate"},
		{strings.Replace(example, "[reliability, source-preservation]", "[false]", 1), "requirement.attributes"},
		{draft + "[missing](#not-declared)", "reference.unresolved"},
	} {
		_, err := Parse([]byte(tt.text))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line == 0 {
			t.Errorf("want %s: %v", tt.rule, err)
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, seed := range []string{example, Template, "---\ntype: requirement\nquality_attributes: [null]\n---\n# Draft\n"} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if r, err := Parse(data); err == nil {
			_ = Check(r)
		}
	})
}
