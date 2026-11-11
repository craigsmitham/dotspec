package specfile

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/yuin/goldmark/ast"
)

var obligationID = regexp.MustCompile(`^\*\*ID:\*\*\s*(.*?)\s*$`)
var obligationQualityAttributes = regexp.MustCompile(`^\*\*Quality attributes:\*\*\s*(.*?)\s*$`)

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
					var qualityLocation *diagnostic.Location
					for _, node := range o.Nodes(clause) {
						p, ok := node.(*ast.Paragraph)
						if !ok {
							continue
						}
						for i := 0; i < p.Lines().Len(); i++ {
							line := p.Lines().At(i)
							if match := obligationQualityAttributes.FindStringSubmatch(strings.TrimSpace(string(line.Value(o.Source)))); match != nil {
								location := o.Location(line.Start, clause.Title+" / Quality attributes")
								if qualityLocation != nil {
									return diagnostic.Error("structure.duplicate-field", "Duplicate Quality attributes field on obligation.", location, *qualityLocation)
								}
								qualityLocation = &location
								var labels Identifiers
								for _, label := range strings.Split(match[1], ",") {
									label = strings.TrimSpace(label)
									if label == "" {
										return diagnostic.Error("identity.invalid", "Quality attributes must contain nonempty IDs separated by commas.", location)
									}
									if err := labels.Register(label, "quality attribute", location); err != nil {
										return err
									}
								}
							}
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
