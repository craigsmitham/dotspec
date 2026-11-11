package model

func (m *domainMapper) businessRule(r Record, s *Section) bool {
	switch s.Title {
	case "Statement":
		m.section(r, "statement", s)
	case "Exceptions":
		m.section(r, "exceptions", s)
	case "Applies to", "Rationale":
		return m.requirement(r, s, "business-rule")
	default:
		return false
	}
	return true
}
