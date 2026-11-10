package usecase

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const draft = "---\ntype: use-case\n---\n# Check specifications\n"

const example = `---
type: use-case
id: check-specifications
---
# Check specifications

## Goal
Identify structural problems before committing a change.

## System boundary
The dotspec CLI.

## Primary actor
A contributor preparing a change.

## Goal level
User task (sea level).

## Trigger
The contributor requests a check.

## Preconditions
The selected project directory exists and is readable.

## Success criterion
The contributor receives a complete report for the selected scope.

## Main success scenario
1. The contributor selects a project directory.
2. Dotspec inspects its supported specifications.
3. Dotspec reports findings and scan completion.

## Extensions
### Unreadable input
Dotspec reports an incomplete scan.

### Preconditions
This nested heading describes an extension; it is not the main preconditions.

## Sequence
Representative. The [reporting contract](../../README.md#check-specifications) defines required behavior.

## Related jobs
[Evaluate requirements](../jobs/develop-requirements.job.spec.md#evaluate-requirements).

## Notes
See [the goal](#goal), [this use case](#check%2Dspecifications), and [external evidence](https://example.com).
`

func TestParsePreservesUseCaseAndDistinctFields(t *testing.T) {
	got, err := Parse([]byte(example))
	if err != nil {
		t.Fatal(err)
	}
	body := example[strings.Index(example, "# Check specifications"):]
	if got.Markdown != strings.TrimSpace(body) {
		t.Fatalf("authored Markdown changed: %q", got.Markdown)
	}
	if got.ID != "check-specifications" || got.Name != "Check specifications" || got.System != nil {
		t.Fatalf("identity or association: %+v", got)
	}
	want := map[string]string{
		"goal":              "Identify structural problems before committing a change.",
		"system boundary":   "The dotspec CLI.",
		"primary actor":     "A contributor preparing a change.",
		"goal level":        "User task (sea level).",
		"trigger":           "The contributor requests a check.",
		"preconditions":     "The selected project directory exists and is readable.",
		"success criterion": "The contributor receives a complete report for the selected scope.",
	}
	actual := map[string]string{"goal": got.Goal, "system boundary": got.SystemBoundary, "primary actor": got.PrimaryActor, "goal level": got.GoalLevel, "trigger": got.Trigger, "preconditions": got.Preconditions, "success criterion": got.SuccessCriterion}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("fields: %v; want %v", actual, want)
	}
	if !strings.HasPrefix(got.MainSuccessScenario, "1. The contributor") || !strings.Contains(got.Extensions, "### Preconditions") || !strings.Contains(got.RelatedJobs, "#evaluate-requirements") {
		t.Fatalf("lost scenario, extension, or job link: %+v", got)
	}
	if len(Check(got)) != 0 {
		t.Fatalf("unexpected advice: %+v", Check(got))
	}
}

func TestDraftsAndAdvice(t *testing.T) {
	for _, source := range []string{draft, draft + "\n## Primary actor\n", draft + "\n## Extensions\nAn unresolved failure condition.\n"} {
		spec, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		spec.Path = "specs/use-cases/check-specifications.use-case.spec.md"
		advice := Check(spec)
		if len(advice) != 5 {
			t.Fatalf("missing content advice: %+v", advice)
		}
		for _, d := range advice {
			if d.Severity != "advice" || d.Location.Path != spec.Path || d.Location.Line < 4 || d.Location.Column != 1 {
				t.Fatalf("invalid advice location/severity: %+v", d)
			}
			if d.Rule == "use-case.primary-actor" && strings.Contains(source, "## Primary actor") && d.Location.Line != 6 {
				t.Fatalf("empty section should point at its heading: %+v", d)
			}
		}
	}
	spec, err := Parse([]byte(strings.Replace(example, "## Sequence", "## Ordering notes", 1)))
	if err != nil {
		t.Fatal(err)
	}
	if advice := Check(spec); len(advice) != 1 || advice[0].Rule != "use-case.sequence" {
		t.Fatalf("scenario order advice: %+v", advice)
	}
}

func TestStructureAndSourceLocations(t *testing.T) {
	for _, tt := range []struct {
		name, source, rule    string
		line, column, related int
	}{
		{"frontmatter", "# Draft", "metadata.frontmatter", 1, 1, 0},
		{"type", strings.Replace(draft, "type: use-case", "type: job", 1), "metadata.type", 2, 7, 0},
		{"unknown field", strings.Replace(draft, "type: use-case", "type: use-case\nextra: true", 1), "metadata.invalid", 3, 1, 0},
		{"invalid ID", strings.Replace(draft, "type: use-case", "type: use-case\nid: bad/id", 1), "identity.invalid", 3, 5, 0},
		{"duplicate section", draft + "## Goal\n\n## Goal\n", "structure.duplicate-section", 7, 1, 5},
		{"wrong level", draft + "### Goal\n", "structure.section-level", 5, 1, 0},
		{"second title", draft + "# Another\n", "structure.title", 5, 1, 4},
		{"before title", "---\ntype: use-case\n---\n## Goal\n# Title\n", "structure.section-order", 4, 1, 0},
		{"unresolved reference", draft + "## Notes\n[Unknown](#missing)\n", "reference.unresolved", 6, 1, 0},
		{"unicode reference", draft + "## Notes\né [Unknown](#missing)\n", "reference.unresolved", 6, 3, 0},
		{"malformed fragment", draft + "## Notes\n[Unknown](#bad%XX)\n", "reference.invalid", 6, 1, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, source := range []string{tt.source, "\ufeff" + strings.ReplaceAll(tt.source, "\n", "\r\n")} {
				_, err := Parse([]byte(source))
				var d *diagnostic.Diagnostic
				if !errors.As(err, &d) || d.Rule != tt.rule || d.Location.Line != tt.line || d.Location.Column != tt.column {
					t.Fatalf("diagnostic: %+v, %v; want %s %d:%d", d, err, tt.rule, tt.line, tt.column)
				}
				if tt.related > 0 && (len(d.Related) != 1 || d.Related[0].Line != tt.related) {
					t.Fatalf("missing previous definition: %+v", d)
				}
			}
		})
	}
}

func TestEveryConventionalSectionRejectsDuplicates(t *testing.T) {
	for _, name := range []string{"Goal", "System boundary", "Primary actor", "Goal level", "Trigger", "Preconditions", "Success criterion", "Main success scenario", "Extensions", "Sequence", "Related jobs", "Supporting actors", "Stakeholder interests", "Minimum guarantees"} {
		_, err := Parse([]byte(draft + "## " + name + "\n\n## " + name + "\n"))
		var d *diagnostic.Diagnostic
		if !errors.As(err, &d) || d.Rule != "structure.duplicate-section" {
			t.Fatalf("duplicate %s: %v", name, err)
		}
	}
}

func TestExamplesAndUnknownProseRemainMarkdown(t *testing.T) {
	extra := "\n## Notes\n\n```markdown\n# Example\n## Goal\n[not checked](#missing)\n```\n\n> ## Goal\n> quoted content\n\n## goal\nLowercase headings are ordinary prose.\n"
	source := strings.Replace(example, "## Notes", "## Other notes", 1) + extra
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if spec.Goal != "Identify structural problems before committing a change." || !strings.Contains(spec.Markdown, strings.TrimSpace(extra)) {
		t.Fatalf("examples changed parsed content: %+v", spec)
	}
	// A link to another file is navigation, even if that file or fragment is missing.
	if _, err := Parse([]byte(draft + "[elsewhere](missing.spec.md#unknown)\n")); err != nil {
		t.Fatal(err)
	}
}

func FuzzParse(f *testing.F) {
	for _, source := range []string{draft, example, Template, ExpandedTemplate, draft + "## Goal\n\n## Goal\n", draft + "[bad](#missing)\n"} {
		f.Add([]byte(source))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if spec, err := Parse(data); err == nil {
			_ = Check(spec)
		}
	})
}

func TestProgressiveExtensionsAndGuarantees(t *testing.T) {
	source := draft + "## Main success scenario\nThe actor requests service; the system supplies it.\n## Extensions\n"
	spec, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range Check(spec) {
		if d.Rule == "use-case.extensions" {
			found = true
			if d.Severity != "advice" || d.Location.Line != 7 {
				t.Fatalf("extension advice: %+v", d)
			}
		}
	}
	if !found {
		t.Fatal("missing extension advice after a scenario")
	}
	source += "### Alternative service\nCondition: preferred service unavailable. Handling unresolved.\n" +
		"## Minimum guarantees\nNo unconfirmed promise.\n## Supporting actors\nExternal provider.\n## Stakeholder interests\nOperator: accurate commitments.\n" +
		"## Evidence\n[Branch](#alternative-service)\n"
	spec, err = Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if spec.MinimumGuarantees != "No unconfirmed promise." || spec.SupportingActors != "External provider." || spec.StakeholderInterests != "Operator: accurate commitments." {
		t.Fatalf("lost context: %+v", spec)
	}
	for _, d := range Check(spec) {
		if d.Rule == "use-case.extensions" {
			t.Fatal("condition-only branch should be accepted")
		}
	}
	for _, template := range []string{Template, ExpandedTemplate} {
		if _, err := Parse([]byte(template)); err != nil {
			t.Fatal(err)
		}
	}
}
