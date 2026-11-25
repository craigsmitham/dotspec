package specdoc

import (
	"slices"
	"strings"
)

// Each direction has its own wording and navigation. A directory link can be
// served by the overview's existing listing when that listing contains its target.
type relationCaption struct {
	label        string
	useDirectory bool
}

func onDetail(label string) relationCaption { return relationCaption{label: label} }
func inDirectory(label string) relationCaption {
	return relationCaption{label: label, useDirectory: true}
}

type relation struct{ forward, reverse relationCaption }

var relations = map[string]relation{
	"provided_by":         {onDetail("Provided by"), inDirectory("Provides")},
	"required_interfaces": {inDirectory("Required interfaces"), onDetail("Required by")},
	"provided_interfaces": {inDirectory("Provided interfaces"), onDetail("Provided by")},
	"subattribute_of":     {onDetail("Parent attribute"), onDetail("Subattributes")},
	"groups":              {onDetail("Groups"), onDetail("Feature")},
	"refines":             {onDetail("Refines"), onDetail("Refined by")},
	"supports":            {onDetail("Supports"), onDetail("Supported by")},
	"applies_to":          {onDetail("Applies to"), inDirectory("Applies here")},
	"quality_attributes":  {onDetail("Quality attribute"), onDetail("Classifies")},
	"derived_from":        {onDetail("Derived from"), onDetail("Basis for")},
	"relies_on":           {onDetail("Relies on"), onDetail("Relied on by")},
	"involves":            {onDetail("Involves"), onDetail("Used in")},
	"illustrates":         {onDetail("Illustrates"), onDetail("Illustrated by")},
	"definition":          {onDetail("Uses"), onDetail("Used by")},
	"part_of":             {onDetail("Part of"), inDirectory("Parts")},
}

func reportRelation(projection string) (relation, bool) {
	// These two fields own declarations. Classification of a requirement
	// uses the same field name and does supply related-subject links.
	if projection == "QualityModel.quality_attributes" || projection == "System.quality_attributes" {
		return relation{}, false
	}
	_, field, _ := strings.Cut(projection, ".")
	rel, ok := relations[field]
	return rel, ok
}

func (c relationCaption) directoryHome(s *subject, target string) string {
	if c.useDirectory && s.isOverview() {
		for title, idx := range s.indexes {
			if slices.ContainsFunc(idx.members, func(member entry) bool { return member.target == target }) {
				return s.sectionIndex(title)
			}
		}
	}
	return ""
}

type route struct {
	relationCaption
	title, target string
	forward       bool
}
