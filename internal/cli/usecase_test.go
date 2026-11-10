package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/check"
)

func TestUseCaseCommandsAndJSON(t *testing.T) {
	prepareExample(t, "dotspec")
	code, out, errOut := invoke("use-case", "list", "--json")
	if code != 0 || errOut != "" || !strings.Contains(out, `"groups": []`) {
		t.Fatalf("empty JSON: %d %s %s", code, out, errOut)
	}
	if code, out, errOut := invoke("use-case", "list"); code != 0 || errOut != "" || !strings.Contains(out, "No use cases found") {
		t.Fatalf("empty text: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("use-case", "init", "--path", ".", "review", "--json")
	var created createdResult
	if err := json.Unmarshal([]byte(out), &created); err != nil || code != 0 || errOut != "" || !filepath.IsAbs(created.Path) {
		t.Fatalf("created JSON: %d %s %s %v", code, out, errOut, err)
	}
	if code, _, _ := invoke("system", "init"); code != 0 {
		t.Fatal("system init failed")
	}
	code, out, errOut = invoke("use-case", "show", "--format", "json", "--", "review")
	var result useCaseResult
	if err := json.Unmarshal([]byte(out), &result); err != nil || code != 0 || errOut != "" || result.UseCase == nil {
		t.Fatalf("show JSON: %d %s %s %v", code, out, errOut, err)
	}
	if result.SchemaVersion != 1 || result.UseCase.Slug != "review" || result.UseCase.System == nil || result.UseCase.System.Name != "New system" || result.UseCase.Trigger == result.UseCase.Preconditions {
		t.Fatalf("lost fields or association: %+v", result)
	}
	code, out, errOut = invoke("use-case", "show", "review")
	if code != 0 || errOut != "" || !strings.Contains(out, result.UseCase.Markdown) || !strings.Contains(out, "System: New system") {
		t.Fatalf("show text: %d %s %s", code, out, errOut)
	}
	code, out, errOut = invoke("use-case", "list", "--json")
	var groups useCasesResult
	if err := json.Unmarshal([]byte(out), &groups); err != nil || code != 0 || errOut != "" || len(groups.Groups) != 1 || len(groups.Groups[0].UseCases) != 1 {
		t.Fatalf("list JSON: %d %s %s %v", code, out, errOut, err)
	}
	code, out, errOut = invoke("use-case", "check", "review", "--json")
	var report check.Report
	if err := json.Unmarshal([]byte(out), &report); err != nil || code != 0 || errOut != "" || !report.Complete || report.Summary.FilesChecked != 1 || report.Scope.Types[0] != "use-case" {
		t.Fatalf("focused report: %d %s %s %v", code, out, errOut, err)
	}
	before, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatal(err)
	}
	if code, _, _ := invoke("use-case", "init", "review"); code != 2 {
		t.Fatal("init overwrote an existing specification")
	}
	if after, err := os.ReadFile(created.Path); err != nil || string(before) != string(after) {
		t.Fatal("commands changed existing content")
	}
}

func TestUseCaseUsageAndFailureChannels(t *testing.T) {
	prepareExample(t, "use cases")
	for _, args := range [][]string{
		{"init"}, {"init", "../escape"}, {"show", "UPPER"}, {"show", "work", "extra"},
		{"list", "extra"}, {"check"}, {"check", "../escape"}, {"list", "--format", "mermaid"},
		{"init", "unused", "--json", "--format", "text"}, {"init", "unused", "--path", ""},
		{"show", "absent"},
	} {
		if code, out, errOut := invoke(append([]string{"use-case"}, args...)...); code != 2 || out != "" || errOut == "" {
			t.Errorf("usage %v: %d %q %q", args, code, out, errOut)
		}
	}
	if _, err := os.Stat("specs/use-cases/unused.use-case.spec.md"); !os.IsNotExist(err) {
		t.Fatal("invalid invocation created a specification")
	}
	for _, format := range []string{"text", "json"} {
		code, out, errOut := invoke("use-case", "check", "absent", "--format", format)
		if code != 2 || out == "" || errOut != "" {
			t.Fatalf("missing check input: %d %q %q", code, out, errOut)
		}
	}
	writeFixture(t, "specs/use-cases/broken.use-case.spec.md", "broken")
	for _, args := range [][]string{{"show", "broken"}, {"list"}, {"check", "broken"}} {
		code, out, errOut := invoke(append([]string{"use-case"}, args...)...)
		if args[0] == "check" {
			out, errOut = errOut, out
		}
		if code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
			t.Fatalf("structural error %v: %d %q %q", args, code, out, errOut)
		}
	}
	writeFixture(t, "specs/use-cases/draft.use-case.spec.md", "---\ntype: use-case\n---\n# Draft\n")
	if code, out, errOut := invoke("use-case", "check", "draft"); code != 0 || errOut != "" || !strings.Contains(out, "advice [use-case.goal]") {
		t.Fatalf("draft should succeed with advice: %d %q %q", code, out, errOut)
	}
}
