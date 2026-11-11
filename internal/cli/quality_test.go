package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

const qualitySystem = `---
type: system
---
# Quality project
## Purpose
Help contributors.
## Scope
CLI.
## Quality attributes
### Usability
**ID:** usability
Effective use.
### Diagnostics
**ID:** diagnostics
**Sub-attribute of:** usability
Understandable findings.
### Accessibility
**ID:** accessibility
**Sub-attribute of:** usability
Accessible output.
`

func TestQualityCLIResolutionAndFiltering(t *testing.T) {
	prepareExample(t, "dotspec")
	if code, out, errOut := invoke("quality", "list", "--json"); code != 0 || errOut != "" || !strings.Contains(out, `"groups": []`) {
		t.Fatalf("empty list: %d %s %s", code, out, errOut)
	}
	if code, out, errOut := invoke("quality", "init", "--path", ".", "--", "work"); code != 0 || errOut != "" || !strings.Contains(out, "Created specs/quality/work.quality.spec.md") {
		t.Fatalf("init: %d %s %s", code, out, errOut)
	}
	writeFixture(t, "specs/system.spec.md", qualitySystem)
	writeFixture(t, "specs/quality/work.quality.spec.md", "---\ntype: quality\nattributes: [diagnostics, accessibility]\n---\n# Accessible diagnostics\n")
	code, out, errOut := invoke("quality", "show", "work", "--json")
	var result qualityResult
	if err := json.Unmarshal([]byte(out), &result); err != nil || code != 0 || errOut != "" || result.Quality == nil {
		t.Fatalf("show: %d %s %s %v", code, out, errOut, err)
	}
	r := result.Quality
	if len(r.Attributes) != 2 || len(r.Classification.Resolved) != 2 || len(r.Classification.Ancestors) != 1 || r.Classification.Ancestors[0] != "usability" {
		t.Fatalf("classification: %+v", r)
	}
	code, out, errOut = invoke("quality", "show", "work")
	if code != 0 || errOut != "" || !strings.Contains(out, "Attributes: diagnostics, accessibility") || !strings.Contains(out, "Ancestor attributes: usability") {
		t.Fatalf("text: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("quality", "list", "--attribute", "usability", "--json")
	var list qualityListResult
	if err := json.Unmarshal([]byte(out), &list); err != nil || code != 0 || errOut != "" || len(list.Groups) != 1 || len(list.Groups[0].Requirements) != 1 {
		t.Fatalf("filter duplicates: %d %s %s %v", code, out, errOut, err)
	}
	code, out, errOut = invoke("quality", "list", "--attribute", "absent")
	if code != 0 || errOut != "" || !strings.Contains(out, "No quality requirements") {
		t.Fatalf("empty filter: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("system", "show", "--json")
	var sys systemResult
	if err := json.Unmarshal([]byte(out), &sys); err != nil || code != 0 || errOut != "" || sys.System.QualityAttributes == nil || len(sys.System.QualityAttributes) != 3 {
		t.Fatalf("model JSON: %d %s %s %v", code, out, errOut, err)
	}
	code, out, errOut = invoke("system", "show")
	if code != 0 || errOut != "" || !strings.Contains(out, "**Sub-attribute of:** usability") {
		t.Fatalf("model text: %d %s %s", code, out, errOut)
	}
	writeFixture(t, "specs/quality/work.quality.spec.md", "---\ntype: quality\nattributes: [unknown]\n---\n# Draft\n")
	code, out, errOut = invoke("quality", "check", "work", "--json")
	if code != 0 || errOut != "" || !strings.Contains(out, "quality.undefined-attribute") || !strings.Contains(out, `"context_files_checked": 1`) {
		t.Fatalf("advice: %d %s %s", code, out, errOut)
	}
}

func TestQualityCLIErrorChannels(t *testing.T) {
	prepareExample(t, "read")
	for _, args := range [][]string{
		{"init"}, {"init", "../escape"}, {"show", "UPPER"}, {"list", "extra"},
		{"list", "--attribute", ""}, {"list", "--attribute", "bad.id"},
		{"check"}, {"init", "unused", "--format", "mermaid"},
		{"show", "diagnostic-severity", "--json", "--format", "text"},
		{"show", "missing"}, {"init", "diagnostic-severity"},
	} {
		code, out, errOut := invoke(append([]string{"quality"}, args...)...)
		if code != 2 || out != "" || errOut == "" {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
	code, out, errOut := invoke("quality", "check", "missing", "--json")
	if code != 2 || errOut != "" || !strings.Contains(out, `"complete": false`) {
		t.Fatalf("missing check: %d %s %s", code, out, errOut)
	}
	writeFixture(t, "specs/quality/broken.quality.spec.md", "broken")
	for _, args := range [][]string{{"show", "broken"}, {"list"}, {"check", "broken"}} {
		code, out, errOut := invoke(append([]string{"quality"}, args...)...)
		if args[0] == "check" {
			out, errOut = errOut, out
		}
		if code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
}
