package job

import (
	"strings"
	"testing"
)

const validJob = "---\ntype: job\n---\n\n# Develop Requirements\n\n## Situation\n\nA capability is changing.\n\n## Motivation\n\nUnderstand stakeholder needs.\n\n## Desired outcome\n\nAgree on what to build.\n"

func TestParseJob(t *testing.T) {
	for _, source := range []string{
		validJob,
		strings.Replace(validJob, "## Desired outcome", "## Intended progress", 1),
		"\xef\xbb\xbf" + strings.ReplaceAll(validJob, "\n", "\r\n"),
		strings.Replace(validJob, "type: job", "type: 'job' # JTBD", 1),
		validJob + "\n## Notes\nAn additional section.\n",
	} {
		got, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "Develop Requirements" || got.Situation != "A capability is changing." || got.Motivation != "Understand stakeholder needs." || got.IntendedProgress != "Agree on what to build." {
			t.Fatalf("unexpected job: %+v", got)
		}
	}
}

func TestIntendedProgressAliases(t *testing.T) {
	for _, heading := range []string{"Intended progress", "Desired outcome"} {
		t.Run(heading, func(t *testing.T) {
			body := "# A job\n\n## " + heading + "\n\nConfidence in the agreed change.\n\n## Notes\n\nSee [the progress](#" + strings.ReplaceAll(strings.ToLower(heading), " ", "-") + ").\n"
			got, err := Parse([]byte("---\ntype: job\n---\n" + body))
			if err != nil {
				t.Fatal(err)
			}
			if got.IntendedProgress != "Confidence in the agreed change." || got.Markdown != strings.TrimSpace(body) {
				t.Fatalf("progress or authored Markdown changed: %+v", got)
			}
			if len(got.Outcomes) != 0 {
				t.Fatalf("narrative progress became an outcome: %+v", got.Outcomes)
			}
			for _, duplicate := range []string{"Intended progress", "Desired outcome"} {
				_, err := Parse([]byte("---\ntype: job\n---\n" + body + "\n## " + duplicate + "\n"))
				if err == nil || !strings.Contains(err.Error(), "duplicate Intended progress section") {
					t.Fatalf("%s + %s: expected duplicate error even for an empty section, got %v", heading, duplicate, err)
				}
			}
		})
	}
}

func TestParseJobPreservesMarkdown(t *testing.T) {
	body := "- Understand needs.\n\n### Details\n\n```markdown\n# Not a title\n## Desired outcome\n```\n\n> ## Situation\n> Quoted example."
	source := strings.Replace(validJob, "Understand stakeholder needs.", body, 1)
	got, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if got.Motivation != body {
		t.Fatalf("Markdown changed: %q", got.Motivation)
	}
}

func TestParseJobRejectsMalformedContent(t *testing.T) {
	for _, section := range []string{"Situation", "Motivation", "Desired outcome"} {
		if _, err := Parse([]byte(validJob + "\n## " + section + "\nDuplicate.\n")); err == nil || !strings.Contains(err.Error(), "duplicate") {
			t.Fatalf("expected duplicate section error, got %v", err)
		}
	}
	for _, tt := range []struct{ source, want string }{
		{"", "frontmatter"},
		{strings.Replace(validJob, "type: job", "type: system", 1), "type: job"},
		{strings.Replace(validJob, "type: job", "type: [", 1), "invalid YAML"},
		{strings.Replace(validJob, "type: job", "type: job\nowner: someone", 1), "invalid YAML"},
		{strings.Replace(validJob, "# Develop Requirements", "#", 1), "nonempty"},
		{validJob + "\n# Another Job\n", "exactly one"},
		{strings.Replace(validJob, "# Develop Requirements", "```\n# Develop Requirements\n```", 1), "name heading"},
	} {
		if _, err := Parse([]byte(tt.source)); err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Fatalf("wanted %q error, got %v", tt.want, err)
		}
	}
}
