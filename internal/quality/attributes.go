package quality

import (
	"fmt"
	"slices"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/system"
)

// Classification separates authored IDs from model resolution and grouping.
// Without a model, labels stay unbound and are not reported as undefined.
type Classification struct {
	ModelPath string             `json:"model_path"`
	Resolved  []system.Attribute `json:"resolved"`
	Broader   []string           `json:"broader"`
	Undefined []string           `json:"undefined"`
}

func Classify(ids []string, sys *system.System) Classification {
	c := Classification{Resolved: []system.Attribute{}, Broader: []string{}, Undefined: []string{}}
	if sys == nil || sys.QualityModel == nil {
		return c
	}
	c.ModelPath = sys.Path
	for _, id := range ids {
		if a := sys.QualityModel.Find(id); a != nil {
			c.Resolved = append(c.Resolved, *a)
		} else {
			c.Undefined = append(c.Undefined, id)
		}
	}
	c.Broader = sys.QualityModel.Ancestors(ids)
	return c
}

func CheckAttributes(spec *Requirement, sys *system.System) []diagnostic.Diagnostic {
	var findings []diagnostic.Diagnostic
	for _, id := range Classify(spec.Attributes, sys).Undefined {
		location := spec.AttributeLocations[id]
		location.Path = spec.Path
		findings = append(findings, diagnostic.Advice("quality.undefined-attribute", fmt.Sprintf("attribute %q is not defined in the associated system's quality model (%s)", id, sys.Path), location))
	}
	return findings
}

// Filter includes exact labels and descendant classifications. Each requirement
// appears once, even when more than one label reaches the requested ancestor.
func Filter(groups []*Group, attribute string) []*Group {
	if attribute == "" {
		return groups
	}
	filtered := []*Group{}
	for _, group := range groups {
		selected := &Group{System: group.System, Requirements: []*Requirement{}}
		for _, spec := range group.Requirements {
			if slices.Contains(spec.Attributes, attribute) || slices.Contains(spec.Classification.Broader, attribute) {
				selected.Requirements = append(selected.Requirements, spec)
			}
		}
		if len(selected.Requirements) != 0 {
			filtered = append(filtered, selected)
		}
	}
	return filtered
}
