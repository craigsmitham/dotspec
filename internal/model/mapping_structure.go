package model

// Structural links and feature membership are explicit source projections.
// Declaration containment does not supply a missing structural parent.
func (m *domainMapper) structureReference(r Record, field string, s *Section) {
	m.section(r, field+"_content", s)
	refs := append([]string{}, Strings(r[field+"_references"])...)
	for _, id := range m.references(s.Content) {
		for _, ref := range m.doc.References {
			if String(ref, "id") == id && (String(ref, "form") == "inline_link" || String(ref, "form") == "reference_link") {
				refs = append(refs, id)
			}
		}
	}
	r[field+"_references"] = refs
	m.bindDesignReferences(r, field+"_references")
}

func (m *domainMapper) systemDeclarations(system Record, s *Section) {
	kind, field := "Component", "components"
	if s.Title == "Features" {
		kind, field = "Feature", "features"
	}
	if s.Title == "Capabilities" {
		kind, field = "Capability", "capabilities"
	}
	m.section(system, field+"_content", s)
	system[field] = []string{}
	for _, child := range s.Children {
		r := m.named(kind, system, child)
		m.section(r, "content", child)
		if id := m.fieldText(child, "ID"); id != "" {
			r["authored_id"] = id
			m.bind(r, "authored_id", m.fieldNodes(child, "ID"), "record_identity", "structure-id")
			m.bind(r, "authored_id", m.fieldLabelNodes(child, "ID"), "structural_wrapper", "structure-id-label")
		}
		if kind == "Feature" {
			r["system"] = system["id"]
		}
		for _, sub := range child.Children {
			if m.traceSection(r, sub) {
				continue
			}
			switch {
			case openQuestionHeading(sub.Title):
				m.openQuestions(r, []*Section{sub})
			case sub.Title == "Purpose":
				m.section(r, "purpose", sub)
			case kind == "Capability" && sub.Title == "Enablers":
				m.section(r, "enablers", sub)
			case kind == "Feature" && sub.Title == "Groups":
				m.structureReference(r, "groups", sub)
			case kind == "Component" && sub.Title == "Part of":
				m.structureReference(r, "part_of", sub)
			case kind == "Component" && sub.Title == "Responsibilities":
				m.section(r, "responsibilities", sub)
			case kind == "Component" && sub.Title == "Boundary":
				m.section(r, "boundary_description", sub)
			case kind == "Component" && (sub.Title == "Contracts" || sub.Title == "Architecture" || sub.Title == "Implementation"):
				field := map[string]string{"Contracts": "contracts", "Architecture": "architecture", "Implementation": "implementation"}[sub.Title]
				m.contextField(r, field, sub.Content)
				m.sectionBindings(r, field+".content", sub)
			default:
				m.extension(r, sub)
			}
		}
		appendID(system, field, String(r, "id"))
	}
}
