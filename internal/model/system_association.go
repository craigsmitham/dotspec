package model

func specificationSystemOrigin(r Record) string {
	if _, explicit := r["specification_system_content"]; explicit {
		return "explicit"
	}
	if String(r, "system_context") != "" {
		return "directory"
	}
	return ""
}

// Replay explicit selection or the recorded discovery default. Never use an
// applies_to target, structural parent, or feature membership as context.
func (d *designResolver) resolveSpecificationSystems() {
	for _, r := range d.elements {
		if !designKind(String(r, "kind"), "SystemContextSpecification") {
			continue
		}
		d.resolve(r, "specification_system", "System", false, nil)
		if specificationSystemOrigin(r) == "directory" {
			values := []string{String(r, "system_context")}
			d.projections[len(d.projections)-1].values = values
			d.projected[String(r, "id")]["specification_system"] = values
		} else if specificationSystemOrigin(r) == "explicit" && len(Strings(r["specification_system_references"])) == 0 {
			d.issue(r, "specification_system", "error", "Specification system requires exactly one explicit link to a system.")
		}
	}
}

func (b *snapshotBuilder) associateSpecificationSystems() {
	d := newDesignResolver(b.snapshot)
	d.resolveSpecificationSystems()
	for _, p := range d.projections {
		origin := specificationSystemOrigin(p.owner)
		if origin != "" {
			p.owner["specification_system_origin"] = origin
		}
		if len(p.values) == 1 {
			p.owner[p.field] = p.values[0]
		}
		if origin == "explicit" {
			path, _ := b.documentForSource(String(p.owner, "source_document"))
			b.contextOK[path] = len(p.values) == 1
		}
	}
}
