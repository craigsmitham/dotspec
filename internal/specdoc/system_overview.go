package specdoc

import (
	"maps"
	"slices"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

// Each authored passage has one home in the overview. The same outline drives
// rendering, source content order, fragment allocation and the navigation map.
type overviewSection struct {
	title, anchor string
	level         int
	sources       []string
	passages      [][]*heading
	detail        string
	detailLinks   []string
	preview       []entry
	hasBody       bool
}

// These System fields each have their own section in the report, in reading order.
var systemContextFields = []string{"scope", "out_of_scope", "environment", "assumptions", "constraints"}

func overviewOutline(kind string) []overviewSection {
	sections := []overviewSection{
		{sources: []string{"description"}},
		{title: "Framing", level: 2, sources: []string{"framing"}},
		{title: "Purpose", level: 2, sources: []string{"purpose"}},
		{title: "Goals", level: 2, sources: []string{"goals"}},
		{title: "Stakeholders", level: 2, sources: []string{"stakeholder_roles"}},
		{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
	}
	if kind == "System" {
		for _, field := range systemContextFields {
			title := sentenceCase(strings.ReplaceAll(field, "_", " "))
			sections = append(sections, overviewSection{title: title, level: 2, sources: []string{field}})
		}
	} else {
		sections = append(sections, overviewSection{title: "Scope", level: 2, sources: []string{"boundary_description", "out_of_scope"}})
	}
	return append(sections, []overviewSection{
		{title: "System topology", level: 2, sources: []string{"part_of", "interacts_with", "required_interfaces", "provided_interfaces"}},
		{title: "Features", level: 2},
		{title: "Concepts", level: 2},
		{title: "Concept uses", level: 3, sources: []string{"concept_uses"}},
		{title: "Synchronizations", level: 3, sources: []string{"synchronizations"}},
		{title: "Quality Attributes", level: 2},
		{title: "Adopts", level: 2, sources: []string{"adopts"}},
		{title: "Contracts", level: 2, sources: []string{"contracts"}},
		{title: "Architecture", level: 2, sources: []string{"architecture"}},
		{title: "Implementation", level: 2, sources: []string{"implementation"}},
		{title: "Open questions", level: 2, sources: []string{"open_questions"}},
		{title: "Glossary", level: 2, sources: []string{"glossary"}},
		{title: "Project references", level: 2, sources: []string{"project_references"}},
		{title: "Reference", level: 2, sources: []string{"related_specifications", "specification_system", "authority_and_evidence", "evidence", "supports", "derived_from", "relies_on"}},
	}...)
}

func (p *planner) planOverview(s *subject) []overviewSection {
	sections := overviewOutline(s.kind())
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
		section.detail = detail.out.String()
		section.detailLinks = detail.page.RequiredLinks
		section.hasBody = len(section.passages) > 0 || section.detail != "" || len(section.preview) > 0
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
		if section.level > 0 {
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
					anchor = string(ids.Generate([]byte(h.reportTitle()), ast.KindHeading))
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

func (p *planner) renderOverview(w *writer, s *subject) {
	for _, title := range slices.Sorted(maps.Keys(s.indexes)) {
		p.renderIndex(s, title, s.indexes[title])
	}
	for _, section := range s.overview {
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
	switch section.title {
	case "Quality Attributes":
		p.attributeOverview(w, s)
	case "System topology":
		p.renderTopology(w, s)
	case "Features", "Concept uses", "Synchronizations":
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
		if s.kind() == "System" && (title == "Quality attributes" || title == "Quality models") {
			continue
		}
		idx := s.indexes[title]
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
