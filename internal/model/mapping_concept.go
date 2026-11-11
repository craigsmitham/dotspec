package model

func (m *domainMapper) concept(r Record, s *Section) bool {
	switch s.Title {
	case "Purpose":
		m.section(r, "purpose", s)
	case "Operational principle":
		m.section(r, "operational_principle", s)
	case "State":
		v := Record{"content": m.doc.Rich(s.Content)}
		r["abstract_state"] = v
		m.sectionBindings(r, "abstract_state.content", s)
		elements := []Record{}
		for _, id := range m.listItems(s.Intro, false) {
			elements = append(elements, Record{"description": m.doc.Rich([]string{id})})
		}
		for _, sub := range s.Children {
			if sub.Title == "Invariants" {
				old, _ := v["invariants"].([]Record)
				next, _ := Object(m.conditions(sub))["conditions"].([]Record)
				v["invariants"] = append(append([]Record{}, old...), next...)
				continue
			}
			x := Record{"name": m.doc.Rich([]string{sub.HeadingID}), "description": m.doc.Rich(sub.Content)}
			for _, detail := range sub.Children {
				if detail.Title == "Value" || detail.Title == "Values" || detail.Title == "Type" {
					m.setRich(x, "value_description", detail.Content)
				} else {
					m.extension(r, detail)
				}
			}
			elements = append(elements, x)
		}
		if len(elements) > 0 {
			v["elements"] = elements
		}
	case "Actions":
		v := m.collection(r, "actions", "actions", s)
		for _, sub := range s.Children {
			x := m.named("ConceptAction", r, sub)
			m.section(x, "description", sub)
			for _, detail := range sub.Children {
				field := map[string]string{"Inputs": "inputs", "Outputs": "outputs", "Effects": "effects", "Rejections": "rejections", "Rejection behavior": "rejections"}[detail.Title]
				if field != "" {
					m.section(x, field, detail)
					continue
				}
				switch detail.Title {
				case "Caller assumptions", "Assumptions":
					m.conditionField(x, "caller_assumptions", detail)
				case "Guards":
					m.conditionField(x, "guards", detail)
				default:
					if !openDecisionHeading(detail.Title) {
						m.extension(x, detail)
					}
				}
			}
			m.openDecisions(x, sub.Children)
			appendID(v, "actions", String(x, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) requirement(r Record, s *Section, kind string) bool {
	switch s.Title {
	case "Rationale":
		m.section(r, "rationale", s)
	case "Applies to":
		m.wrapper(r, "applies_to", s, map[string]string{"Conditions": "conditions", "Workload": "workload", "Population": "population", "Exclusions": "exclusions"}, "Subject", "Subjects")
		for _, sub := range s.Children {
			if sub.Title == "Subjects" || sub.Title == "Subject" {
				m.contextField(Object(r["applies_to"]), "subjects", sub.Content)
				m.sectionBindings(r, "applies_to.subjects.content", sub)
			}
		}
	case "Statement":
		if kind != "requirement" {
			return false
		}
		m.section(r, "statement", s)
	case "Verification":
		if kind != "requirement" {
			return false
		}
		m.wrapper(r, "verification", s, map[string]string{"Method": "method", "Procedure": "method", "Success criteria": "success_criteria", "Conditions": "conditions", "Evidence to retain": "evidence_to_retain", "Evidence retention": "evidence_to_retain"}, "Cases")
		for _, sub := range s.Children {
			if sub.Title == "Cases" {
				cases, _ := Object(r["verification"])["cases"].([]Record)
				if cases == nil {
					cases = []Record{}
				}
				for _, child := range sub.Children {
					v := Record{"name": m.doc.Rich([]string{child.HeadingID}), "procedure": m.doc.Rich(child.Intro)}
					for _, detail := range child.Children {
						field := map[string]string{"Conditions": "conditions", "Procedure": "procedure", "Expected result": "expected_result", "Evidence to retain": "evidence_to_retain"}[detail.Title]
						if field != "" {
							m.setRich(v, field, detail.Content)
						} else {
							m.extension(r, detail)
						}
					}
					cases = append(cases, v)
				}
				Object(r["verification"])["cases"] = cases
				m.bind(r, "verification.cases", []string{sub.HeadingID}, "structural_wrapper", "verification-cases")
			}
		}
	case "Criterion":
		if kind != "quality" {
			return false
		}
		m.wrapper(r, "criterion", s, map[string]string{"Statement": "statement", "Measure": "measure", "Comparison": "comparison", "Target": "target", "Units": "units", "Aggregation": "aggregation", "Window": "window", "Observation window": "window", "Conditions": "conditions", "Workload": "workload", "Population": "population", "Exclusions": "exclusions"}, "Open decision", "Open decisions", "Open question", "Open questions")
		m.openDecisions(r, s.Children)
	case "Measurement":
		if kind != "quality" {
			return false
		}
		m.wrapper(r, "measurement", s, map[string]string{"Sampling": "sampling", "Conditions": "conditions", "Procedure": "procedure", "Evidence to retain": "evidence_to_retain", "Evidence retention": "evidence_to_retain", "Workload": "workload", "Population": "population", "Exclusions": "exclusions"})
	default:
		return false
	}
	return true
}

func (m *domainMapper) wrapper(owner Record, field string, s *Section, details map[string]string, special ...string) {
	v := Record{"content": m.doc.Rich(s.Content)}
	owner[field] = v
	m.sectionBindings(owner, field+".content", s)
	for _, sub := range s.Children {
		if nested, ok := details[sub.Title]; ok {
			m.setRich(v, nested, sub.Content)
			m.sectionBindings(owner, field+"."+nested, sub)
		} else {
			accepted := false
			for _, title := range special {
				if sub.Title == title {
					accepted = true
					break
				}
			}
			if !accepted {
				m.extension(owner, sub)
			}
		}
	}
}
