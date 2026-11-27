package specfile

import (
	"slices"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

// RegisterDeclarations reads IDs only from the introduction of named entries.
// Examples in code, quotations, or nested detail sections do not declare IDs.
func (ids *Identifiers) RegisterDeclarations(o *Outline, headings ...string) error {
	for _, section := range o.Sections {
		if !slices.Contains(headings, section.Title) {
			continue
		}
		for _, entry := range section.Children {
			fields, err := o.ReadFields(entry.Title, entry, "ID")
			if err != nil {
				return err
			}
			location, declared := fields.Locations["ID"]
			if !declared {
				continue
			}
			if fields.Values["ID"] == "" {
				return diagnostic.Error("identity.invalid", "An explicit declaration ID must be nonempty.", location)
			}
			if err := ids.Register(fields.Values["ID"], entry.Title, location); err != nil {
				return err
			}
		}
	}
	return nil
}
