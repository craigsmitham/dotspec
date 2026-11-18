package model

func (m *domainMapper) concept(r Record, s *Section) bool {
	switch s.Title {
	case "Effects table":
		m.effectsTable(r, s)
	case "Purpose":
		m.section(r, "purpose", s)
	case "Type parameters":
		m.conceptTypes(r, s, "type_parameters")
	case "Introduced types":
		m.conceptTypes(r, s, "introduced_types")
	case "Operational principles":
		m.section(r, "operational_principles", s)
	case "State":
		v := Record{"content": m.doc.Rich(s.Content)}
		r["abstract_state"] = v
		m.sectionBindings(r, "abstract_state.content", s)
		elements := []Record{}
		for _, id := range m.listItems(s.Intro, false) {
			elements = append(elements, Record{"description": m.doc.Rich([]string{id})})
		}
		for _, sub := range s.Children {
			if sub.Title == "Initial state" {
				m.setRich(v, "initial_state", sub.Content)
				continue
			}
			if sub.Title == "Invariants" {
				old, _ := v["invariants"].([]Record)
				v["invariants"] = append(append([]Record{}, old...), m.conditions(sub)...)
				continue
			}
			x := Record{"name": m.doc.Rich([]string{sub.HeadingID}), "description": m.doc.Rich(sub.Content)}
			for _, detail := range sub.Children {
				switch detail.Title {
				case "Value", "Values", "Type":
					m.setRich(x, "value_description", detail.Content)
				case "Derivation":
					m.setRich(x, "derivation", detail.Content)
				default:
					m.extension(r, detail)
				}
			}
			elements = append(elements, x)
		}
		if len(elements) > 0 {
			v["elements"] = elements
		}
	case "Actions":
		m.declared(r, "actions", s)
		for _, sub := range s.Children {
			x := m.named("ConceptAction", r, sub)
			m.section(x, "description", sub)
			for _, detail := range sub.Children {
				field := map[string]string{"Inputs": "inputs", "Outputs": "outputs", "Initiation": "initiation", "Effects": "effects", "Rejections": "rejections", "Rejection behavior": "rejections"}[detail.Title]
				if field != "" {
					m.section(x, field, detail)
					if field == "inputs" {
						m.actionParameters(x, detail, "input_parameters")
					}
					if field == "outputs" {
						m.actionParameters(x, detail, "output_parameters")
					}
					continue
				}
				switch detail.Title {
				case "Caller assumptions", "Assumptions":
					m.conditionField(x, "caller_assumptions", detail)
				case "Guards":
					m.conditionField(x, "guards", detail)
				default:
					if !openQuestionHeading(detail.Title) {
						m.extension(x, detail)
					}
				}
			}
			m.openQuestions(x, sub.Children)
			appendID(r, "actions", String(x, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) requirement(r Record, s *Section) bool {
	switch s.Title {
	case "Rationale":
		m.section(r, "rationale", s)
	case "Applies to":
		m.designField(r, "applies_to", s, "Targets")
		m.wrapper(r, "applicability", s, map[string]string{"Conditions": "conditions", "Workload": "workload", "Population": "population", "Exclusions": "exclusions"})
	case "Statement":
		m.section(r, "statement", s)
	case "Verification plan", "Validation plan":
		return m.assessmentSection(r, s)
	case "Criterion":
		m.wrapper(r, "criterion", s, map[string]string{"Measure": "measure", "Target": "target"}, "Open question", "Open questions")
		m.openQuestions(r, s.Children)
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
