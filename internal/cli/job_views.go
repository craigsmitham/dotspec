package cli

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/craigsmitham/dotspec/internal/job"
)

func stepLabel(step job.Step) string {
	label := step.Name
	if step.ID != "" {
		label += " [" + step.ID + "]"
	}
	if step.Stage != "" {
		label += " (" + step.Stage + ")"
	}
	return label
}

func outcomeLabel(outcome job.Outcome) string {
	if outcome.ID == "" {
		return outcome.Statement
	}
	return outcome.ID + ": " + outcome.Statement
}

func listSteps(spec *job.Job, out io.Writer) error {
	if len(spec.Steps) == 0 {
		_, err := fmt.Fprintln(out, "No steps defined.")
		return err
	}
	var output strings.Builder
	fmt.Fprintf(&output, "Steps: %s\n", spec.Name)
	for i, step := range spec.Steps {
		fmt.Fprintf(&output, "\n%d. %s\n", i+1, stepLabel(step))
		if step.Description != "" {
			fmt.Fprintln(&output, indent(step.Description))
		}
	}
	_, err := io.WriteString(out, output.String())
	return err
}

func listOutcomes(spec *job.Job, stepID string, out io.Writer) error {
	var output strings.Builder
	count := 0
	write := func(label string, outcomes []job.Outcome) {
		if len(outcomes) == 0 {
			return
		}
		fmt.Fprintf(&output, "\n%s\n", label)
		for _, outcome := range outcomes {
			fmt.Fprintln(&output, indent(outcomeLabel(outcome)))
			count++
		}
	}
	fmt.Fprintf(&output, "Outcomes: %s\n", spec.Name)
	if stepID == "" {
		write("Job", spec.Outcomes)
	}
	found := stepID == ""
	for _, step := range spec.Steps {
		if stepID != "" && step.ID != stepID {
			continue
		}
		found = true
		write("Step: "+stepLabel(step), step.Outcomes)
	}
	if !found {
		return fmt.Errorf("step %q not found in job %s", stepID, spec.Name)
	}
	if count == 0 {
		_, err := fmt.Fprintln(out, "No outcomes defined for the selection.")
		return err
	}
	_, err := io.WriteString(out, output.String())
	return err
}

func writeJobMap(spec *job.Job, options jobOptions, out io.Writer) error {
	if len(spec.Steps) == 0 {
		_, err := fmt.Fprintln(out, "No steps defined.")
		return err
	}
	var output strings.Builder
	if options.format == "mermaid" {
		output.WriteString("flowchart TD\n")
		// Internal node names avoid interpreting author IDs as Mermaid syntax.
		for i, step := range spec.Steps {
			fmt.Fprintf(&output, "  s%d[\"%s\"]\n", i+1, mermaidLabel(stepLabel(step)))
			if i > 0 {
				fmt.Fprintf(&output, "  s%d --> s%d\n", i, i+1)
			}
			if options.outcomes {
				for j, outcome := range step.Outcomes {
					fmt.Fprintf(&output, "  o%d_%d[\"%s\"]\n  s%d -.-> o%d_%d\n", i+1, j+1, mermaidLabel(outcomeLabel(outcome)), i+1, i+1, j+1)
				}
			}
		}
		if options.outcomes && len(spec.Outcomes) > 0 {
			fmt.Fprintf(&output, "  job[\"%s\"]\n  job --> s1\n", mermaidLabel(spec.Name))
			for i, outcome := range spec.Outcomes {
				fmt.Fprintf(&output, "  jo%d[\"%s\"]\n  job -.-> jo%d\n", i+1, mermaidLabel(outcomeLabel(outcome)), i+1)
			}
		}
	} else {
		fmt.Fprintf(&output, "Job map: %s\n", spec.Name)
		if options.outcomes && len(spec.Outcomes) > 0 {
			output.WriteString("\nJob outcomes\n")
			for _, outcome := range spec.Outcomes {
				fmt.Fprintln(&output, indent(outcomeLabel(outcome)))
			}
		}
		for i, step := range spec.Steps {
			if i == 0 {
				output.WriteByte('\n')
			} else {
				output.WriteString("   |\n   v\n")
			}
			fmt.Fprintf(&output, "%d. %s\n", i+1, stepLabel(step))
			if options.outcomes {
				for _, outcome := range step.Outcomes {
					fmt.Fprintln(&output, indent(outcomeLabel(outcome)))
				}
			}
		}
	}
	_, err := io.WriteString(out, output.String())
	return err
}

// Mermaid's decimal entity notation prevents labels from becoming diagram syntax.
func mermaidLabel(label string) string {
	var output strings.Builder
	for _, r := range strings.Join(strings.Fields(label), " ") {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '.' || r == ',' || r == ':' {
			output.WriteRune(r)
		} else {
			fmt.Fprintf(&output, "#%d;", r)
		}
	}
	return output.String()
}
