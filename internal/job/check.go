package job

import (
	"fmt"
	"regexp"
	"strings"
)

type Suggestion struct {
	Location, Message string
}

var boundaryWords = regexp.MustCompile(`(?i)\b(for|when|during|until|from|within)\b`)
var directionWords = regexp.MustCompile(`(?i)\b(minimize|maximize|reduce|increase|decrease|lower|improve)\b`)
var measureWords = regexp.MustCompile(`(?i)\b(time|likelihood|probability|number|proportion|percentage|amount|cost|effort|frequency|rate|duration|accuracy|completeness)\b`)

// Check provides transparent writing heuristics, not proof of semantic quality.
// Structural errors are rejected by Parse; suggestions never invalidate a draft.
func Check(job *Job) []Suggestion {
	var suggestions []Suggestion
	add := func(location, message string) { suggestions = append(suggestions, Suggestion{location, message}) }
	if strings.TrimSpace(job.Statement) == "" {
		add("Job", "Add a ## Job statement describing the accomplishment and where it begins and ends.")
	} else if !boundaryWords.MatchString(job.Statement) {
		add("Job", "Consider clarifying the execution boundary; no context marker such as 'for', 'when', or 'during' was found.")
	}
	if strings.TrimSpace(job.Performer) == "" {
		add("Performer", "Identify who is trying to accomplish the job.")
	}
	checkOutcomes := func(outcomes []Outcome, location string) {
		for i, outcome := range outcomes {
			label := outcome.ID
			if label == "" {
				label = fmt.Sprintf("outcome %d", i+1)
			}
			if strings.TrimSpace(outcome.Statement) == "" {
				add(location+" / "+label, "Describe the success criterion.")
			} else if !directionWords.MatchString(outcome.Statement) || !measureWords.MatchString(outcome.Statement) {
				add(location+" / "+label, "Consider a direction of improvement and a measure of success; the wording does not contain both recognizable elements.")
			}
		}
	}
	checkOutcomes(job.Outcomes, "Overall job outcomes")
	for _, step := range job.Steps {
		checkOutcomes(step.Outcomes, "Outcomes for step: "+step.Name)
	}
	return suggestions
}
