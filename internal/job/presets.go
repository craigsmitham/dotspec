package job

import "fmt"

// Presets guide authoring; they do not classify jobs or impose a methodology.
func ValidatePreset(preset string) error {
	switch preset {
	case "", "general", "odi", "job-story":
		return nil
	}
	return fmt.Errorf("unknown JTBD preset %q: choose general, odi, or job-story", preset)
}

func TemplateFor(preset string) (string, error) {
	if err := ValidatePreset(preset); err != nil {
		return "", err
	}
	switch preset {
	case "odi":
		return Template + `
## Steps

This is a partial functional map, not a required execution sequence.
Use Define, Locate, Prepare, Confirm, Execute, Monitor, Modify and Conclude
as elicitation prompts where useful; stages may repeat or be omitted.

### Accomplishment within the job

**ID:** first-step

Describe what the performer needs to accomplish, independently of the solution.

## Outcomes

### Success criterion

**ID:** success-criterion

Describe a direction of improvement, measure, object and relevant circumstances.
Keep quantified product commitments in requirements and research ratings in dated evidence.

#### Applies to

[Accomplishment within the job](#first-step)
`, nil
	case "job-story":
		return `---
type: job-spec
---

# New job

## Job

When [situation], I want to [motivation], so I can [intended progress].

## Performer

Describe who seeks this progress.

## Situation

Describe when and where the need arises.

## Motivation

Describe what the person wants to do and why.

## Intended progress

Describe the result they seek, independently of the chosen solution.
`, nil
	default:
		return Template, nil
	}
}
