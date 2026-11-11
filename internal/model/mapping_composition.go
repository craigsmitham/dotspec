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
				case "Definition":
					x["definition_references"] = m.references(sub.Content)
					m.bindDesignReferences(x, "definition_references")
					m.bind(x, "definition_references", []string{sub.HeadingID}, "structural_wrapper", "concept-definition-heading")
				case "Configuration":
					m.section(x, "configuration", sub)
				case "Dependencies":
					m.section(x, "dependencies_content", sub)
					x["dependencies_references"] = m.references(sub.Content)
				case "Bindings":
					m.section(x, "bindings_content", sub)
					x["bindings"] = []string{}
					for _, binding := range sub.Children {
						v := m.named("TypeBinding", x, binding)
						m.section(v, "content", binding)
						m.designField(v, "parameter", binding, "Parameter")
						m.designField(v, "type", binding, "Type")
						m.designField(v, "type_use", binding, "Use")
						m.primitive(v, binding)
						for _, detail := range binding.Children {
							m.extension(v, detail)
						}
						appendID(x, "bindings", String(v, "id"))
					}
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
				for _, item := range sub.Children {
					if field != "when" && field != "then" {
						m.extension(x, item)
						continue
					}
					// Only explicit Use/Action fields distinguish an action reference.
					// Other explanatory headings remain extension sections.
					if m.fieldNodes(item, "Use") == nil && m.fieldNodes(item, "Action") == nil {
						m.extension(x, item)
						continue
					}
					v := m.named("SynchronizationAction", x, item)
					m.section(v, "content", item)
					m.designField(v, "concept_use", item, "Use")
					m.designField(v, "action", item, "Action")
					for _, detail := range item.Children {
						slot := map[string]string{"Arguments": "arguments", "Result": "result"}[detail.Title]
						if slot != "" {
							m.section(v, slot, detail)
						} else {
							m.extension(v, detail)
						}
					}
					appendID(x, field+"_actions", String(v, "id"))
				}
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

func (m *domainMapper) primitive(r Record, s *Section) {
	if nodes := m.fieldNodes(s, "Primitive"); nodes != nil {
		r["primitive"] = m.fieldText(s, "Primitive")
		m.bind(r, "primitive", nodes, "named_field", "primitive-type")
		m.bind(r, "primitive", m.fieldLabelNodes(s, "Primitive"), "structural_wrapper", "primitive-type-label")
	}
}

func (m *domainMapper) conceptTypes(r Record, s *Section, field string) {
	m.section(r, field+"_content", s)
	r[field] = []string{}
	for _, child := range s.Children {
		x := m.named("ConceptType", r, child)
		m.section(x, "content", child)
		for _, detail := range child.Children {
			m.extension(x, detail)
		}
		appendID(r, field, String(x, "id"))
	}
}

func (m *domainMapper) actionParameters(r Record, s *Section, field string) {
	for _, child := range s.Children {
		x := m.named("ActionParameter", r, child)
		m.section(x, "content", child)
		m.designField(x, "type", child, "Type")
		m.primitive(x, child)
		for _, detail := range child.Children {
			if detail.Title == "Type" {
				m.section(x, "type_description", detail)
				// A compound type description is prose; only the explicit Type field resolves a type.
			} else {
				m.extension(x, detail)
			}
		}
		appendID(r, field, String(x, "id"))
	}
}

func (m *domainMapper) bindDesignReferences(r Record, field string) {
	for _, id := range Strings(r[field]) {
		for _, ref := range m.doc.References {
			if String(ref, "id") == id {
				m.bind(r, field, []string{String(ref, "content_node")}, "named_field", "design-reference")
			}
		}
	}
}
