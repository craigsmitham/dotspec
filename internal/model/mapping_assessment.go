package model

import "fmt"

// Assessment structure is shared; field, subject and basis retain the purpose.
func (m *domainMapper) assessmentSection(r Record, s *Section) bool {
	field := map[string]string{"Verification plan": "verification_plan", "Validation plan": "validation_plan"}[s.Title]
	if field == "" {
		return false
	}
	if _, ok := compiledModel().catalog.Classes[String(r, "kind")].Slots[field]; !ok {
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
	v := Object(r[field])
	if v == nil {
		v = Record{}
		r[field] = v
	}
	m.setRich(v, "content", s.Content)
	m.sectionBindings(r, field+".content", s)
	details := map[string]string{"Subject": "subject", "Participants": "participants", "Method": "method", "Procedure": "method", "Success criteria": "success_criteria", "Conditions": "conditions", "Evidence to retain": "evidence_to_retain", "Evidence retention": "evidence_to_retain"}
	for _, sub := range s.Children {
		if f := details[sub.Title]; f != "" {
			m.setRich(v, f, sub.Content)
			m.sectionBindings(r, field+"."+f, sub)
			continue
		}
		switch sub.Title {
		case "Basis":
			m.contextField(v, "basis", sub.Content)
			m.sectionBindings(r, field+".basis.content", sub)
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
