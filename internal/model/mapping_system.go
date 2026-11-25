package model

import (
	"fmt"
	"strings"
)

func (m *domainMapper) system(r Record, s *Section) bool {
	if m.composition(r, s) {
		return true
	}
	switch s.Title {
	case "Part of":
		m.structureReference(r, "part_of", s)
	case "Framing", "Description":
		m.section(r, strings.ToLower(s.Title), s)
	case "Purpose":
		m.section(r, "purpose", s)
	case "Project references":
		m.contextField(r, "project_references", s.Content)
		m.markSection(r, "project_references.content", s)
	case "Scope":
		m.section(r, "scope", s)
	case "Environment":
		m.section(r, "environment", s)
	case "Stakeholders":
		m.section(r, "stakeholders", s)
	case "Constraints":
		// Constraints stay narrative; an addressable constraint is a requirement or business rule.
		m.section(r, "constraints", s)
	case "Goals", "Assumptions":
		field, kind := "goals", "Goal"
		if s.Title == "Assumptions" {
			field, kind = "assumptions", "Assumption"
		}
		m.declared(r, field, s)
		for _, child := range s.Children {
			x := m.named(kind, r, child)
			if kind == "Assumption" {
				x["statement"] = m.doc.Rich(child.Intro)
				m.markField(x, "statement", child.Intro)
				m.assumptionDetails(x, child)
			} else {
				m.section(x, "statement", child)
			}
			appendID(r, field, String(x, "id"))
		}
	case "Glossary":
		m.declared(r, "glossary", s)
		m.glossaryTerms(r, "glossary", s.Children)
	case "Quality attributes":
		m.declared(r, "quality_attributes", s)
		if String(r, "kind") == "System" {
			break
		}
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

// openQuestions only recognizes explicitly headed questions. A sentence about
// missing evidence or a target is not automatically promoted to a question.
func (m *domainMapper) openQuestions(owner Record, sections []*Section) {
	for _, s := range sections {
		if !openQuestionHeading(s.Title) {
			continue
		}
		m.markHeading("open_questions", s.HeadingID)
		values, _ := owner["open_questions"].([]Record)
		if values == nil {
			values = []Record{}
			owner["open_questions"] = values
		}
		add := func(question []string, item *Section) {
			v := Record{"question": m.doc.Rich(question)}
			index := len(values)
			values = append(values, v)
			owner["open_questions"] = values
			path := fmt.Sprintf("open_questions.%d", index)
			m.markField(owner, path+".question", question)
			if item == nil {
				return
			}
			m.markHeading(path+".question", item.HeadingID)
		}
		if len(s.Children) > 0 {
			for _, child := range s.Children {
				add(append([]string{child.HeadingID}, child.Content...), child)
			}
		} else if items := m.listItems(s.Intro, false); len(items) > 0 {
			for _, item := range items {
				add([]string{item}, nil)
			}
		} else {
			add(s.Content, s)
		}
	}
}

func openQuestionHeading(title string) bool {
	return title == "Open questions"
}

func (m *domainMapper) glossaryTerms(r Record, field string, terms []*Section) {
	for _, child := range terms {
		x := m.named("GlossaryTerm", r, child)
		m.section(x, "definition", child)
		for _, sub := range child.Children {
			switch sub.Title {
			case "Alternative names":
				names, _ := x["alternative_names"].([]Record)
				items := m.listItems(sub.Intro, false)
				for _, item := range items {
					names = append(names, m.doc.Rich([]string{item}))
					x["alternative_names"] = names
					m.markField(x, fmt.Sprintf("alternative_names.%d", len(names)-1), []string{item})
				}
				if len(items) == 0 {
					names = append(names, m.doc.Rich(sub.Intro))
					x["alternative_names"] = names
					m.markField(x, fmt.Sprintf("alternative_names.%d", len(names)-1), sub.Intro)
				}
				x["alternative_names"] = names
				m.markHeading("alternative_names", sub.HeadingID)
				for _, detail := range sub.Children {
					m.subjectSection(x, detail)
				}
			case "Topics":
				m.glossaryTopics(x, sub)
			case "Relationships", "Related terms":
				detail := "relationships"
				if sub.Title == "Related terms" {
					detail = "related_terms"
				}
				m.contextField(x, detail, sub.Content)
				m.markSection(x, detail+".content", sub)
				m.bindReferenceOccurrences(x, detail+".references", m.references(sub.Content))
			}
		}
		appendID(r, field, String(x, "id"))
	}
}

func (m *domainMapper) glossaryTopics(term Record, section *Section) {
	for _, item := range m.listItems(section.Intro, false) {
		for _, block := range Strings(m.doc.NodeByID[item]["children"]) {
			if String(m.doc.NodeByID[block], "node_kind") != "paragraph" {
				continue
			}
			topic := strings.Join(strings.Fields(m.doc.Text(block)), " ")
			if topic != "" {
				index := len(Strings(term["topics"]))
				term["topics"] = append(Strings(term["topics"]), topic)
				m.markField(term, fmt.Sprintf("topics.%d", index), []string{block})
			}
			break
		}
	}
}

// SystemGlossary returns the terms in the inline glossary and each declared
// source. A shared vocabulary retains one term identity across consuming systems.
func SystemGlossary(system Record, index map[string]Record) []string {
	terms := append([]string(nil), Strings(system["glossary"])...)
	for _, doc := range Strings(system["glossary_sources"]) {
		terms = append(terms, Strings(index[doc]["glossary_terms"])...)
	}
	return terms
}
