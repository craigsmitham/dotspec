package job

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

var boundaryWords = regexp.MustCompile(`(?i)\b(for|when|during|until|from|within)\b`)
var directionWords = regexp.MustCompile(`(?i)\b(minimize|maximize|reduce|increase|decrease|lower|improve)\b`)
var measureWords = regexp.MustCompile(`(?i)\b(time|likelihood|probability|number|proportion|percentage|amount|cost|effort|frequency|rate|duration|accuracy|completeness)\b`)

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
		add("job.statement", "Job", "Add a ## Job statement describing the accomplishment and where it begins and ends.", sectionLocation("Job"))
	} else if !boundaryWords.MatchString(job.Statement) {
		add("job.boundary", "Job", "Consider clarifying the execution boundary; no context marker such as 'for', 'when', or 'during' was found.", sectionLocation("Job"))
	}
	if strings.TrimSpace(job.Performer) == "" {
		add("job.performer", "Performer", "Identify who is trying to accomplish the job.", sectionLocation("Performer"))
	}
	checkOutcomes := func(outcomes []Outcome, location string) {
		for i, outcome := range outcomes {
			label := outcome.ID
			if label == "" {
				label = fmt.Sprintf("outcome %d", i+1)
			}
			if strings.TrimSpace(outcome.Statement) == "" {
				add("job.outcome-empty", location+" / "+label, "Describe the success criterion.", outcome.Location)
			} else if !directionWords.MatchString(outcome.Statement) || !measureWords.MatchString(outcome.Statement) {
				add("job.outcome-measure", location+" / "+label, "Consider a direction of improvement and a measure of success; the wording does not contain both recognizable elements.", outcome.Location)
			}
		}
	}
	checkOutcomes(job.Outcomes, "Overall job outcomes")
	for _, step := range job.Steps {
		checkOutcomes(step.Outcomes, "Outcomes for step: "+step.Name)
	}
	return suggestions
}
