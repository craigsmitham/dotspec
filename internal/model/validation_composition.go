package model

import "slices"

func (v *graphValidator) checkComposition() {
	d := &designResolver{index: map[string]Record{}, elements: records(v.root["elements"]), projected: map[string]map[string][]string{}}
	for id, r := range v.records {
		d.index[id] = r.value
	}
	d.run()
	for _, p := range d.projections {
		actual := []string{}
		if p.many {
			actual = Strings(p.owner[p.field])
		} else if id := String(p.owner, p.field); id != "" {
			actual = []string{id}
		}
		if !slices.Equal(actual, p.values) {
			v.add("concept.projection", v.records[String(p.owner, "id")].path+"/"+p.field, "Resolved relationship must match explicit references and their concept/system scope.")
		}
	}
	for _, located := range v.sortedRecords() {
		r := located.value
		if designKind(String(r, "kind"), "SystemContextSpecification") && String(r, "specification_system_origin") != specificationSystemOrigin(r) {
			v.add("relationship.system-origin", located.path+"/specification_system_origin", "System association provenance must match the explicit selection or directory default.")
		}
		if String(r, "kind") == "Feature" && String(r, "system") != String(r, "owner") {
			v.add("relationship.feature-system", located.path+"/system", "A feature belongs to its explicitly declaring system.")
		}
		for field, slot := range v.catalog.Classes[String(r, "kind")].Slots {
			if slot.Annotations["trace_targets"] == "" && slot.Annotations["relationship_source"] == "" {
				continue
			}
			roots := Strings(Object(r[field+"_content"])["content_nodes"])
			for _, id := range Strings(r[field+"_references"]) {
				node := String(v.records[id].value, "content_node")
				matched := false
				for _, root := range roots {
					matched = matched || v.descendsFrom(node, root)
				}
				if !matched {
					v.add("relationship.provenance", located.path+"/"+field+"_references", "Traceability occurrences must come from their declared relationship section.")
				}
			}
		}

		// A single declaration/occurrence has one role. Separate records can of
		// course describe corresponding inputs/results or trigger/invoked actions.
		rolePairs := map[string][2]string{"Concept": {"type_parameters", "introduced_types"}, "ConceptAction": {"input_parameters", "output_parameters"}, "Synchronization": {"when_actions", "then_actions"}}
		if pair, ok := rolePairs[String(r, "kind")]; ok {
			for _, id := range Strings(r[pair[0]]) {
				if slices.Contains(Strings(r[pair[1]]), id) {
					v.add("concept.role", located.path+"/"+pair[1], "A declaration or action occurrence cannot have both roles; use distinct records.")
				}
			}
		}
	}
}
