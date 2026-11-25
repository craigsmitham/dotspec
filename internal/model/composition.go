package model

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

// Design relationships are projections of explicit Markdown reference occurrences.
// They never fetch additional files or interpret behavior, expressions or names.
type designIssue struct {
	owner                    Record
	field, severity, message string
}
type designProjection struct {
	owner  Record
	field  string
	values []string
	many   bool
}
type designResolver struct {
	documentDefinitions map[string]string
	links               map[string]map[string][]string // owner ID, then authored role
	index               map[string]Record
	elements            []Record
	projected           map[string]map[string][]string
	projections         []designProjection
	issues              []designIssue
	occurrences         map[string]string
}

func newDesignResolver(s *Snapshot) *designResolver {
	d := &designResolver{index: map[string]Record{}, elements: s.Elements, projected: map[string]map[string][]string{}}
	for _, collection := range [][]Record{s.Elements, s.Sources, s.Documents, s.ContentNodes, s.References, s.Fragments} {
		for _, r := range collection {
			d.index[String(r, "id")] = r
		}
	}
	d.links = map[string]map[string][]string{}
	for _, ref := range s.References {
		if role := String(ref, "role"); role != "" {
			owner := String(ref, "context_owner")
			if d.links[owner] == nil {
				d.links[owner] = map[string][]string{}
			}
			d.links[owner][role] = append(d.links[owner][role], String(ref, "id"))
		}
	}
	return d
}
func (d *designResolver) issue(r Record, field, severity, message string) {
	d.issues = append(d.issues, designIssue{r, field, severity, message})
}
func (d *designResolver) values(r Record, field string) []string {
	return d.projected[String(r, "id")][field]
}

// A heading reference resolves to its explicitly named domain record; a file
// reference to its root concept. Ordinary reference resolution is kept intact.
func (d *designResolver) refs(r Record, field string) []string {
	return d.links[String(r, "id")][field]
}

func (d *designResolver) promote(target Record, ranges []string) []Record {
	if String(target, "kind") != "SourceFragment" && designKind(String(target, "kind"), ranges...) {
		return []Record{target}
	}
	matches := []Record{}
	for _, r := range d.elements {
		if !designKind(String(r, "kind"), ranges...) {
			continue
		}
		if String(target, "kind") == "SourceDocument" && String(r, "source_document") == String(target, "id") && String(r, "owner") == "" {
			matches = append(matches, r)
		}
		if String(target, "kind") == "SourceFragment" && namesFragment(r, target) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 0 && designKind(String(target, "kind"), ranges...) {
		return []Record{target}
	}
	return matches
}
func (d *designResolver) resolve(r Record, field string) {
	slot := compiledModel().catalog.Classes[String(r, "kind")].Slots[field]
	refs := d.refs(r, field)
	values := []string{}
	if !slot.Multivalued && len(refs) > 1 {
		d.issue(r, field, "error", "Expected one explicit reference.")
	} else {
		for _, id := range refs {
			ref := d.index[id]
			if field == "specification_system" {
				u, err := url.Parse(String(ref, "authored_target"))
				if err != nil || u.RawQuery != "" || u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") {
					d.issue(r, field, "error", "Specification system must link to a repository-local system document.")
					continue
				}
			}
			if String(ref, "source_document") != String(r, "source_document") {
				d.issue(r, field, "error", "Reference must be authored in this record's source.")
				continue
			}
			resolution := Object(ref["resolution"])
			status := String(resolution, "status")
			if status != "resolved" {
				severity := "error"
				if status == "external" || status == "not_attempted" || status == "out_of_scope" {
					severity = "advice"
				}
				d.issue(r, field, severity, fmt.Sprintf("Reference remains %s; no %s relationship was inferred.", status, strings.Join(slot.ReferenceRanges, " or ")))
				continue
			}
			matches := d.promote(d.index[String(resolution, "target")], slot.ReferenceRanges)
			if len(matches) != 1 {
				d.issue(r, field, "error", fmt.Sprintf("Reference must identify exactly one %s.", strings.Join(slot.ReferenceRanges, " or ")))
				continue
			}
			if String(r, "kind") == "ConceptUse" && field == "dependencies" && !d.sameSystem(r, matches[0]) {
				d.issue(r, field, "error", "Reference has the wrong declaring concept or system scope.")
				continue
			}
			target := String(matches[0], "id")
			if d.occurrences != nil {
				d.occurrences[id] = target
			}
			values = appendUnique(values, target)
		}
	}
	if d.projected[String(r, "id")] == nil {
		d.projected[String(r, "id")] = map[string][]string{}
	}
	d.projected[String(r, "id")][field] = values
	d.projections = append(d.projections, designProjection{r, field, values, slot.Multivalued})
}
func (d *designResolver) sameSystem(r Record, candidate Record) bool {
	system := governingSystem(r, d.index)
	return system != "" && system == governingSystem(candidate, d.index)
}
func (d *designResolver) run() {
	d.resolveSpecificationSystems()
	for _, r := range d.elements {
		if String(r, "kind") == "ConceptUse" {
			d.resolve(r, "definition")
		}
	}
	// Resolve explicit structural links before trace headings. This preserves
	// diagnostic order and lets cycle checks examine all proposed edges.
	for _, trace := range []bool{false, true} {
		for _, r := range d.elements {
			slots := compiledModel().catalog.Classes[String(r, "kind")].Slots
			for _, field := range sortedKeys(d.links[String(r, "id")]) {
				role := slots[field].Annotations
				if field == "specification_system" || field == "definition" || role["occurrence_role"] != "link" {
					continue
				}
				if (role["trace_heading"] != "") == trace {
					d.resolve(r, field)
				}
			}
		}
	}
	d.withholdCyclicLinks()
	d.checkDeclarations()
	d.traceReview()
}
func (d *designResolver) text(value any) string {
	var out strings.Builder
	previousSource := ""
	previousEnd := -1
	spans := records(Object(value)["spans"])
	if spans == nil {
		for _, id := range Strings(Object(value)["content_nodes"]) {
			spans = append(spans, Object(d.index[id]["span"]))
		}
	}
	for _, span := range spans {
		source := d.index[String(span, "source")]
		data := String(source, "markdown")
		start, end := Number(span["start"]), Number(span["end"])
		if start >= 0 && end >= start && end <= len(data) {
			sourceID := String(span, "source")
			if out.Len() > 0 {
				if sourceID == previousSource && start >= previousEnd && strings.TrimSpace(string(data[previousEnd:start])) == "" {
					out.WriteString(data[previousEnd:start])
				} else {
					out.WriteString("\n\n")
				}
			}
			out.WriteString(data[start:end])
			previousSource, previousEnd = sourceID, end
		}
	}
	return out.String()
}
func (d *designResolver) name(r Record) string {
	return strings.TrimSpace(strings.TrimLeft(d.text(r["name"]), "#"))
}
func (d *designResolver) checkDeclarations() {
	names := map[string]bool{}
	for _, r := range d.elements {
		kind := String(r, "kind")
		switch kind {
		case "Interface", "Component", "ConceptUse", "ConceptAction", "Synchronization":
			name := d.name(r)
			key := String(r, "owner") + "/" + kind + "/" + name
			if name != "" && names[key] {
				d.issue(r, "name", "error", "Duplicate name within this declaration scope: "+name)
			}
			names[key] = true
		}
		if kind == "ConceptUse" {
			if len(d.refs(r, "definition")) == 0 {
				d.issue(r, "definition", "advice", "Link the concept definition when this use is ready to resolve.")
			}

		}
	}
}

func (b *snapshotBuilder) resolveComposition() bool {
	d := newDesignResolver(b.snapshot)
	d.occurrences = map[string]string{}
	d.run()
	for _, p := range d.projections {
		for _, refID := range d.refs(p.owner, p.field) {
			ref := d.index[refID]

			resolution := Object(ref["resolution"])
			if String(resolution, "status") != "resolved" {
				continue
			}
			if target := d.occurrences[refID]; target != "" && slices.Contains(p.values, target) {
				resolution["target"] = target
			} else {
				resolution["status"] = "invalid_target"
				delete(resolution, "target")
			}
		}
		if len(p.values) == 0 {
			continue
		}
		if p.many {
			p.owner[p.field] = p.values
		} else {
			p.owner[p.field] = p.values[0]
		}
	}
	hasError := false
	for _, issue := range d.issues {
		r := Record{"id": b.next("finding"), "kind": "Diagnostic", "code": designIssueCode(issue.field), "severity": issue.severity, "message": issue.message, "phase": "resolve", "source": issue.owner["source_document"], "record": issue.owner["id"]}
		b.snapshot.Diagnostics = append(b.snapshot.Diagnostics, r)
		hasError = hasError || issue.severity == "error"
	}
	return hasError
}

// Definition and applicability references may name a concrete subtype of an abstract range.
func designKind(actual string, expected ...string) bool {
	ancestors := compiledModel().catalog.Classes[actual].Ancestors
	for _, kind := range expected {
		if slices.Contains(ancestors, kind) {
			return true
		}
	}
	return false
}
func designIssueCode(field string) string {
	switch field {
	case "measurement_basis", "illustrates", "interacts_with", "depends_on", "provided_by", "part_of", "groups", "specification_system", "subject_system", "supports", "applies_to", "derived_from", "refines", "relies_on", "involves", "motivation":
		return "relationship." + field
	default:
		return "concept." + field
	}
}
