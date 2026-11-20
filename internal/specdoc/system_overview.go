package specdoc

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

// Each authored passage has one home in the overview. The same outline drives
// rendering, source content order, fragment allocation and the navigation map.
type overviewSection struct {
	title, anchor        string
	level                int
	sources, collections []string
	passages             [][]*heading
	detail               string
	detailLinks          []string
	attributes           []*subject
	preview              []entry
	hasBody              bool
}

func overviewOutline() []overviewSection {
	return []overviewSection{
		{sources: []string{"description"}},
		{title: "Purpose", level: 2, sources: []string{"purpose"}},
		{title: "Goals", level: 2, sources: []string{"goals"}, collections: []string{"Goals"}},
		{title: "Business rules", level: 2, collections: []string{"Business rules"}},
		{title: "User needs", level: 2, sources: []string{"stakeholders"}, collections: []string{"Jobs"}},
		{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
		{title: "Scope", level: 2, sources: []string{"scope", "boundary_description", "out_of_scope", "environment", "assumptions", "constraints", "project_context", "glossary"}, collections: []string{"Assumptions", "Glossary"}},
		{title: "System topology", level: 2, sources: []string{"part_of", "interacts_with", "required_interfaces", "provided_interfaces"}, collections: []string{"Systems", "Components", "Interfaces"}},
		{title: "Features", level: 2, collections: []string{"Features"}},
		{title: "Capabilities", level: 2, collections: []string{"Capabilities"}},
		{title: "Use cases", level: 2, collections: []string{"Use cases"}},
		{title: "Concepts", level: 2, collections: []string{"Concepts", "Concept uses", "Synchronizations"}},
		{title: "Concept uses", level: 3, sources: []string{"concept_uses"}, collections: []string{"Concept uses"}},
		{title: "Synchronizations", level: 3, sources: []string{"synchronizations"}, collections: []string{"Synchronizations"}},
		{title: "Requirements", level: 2, collections: []string{"Requirements"}},
		{title: "Quality attributes", level: 2, collections: []string{"Quality attributes", "Quality models"}},
		{title: "Contracts", level: 2, sources: []string{"contracts"}},
		{title: "Architecture", level: 2, sources: []string{"architecture"}},
		{title: "Implementation", level: 2, sources: []string{"implementation"}},
		{title: "Open questions", level: 2, sources: []string{"open_questions"}, collections: []string{"Open questions"}},
		{title: "Reference", level: 2, sources: []string{"related_specifications", "specification_system", "authority_and_evidence", "evidence", "supports", "derived_from", "relies_on"}},
	}
}

func (p *planner) planOverview(s *subject) []overviewSection {
	sections := overviewOutline()
	groups := map[string][][]*heading{}
	for _, passage := range s.passages {
		h := passage[0]
		if h == s.name {
			sections[0].passages = append(sections[0].passages, passage)
		} else {
			field, _, _ := strings.Cut(h.field, ".")
			field = strings.TrimSuffix(field, "_content")
			groups[field] = append(groups[field], passage)
		}
	}
	for i := range sections {
		section := &sections[i]
		for _, title := range section.sources {
			section.passages = append(section.passages, groups[title]...)
			delete(groups, title)
		}
		section.collections = slices.DeleteFunc(section.collections, func(title string) bool { return s.indexes[title] == nil })
	}
	// Decide which headings exist before allocating their Markdown fragments.
	// Keep a collection in the map even when it supplies no body content.
	var parent *overviewSection
	for i := range sections {
		section := &sections[i]
		detail := &writer{file: s.file, page: &Page{}}
		p.overviewDetails(detail, s, section)
		section.detail = detail.out.String()
		section.detailLinks = detail.page.RequiredLinks
		section.hasBody = len(section.passages) > 0 || section.detail != "" || len(section.preview) > 0 || len(section.attributes) > 0
		if section.level == 2 {
			parent = section
		} else if section.level == 3 && section.hasBody {
			parent.hasBody = true
		}
	}
	return slices.DeleteFunc(sections, func(section overviewSection) bool {
		return !section.hasBody && len(section.collections) == 0
	})
}

func (section overviewSection) merges(i int) bool {
	h := section.passages[i][0]
	if section.level == 0 {
		return h.title == "Description" && i == 1
	}
	return i == 0 && h.title == section.title
}

func (p *planner) indexOverviewHeadings(s *subject, ids parser.IDs) {
	// Markdown renderers derive automatic fragments from heading markup. A
	// rebased title link can change both its own fragment and later duplicate
	// suffixes. Use explicit addresses on outlines containing rich headings.
	explicit := slices.ContainsFunc(s.source.headings, func(h *heading) bool {
		return h.owner == s && h.detailEntry == nil && needsHeadingAnchor(s, h)
	})
	titleID := string(ids.Generate([]byte(s.name.title), ast.KindHeading))
	if explicit {
		s.name.explicitAnchor = "overview.lead"
		titleID = s.name.explicitAnchor
	}
	for i := range s.overview {
		section := &s.overview[i]
		if section.level > 0 && section.hasBody {
			section.anchor = string(ids.Generate([]byte(section.title), ast.KindHeading))
			if explicit {
				section.anchor = "overview." + section.anchor
			}
		}
		for j, group := range section.passages {
			for k, h := range group {
				if h == s.name || h.detailEntry != nil {
					continue
				}
				anchor := section.anchor
				if k != 0 || !section.merges(j) {
					anchor = string(ids.Generate([]byte(h.title), ast.KindHeading))
					if explicit {
						h.explicitAnchor = "passage." + h.anchor
						anchor = h.explicitAnchor
					}
				} else if section.level == 0 {
					anchor = titleID
				}
				p.targets[h.node] = s.file + "#" + anchor
			}
		}
	}
}

func overviewIndexLabel(title string) string {
	switch title {
	case "Goals", "Assumptions", "Constraints", "Business rules", "Features", "Capabilities", "Use cases", "Requirements", "Open questions", "Quality attributes":
		return "All " + strings.ToLower(title)
	}
	return collectionTitle(title)
}

func (p *planner) renderOverview(w *writer, s *subject) {
	for _, title := range slices.Sorted(maps.Keys(s.indexes)) {
		p.renderIndex(s, title, s.indexes[title])
	}
	for i, section := range s.overview {
		if !section.hasBody {
			continue
		}
		if section.level > 0 {
			if s.name.explicitAnchor != "" {
				w.text(`<a id="` + section.anchor + `"></a>`)
			}
			w.text(strings.Repeat("#", section.level) + " " + section.title)
		}
		for j, group := range section.passages {
			merged := section.merges(j)
			offset := section.level - 1
			if merged {
				offset--
			}
			p.renderPassage(w, s, group, offset, merged)
		}
		if i == 0 {
			if s.parent != nil {
				w.text(w.link(s.parent.name.title, s.parent.file))
			}
			p.overviewNavigation(w, s)
			continue
		}
		w.out.WriteString(section.detail)
		for _, attribute := range section.attributes {
			w.text(attributeEntry(w, attribute, true))
		}
		if len(section.preview) > 0 {
			var lines []string
			for _, e := range section.preview {
				lines = append(lines, p.collectionEntryLine(w, e))
			}
			w.text(strings.Join(lines, "\n"))
		}
		for _, dest := range section.detailLinks {
			w.requireLink(dest)
		}
	}
}

func (p *planner) overviewDetails(w *writer, s *subject, section *overviewSection) {
	switch section.title {
	case "Quality attributes":
		for _, key := range model.Strings(s.r["quality_attributes"]) {
			if a := p.pages[key]; a != nil && len(model.Strings(a.r["subattribute_of"])) == 0 {
				section.attributes = append(section.attributes, a)
			}
		}
	case "System topology":
		p.renderTopology(w, s)
	case "Features", "Concept uses", "Synchronizations":
		for _, title := range section.collections {
			idx := s.indexes[title]
			for _, group := range slices.Sorted(maps.Keys(idx.groups)) {
				if group == "By component" {
					continue
				}
				entries := idx.groups[group]
				sortEntries(entries)
				for _, e := range entries {
					if e.preview == "" && len(p.entryRelations[e.target]) == 0 {
						continue
					}
					section.preview = append(section.preview, e)
				}
			}
		}
	}
}

func entryLine(w *writer, e entry) string {
	line := "- " + w.link(inlineText(e.title), e.target)
	if e.preview != "" {
		line += " — " + inlineText(e.preview)
	}
	return line
}

func (p *planner) overviewNavigation(w *writer, s *subject) {
	if len(s.overview) <= 1 {
		return
	}
	w.text("| Section | Explore |\n| --- | --- |")
	for _, section := range s.overview {
		if section.level != 2 {
			continue
		}
		var links []string
		for _, title := range section.collections {
			idx := s.indexes[title]
			links = append(links, w.link(overviewIndexLabel(title), idx.file))
			for _, group := range slices.Sorted(maps.Keys(idx.sources)) {
				if group != "" {
					links = append(links, w.link(group, s.sectionIndex(group)))
				}
			}
		}

		if len(links) == 0 {
			for _, group := range section.passages {
				h := group[0]
				if p.targets[h.node] == s.file+"#"+section.anchor && h.title == section.title {
					continue
				}
				links = append(links, inlineText(h.title))
			}
		}
		if len(links) == 0 && len(section.passages) > 0 {
			h := section.passages[0][0]
			prose := sourceProse(s.source, h)
			preview := paragraphPreview(prose)
			if preview == "" {
				preview = contentText(prose)
			}
			links = append(links, inlineText(preview))
		}
		label := section.title
		if section.hasBody {
			label = w.link(label, s.file+"#"+section.anchor)
		}
		fmt.Fprintf(&w.out, "| %s | %s |\n", label, strings.Join(links, " · "))
	}
}
