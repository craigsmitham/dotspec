package model

import (
	"fmt"
	"slices"
)

// Traceability uses the schema's heading and endpoint contracts. Ordinary links
// elsewhere remain citations, even when they name the same target.
func (m *domainMapper) traceSection(r Record, s *Section) bool {
	cls := compiledModel().catalog.Classes[String(r, "kind")]
	if !slices.Contains(cls.Ancestors, "TraceableElement") {
		return false
	}
	field := ""
	for name, slot := range cls.Slots {
		if slot.Annotations["trace_heading"] == s.Title {
			field = name
			break
		}
	}
	if field == "" && s.Title != "Obligations" {
		return false
	}
	if m.traceMapped == nil {
		m.traceMapped = map[string]bool{}
	}
	key := String(r, "id") + "/" + s.ID
	if m.traceMapped[key] {
		return true
	}
	m.traceMapped[key] = true
	if field != "" {
		m.section(r, field+"_content", s)
		r[field+"_references"] = append(Strings(r[field+"_references"]), m.references(s.Content)...)
		m.bindDesignReferences(r, field+"_references")
		return true
	}
	m.section(r, "obligations_content", s)
	if r["obligations"] == nil {
		r["obligations"] = []string{}
	}
	for _, child := range s.Children {
		x := m.named("ObligationClause", r, child)
		if id := m.fieldText(child, "ID"); id != "" {
			x["authored_id"] = id
			m.bind(x, "authored_id", m.fieldNodes(child, "ID"), "record_identity", "obligation-id")
			m.bind(x, "authored_id", m.fieldLabelNodes(child, "ID"), "structural_wrapper", "obligation-id-label")
		}
		nodes := m.withoutFields(child, child.Intro, "ID")
		x["statement"] = m.doc.Rich(nodes)
		m.bind(x, "statement", nodes, "named_field", "obligation-statement")
		for _, sub := range child.Children {
			if m.traceSection(x, sub) {
				continue
			}
			if sub.Title == "Verification" || sub.Title == "Applies to" || sub.Title == "Rationale" {
				m.requirement(x, sub, "requirement")
			} else {
				m.extension(x, sub)
			}
		}
		appendID(r, "obligations", String(x, "id"))
	}
	return true
}

// Some nested records retain a complete section without parsing its children.
// Recognize trace sections there too, without decomposing arbitrary prose.
func (m *domainMapper) nestedTraceability() {
	sections := map[string]*Section{}
	var visit func([]*Section)
	visit = func(items []*Section) {
		for _, s := range items {
			sections[s.HeadingID] = s
			visit(s.Children)
		}
	}
	visit(m.doc.Sections)
	for i := 0; i < len(m.result.Elements); i++ {
		r := m.result.Elements[i]
		for _, id := range Strings(Object(r["name"])["content_nodes"]) {
			if s := sections[id]; s != nil {
				for _, child := range s.Children {
					m.traceSection(r, child)
				}
			}
		}
	}
}

func (d *designResolver) resolveTraceability() {
	for _, r := range d.elements {
		cls := compiledModel().catalog.Classes[String(r, "kind")]
		for _, field := range sortedKeys(cls.Slots) {
			slot := cls.Slots[field]
			if targets := slot.Annotations["trace_targets"]; targets != "" {
				d.resolve(r, field, targets, true, nil)
			}
		}
	}
	// Preserve authored cycles as occurrences and diagnostics, but do not export
	// a cyclic refinement claim. Every edge participating in a cycle is withheld.
	// Snapshot validation replays the same projection; forged cycles are rejected.
	for i := range d.projections {
		p := &d.projections[i]
		if p.field != "refines" {
			continue
		}
		source := String(p.owner, "id")
		kept := []string{}
		for _, target := range p.values {
			seen := map[string]bool{}
			var reaches func(string) bool
			reaches = func(id string) bool {
				if id == source {
					return true
				}
				if seen[id] {
					return false
				}
				seen[id] = true
				for _, next := range d.values(d.index[id], "refines") {
					if reaches(next) {
						return true
					}
				}
				return false
			}
			if reaches(target) {
				d.issue(p.owner, "refines", "error", "Refinement cycle: no cyclic refinement relationship was projected.")
			} else {
				kept = append(kept, target)
			}
		}
		p.values = kept
	}
}

func (d *designResolver) traceReview() {
	for _, r := range d.elements {
		if !designKind(String(r, "kind"), "Requirement or ObligationClause") {
			continue
		}
		if len(Strings(r["supports_references"]))+len(Strings(r["derived_from_references"]))+len(Strings(r["refines_references"])) == 0 {
			d.issue(r, "motivation", "advice", "No explicit motivation or derivation is recorded for this obligation; review its justification.")
		}
		field := "verification"
		if String(r, "kind") == "QualityRequirement" {
			field = "measurement"
		}
		if d.text(Object(r[field])["content"]) == "" {
			d.issue(r, "verification", "advice", fmt.Sprintf("No %s guidance is recorded for this obligation; review how it will be checked.", field))
		}
	}
}
