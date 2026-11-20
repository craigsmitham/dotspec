package model

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/assessment"
)

// Each authored heading starts a distinct plan in source order.
func (m *domainMapper) assessmentSection(r Record, s *Section) bool {
	definition, known := assessment.FromHeading(s.Title)
	if !known {
		return false
	}
	if _, ok := compiledModel().catalog.Classes[String(r, "kind")].Slots["assessment_plans"]; !ok {
		return false
	}
	if m.traceMapped == nil {
		m.traceMapped = map[string]bool{}
	}
	key := String(r, "id") + "/assessment/" + s.ID
	if m.traceMapped[key] {
		return true
	}
	m.traceMapped[key] = true
	plans := records(r["assessment_plans"])
	field := fmt.Sprintf("assessment_plans.%d", len(plans))
	v := Record{"purpose": string(definition.Purpose)}
	r["assessment_plans"] = append(plans, v)
	m.setRich(v, "content", s.Content)
	m.bind(r, field+".content", s.Content, "named_field", "assessment-content")
	m.bind(r, field+".purpose", []string{s.HeadingID}, "named_field", "assessment-purpose")
	details := map[string]string{"Participants": "participants", "Method": "method", "Procedure": "method", "Success criteria": "success_criteria", "Conditions": "conditions", "Evidence to retain": "evidence_to_retain", "Evidence retention": "evidence_to_retain"}
	for _, sub := range s.Children {
		if f := details[sub.Title]; f != "" {
			m.setRich(v, f, sub.Content)
			m.sectionBindings(r, field+"."+f, sub)
			continue
		}
		switch sub.Title {
		case "Subject", "Basis":
			role := strings.ToLower(sub.Title)
			m.contextField(v, role, sub.Content)
			m.sectionBindings(r, field+"."+role+".content", sub)
			m.bindReferenceOccurrences(r, field+"."+role+".references", m.references(sub.Content))
		case "Cases":
			cases := reviewObjects(v["cases"])
			if cases == nil {
				cases = []Record{}
			}
			for _, child := range sub.Children {
				path := fmt.Sprintf("%s.cases.%d", field, len(cases))
				c := Record{"name": m.doc.Rich([]string{child.HeadingID}), "procedure": m.doc.Rich(child.Intro)}
				m.bind(r, path+".name", []string{child.HeadingID}, "named_field", "assessment-case-name")
				m.bind(r, path+".procedure", child.Intro, "named_field", "assessment-case-procedure")
				for _, detail := range child.Children {
					f := map[string]string{"Conditions": "conditions", "Procedure": "procedure", "Expected result": "expected_result", "Evidence to retain": "evidence_to_retain", "Evidence retention": "evidence_to_retain"}[detail.Title]
					if f != "" {
						m.setRich(c, f, detail.Content)
						m.sectionBindings(r, path+"."+f, detail)
					} else {
						m.extension(r, detail)
					}
				}
				cases = append(cases, c)
			}
			v["cases"] = cases
			m.bind(r, field+".cases", []string{sub.HeadingID}, "structural_wrapper", "assessment-cases")
		default:
			m.extension(r, sub)
		}
	}
	return true
}

func plansFor(r Record, purpose assessment.Purpose) []Record {
	var result []Record
	for _, plan := range records(r["assessment_plans"]) {
		if String(plan, "purpose") == string(purpose) {
			result = append(result, plan)
		}
	}
	return result
}

func (d *designResolver) hasPlanGuidance(r Record, purpose assessment.Purpose) bool {
	for _, plan := range plansFor(r, purpose) {
		if strings.TrimSpace(d.text(plan["content"])) != "" {
			return true
		}
	}
	return false
}

// AssessmentHeadings locates authored plan headings for navigation in a snapshot.
// Group them by the purpose recorded on the plan, including after JSON reload.
func (s *Snapshot) AssessmentHeadings() map[string]assessment.Purpose {
	type bindingKey struct{ owner, field string }
	purposes := map[bindingKey]assessment.Purpose{}
	for _, owner := range s.Elements {
		for i, plan := range records(owner["assessment_plans"]) {
			key := bindingKey{String(owner, "id"), fmt.Sprintf("assessment_plans.%d.purpose", i)}
			purposes[key] = assessment.Purpose(String(plan, "purpose"))
		}
	}
	headings := map[string]assessment.Purpose{}
	for _, binding := range s.ContentBindings {
		if purpose, ok := purposes[bindingKey{String(binding, "owner"), String(binding, "field_path")}]; ok {
			headings[String(binding, "content_node")] = purpose
		}
	}
	return headings
}
