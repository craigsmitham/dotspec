// Package requirement reads and discovers requirement specifications.
package requirement

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
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

// QualityTemplate is the quality preset. It prompts for an explicit criterion;
// the result is an ordinary requirement, classified by its quality attributes.
var QualityTemplate = `---
type: requirement-spec
---

# New quality requirement

## Applies to

Identify the subject and operating conditions, workload, population, and exclusions.
For outcomes of use, identify the users, their task, available assistance, and environment; link relevant actors, jobs, or use cases.

## Statement

State one coherent quality requirement. What observable result would violate it?

## Criterion

Define or reference what counts as satisfaction. For a measured property, state the measure (including units, aggregation, and observation window as needed) and the required bound. Keep each required bound in one place; reference it if already stated above or in an authoritative source.
If a bound is unknown, record an open question and the evidence needed to choose it. Qualitative criteria are valid; keep preferred targets distinct from required bounds.

## Rationale

Explain why this quality matters and any relevant tradeoffs.

`

// Presets guide authoring; they do not create another specification type.
func ValidatePreset(preset string) error {
	switch preset {
	case "general", "quality":
		return nil
	}
	return fmt.Errorf("unknown requirement preset %q: choose general or quality", preset)
}

func TemplateFor(preset string) (string, error) {
	if err := ValidatePreset(preset); err != nil {
		return "", err
	}
	if preset == "quality" {
		return QualityTemplate, nil
	}
	return Template, nil
}

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
	Classification      Classification                 `json:"classification"`
	Markdown            string                         `json:"markdown"`
	Path                string                         `json:"path"`
	System              *system.System                 `json:"system"`
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
