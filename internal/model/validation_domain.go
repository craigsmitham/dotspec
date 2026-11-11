package model

import (
	"fmt"
	"slices"
	"strings"
)

// modelPaths is deliberately a small grammar shared by projection and ownership
// annotations: Class.field.path [or Class.field.path ...]. Endpoints are always
// checked from induced catalog ranges, never from a second relationship table.
func modelPaths(expression string) [][2]string {
	var paths [][2]string
	for item := range strings.SplitSeq(expression, " or ") {
		class, path, found := strings.Cut(strings.TrimSpace(item), ".")
		if found && class != "" && path != "" {
			paths = append(paths, [2]string{class, path})
		}
	}
	return paths
}

func (v *graphValidator) checkDomain() {
	v.checkAuthoredIDs()
	v.checkMemberships()
	v.checkResolutions()
	v.checkLabels()
	v.checkAssertions()
	v.checkCycles("quality.broader-cycle", "QualityAttribute", func(r Record) []string { return Strings(r["broader"]) })
	v.checkCycles("ownership.cycle", "DomainElement", func(r Record) []string { return valuesAt(r, "owner") })
	v.checkCycles("system.parent-cycle", "System", func(r Record) []string { return valuesAt(r, "parent_system") })
	for _, located := range v.sortedRecords() {
		r, path := located.value, located.path
		kind := String(r, "kind")
		if kind == "QualityAttribute" {
			for _, id := range Strings(r["broader"]) {
				if other := v.records[id].value; other != nil && String(other, "owner") != String(r, "owner") {
					v.add("quality.vocabulary", path+"/broader", "broader attributes must belong to the same declared vocabulary")
				}
			}
		}
		if kind == "System" && String(r, "owning_system") == String(r, "id") {
			v.add("system.self-owner", path+"/owning_system", "a system cannot be its own enclosing system")
		}
	}
}

func (v *graphValidator) checkAuthoredIDs() {
	docs, vocabularies := map[string]string{}, map[string]string{}
	for _, located := range v.sortedRecords() {
		r := located.value
		if !v.isA(String(r, "kind"), "DomainElement") {
			continue
		}
		authored := String(r, "authored_id")
		if authored == "" {
			continue
		}
		key := String(r, "source_document") + "\x00" + authored
		if prior, ok := docs[key]; ok {
			v.add("identity.authored", located.path+"/authored_id", fmt.Sprintf("authored ID %q also occurs at %s in this document", authored, prior))
		} else {
			docs[key] = located.path
		}
		if String(r, "kind") == "QualityAttribute" {
			key = String(r, "owner") + "\x00" + authored
			if prior, ok := vocabularies[key]; ok {
				v.add("quality.identity", located.path+"/authored_id", fmt.Sprintf("attribute ID %q also occurs at %s in this vocabulary", authored, prior))
			} else {
				vocabularies[key] = located.path
			}
		}
		if fragment := v.records[String(r, "source_fragment")].value; fragment != nil && String(fragment, "authored_anchor") != "" && String(fragment, "authored_anchor") != authored {
			v.add("identity.fragment", located.path+"/source_fragment", "authored-ID fragment must identify this element's authored ID")
		}
	}
}

func (v *graphValidator) checkMemberships() {
	expected := map[string]string{}
	for _, located := range v.sortedRecords() {
		r := located.value
		seenPaths := map[string]bool{}
		for _, ancestor := range v.catalog.Classes[String(r, "kind")].Ancestors {
			for _, member := range modelPaths(v.catalog.Classes[ancestor].Annotations["owned_memberships"]) {
				if !v.isA(String(r, "kind"), member[0]) || seenPaths[member[1]] {
					continue
				}
				seenPaths[member[1]] = true
				seen := map[string]bool{}
				for _, id := range valuesAt(r, member[1]) {
					path := located.path + "/" + strings.ReplaceAll(member[1], ".", "/")
					if seen[id] {
						v.add("membership.duplicate", path, fmt.Sprintf("member %q is repeated", id))
					}
					seen[id] = true
					child := v.records[id].value
					if child == nil {
						continue
					}
					if old, ok := expected[id]; ok && old != String(r, "id") {
						v.add("membership.multiple", path, fmt.Sprintf("member %q also belongs to %q", id, old))
					}
					expected[id] = String(r, "id")
					if String(child, "owner") != String(r, "id") {
						v.add("membership.owner", path, fmt.Sprintf("member %q does not name this record as its owner", id))
					}
					if v.sourceOf(child) != v.sourceOf(r) {
						v.add("membership.source", path, "owned members must be expressed in the owner's source document")
					}
				}
			}
		}
	}
	for _, located := range v.sortedRecords() {
		r := located.value
		owner := String(r, "owner")
		if owner == "" || !v.isA(String(r, "kind"), "DomainElement") {
			continue
		}
		parent := v.records[owner].value
		// A Document indexes all mapped elements; only an explicitly document-owned
		// unknown section is an ownership claim in that otherwise broad index.
		if String(r, "kind") == "ExtensionSection" && String(parent, "kind") == "Document" {
			if !slices.Contains(Strings(parent["elements"]), String(r, "id")) {
				v.add("membership.missing", located.path+"/owner", "document-owned extension is absent from the document's element index")
			}
			if v.sourceOf(r) != v.sourceOf(parent) {
				v.add("membership.source", located.path+"/owner", "document-owned extension has another source")
			}
			continue
		}
		if expected[String(r, "id")] != owner {
			v.add("membership.missing", located.path+"/owner", "owner does not contain this record in a canonical membership field")
		}
	}
}

func (v *graphValidator) checkResolutions() {
	for _, located := range v.sortedRecords() {
		r, path := located.value, located.path
		switch String(r, "kind") {
		case "SourceReference", "LabelOccurrence":
			if id := String(r, "resolution"); id != "" {
				result := v.records[id].value
				if result != nil && String(result, "occurrence") != String(r, "id") {
					v.add("resolution.reciprocal", path+"/resolution", "resolution must point back to this occurrence")
				}
			}
			node := v.records[String(r, "content_node")].value
			source := String(r, "source_document")
			if String(r, "kind") == "LabelOccurrence" {
				source = v.sourceOf(v.records[String(r, "owner")].value)
			}
			if node != nil && v.sourceOf(node) != source {
				v.add("reference.source", path+"/content_node", "authored occurrence must be in its source document")
			}
			if owner := v.records[String(r, "context_owner")].value; owner != nil && v.sourceOf(owner) != source {
				v.add("reference.context", path+"/context_owner", "context owner must be expressed in the occurrence's source")
			}
			if definition := v.records[String(r, "definition_node")].value; definition != nil && v.sourceOf(definition) != source {
				v.add("reference.definition", path+"/definition_node", "reference definition must be in the authored occurrence's document")
			}
		case "Resolution":
			occurrence := v.records[String(r, "occurrence")].value
			if occurrence != nil && String(occurrence, "resolution") != String(r, "id") {
				v.add("resolution.reciprocal", path+"/occurrence", "occurrence must point back to this resolution")
			}
			candidates := Strings(r["candidates"])
			if String(r, "status") == "ambiguous" {
				if len(uniqueStrings(candidates)) < 2 {
					v.add("resolution.ambiguity", path+"/candidates", "an ambiguous resolution requires at least two distinct candidates")
				}
			} else if len(candidates) > 0 {
				v.add("resolution.candidates", path+"/candidates", "only an ambiguous resolution has candidates")
			}
		}
	}
}

func (v *graphValidator) checkLabels() {
	actual := map[string][]string{}
	for _, located := range v.sortedRecords() {
		label, path := located.value, located.path
		if String(label, "kind") != "LabelOccurrence" {
			continue
		}
		owner := v.records[String(label, "owner")].value
		if owner == nil {
			continue
		}
		field, targets := "attribute_labels", "attributes"
		vocabulary, system := "", ""
		switch String(label, "role") {
		case "quality_classification":
			if !v.isA(String(owner, "kind"), "Requirement") {
				v.add("label.owner", path+"/owner", "quality-classification labels belong to requirements")
				continue
			}
			system = String(owner, "owning_system")
			if system == "" {
				system = String(label, "system_context")
			}
			vocabulary = String(v.records[system].value, "quality_model")
		case "broader_attribute":
			field, targets = "broader_labels", "broader"
			if String(owner, "kind") != "QualityAttribute" {
				v.add("label.owner", path+"/owner", "broader labels belong to quality attributes")
				continue
			}
			vocabulary = String(owner, "owner")
			system = String(v.records[vocabulary].value, "owner")
		}
		if !slices.Contains(Strings(owner[field]), String(label, "id")) {
			v.add("label.membership", path+"/owner", "owner must retain this authored label occurrence")
		}
		if context := String(label, "system_context"); context != "" && context != system {
			v.add("label.context", path+"/system_context", "label system context disagrees with its owner")
		}
		if context := String(label, "vocabulary_context"); context != "" && context != vocabulary {
			v.add("label.context", path+"/vocabulary_context", "label vocabulary context disagrees with the declared vocabulary")
		}
		result := v.records[String(label, "resolution")].value
		status := String(result, "status")
		if result == nil || status == "not_attempted" {
			continue
		}
		if vocabulary == "" {
			if status != "unbound_label" {
				v.add("label.unbound", path+"/resolution", "a label with no declared vocabulary must remain unbound")
			}
			continue
		}
		matches := []string{}
		for _, id := range Strings(v.records[vocabulary].value["attributes"]) {
			if String(v.records[id].value, "authored_id") == String(label, "label") {
				matches = append(matches, id)
			}
		}
		if len(matches) == 0 {
			if status != "undefined_label" {
				v.add("label.undefined", path+"/resolution", "a label missing from its declared vocabulary must be undefined")
			}
		} else if len(matches) == 1 {
			if status != "resolved" || String(result, "target") != matches[0] {
				v.add("label.target", path+"/resolution", "label must resolve to the matching attribute in its declared vocabulary")
			}
		}
		if status == "resolved" {
			key := String(owner, "id") + "\x00" + targets
			actual[key] = append(actual[key], String(result, "target"))
		}
	}
	for _, located := range v.sortedRecords() {
		r := located.value
		field, labelField := "", ""
		if v.isA(String(r, "kind"), "Requirement") {
			field, labelField = "attributes", "attribute_labels"
		} else if String(r, "kind") == "QualityAttribute" {
			field, labelField = "broader", "broader_labels"
		}
		if field == "" {
			continue
		}
		for _, id := range Strings(r[labelField]) {
			label := v.records[id].value
			if label != nil && String(label, "owner") != String(r, "id") {
				v.add("label.membership", located.path+"/"+labelField, "label occurrence names a different owner")
			}
		}
		projected := Strings(r[field])
		expected := actual[String(r, "id")+"\x00"+field]
		if len(uniqueStrings(projected)) != len(projected) || !sameStrings(projected, expected) {
			v.add("label.projection", located.path+"/"+field, "resolved semantic targets must agree exactly with authored label resolutions")
		}
	}
}

func (v *graphValidator) checkAssertions() {
	seen := map[string]bool{}
	for _, located := range v.sortedRecords() {
		r, path := located.value, located.path
		kind := String(r, "kind")
		if !v.isA(kind, "Assertion") {
			continue
		}
		key := kind + "\x00" + String(r, "subject") + "\x00" + String(r, "object")
		if seen[key] {
			v.add("assertion.duplicate", path, "a materialized semantic relation is repeated")
		}
		seen[key] = true
		if String(r, "origin") == "derived" && (String(r, "mapping_rule") == "" || len(Strings(r["inputs"])) == 0) {
			v.add("assertion.derivation", path, "derived assertions require a mapping rule and input records")
		}
		provenance := len(Strings(r["source_nodes"])) > 0 || len(records(r["source_spans"])) > 0
		for _, id := range Strings(r["inputs"]) {
			if v.sourceOf(v.records[id].value) != "" {
				provenance = true
			}
		}
		if !provenance {
			v.add("assertion.provenance", path, "assertion must retain source provenance directly or through input records")
		}
		expression := v.catalog.Classes[kind].Annotations["projection_of"]
		if expression == "" {
			continue
		}
		subject := v.records[String(r, "subject")].value
		agrees := false
		for _, projection := range modelPaths(expression) {
			if v.isA(String(subject, "kind"), projection[0]) && slices.Contains(valuesAt(subject, projection[1]), String(r, "object")) {
				agrees = true
			}
		}
		if !agrees {
			v.add("assertion.projection", path, "materialized relation disagrees with its canonical typed membership or reference field")
		}
	}
}

func uniqueStrings(values []string) []string {
	set := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if !set[value] {
			set[value] = true
			out = append(out, value)
		}
	}
	return out
}
func sameStrings(a, b []string) bool {
	a, b = uniqueStrings(a), uniqueStrings(b)
	if len(a) != len(b) {
		return false
	}
	for _, value := range a {
		if !slices.Contains(b, value) {
			return false
		}
	}
	return true
}

func (v *graphValidator) sourceOf(r Record) string {
	if r == nil {
		return ""
	}
	if source := String(r, "source_document"); source != "" {
		return source
	}
	if source := String(r, "source"); source != "" {
		return source
	}
	switch String(r, "kind") {
	case "SourceDocument":
		return String(r, "id")
	case "ContentNode":
		return String(v.records[String(r, "document")].value, "source")
	case "SourceReference":
		return String(r, "source_document")
	case "LabelOccurrence":
		return String(v.records[String(r, "owner")].value, "source_document")
	}
	return ""
}
