package requirement

import (
	"slices"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/obligation"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Classification = obligation.Classification

func Classify(ids []string, sys *system.System) Classification { return obligation.Classify(ids, sys) }
func CheckAttributes(spec *Requirement, sys *system.System) []diagnostic.Diagnostic {
	return obligation.CheckAttributes(spec.Content, sys, Layout.Kind)
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
