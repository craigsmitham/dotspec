// Package specfile provides shared Markdown parsing and repository file access.
package specfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
	"go.yaml.in/yaml/v3"
)

type Document struct {
	Name     string
	Sections map[string]string
	Path     string
}

// Metadata is the common type declaration. Type packages embed it with
// yaml:",inline" to define their own strictly decoded frontmatter schema.
type Metadata struct {
	Type string `yaml:"type"`
}

func (m Metadata) SpecType() string { return m.Type }

// TypedMetadata must be a pointer to the type's metadata destination.
type TypedMetadata interface {
	SpecType() string
}

func Read(path, kind string, metadata TypedMetadata, required ...string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec, err := Parse(data, kind, metadata, required...)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	spec.Path = path
	return spec, nil
}

// Parse validates a spec type and its required H2 sections. Section bodies
// retain their original Markdown; headings in code and blockquotes are ignored.
func Parse(data []byte, kind string, metadata TypedMetadata, required ...string) (*Document, error) {
	outline, err := ParseOutline(data, kind, metadata)
	if err != nil {
		return nil, err
	}
	spec := &Document{Name: outline.Name, Sections: make(map[string]string)}
	wanted := make(map[string]bool, len(required))
	for _, name := range required {
		wanted[name] = true
	}
	locations := make(map[string]diagnostic.Location)
	for _, section := range outline.Sections {
		if section.Level != 2 || !wanted[section.Title] {
			continue
		}
		name := section.Title
		location := outline.Location(section.Start, name)
		if previous, exists := locations[name]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", name), location, previous)
		}
		spec.Sections[name] = outline.Body(section)
		locations[name] = location
	}
	for _, name := range required {
		if spec.Sections[name] != "" {
			continue
		}
		location, exists := locations[name]
		if !exists {
			location = outline.NameLocation
			location.Element = name
		}
		return nil, diagnostic.Error("structure.required-section", fmt.Sprintf("missing or empty %s section (use ## %s)", name, name), location)
	}
	if err := ValidateReferences(outline.Original, outline.Source, outline.AST, outline.Anchors); err != nil {
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

func frontmatter(data []byte, kind string, metadata TypedMetadata) ([]byte, error) {
	data = Normalize(data)
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) == 0 || string(bytes.TrimSpace(lines[0])) != "---" {
		return nil, diagnostic.Error("metadata.frontmatter", "expected YAML frontmatter beginning with ---", diagnostic.At(data, 0))
	}
	end := 1
	for end < len(lines) && string(bytes.TrimSpace(lines[end])) != "---" {
		end++
	}
	if end == len(lines) {
		return nil, diagnostic.Error("metadata.frontmatter", "YAML frontmatter must end with ---", diagnostic.At(data, 0))
	}
	decoder := yaml.NewDecoder(bytes.NewReader(bytes.Join(lines[1:end], []byte("\n"))))
	decoder.KnownFields(true)
	if err := decoder.Decode(metadata); err != nil {
		location := diagnostic.Location{Line: 2, Column: 1}
		if match := yamlLine.FindStringSubmatch(err.Error()); match != nil {
			line, _ := strconv.Atoi(match[1])
			location.Line = line + 1
		}
		return nil, diagnostic.Error("metadata.invalid", fmt.Sprintf("invalid YAML frontmatter: %v", err), location)
	}
	if metadata.SpecType() != kind {
		return nil, diagnostic.Error("metadata.type", fmt.Sprintf("frontmatter must declare type: %s", kind), MetadataLocation(data, "type"))
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, diagnostic.Error("metadata.invalid", "frontmatter must contain exactly one YAML document", diagnostic.Location{Line: 2, Column: 1})
	}

	body := bytes.Join(lines[end+1:], []byte("\n"))
	return body, nil
}
