package model

func (m *domainMapper) job(r Record, s *Section) bool {
	fields := map[string]string{"Job": "statement", "Performer": "performer", "Situation": "situation", "Motivation": "motivation", "Intended progress": "intended_progress"}
	if field, ok := fields[s.Title]; ok {
		m.section(r, field, s)
		return true
	}
	switch s.Title {
	case "Steps":
		m.declared(r, "steps", s)
		for _, child := range s.Children {
			x := m.named("JobStep", r, child)
			description := m.withoutFields(child, child.Intro, "ID", "Stage")
			x["description"] = m.doc.Rich(description)
			m.bind(x, "description", description, "named_field", "step-description")
			if id := m.fieldText(child, "ID"); id != "" {
				x["authored_id"] = id
				m.bind(x, "authored_id", m.fieldNodes(child, "ID"), "record_identity", "step-id")
			}
			if nodes := m.fieldNodes(child, "Stage"); nodes != nil {
				x["stage"] = m.doc.Rich(nodes)
				m.bind(x, "stage", nodes, "named_field", "step-stage")
				m.bind(x, "stage", m.fieldLabelNodes(child, "Stage"), "structural_wrapper", "step-stage-label")
			}
			for _, sub := range child.Children {
				if sub.Title == "Outcomes" {
					m.outcomes(x, sub)
				} else {
					m.extension(x, sub)
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
	for _, item := range m.listItems(s.Intro, false) {
		x := m.record("DesiredOutcome", String(owner, "id"))
		x["statement"] = m.doc.Rich([]string{item})
		m.bind(x, "statement", []string{item}, "structured_item", "outcome-list-item")
		if match := outcomeID.FindStringSubmatch(m.doc.Text(item)); len(match) > 1 {
			x["authored_id"] = match[1]
			m.bind(x, "authored_id", []string{item}, "record_identity", "outcome-id")
		}
		children := Strings(m.doc.NodeByID[item]["children"])
		if len(children) > 1 {
			x["supporting_content"] = m.doc.Rich(children[1:])
			m.bind(x, "supporting_content", children[1:], "named_field", "outcome-supporting-content")
		}
		appendID(owner, "outcomes", String(x, "id"))
	}
	for _, sub := range s.Children {
		if m.fieldText(sub, "ID") == "" {
			m.extension(owner, sub)
			continue
		}
		x := m.named("DesiredOutcome", owner, sub)
		if id := m.fieldText(sub, "ID"); id != "" {
			x["authored_id"] = id
			m.bind(x, "authored_id", m.fieldNodes(sub, "ID"), "record_identity", "outcome-id")
			m.bind(x, "authored_id", m.fieldLabelNodes(sub, "ID"), "structural_wrapper", "outcome-id-label")
		}
		nodes := m.withoutFields(sub, sub.Intro, "ID")
		x["statement"] = m.doc.Rich(nodes)
		m.bind(x, "statement", nodes, "named_field", "outcome-statement")
		for _, detail := range sub.Children {
			if m.traceSection(x, detail) {
				continue
			}
			if detail.Title == "Supporting content" {
				m.section(x, "supporting_content", detail)
			} else {
				m.extension(x, detail)
			}
		}
		appendID(owner, "outcomes", String(x, "id"))
	}
}
