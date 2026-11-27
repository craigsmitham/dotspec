package specfile

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// BlankSections collects writing advice for sections whose captured body is blank.
type BlankSections struct {
	path     string
	name     diagnostic.Location
	sections map[string]diagnostic.Location
	Advice   []diagnostic.Diagnostic
}

// NewBlankSections locates advice in the file at path, using the section
// locations recorded by Capture and the name heading for absent sections.
func NewBlankSections(path string, name diagnostic.Location, sections map[string]diagnostic.Location) *BlankSections {
	return &BlankSections{path: path, name: name, sections: sections}
}

// Suggest adds advice for section when body is blank.
func (b *BlankSections) Suggest(rule, section, body, message string) {
	if strings.TrimSpace(body) != "" {
		return
	}
	location, found := b.sections[section]
	if !found {
		location = b.name
	}
	location.Path, location.Element = b.path, section
	b.Advice = append(b.Advice, diagnostic.Advice(rule, message, location))
}
