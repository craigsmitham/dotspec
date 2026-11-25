package model

import (
	"encoding/json"
	"maps"
)

// The parser keeps node identities for link rewriting and Markdown rendering.
// A saved snapshot addresses those same passages using normalized byte ranges.
func (s Snapshot) MarshalJSON() ([]byte, error) {
	type wireSnapshot Snapshot
	out := wireSnapshot(s)
	if s.ContentNodes == nil {
		return json.Marshal(out)
	}
	out.Elements = make([]Record, len(s.Elements))
	for i, element := range s.Elements {
		r := maps.Clone(element)
		for field, slot := range compiledModel().catalog.Classes[String(r, "kind")].Slots {
			if slot.Annotations["occurrence_role"] == "" {
				continue
			}
			delete(r, field)
		}
		if element["specification_system_origin"] == "directory" {
			r["directory_system"] = element["specification_system"]
		}
		delete(r, "specification_system_origin")
		out.Elements[i] = r
	}
	out.Documents = make([]Record, len(s.Documents))
	for i, doc := range s.Documents {
		r := maps.Clone(doc)
		delete(r, "root")
		out.Documents[i] = r
	}
	located := func(records []Record) []Record {
		result := make([]Record, len(records))
		for i, record := range records {
			result[i] = locatedRecordValue(record)
		}
		return result
	}
	out.Fragments = located(s.Fragments)
	out.References = located(s.References)
	out.LabelOccurrences = located(s.LabelOccurrences)
	return json.Marshal(out)
}

func locatedRecordValue(record Record) Record {
	r := maps.Clone(record)
	delete(r, "content_node")
	delete(r, "definition_node")
	return r
}

// Rich fields retain only the source ranges in interchange.
func (r Record) MarshalJSON() ([]byte, error) {
	if _, rich := r["spans"]; rich && r["content_nodes"] != nil {
		return json.Marshal(map[string]any{"spans": r["spans"]})
	}
	return json.Marshal(map[string]any(r))
}

// navigation supplies indexed endpoints to readers of an offline snapshot.
func (s *Snapshot) navigation() *Snapshot {
	if s.ContentNodes != nil {
		return s
	}
	view := *s
	view.Elements = make([]Record, len(s.Elements))
	owners := map[string]Record{}
	for i, element := range s.Elements {
		e := maps.Clone(element)
		if system := String(e, "directory_system"); system != "" {
			e["specification_system"] = system
			e["specification_system_origin"] = "directory"
		}
		view.Elements[i] = e
		owners[String(e, "id")] = e
	}
	for _, group := range [][]Record{s.References, s.LabelOccurrences} {
		for _, occurrence := range group {
			field := String(occurrence, "role")
			owner := owners[String(occurrence, "context_owner")]
			if String(occurrence, "kind") == "LabelOccurrence" {
				owner = owners[String(occurrence, "owner")]
				if field == "quality_classification" {
					field = "quality_attributes"
				}
			}
			if owner == nil || field == "" {
				continue
			}
			if field == "specification_system" {
				owner["specification_system_origin"] = "explicit"
			}
			target := String(Object(occurrence["resolution"]), "target")
			if target == "" {
				continue
			}
			if compiledModel().catalog.Classes[String(owner, "kind")].Slots[field].Multivalued {
				owner[field] = appendUnique(Strings(owner[field]), target)
			} else {
				owner[field] = target
			}
		}
	}
	return &view
}

func namesFragment(r, fragment Record) bool {
	f := Object(fragment["span"])
	for _, span := range records(Object(r["name"])["spans"]) {
		if String(span, "source") == String(f, "source") && Number(span["start"]) == Number(f["start"]) && Number(span["end"]) == Number(f["end"]) {
			return true
		}
	}
	return false
}
