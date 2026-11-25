package model

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

// ProductMember identifies an independent declaration and its source-ordered parts.
type ProductMember struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Source   string   `json:"source"`
	System   string   `json:"system"`
	Basis    string   `json:"basis"`
	Elements []string `json:"elements,omitempty"`
}

type ProductSpecification struct {
	System                string                `json:"system,omitempty"`
	Scope                 Record                `json:"scope"`
	Complete              bool                  `json:"complete"`
	Members               []ProductMember       `json:"members"`
	Relationships         []Record              `json:"relationships,omitempty"`
	References            []Record              `json:"references,omitempty"`
	ReferenceCatalog      []Record              `json:"reference_catalog,omitempty"`
	Adoptions             []Record              `json:"adoptions,omitempty"`
	ReferencedSources     []Record              `json:"referenced_sources,omitempty"`
	Verification          []ProcedureAssessment `json:"verification,omitempty"`
	VerificationSelection *VerificationContext  `json:"verification_selection,omitempty"`
	Diagnostics           []Record              `json:"diagnostics,omitempty"`
}

// Product selects by system record ID, authored ID, or source basename.
// Acquiring a direct target supplies context without enrolling it in a product.
func (s *Snapshot) Product(selector string) (*ProductSpecification, error) {
	s = s.navigation()
	resolver := newDesignResolver(s)
	index := resolver.index
	selected := map[string]bool{}
	project := String(Object(s.Scope["project"]), "root")
	for _, src := range s.Sources {
		key := String(src, "id")
		selected[key] = src["selected"] == true && (project == "" || String(Object(src["project"]), "root") == project)
	}
	roots := Strings(s.Scope["roots"])
	view := &ProductSpecification{Scope: s.Scope, Members: []ProductMember{}}
	view.Complete = s.Scope["inspection_complete"] == true && len(roots) == 1 && roots[0] == project && len(Strings(s.Scope["selected_spec_types"])) == len(compiledModel().catalog.Enums["SpecificationType"])
	for _, diagnostic := range s.Diagnostics {
		finding := maps.Clone(diagnostic)
		if source := index[String(diagnostic, "source")]; source != nil {
			finding["locator"] = source["locator"]
		}
		view.Diagnostics = append(view.Diagnostics, finding)
		if diagnostic["severity"] == "error" {
			view.Complete = false
		}
	}
	associated := map[string]bool{}
	for _, r := range s.Elements {
		if selected[String(r, "source_document")] {
			associated[String(r, "specification_system")] = true
		}
	}
	var choices, matches []string
	for _, r := range s.Elements {
		id := String(r, "id")
		source := String(r, "source_document")
		if r["kind"] == "QualityModel" && selected[source] {
			view.ReferenceCatalog = append(view.ReferenceCatalog, Record{"id": id, "source": index[source]["locator"]})
		}
		if r["kind"] != "System" || (project != "" && String(Object(index[source]["project"]), "root") != project) {
			continue
		}
		if !selected[source] && !associated[id] {
			continue
		}
		slug := strings.TrimSuffix(filepath.Base(String(index[source], "locator")), ".system.spec.md")
		parent := index[String(r, "part_of")]
		if selector == "" && selected[String(parent, "source_document")] {
			continue
		}
		choices = append(choices, slug)
		if selector == "" || selector == id || selector == slug || selector == String(r, "authored_id") {
			matches = append(matches, id)
		}
	}
	if len(matches) == 0 && selector == "" && len(choices) == 0 {
		for _, r := range s.Elements {
			source := String(r, "source_document")
			if selected[source] && designKind(String(r, "kind"), "Specification") {
				view.Complete = false
				view.Diagnostics = append(view.Diagnostics, Record{"source": index[source]["locator"], "record": r["id"], "message": "No root system selected in the inspected scope."})
			}
		}
		return view, nil
	}
	if len(matches) != 1 {
		slices.Sort(choices)
		return nil, fmt.Errorf("select one system by record ID, authored ID, or basename (available: %s)", strings.Join(choices, ", "))
	}
	view.System = matches[0]
	systems := map[string]bool{view.System: true}
	for _, r := range s.Elements {
		if r["kind"] != "System" || !selected[String(r, "source_document")] {
			continue
		}
		var path []string
		for parent := String(r, "id"); parent != "" && !slices.Contains(path, parent); parent = String(index[parent], "part_of") {
			if systems[parent] {
				for _, child := range path {
					systems[child] = true
				}
				break
			}
			if !selected[String(index[parent], "source_document")] {
				break
			}
			path = append(path, parent)
		}
	}
	memberSources := map[string]bool{}
	memberIDs := map[string]bool{}
	for _, doc := range s.Documents {
		src := String(doc, "source")
		for _, id := range Strings(doc["elements"]) {
			r := index[id]
			if !designKind(String(r, "kind"), "Specification") || (!selected[src] && id != view.System) {
				continue
			}
			system, basis := String(r, "specification_system"), specificationSystemOrigin(r)
			if r["kind"] == "System" {
				system, basis = id, "part_of"
				if id == view.System {
					basis = "selected"
				}
			}
			if system == "" {
				view.Complete = false
				view.Diagnostics = append(view.Diagnostics, Record{"source": index[src]["locator"], "record": id, "message": "Specification system is unresolved."})
			}
			if !systems[system] {
				continue
			}
			m := ProductMember{ID: id, Name: resolver.name(r), Kind: String(r, "kind"), Source: String(index[src], "locator"), System: system, Basis: basis}
			memberSources[src], memberIDs[id] = true, true
			for _, part := range Strings(doc["elements"]) {
				if part != id {
					m.Elements = append(m.Elements, part)
					memberIDs[part] = true
				}
			}
			view.Members = append(view.Members, m)
		}
	}
	slices.SortFunc(view.Members, func(a, b ProductMember) int { return strings.Compare(a.Source, b.Source) })
	for _, edge := range s.Relationships() {
		if memberIDs[String(edge, "subject")] {
			view.Relationships = append(view.Relationships, edge)
		}
	}
	reached := map[string]bool{}
	for _, r := range s.Elements {
		if r["kind"] != "AdoptedProvision" || !memberIDs[String(r, "id")] {
			continue
		}
		ref := index[String(r, "provision")]
		adoption := Record{"id": r["id"], "owner": r["owner"], "source": index[String(r, "source_document")]["locator"], "name": resolver.name(r), "provision": ref["authored_target"], "edition": resolver.text(r["edition"]), "applicability": resolver.text(r["applicability"]), "resolution": ref["resolution"]}
		if digest := r["observed_source_sha256"]; digest != nil {
			adoption["observed_source_sha256"] = digest
		}
		view.Adoptions = append(view.Adoptions, adoption)
	}
	for _, ref := range s.References {
		if memberSources[String(ref, "source_document")] {
			view.References = append(view.References, locatedRecordValue(ref))
			resolution := Object(ref["resolution"])
			target := index[String(resolution, "target")]
			if resolution["status"] == "external" {
				target = index[resourceID("external", String(ref, "authored_target"))]
			}
			if source := String(target, "source_document"); source != "" {
				target = index[source]
			} else if source := String(target, "source"); source != "" {
				target = index[source]
			}
			key := String(target, "id")
			if key != "" && !memberSources[key] && !reached[key] && target["locator"] != nil {
				reached[key] = true
				view.ReferencedSources = append(view.ReferencedSources, Record{"id": key, "kind": target["kind"], "locator": target["locator"]})
			}
		}
	}
	return view, nil
}

// IncludeEvidence attaches observations for the selected product without
// introducing a verification run or its procedure into product membership.
func (v *ProductSpecification) IncludeEvidence(s *Snapshot, root string, evidence VerificationEvidence) error {
	assessments, err := s.Assessments(root, evidence)
	if err != nil {
		return err
	}
	v.Verification = nil
	v.VerificationSelection = nil
	if evidence.Selection.ProductRevision != "" {
		v.VerificationSelection = &evidence.Selection
	}
	for _, member := range v.Members {
		for _, a := range assessments[member.ID] {
			if len(a.Procedures) > 0 || len(a.Observations) > 0 {
				v.Verification = append(v.Verification, *a)
			}
		}
	}
	return nil
}
