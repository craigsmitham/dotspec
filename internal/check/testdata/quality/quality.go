// Package quality is an illustrative test integration, not a supported public
// specification format. It deliberately depends only on the shared machinery.
package quality

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var Layout = specfile.Layout{Kind: "quality", Collection: "quality", Suffix: ".quality.spec.md"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	Unit              string `yaml:"unit"`
}

type Requirement struct {
	Unit     string
	Location diagnostic.Location
}

func ParseMarkdown(input *specfile.Markdown) (*Requirement, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	return &Requirement{Unit: meta.Unit, Location: doc.NameLocation}, nil
}

func CheckType() check.Type {
	return check.Type{Layout: Layout, Inspect: func(files *specfile.SourceSet, path string) ([]diagnostic.Diagnostic, error) {
		requirement, _, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
		if err != nil {
			return nil, err
		}
		if requirement.Unit == "" {
			location := requirement.Location
			location.Path = path
			return []diagnostic.Diagnostic{diagnostic.Advice("quality.unit", "Identify the unit of measurement.", location)}, nil
		}
		return nil, nil
	}}
}
