package specdoc

import (
	"cmp"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Each target is a canonical subject page. Equal titles may name distinct subjects.
type collectionIndex struct {
	file     string
	overview *overviewSection
	members  []entry
	sources  []*heading
	passages [][]*heading
}

func (s *subject) index(title string) *collectionIndex {
	if s.indexes == nil {
		s.indexes = map[string]*collectionIndex{}
	}
	if s.indexes[title] == nil {
		s.indexes[title] = &collectionIndex{file: path.Join(s.dir, specfile.NameSlug(title)+".md")}
	}
	return s.indexes[title]
}

func (idx *collectionIndex) add(member entry) {
	if !slices.ContainsFunc(idx.members, func(e entry) bool { return e.target == member.target }) {
		idx.members = append(idx.members, member)
	}
}

func (s *subject) collection() string {
	if s.externalSystem {
		return s.collectionType() + " for other systems"
	}
	return s.collectionType()
}

func (s *subject) collectionType() string {
	if s.name.qualityAttribute {
		return "Quality attributes"
	}
	return sentenceCase(strings.ReplaceAll(families[s.kind()], "-", " "))
}

// Give the owner and its enclosing components the same canonical entry.
// A system (including a subsystem) is the outermost scope.
func (s *subject) includeMember(title string, member entry) {
	for owner := s; owner != nil; owner = owner.navigationParent {
		owner.index(title).add(member)
		if owner.kind() == "System" {
			break
		}
	}
}

func collectionName(title string) string {
	switch title {
	case "Child systems":
		return "Systems"
	case "Required interfaces", "Provided interfaces":
		return "Interfaces"
	case "Applicable requirements":
		return "Requirements"
	}
	return title
}

func (s *subject) sectionIndex(title string) string {
	idx := s.indexes[collectionName(title)]
	if idx == nil {
		return ""
	}
	if idx.overview != nil {
		return s.file + "#" + idx.overview.anchor
	}
	return idx.file
}

// headingAt returns the authored passage containing a content node.
func (sf *sourceFile) headingAt(pos int) *heading {
	i, exact := slices.BinarySearchFunc(sf.headings, pos, func(h *heading, pos int) int {
		return cmp.Compare(h.start, pos)
	})
	if !exact {
		i--
	}
	if i < 0 {
		return nil
	}
	return sf.headings[i]
}

// Each mapped question supplies its wording and owning subject. Multiple list
// items can share the same section link and must remain separate entries.
func (p *planner) indexOpenQuestions() {
	for _, r := range p.snapshot.Elements {
		src := p.sources[attr(r, "source_document")]
		if src == nil {
			continue
		}
		for _, q := range model.OpenQuestions(r) {
			nodes := model.Strings(model.Object(q["question"])["content_nodes"])
			if len(nodes) == 0 {
				continue
			}
			start := position(p.nodes[nodes[0]], "start")
			h := src.headingAt(start)
			if h == nil || h.owner == nil {
				continue
			}
			end := position(p.nodes[nodes[len(nodes)-1]], "end")
			question := contentText(src.passageMarkdown(start, end))
			if question != "" {
				system := h.owner
				if system.kind() != "System" {
					system = system.system
				}
				if system == nil {
					continue
				}
				for parent := system.navigationParent; parent != nil && parent.kind() == "System"; parent = parent.navigationParent {
					system = parent
				}
				idx := system.index("Open questions")
				owner := attr(p.nodes[nameNode(r)], "decoded_text")
				idx.members = append(idx.members, entry{question, h.node, owner})
			}
		}
	}
}

// A collection heading whose children all have their own reports belongs to
// the directory. Route citations to its authored fragment there.
func (p *planner) planPassages(s *subject) {
	var groups [][]*heading
	for _, h := range s.source.headings {
		if h.owner != s || h.inGlossary {
			continue
		}
		if len(groups) == 0 || (h.level <= s.name.level+1 && (h.field != "" || h.subjectKind != "")) {
			groups = append(groups, []*heading{h})
		} else {
			i := len(groups) - 1
			groups[i] = append(groups[i], h)
		}
	}
	for _, group := range groups {
		h := group[0]
		field, _, _ := strings.Cut(h.field, ".")
		field = strings.TrimSuffix(field, "_content")
		omitEmpty := s.kind() == "System" && slices.Contains(systemContextFields, field)
		if s.isOverview() && (field == "architecture" || field == "part_of" || field == "interacts_with") {
			omitEmpty = true
		}
		if omitEmpty && !slices.ContainsFunc(group, func(part *heading) bool {
			return part.subjectKind != "" || needsHeadingAnchor(s, part) || strings.TrimSpace(s.source.data[part.end:part.bodyEnd]) != ""
		}) {
			for _, part := range group {
				p.keepEmptyHeading(s, part)
			}
			continue
		}
		if h != s.name && len(group) == 1 && h.subjectKind == "" && !needsHeadingAnchor(s, h) && strings.TrimSpace(sourceProse(s.source, h)) == "" {
			if dest := s.sectionIndex(h.title); dest != "" {
				index := collectionName(h.title)
				idx := s.indexes[index]
				idx.sources = append(idx.sources, h)
				p.targets[h.node] = dest
				continue
			}
			if specfile.IsCollectionHeading(h.title) || h.title == "Open questions" {
				p.keepEmptyHeading(s, h)
				continue
			}
		}
		s.passages = append(s.passages, group)
	}
}

// Authored links can still reach an empty heading omitted from the overview.
func (p *planner) keepEmptyHeading(s *subject, h *heading) {
	anchor := "empty." + h.anchor
	s.emptyHeadingAnchors = append(s.emptyHeadingAnchors, anchor)
	p.targets[h.node] = s.file + "#" + anchor
}

func sortEntries(entries []entry) {
	slices.SortFunc(entries, func(a, b entry) int {
		return cmp.Or(strings.Compare(strings.ToLower(a.title), strings.ToLower(b.title)), strings.Compare(a.title, b.title), strings.Compare(a.target, b.target))
	})
}

func collectionTitle(title string) string {
	if title == "Concepts" {
		return "Concept definitions"
	}
	return title
}

func (p *planner) renderIndex(owner *subject, title string, idx *collectionIndex) {
	w := p.newPage(idx.file, Document{}, collectionTitle(title))
	w.setBreadcrumb(owner, nil)
	if title == "Glossary" {
		p.renderGlossary(w, owner)
	}
	renderIndexMembers(w, title, idx.members)
	for _, passage := range idx.passages {
		p.renderPassage(w, owner, passage, 1+owner.name.level-passage[0].level, false)
	}
	if len(idx.passages) > 0 {
		p.appendFootnotes(w, owner)
	}
	w.finish()
}

func renderIndexMembers(w *writer, title string, members []entry) {
	if len(members) == 0 {
		return
	}
	sortEntries(members)
	w.out.WriteByte('\n')
	for _, e := range members {
		e.preview = directoryPreview(title, e.preview)
		fmt.Fprintln(&w.out, entryLine(w, e))
	}
}

func directoryPreview(title, preview string) string {
	if title == "Open questions" {
		return preview
	}
	return indexPreview(preview)
}

// An overview already has a route to these members through its typed directory.
// The member's own report explains its relationship to the owner.
func (s *subject) relationshipDirectory(label, target string) string {
	if !s.isOverview() {
		return ""
	}
	switch label {
	case "Parts", "Provides", "Applies here", "Provided interfaces", "Required interfaces":
		for title, idx := range s.indexes {
			if slices.ContainsFunc(idx.members, func(member entry) bool { return member.target == target }) {
				return s.sectionIndex(title)
			}
		}
	}
	return ""
}
