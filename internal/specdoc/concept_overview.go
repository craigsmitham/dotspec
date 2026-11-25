package specdoc

import (
	"cmp"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
)

// Key groups by the resolved definition identity. Two definitions may share a
// title, and several components may each declare a use of the same definition.
type conceptGroup struct {
	definition entryRelation
	uses       []entry
}

func (p *planner) conceptOverview(owner *subject) []overviewSection {
	sections := []overviewSection{{title: "Concepts", collection: "Concepts", level: 2, sources: []string{"concept_uses"}}}
	groups := map[string]*conceptGroup{}
	if idx := owner.indexes["Concepts"]; idx != nil {
		for _, e := range idx.groups[""] {
			key := id(p.subjectsByFile[e.target].r)
			groups[key] = &conceptGroup{definition: entryRelation{element: key, title: e.title}}
		}
	}
	if idx := owner.indexes["Concept uses"]; idx != nil {
		for _, e := range idx.groups[""] {
			var def entryRelation
			for _, r := range p.entryRelations[e.target] {
				if r.label == relations["definition"].forward {
					def = r
					break
				}
			}
			if groups[def.element] == nil {
				groups[def.element] = &conceptGroup{definition: def}
			}
			groups[def.element].uses = append(groups[def.element].uses, e)
		}
	}
	var children []overviewSection
	for key, group := range groups {
		children = append(children, p.conceptGroupSection(owner, key, group))
	}
	// The linked caption breaks title ties with a project-relative destination.
	// Snapshot record IDs can differ between checkouts.
	slices.SortFunc(children, func(a, b overviewSection) int {
		if a.caption == "" || b.caption == "" {
			return strings.Compare(b.caption, a.caption)
		}
		return cmp.Or(strings.Compare(a.title, b.title), strings.Compare(a.caption, b.caption))
	})
	for _, child := range children {
		sections[0].detailLinks = append(sections[0].detailLinks, child.detailLinks...)
	}
	return append(sections, children...)
}

func (p *planner) conceptGroupSection(owner *subject, key string, g *conceptGroup) overviewSection {
	w := &writer{file: owner.file, page: &Page{}}
	section := overviewSection{title: "Concept uses", level: 3}
	if key != "" {
		section.title = g.definition.title
		section.caption = p.entryRelationLink(w, g.definition)
		if def := p.pages[key]; def != nil {
			section.conceptMembers = append(section.conceptMembers, def.file)
		}
		w.text(inlineText(indexPreview(p.conceptPurpose(key))))
		if len(g.uses) > 0 {
			label := "Use in this "
			if len(g.uses) > 1 {
				label = "Uses in this "
			}
			w.text("**" + label + strings.ToLower(owner.kind()) + "**")
		}
	}
	sortEntries(g.uses)
	var lines []string
	for _, use := range g.uses {
		line := entryLine(w, use)
		declaring := p.subjectsByFile[use.target].navigationParent
		if declaring != owner && declaring.kind() == "Component" {
			line += "\\\n  **Declared in** " + w.link(inlineText(declaring.name.title), declaring.file)
		}
		lines = append(lines, line)
		section.conceptMembers = append(section.conceptMembers, use.target)
	}
	w.text(strings.Join(lines, "\n"))
	section.detail = w.out.String()
	section.detailLinks = w.page.RequiredLinks
	return section
}

func (p *planner) conceptPurpose(key string) string {
	if def := p.pages[key]; def != nil {
		return subjectPreview(def)
	}
	// Reference context supplies the external definition without publishing a
	// report for it. Read its purpose from that same snapshot.
	r := p.elements[key]
	for _, source := range p.snapshot.Sources {
		if id(source) != attr(r, "source_document") {
			continue
		}
		data := attr(source, "markdown")
		for _, n := range model.Strings(model.Object(r["purpose"])["content_nodes"]) {
			node := p.nodes[n]
			if attr(node, "node_kind") == "paragraph" {
				start, end := position(node, "start"), position(node, "end")
				if start >= 0 && end >= start && end <= len(data) {
					return paragraphPreview(data[start:end])
				}
			}
		}
	}
	return ""
}
