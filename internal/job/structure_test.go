package job

import (
	"strings"
	"testing"
)

const structuredJob = `---
type: job
id: prepare-change
---

# Prepare a change

## Job

Define requirements for a proposed change until the team agrees on a baseline.

## Performer

A product engineer working with stakeholders.

## Outcomes

- **OUT-001:** Minimize the likelihood of missing a stakeholder need.

## Steps

### Identify stakeholders

**ID:** identify-stakeholders
**Stage:** locate

Identify people whose needs may be affected.

#### Outcomes

- **OUT-002:** Minimize the time needed to identify
  affected stakeholders.

### Agree on the baseline

**ID:** agree-baseline
**Stage:** team-specific-stage

Reach agreement with stakeholders.

#### Outcomes

- **OUT-003:** Minimize the time needed to reach agreement.

## Progress cases

### Trying a new approach

#### Outcomes

- **CASE-NOTE:** This is ordinary case prose, not a job outcome.

## Evidence

Working hypotheses. See [the first step](#identify-stakeholders),
[its outcome](#OUT-002), and [the job statement](#job).
`

func TestStructuredJob(t *testing.T) {
	got, err := Parse([]byte(structuredJob))
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "prepare-change" || got.Performer == "" || got.Statement == "" {
		t.Fatalf("missing job fields: %+v", got)
	}
	if len(got.Outcomes) != 1 || got.Outcomes[0].ID != "OUT-001" {
		t.Fatalf("job outcomes: %+v", got.Outcomes)
	}
	if len(got.Steps) != 2 {
		t.Fatalf("steps: %+v", got.Steps)
	}
	first := got.Steps[0]
	if first.ID != "identify-stakeholders" || first.Stage != "locate" || first.Description != "Identify people whose needs may be affected." {
		t.Fatalf("step fields: %+v", first)
	}
	if len(first.Outcomes) != 1 || first.Outcomes[0].ID != "OUT-002" || !strings.Contains(first.Outcomes[0].Statement, "affected stakeholders.") {
		t.Fatalf("step outcomes: %+v", first.Outcomes)
	}
	if got.Steps[1].Stage != "team-specific-stage" {
		t.Fatal("custom stage lost")
	}
	if !strings.Contains(got.ProgressCases, "CASE-NOTE") || !strings.Contains(got.Markdown, "CASE-NOTE") {
		t.Fatal("case prose lost")
	}
	if suggestions := CheckWithPreset(got, "odi"); len(suggestions) != 0 {
		t.Fatalf("unexpected suggestions: %+v", suggestions)
	}
}

func TestDraftsAndLegacyStoriesRemainValid(t *testing.T) {
	for _, source := range []string{
		"---\ntype: job\n---\n# A job\n",
		validJob,
		"---\ntype: job\n---\n# A job\n## Steps\n### An incomplete step\n",
		"---\ntype: job\n---\n# A job\n## Job\n\n## Outcomes\n- A preliminary success criterion.\n",
	} {
		got, err := Parse([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		if len(Check(got)) == 0 {
			t.Fatal("expected advisory suggestions for the draft")
		}
	}
}

func TestStepIdentitySurvivesRenameAndReorder(t *testing.T) {
	firstStart := strings.Index(structuredJob, "### Identify stakeholders")
	secondStart := strings.Index(structuredJob, "### Agree on the baseline")
	end := strings.Index(structuredJob, "## Progress cases")
	source := structuredJob[:firstStart] + structuredJob[secondStart:end] + structuredJob[firstStart:secondStart] + structuredJob[end:]
	source = strings.Replace(source, "### Identify stakeholders", "### Identify affected people", 1)
	got, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if got.Steps[0].ID != "agree-baseline" || got.Steps[1].ID != "identify-stakeholders" || got.Steps[1].Outcomes[0].ID != "OUT-002" {
		t.Fatal("identities changed with document order")
	}
}

func TestHeadingsAndMetadataInExamplesAreNotStructure(t *testing.T) {
	example := "\n```markdown\n### Fake step\n**ID:** identify-stakeholders\n#### Outcomes\n- **OUT-001:** Fake outcome.\n```\n\n> ### Quoted step\n> **ID:** agree-baseline\n"
	source := strings.Replace(structuredJob, "Identify people whose needs may be affected.", "Identify people whose needs may be affected.\n"+example, 1)
	got, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Steps) != 2 || len(got.Steps[0].Outcomes) != 1 {
		t.Fatal("parsed examples as structure")
	}
	if !strings.Contains(got.Steps[0].Description, "**ID:** identify-stakeholders") {
		t.Fatal("removed metadata from code example")
	}
}

func TestStructuralErrors(t *testing.T) {
	tests := []struct{ source, want string }{
		{strings.Replace(structuredJob, "**ID:** agree-baseline", "**ID:** identify-stakeholders", 1), "duplicate ID"},
		{strings.Replace(structuredJob, "**OUT-003:**", "**OUT-001:**", 1), "duplicate ID"},
		{strings.Replace(structuredJob, "**OUT-003:**", "**identify-stakeholders:**", 1), "duplicate ID"},
		{strings.Replace(structuredJob, "**ID:** agree-baseline", "**ID:** invalid id", 1), "invalid ID"},
		{strings.Replace(structuredJob, "id: prepare-change", "id: bad/id", 1), "invalid ID"},
		{strings.Replace(structuredJob, "**ID:** agree-baseline", "**ID:** agree-baseline\n**ID:** another", 1), "duplicate ID field"},
		{strings.Replace(structuredJob, "**Stage:** locate", "**Stage:** locate\n**Stage:** define", 1), "duplicate Stage field"},
		{structuredJob + "\n## Steps\n", "duplicate Steps"},
		{strings.Replace(structuredJob, "### Identify stakeholders", "#### Identify stakeholders", 1), "level-three"},
		{strings.Replace(structuredJob, "### Identify stakeholders", "###", 1), "nonempty level-three"},
		{strings.Replace(structuredJob, "#### Outcomes", "##### Outcomes", 1), "#### Outcomes"},
		{strings.Replace(structuredJob, "### Agree on the baseline", "#### Outcomes\n\n### Agree on the baseline", 1), "duplicate Outcomes"},
		{strings.Replace(structuredJob, "(#OUT-002)", "(#missing-outcome)", 1), "unresolved local reference"},
		{strings.Replace(structuredJob, "type: job", "type: job\nprofiles: [odi]", 1), "invalid YAML"},
	}
	for _, tt := range tests {
		_, err := Parse([]byte(tt.source))
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("want %q, got %v", tt.want, err)
		}
	}
}

func TestReferencesAndUnknownSections(t *testing.T) {
	source := structuredJob + "\n## Additional notes\n\nText, [external](https://example.com), [file](notes.md), [local heading](#additional-notes), and [encoded ID](#OUT%2D002).\n"
	got, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.Markdown, "## Additional notes") {
		t.Fatal("unknown section lost")
	}
}

func TestCheckExplainsHeuristics(t *testing.T) {
	source := "---\ntype: job\n---\n# Maintain requirements\n## Job\nManage requirements.\n## Outcomes\n- **O1:** Make things better.\n## Steps\n### Review them\n"
	got, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	suggestions := CheckWithPreset(got, "odi")
	if len(suggestions) != 3 {
		t.Fatalf("want boundary, performer, and authored outcome suggestions; got %+v", suggestions)
	}
}

func TestCheckAllowsOptionalStructure(t *testing.T) {
	base := "---\ntype: job\n---\n# Prepare a change\n\n## Job\n\nPrepare a product change from proposal to review.\n\n## Performer\n\nA product engineer.\n"
	for _, optional := range []string{
		"",
		"\n## Outcomes\n\n## Steps\n",
		"\n## Outcomes\n\n- Minimize the time needed to assess the proposed change.\n",
		"\n## Steps\n\n### Locate commitments\n\n**Stage:** locate\n\n### Locate supporting evidence\n\n**Stage:** locate\n\n#### Outcomes\n",
		"\n## Steps\n\n### Locate commitments\n\n#### Outcomes\n\n- Minimize the likelihood of overlooking an applicable commitment.\n",
	} {
		got, err := Parse([]byte(base + optional))
		if err != nil {
			t.Fatal(err)
		}
		if suggestions := Check(got); len(suggestions) != 0 {
			t.Fatalf("optional structure should not prompt completion: %q: %+v", optional, suggestions)
		}
	}
}

func FuzzStructuredJob(f *testing.F) {
	for _, source := range []string{structuredJob, validJob, Template, "---\ntype: job\n---\n# Draft\n"} {
		f.Add([]byte(source))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		parsed, err := Parse(data)
		if err == nil {
			_ = Check(parsed)
		}
	})
}
