package model

import (
	"encoding/base64"
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
	index       map[string]Record
	elements    []Record
	projected   map[string]map[string][]string
	projections []designProjection
	issues      []designIssue
}

func newDesignResolver(s *Snapshot) *designResolver {
	d := &designResolver{index: map[string]Record{}, elements: s.Elements, projected: map[string]map[string][]string{}}
	for _, collection := range [][]Record{s.Elements, s.Sources, s.Documents, s.ContentNodes, s.References, s.Resolutions, s.Fragments} {
		for _, r := range collection {
			d.index[String(r, "id")] = r
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
func (d *designResolver) target(r Record, field string) Record {
	values := d.values(r, field)
	if len(values) == 1 {
		return d.index[values[0]]
	}
	return nil
}
func (d *designResolver) isMember(owner Record, field string, r Record) bool {
	return slices.Contains(Strings(owner[field]), String(r, "id"))
}

// A heading reference resolves to its explicitly named domain record; a file
// reference to its root concept. Ordinary reference resolution is kept intact.
func (d *designResolver) promote(target Record, kind string) []Record {
	if String(target, "kind") != "SourceFragment" && designKind(String(target, "kind"), kind) {
		return []Record{target}
	}
	matches := []Record{}
	for _, r := range d.elements {
		if !designKind(String(r, "kind"), kind) {
			continue
		}
		if String(target, "kind") == "SourceDocument" && String(r, "source_document") == String(target, "id") && String(r, "owner") == "" {
			matches = append(matches, r)
		}
		if String(target, "kind") == "SourceFragment" && slices.Contains(Strings(Object(r["name"])["content_nodes"]), String(target, "content_node")) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 0 && designKind(String(target, "kind"), kind) {
		return []Record{target}
	}
	return matches
}
func (d *designResolver) resolve(r Record, field, kind string, many bool, accepts func(Record) bool, scopes ...Record) {
	refs := Strings(r[field+"_references"])
	values := []string{}
	if !many && len(refs) > 1 {
		d.issue(r, field, "error", "Expected one explicit reference.")
	} else {
		for _, id := range refs {
			ref := d.index[id]
			if field == "specification_system" {
				u, err := url.Parse(d.referenceDestination(ref))
				if err != nil || u.RawQuery != "" || u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") {
					d.issue(r, field, "error", "Specification system must link to a repository-local system document.")
					continue
				}
			}
			if String(ref, "source_document") != String(r, "source_document") {
				d.issue(r, field, "error", "Reference must be authored in this record's source.")
				continue
			}
			resolution := d.index[String(ref, "resolution")]
			status := String(resolution, "status")
			if status != "resolved" {
				severity := "error"
				if status == "external" || status == "not_attempted" || status == "out_of_scope" {
					severity = "advice"
				}
				d.issue(r, field, severity, fmt.Sprintf("Reference remains %s; no %s relationship was inferred.", status, kind))
				continue
			}
			matches := d.promote(d.index[String(resolution, "target")], kind)
			if len(matches) != 1 {
				d.issue(r, field, "error", fmt.Sprintf("Reference must identify exactly one %s.", kind))
				continue
			}
			unknownScope := false
			for _, scope := range scopes {
				unknownScope = unknownScope || scope == nil
			}
			if unknownScope {
				d.issue(r, field, "advice", "Scope is not yet verifiable; resolve the declaring definition or use first.")
				continue
			}
			if accepts != nil && !accepts(matches[0]) {
				d.issue(r, field, "error", "Reference has the wrong declaring concept or system scope.")
				continue
			}
			values = appendUnique(values, String(matches[0], "id"))
		}
	}
	if d.projected[String(r, "id")] == nil {
		d.projected[String(r, "id")] = map[string][]string{}
	}
	d.projected[String(r, "id")][field] = values
	d.projections = append(d.projections, designProjection{r, field, values, many})
}
func (d *designResolver) sameSystem(r Record, candidate Record) bool {
	owner := d.index[String(r, "owner")]
	if String(owner, "kind") == "System" {
		return String(candidate, "owner") == String(owner, "id")
	}
	return String(candidate, "owner") == String(owner, "owner")
}
func (d *designResolver) run() {
	d.resolveSpecificationSystems()
	for _, r := range d.elements {
		if String(r, "kind") == "ConceptUse" {
			d.resolve(r, "definition", "Concept", false, nil)
		}
	}
	for _, r := range d.elements {
		switch String(r, "kind") {
		case "System":
			d.resolve(r, "part_of", "System", false, nil)
		case "Component":
			d.resolve(r, "part_of", "System or Component", false, nil)
		case "Feature":
			d.resolve(r, "groups", "Requirement or ObligationClause or UseCase or Concept or ConceptUse or Synchronization or Interface", true, nil)
		case "UseCase":
			d.resolve(r, "subject_system", "System", false, nil)
		case "GeneralRequirement", "QualityRequirement", "BusinessRule", "ObligationClause":
			d.resolve(r, "applies_to", "DomainElement", true, nil)
		case "ConceptUse":
			d.resolve(r, "dependencies", "ConceptUse", true, func(t Record) bool { return d.sameSystem(r, t) })
		case "TypeBinding":
			use := d.index[String(r, "owner")]
			definition := d.target(use, "definition")
			d.resolve(r, "parameter", "ConceptType", false, func(t Record) bool { return d.isMember(definition, "type_parameters", t) }, definition)
			d.resolve(r, "supplying_concept_use", "ConceptUse", false, func(t Record) bool { return d.sameSystem(r, t) })
			typeUse := d.target(r, "supplying_concept_use")
			typeDefinition := d.target(typeUse, "definition")
			d.resolve(r, "bound_type", "ConceptType", false, func(t Record) bool { return d.isMember(typeDefinition, "introduced_types", t) }, typeDefinition)
			if String(r, "primitive") != "" && (len(Strings(r["bound_type_references"])) > 0 || len(Strings(r["supplying_concept_use_references"])) > 0) {
				d.issue(r, "bound_type", "error", "Choose a primitive or a type from a concept use, not both.")
			}
		case "ActionParameter":
			action := d.index[String(r, "owner")]
			d.resolve(r, "type", "ConceptType", false, func(t Record) bool { return String(t, "owner") == String(action, "owner") })
			if String(r, "primitive") != "" && len(Strings(r["type_references"])) > 0 {
				d.issue(r, "type", "error", "Choose a primitive or a concept type, not both.")
			}
		case "SynchronizationAction":
			d.resolve(r, "concept_use", "ConceptUse", false, func(t Record) bool { return d.sameSystem(r, t) })
			definition := d.target(d.target(r, "concept_use"), "definition")
			d.resolve(r, "action", "ConceptAction", false, func(t Record) bool { return String(t, "owner") == String(definition, "id") }, definition)
		}
	}
	d.resolveTraceability()
	d.checkDeclarations()
	d.traceReview()
}
func (d *designResolver) text(value any) string {
	var out strings.Builder
	previousSource := ""
	previousEnd := -1
	for _, id := range Strings(Object(value)["content_nodes"]) {
		n := d.index[id]
		span := Object(n["span"])
		source := d.index[String(span, "source")]
		data, _ := base64.StdEncoding.DecodeString(String(Object(source["payload"]), "bytes_base64"))
		start, end := Number(span["start"]), Number(span["end"])
		if start >= 0 && end >= start && end <= len(data) {
			sourceID := String(span, "source")
			if out.Len() > 0 {
				if sourceID == previousSource && start >= previousEnd && strings.TrimSpace(string(data[previousEnd:start])) == "" {
					out.Write(data[previousEnd:start])
				} else {
					out.WriteString("\n\n")
				}
			}
			out.Write(data[start:end])
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
		case "Interface", "Component", "Feature", "Capability", "ConceptUse", "ConceptType", "ConceptAction", "ActionParameter", "Synchronization":
			role := ""
			if kind == "ActionParameter" && d.isMember(d.index[String(r, "owner")], "output_parameters", r) {
				role = "output"
			}
			name := d.name(r)
			key := String(r, "owner") + "/" + kind + "/" + role + "/" + name
			if name != "" && names[key] {
				d.issue(r, "name", "error", "Duplicate name within this declaration scope: "+name)
			}
			names[key] = true
		}
		if kind == "ConceptUse" {
			if len(Strings(r["definition_references"])) == 0 {
				d.issue(r, "definition", "advice", "Link the concept definition when this use is ready to resolve.")
			}
			seen := map[string]bool{}
			for _, id := range Strings(r["bindings"]) {
				binding := d.index[id]
				parameter := d.target(binding, "parameter")
				key := String(parameter, "id")
				if key != "" && seen[key] {
					d.issue(binding, "parameter", "error", "This parameter already has a binding in the same concept use.")
				}
				seen[key] = true
				if parameter == nil || (String(binding, "primitive") == "" && (d.target(binding, "bound_type") == nil || d.target(binding, "supplying_concept_use") == nil)) {
					d.issue(binding, "parameter", "advice", "Complete the parameter and its primitive or named-use type binding.")
				}
			}
			for _, id := range Strings(d.target(r, "definition")["type_parameters"]) {
				if !seen[id] {
					d.issue(r, "bindings", "advice", "No resolved binding for parameter "+d.name(d.index[id])+".")
				}
			}
		}
		if kind == "SynchronizationAction" && (d.target(r, "concept_use") == nil || d.target(r, "action") == nil) {
			d.issue(r, "action", "advice", "Qualify the action with its named concept use and definition's action.")
		}
	}
}

func (b *snapshotBuilder) resolveComposition() bool {
	d := newDesignResolver(b.snapshot)
	d.run()
	for _, p := range d.projections {
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
func designKind(actual, expected string) bool {
	ancestors := compiledModel().catalog.Classes[actual].Ancestors
	for _, kind := range strings.Split(expected, " or ") {
		if slices.Contains(ancestors, kind) {
			return true
		}
	}
	return false
}
func designIssueCode(field string) string {
	switch field {
	case "measurement_basis", "illustrates", "interacts_with", "depends_on", "provided_by", "part_of", "groups", "specification_system", "subject_system", "supports", "applies_to", "derived_from", "refines", "relies_on", "involves", "motivation", "verification_plan":
		return "relationship." + field
	default:
		return "concept." + field
	}
}

func (d *designResolver) referenceDestination(ref Record) string {
	if String(ref, "form") == "reference_link" {
		if destination := String(d.index[String(ref, "definition_node")], "destination"); destination != "" {
			return destination
		}
		if destination := String(d.index[String(ref, "content_node")], "destination"); destination != "" {
			return destination
		}
	}
	return String(ref, "authored_target")
}
