package obligation

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// MeasurementGuidance describes recorded sections, never their adequacy.
// Narrative criteria may carry or reference the complete comparison; structured
// criteria identify a comparison and target separately. Plans collect evidence
// against those criteria; their presence never substitutes for an obligation's
// criterion or applicability.
type MeasurementGuidance struct {
	Measurement, Criterion, Comparison, Conditions bool
}

func (g MeasurementGuidance) Missing() []string {
	if !g.Measurement {
		return nil
	}
	missing := []string{}
	if !g.Criterion {
		missing = append(missing, "criterion")
	} else if !g.Comparison {
		missing = append(missing, "comparison rule and target")
	}
	if !g.Conditions {
		missing = append(missing, "conditions or applicability")
	}
	return missing
}
func MeasurementMessage(missing []string) string {
	return "A measure or measurement plan is recorded, but no separately identified " + strings.Join(missing, ", ") + " is recorded for the obligation. Define or reference the required comparison under Criterion and its operating context under Applies to or Criterion. Assessment conditions describe how evidence is collected. Prose may already supply the information; this advice does not assess adequacy or establish a verification result."
}
func measurementGuidance(doc *specfile.Outline) MeasurementGuidance {
	sections := map[string]*specfile.Section{}
	for _, s := range doc.Sections {
		sections[s.Title] = s
	}
	body := func(s *specfile.Section) string {
		if s == nil {
			return ""
		}
		return doc.Body(s)
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
	measurement, criterion, applicability := sections["Measurement plan"], sections["Criterion"], sections["Applies to"]
	return MeasurementGuidance{
		Measurement: body(measurement) != "" || child(criterion, "Measure") != "",
		Criterion:   body(criterion) != "",
		Comparison:  intro(criterion) != "" || (child(criterion, "Comparison") != "" && child(criterion, "Target") != ""),
		Conditions:  intro(applicability) != "" || child(applicability, "Conditions") != "" || child(criterion, "Conditions") != "",
	}
}
