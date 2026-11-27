package cli

import (
	"os"
	"strings"
	"testing"
)

func TestSpecKindUsageAndFailureChannels(t *testing.T) {
	for _, kind := range []struct {
		command, collection, specType, advice string
		extraUsage                            [][]string
	}{
		{command: "concept", collection: "concepts", specType: "concept-spec", advice: "concept.purpose"},
		{command: "use-case", collection: "use-cases", specType: "use-case-spec", advice: "use-case.goal",
			extraUsage: [][]string{{"init", "unused", "--template", "unknown"}}},
		{command: "business-rule", collection: "business-rules", specType: "business-rule-spec", advice: "business-rule.statement"},
	} {
		t.Run(kind.command, func(t *testing.T) {
			prepareExample(t, kind.collection)
			usage := append([][]string{
				{"init"}, {"init", "../escape"}, {"show", "UPPER"}, {"show", "work", "extra"},
				{"list", "extra"}, {"check"}, {"check", "../escape"}, {"list", "--format", "mermaid"},
				{"init", "unused", "--json", "--format", "text"}, {"init", "unused", "--path", ""},
				{"show", "absent"},
			}, kind.extraUsage...)
			for _, args := range usage {
				if code, out, errOut := invoke(append([]string{kind.command}, args...)...); code != 2 || out != "" || errOut == "" {
					t.Errorf("usage %v: %d %q %q", args, code, out, errOut)
				}
			}
			source := func(slug string) string {
				return "specs/" + kind.collection + "/" + slug + "." + kind.command + ".spec.md"
			}
			if _, err := os.Stat(source("unused")); !os.IsNotExist(err) {
				t.Fatal("invalid invocation created a specification")
			}
			for _, format := range []string{"text", "json"} {
				code, out, errOut := invoke(kind.command, "check", "absent", "--format", format)
				if code != 2 || out == "" || errOut != "" {
					t.Fatalf("missing check input: %d %q %q", code, out, errOut)
				}
			}
			writeFixture(t, source("broken"), "broken")
			for _, args := range [][]string{{"show", "broken"}, {"list"}, {"check", "broken"}} {
				code, out, errOut := invoke(append([]string{kind.command}, args...)...)
				if args[0] == "check" {
					out, errOut = errOut, out
				}
				if code != 1 || out != "" || !strings.Contains(errOut, "metadata.frontmatter") {
					t.Fatalf("structural error %v: %d %q %q", args, code, out, errOut)
				}
			}
			writeFixture(t, source("draft"), "---\ntype: "+kind.specType+"\n---\n# Draft\n")
			if code, out, errOut := invoke(kind.command, "check", "draft"); code != 0 || errOut != "" || !strings.Contains(out, "advice ["+kind.advice+"]") {
				t.Fatalf("draft should succeed with advice: %d %q %q", code, out, errOut)
			}
		})
	}
}
