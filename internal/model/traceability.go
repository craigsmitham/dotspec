package model

import (
	"slices"

	"github.com/craigsmitham/dotspec/internal/requirement"
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
	if field == "" {
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
	m.section(r, field+"_content", s)
	m.relationship(r, field, s.Content)
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
		if p.field != "refines" && p.field != "part_of" {
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
				for _, next := range d.values(d.index[id], p.field) {
					if reaches(next) {
						return true
					}
				}
				return false
			}
			if reaches(target) {
				d.issue(p.owner, p.field, "error", "Relationship cycle: no cyclic "+p.field+" relationship was projected.")
			} else {
				kept = append(kept, target)
			}
		}
		p.values = kept
	}
}

func (d *designResolver) traceReview() {
	for _, r := range d.elements {
		if !designKind(String(r, "kind"), "Requirement") {
			continue
		}
		if len(d.refs(r, "supports"))+len(d.refs(r, "derived_from"))+len(d.refs(r, "refines")) == 0 {
			d.issue(r, "motivation", "advice", "No explicit motivation or derivation is recorded for this requirement; review its justification.")
		}

		if missing := d.measurementGuidance(r).Missing(); len(missing) > 0 {
			d.issue(r, "measurement_basis", "advice", requirement.MeasurementMessage(missing))
		}
	}
}

// Requirement classification follows declaration ownership, never semantic links.
func governingSystem(r Record, index map[string]Record) string {
	seen := map[string]bool{}
	for r != nil && !seen[String(r, "id")] {
		seen[String(r, "id")] = true
		if String(r, "kind") == "System" {
			return String(r, "id")
		}
		if _, explicit := r["specification_system_content"]; explicit {
			return String(r, "specification_system")
		}
		if id := String(r, "specification_system"); id != "" {
			return id
		}
		r = index[String(r, "owner")]
	}
	return ""
}
