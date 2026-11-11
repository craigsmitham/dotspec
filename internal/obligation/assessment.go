package obligation

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// MeasurementGuidance describes recorded sections, never their adequacy.
// Narrative criteria may carry the complete comparison; structured criteria
// identify a comparison and target separately. Neither requires duplication.
type MeasurementGuidance struct {
	Measurement, Verification, Criterion, Comparison, Conditions bool
}

func (g MeasurementGuidance) Missing() []string {
	if !g.Measurement || g.Verification {
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
	return "Measurement guidance is recorded, but no separately identified " + strings.Join(missing, ", ") + " is recorded. Review how observations will be compared with the obligation; prose may already supply the information. Guidance presence is not a verification result."
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
	success := child(measurement, "Success criteria") != ""
	return MeasurementGuidance{
		Measurement: body(measurement) != "", Verification: body(sections["Verification plan"]) != "",
		Criterion:  body(criterion) != "" || success,
		Comparison: intro(criterion) != "" || success || (child(criterion, "Comparison") != "" && child(criterion, "Target") != ""),
		Conditions: intro(applicability) != "" || child(applicability, "Conditions") != "" || child(criterion, "Conditions") != "" || child(measurement, "Conditions") != "",
	}
}
