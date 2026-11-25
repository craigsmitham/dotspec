package model

import (
	"net/url"
	"path"
	"slices"
	"strings"
)

// Populate navigation from canonical occurrences after validating the wire shape.
// This index is also used to check membership, classification scope and cycles.
func (v *graphValidator) occurrenceRelationships() {
	for _, loc := range v.sortedRecords() {
		r := loc.value
		if system := String(r, "directory_system"); system != "" {
			r["specification_system"] = system
		}
	}
	for _, loc := range v.sortedRecords() {
		r := loc.value
		kind := String(r, "kind")
		if kind != "SourceReference" && kind != "LabelOccurrence" {
			continue
		}
		resolution := Object(r["resolution"])
		target := String(resolution, "target")
		status := String(resolution, "status")
		if (target != "") != (status == "resolved") {
			v.add("resolution.state", loc.path+"/resolution", "A resolved occurrence must identify exactly one target.")
		}
		candidates := Strings(resolution["candidates"])
		if status == "ambiguous" && len(uniqueStrings(candidates)) < 2 || status != "ambiguous" && len(candidates) > 0 {
			v.add("resolution.ambiguity", loc.path+"/resolution", "Ambiguous occurrences require distinct candidates.")
		}
		owner := v.records[String(r, "context_owner")].value
		role := String(r, "role")
		if kind == "LabelOccurrence" {
			owner = v.records[String(r, "owner")].value
			if role == "quality_classification" {
				role = "quality_attributes"
			}
		} else if String(r, "source_document") != String(Object(r["span"]), "source") {
			v.add("reference.context", loc.path+"/span", "Occurrence must lie in its declared source document.")
		}
		if owner != nil && v.sourceOf(owner) != String(Object(r["span"]), "source") {
			v.add("reference.context", loc.path+"/span", "Occurrence must be expressed in its owner's source.")
		}
		if role == "" {
			continue
		}
		ownerKind := String(owner, "kind")
		slot := v.catalog.Classes[ownerKind].Slots[role]
		expected := "link"
		if kind == "LabelOccurrence" {
			expected = "label"
		}
		if slot.Annotations["occurrence_role"] != expected {
			v.add("relationship.role", loc.path+"/role", "The owner does not declare this relationship.")
			continue
		}
		if role == "specification_system" && owner["directory_system"] != nil {
			v.add("relationship.system-origin", loc.path, "An explicitly selected system cannot also have a directory default.")
		}
		if target == "" {
			continue
		}
		class := String(v.records[target].value, "kind")
		if !slices.ContainsFunc(slot.ReferenceRanges, func(parent string) bool { return v.isA(class, parent) }) {
			v.add("relationship.target", loc.path+"/resolution/target", "Expected "+strings.Join(slot.ReferenceRanges, " or ")+".")
		}
		if kind == "SourceReference" && v.isA(class, "DomainElement") {
			v.checkLinkTarget(r, v.records[target].value, loc.path)
		}
		if role == "dependencies" && String(owner, "owner") != String(v.records[target].value, "owner") {
			v.add("relationship.scope", loc.path+"/resolution/target", "Concept uses must belong to the same system.")
		}
		if slot.Multivalued {
			owner[role] = appendUnique(Strings(owner[role]), target)
		} else {
			if owner[role] != nil {
				v.add("relationship.cardinality", loc.path+"/resolution/target", "The owner already has a target for "+role+".")
			}
			owner[role] = target
		}
	}
}

func (v *graphValidator) checkLinkTarget(ref, target Record, location string) {
	u, err := url.Parse(String(ref, "authored_target"))
	source := String(v.records[String(ref, "source_document")].value, "locator")
	targetPath := String(v.records[v.sourceOf(target)].value, "locator")
	if err == nil && u.Path != "" {
		source = path.Join(path.Dir(source), u.Path)
	}
	if err == nil && !u.IsAbs() && u.Host == "" && source == targetPath {
		if u.Fragment == "" && v.isA(String(target, "kind"), "AuthoredDocumentRoot") ||
			u.Fragment != "" && String(target, "authored_id") == u.Fragment {
			return
		}
		for _, record := range v.sortedRecords() {
			f := record.value
			if String(f, "kind") == "SourceFragment" && u.Fragment != "" && String(f, "generated_anchor") == u.Fragment && namesFragment(target, f) {
				return
			}
		}
	}
	v.add("relationship.destination", location+"/resolution/target", "Resolved subject must match the authored file and fragment.")
}
