package model

import (
	"fmt"
	"strings"
)

var questionDetails = map[string]string{
	"Raised by": "raised_by", "Raised on": "raised_on", "Resolve by": "resolve_by",
	"Importance": "importance", "Evidence needed": "evidence_needed", "Tracking reference": "tracking_reference",
}

func (m *domainMapper) explanatory(r Record, s *Section) bool {
	field := map[string]string{"Explanatory notes": "explanatory_notes", "Implementation suggestions": "implementation_suggestions", "Implementation hints": "implementation_suggestions", "Illustrative examples": "illustrative_examples"}[s.Title]
	if field == "" {
		return false
	}
	if _, ok := compiledModel().catalog.Classes[String(r, "kind")].Slots[field]; !ok {
		return false
	}
	if m.traceMapped == nil {
		m.traceMapped = map[string]bool{}
	}
	key := String(r, "id") + "/explanation/" + s.ID
	if m.traceMapped[key] {
		return true
	}
	m.traceMapped[key] = true
	m.section(r, field, s)
	return true
}

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
		v := m.collection(r, field, "interfaces", s)
		for _, child := range s.Children {
			item := m.named("Interface", r, child)
			m.section(item, "content", child)
			m.setRich(item, "description", child.Intro)
			m.bind(item, "description", child.Intro, "named_field", "interface-description")
			for _, detail := range child.Children {
				switch detail.Title {
				case "Contracts", "Counterpart":
					f := strings.ToLower(detail.Title)
					m.contextField(item, f, detail.Content)
					m.sectionBindings(item, f+".content", detail)
				default:
					m.extension(item, detail)
				}
			}
			appendID(v, "interfaces", String(item, "id"))
		}
	default:
		return false
	}
	return true
}

func (m *domainMapper) assumptionDetails(r Record, s *Section) {
	for _, sub := range s.Children {
		if m.traceSection(r, sub) || m.explanatory(r, sub) {
			continue
		}
		switch sub.Title {
		case "Basis", "Applicability":
			m.section(r, strings.ToLower(sub.Title), sub)
		case "Validation plan":
			m.assessmentSection(r, sub)
		case "Open question", "Open questions":
			m.openQuestions(r, []*Section{sub})
		default:
			m.extension(r, sub)
		}
	}
}

func (m *domainMapper) effectsTable(r Record, s *Section) {
	v := Record{"content": m.doc.Rich(s.Content), "rows": []Record{}}
	r["effects_table"] = v
	m.sectionBindings(r, "effects_table.content", s)
	for _, id := range s.Intro {
		table := m.doc.NodeByID[id]
		if String(table, "node_kind") != "table" {
			continue
		}
		rows := []Record{}
		for _, rid := range Strings(table["children"]) {
			row := m.doc.NodeByID[rid]
			if String(row, "node_kind") != "table_row" {
				continue
			}
			cells := []Record{}
			for _, cid := range Strings(row["children"]) {
				if String(m.doc.NodeByID[cid], "node_kind") != "table_cell" {
					continue
				}
				cells = append(cells, m.doc.Rich([]string{cid}))
				m.bind(r, fmt.Sprintf("effects_table.rows.%d.cells.%d", len(rows), len(cells)-1), []string{cid}, "structured_item", "effects-cell")
			}
			rows = append(rows, Record{"cells": cells})
		}
		v["rows"] = rows
		break // One authored matrix; the complete section always remains available.
	}
	for _, child := range s.Children {
		m.extension(r, child)
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
				r["illustrates_references"] = append(append([]string{}, Strings(r["illustrates_references"])...), m.references([]string{id})...)
				m.bind(r, "illustrates_content", []string{id}, "named_field", "step-illustrates")
				m.bindDesignReferences(r, "illustrates_references")
				return
			}
		}
		for _, child := range Strings(n["children"]) {
			visit(child)
		}
	}
	visit(item)
}
