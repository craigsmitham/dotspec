package specfile

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Section spans original Markdown, including any nested subsections.
type Section struct {
	Title                    string
	Level                    int
	Start, ContentStart, End int
	Children                 []*Section
}

type Outline struct {
	Name, ID     string
	NameLocation diagnostic.Location
	Source       []byte
	Original     []byte
	AST          ast.Node
	Sections     []*Section
	Anchors      map[string]bool
}

// ParseOutline retains document structure for types with nested sections.
func ParseOutline(data []byte, kind string) (*Outline, error) {
	data = Normalize(data)
	id, body, err := frontmatter(data, kind, true)
	if err != nil {
		return nil, err
	}
	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	doc := md.Parser().Parse(text.NewReader(body))
	out := &Outline{ID: id, Source: body, Original: data, AST: doc, Anchors: make(map[string]bool)}
	var stack []*Section
	titles := 0
	var titleLocation diagnostic.Location
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok {
			continue
		}
		name := strings.TrimSpace(string(h.Text(body)))
		location := out.Location(lineStart(body, h.Pos()), name)
		if anchor, ok := h.AttributeString("id"); ok {
			out.Anchors[string(anchor.([]byte))] = true
		}
		if h.Level == 1 {
			if titles > 0 || name == "" {
				return nil, diagnostic.Error("structure.title", fmt.Sprintf("expected exactly one nonempty top-level heading naming the %s", kind), location, titleLocation)
			}
			titles++
			titleLocation = location
			out.Name = name
			out.NameLocation = location
			continue
		}
		if titles == 0 {
			return nil, diagnostic.Error("structure.section-order", fmt.Sprintf("%s name heading must precede %s", kind, name), location)
		}
		section := &Section{Title: name, Level: h.Level, Start: lineStart(body, h.Pos()), ContentStart: headingEnd(body, h), End: len(body)}
		for len(stack) > 0 && stack[len(stack)-1].Level >= h.Level {
			stack[len(stack)-1].End = section.Start
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			out.Sections = append(out.Sections, section)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, section)
		}
		stack = append(stack, section)
	}
	if titles != 1 || out.Name == "" {
		return nil, diagnostic.Error("structure.title", fmt.Sprintf("expected exactly one nonempty top-level heading naming the %s", kind), out.Location(0, ""))
	}
	return out, nil
}

func (o *Outline) Location(offset int, element string) diagnostic.Location {
	location := BodyLocation(o.Original, o.Source, offset)
	location.Element = element
	return location
}

func (o *Outline) Body(s *Section) string {
	return strings.TrimSpace(string(o.Source[s.ContentStart:s.End]))
}

func (o *Outline) IntroEnd(s *Section) int {
	if len(s.Children) != 0 {
		return s.Children[0].Start
	}
	return s.End
}

// Nodes returns top-level Markdown blocks in the section's introductory body.
func (o *Outline) Nodes(s *Section) []ast.Node {
	var nodes []ast.Node
	for n := o.AST.FirstChild(); n != nil; n = n.NextSibling() {
		if n.Pos() >= s.ContentStart && n.Pos() < o.IntroEnd(s) {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
