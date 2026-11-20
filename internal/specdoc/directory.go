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

const applicableRules = "Applies to this system"

// An empty group is the collection's main list, directly beneath its title.
// Named groups retain relationships such as required versus provided interfaces.
type collectionIndex struct {
	file    string
	groups  map[string][]entry
	sources map[string][]*heading
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
		slices.SortFunc(members, func(a, b entry) int {
			return cmp.Or(strings.Compare(a.title, b.title), strings.Compare(a.target, b.target))
		})
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
	case "Applies here":
		if idx := owner.indexes["Business rules"]; idx.contains(applicableRules, target) {
			return idx
		}
		if idx := owner.indexes["Requirements"]; idx.contains("Applicable requirements", target) {
			return idx
		}
	}
	return nil
}
