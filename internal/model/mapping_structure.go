package model

import "strings"

// Structural links and feature membership are explicit source projections.
// Declaration containment does not supply a missing structural parent.
func (m *domainMapper) structureReference(r Record, field string, s *Section) {
	m.section(r, field+"_content", s)
	m.relationship(r, field, s.Content)
}

func (m *domainMapper) componentSection(r Record, s *Section) bool {
	if m.composition(r, s) {
		return true
	}
	switch s.Title {
	case "Purpose", "Description", "Overview", "Responsibilities":
		m.section(r, strings.ToLower(s.Title), s)
	case "Part of":
		m.structureReference(r, "part_of", s)
	case "Boundary":
		m.section(r, "boundary_description", s)
	case "Contracts", "Architecture", "Implementation":
		field := strings.ToLower(s.Title)
		m.contextField(r, field, s.Content)
		m.markSection(r, field+".content", s)
	default:
		return false
	}
	return true
}
