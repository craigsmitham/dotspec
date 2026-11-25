package specdoc

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/mermaid"
	"github.com/craigsmitham/dotspec/internal/model"
)

// Each view opens one declared boundary. The file tree only determines where
// the reader can find a report, never which subjects belong inside the diagram.
type architectureView struct {
	owner                  *subject
	subsystems, components []*subject
	provided, required     []*subject
	partners               []string
	capabilities           map[string][]*subject
}

func (p *planner) architecture(owner *subject) architectureView {
	v := architectureView{owner: owner, capabilities: map[string][]*subject{}}
	for _, s := range p.subjects {
		if attr(s.r, "part_of") == id(owner.r) {
			switch s.kind() {
			case "System":
				v.subsystems = append(v.subsystems, s)
			case "Component":
				v.components = append(v.components, s)
			}
		}
		if s.kind() == "Capability" {
			for _, provider := range model.Strings(s.r["provided_by"]) {
				v.capabilities[provider] = append(v.capabilities[provider], s)
			}
		}
	}
	interfaces := func(field string) []*subject {
		var entries []*subject
		for _, target := range model.Strings(owner.r[field]) {
			if s := p.pages[target]; s != nil {
				entries = append(entries, s)
			}
		}
		return entries
	}
	v.provided = interfaces("provided_interfaces")
	v.required = interfaces("required_interfaces")
	for _, target := range model.Strings(owner.r["interacts_with"]) {
		// A declared interaction with a constituent is preserved in the text.
		// Only partners outside the structural boundary go outside the box.
		if p.insideBoundary(target, id(owner.r)) {
			continue
		}
		if name := attr(p.nodes[nameNode(p.elements[target])], "decoded_text"); name != "" {
			v.partners = append(v.partners, name)
		}
	}
	return v
}

func (p *planner) insideBoundary(target, owner string) bool {
	seen := map[string]bool{}
	for target != "" && !seen[target] {
		if target == owner {
			return true
		}
		seen[target] = true
		target = attr(p.elements[target], "part_of")
	}
	return false
}

func (v architectureView) diagram(structure, boundary bool) string {
	var inside, outside strings.Builder
	node := 0
	add := func(parts []*subject, role string, port bool) {
		for _, s := range parts {
			node++
			label := mermaid.Label(role + ": " + s.name.title)
			if port {
				fmt.Fprintf(&inside, "        n%d([\"%s\"])\n", node, label)
			} else {
				fmt.Fprintf(&inside, "        n%d[\"%s\"]\n", node, label)
			}
		}
	}
	if structure {
		add(v.subsystems, "Subsystem", false)
		add(v.components, "Component", false)
	}
	if boundary {
		add(v.provided, "Provided interface", true)
		add(v.required, "Required interface", true)
		for _, name := range v.partners {
			node++
			fmt.Fprintf(&outside, "    owner ---|interacts with| n%d[\"%s\"]\n", node, mermaid.Label(name))
		}
	}
	if node == 0 {
		return ""
	}
	label := mermaid.Label(v.owner.name.title)
	var diagram strings.Builder
	diagram.WriteString("\n```mermaid\nflowchart LR\n")
	if inside.Len() > 0 {
		fmt.Fprintf(&diagram, "    subgraph owner[\"%s\"]\n%s    end\n", label, inside.String())
	} else {
		fmt.Fprintf(&diagram, "    owner[\"%s\"]\n", label)
	}
	diagram.WriteString(outside.String())
	diagram.WriteString("```\n")
	return diagram.String()
}

func architectureLink(w *writer, s *subject) string {
	preview := subjectPreview(s)
	if s.kind() == "System" || (s.kind() == "Component" && preview == "") {
		field := "Purpose"
		if s.kind() == "Component" {
			field = "Responsibilities"
		}
		if h := s.section(field); h != nil {
			if summary := paragraphPreview(sourceProse(s.source, h)); summary != "" {
				preview = summary
			}
		}
	}
	name := w.link(inlineText(s.name.title), s.file)
	if preview != "" {
		return name + " — " + inlineText(indexPreview(preview))
	}
	return name
}

func (v architectureView) sections() []overviewSection {
	sections := []overviewSection{
		{title: "Architecture", level: 2, sources: []string{"architecture"}},
		{title: "Part of", level: 3, sources: []string{"part_of"}},
	}
	hasParts := len(v.components)+len(v.subsystems) > 0
	kind := v.owner.kind()
	if hasParts && len(v.partners) > 0 {
		sections = append(sections,
			overviewSection{title: kind + " structure", level: 3, detail: v.diagram(true, false)},
			overviewSection{title: kind + " context", level: 3, detail: v.diagram(false, true)})
	} else if diagram := v.diagram(true, true); diagram != "" {
		// A caption keeps the diagram outside any preceding authored subsection.
		if v.owner.section("Part of") != nil || v.owner.section("Architecture") != nil {
			title := kind + " context"
			if hasParts {
				title = kind + " structure"
			}
			sections = append(sections, overviewSection{title: title, level: 3, detail: diagram})
		} else {
			sections[0].detail = diagram
		}
	}

	for _, group := range []struct {
		title string
		items []*subject
	}{
		{"Subsystems", v.subsystems},
		{"Components", v.components},
		{"Interfaces", nil},
		{"Interactions", nil},
		{"Provided capabilities", v.capabilities[id(v.owner.r)]},
	} {
		section := overviewSection{title: group.title, level: 3}
		w := &writer{file: v.owner.file, page: &Page{}}
		var lines []string
		switch group.title {
		case "Interactions":
			section.sources = []string{"interacts_with"}
		case "Interfaces":
			section.sources = []string{"provided_interfaces", "required_interfaces"}
			for _, s := range v.provided {
				lines = append(lines, "- **Provided** · "+architectureLink(w, s))
			}
			for _, s := range v.required {
				lines = append(lines, "- **Required** · "+architectureLink(w, s))
			}
		default:
			for _, s := range group.items {
				lines = append(lines, "- "+architectureLink(w, s))
				for _, capability := range v.capabilities[id(s.r)] {
					lines = append(lines, "  - Provides "+architectureLink(w, capability))
				}
			}
		}
		w.text(strings.Join(lines, "\n"))
		section.detail = w.out.String()
		section.detailLinks = w.page.RequiredLinks
		sections = append(sections, section)
	}
	return sections
}
