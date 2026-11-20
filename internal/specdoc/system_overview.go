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
	hasBody              bool
}

func (p *planner) planOverview(s *subject) []overviewSection {
	sections := []overviewSection{
		{sources: []string{"Description"}},
		{title: "Purpose", level: 2, sources: []string{"Purpose"}},
		{title: "Goals", level: 2, sources: []string{"Goals"}},
		{title: "Business rules", level: 2, sources: []string{"Business rules"}, collections: []string{"Business rules"}},
		{title: "User needs", level: 2, sources: []string{"User needs", "Stakeholders", "Jobs"}, collections: []string{"Jobs"}},
		{title: "Scope", level: 2, sources: []string{"Scope", "Out of scope", "Environment", "Assumptions", "Constraints", "Project context", "Glossary"}, collections: []string{"Glossary"}},
		{title: "System topology", level: 2, sources: []string{"System topology", "Part of", "Child systems", "Systems", "Components", "Interacts with", "Required interfaces", "Provided interfaces", "Interfaces"}, collections: []string{"Systems", "Components", "Interfaces"}},
		{title: "Requirements", level: 2, sources: []string{"Requirements"}},
		{title: "Behavior", level: 3, sources: []string{"Behavior", "Features", "Capabilities", "Use cases"}, collections: []string{"Features", "Capabilities", "Use cases"}},
		{title: "Concepts", level: 3, sources: []string{"Concepts", "Concept uses", "Synchronizations"}, collections: []string{"Concepts", "Concept uses", "Synchronizations"}},
		{title: "Obligations", level: 3, sources: []string{"Obligations", "Applicable requirements"}, collections: []string{"Requirements"}},
		{title: "Quality attributes", level: 2, sources: []string{"Quality attributes", "Quality models"}, collections: []string{"Quality attributes", "Quality models"}},
		{title: "Open questions", level: 2, sources: []string{"Open questions"}},
		{title: "Reference", level: 2, sources: []string{"Reference", "Evidence"}},
	}
	groups := map[string][][]*heading{}
	for _, passage := range s.passages {
		h := passage[0]
		if h == s.name {
			sections[0].passages = append(sections[0].passages, passage)
		} else {
			groups[h.title] = append(groups[h.title], passage)
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
	// Preserve unfamiliar sections in source order, including duplicate titles.
	var extra [][]*heading
	for _, gs := range groups {
		extra = append(extra, gs...)
	}
	slices.SortFunc(extra, func(a, b []*heading) int { return a[0].start - b[0].start })
	last := &sections[len(sections)-1]
	last.passages = append(last.passages, extra...)

	// Decide which headings exist before allocating their Markdown fragments.
	// Keep a collection in the map even when it supplies no body content.
	for i := len(sections) - 1; i >= 0; i-- {
		section := &sections[i]
		detail := &writer{file: s.file, page: &Page{}}
		p.overviewDetails(detail, s, *section)
		section.detail = detail.out.String()
		section.detailLinks = detail.page.RequiredLinks
		section.hasBody = len(section.passages) > 0 || section.detail != ""
		if section.level == 2 {
			for _, sub := range sections[i+1:] {
				if sub.level != 3 {
					break
				}
				section.hasBody = section.hasBody || sub.hasBody
				section.collections = append(section.collections, sub.collections...)
			}
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
		return h.owner == s && needsHeadingAnchor(s, h)
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
				if h == s.name {
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
	case "Open questions":
		return "All open questions"
	case "Business rules":
		return "All business rules"
	case "Quality attributes":
		return "Attribute index"
	case "Requirements":
		return "All requirements"
	case "Concepts":
		return "Concept definitions"
	}
	return title
}

func (p *planner) renderSystem(w *writer, s *subject) {
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
		for _, dest := range section.detailLinks {
			w.requireLink(dest)
		}
	}
}

func (p *planner) overviewDetails(w *writer, s *subject, section overviewSection) {
	switch section.title {
	case "Quality attributes":
		p.attributeOverview(w, s)
	case "System topology":
		p.renderTopology(w, s)
	case "Obligations":
		p.overviewRefinement(w, s)
	case "Business rules", "Behavior", "Concepts", "User needs":
		for _, title := range section.collections {
			group := ""
			if title == "Business rules" {
				group = applicableRules
			}
			entries := s.indexes[title].groups[group]
			if len(entries) == 0 || len(entries) > 3 {
				continue
			}
			if section.title == "Behavior" || section.title == "Concepts" {
				w.text("**" + overviewIndexLabel(title) + "**")
			}
			sortEntries(entries)
			var lines []string
			for _, e := range entries {
				lines = append(lines, entryLine(w, e))
			}
			w.text(strings.Join(lines, "\n"))
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
	questions := s.indexes["Open questions"]
	if len(s.overview) <= 1 && questions == nil {
		return
	}
	w.text("| Section | Explore |\n| --- | --- |")
	for _, section := range s.overview {
		if section.level != 2 {
			continue
		}
		var links []string
		for _, title := range section.collections {
			links = append(links, w.link(overviewIndexLabel(title), s.indexes[title].file))
		}
		if section.title == "Open questions" && questions != nil {
			links = append(links, w.link(overviewIndexLabel(section.title), questions.file))
			questions = nil
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
			links = append(links, inlineText(paragraphPreview(sourceProse(s.source, h))))
		}
		label := section.title
		if section.hasBody {
			label = w.link(label, s.file+"#"+section.anchor)
		} else if len(links) == 1 {
			label, links = links[0], nil
		}
		fmt.Fprintf(&w.out, "| %s | %s |\n", label, strings.Join(links, " · "))
	}
	if questions != nil {
		fmt.Fprintf(&w.out, "| %s | |\n", w.link(overviewIndexLabel("Open questions"), questions.file))
	}
}

// Show the local use-case promise and its refinement when there is one
// unambiguous example. Larger collections have their own complete index.
func (p *planner) overviewRefinement(w *writer, s *subject) {
	var refinement *subject
	var parent model.Record
	for _, r := range p.subjects {
		if r.parent != s || r.kind() != "Requirement" {
			continue
		}
		for _, target := range model.Strings(r.r["refines"]) {
			candidate := p.elements[target]
			owner := p.ownerPage(candidate)
			if owner == nil || owner.kind() != "UseCase" || p.targets[nameNode(candidate)] == "" {
				continue
			}
			if refinement != nil {
				return
			}
			refinement, parent = r, candidate
		}
	}
	if refinement != nil {
		w.text(w.link(refinement.name.title, refinement.file) + " refines " + w.link(attr(p.nodes[nameNode(parent)], "decoded_text"), p.targets[nameNode(parent)]) + ".")
	}
}
