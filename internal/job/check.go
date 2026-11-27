package job

import (
	"fmt"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
)

// Check provides transparent writing heuristics, not proof of semantic quality.
// Structural errors are rejected by Parse; suggestions never invalidate a draft.
func Check(job *Job) []diagnostic.Diagnostic {
	var suggestions []diagnostic.Diagnostic
	add := func(rule, element, message string, location diagnostic.Location) {
		location.Path, location.Element = job.Path, element
		suggestions = append(suggestions, diagnostic.Advice(rule, message, location))
	}
	sectionLocation := func(name string) diagnostic.Location {
		if location, found := job.SectionLocations[name]; found {
			return location
		}
		return job.Location
	}
	if strings.TrimSpace(job.Statement) == "" {
		add("job.statement", "Job", "Add a ## Job statement describing what the performer seeks to accomplish or change.", sectionLocation("Job"))
	}
	if strings.TrimSpace(job.Performer) == "" {
		add("job.performer", "Performer", "Identify who is trying to accomplish the job.", sectionLocation("Performer"))
	}
	checkOutcomes := func(outcomes []DesiredOutcome, location string) {
		for i, outcome := range outcomes {
			label := outcome.ID
			if label == "" {
				label = fmt.Sprintf("outcome %d", i+1)
			}
			if strings.TrimSpace(outcome.Statement) == "" {
				add("job.outcome-empty", location+" / "+label, "Describe the success criterion.", outcome.Location)
			}
		}
	}
	checkOutcomes(job.Outcomes, "Outcomes declared by job")
	for _, step := range job.Steps {
		checkOutcomes(step.Outcomes, "Outcomes declared by step: "+step.Name)
	}
	return suggestions
}
