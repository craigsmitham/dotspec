package specdoc

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// A preview selects a complete authored paragraph. The destination retains its
// formatting, links, subsequent paragraphs and qualifications.
func paragraphPreview(md string) string {
	source := []byte(md)
	doc := markdownParser.Parse(text.NewReader(source))
	for block := doc.FirstChild(); block != nil; block = block.NextSibling() {
		if block.Kind() != ast.KindParagraph {
			continue
		}
		var lines []string
		for i := 0; i < block.Lines().Len(); i++ {
			segment := block.Lines().At(i)
			line := strings.TrimSpace(string(segment.Value(source)))
			if attributeField(line) {
				continue
			}
			lines = append(lines, line)
		}
		if len(lines) > 0 {
			return contentText(strings.Join(lines, "\n"))
		}
	}
	return ""
}

func attributeField(line string) bool {
	return strings.HasPrefix(line, "**ID:**") || strings.HasPrefix(line, "**Sub-attribute of:**")
}

// Previews are visible text, including when embedded in a Markdown table.
var inlineEscapes = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;",
	"\\", "&#92;", "|", "&#124;", "[", "&#91;", "]", "&#93;",
	"*", "&#42;", "_", "&#95;", "`", "&#96;", "~", "&#126;")

func inlineText(s string) string {
	return inlineEscapes.Replace(s)
}

func (p *planner) attributeOverview(w *writer, system *subject) {
	var roots []*subject
	children := map[string][]*subject{}
	for _, id := range model.Strings(system.r["quality_attributes"]) {
		a := p.pages[id]
		if a == nil {
			continue
		}
		parents := model.Strings(a.r["subattribute_of"])
		if len(parents) == 0 {
			roots = append(roots, a)
		}
		for _, parent := range parents {
			children[parent] = append(children[parent], a)
		}
	}

	var renderChildren func(parent string, indent string)
	renderChildren = func(parent string, indent string) {
		for _, child := range children[parent] {
			fmt.Fprintf(&w.out, "%s- %s\n", indent, attributeEntry(w, child, false))
			// Show every declared grouping, but expand a shared subtree only
			// under the first parent, where its report also lives.
			if model.Strings(child.r["subattribute_of"])[0] == parent {
				renderChildren(id(child.r), indent+"  ")
			}
		}
	}
	for _, root := range roots {
		w.text(attributeEntry(w, root, true))
		if key := id(root.r); len(children[key]) > 0 {
			w.out.WriteByte('\n')
			renderChildren(key, "")
		}
	}
}

func attributeEntry(w *writer, a *subject, root bool) string {
	name := w.link(inlineText(a.name.title), a.file)
	if root {
		name = "**" + name + "**"
	}
	heading := a.section("Definition")
	if heading == nil {
		heading = a.name
	}
	if preview := paragraphPreview(sourceProse(a.source, heading)); preview != "" {
		return name + " — " + inlineText(preview)
	}
	return name
}

func (s *subject) section(title string) *heading {
	for _, h := range s.source.headings {
		if h.owner == s && h.level == s.name.level+1 && h.title == title {
			return h
		}
	}
	return nil
}

func hasDescriptionSection(kind string) bool {
	switch kind {
	case "System", "Component", "Feature", "Capability", "ConceptUse", "QualityModel":
		return true
	}
	return false
}

func subjectPreview(s *subject) string {
	if hasDescriptionSection(s.kind()) {
		if h := s.section("Description"); h != nil {
			return paragraphPreview(sourceProse(s.source, h))
		}
		return ""
	}
	title := map[string]string{
		"Concept": "Purpose", "UseCase": "Goal", "Job": "Job",
		"Requirement": "Statement", "BusinessRule": "Statement",
	}[s.kind()]
	h := s.name
	if title != "" {
		h = s.section(title)
		if h == nil {
			return ""
		}
	}
	if s.kind() == "Interface" {
		return paragraphPreview(sourceProse(s.source, h))
	}
	// A bare metadata field or cross-reference gives little context in a list.
	raw := sourceProse(s.source, h)
	if title == "" && strings.HasPrefix(raw, "**") && strings.Contains(strings.SplitN(raw, "\n", 2)[0], ":**") {
		raw = ""
	}
	summary := paragraphPreview(raw)
	if s.kind() == "Synchronization" && summary == "" {
		if h := s.section("Semantics"); h != nil {
			summary = paragraphPreview(sourceProse(s.source, h))
		}
		if summary == "" {
			var clauses []string
			for _, field := range []string{"When", "Where", "Then"} {
				if h := s.section(field); h != nil {
					if preview := paragraphPreview(sourceProse(s.source, h)); preview != "" {
						clauses = append(clauses, field+" — "+preview)
					}
				}
			}
			summary = strings.Join(clauses, " · ")
		}
	}
	if summary == s.name.title {
		return ""
	}
	return summary
}

func (p *planner) featureGoals(w *writer, s *subject) {
	if paragraphPreview(sourceProse(s.source, s.name)) != "" || s.section("Purpose") != nil || s.section("Description") != nil {
		return
	}
	var lines []string
	for _, id := range model.Strings(s.r["groups"]) {
		member := p.pages[id]
		if member == nil || member.kind() != "UseCase" {
			continue
		}
		if goal := subjectPreview(member); goal != "" {
			lines = append(lines, entryLine(w, entry{member.name.title, p.targets[member.section("Goal").node], goal}))
		}
	}
	w.text(strings.Join(lines, "\n"))
}
