package model

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/requirement"
)

func (d *designResolver) richIntro(value any) string {
	var intro []Record
	for _, span := range records(Object(value)["spans"]) {
		raw := d.text(Record{"spans": []Record{span}})
		if strings.HasPrefix(strings.TrimSpace(raw), "#") {
			break
		}
		intro = append(intro, span)
	}
	return strings.TrimSpace(d.text(Record{"spans": intro}))
}
func (d *designResolver) measurementGuidance(r Record) requirement.MeasurementGuidance {
	criterion, applicability := Object(r["criterion"]), Object(r["applicability"])
	present := func(v any) bool { return strings.TrimSpace(d.text(v)) != "" }
	return requirement.MeasurementGuidance{
		Measurement: present(criterion["measure"]),
		Comparison:  d.richIntro(criterion["content"]) != "" || present(criterion["target"]),
		Conditions:  d.richIntro(applicability["content"]) != "" || present(applicability["conditions"]),
	}
}
