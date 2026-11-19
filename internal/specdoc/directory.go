package specdoc

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
)

type directory struct {
	section *heading // reuse the authored section when one exists
	anchor  string
	entries []entry
}

func (s *subject) addMember(label string, member entry) {
	if s.directories == nil {
		s.directories = map[string]*directory{}
	}
	d := s.directories[label]
	if d == nil {
		d = &directory{}
		s.directories[label] = d
	}
	if !d.contains(member.target) {
		d.entries = append(d.entries, member)
	}
}

func (d *directory) contains(target string) bool {
	return d != nil && slices.ContainsFunc(d.entries, func(e entry) bool { return e.target == target })
}

func (p *planner) placeDirectoryEntry(s *subject) {
	owner := s.system
	label := sentenceCase(strings.ReplaceAll(families[s.kind()], "-", " "))
	switch s.kind() {
	case "System":
		label = "Child systems"
	case "Component":
		if parent := p.pages[attr(s.r, "part_of")]; parent != nil {
			owner = parent
		}
	case "Interface":
		if declared := p.pages[attr(s.r, "owner")]; declared != nil {
			owner = declared
		}
		label = "Provided interfaces"
		if slices.Contains(model.Strings(owner.r["required_interfaces"]), id(s.r)) {
			label = "Required interfaces"
		}
	default:
		// Some declarations (including concept uses and synchronizations) can
		// belong to a component rather than directly to its system.
		if declaring := p.pages[attr(s.r, "owner")]; declaring != nil && declaring.kind() == "Component" {
			owner = declaring
		}
	}
	owner.addMember(label, entry{s.name.title, s.file, s})
}

func (d *directory) render(w *writer) {
	slices.SortFunc(d.entries, func(a, b entry) int {
		return cmp.Or(strings.Compare(a.title, b.title), strings.Compare(a.target, b.target))
	})
	w.out.WriteByte('\n')
	for _, e := range d.entries {
		line := "- " + w.link(e.title, e.target)
		if len(d.entries) > 1 && e.subject != nil {
			if summary := subjectPreview(e.subject); summary != "" {
				line += " — " + inlineText(summary)
			}
		}
		fmt.Fprintln(&w.out, line)
	}
}
