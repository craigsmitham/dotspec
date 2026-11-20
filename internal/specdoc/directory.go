package specdoc

import (
	"cmp"
	"fmt"
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func applicableRules(owner *subject) string {
	return "Applies to this " + strings.ToLower(owner.kind())
}

// An empty group is the collection's main list, directly beneath its title.
// Named groups retain relationships such as required versus provided interfaces.
type collectionIndex struct {
	file     string
	groups   map[string][]entry
	sources  map[string][]*heading
	passages [][]*heading
}

func (s *subject) index(title string) *collectionIndex {
	if s.indexes == nil {
		s.indexes = map[string]*collectionIndex{}
	}
	if s.indexes[title] == nil {
		s.indexes[title] = &collectionIndex{file: path.Join(s.dir, specfile.NameSlug(title)+".md"), groups: map[string][]entry{}, sources: map[string][]*heading{}}
	}
	return s.indexes[title]
}

func (idx *collectionIndex) contains(group, target string) bool {
	return idx != nil && slices.ContainsFunc(idx.groups[group], func(e entry) bool { return e.target == target })
}

func (idx *collectionIndex) add(group string, member entry) {
	if !idx.contains(group, member.target) {
		idx.groups[group] = append(idx.groups[group], member)
	}
}

func (p *planner) placeDirectoryEntry(s *subject) {
	title := sentenceCase(strings.ReplaceAll(families[s.kind()], "-", " "))
	group := ""
	if s.kind() == "Interface" {
		group = "Provided interfaces"
		if slices.Contains(model.Strings(s.parent.r["required_interfaces"]), id(s.r)) {
			group = "Required interfaces"
		}
	}
	s.parent.index(title).add(group, entry{s.name.title, s.file, subjectPreview(s)})
}

// An explicitly related member appears once, in the group that explains why it is here.
func (idx *collectionIndex) classify(group string, member entry) {
	idx.groups[""] = slices.DeleteFunc(idx.groups[""], func(e entry) bool { return e.target == member.target })
	if len(idx.groups[""]) == 0 {
		delete(idx.groups, "")
	}
	idx.add(group, member)
}

// An authored overview section keeps its narrative and points to the matching list.
func collectionGroup(title string) (string, string) {
	group := ""
	switch title {
	case "Child systems":
		title = "Systems"
	case "Required interfaces", "Provided interfaces":
		title, group = "Interfaces", title
	case "Applicable requirements":
		title, group = "Requirements", title
	}
	return title, group
}

func (s *subject) sectionIndex(title string) string {
	title, group := collectionGroup(title)
	idx := s.indexes[title]
	if idx == nil {
		return ""
	}
	if group != "" {
		if len(idx.groups[group]) == 0 {
			return ""
		}
		return idx.file + "#" + specfile.NameSlug(group)
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
				idx := p.governing.index("Open questions")
				owner := attr(p.nodes[nameNode(r)], "decoded_text")
				idx.groups[""] = append(idx.groups[""], entry{question, h.node, owner})
			}
		}
	}
}

// A collection heading whose children all have their own reports belongs to
// the directory. Preserve links to its authored fragment there.
func (p *planner) planPassages(s *subject) {
	var groups [][]*heading
	for _, h := range s.source.headings {
		if h.owner != s || h.inGlossary {
			continue
		}
		if len(groups) == 0 || h.level <= s.name.level+1 {
			groups = append(groups, []*heading{h})
		} else {
			i := len(groups) - 1
			groups[i] = append(groups[i], h)
		}
	}
	for _, group := range groups {
		h := group[0]
		if h != s.name && len(group) == 1 && h.subjectKind == "" && !needsHeadingAnchor(s, h) && strings.TrimSpace(sourceProse(s.source, h)) == "" {
			if dest := s.sectionIndex(h.title); dest != "" {
				index, part := collectionGroup(h.title)
				idx := s.indexes[index]
				idx.sources[part] = append(idx.sources[part], h)
				p.targets[h.node] = dest
				continue
			}
		}
		s.passages = append(s.passages, group)
	}
}

func (p *planner) requirementPreview(r model.Record) string {
	if s := p.pages[id(r)]; s != nil {
		return subjectPreview(s)
	}
	if page := p.ownerPage(r); page != nil {
		for _, h := range page.source.headings {
			if h.node == nameNode(r) {
				return paragraphPreview(sourceProse(page.source, h))
			}
		}
	}
	return ""
}

func sortEntries(entries []entry) {
	slices.SortFunc(entries, func(a, b entry) int {
		return cmp.Or(strings.Compare(a.title, b.title), strings.Compare(a.target, b.target))
	})
}

func (p *planner) renderIndex(owner *subject, title string, idx *collectionIndex) {
	w := p.newPage(idx.file, Document{}, title)
	w.text(w.link(owner.name.title, owner.file))
	if title == "Glossary" {
		p.renderGlossary(w, owner)
	}
	if title == "Quality attributes" {
		p.attributeOverview(w, owner)
	}
	for _, group := range slices.Sorted(maps.Keys(idx.groups)) {
		if group != "" {
			w.heading(group)
		}
		members := idx.groups[group]
		sortEntries(members)
		w.out.WriteByte('\n')
		for _, e := range members {
			fmt.Fprintln(&w.out, entryLine(w, e))
		}
		for _, h := range idx.sources[group] {
			w.text(link(h.sourceLabel(), relativeURL(w.file, owner.source.path, h.anchor)))
		}
	}
	if _, ok := idx.groups[""]; !ok {
		for _, h := range idx.sources[""] {
			w.text(link(h.sourceLabel(), relativeURL(w.file, owner.source.path, h.anchor)))
		}
	}
	for _, passage := range idx.passages {
		p.renderPassage(w, owner, passage, 1+owner.name.level-passage[0].level, false)
	}
	if len(idx.passages) > 0 {
		p.appendFootnotes(w, owner)
	}
	w.finish()
}

// Structural and applicability relationships are presented in scoped indexes.
func relationshipIndex(owner *subject, label, target string) *collectionIndex {
	switch label {
	case "Parts":
		for _, title := range []string{"Systems", "Components"} {
			if idx := owner.indexes[title]; idx.contains("", target) {
				return idx
			}
		}
	case "Provides":
		if idx := owner.indexes["Capabilities"]; idx.contains("Provided capabilities", target) {
			return idx
		}
	case "Applies here":
		if idx := owner.indexes["Business rules"]; idx.contains(applicableRules(owner), target) {
			return idx
		}
		if idx := owner.indexes["Requirements"]; idx.contains("Applicable requirements", target) {
			return idx
		}
	}
	return nil
}
