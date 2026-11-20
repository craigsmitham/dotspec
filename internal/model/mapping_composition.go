package model

// Composition is explicitly authored in system sections. It is independent of
// directory-derived system association and does not change reference traversal.
func (m *domainMapper) composition(r Record, s *Section) bool {
	switch s.Title {
	case "Concept uses":
		m.section(r, "concept_uses_content", s)
		r["concept_uses"] = []string{}
		for _, child := range s.Children {
			x := m.named("ConceptUse", r, child)
			m.section(x, "content", child)
			for _, sub := range child.Children {
				switch sub.Title {
				case "Description":
					m.section(x, "description", sub)
				case "Definition":
					x["definition_references"] = m.references(sub.Content)
					m.bindDesignReferences(x, "definition_references")
					m.bind(x, "definition_references", []string{sub.HeadingID}, "structural_wrapper", "concept-definition-heading")
				case "Configuration":
					m.section(x, "configuration", sub)
				case "Dependencies":
					m.section(x, "dependencies_content", sub)
					x["dependencies_references"] = m.references(sub.Content)
				default:
					m.extension(x, sub)
				}
			}
			appendID(r, "concept_uses", String(x, "id"))
		}
	case "Synchronizations":
		m.section(r, "synchronizations_content", s)
		r["synchronizations"] = []string{}
		for _, child := range s.Children {
			x := m.named("Synchronization", r, child)
			m.section(x, "content", child)
			for _, sub := range child.Children {
				field := map[string]string{"When": "when", "Where": "where", "Then": "then", "Semantics": "semantics"}[sub.Title]
				if field == "" {
					m.extension(x, sub)
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
		r[field+"_references"] = m.references(nodes)
		m.bindDesignReferences(r, field+"_references")
		m.bind(r, field+"_references", m.fieldLabelNodes(s, label), "structural_wrapper", "design-reference-label/"+label)
	}
}

func (m *domainMapper) bindDesignReferences(r Record, field string) {
	m.bindReferenceOccurrences(r, field, Strings(r[field]))
}

func (m *domainMapper) bindReferenceOccurrences(owner Record, field string, ids []string) {
	for _, id := range ids {
		for _, ref := range m.doc.References {
			if String(ref, "id") == id {
				m.bind(owner, field, []string{String(ref, "content_node")}, "named_field", "reference-occurrence")
			}
		}
	}
}
