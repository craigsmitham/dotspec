package model

func (m *domainMapper) concept(r Record, s *Section) bool {
	switch s.Title {
	case "Effects table":
		m.section(r, "effects_table", s)
	case "Purpose":
		m.section(r, "purpose", s)
	case "Operational principles":
		m.section(r, "operational_principles", s)
	case "State":
		m.section(r, "abstract_state", s)
	case "Actions":
		m.declared(r, "actions", s)
		for _, sub := range s.Children {
			x := m.named("ConceptAction", r, sub)
			m.section(x, "description", sub)
			for _, detail := range sub.Children {
				field := map[string]string{"Inputs": "inputs", "Outputs": "outputs", "Initiation": "initiation", "Effects": "effects", "Rejections": "rejections"}[detail.Title]
				if field != "" {
					m.section(x, field, detail)
					continue
				}
				switch detail.Title {
				case "Caller assumptions":
					m.section(x, "caller_assumptions", detail)
				case "Guards":
					m.section(x, "guards", detail)
				default:
					if !openQuestionHeading(detail.Title) {
						m.subjectSection(x, detail)
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
	case "Verification procedures":
		m.procedures(r, r, "verification_procedures", s)
	case "Criterion":
		m.wrapper(r, "criterion", s, map[string]string{"Measure": "measure", "Target": "target"})
		for _, sub := range s.Children {
			if sub.Title == "Verification procedures" {
				m.procedures(r, Object(r["criterion"]), "criterion.verification_procedures", sub)
			}
		}
		m.openQuestions(r, s.Children)
	default:
		return false
	}
	return true
}

func (m *domainMapper) wrapper(owner Record, field string, s *Section, details map[string]string) {
	v := Record{"content": m.doc.Rich(s.Content)}
	owner[field] = v
	m.markSection(owner, field+".content", s)
	for _, sub := range s.Children {
		if nested, ok := details[sub.Title]; ok {
			m.setRich(v, nested, sub.Content)
			m.markSection(owner, field+"."+nested, sub)
		}
	}
}
