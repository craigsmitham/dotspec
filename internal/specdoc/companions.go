package specdoc

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// A companion as presented: its authored target is relative to the declaring
// specification's source file. line is the declaring line of a named native
// definition, as the snapshot resolved it, or zero.
type companionView struct {
	role, note, target, source string
	line                       int
}

// readCompanions collects each subject's frontmatter companions in authored order.
func (p *planner) readCompanions() {
	refs := map[string]model.Record{}
	for _, r := range p.snapshot.References {
		refs[id(r)] = r
	}
	p.companions = map[string][]companionView{}
	for _, s := range p.subjects {
		for _, cid := range model.Strings(s.r["companions"]) {
			c := p.elements[cid]
			if c == nil {
				continue
			}
			ref := refs[attr(c, "target")]
			p.companions[id(s.r)] = append(p.companions[id(s.r)], companionView{
				role:   attr(c, "role"),
				note:   attr(c, "note"),
				target: attr(ref, "authored_target"),
				source: s.source.path,
				line:   model.Number(model.Object(ref["resolution"])["line"]),
			})
		}
	}
}

// companionLink names a file companion by its base name and definition. A
// named definition in a supported native format links to its declaring line,
// which code hosts can scroll to.
func (p *planner) companionLink(page string, c companionView) string {
	u, err := url.Parse(c.target)
	if err != nil || u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") || u.Path == "" {
		return link(c.target, c.target)
	}
	title := path.Base(u.Path)
	raw := c.target
	if u.Fragment != "" {
		title += " › " + u.Fragment
		if c.line > 0 {
			raw = (&url.URL{Path: u.Path, Fragment: "L" + strconv.Itoa(c.line)}).String()
		}
	}
	return link(title, p.documentLink(page, c.source, raw))
}

func companionLine(p *planner, page string, c companionView) string {
	line := "- " + p.companionLink(page, c)
	if c.note != "" {
		line += " — " + inlineText(c.note)
	}
	return line
}

// companionBody groups a subject's companions by role in vocabulary order.
func (p *planner) companionBody(page string, s *subject) string {
	views := p.companions[id(s.r)]
	var b strings.Builder
	for _, role := range specfile.CompanionRoles {
		var lines []string
		for _, c := range views {
			if c.role == role {
				lines = append(lines, companionLine(p, page, c))
			}
		}
		if len(lines) > 0 {
			b.WriteString("\n**" + specfile.CompanionRoleLabel(role) + "**\n\n" + strings.Join(lines, "\n") + "\n")
		}
	}
	return b.String()
}

func (p *planner) renderCompanions(w *writer, s *subject) {
	if body := p.companionBody(w.file, s); body != "" {
		w.heading("Companions")
		w.text(body)
	}
}

// The outermost system indexes every companion declared within it, so its
// section map leads to the whole intent layer, not only Markdown specifications.
func (p *planner) companionSystem(s *subject) *subject {
	system := s
	if system.kind() != "System" {
		system = system.system
	}
	for system != nil && system.navigationParent != nil && system.navigationParent.kind() == "System" {
		system = system.navigationParent
	}
	return system
}

func (p *planner) indexCompanions() {
	for _, s := range p.subjects {
		if len(p.companions[id(s.r)]) == 0 {
			continue
		}
		if system := p.companionSystem(s); system != nil {
			system.index("Companions")
		}
	}
}

func (p *planner) renderCompanionIndex(w *writer, owner *subject) {
	for _, role := range specfile.CompanionRoles {
		var lines []string
		for _, s := range p.subjects {
			if p.companionSystem(s) != owner {
				continue
			}
			for _, c := range p.companions[id(s.r)] {
				if c.role == role {
					lines = append(lines, companionLine(p, w.file, c)+" · "+w.link(inlineText(s.name.title), s.file))
				}
			}
		}
		if len(lines) > 0 {
			w.heading(specfile.CompanionRoleLabel(role))
			w.text(strings.Join(lines, "\n"))
		}
	}
}
