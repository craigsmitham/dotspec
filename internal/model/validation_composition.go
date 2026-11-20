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

	}
}
