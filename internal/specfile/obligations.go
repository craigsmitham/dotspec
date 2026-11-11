package specfile

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

var obligationID = regexp.MustCompile(`^\*\*ID:\*\*\s*(.*?)\s*$`)

// RegisterObligations makes explicitly named clauses addressable to the same
// local and cross-file reference checks as document and other declared IDs.
// Fenced examples, quoted examples and fields below a clause's intro are excluded.
func (ids *Identifiers) RegisterObligations(o *Outline) error {
	var visit func([]*Section) error
	visit = func(sections []*Section) error {
		for _, s := range sections {
			if s.Title == "Obligations" {
				for _, clause := range s.Children {
					var previous *diagnostic.Location
					for _, node := range o.Nodes(clause) {
						p, ok := node.(*ast.Paragraph)
						if !ok {
							continue
						}
						for i := 0; i < p.Lines().Len(); i++ {
							line := p.Lines().At(i)
							match := obligationID.FindStringSubmatch(strings.TrimSpace(string(line.Value(o.Source))))
							if match == nil {
								continue
							}
							location := o.Location(line.Start, clause.Title+" / ID")
							if previous != nil {
								return diagnostic.Error("structure.duplicate-field", fmt.Sprintf("obligation %s: duplicate ID field", clause.Title), location, *previous)
							}
							if match[1] == "" {
								return diagnostic.Error("identity.invalid", "An explicit obligation ID must be nonempty.", location)
							}
							if err := ids.Register(match[1], clause.Title, location); err != nil {
								return err
							}
							previous = &location
						}
					}
				}
			}
			if err := visit(s.Children); err != nil {
				return err
			}
		}
		return nil
	}
	return visit(o.Sections)
}
