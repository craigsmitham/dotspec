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
		{sources: []string{"description"}},
		{title: "Purpose", level: 2, sources: []string{"purpose"}},
		{title: "Goals", level: 2, sources: []string{"goals"}, collections: []string{"Goals"}},
		{title: "Business rules", level: 2, collections: []string{"Business rules"}},
		{title: "User needs", level: 2, sources: []string{"stakeholders"}, collections: []string{"Jobs"}},
		{title: "Scope", level: 2, sources: []string{"scope", "out_of_scope", "environment", "assumptions", "constraints", "project_context", "glossary"}, collections: []string{"Assumptions", "Glossary"}},
		{title: "System topology", level: 2, sources: []string{"part_of", "interacts_with", "required_interfaces", "provided_interfaces"}, collections: []string{"Systems", "Components", "Interfaces"}},
		{title: "Features", level: 2, collections: []string{"Features"}},
		{title: "Capabilities", level: 2, collections: []string{"Capabilities"}},
		{title: "Use cases", level: 2, collections: []string{"Use cases"}},
		{title: "Concepts", level: 2, sources: []string{"concept_uses", "synchronizations"}, collections: []string{"Concepts", "Concept uses", "Synchronizations"}},
		{title: "Requirements", level: 2, sources: []string{"obligations"}, collections: []string{"Requirements"}},
		{title: "Quality attributes", level: 2, collections: []string{"Quality attributes", "Quality models"}},
		{title: "Open questions", level: 2, sources: []string{"open_questions"}},
		{title: "Reference", level: 2, sources: []string{"related_specifications", "authority_and_evidence", "evidence", "supports", "derived_from", "relies_on"}},
	}
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
	topologyCapabilities := map[string]bool{}
	for i := range sections {
		section := &sections[i]
		detail := &writer{file: s.file, page: &Page{}}
		p.overviewDetails(detail, s, *section, topologyCapabilities)
		section.detail = detail.out.String()
		section.detailLinks = detail.page.RequiredLinks
		section.hasBody = len(section.passages) > 0 || section.detail != ""
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
	return i == 0 && (h.title == section.title || (section.title == "Requirements" && h.title == "Obligations"))
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
	case "Goals", "Assumptions", "Constraints", "Business rules", "Features", "Capabilities", "Use cases", "Requirements", "Open questions":
		return "All " + strings.ToLower(title)
	case "Quality attributes":
		return "Attribute index"
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

func (p *planner) overviewDetails(w *writer, s *subject, section overviewSection, topologyCapabilities map[string]bool) {
	switch section.title {
	case "Quality attributes":
		for _, key := range model.Strings(s.r["quality_attributes"]) {
			if a := p.pages[key]; a != nil && len(model.Strings(a.r["subattribute_of"])) == 0 {
				w.text(attributeEntry(w, a, true))
				var classified []entry
				for _, r := range p.snapshot.Elements {
					if !slices.Contains(model.Strings(r["quality_attributes"]), key) || (attr(r, "kind") != "Requirement" && attr(r, "kind") != "ObligationClause") {
						continue
					}
					n := nameNode(r)
					if target := p.targets[n]; target != "" {
						classified = append(classified, entry{title: attr(p.nodes[n], "decoded_text"), target: target})
					}
				}
				sortEntries(classified)
				var links []string
				for _, e := range classified {
					links = append(links, w.link(inlineText(e.title), e.target))
				}
				if len(links) > 0 {
					w.text("Requirements classified under " + inlineText(a.name.title) + ": " + strings.Join(links, ", ") + ".")
				}
			}
		}
	case "System topology":
		p.renderTopology(w, s, topologyCapabilities)
	case "Requirements":
		p.overviewRefinement(w, s)
	case "Business rules", "Features", "Capabilities", "Use cases", "Concepts", "User needs":
		definitions := map[string]*subject{}
		if section.title == "Concepts" {
			for _, useID := range model.Strings(s.r["concept_uses"]) {
				if use := p.pages[useID]; use != nil {
					definitions[use.file] = p.pages[attr(use.r, "definition")]
				}
			}
		}
		for _, title := range section.collections {
			group := ""
			if title == "Business rules" {
				group = applicableRules(s)
			}
			entries := s.indexes[title].groups[group]
			if title == "Capabilities" {
				entries = slices.Concat(entries, s.indexes[title].groups["Provided capabilities"])
			}
			if len(entries) == 0 || len(entries) > 3 {
				continue
			}
			sortEntries(entries)
			var lines []string
			for _, e := range entries {
				def := definitions[e.target]
				if topologyCapabilities[e.target] || (e.preview == "" && def == nil) {
					continue
				}
				line := entryLine(w, e)
				if def != nil {
					line += "\n  - Uses " + w.link(inlineText(def.name.title), def.file)
				}
				lines = append(lines, line)
			}
			if len(lines) > 0 && section.title == "Concepts" && len(section.collections) > 1 {
				w.text("**" + overviewIndexLabel(title) + "**")
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
		useCase := p.ownerPage(parent)
		parentLink := w.link(inlineText(attr(p.nodes[nameNode(parent)], "decoded_text")), p.targets[nameNode(parent)])
		w.text(w.link(inlineText(useCase.name.title), useCase.file) + " includes " + parentLink + ". " + inlineText(p.requirementPreview(parent)))
		w.text(w.link(inlineText(refinement.name.title), refinement.file) + " refines that promise. " + inlineText(subjectPreview(refinement)))
	}
}
