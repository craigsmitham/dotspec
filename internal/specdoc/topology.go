package specdoc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/mermaid"
	"github.com/craigsmitham/dotspec/internal/model"
)

// Containment and boundary interactions have distinct meanings. Only explicit
// model relationships supply edges; directory membership is navigation.
func (p *planner) renderTopology(w *writer, owner *subject, shownCapabilities map[string]bool) {
	type connection struct{ target, edge, role string }
	var connections []connection
	if owner.parent != nil && attr(owner.r, "part_of") == id(owner.parent.r) {
		connections = append(connections, connection{id(owner.parent.r), "part of", "Parent " + strings.ToLower(owner.parent.kind())})
	}
	providers := map[string][]*subject{}
	for _, child := range p.subjects {
		if attr(child.r, "part_of") == id(owner.r) {
			connections = append(connections, connection{id(child.r), "contains", child.kind()})
		}
		if child.kind() == "Capability" {
			for _, provider := range model.Strings(child.r["provided_by"]) {
				providers[provider] = append(providers[provider], child)
			}
		}
	}
	for _, role := range []struct{ field, edge, label string }{
		{"required_interfaces", "requires interface", "Required interface"},
		{"provided_interfaces", "provides interface", "Provided interface"},
		{"interacts_with", "interacts with", "Boundary partner"},
	} {
		for _, target := range model.Strings(owner.r[role.field]) {
			connections = append(connections, connection{target, role.edge, role.label})
		}
	}
	nodes := map[string]int{id(owner.r): 0}
	var diagram strings.Builder
	diagram.WriteString("```mermaid\nflowchart LR\n")
	fmt.Fprintf(&diagram, "    n0[\"%s\"]\n", mermaid.Label(owner.name.title))
	var lines []string
	for _, c := range connections {
		name, detail := p.topologyDetail(w, c.target)
		if detail == "" {
			continue
		}
		n, seen := nodes[c.target]
		if !seen {
			n = len(nodes)
			nodes[c.target] = n
			fmt.Fprintf(&diagram, "    n%d[\"%s\"]\n", n, mermaid.Label(name))
		}
		fmt.Fprintf(&diagram, "    n0 -->|%s| n%d\n", c.edge, n)
		lines = append(lines, "- **"+c.role+"** · "+detail)
		if c.edge == "contains" && len(providers[c.target]) > 0 {
			for _, capability := range providers[c.target] {
				line := "  - Provides " + w.link(inlineText(capability.name.title), capability.file)
				if preview := subjectPreview(capability); preview != "" {
					line += " — " + inlineText(preview)
				}
				lines = append(lines, line)
				shownCapabilities[capability.file] = true
			}
		}
	}
	if len(lines) == 0 {
		return
	}
	diagram.WriteString("```")
	w.text(diagram.String())
	w.text(strings.Join(lines, "\n"))
}

// A resolved boundary partner in another project keeps its authored destination.
func (p *planner) topologyDetail(w *writer, target string) (string, string) {
	if s := p.pages[target]; s != nil {
		detail := w.link(inlineText(s.name.title), s.file)
		if summary := subjectPreview(s); summary != "" {
			detail += " — " + inlineText(summary)
		}
		return s.name.title, detail
	}
	r := p.elements[target]
	name := attr(p.nodes[nameNode(r)], "decoded_text")
	if name == "" {
		return "", ""
	}
	for _, source := range p.snapshot.Sources {
		if id(source) != attr(r, "source_document") {
			continue
		}
		rel, err := filepath.Rel(p.root, filepath.FromSlash(attr(source, "locator")))
		if err != nil {
			break
		}
		dest := relativeURL(w.file, filepath.ToSlash(rel), "")
		w.requireLink(dest)
		return name, link(inlineText(name), dest)
	}
	return "", ""
}
