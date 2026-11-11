// Package obligation shares the authoring model for general and quality requirements.
package obligation

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"go.yaml.in/yaml/v3"
)

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string    `yaml:"id"`
	QualityAttributes yaml.Node `yaml:"quality_attributes"`
}

type Content struct {
	ValidationPlan      string                         `json:"validation_plan,omitempty"`
	MeasurementGuidance MeasurementGuidance            `json:"-"`
	References          *specfile.ReferenceDocument    `json:"-"`
	ID                  string                         `json:"id,omitempty"`
	Slug                string                         `json:"slug"`
	Name                string                         `json:"name"`
	Applicability       string                         `json:"applicability"`
	Statement           string                         `json:"statement"`
	Criterion           string                         `json:"criterion,omitempty"`
	VerificationPlan    string                         `json:"verification_plan"`
	MeasurementPlan     string                         `json:"measurement_plan,omitempty"`
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

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte, kind string) (*Content, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &Content{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location), QualityAttributes: []string{}, AttributeLocations: make(map[string]diagnostic.Location)}
	if err := parseAttributes(spec, &meta.QualityAttributes, kind); err != nil {
		return nil, err
	}
	spec.Classification = Classify(spec.QualityAttributes, nil)
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, kind, specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterObligations(doc); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Applies to": &spec.Applicability, "Statement": &spec.Statement,
		"Validation plan": &spec.ValidationPlan, "Verification plan": &spec.VerificationPlan, "Measurement plan": &spec.MeasurementPlan,
		"Criterion": &spec.Criterion, "Rationale": &spec.Rationale,
	}
	for _, section := range doc.Sections {
		target, known := fields[section.Title]
		if !known {
			continue
		}
		location := doc.Location(section.Start, section.Title)
		if section.Level != 2 {
			return nil, diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-two heading (## %s)", section.Title, section.Title), location)
		}
		if previous, exists := spec.SectionLocations[section.Title]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		spec.SectionLocations[section.Title] = location
		*target = doc.Body(section)
	}
	spec.MeasurementGuidance = measurementGuidance(doc)
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}

func parseAttributes(spec *Content, node *yaml.Node, kind string) error {
	if node.Kind == 0 {
		return nil
	}
	location := diagnostic.Location{Line: node.Line + 1, Column: node.Column, Element: "quality_attributes"}
	if node.Kind != yaml.SequenceNode {
		return diagnostic.Error(kind+".attributes", "attributes must be a list of unique, nonempty string IDs", location)
	}
	var ids specfile.Identifiers
	for _, value := range node.Content {
		location := diagnostic.Location{Line: value.Line + 1, Column: value.Column, Element: "quality_attributes"}
		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
			return diagnostic.Error(kind+".attributes", "each attribute must be a nonempty string ID", location)
		}
		if err := ids.Register(value.Value, "attribute", location); err != nil {
			return err
		}
		spec.QualityAttributes = append(spec.QualityAttributes, value.Value)
		spec.AttributeLocations[value.Value] = location
	}
	return nil
}
