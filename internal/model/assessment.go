package model

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/obligation"
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
func (d *designResolver) measurementGuidance(r Record) obligation.MeasurementGuidance {
	criterion, applicability := Object(r["criterion"]), Object(r["applicability"])
	present := func(v any) bool { return strings.TrimSpace(d.text(v)) != "" }
	return obligation.MeasurementGuidance{
		Measurement: present(criterion["measure"]),
		Criterion:   present(criterion["content"]),
		Comparison:  d.richIntro(criterion["content"]) != "" || present(criterion["target"]),
		Conditions:  d.richIntro(applicability["content"]) != "" || present(applicability["conditions"]),
	}
}
