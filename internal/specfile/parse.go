// Package specfile provides shared Markdown parsing and repository file access.
package specfile

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

// Metadata is the common type declaration. Type packages embed it with
// yaml:",inline" to define their own strictly decoded frontmatter schema.
type Metadata struct {
	Type string `yaml:"type"`
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
