package specfile

import (
	"fmt"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

// Section spans original Markdown, including any nested subsections.
type Section struct {
	Title                    string
	Level                    int
	Start, ContentStart, End int
	Children                 []*Section
}

type Outline struct {
	Name         string
	NameLocation diagnostic.Location
	Source       []byte
	Original     []byte
	AST          ast.Node
	Sections     []*Section
	Anchors      map[string]bool
	// Companions lists the validated frontmatter companions of a typed specification.
	Companions []Companion
}

// TypedOutline is the shared Markdown extraction path. The caller supplies a
// fresh metadata destination defining its accepted fields and interprets the
// resulting sections and metadata, including any explicit identifiers.
func (md *Markdown) TypedOutline(kind string, metadata TypedMetadata) (*Outline, error) {
	companions, err := md.decode(kind, metadata)
	if err != nil {
		return nil, err
	}
	out, err := md.Outline(kind)
	if err == nil {
		out.Companions = companions
		err = out.checkFields(DocumentRootClass(kind), out.Sections, 2)
	}
	return out, err
}

func (md *Markdown) Outline(kind string) (*Outline, error) {
	return md.outline().parseHeadings(kind)
}

func (out *Outline) parseHeadings(kind string) (*Outline, error) {
	body := out.Source
	var stack []*Section
	titles := 0
	var titleLocation diagnostic.Location
	for n := out.AST.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok {
			continue
		}
		name := strings.TrimSpace(inlineText(h, body))
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
	location := diagnostic.At(o.Original, len(o.Original)-len(o.Source)+offset)
	location.Element = element
	return location
}

// Capture copies the body of each top-level section named in fields and
// records where it starts. A repeated section is rejected at its second heading.
func (o *Outline) Capture(fields map[string]*string, locations map[string]diagnostic.Location) error {
	for _, section := range o.Sections {
		target, known := fields[section.Title]
		if !known {
			continue
		}
		location := o.Location(section.Start, section.Title)
		if previous, exists := locations[section.Title]; exists {
			return diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		locations[section.Title] = location
		*target = o.Body(section)
	}
	return nil
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

// inlineText extracts labels without formatting markers, preserving the source
// spelling of text and HTML. This replaces Goldmark's deprecated Node.Text API.
func inlineText(node ast.Node, source []byte) string {
	var out strings.Builder
	var write func(ast.Node)
	write = func(node ast.Node) {
		switch n := node.(type) {
		case *ast.Text:
			out.Write(n.Value(source))
			if n.SoftLineBreak() {
				out.WriteByte('\n')
			}
		case *ast.String:
			out.Write(n.Value)
		case *ast.AutoLink:
			out.Write(n.Label(source))
		case *ast.RawHTML:
			out.Write(n.Segments.Value(source))
		default:
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				write(child)
			}
		}
	}
	write(node)
	return out.String()
}
