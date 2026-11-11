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
	case "Purpose":
		m.section(r, "purpose", s)
	case "Scope":
		m.section(r, "scope", s)
	case "Environment":
		r["environment"] = Record{"content": m.doc.Rich(s.Content)}
		m.sectionBindings(r, "environment.content", s)
	case "Goals", "Assumptions", "Constraints":
		field, kind := "goals", "SystemGoal"
		if s.Title == "Assumptions" {
			field, kind = "assumptions", "Assumption"
		}
		if s.Title == "Constraints" {
			field, kind = "constraints", "Constraint"
		}
		v := m.collection(r, field, "items", s)
		for _, item := range m.listItems(s.Intro, false) {
			x := m.record(kind, String(r, "id"))
			x["statement"] = m.doc.Rich([]string{item})
			m.bind(x, "statement", []string{item}, "structured_item", "system-list-item")
			appendID(v, "items", String(x, "id"))
		}
		for _, child := range s.Children {
			x := m.named(kind, r, child)
			m.section(x, "statement", child)
			appendID(v, "items", String(x, "id"))
		}
	case "Glossary":
		v := m.collection(r, "glossary", "terms", s)
		for _, child := range s.Children {
			x := m.named("GlossaryTerm", r, child)
			m.section(x, "definition", child)
			if refs := m.references(child.Content); len(refs) > 0 {
				x["references"] = refs
			}
			for _, sub := range child.Children {
				switch sub.Title {
				case "Alternative names", "Aliases":
					names, _ := x["alternative_names"].([]Record)
					items := m.listItems(sub.Intro, false)
					for _, item := range items {
						names = append(names, m.doc.Rich([]string{item}))
						x["alternative_names"] = names
						m.bind(x, fmt.Sprintf("alternative_names.%d", len(names)-1), []string{item}, "named_field", "glossary-alias")
					}
					if len(items) == 0 {
						names = append(names, m.doc.Rich(sub.Intro))
						x["alternative_names"] = names
						m.bind(x, fmt.Sprintf("alternative_names.%d", len(names)-1), sub.Intro, "named_field", "glossary-alias")
					}
					x["alternative_names"] = names
					m.bind(x, "alternative_names", []string{sub.HeadingID}, "structural_wrapper", "glossary-aliases-heading")
					for _, detail := range sub.Children {
						m.extension(x, detail)
					}
				case "Related terms":
					m.contextField(x, "related_terms", sub.Content)
					m.sectionBindings(x, "related_terms.content", sub)
				}
			}
			appendID(v, "terms", String(x, "id"))
		}
	case "Quality attributes":
		m.section(r, "quality_attributes_content", s)
		r["quality_attributes"] = []string{}
		m.bind(r, "quality_attributes", []string{s.HeadingID}, "structural_wrapper", "quality-attributes-heading")
		for _, child := range s.Children {
			a := m.named("QualityAttribute", r, child)
			a["authored_id"] = m.fieldText(child, "ID")
			definition := m.withoutFields(child, child.Content, "ID", "Sub-attribute of")
			a["definition"] = m.doc.Rich(definition)
			m.bind(a, "definition", definition, "named_field", "attribute-definition")
			m.bind(a, "authored_id", m.fieldNodes(child, "ID"), "record_identity", "attribute-id")
			a["subattribute_of_labels"] = []string{}
			if ids := m.fieldNodes(child, "Sub-attribute of"); len(ids) > 0 {
				for _, label := range strings.Split(m.fieldText(child, "Sub-attribute of"), ",") {
					label = strings.TrimSpace(label)
					if label != "" {
						appendID(a, "subattribute_of_labels", m.label(a, label, ids[0], "subattribute_of"))
					}
				}
				m.bind(a, "subattribute_of_labels", ids, "named_field", "subattribute-labels")
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
		collection := Object(owner["open_questions"])
		if collection == nil {
			collection = Record{"questions": []Record{}}
			owner["open_questions"] = collection
		}
		m.setRich(collection, "content", s.Content)
		m.sectionBindings(owner, "open_questions.content", s)
		values, _ := collection["questions"].([]Record)
		add := func(question []string, item *Section) {
			v := Record{"question": m.doc.Rich(question)}
			index := len(values)
			values = append(values, v)
			collection["questions"] = values
			path := fmt.Sprintf("open_questions.questions.%d", index)
			m.bind(owner, path+".question", question, "named_field", "open-question-question")
			if item == nil {
				return
			}
			m.bind(owner, path+".question", []string{item.HeadingID}, "structural_wrapper", "open-question-heading")
			for _, sub := range item.Children {
				if sub.Title == "Evidence needed" {
					m.setRich(v, "evidence_needed", sub.Content)
					m.sectionBindings(owner, path+".evidence_needed", sub)
				} else {
					m.extension(owner, sub)
				}
			}
		}
		plural := s.Title == "Open questions"
		namedChildren := []*Section{}
		for _, child := range s.Children {
			if child.Title != "Evidence needed" {
				namedChildren = append(namedChildren, child)
			}
		}
		if plural && len(namedChildren) > 0 {
			for _, child := range namedChildren {
				add(append([]string{child.HeadingID}, child.Intro...), child)
			}
			for _, child := range s.Children {
				if child.Title == "Evidence needed" {
					m.extension(owner, child)
				}
			}
		} else if items := m.listItems(s.Intro, false); plural && len(items) > 0 && len(s.Children) == 0 {
			for _, item := range items {
				add([]string{item}, nil)
			}
		} else {
			add(s.Intro, s)
		}
	}
}

func openQuestionHeading(title string) bool {
	return title == "Open question" || title == "Open questions"
}
