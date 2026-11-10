// Package specfile provides shared Markdown parsing and repository file access.
package specfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.yaml.in/yaml/v3"
)

type Document struct {
	Name     string
	Sections map[string]string
	Path     string
}

func Read(path, kind string, required ...string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec, err := Parse(data, kind, required...)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	spec.Path = path
	return spec, nil
}

// Parse validates a spec type and its required H2 sections. Section bodies
// retain their original Markdown; headings in code and blockquotes are ignored.
func Parse(data []byte, kind string, required ...string) (*Document, error) {
	data = Normalize(data)
	_, body, err := frontmatter(data, kind, false)
	if err != nil {
		return nil, err
	}
	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	doc := md.Parser().Parse(text.NewReader(body))
	anchors := make(map[string]bool)
	var headings []*ast.Heading
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if h, ok := node.(*ast.Heading); ok {
			if anchor, ok := h.AttributeString("id"); ok {
				anchors[string(anchor.([]byte))] = true
			}
			if h.Level <= 2 {
				headings = append(headings, h)
			}
		}
	}
	spec := &Document{}
	wanted := make(map[string]bool, len(required))
	for _, name := range required {
		wanted[name] = true
	}
	titles := 0
	var titleLocation diagnostic.Location
	sections := make(map[string]string)
	locations := make(map[string]diagnostic.Location)
	for i, h := range headings {
		name := strings.TrimSpace(string(h.Text(body)))
		location := BodyLocation(data, body, lineStart(body, h.Pos()))
		location.Element = name
		if h.Level == 1 {
			if titles > 0 || name == "" {
				return nil, diagnostic.Error("structure.title", fmt.Sprintf("expected exactly one nonempty top-level heading naming the %s", kind), location, titleLocation)
			}
			titles++
			titleLocation = location
			spec.Name = name
			continue
		}
		if !wanted[name] {
			continue
		}
		if titles == 0 {
			return nil, diagnostic.Error("structure.section-order", fmt.Sprintf("%s name heading must precede %s", kind, name), location)
		}
		if _, exists := sections[name]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", name), location, locations[name])
		}
		stop := len(body)
		if i+1 < len(headings) {
			stop = lineStart(body, headings[i+1].Pos())
		}
		sections[name] = strings.TrimSpace(string(body[headingEnd(body, h):stop]))
		locations[name] = location
	}
	if titles != 1 || spec.Name == "" {
		return nil, diagnostic.Error("structure.title", fmt.Sprintf("expected exactly one nonempty top-level heading naming the %s", kind), BodyLocation(data, body, 0))
	}
	for _, name := range required {
		if sections[name] == "" {
			location, exists := locations[name]
			if !exists {
				location = titleLocation
				location.Element = name
			}
			return nil, diagnostic.Error("structure.required-section", fmt.Sprintf("missing or empty %s section (use ## %s)", name, name), location)
		}
	}
	spec.Sections = sections
	if err := ValidateReferences(data, body, doc, anchors); err != nil {
		return nil, err
	}
	return spec, nil
}

func lineStart(source []byte, pos int) int {
	return bytes.LastIndexByte(source[:pos], '\n') + 1
}

func lineEnd(source []byte, pos int) int {
	if n := bytes.IndexByte(source[pos:], '\n'); n >= 0 {
		return pos + n + 1
	}
	return len(source)
}

func headingEnd(source []byte, h *ast.Heading) int {
	last := h.Pos()
	if h.Lines().Len() > 0 {
		last = h.Lines().At(h.Lines().Len() - 1).Start
	}
	end := lineEnd(source, last)
	// Setext headings have an underline following their text lines.
	if !bytes.HasPrefix(bytes.TrimSpace(source[lineStart(source, h.Pos()):]), []byte("#")) {
		end = lineEnd(source, end)
	}
	return end
}

func Normalize(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
}

func BodyLocation(data, body []byte, offset int) diagnostic.Location {
	return diagnostic.At(data, len(data)-len(body)+offset)
}

// MetadataLocation points to a YAML value; callers use it only after successful
// metadata parsing. YAML's line numbers start after the opening fence.
func MetadataLocation(data []byte, key string) diagnostic.Location {
	lines := bytes.Split(Normalize(data), []byte("\n"))
	end := 1
	for end < len(lines) && string(bytes.TrimSpace(lines[end])) != "---" {
		end++
	}
	var node yaml.Node
	if end < len(lines) && yaml.Unmarshal(bytes.Join(lines[1:end], []byte("\n")), &node) == nil && len(node.Content) > 0 {
		fields := node.Content[0].Content
		for i := 0; i+1 < len(fields); i += 2 {
			if fields[i].Value == key {
				return diagnostic.Location{Line: fields[i+1].Line + 1, Column: fields[i+1].Column, Element: key}
			}
		}
	}
	return diagnostic.Location{Line: 1, Column: 1, Element: key}
}

var yamlLine = regexp.MustCompile(`line (\d+)`)

func frontmatter(data []byte, kind string, allowID bool) (string, []byte, error) {
	data = Normalize(data)
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) == 0 || string(bytes.TrimSpace(lines[0])) != "---" {
		return "", nil, diagnostic.Error("metadata.frontmatter", "expected YAML frontmatter beginning with ---", diagnostic.At(data, 0))
	}
	end := 1
	for end < len(lines) && string(bytes.TrimSpace(lines[end])) != "---" {
		end++
	}
	if end == len(lines) {
		return "", nil, diagnostic.Error("metadata.frontmatter", "YAML frontmatter must end with ---", diagnostic.At(data, 0))
	}
	var metadata struct {
		Type string `yaml:"type"`
		ID   string `yaml:"id"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(bytes.Join(lines[1:end], []byte("\n"))))
	decoder.KnownFields(true)
	var systemMetadata struct {
		Type string `yaml:"type"`
	}
	var target any = &metadata
	if !allowID {
		target = &systemMetadata
	}
	if err := decoder.Decode(target); err != nil {
		location := diagnostic.Location{Line: 2, Column: 1}
		if match := yamlLine.FindStringSubmatch(err.Error()); match != nil {
			line, _ := strconv.Atoi(match[1])
			location.Line = line + 1
		}
		return "", nil, diagnostic.Error("metadata.invalid", fmt.Sprintf("invalid YAML frontmatter: %v", err), location)
	}
	if !allowID {
		metadata.Type = systemMetadata.Type
	}
	if metadata.Type != kind {
		return "", nil, diagnostic.Error("metadata.type", fmt.Sprintf("frontmatter must declare type: %s", kind), MetadataLocation(data, "type"))
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return "", nil, diagnostic.Error("metadata.invalid", "frontmatter must contain exactly one YAML document", diagnostic.Location{Line: 2, Column: 1})
	}

	body := bytes.Join(lines[end+1:], []byte("\n"))
	return metadata.ID, body, nil
}
