package specfile

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
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
	// Role is the author-stated role of a frontmatter companion; authored body
	// links have none.
	Role string
}

func (o *Outline) References() *ReferenceDocument {
	refs := &ReferenceDocument{Anchors: o.Anchors, Links: []Reference{}}
	for _, c := range o.Companions {
		refs.Links = append(refs.Links, Reference{Destination: c.Target, Location: c.Location, Role: c.Role})
	}
	outlineLinks(o, o.AST, refs)
	return refs
}

func outlineLinks(o *Outline, root ast.Node, refs *ReferenceDocument) {
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if link, ok := n.(*ast.Link); entering && ok {
			refs.Links = append(refs.Links, Reference{Destination: string(link.Destination), Location: o.Location(link.Pos(), string(link.Destination))})
		}
		return ast.WalkContinue, nil
	})
}

// GlossarySourcesField identifies a system's explicitly included vocabulary files.
const GlossarySourcesField = "glossary_sources"

// GlossaryDetailReferences retains outgoing links in explicitly modeled passages.
// Ordinary citations in the rest of an imported document keep the normal one-hop scope.
func (md *Markdown) GlossaryDetailReferences() (*ReferenceDocument, error) {
	outline, err := md.Outline("glossary")
	if err != nil {
		return nil, err
	}
	refs := &ReferenceDocument{Anchors: outline.Anchors}
	// Walk blocks in document order alongside the ordered detail sections.
	block := outline.AST.FirstChild()
	for _, term := range outline.Sections {
		for _, sub := range term.Children {
			switch sub.Title {
			case "Relationships", "Related terms":
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

var anchorAttribute = regexp.MustCompile(`(?i)^<a\s+(?:[^>]*?\s)?(?:id|name)=["']([^"']+)["']`)

// HTMLAnchor reads an explicit fragment from an HTML anchor node in Markdown.
func HTMLAnchor(n ast.Node, source []byte) string {
	var raw []byte
	switch node := n.(type) {
	case *ast.RawHTML:
		raw = node.Segments.Value(source)
	case *ast.HTMLBlock:
		raw = node.Lines().Value(source)
	default:
		return ""
	}
	if match := anchorAttribute.FindSubmatch(bytes.TrimSpace(raw)); match != nil {
		return string(match[1])
	}
	return ""
}

// RenderedHeadingIDs allocates GitHub-style fragments from visible heading text.
// Each document supplies its headings in rendering order, including repeated names.
type RenderedHeadingIDs map[string]int

func (ids RenderedHeadingIDs) Next(title string) string {
	base := strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r), unicode.IsMark(r), unicode.Is(unicode.Pc, r), r == '-':
			return unicode.ToLower(r)
		case r == ' ':
			return '-'
		}
		return -1
	}, strings.TrimSpace(title))
	id := base
	for ids[id] > 0 {
		id = base + "-" + strconv.Itoa(ids[base])
		ids[base]++
	}
	ids[id] = 1
	return id
}

// RenderedText extracts the visible text of a Markdown node, resolving escapes
// and character references outside of code spans.
func RenderedText(node ast.Node, source []byte) string {
	var out strings.Builder
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Type() == ast.TypeBlock {
			out.WriteByte(' ')
		}
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Text:
			value := n.Value(source)
			if !n.IsRaw() && n.Parent().Kind() != ast.KindCodeSpan {
				value = util.ResolveEntityNames(util.ResolveNumericReferences(util.UnescapePunctuations(value)))
			}
			out.Write(value)
			if n.SoftLineBreak() || n.HardLineBreak() {
				out.WriteByte(' ')
			}
		case *ast.String:
			out.Write(n.Value)
		case *ast.AutoLink:
			out.Write(n.Label(source))
		case *ast.CodeBlock, *ast.FencedCodeBlock:
			out.Write(n.Lines().Value(source))
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(out.String())
}

// MarkdownReferences reads heading and HTML anchor targets in native Markdown documents without
// imposing a dotspec schema, title, or layout. Explicit dotspec IDs come only
// from the owning specification parser.
func MarkdownReferences(data []byte) *ReferenceDocument {
	md := NewMarkdown(data)
	refs := md.AnchorReferences()
	outlineLinks(md.outline(), md.Tree, refs)
	return refs
}
func (md *Markdown) AnchorReferences() *ReferenceDocument {
	outline := md.outline()
	doc := outline.AST
	ids := RenderedHeadingIDs{}
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if anchor := HTMLAnchor(n, md.Body); anchor != "" {
			outline.Anchors[anchor] = true
		}
		if heading, ok := n.(*ast.Heading); ok {
			outline.Anchors[ids.Next(RenderedText(n, md.Body))] = true
			if id, ok := heading.AttributeString("id"); ok {
				outline.Anchors[string(id.([]byte))] = true
			}
		}
		return ast.WalkContinue, nil
	})
	return &ReferenceDocument{Anchors: outline.Anchors}
}
