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
	var group []*heading
	save := func() {
		if len(group) == 0 {
			return
		}
		groups[group[0].title] = append(groups[group[0].title], group)
	}
	for _, h := range s.source.headings {
		if h.owner != s || h == s.name || h.inGlossary {
			continue
		}
		if h.level == s.name.level+1 || len(group) == 0 {
			save()
			group = nil
		}
		group = append(group, h)
	}
	save()
	sections[0].passages = [][]*heading{{s.name}}
	for i := range sections {
		section := &sections[i]
		for _, title := range section.sources {
			for _, group := range groups[title] {
				h := group[0]
				// An empty heading that only declares lifted subjects belongs to
				// their index. Preserve incoming links and its precise source route there.
				if dest := s.sectionIndex(title); dest != "" && title != section.title && len(group) == 1 && h.subjectKind == "" && strings.TrimSpace(sourceProse(s.source, h)) == "" {
					index, part := collectionGroup(title)
					s.indexes[index].sources[part] = append(s.indexes[index].sources[part], h)
					p.targets[h.node] = dest
					continue
				}
				section.passages = append(section.passages, group)
			}
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

	var result []overviewSection
	for i, section := range sections {
		populated := len(section.passages) > 0 || len(section.collections) > 0
		if section.title == "Requirements" {
			for _, sub := range sections[i+1:] {
				if sub.level != 3 {
					break
				}
				populated = populated || len(sub.passages) > 0 || len(sub.collections) > 0
			}
		}
		if populated {
			result = append(result, section)
		}
	}
	return result
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
		if section.level > 0 {
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

var collectionDescriptions = map[string]string{
	"Glossary":           "Definitions of terms used in this system.",
	"Business rules":     "Policies that govern the system.",
	"Jobs":               "The work people need to accomplish and the outcomes they seek.",
	"Systems":            "Immediate child systems and their boundaries.",
	"Components":         "Parts of this system and their responsibilities.",
	"Interfaces":         "Required and provided interfaces and their contracts.",
	"Features":           "Product behavior and the specifications each feature groups.",
	"Capabilities":       "What the system can provide and who provides it.",
	"Use cases":          "Actor goals, interactions, and outcomes.",
	"Concepts":           "Definitions of purposeful behavior, state, and actions.",
	"Concept uses":       "How this system uses each concept.",
	"Synchronizations":   "How actions across concept uses are coordinated.",
	"Requirements":       "Standalone requirements and embedded obligations, including routes by component.",
	"Quality attributes": "Declared attributes and their definitions.",
	"Quality models":     "Reference classifications and their definitions.",
}

func overviewIndexLabel(title string) string {
	switch title {
	case "Open questions":
		return "All open questions"
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
		switch section.title {
		case "Quality attributes":
			p.attributeOverview(w, s)
		case "System topology":
			p.renderTopology(w, s)
		case "Obligations":
			p.overviewRefinement(w, s)
		}
		for _, title := range section.collections {
			idx := s.indexes[title]
			line := w.link(overviewIndexLabel(title), idx.file)
			if title != "Quality attributes" {
				line += " — " + collectionDescriptions[title]
			}
			w.text(line)
			entries := idx.groups[""]
			switch section.title {
			case "Behavior", "Concepts", "User needs":
				if len(entries) <= 3 {
					var lines []string
					for _, e := range entries {
						lines = append(lines, entryLine(w, e))
					}
					w.text(strings.Join(lines, "\n"))
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
	questions := s.indexes["Open questions"]
	if len(s.overview) <= 1 && questions == nil {
		return
	}
	w.text("| Section | Explore |\n| --- | --- |")
	for i, section := range s.overview {
		if section.level != 2 {
			continue
		}
		var links []string
		for _, sub := range s.overview[i:] {
			if sub.anchor != section.anchor && sub.level == 2 {
				break
			}
			for _, title := range sub.collections {
				links = append(links, w.link(overviewIndexLabel(title), s.indexes[title].file))
			}
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
		fmt.Fprintf(&w.out, "| %s | %s |\n", w.link(section.title, s.file+"#"+section.anchor), strings.Join(links, " · "))
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
			if owner == nil || owner.kind() != "UseCase" || p.targets[target] == "" {
				continue
			}
			if refinement != nil {
				return
			}
			refinement, parent = r, candidate
		}
	}
	if refinement != nil {
		w.text(w.link(refinement.name.title, refinement.file) + " refines " + w.link(attr(p.nodes[nameNode(parent)], "decoded_text"), p.targets[id(parent)]) + ".")
	}
}
