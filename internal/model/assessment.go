package model

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/obligation"
)

func (d *designResolver) richIntro(value any) string {
	nodes := []string{}
	for _, id := range Strings(Object(value)["content_nodes"]) {
		if String(d.index[id], "node_kind") == "heading" {
			break
		}
		nodes = append(nodes, id)
	}
	return strings.TrimSpace(d.text(Record{"content_nodes": nodes}))
}
func (d *designResolver) measurementGuidance(r Record) obligation.MeasurementGuidance {
	measurement, criterion, applicability := Object(r["measurement_plan"]), Object(r["criterion"]), Object(r["applicability"])
	present := func(v any) bool { return strings.TrimSpace(d.text(v)) != "" }
	return obligation.MeasurementGuidance{
		Measurement: present(measurement["content"]) || present(criterion["measure"]),
		Criterion:   present(criterion["content"]),
		Comparison:  d.richIntro(criterion["content"]) != "" || (present(criterion["comparison"]) && present(criterion["target"])),
		Conditions:  d.richIntro(applicability["content"]) != "" || present(applicability["conditions"]) || present(criterion["conditions"]),
	}
}
