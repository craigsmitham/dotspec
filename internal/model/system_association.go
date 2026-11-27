package model

// An explicit selection is identified by its authored section, even when the
// section's link is invalid.
func explicitSpecificationSystem(r Record) bool {
	_, explicit := r["specification_system_content"]
	return explicit
}

// Any association other than an explicit selection is the directory default
// recorded at export.
func specificationSystemOrigin(r Record) string {
	if explicitSpecificationSystem(r) {
		return "explicit"
	}
	if String(r, "specification_system") != "" {
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
		d.resolve(r, "specification_system")
		if specificationSystemOrigin(r) == "directory" {
			// The directory default has no authored reference to replay.
			values := []string{String(r, "specification_system")}
			d.projections[len(d.projections)-1].values = values
			d.projected[String(r, "id")]["specification_system"] = values
		} else if specificationSystemOrigin(r) == "explicit" && len(d.refs(r, "specification_system")) == 0 {
			d.issue(r, "specification_system", "error", "Specification system requires exactly one explicit link to a system.")
		}
	}
}

func (b *snapshotBuilder) associateSpecificationSystems() {
	// The nearest enclosing system applies only where no selection was authored;
	// an invalid explicit selection never falls back to it.
	for id, system := range b.directorySystem {
		if r := b.elements[id]; designKind(String(r, "kind"), "SystemContextSpecification") {
			if !explicitSpecificationSystem(r) {
				r["specification_system"] = system
			}
		}
	}
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

func qualityScope(attribute Record) string {
	if String(attribute, "kind") == "SystemQualityAttribute" {
		return String(attribute, "specification_system")
	}
	return String(attribute, "owner")
}
