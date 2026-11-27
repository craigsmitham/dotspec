// Package requirement reads and discovers requirement specifications.
package requirement

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"go.yaml.in/yaml/v3"
)

var Layout = specfile.Layout{Kind: "requirement", Collection: "requirements", Suffix: ".requirement.spec.md", SlugExample: "reject-existing-target"}

var Template = `---
type: requirement-spec
---

# New requirement

## Applies to

Identify the subject and operating conditions, including any exclusions. Distinguish assumptions from conditions the system must accommodate. For outcomes of use, identify the users, task, and environment.

## Statement

State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.
What observable result would violate this requirement? Add a Criterion section if an explicit assessment rule helps, defining or referencing each required bound once.

## Rationale

Explain why this requirement is needed and reference its source or governing goal.

`

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string    `yaml:"id"`
	QualityAttributes yaml.Node `yaml:"quality_attributes"`
}

type Requirement struct {
	MeasurementGuidance MeasurementGuidance            `json:"-"`
	References          *specfile.ReferenceDocument    `json:"-"`
	ID                  string                         `json:"id,omitempty"`
	Slug                string                         `json:"slug"`
	Name                string                         `json:"name"`
	Applicability       string                         `json:"applicability"`
	Statement           string                         `json:"statement"`
	Criterion           string                         `json:"criterion,omitempty"`
	Rationale           string                         `json:"rationale"`
	QualityAttributes   []string                       `json:"quality_attributes"`
	Markdown            string                         `json:"markdown"`
	Path                string                         `json:"path"`
	Location            diagnostic.Location            `json:"location"`
	SectionLocations    map[string]diagnostic.Location `json:"-"`
	AttributeLocations  map[string]diagnostic.Location `json:"-"`
}

func Read(path string) (*Requirement, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*Requirement, error) {
	spec, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	spec.Slug, spec.Path = slug, path
	spec.Location.Path = path
	return spec, nil
}

// Parse interprets the requirement form without filesystem access.
func Parse(data []byte) (*Requirement, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*Requirement, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &Requirement{
		ID: meta.ID, Name: doc.Name, Location: doc.NameLocation,
		Markdown:           strings.TrimSpace(string(doc.Source)),
		QualityAttributes:  []string{},
		SectionLocations:   make(map[string]diagnostic.Location),
		AttributeLocations: make(map[string]diagnostic.Location),
	}
	if err := parseAttributes(spec, &meta.QualityAttributes); err != nil {
		return nil, err
	}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, Layout.Kind, input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Applies to": &spec.Applicability, "Statement": &spec.Statement,
		"Criterion": &spec.Criterion, "Rationale": &spec.Rationale,
	}
	if err := doc.Capture(fields, spec.SectionLocations); err != nil {
		return nil, err
	}
	spec.MeasurementGuidance = measurementGuidance(doc)
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}

func parseAttributes(spec *Requirement, node *yaml.Node) error {
	if node.Kind == 0 {
		return nil
	}
	location := diagnostic.Location{Line: node.Line + 1, Column: node.Column, Element: "quality_attributes"}
	if node.Kind != yaml.SequenceNode {
		return diagnostic.Error("requirement.attributes", "attributes must be a list of unique, nonempty string IDs", location)
	}
	var ids specfile.Identifiers
	for _, value := range node.Content {
		location := diagnostic.Location{Line: value.Line + 1, Column: value.Column, Element: "quality_attributes"}
		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
			return diagnostic.Error("requirement.attributes", "each attribute must be a nonempty string ID", location)
		}
		if err := ids.Register(value.Value, "attribute", location); err != nil {
			return err
		}
		spec.QualityAttributes = append(spec.QualityAttributes, value.Value)
		spec.AttributeLocations[value.Value] = location
	}
	return nil
}
