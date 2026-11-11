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
func Check(job *Job) []diagnostic.Diagnostic { return CheckWithPreset(job, "general") }

func CheckWithPreset(job *Job, preset string) []diagnostic.Diagnostic {
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
	} else if preset == "odi" && !boundaryWords.MatchString(job.Statement) {
		add("job.boundary", "Job", "Consider clarifying the execution boundary; no context marker such as 'for', 'when', or 'during' was found.", sectionLocation("Job"))
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
			} else if preset == "odi" && (!directionWords.MatchString(outcome.Statement) || !measureWords.MatchString(outcome.Statement)) {
				add("job.outcome-measure", location+" / "+label, "Consider a direction of improvement and a measure of success; the wording does not contain both recognizable elements.", outcome.Location)
			}
		}
	}
	checkOutcomes(job.Outcomes, "Outcomes declared by job")
	for _, step := range job.Steps {
		checkOutcomes(step.Outcomes, "Outcomes declared by step: "+step.Name)
	}
	if preset == "job-story" {
		for _, field := range []struct{ name, value string }{{"Situation", job.Situation}, {"Motivation", job.Motivation}, {"Intended progress", job.IntendedProgress}} {
			if strings.TrimSpace(field.value) == "" {
				add("job.story-context", field.name, "Describe "+strings.ToLower(field.name)+" for this job story.", sectionLocation(field.name))
			}
		}
	}
	if preset == "odi" && len(job.Steps) == 0 {
		add("job.map", "Steps", "Consider an optional functional job map; Define, Locate, Prepare, Confirm, Execute, Monitor, Modify and Conclude are elicitation prompts, not required schema slots.", sectionLocation("Steps"))
	}
	if preset == "progress" {
		if len(job.Cases) == 0 {
			add("job.progress-case", "Progress cases", "Consider a progress case recording circumstances, a struggling moment, alternatives and trade-offs; distinguish observation from hypothesis in Evidence.", sectionLocation("Progress cases"))
		}
		for _, c := range job.Cases {
			for _, field := range []struct{ name, value string }{{"Circumstances", c.Circumstances}, {"Struggling moment", c.StrugglingMoment}, {"Participants", c.Participants}, {"Alternatives considered", c.AlternativesConsidered}, {"Trade-offs and selection criteria", c.TradeoffsAndSelectionCriteria}} {
				if strings.TrimSpace(field.value) == "" {
					add("job.progress-context", c.Name+" / "+field.name, "Consider recording "+strings.ToLower(field.name)+" for this case.", sectionLocation("Progress cases"))
				}
			}
		}
	}
	return suggestions
}
