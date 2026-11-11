package obligation

import (
	"fmt"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/system"
)

// Classification separates authored IDs from system attribute resolution and grouping.
// Without a quality attribute declaration, labels stay unbound and are not reported as undefined.
type Classification struct {
	SystemPath string             `json:"system_path"`
	Resolved   []system.Attribute `json:"resolved"`
	Broader    []string           `json:"broader"`
	Undefined  []string           `json:"undefined"`
}

func Classify(ids []string, sys *system.System) Classification {
	c := Classification{Resolved: []system.Attribute{}, Broader: []string{}, Undefined: []string{}}
	if sys == nil || sys.QualityAttributes == nil {
		return c
	}
	c.SystemPath = sys.Path
	for _, id := range ids {
		if a := sys.FindQualityAttribute(id); a != nil {
			c.Resolved = append(c.Resolved, *a)
		} else {
			c.Undefined = append(c.Undefined, id)
		}
	}
	c.Broader = sys.QualityAttributeAncestors(ids)
	return c
}

func CheckAttributes(spec *Content, sys *system.System, kind string) []diagnostic.Diagnostic {
	var findings []diagnostic.Diagnostic
	for _, id := range Classify(spec.Attributes, sys).Undefined {
		location := spec.AttributeLocations[id]
		location.Path = spec.Path
		findings = append(findings, diagnostic.Advice(kind+".undefined-attribute", fmt.Sprintf("attribute %q is not declared by the associated system (%s)", id, sys.Path), location))
	}
	return findings
}
