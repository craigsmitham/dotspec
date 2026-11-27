package system

import (
	"fmt"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

type Attribute struct {
	ID                   string              `json:"id"`
	Name                 string              `json:"name"`
	Definition           string              `json:"definition"`
	SubattributeOf       []string            `json:"subattribute_of"`
	Location             diagnostic.Location `json:"location"`
	subattributeLocation diagnostic.Location
}

// ParseQualityAttributes reads the attribute section syntax of a reusable quality model.
func ParseQualityAttributes(doc *specfile.Outline, section *specfile.Section, ids *specfile.Identifiers) ([]Attribute, error) {
	attributes := []Attribute{}
	for _, child := range section.Children {
		attribute, err := parseAttribute(doc, child, ids)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, attribute)
	}
	if err := ValidateQualityAttributes(attributes); err != nil {
		return nil, err
	}
	return attributes, nil
}

func parseAttribute(doc *specfile.Outline, section *specfile.Section, ids *specfile.Identifiers) (Attribute, error) {
	a := Attribute{Name: section.Title, SubattributeOf: []string{}, Location: doc.Location(section.Start, section.Title)}
	fields, err := doc.ReadFields("attribute "+a.Name, section, "ID", "Sub-attribute of")
	if err != nil {
		return a, err
	}
	a.ID = fields.Values["ID"]
	if a.ID == "" {
		return a, diagnostic.Error("identity.required", "each quality attribute needs a nonempty **ID:** field", a.Location)
	}
	if err := ids.Register(a.ID, "attribute "+a.Name, fields.Locations["ID"]); err != nil {
		return a, err
	}
	if value, exists := fields.Values["Sub-attribute of"]; exists {
		a.subattributeLocation = fields.Locations["Sub-attribute of"]
		var references specfile.Identifiers
		for _, id := range strings.Split(value, ",") {
			id = strings.TrimSpace(id)
			if id == "" {
				return a, diagnostic.Error("quality.subattribute", "Sub-attribute of must contain a comma-separated list of nonempty attribute IDs; omit the field for a root attribute", a.subattributeLocation)
			}
			if err := references.Register(id, "grouping attribute", a.subattributeLocation); err != nil {
				return a, err
			}
			a.SubattributeOf = append(a.SubattributeOf, id)
		}
	}
	// Remove only the recognized field lines. Preserve examples, blockquotes,
	// and all deeper Markdown sections as part of the definition.
	a.Definition = fields.Strip(doc, section.ContentStart, section.End)
	return a, nil
}

func ValidateQualityAttributes(attributes []Attribute) error {
	byID := make(map[string]*Attribute, len(attributes))
	for i := range attributes {
		byID[attributes[i].ID] = &attributes[i]
	}
	// Iterative topological traversal avoids recursion limits on authored graphs.
	remaining := make(map[string]int, len(byID))
	children := make(map[string][]string)
	var roots []string
	for _, a := range attributes {
		remaining[a.ID] = len(a.SubattributeOf)
		if len(a.SubattributeOf) == 0 {
			roots = append(roots, a.ID)
		}
		for _, id := range a.SubattributeOf {
			if byID[id] == nil {
				return diagnostic.Error("quality.subattribute", fmt.Sprintf("attribute %q refers to undefined parent attribute %q", a.ID, id), a.subattributeLocation)
			}
			children[id] = append(children[id], a.ID)
		}
	}
	for i := 0; i < len(roots); i++ {
		for _, child := range children[roots[i]] {
			remaining[child]--
			if remaining[child] == 0 {
				roots = append(roots, child)
			}
		}
	}
	if len(roots) != len(attributes) {
		for _, a := range attributes {
			if remaining[a.ID] != 0 {
				return diagnostic.Error("quality.cycle", fmt.Sprintf("attribute %q has a cyclic parent relationship or depends on one", a.ID), a.subattributeLocation)
			}
		}
	}
	return nil
}

func (s *System) FindQualityAttribute(id string) *Attribute {
	for i := range s.QualityAttributes {
		if s.QualityAttributes[i].ID == id {
			return &s.QualityAttributes[i]
		}
	}
	return nil
}

// QualityAttributeAncestors returns distinct ancestor IDs in declaration order. It does not
// turn those classifications into authored labels or inherited requirements.
func (s *System) QualityAttributeAncestors(ids []string) []string {
	seen := make(map[string]bool)
	queue := append([]string(nil), ids...)
	for _, id := range ids {
		seen[id] = true
	}
	ancestors := make(map[string]bool)
	for i := 0; i < len(queue); i++ {
		if a := s.FindQualityAttribute(queue[i]); a != nil {
			for _, parent := range a.SubattributeOf {
				ancestors[parent] = true
				if !seen[parent] {
					seen[parent] = true
					queue = append(queue, parent)
				}
			}
		}
	}
	result := []string{}
	for _, a := range s.QualityAttributes {
		if ancestors[a.ID] {
			result = append(result, a.ID)
		}
	}
	return result
}

// Check returns writing advice for a checked system declaration.
func Check(spec *System) []diagnostic.Diagnostic {
	return append(withPath(spec.advice, spec.Path), qualityAdvice(spec)...)
}

// qualityAdvice also applies when the system is context for another file,
// because its attributes classify that file's requirements.
func qualityAdvice(spec *System) []diagnostic.Diagnostic {
	var findings []diagnostic.Diagnostic
	for _, a := range spec.QualityAttributes {
		if a.Definition == "" {
			findings = append(findings, diagnostic.Advice("quality.definition", "Define what this quality attribute means for the system.", a.Location))
		}
	}
	return findings
}
