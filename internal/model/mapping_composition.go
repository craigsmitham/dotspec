package model

// Composition is explicitly authored in system sections. It is independent of
// directory-derived system association and does not change reference traversal.
func (m *domainMapper) composition(r Record, s *Section) bool {
	switch s.Title {
	case "Concept uses":
		m.declared(r, "concept_uses", s)
		for _, child := range s.Children {
			x := m.named("ConceptUse", r, child)
			m.section(x, "content", child)
			for _, sub := range child.Children {
				switch sub.Title {
				case "Description":
					m.section(x, "description", sub)
				case "Definition":
					m.relationship(x, "definition", sub.Content)
					m.markHeading("definition", sub.HeadingID)
				case "Configuration":
					m.section(x, "configuration", sub)
				case "Dependencies":
					m.section(x, "dependencies_content", sub)
					m.relationship(x, "dependencies", sub.Content)
				default:
					m.subjectSection(x, sub)
				}
			}
			appendID(r, "concept_uses", String(x, "id"))
		}
	case "Synchronizations":
		m.declared(r, "synchronizations", s)
		for _, child := range s.Children {
			x := m.named("Synchronization", r, child)
			m.section(x, "content", child)
			for _, sub := range child.Children {
				field := map[string]string{"When": "when", "Where": "where", "Then": "then", "Semantics": "semantics"}[sub.Title]
				if field == "" {
					m.subjectSection(x, sub)
					continue
				}
				m.section(x, field, sub)
			}
			appendID(r, "synchronizations", String(x, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) designField(r Record, field string, s *Section, label string) {
	if nodes := m.fieldNodes(s, label); nodes != nil {
		m.relationship(r, field, nodes)
	}
}

func (m *domainMapper) relationship(owner Record, field string, nodes []string) {
	selected := m.descendants(nodes)
	for _, ref := range m.doc.References {
		form := String(ref, "form")
		if selected[String(ref, "content_node")] && (form == "inline_link" || form == "reference_link") {
			ref["role"], ref["context_owner"] = field, owner["id"]
		}
	}
}

func (m *domainMapper) bindReferenceOccurrences(owner Record, field string, ids []string) {
	for _, id := range ids {
		for _, ref := range m.doc.References {
			if String(ref, "id") == id {
				node := m.doc.NodeByID[String(ref, "content_node")]
				node["link_owner"], node["link_field"] = owner["id"], field
				ref["context_owner"], ref["context_field"] = owner["id"], field
			}
		}
	}
}
