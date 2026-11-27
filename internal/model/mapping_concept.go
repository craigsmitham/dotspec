package model

func (m *domainMapper) concept(r Record, s *Section) bool {
	switch s.Title {
	case "Actions":
		m.declared(r, "actions", s)
		for _, sub := range s.Children {
			x := m.named("ConceptAction", r, sub)
			m.section(x, "description", sub)
			for _, detail := range sub.Children {
				if !openQuestionHeading(detail.Title) {
					m.subjectSection(x, detail)
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
	case "Applies to":
		m.designField(r, "applies_to", s, "Targets")
		m.wrapper(r, "applicability", s, map[string]string{"Conditions": "conditions", "Workload": "workload", "Population": "population", "Exclusions": "exclusions"})
	case "Criterion":
		m.wrapper(r, "criterion", s, map[string]string{"Measure": "measure", "Target": "target"})
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
