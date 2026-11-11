package job

import "fmt"

// Presets guide authoring; they do not classify jobs or impose a methodology.
func ValidatePreset(preset string) error {
	switch preset {
	case "", "general", "odi", "progress", "job-story":
		return nil
	}
	return fmt.Errorf("unknown JTBD preset %q: choose general, odi, progress, or job-story", preset)
}

type ProgressCase struct {
	Name                          string `json:"name"`
	Circumstances                 string `json:"circumstances,omitempty"`
	StrugglingMoment              string `json:"struggling_moment,omitempty"`
	Participants                  string `json:"participants,omitempty"`
	AlternativesConsidered        string `json:"alternatives_considered,omitempty"`
	TradeoffsAndSelectionCriteria string `json:"tradeoffs_and_selection_criteria,omitempty"`
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
	case "progress":
		return Template + `
## Progress cases

### Candidate change

#### Participants

Distinguish the performer, buyer and decision maker where relevant.

#### Circumstances

Describe the particular context in which this case arose.

#### Struggling moment

Describe what made a change worth considering.

#### Current approach

Describe how the person currently makes progress.

#### Desired progress

Describe the change they seek, including emotional and social dimensions where relevant.

#### Alternatives considered

Include continuing the current approach or doing nothing where relevant.

#### Trade-offs and selection criteria

Describe criteria for choosing, hiring or rejecting an approach.

#### Forces

- Push: Describe what makes the current situation untenable.
- Pull: Describe what attracts the person to a new approach.
- Anxiety: Describe concerns about changing.
- Habit: Describe what keeps the current approach in place.

#### Timeline

Describe the decision to change; this is distinct from the functional job map.
Dates may be unknown.

#### Evidence

##### Research basis

State whether this case is observed or hypothesized and link dated sources.

##### Limitations

Describe what remains unknown and which population or circumstances the evidence covers.
`, nil
	case "job-story":
		return `---
type: job
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
