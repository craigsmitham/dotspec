package specfile

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// ReferenceDocument exposes only declared identities and authored Markdown links.
// It carries no inference about authority, applicability, or satisfaction.
type ReferenceDocument struct {
	Anchors map[string]bool
	Links   []Reference
}

type Reference struct {
	Destination string
	Location    diagnostic.Location
}

func (o *Outline) References() *ReferenceDocument {
	refs := &ReferenceDocument{Anchors: o.Anchors, Links: []Reference{}}
	outlineLinks(o, o.AST, refs)
	return refs
}

func outlineLinks(o *Outline, root ast.Node, refs *ReferenceDocument) {
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if link, ok := n.(*ast.Link); entering && ok {
			refs.Links = append(refs.Links, Reference{string(link.Destination), o.Location(link.Pos(), string(link.Destination))})
		}
		return ast.WalkContinue, nil
	})
}

// GlossarySourcesField identifies a system's explicitly included vocabulary files.
const GlossarySourcesField = "glossary_sources"

// GlossaryDetailReferences retains outgoing links in explicitly modeled passages.
// Ordinary citations in the rest of an imported document keep the normal one-hop scope.
func GlossaryDetailReferences(data []byte) (*ReferenceDocument, error) {
	outline, err := MarkdownOutline(data, "glossary")
	if err != nil {
		return nil, err
	}
	refs := &ReferenceDocument{Anchors: outline.Anchors}
	// Walk blocks in document order alongside the ordered detail sections.
	block := outline.AST.FirstChild()
	for _, term := range outline.Sections {
		for _, sub := range term.Children {
			switch sub.Title {
			case "Relationships", "Related terms", "Example", "Examples", "Illustrative examples":
				for block != nil && block.Pos() < sub.ContentStart {
					block = block.NextSibling()
				}
				for block != nil && block.Pos() < sub.End {
					outlineLinks(outline, block, refs)
					block = block.NextSibling()
				}
			}
		}
	}
	return refs, nil
}

// MarkdownReferences reads heading targets in native Markdown documents without
// imposing a dotspec schema, title, or layout. Explicit dotspec IDs come only
// from the owning specification parser.
func MarkdownReferences(data []byte) *ReferenceDocument {
	source := Normalize(data)
	body := source
	if strings.HasPrefix(string(body), "---\n") {
		if end := strings.Index(string(body[4:]), "\n---\n"); end >= 0 {
			body = body[4+end+5:]
		}
	}
	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	doc := md.Parser().Parse(text.NewReader(body))
	outline := &Outline{Original: source, Source: body, AST: doc, Anchors: map[string]bool{}}
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if heading, ok := n.(*ast.Heading); entering && ok {
			if id, ok := heading.AttributeString("id"); ok {
				outline.Anchors[string(id.([]byte))] = true
			}
		}
		return ast.WalkContinue, nil
	})
	return outline.References()
}
