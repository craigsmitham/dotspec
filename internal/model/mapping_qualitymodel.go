package model

import "strings"

func (m *domainMapper) qualityModel(r Record, s *Section) bool {
	switch s.Title {
	case "Sources":
		m.contextField(r, "sources", s.Content)
		m.markSection(r, "sources.content", s)
	case "Quality attributes":
		m.declared(r, "quality_attributes", s)
		for _, child := range s.Children {
			a := m.named("QualityModelAttribute", r, child)
			a["authored_id"] = m.fieldText(child, "ID")
			definition := m.withoutFields(child, child.Content, "ID", "Sub-attribute of")
			a["definition"] = m.doc.Rich(definition)
			m.markField(a, "definition", definition)
			m.authoredFragment(a, m.fieldNodes(child, "ID"))
			if ids := m.fieldNodes(child, "Sub-attribute of"); len(ids) > 0 {
				for _, label := range strings.Split(m.fieldText(child, "Sub-attribute of"), ",") {
					label = strings.TrimSpace(label)
					if label != "" {
						m.label(a, label, ids[0], "subattribute_of")
					}
				}
			}
			m.openQuestions(a, child.Children)
			appendID(r, "quality_attributes", String(a, "id"))
		}
	default:
		return false
	}
	return true
}
