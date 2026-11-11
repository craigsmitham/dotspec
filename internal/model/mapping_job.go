package model

import "strings"

func (m *domainMapper) job(r Record, s *Section) bool {
	fields := map[string]string{"Job": "statement", "Performer": "performer", "Situation": "situation", "Motivation": "motivation", "Intended progress": "intended_progress", "Desired outcome": "intended_progress"}
	if field, ok := fields[s.Title]; ok {
		m.section(r, field, s)
		return true
	}
	switch s.Title {
	case "Steps":
		v := m.collection(r, "steps", "steps", s)
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
			appendID(v, "steps", String(x, "id"))
		}
	case "Outcomes":
		m.outcomes(r, s)
	case "Progress cases":
		v := m.collection(r, "progress_cases", "cases", s)
		for _, child := range s.Children {
			x := m.named("ProgressCase", r, child)
			m.section(x, "description", child)
			for _, sub := range child.Children {
				switch sub.Title {
				case "Current approach":
					m.section(x, "current_approach", sub)
				case "Desired progress":
					m.section(x, "desired_progress", sub)
				case "Forces":
					m.forces(x, sub)
				case "Timeline":
					m.timeline(x, sub)
				case "Evidence":
					m.evidence(x, "evidence", sub)
				default:
					m.extension(x, sub)
				}
			}
			appendID(v, "cases", String(x, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) outcomes(owner Record, s *Section) {
	v := m.collection(owner, "outcomes", "outcomes", s)
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
		appendID(v, "outcomes", String(x, "id"))
	}
	for _, sub := range s.Children {
		m.extension(owner, sub)
	}
}

func (m *domainMapper) forces(owner Record, s *Section) {
	items := []Record{}
	for _, id := range m.listItems(s.Intro, false) {
		text := strings.TrimSpace(m.doc.Text(id))
		text = strings.TrimLeft(text, "-*+0123456789.) ")
		text = strings.TrimPrefix(text, "**")
		label := strings.ToLower(strings.TrimSpace(strings.SplitN(text, ":", 2)[0]))
		item := Record{"description": m.doc.Rich([]string{id})}
		switch label {
		case "push", "pull", "anxiety", "habit":
			item["kind"] = label
		}
		items = append(items, item)
	}
	for _, sub := range s.Children {
		item := Record{"description": m.doc.Rich(sub.Content)}
		switch strings.ToLower(sub.Title) {
		case "push", "pull", "anxiety", "habit":
			item["kind"] = strings.ToLower(sub.Title)
		}
		items = append(items, item)
	}
	v := Object(owner["forces"])
	if v == nil {
		v = Record{}
	}
	previous, _ := v["items"].([]Record)
	v["items"] = append(append([]Record{}, previous...), items...)
	m.setRich(v, "content", s.Content)
	owner["forces"] = v
	m.sectionBindings(owner, "forces.content", s)
}

func (m *domainMapper) timeline(owner Record, s *Section) {
	events := []Record{}
	for _, id := range m.listItems(s.Intro, false) {
		events = append(events, Record{"description": m.doc.Rich([]string{id})})
	}
	for _, child := range s.Children {
		v := Record{"label": m.doc.Rich([]string{child.HeadingID}), "description": m.doc.Rich(child.Content)}
		for _, sub := range child.Children {
			if sub.Title == "Time" || sub.Title == "Date" {
				m.setRich(v, "time", sub.Content)
			} else {
				m.extension(owner, sub)
			}
		}
		events = append(events, v)
	}
	v := Object(owner["timeline"])
	if v == nil {
		v = Record{}
	}
	previous, _ := v["events"].([]Record)
	v["events"] = append(append([]Record{}, previous...), events...)
	m.setRich(v, "content", s.Content)
	owner["timeline"] = v
	m.sectionBindings(owner, "timeline.content", s)
}

func (m *domainMapper) evidence(owner Record, field string, s *Section) {
	v := m.collection(owner, field, "statements", s)
	for _, id := range m.listItems(s.Intro, false) {
		x := m.record("EvidenceStatement", String(owner, "id"))
		x["statement"] = m.doc.Rich([]string{id})
		x["references"] = m.references([]string{id})
		m.bind(x, "statement", []string{id}, "structured_item", "evidence-list-item")
		appendID(v, "statements", String(x, "id"))
	}
	for _, sub := range s.Children {
		switch sub.Title {
		case "Research basis":
			m.setRich(v, "research_basis", sub.Content)
			m.sectionBindings(owner, field+".research_basis", sub)
		case "Limitations":
			m.setRich(v, "limitations", sub.Content)
			m.sectionBindings(owner, field+".limitations", sub)
		default:
			x := m.named("EvidenceStatement", owner, sub)
			m.section(x, "statement", sub)
			x["references"] = m.references(sub.Content)
			for _, detail := range sub.Children {
				switch detail.Title {
				case "Type":
					m.section(x, "evidence_type", detail)
				case "Research status":
					m.section(x, "research_status", detail)
				case "Limitations":
					m.section(x, "limitations", detail)
				default:
					m.extension(x, detail)
				}
			}
			appendID(v, "statements", String(x, "id"))
		}
	}
}
