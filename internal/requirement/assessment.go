package requirement

import (
	"strings"

	"github.com/dotspec-md/dotspec/internal/specfile"
)

// MeasurementGuidance describes recorded sections, never their adequacy.
// Narrative criteria may carry or reference the complete rule; structured
// criteria identify a measure and its target separately.
type MeasurementGuidance struct {
	Measurement, Comparison, Conditions bool
}

func (g MeasurementGuidance) Missing() []string {
	if !g.Measurement {
		return nil
	}
	missing := []string{}
	if !g.Comparison {
		missing = append(missing, "target")
	}
	if !g.Conditions {
		missing = append(missing, "conditions or applicability")
	}
	return missing
}
func MeasurementMessage(missing []string) string {
	return "A measure is recorded, but no separately identified " + strings.Join(missing, ", ") + " is recorded for the requirement. Define or reference the required bound under Criterion and its operating context under Applies to."
}
func measurementGuidance(doc *specfile.Outline) MeasurementGuidance {
	sections := map[string]*specfile.Section{}
	for _, s := range doc.Sections {
		sections[s.Title] = s
	}
	intro := func(s *specfile.Section) string {
		if s == nil {
			return ""
		}
		return strings.TrimSpace(string(doc.Source[s.ContentStart:doc.IntroEnd(s)]))
	}
	child := func(s *specfile.Section, title string) string {
		if s == nil {
			return ""
		}
		values := []string{}
		for _, c := range s.Children {
			if c.Title == title {
				values = append(values, doc.Body(c))
			}
		}
		return strings.TrimSpace(strings.Join(values, "\n"))
	}
	criterion, applicability := sections["Criterion"], sections["Applies to"]
	return MeasurementGuidance{
		Measurement: child(criterion, "Measure") != "",
		Comparison:  intro(criterion) != "" || child(criterion, "Target") != "",
		Conditions:  intro(applicability) != "" || child(applicability, "Conditions") != "",
	}
}
