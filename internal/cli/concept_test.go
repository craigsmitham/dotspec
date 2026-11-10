package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
)

func TestConceptCommandsAndJSON(t *testing.T) {
	prepareExample(t, "dotspec")
	code, out, errOut := invoke("concept", "list", "--json")
	if code != 0 || errOut != "" || !strings.Contains(out, `"groups": []`) {
		t.Fatalf("empty JSON: %d %s %s", code, out, errOut)
	}
	if code, out, errOut := invoke("concept", "list"); code != 0 || errOut != "" || !strings.Contains(out, "No concepts found") {
		t.Fatalf("empty text: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("concept", "init", "--path", ".", "review", "--json")
	var created createdResult
	if err := json.Unmarshal([]byte(out), &created); err != nil || code != 0 || errOut != "" || !filepath.IsAbs(created.Path) {
		t.Fatalf("created JSON: %d %s %s %v", code, out, errOut, err)
	}
	if code, _, _ := invoke("system", "init"); code != 0 {
		t.Fatal("system init failed")
	}
	code, out, errOut = invoke("concept", "show", "--format", "json", "--", "review")
	var result struct {
		SchemaVersion int `json:"schema_version"`
		Concept       *struct {
			Slug                 string `json:"slug"`
			Markdown             string `json:"markdown"`
			Purpose              string `json:"purpose"`
			State                string `json:"state"`
			Actions              string `json:"actions"`
			OperationalPrinciple string `json:"operational_principle"`
			System               *struct {
				Name string `json:"name"`
			} `json:"system"`
		} `json:"concept"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil || code != 0 || errOut != "" || result.Concept == nil {
		t.Fatalf("show JSON: %d %s %s %v", code, out, errOut, err)
	}
	if result.SchemaVersion != 1 || result.Concept.Slug != "review" || result.Concept.System == nil || result.Concept.System.Name != "New system" || result.Concept.Purpose == "" || result.Concept.State == "" || result.Concept.Actions == "" || result.Concept.OperationalPrinciple == "" {
		t.Fatalf("lost fields or association: %+v", result)
	}
	code, out, errOut = invoke("concept", "show", "review")
	if code != 0 || errOut != "" || !strings.Contains(out, result.Concept.Markdown) || !strings.Contains(out, "System: New system") {
		t.Fatalf("show text: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("concept", "list", "--json")
	var groups struct {
		Groups []struct {
			Concepts []any `json:"concepts"`
		} `json:"groups"`
	}
	if err := json.Unmarshal([]byte(out), &groups); err != nil || code != 0 || errOut != "" || len(groups.Groups) != 1 || len(groups.Groups[0].Concepts) != 1 {
		t.Fatalf("list JSON: %d %s %s %v", code, out, errOut, err)
	}
	code, out, errOut = invoke("concept", "check", "review", "--json")
	var report check.Report
	if err := json.Unmarshal([]byte(out), &report); err != nil || code != 0 || errOut != "" || !report.Complete || report.Summary.FilesChecked != 1 || report.Scope.Types[0] != "concept" {
		t.Fatalf("focused report: %d %s %s %v", code, out, errOut, err)
	}
	before, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatal(err)
	}
	if code, _, _ := invoke("concept", "init", "review"); code != 2 {
		t.Fatal("init overwrote an existing specification")
	}
	if after, err := os.ReadFile(created.Path); err != nil || string(before) != string(after) {
		t.Fatal("commands changed existing content")
	}
}

func TestConceptUsageAndFailureChannels(t *testing.T) {
	prepareExample(t, "concepts")
	for _, args := range [][]string{
		{"init"}, {"init", "../escape"}, {"show", "UPPER"}, {"show", "work", "extra"},
		{"list", "extra"}, {"check"}, {"check", "../escape"}, {"list", "--format", "mermaid"},
		{"init", "unused", "--json", "--format", "text"}, {"init", "unused", "--path", ""},
		{"show", "absent"},
	} {
		if code, out, errOut := invoke(append([]string{"concept"}, args...)...); code != 2 || out != "" || errOut == "" {
			t.Errorf("usage %v: %d %q %q", args, code, out, errOut)
		}
	}
	if _, err := os.Stat("specs/concepts/unused.concept.spec.md"); !os.IsNotExist(err) {
		t.Fatal("invalid invocation created a specification")
	}
	for _, format := range []string{"text", "json"} {
		code, out, errOut := invoke("concept", "check", "absent", "--format", format)
		if code != 2 || out == "" || errOut != "" {
			t.Fatalf("missing check input: %d %q %q", code, out, errOut)
		}
	}
	writeFixture(t, "specs/concepts/broken.concept.spec.md", "broken")
	for _, args := range [][]string{{"show", "broken"}, {"list"}, {"check", "broken"}} {
		code, out, errOut := invoke(append([]string{"concept"}, args...)...)
		if args[0] == "check" {
			out, errOut = errOut, out
		}
		if code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
			t.Fatalf("structural error %v: %d %q %q", args, code, out, errOut)
		}
	}
	writeFixture(t, "specs/concepts/draft.concept.spec.md", "---\ntype: concept\n---\n# Draft\n")
	if code, out, errOut := invoke("concept", "check", "draft"); code != 0 || errOut != "" || !strings.Contains(out, "advice [concept.purpose]") {
		t.Fatalf("draft should succeed with advice: %d %q %q", code, out, errOut)
	}
}
