package model

func (m *domainMapper) qualityModel(r Record, s *Section) bool {
	switch s.Title {
	case "Description":
		m.section(r, "description", s)
	case "Edition":
		m.section(r, "edition", s)
	case "Sources":
		m.contextField(r, "sources", s.Content)
		m.sectionBindings(r, "sources.content", s)
	case "Quality attributes":
		return m.system(r, s)
	default:
		return false
	}
	return true
}
