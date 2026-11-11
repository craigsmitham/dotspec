package specfile

import (
	"fmt"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

// RegisterDeclarations reads IDs only from the introduction of named entries.
// Examples in code, quotations, or nested detail sections do not declare IDs.
func (ids *Identifiers) RegisterDeclarations(o *Outline, headings ...string) error {
	for _, section := range o.Sections {
		if section.Level != 2 || !slices.Contains(headings, section.Title) {
			continue
		}
		for _, entry := range section.Children {
			var previous *diagnostic.Location
			for _, node := range o.Nodes(entry) {
				paragraph, ok := node.(*ast.Paragraph)
				if !ok {
					continue
				}
				for i := 0; i < paragraph.Lines().Len(); i++ {
					line := paragraph.Lines().At(i)
					match := obligationID.FindStringSubmatch(strings.TrimSpace(string(line.Value(o.Source))))
					if match == nil {
						continue
					}
					location := o.Location(line.Start, entry.Title+" / ID")
					if previous != nil {
						return diagnostic.Error("structure.duplicate-field", fmt.Sprintf("%s: duplicate ID field", entry.Title), location, *previous)
					}
					if match[1] == "" {
						return diagnostic.Error("identity.invalid", "An explicit declaration ID must be nonempty.", location)
					}
					if err := ids.Register(match[1], entry.Title, location); err != nil {
						return err
					}
					previous = &location
				}
			}
		}
	}
	return nil
}
