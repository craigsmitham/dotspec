package specfile

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func markdownOutline(data []byte) *Outline {
	source := Normalize(data)
	body := source
	if strings.HasPrefix(string(body), "---\n") {
		if end := strings.Index(string(body[4:]), "\n---\n"); end >= 0 {
			body = body[4+end+5:]
		}
	}
	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	return &Outline{Original: source, Source: body, AST: md.Parser().Parse(text.NewReader(body)), Anchors: map[string]bool{}}
}

// SpecificationSystem reads the explicit association without resolving it.
// Presence remains distinct from absence, including when the selection is invalid.
func SpecificationSystem(data []byte) (Reference, bool, error) {
	return markdownOutline(data).SpecificationSystem()
}

func (o *Outline) SpecificationSystem() (Reference, bool, error) {
	var location diagnostic.Location
	var links []Reference
	found, active := false, false
	for n := o.AST.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok {
			if h.Level <= 2 {
				active = false
			}
			if h.Level == 2 && strings.TrimSpace(inlineText(h, o.Source)) == "Specification system" {
				if found {
					return Reference{}, true, diagnostic.Error("structure.duplicate-section", "Duplicate Specification system section.", o.Location(h.Pos(), "Specification system"), location)
				}
				found, active = true, true
				location = o.Location(h.Pos(), "Specification system")
			}
		}
		if !active {
			continue
		}
		_ = ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			if link, ok := node.(*ast.Link); entering && ok {
				links = append(links, Reference{string(link.Destination), o.Location(link.Pos(), "Specification system")})
			}
			return ast.WalkContinue, nil
		})
	}
	if !found {
		return Reference{}, false, nil
	}
	if len(links) != 1 {
		return Reference{}, true, diagnostic.Error("relationship.specification_system", "Specification system must contain exactly one Markdown link to a system.", location)
	}
	return links[0], true, nil
}

// MarkdownTitleAnchor returns only the document title's generated anchor.
func MarkdownTitleAnchor(data []byte) string {
	o := markdownOutline(data)
	for n := o.AST.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
			if id, ok := h.AttributeString("id"); ok {
				return string(id.([]byte))
			}
		}
	}
	return ""
}
