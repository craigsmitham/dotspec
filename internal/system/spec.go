// Package system reads and discovers system specifications.
package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Filename = "system.spec.md"

var Layout = specfile.Layout{Kind: "system", Filename: Filename}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id,omitempty"`
}

const Template = `---
type: system
---

# New system

## Purpose

Describe why this system exists and whom it serves.

## Scope

Describe the system's boundary and what it includes.
`

type System struct {
	ValidationPlan           string                      `json:"validation_plan,omitempty"`
	OutOfScope               string                      `json:"out_of_scope,omitempty"`
	InteractsWith            string                      `json:"interacts_with,omitempty"`
	RequiredInterfaces       string                      `json:"required_interfaces,omitempty"`
	ProvidedInterfaces       string                      `json:"provided_interfaces,omitempty"`
	ProjectContext           string                      `json:"project_context,omitempty"`
	ID                       string                      `json:"id,omitempty"`
	Capabilities             string                      `json:"capabilities,omitempty"`
	Components               string                      `json:"components,omitempty"`
	Features                 string                      `json:"features,omitempty"`
	PartOf                   string                      `json:"part_of,omitempty"`
	ConceptUses              string                      `json:"concept_uses,omitempty"`
	Synchronizations         string                      `json:"synchronizations,omitempty"`
	References               *specfile.ReferenceDocument `json:"-"`
	Name                     string                      `json:"name"`
	Purpose                  string                      `json:"purpose"`
	Scope                    string                      `json:"scope"`
	Goals                    string                      `json:"goals,omitempty"`
	Environment              string                      `json:"environment,omitempty"`
	Assumptions              string                      `json:"assumptions,omitempty"`
	Constraints              string                      `json:"constraints,omitempty"`
	Glossary                 string                      `json:"glossary,omitempty"`
	Path                     string                      `json:"path"`
	QualityAttributes        []Attribute                 `json:"quality_attributes"`
	QualityAttributesContent *string                     `json:"quality_attributes_content,omitempty"`
}

func Read(path string) (*System, error) {
	if filepath.Base(path) != Filename {
		return nil, diagnostic.Error("layout.system-location", fmt.Sprintf("expected a %s filename", Filename), diagnostic.Location{Path: path})
	}
	if _, err := Layout.Owner(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec, err := Parse(data)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	spec.Path = path
	for i := range spec.QualityAttributes {
		spec.QualityAttributes[i].Location.Path = path
	}
	return spec, nil
}

func Parse(data []byte) (*System, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &System{Name: doc.Name, ID: meta.ID}
	fields := map[string]*string{"Validation plan": &spec.ValidationPlan, "Out of scope": &spec.OutOfScope, "Interacts with": &spec.InteractsWith, "Required interfaces": &spec.RequiredInterfaces, "Provided interfaces": &spec.ProvidedInterfaces, "Project context": &spec.ProjectContext, "Capabilities": &spec.Capabilities, "Components": &spec.Components, "Features": &spec.Features, "Part of": &spec.PartOf, "Purpose": &spec.Purpose, "Scope": &spec.Scope, "Goals": &spec.Goals, "Environment": &spec.Environment, "Assumptions": &spec.Assumptions, "Constraints": &spec.Constraints, "Glossary": &spec.Glossary, "Concept uses": &spec.ConceptUses, "Synchronizations": &spec.Synchronizations}
	locations := make(map[string]diagnostic.Location)
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "system", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Components", "Features", "Capabilities", "Concept uses", "Synchronizations"); err != nil {
		return nil, err
	}
	if err := ids.RegisterInterfaces(doc); err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		if section.Level != 2 {
			if _, known := fields[section.Title]; known || section.Title == "Quality attributes" {
				return nil, diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-two heading (## %s)", section.Title, section.Title), doc.Location(section.Start, section.Title))
			}
			continue
		}
		target, known := fields[section.Title]
		if !known && section.Title != "Quality attributes" {
			continue
		}
		location := doc.Location(section.Start, section.Title)
		if previous, exists := locations[section.Title]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		locations[section.Title] = location
		if known {
			*target = doc.Body(section)
			continue
		}
		content := doc.Body(section)
		spec.QualityAttributesContent = &content
		spec.QualityAttributes, err = ParseQualityAttributes(doc, section, &ids)
		if err != nil {
			return nil, err
		}
	}
	for _, name := range []string{"Purpose", "Scope"} {
		if *fields[name] == "" {
			location, found := locations[name]
			if !found {
				location = doc.NameLocation
				location.Element = name
			}
			return nil, diagnostic.Error("structure.required-section", fmt.Sprintf("missing or empty %s section (use ## %s)", name, name), location)
		}
	}
	if err := ids.RegisterObligations(doc); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}
