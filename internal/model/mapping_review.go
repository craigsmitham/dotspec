package model

import (
	"strings"
)

func (m *domainMapper) boundarySection(r Record, s *Section) bool {
	if !designKind(String(r, "kind"), "System or Component") {
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
