// Package specfile provides shared Markdown parsing and repository file access.
package specfile

import (
	"bytes"
	"fmt"
	"os"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

type Document struct {
	Name     string
	Sections map[string]string
	Path     string
}

// Metadata is the common type declaration. Type packages embed it with
// yaml:",inline" to define their own strictly decoded frontmatter schema.
type Metadata struct {
	Type       string      `yaml:"type"`
	Companions []Companion `yaml:"companions"`
}

const QualityModelSuffix = ".quality-model.md"

// AuthoredType returns the Markdown frontmatter type for a subject kind.
func AuthoredType(kind string) string {
	if kind == "quality-model" {
		return kind
	}
	return kind + "-spec"
}

// AuthoredSuffix distinguishes reusable reference documents from specification files.
func AuthoredSuffix(kind string) string {
	if kind == "quality-model" {
		return QualityModelSuffix
	}
	return "." + kind + ".spec.md"
}

func (m Metadata) DocumentType() string { return m.Type }

// TypedMetadata must be a pointer to the type's metadata destination.
type TypedMetadata interface {
	DocumentType() string
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
	return NewMarkdown(data).MetadataLocation(key)
}
