package specdoc

import (
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

// A directory excerpt ends at a word boundary. Mark omitted wording explicitly.
func indexPreview(preview string) string {
	const limit = 36
	words := strings.Fields(preview)
	if len(words) <= 24 {
		return preview
	}
	end := min(limit, len(words))
	for i, word := range words[:end] {
		if strings.HasSuffix(word, ".") || strings.HasSuffix(word, "?") || strings.HasSuffix(word, "!") {
			end = i + 1
			break
		}
	}
	if end == len(words) {
		return preview
	}
	return strings.Join(words[:end], " ") + " …"
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
	if s.kind() == "SystemQualityAttribute" {
		h := s.section("Definition")
		if h == nil {
			h = s.name
		}
		return paragraphPreview(sourceProse(s.source, h))
	}
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

func (p *planner) featureGoals(s *subject) []entry {
	if paragraphPreview(sourceProse(s.source, s.name)) != "" || s.section("Purpose") != nil || s.section("Description") != nil {
		return nil
	}
	var idx collectionIndex
	for _, id := range model.Strings(s.r["groups"]) {
		member := p.pages[id]
		if member == nil || member.kind() != "UseCase" {
			continue
		}
		idx.add(entry{member.name.title, member.file, subjectPreview(member)})
	}
	return idx.members
}
