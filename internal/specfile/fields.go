package specfile

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

var fieldLine = regexp.MustCompile(`^\*\*(.+?):\*\*\s*(.*?)\s*$`)

// Fields holds the `**Name:** value` lines recognized in a section's introduction.
type Fields struct {
	Values    map[string]string
	Locations map[string]diagnostic.Location
	lines     map[int]bool
}

// ReadFields reads `**Name:** value` lines for the accepted names from the
// paragraphs of a section's introduction. Code, quotations and deeper sections
// never declare fields. A repeated field is rejected at its second line, and
// owner labels the section in that message.
func (o *Outline) ReadFields(owner string, section *Section, names ...string) (Fields, error) {
	f := Fields{Values: map[string]string{}, Locations: map[string]diagnostic.Location{}, lines: map[int]bool{}}
	for _, node := range o.Nodes(section) {
		if _, ok := node.(*ast.Paragraph); !ok {
			continue
		}
		for i := 0; i < node.Lines().Len(); i++ {
			segment := node.Lines().At(i)
			match := fieldLine.FindStringSubmatch(strings.TrimSpace(string(segment.Value(o.Source))))
			if match == nil || !slices.Contains(names, match[1]) {
				continue
			}
			location := o.Location(segment.Start, section.Title+" / "+match[1])
			if previous, exists := f.Locations[match[1]]; exists {
				return f, diagnostic.Error("structure.duplicate-field", fmt.Sprintf("%s: duplicate %s field", owner, match[1]), location, previous)
			}
			f.Values[match[1]], f.Locations[match[1]] = match[2], location
			f.lines[bytes.LastIndexByte(o.Source[:segment.Start], '\n')+1] = true
		}
	}
	return f, nil
}

// Strip returns the Markdown between start and end without the recognized
// field lines. Examples that only look like fields remain.
func (f Fields) Strip(o *Outline, start, end int) string {
	var kept []string
	offset := start
	for _, line := range strings.Split(string(o.Source[start:end]), "\n") {
		if !f.lines[offset] {
			kept = append(kept, line)
		}
		offset += len(line) + 1
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}
