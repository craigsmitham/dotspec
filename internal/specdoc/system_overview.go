package specdoc

import (
	"cmp"
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Each authored passage has one home in the overview. The same outline drives
// rendering, source content order, fragment allocation and the navigation map.
type overviewSection struct {
	collection     string // directory or interface group represented by this section
	title, anchor  string
	caption        string // linked heading for a concept definition
	conceptMembers []string
	level          int
	sources        []string
	passages       [][]*heading
	detail         string
	detailLinks    []string
	preview        []entry
	hasBody        bool
}

// These System fields each have their own section in the report, in reading order.
var systemContextFields = []string{"scope", "out_of_scope", "environment", "assumptions", "constraints"}

func (p *planner) overviewOutline(s *subject) []overviewSection {
	sections := []overviewSection{
		{sources: []string{"description"}},
		{title: "Framing", level: 2, sources: []string{"framing"}},
		{title: "Purpose", level: 2, sources: []string{"purpose"}},
		{title: "Goals", level: 2, sources: []string{"goals"}},
		{title: "Stakeholders", level: 2, sources: []string{"stakeholder_roles"}},
		{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
	}
	if s.kind() == "System" {
		for _, field := range systemContextFields {
			title := sentenceCase(strings.ReplaceAll(field, "_", " "))
			level := 2
			if field == "out_of_scope" {
				level = 3
			}
			sections = append(sections, overviewSection{title: title, level: level, sources: []string{field}})
		}
	} else {
		sections = append(sections, overviewSection{title: "Scope", level: 2, sources: []string{"boundary_description", "out_of_scope"}})
	}
	sections = append(sections, p.architecture(s).sections()...)
	capabilityLevel := 2
	if s.kind() == "Component" {
		capabilityLevel = 3
	}
	sections = append(sections, overviewSection{title: "Capabilities", collection: "Capabilities", level: capabilityLevel, sources: []string{"capabilities"}})
	sections = append(sections, overviewSection{title: "Features", collection: "Features", level: 2})
	sections = append(sections, p.conceptOverview(s)...)
	return append(sections, []overviewSection{
		{title: "Synchronizations", level: 3, sources: []string{"synchronizations"}},
		{title: "Quality Attributes", collection: "Quality attributes", level: 2},
		{title: "Adopts", level: 2, sources: []string{"adopts"}},
		{title: "Contracts", level: 2, sources: []string{"contracts"}},
		{title: "Implementation", level: 2, sources: []string{"implementation"}},
		{title: "Open questions", level: 2, sources: []string{"open_questions"}},
		{title: "Glossary", level: 2, sources: []string{"glossary"}},
		{title: "Project references", level: 2, sources: []string{"project_references"}},
		{title: "Reference", level: 2, sources: []string{"related_specifications", "specification_system", "authority_and_evidence", "evidence", "supports", "derived_from", "relies_on"}},
	}...)
}

func (p *planner) planOverview(s *subject) []overviewSection {
	sections := p.overviewOutline(s)
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
	}
	// Decide which headings exist before allocating their Markdown fragments.
	var parent *overviewSection
	for i := range sections {
		section := &sections[i]
		detail := &writer{file: s.file, page: &Page{}}
		p.overviewDetails(detail, s, section)
		section.detail += detail.out.String()
		section.detailLinks = append(section.detailLinks, detail.page.RequiredLinks...)
		section.hasBody = len(section.passages) > 0 || section.detail != "" || section.caption != "" || len(section.preview) > 0
		if section.level == 2 {
			parent = section
		} else if section.level == 3 && section.hasBody {
			parent.hasBody = true
		}
	}
	return slices.DeleteFunc(sections, func(section overviewSection) bool {
		return !section.hasBody
	})
}

func (section overviewSection) merges(i int) bool {
	h := section.passages[i][0]
	if section.level == 0 {
		return h.title == "Description" && i == 1
	}
	return i == 0 && (h.title == section.title || (section.title == "Interactions" && h.title == "Interacts with"))
}

func (p *planner) indexOverviewHeadings(s *subject) {
	ids := specfile.RenderedHeadingIDs{}
	explicit := slices.ContainsFunc(s.source.headings, func(h *heading) bool {
		return h.owner == s && h.detailEntry == nil && needsHeadingAnchor(s, h)
	})
	titleID := ids.Next(s.name.title)
	if explicit {
		s.name.explicitAnchor = "overview.lead"
		titleID = s.name.explicitAnchor
	}
	for i := range s.overview {
		section := &s.overview[i]
		if section.level > 0 {
			section.anchor = ids.Next(section.title)
			if explicit && section.caption == "" {
				section.anchor = "overview." + section.anchor
			}
		}
		for j, group := range section.passages {
			for k, h := range group {
				if h.detailEntry != nil {
					continue
				}
				anchor := section.anchor
				if h != s.name && (k != 0 || !section.merges(j)) {
					anchor = ids.Next(h.reportTitle())
					if explicit || anchor == "" {
						h.explicitAnchor = "passage." + h.anchor
						anchor = h.explicitAnchor
					}
				} else if section.level == 0 {
					anchor = titleID
				}
				if h != s.name {
					p.targets[h.node] = s.file + "#" + anchor
				}
				for _, title := range h.bodyHeadings {
					ids.Next(title)
				}
			}
		}
	}
}

// Use an overview section as the directory only when every member has a
// destination there. Directory discovery can extend beyond declared parts,
// and component roll-ups lead to collections outside the overview.
func (p *planner) planCollectionHomes(s *subject) {
	byCollection := map[string]map[string]*overviewSection{}
	s.conceptSections = map[string]*overviewSection{}
	var parent *overviewSection
	for i := range s.overview {
		section := &s.overview[i]
		if section.level == 2 {
			parent = section
		}
		for _, file := range section.conceptMembers {
			s.conceptSections[file] = section
		}
		if section.collection == "" || (s.kind() == "Component" && !inlineComponentCollection(section.collection)) {
			continue
		}
		title, group := collectionGroup(section.collection)
		if byCollection[title] == nil {
			byCollection[title] = map[string]*overviewSection{"": section}
		} else if group != "" {
			// When interfaces span two directions, the common Architecture
			// section is the destination for citations to the whole collection.
			byCollection[title][""] = parent
		}
		byCollection[title][group] = section
	}
	if concepts := byCollection["Concepts"]; concepts != nil {
		byCollection["Concept uses"] = concepts
	}
	for _, title := range slices.Sorted(maps.Keys(byCollection)) {
		sections := byCollection[title]
		idx := s.indexes[title]
		if idx == nil || !idx.coveredBy(s, sections) {
			continue
		}
		idx.overviewSections = sections
		for _, group := range slices.Sorted(maps.Keys(idx.sources)) {
			section := sections[group]
			for _, h := range idx.sources[group] {
				p.targets[h.node] = s.file + "#" + section.anchor
			}
		}
	}
}

func (idx *collectionIndex) coveredBy(s *subject, sections map[string]*overviewSection) bool {
	for group, members := range idx.groups {
		section := sections[group]
		if section == nil {
			return false
		}
		links := map[string]bool{}
		for _, dest := range section.detailLinks {
			links[dest] = true
		}
		for _, e := range section.preview {
			links[relative(s.file, e.target)] = true
		}
		for _, e := range members {
			if !links[relative(s.file, e.target)] {
				return false
			}
		}
	}
	return true
}

func (s *subject) browseCollection(title string, idx *collectionIndex) bool {
	if idx.overviewSections != nil {
		return false
	}
	if s.kind() != "System" {
		return true
	}
	if isDetailCollection(title) && len(idx.passages) > 0 && len(idx.groups) == 0 {
		return false
	}
	return title != "Quality models"
}

func (p *planner) renderOverview(w *writer, s *subject) {
	for _, title := range slices.Sorted(maps.Keys(s.indexes)) {
		if idx := s.indexes[title]; idx.overviewSections == nil {
			p.renderIndex(s, title, idx)
		}
	}
	for _, section := range s.overview {
		if section.level > 0 {
			if s.name.explicitAnchor != "" && section.caption == "" {
				w.text(`<a id="` + section.anchor + `"></a>`)
			}
			title := cmp.Or(section.caption, section.title)
			w.text(strings.Repeat("#", section.level) + " " + title)
		}
		for j, group := range section.passages {
			merged := section.merges(j)
			offset := section.level - 1
			if merged {
				offset--
			}
			p.renderPassage(w, s, group, offset, merged)
		}
		if section.level == 0 {
			overviewNavigation(w, s)
			continue
		}
		w.out.WriteString(section.detail)
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
	if len(section.conceptMembers) != 0 {
		return
	}
	switch section.title {
	case "Capabilities":
		if idx := s.indexes["Capabilities"]; idx != nil {
			p.capabilityOverview(w, s, idx.groups[""])
		}
	case "Quality Attributes":
		p.attributeOverview(w, s)
	case "Features", "Synchronizations":
		if idx := s.indexes[section.title]; idx != nil {
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

func (p *planner) capabilityOverview(w *writer, owner *subject, entries []entry) {
	if len(entries) == 0 {
		return
	}
	var rows []string
	if owner.kind() == "System" {
		rows = append(rows, "| Capability | What it enables | Provided by |", "| --- | --- | --- |")
	}
	sortEntries(entries)
	for _, e := range entries {
		e.preview = indexPreview(e.preview)
		if owner.kind() == "Component" {
			rows = append(rows, entryLine(w, e))
			continue
		}
		var providers []string
		for _, r := range p.entryRelations[e.target] {
			if r.label == relations["provided_by"].forward {
				providers = append(providers, p.entryRelationLink(w, r))
			}
		}
		slices.Sort(providers)
		rows = append(rows, "| "+w.link(inlineText(e.title), e.target)+" | "+inlineText(e.preview)+" | "+strings.Join(slices.Compact(providers), " · ")+" |")
	}
	w.text(strings.Join(rows, "\n"))
}

func entryLine(w *writer, e entry) string {
	line := "- " + w.link(inlineText(e.title), e.target)
	if e.preview != "" {
		line += " — " + inlineText(e.preview)
	}
	return line
}

func overviewNavigation(w *writer, s *subject) {
	var sections []string
	for _, section := range s.overview {
		if section.level == 2 {
			sections = append(sections, w.link(section.title, s.file+"#"+section.anchor))
		}
	}
	w.disclosure("On this page", sections)

	titles := slices.SortedFunc(maps.Keys(s.indexes), func(a, b string) int {
		return strings.Compare(collectionTitle(a), collectionTitle(b))
	})
	var collections []string
	for _, title := range titles {
		idx := s.indexes[title]
		if !s.browseCollection(title, idx) {
			continue
		}
		line := w.link(collectionTitle(title), idx.file)
		for _, group := range slices.Sorted(maps.Keys(idx.sources)) {
			if group != "" {
				line += "\n  - " + w.link(group, s.sectionIndex(group))
			}
		}
		collections = append(collections, line)
	}
	w.disclosure("Browse specifications", collections)
}
