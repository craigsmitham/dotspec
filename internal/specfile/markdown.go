package specfile

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.yaml.in/yaml/v3"
)

// Markdown holds normalized source and the one syntax tree used for authoring,
// link lookup and model extraction. Offsets in Tree refer to Body.
type Markdown struct {
	InvalidUTF8   bool
	Source, Body  []byte
	Tree          ast.Node
	Metadata      *yaml.Node
	MetadataError error
}

var yamlLine = regexp.MustCompile(`line ([0-9]+):`)

func yamlError(err error) error {
	loc := diagnostic.Location{Line: 2, Column: 1}
	if line := yamlLine.FindStringSubmatch(err.Error()); line != nil {
		n, _ := strconv.Atoi(line[1])
		loc.Line = n + 1
	}
	return diagnostic.Error("metadata.invalid", "invalid YAML frontmatter: "+err.Error(), loc)
}
func NewMarkdown(data []byte) *Markdown {
	source := Normalize(data)
	invalid := !utf8.Valid(source)
	if invalid {
		source = bytes.ToValidUTF8(source, []byte("�"))
	}
	mdSource := &Markdown{Source: source, Body: source, InvalidUTF8: invalid}
	first, rest, _ := bytes.Cut(source, []byte("\n"))
	if string(bytes.TrimRight(first, " \t")) == "---" {
		mdSource.MetadataError = diagnostic.Error("metadata.frontmatter", "YAML frontmatter must end with ---", diagnostic.At(source, 0))
		start := len(source) - len(rest)
		for len(rest) > 0 {
			line, next, _ := bytes.Cut(rest, []byte("\n"))
			if string(bytes.TrimRight(line, " \t")) == "---" {
				mdSource.Body = next
				node := &yaml.Node{}
				err := yaml.Unmarshal(source[start:len(source)-len(rest)], node)
				mdSource.Metadata, mdSource.MetadataError = node, nil
				if err != nil {
					mdSource.MetadataError = yamlError(err)
				}
				break
			}
			rest = next
		}
	}
	footnotes := &sourceFootnotes{}
	md := goldmark.New(goldmark.WithExtensions(extension.GFM, footnotes), goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	tree := md.Parser().Parse(text.NewReader(mdSource.Body))
	footnotes.restore(tree)
	mdSource.Tree = tree
	return mdSource
}

func (md *Markdown) Value(key string) *yaml.Node {
	if md.Metadata == nil || len(md.Metadata.Content) == 0 {
		return nil
	}
	fields := md.Metadata.Content[0]
	if fields.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(fields.Content); i += 2 {
		if fields.Content[i].Value == key {
			return fields.Content[i+1]
		}
	}
	return nil
}

func (md *Markdown) MetadataLocation(key string) diagnostic.Location {
	if n := md.Value(key); n != nil {
		return diagnostic.Location{Line: n.Line + 1, Column: n.Column, Element: key}
	}
	return diagnostic.Location{Line: 1, Column: 1, Element: key}
}

func (md *Markdown) decode(kind string, into TypedMetadata) error {
	if md.MetadataError != nil {
		return md.MetadataError
	}
	if md.Metadata == nil {
		return diagnostic.Error("metadata.frontmatter", "expected YAML frontmatter beginning with ---", diagnostic.At(md.Source, 0))
	}
	if len(md.Metadata.Content) == 0 || md.Metadata.Content[0].Kind != yaml.MappingNode {
		return diagnostic.Error("metadata.invalid", "expected a YAML mapping", diagnostic.Location{Line: 2, Column: 1})
	}
	// Node.Decode accepts unknown keys. Check the actual destination's declared
	// metadata fields first, including fields promoted by an inline embedding.
	allowed := map[string]bool{}
	for _, f := range reflect.VisibleFields(reflect.TypeOf(into).Elem()) {
		key, _, _ := strings.Cut(f.Tag.Get("yaml"), ",")
		if key != "" {
			allowed[key] = true
		}
	}
	fields := md.Metadata.Content[0].Content
	for i := 0; i < len(fields); i += 2 {
		k := fields[i]
		if !allowed[k.Value] {
			return diagnostic.Error("metadata.invalid", fmt.Sprintf("invalid YAML frontmatter: unknown field %q", k.Value), diagnostic.Location{Line: k.Line + 1, Column: k.Column, Element: k.Value})
		}
	}
	if err := md.Metadata.Decode(into); err != nil {
		return yamlError(err)
	}
	if expected := AuthoredType(kind); into.SpecType() != expected {
		return diagnostic.Error("metadata.type", fmt.Sprintf("frontmatter must declare type: %s", expected), md.MetadataLocation("type"))
	}
	return nil
}

func (md *Markdown) outline() *Outline {
	return &Outline{Original: md.Source, Source: md.Body, AST: md.Tree, Anchors: map[string]bool{}}
}

// Goldmark's HTML-oriented footnote transformer drops unused definitions and
// moves used ones to the end. Keep its syntax parsers but restore source parents
// after inline parsing, so this content model retains definitions and source order.
type footnotePlacement struct{ node, parent ast.Node }
type sourceFootnotes struct{ placements []footnotePlacement }
type sourceFootnoteBlockParser struct {
	parser.BlockParser
	owner *sourceFootnotes
}

func (p *sourceFootnoteBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	p.owner.placements = append(p.owner.placements, footnotePlacement{node, node.Parent()})
	p.BlockParser.Close(node, reader, pc)
}
func (f *sourceFootnotes) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithBlockParsers(util.Prioritized(&sourceFootnoteBlockParser{extension.NewFootnoteBlockParser(), f}, 999)), parser.WithInlineParsers(util.Prioritized(sourceFootnoteInline{extension.NewFootnoteParser()}, 101)))
}
func (f *sourceFootnotes) restore(root ast.Node) {
	affected := map[ast.Node]bool{}
	for _, p := range f.placements {
		affected[p.parent] = true
		p.node.Parent().RemoveChild(p.node.Parent(), p.node)
		p.parent.AppendChild(p.parent, p.node)
	}
	var clean func(ast.Node)
	clean = func(n ast.Node) {
		for c := n.FirstChild(); c != nil; {
			next := c.NextSibling()
			if _, ok := c.(*east.FootnoteList); ok {
				n.RemoveChild(n, c)
			} else {
				clean(c)
			}
			c = next
		}
		if affected[n] {
			nodes := []ast.Node{}
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				nodes = append(nodes, c)
			}
			sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Pos() < nodes[j].Pos() })
			for _, c := range nodes {
				n.RemoveChild(n, c)
				n.AppendChild(n, c)
			}
		}
	}
	clean(root)
}

type sourceFootnoteInline struct{ parser.InlineParser }

func (p sourceFootnoteInline) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	line, _ := reader.PeekLine()
	if !bytes.HasPrefix(line, []byte("[^")) && !bytes.HasPrefix(line, []byte("![^")) {
		return nil
	}
	return p.InlineParser.Parse(parent, reader, pc)
}
