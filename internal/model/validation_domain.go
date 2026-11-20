package model

import (
	"fmt"
	pathpkg "path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
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
	v.checkSpecificationSources()
	v.checkAuthoredIDs()
	v.checkMemberships()
	v.checkLabels()
	v.checkRelationshipCycles()
	for _, located := range v.sortedRecords() {
		r, path := located.value, located.path
		kind := String(r, "kind")
		if kind == "System" {
			v.checkGlossarySources(r, path)
		}
		if kind == "System" || kind == "Component" {
			required := Strings(r["required_interfaces"])
			for _, id := range Strings(r["provided_interfaces"]) {
				if slices.Contains(required, id) {
					v.add("interface.direction", path+"/provided_interfaces", "an interface declaration must have one direction relative to its owner")
				}
			}
		}
		if v.isA(kind, "QualityAttribute") {
			for _, id := range valuesAt(r, "subattribute_of") {
				if other := v.records[id].value; other != nil && String(other, "owner") != String(r, "owner") {
					v.add("quality.scope", path+"/subattribute_of", "parent attributes must belong to the same declaring owner")
				}
			}
		}
	}
}

// A mapped specification must be the sole primary subject of its source.
// A source rejected during parsing can still be exported with its diagnostics.
func (v *graphValidator) checkSpecificationSources() {
	docs := map[string]Record{}
	failed := map[string]bool{}
	for _, record := range v.sortedRecords() {
		r := record.value
		switch String(r, "kind") {
		case "Document":
			docs[String(r, "source")] = r
		case "Diagnostic":
			if r["phase"] == "parse" && r["severity"] == "error" {
				failed[String(r, "source")] = true
			}
		}
	}
	roots := map[string]int{}
	for _, record := range v.sortedRecords() {
		r := record.value
		kind := String(r, "kind")
		if !v.isA(kind, "Specification") {
			continue
		}
		source := String(r, "source_document")
		doc := docs[source]
		form, typed := strings.CutSuffix(String(doc, "declared_type"), "-spec")
		name := String(v.records[source].value, "locator")
		roots[source]++
		if !typed || kind != specfile.SpecificationClass(form) || !strings.HasSuffix(name, "."+form+".spec.md") || !slices.Contains(Strings(doc["elements"]), String(r, "id")) {
			v.add("specification.source", record.path+"/source_document", "a specification must be declared by its corresponding file and document type")
		}
	}
	for _, record := range v.sortedRecords() {
		r := record.value
		if String(r, "kind") != "Document" {
			continue
		}
		source := String(r, "source")
		if roots[source] > 1 || (roots[source] == 0 && v.records[source].value["selected"] == true && !failed[source]) {
			v.add("specification.root", record.path+"/elements", "each successfully mapped specification file declares one primary subject")
		}
	}
}

// Resolved vocabulary documents must match declared paths, in declaration order.
func (v *graphValidator) checkGlossarySources(sys Record, location string) {
	paths := Strings(sys["glossary_source_paths"])
	dir := pathpkg.Dir(String(v.records[v.sourceOf(sys)].value, "locator"))
	for _, id := range Strings(sys["glossary_sources"]) {
		doc := v.records[id].value
		file := String(v.records[String(doc, "source")].value, "locator")
		i := slices.IndexFunc(paths, func(name string) bool { return pathpkg.Join(dir, name) == file })
		if i < 0 || doc["glossary_terms"] == nil {
			v.add("glossary.source", location+"/glossary_sources", "expected a glossary document at a declared path")
			return
		}
		paths = paths[i+1:]
	}
}

func (v *graphValidator) checkAuthoredIDs() {
	docs, attributeScopes := map[string]string{}, map[string]string{}
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
		if v.isA(String(r, "kind"), "QualityAttribute") {
			key = String(r, "owner") + "\x00" + authored
			if prior, ok := attributeScopes[key]; ok {
				v.add("quality.identity", located.path+"/authored_id", fmt.Sprintf("attribute ID %q also occurs at %s under this owner", authored, prior))
			} else {
				attributeScopes[key] = located.path
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
					externalAttribute := member[0] == "System" && member[1] == "quality_attributes" && String(child, "kind") == "SystemQualityAttribute"
					if !externalAttribute && v.sourceOf(child) != v.sourceOf(r) {
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
		if expected[String(r, "id")] != owner {
			v.add("membership.missing", located.path+"/owner", "owner does not contain this record in a canonical membership field")
		}
	}
}

func (v *graphValidator) checkLabels() {
	index := map[string]Record{}
	for id, record := range v.records {
		index[id] = record.value
	}
	for _, located := range v.sortedRecords() {
		label, path := located.value, located.path
		if String(label, "kind") != "LabelOccurrence" {
			continue
		}
		owner := v.records[String(label, "owner")].value
		if owner == nil {
			continue
		}
		system := ""
		switch String(label, "role") {
		case "quality_classification":
			if !v.isA(String(owner, "kind"), "Requirement") {
				v.add("label.owner", path+"/owner", "quality-classification labels belong to requirements")
				continue
			}
			system = governingSystem(owner, index)
		case "subattribute_of":
			if !v.isA(String(owner, "kind"), "QualityAttribute") {
				v.add("label.owner", path+"/owner", "parent labels belong to quality attributes")
				continue
			}
			system = String(owner, "owner")
		}
		if context := String(label, "governing_quality_model"); context != "" && (String(label, "role") != "subattribute_of" || context != system || String(v.records[system].value, "kind") != "QualityModel") {
			v.add("label.context", path+"/governing_quality_model", "governing quality model must be the declaring model of a parent label")
		}
		if context := String(label, "governing_system"); context != "" && context != system {
			v.add("label.context", path+"/governing_system", "governing system disagrees with the label owner's system")
		}
		result := Object(label["resolution"])
		status := String(result, "status")
		if result == nil || status == "not_attempted" {
			continue
		}
		if _, declared := v.records[system].value["quality_attributes"]; !declared {
			if status != "unbound_label" {
				v.add("label.unbound", path+"/resolution", "a label with no declared quality attributes must remain unbound")
			}
			continue
		}
		matches := []string{}
		for _, id := range Strings(v.records[system].value["quality_attributes"]) {
			if String(v.records[id].value, "authored_id") == String(label, "label") {
				matches = append(matches, id)
			}
		}
		if len(matches) == 0 {
			if status != "undefined_label" {
				v.add("label.undefined", path+"/resolution", "a label missing from its owner's declared quality attributes must be undefined")
			}
		} else if len(matches) == 1 {
			if status != "resolved" || String(result, "target") != matches[0] {
				v.add("label.target", path+"/resolution", "label must resolve to the matching attribute in its owner's declared quality attributes")
			}
		}
		if len(matches) > 1 && (status != "ambiguous" || !sameStrings(matches, Strings(result["candidates"]))) {
			v.add("label.ambiguous", path+"/resolution", "ambiguous labels must identify the matching candidates in scope")
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
	case "LabelOccurrence":
		return String(v.records[String(r, "owner")].value, "source_document")
	}
	return ""
}

// Cycle policies come from each reference slot, independently of the overall graph.
// Reciprocal provenance and owner/member links do not make the whole graph a DAG.
func (v *graphValidator) checkRelationshipCycles() {
	fields := map[string]bool{}
	for _, cls := range v.catalog.Classes {
		for field, slot := range cls.Slots {
			if slot.Annotations["relationship_cycles"] == "forbidden" {
				fields[field] = true
			}
		}
	}
	for _, field := range sortedKeys(fields) {
		code := "relationship.cycle"
		switch field {
		case "owner":
			code = "ownership.cycle"
		case "enclosing_system":
			code = "system.enclosing-cycle"
		case "subattribute_of":
			code = "quality.subattribute-cycle"
		}
		v.checkCycles(code, "Record", func(r Record) []string {
			slot := v.catalog.Classes[String(r, "kind")].Slots[field]
			if slot.Annotations["relationship_cycles"] != "forbidden" {
				return nil
			}
			return valuesAt(r, field)
		})
	}
}
