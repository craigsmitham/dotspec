package specdoc

import (
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Each authored passage has one home in the overview. The same outline drives
// rendering, source content order, fragment allocation and the navigation map.
type overviewSection struct {
	group         string // optional label in the system section map
	collection    string // type of specification listed in this section
	title, anchor string
	directory     *collectionIndex // complete index rendered in this overview
	level         int
	sources       []string
	passages      [][]*heading
	detail        string
	detailLinks   []string
	hasBody       bool
}

// These System fields each have their own section in the report, in reading order.
var systemContextFields = []string{"scope", "out_of_scope", "environment", "assumptions", "constraints"}

func (p *planner) overviewOutline(s *subject) []overviewSection {
	sections := []overviewSection{
		{sources: []string{"description"}},
		{group: "Intent", title: "Framing", level: 2, sources: []string{"framing"}},
		{group: "Intent", title: "Purpose", level: 2, sources: []string{"purpose"}},
		{group: "Intent", title: "Goals", level: 2, sources: []string{"goals"}},
		{group: "Intent", title: "Stakeholders", level: 2, sources: []string{"stakeholder_roles"}},
	}
	features := overviewSection{group: "Behavior", title: "Features", collection: "Features", level: 2}
	capabilities := overviewSection{group: "Behavior", title: "Capabilities", collection: "Capabilities", level: 2, sources: []string{"capabilities"}}
	quality := overviewSection{title: "Quality Attributes", collection: "Quality attributes", level: 2}
	architecture := p.architecture(s).sections()
	if s.kind() == "System" {
		sections = append(sections, overviewSection{group: "Intent", title: "Jobs", collection: "Jobs", level: 2})
		for _, field := range systemContextFields {
			title := sentenceCase(strings.ReplaceAll(field, "_", " "))
			level := 2
			if field == "out_of_scope" {
				level = 3
			}
			sections = append(sections, overviewSection{group: "Context", title: title, level: level, sources: []string{field}})
		}
		sections = append(sections,
			overviewSection{group: "Context", title: "Business rules", collection: "Business rules", level: 2},
			features, capabilities,
			overviewSection{group: "Behavior", title: "Use cases", collection: "Use cases", level: 2},
			quality, overviewSection{title: "Requirements", collection: "Requirements", level: 2})
	} else {
		sections = append(sections,
			overviewSection{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
			overviewSection{title: "Scope", level: 2, sources: []string{"boundary_description", "out_of_scope"}})
		sections = append(sections, architecture...)
		capabilities.level = 3
		sections = append(sections, capabilities, features)
	}
	sections = append(sections,
		overviewSection{group: "Design", title: "Concepts", collection: "Concepts", level: 2},
		overviewSection{group: "Design", title: "Concept uses", collection: "Concept uses", level: 2, sources: []string{"concept_uses"}},
		overviewSection{group: "Design", title: "Synchronizations", collection: "Synchronizations", level: 2, sources: []string{"synchronizations"}})
	if s.kind() == "System" {
		architecture[0].group = "Design"
		sections = append(sections, architecture...)
	} else {
		sections = append(sections, quality)
	}
	return append(sections, []overviewSection{
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
		section.directory = s.indexes[section.collection]
		if section.directory != nil {
			renderIndexMembers(detail, section.collection, section.directory.members)
		}
		section.detail += detail.out.String()
		section.detailLinks = append(section.detailLinks, detail.page.RequiredLinks...)
		section.hasBody = len(section.passages) > 0 || section.detail != ""
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
			if explicit {
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

func (p *planner) planCollectionHomes(s *subject) {
	for i := range s.overview {
		section := &s.overview[i]
		idx := section.directory
		if idx == nil {
			continue
		}
		idx.overview = section
		for _, h := range idx.sources {
			p.targets[h.node] = s.file + "#" + section.anchor
		}
	}
}

func (s *subject) browseCollection(title string, idx *collectionIndex) bool {
	if idx.overview != nil {
		return false
	}
	if s.kind() != "System" {
		return true
	}
	if isDetailCollection(title) && len(idx.passages) > 0 && len(idx.members) == 0 {
		return false
	}
	return title != "Quality models"
}

func (p *planner) renderOverview(w *writer, s *subject) {
	for _, title := range slices.Sorted(maps.Keys(s.indexes)) {
		if idx := s.indexes[title]; idx.overview == nil {
			p.renderIndex(s, title, idx)
		}
	}
	for _, section := range s.overview {
		if section.level > 0 {
			if s.name.explicitAnchor != "" {
				w.text(`<a id="` + section.anchor + `"></a>`)
			}
			title := section.title
			w.text(strings.Repeat("#", section.level) + " " + title)
		}
		if section.directory != nil {
			w.out.WriteString(section.detail)
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
		if section.directory == nil {
			w.out.WriteString(section.detail)
		}
		for _, dest := range section.detailLinks {
			w.requireLink(dest)
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
	var previousGroup string
	for _, section := range s.overview {
		if section.level != 2 {
			continue
		}
		link := w.link(section.title, s.file+"#"+section.anchor)
		switch {
		case s.kind() == "Component" || section.group == "":
			sections = append(sections, link)
		case section.group != previousGroup:
			sections = append(sections, "**"+section.group+"**\\\n"+link)
		default:
			sections[len(sections)-1] += " · " + link
		}
		previousGroup = section.group
	}
	if s.kind() == "System" {
		w.disclosureBody("On this page", strings.Join(sections, "\n\n"))
	} else {
		w.disclosure("On this page", sections)
	}

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
		collections = append(collections, line)
	}
	w.disclosure("Browse specifications", collections)
}
