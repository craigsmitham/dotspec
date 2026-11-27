package model

func (m *domainMapper) job(r Record, s *Section) bool {
	switch s.Title {
	case "Steps":
		m.declared(r, "steps", s)
		for _, child := range s.Children {
			x := m.named("JobStep", r, child)
			description := m.withoutFields(child, child.Intro, "ID", "Stage")
			x["description"] = m.doc.Rich(description)
			m.markField(x, "description", description)
			if id := m.fieldText(child, "ID"); id != "" {
				x["authored_id"] = id
				m.authoredFragment(x, m.fieldNodes(child, "ID"))
			}
			if nodes := m.fieldNodes(child, "Stage"); nodes != nil {
				x["stage"] = m.doc.Rich(nodes)
				m.markField(x, "stage", nodes)
			}
			for _, sub := range child.Children {
				if sub.Title == "Outcomes" {
					m.outcomes(x, sub)
				} else {
					m.subjectSection(x, sub)
				}
			}
			appendID(r, "steps", String(x, "id"))
		}
	case "Outcomes":
		m.outcomes(r, s)

	default:
		return false
	}
	return true
}

func (m *domainMapper) outcomes(owner Record, s *Section) {
	m.declared(owner, "outcomes", s)
	for _, sub := range s.Children {
		x := m.named("DesiredOutcome", owner, sub)
		if id := m.fieldText(sub, "ID"); id != "" {
			x["authored_id"] = id
			m.authoredFragment(x, m.fieldNodes(sub, "ID"))
		}
		nodes := m.withoutFields(sub, sub.Intro, "ID")
		x["statement"] = m.doc.Rich(nodes)
		m.markField(x, "statement", nodes)
		for _, detail := range sub.Children {
			m.subjectSection(x, detail)
		}
		appendID(owner, "outcomes", String(x, "id"))
	}
}
