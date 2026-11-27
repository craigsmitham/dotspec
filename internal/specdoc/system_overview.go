package specdoc

import (
	"maps"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/specfile"
)

// Each authored passage has one home in the overview. The same outline drives
// rendering, source content order, fragment allocation and the navigation map.
type overviewSection struct {
	group             string // optional label in the system section map
	collection        string // type of specification listed in this section
	title, anchor     string
	directory         *collectionIndex // complete index rendered in this overview
	level             int
	sources           []string
	passages          [][]*heading
	diagram           string
	hasBody           bool
	omitEmptyPassages bool
	companions        bool // lists the subject's frontmatter companions
}

func (h *heading) overviewField() string {
	field, _, _ := strings.Cut(h.field, ".")
	return strings.TrimSuffix(field, "_content")
}

func (p *planner) planOverview(s *subject) []overviewSection {
	architecture := p.architecture(s).sections()
	var sections []overviewSection
	if s.kind() == "System" {
		sections = systemOutline(architecture)
	} else {
		sections = componentOutline(architecture)
	}
	sections = append(sections, overviewSection{title: "Companions", level: 2, companions: true})
	sections = append(sections, overviewReferences...)
	omitEmpty := map[string]bool{}
	for _, section := range sections {
		if section.omitEmptyPassages {
			for _, field := range section.sources {
				omitEmpty[field] = true
			}
		}
	}
	groups := map[string][][]*heading{}
	for _, passage := range p.planPassages(s, omitEmpty) {
		h := passage[0]
		if h == s.name {
			sections[0].passages = append(sections[0].passages, passage)
		} else {
			field := h.overviewField()
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
		section.directory = s.indexes[section.collection]
		hasMembers := section.directory != nil && len(section.directory.members) > 0
		section.hasBody = len(section.passages) > 0 || section.diagram != "" || hasMembers || (section.companions && len(p.companions[id(s.r)]) > 0)
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
			renderIndexMembers(w, section.collection, section.directory.members)
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
		w.out.WriteString(section.diagram)
		if section.companions {
			w.text(p.companionBody(w.file, s))
		}
	}
}

func entryLine(w *writer, e entry) string {
	line := "- " + w.link(inlineText(e.title), e.target)
	if e.preview != "" {
		line += " — " + inlineText(e.preview)
	}
	if e.note != "" {
		line += " *(" + inlineText(e.note) + ")*"
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
		case section.group == "":
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

	titles := slices.Sorted(maps.Keys(s.indexes))
	var collections []string
	for _, title := range titles {
		idx := s.indexes[title]
		if !s.browseCollection(title, idx) {
			continue
		}
		line := w.link(title, idx.file)
		collections = append(collections, line)
	}
	w.disclosure("Browse specifications", collections)
}
