package model

import (
	"strings"
)

func (m *domainMapper) boundarySection(r Record, s *Section) bool {
	if !designKind(String(r, "kind"), "System", "Component") {
		return false
	}
	switch s.Title {
	case "Out of scope":
		m.section(r, "out_of_scope", s)
	case "Required interfaces", "Provided interfaces":
		field := "required_interfaces"
		if s.Title == "Provided interfaces" {
			field = "provided_interfaces"
		}
		m.declared(r, field, s)
		for _, child := range s.Children {
			item := m.named("Interface", r, child)
			m.section(item, "content", child)
			m.setRich(item, "description", child.Intro)
			m.markField(item, "description", child.Intro)
			for _, detail := range child.Children {
				switch detail.Title {
				case "Adopts":
					m.adoptions(item, detail)
				case "Contracts", "Counterpart":
					f := strings.ToLower(detail.Title)
					m.contextField(item, f, detail.Content)
					m.markSection(item, f+".content", detail)
				default:
					m.subjectSection(item, detail)
				}
			}
			appendID(r, field, String(item, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) assumptionDetails(r Record, s *Section) {
	for _, sub := range s.Children {
		if m.traceSection(r, sub) {
			continue
		}
		switch sub.Title {
		case "Basis", "Applicability":
			m.section(r, strings.ToLower(sub.Title), sub)
		case "Open questions":
			m.openQuestions(r, []*Section{sub})
		default:
			m.subjectSection(r, sub)
		}
	}
}

// Only an explicitly labelled nested bullet supplies a semantic illustration.
// Ordinary links in step prose remain contextual references.
func (m *domainMapper) stepIllustrates(r Record, item string) {
	var visit func(string)
	visit = func(id string) {
		n := m.doc.NodeByID[id]
		if id != item && String(n, "node_kind") == "list_item" {
			text := strings.TrimSpace(m.doc.Text(id))
			text = strings.TrimLeft(text, "-*+ \t")
			if strings.HasPrefix(text, "Illustrates:") || strings.HasPrefix(text, "Illustrates**:") {
				m.setRich(r, "illustrates_content", []string{id})
				m.relationship(r, "illustrates", []string{id})
				m.markField(r, "illustrates_content", []string{id})
				return
			}
		}
		for _, child := range Strings(n["children"]) {
			visit(child)
		}
	}
	visit(item)
}

func (m *domainMapper) adoptions(owner Record, section *Section) {
	m.declared(owner, "adopts", section)
	for _, entry := range section.Children {
		adoption := m.named("AdoptedProvision", owner, entry)
		refs := m.linkOccurrences(entry.Intro)
		if len(refs) == 1 {
			adoption["provision"] = refs[0]
			m.bindReferenceOccurrences(adoption, "provision", refs)
		}
		for _, detail := range entry.Children {
			m.section(adoption, strings.ToLower(detail.Title), detail)
		}
		appendID(owner, "adopts", String(adoption, "id"))
	}
}

// Only links identify procedures and provisions. Images and link definitions
// can occur in the same Markdown passage without designating a target.
func (m *domainMapper) linkOccurrences(nodes []string) []string {
	inside := m.descendants(nodes)
	var ids []string
	for _, ref := range m.doc.References {
		if inside[String(ref, "content_node")] && (ref["form"] == "inline_link" || ref["form"] == "reference_link") {
			ids = append(ids, String(ref, "id"))
		}
	}
	return ids
}

func (m *domainMapper) procedures(owner, value Record, path string, section *Section) {
	ids := m.linkOccurrences(section.Content)
	procedures := Object(value["verification_procedures"])
	if procedures == nil {
		procedures = Record{}
		value["verification_procedures"] = procedures
	}
	m.setRich(procedures, "content", section.Content)
	procedures["references"] = append(Strings(procedures["references"]), ids...)
	m.markSection(owner, path+".content", section)
	m.bindReferenceOccurrences(owner, path+".references", ids)
}
