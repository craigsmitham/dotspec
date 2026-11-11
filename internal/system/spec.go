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
	doc, err := specfile.ParseOutline(data, Layout.Kind, &metadata{})
	if err != nil {
		return nil, err
	}
	spec := &System{Name: doc.Name}
	fields := map[string]*string{"Purpose": &spec.Purpose, "Scope": &spec.Scope, "Goals": &spec.Goals, "Environment": &spec.Environment, "Assumptions": &spec.Assumptions, "Constraints": &spec.Constraints, "Glossary": &spec.Glossary}
	locations := make(map[string]diagnostic.Location)
	var ids specfile.Identifiers
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
	ids.AddAnchors(doc.Anchors)
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}
