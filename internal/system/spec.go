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
	Name         string        `json:"name"`
	Purpose      string        `json:"purpose"`
	Scope        string        `json:"scope"`
	Path         string        `json:"path"`
	QualityModel *QualityModel `json:"quality_model"`
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
	if spec.QualityModel != nil {
		for i := range spec.QualityModel.Attributes {
			spec.QualityModel.Attributes[i].Location.Path = path
		}
	}
	return spec, nil
}

func Parse(data []byte) (*System, error) {
	doc, err := specfile.ParseOutline(data, Layout.Kind, &metadata{})
	if err != nil {
		return nil, err
	}
	spec := &System{Name: doc.Name}
	fields := map[string]*string{"Purpose": &spec.Purpose, "Scope": &spec.Scope}
	locations := make(map[string]diagnostic.Location)
	var ids specfile.Identifiers
	for _, section := range doc.Sections {
		if section.Level != 2 {
			if section.Title == "Quality attributes" {
				return nil, diagnostic.Error("structure.section-level", "quality attributes must use ## Quality attributes", doc.Location(section.Start, section.Title))
			}
			continue
		}
		target, required := fields[section.Title]
		if !required && section.Title != "Quality attributes" {
			continue
		}
		location := doc.Location(section.Start, section.Title)
		if previous, exists := locations[section.Title]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		locations[section.Title] = location
		if required {
			*target = doc.Body(section)
			continue
		}
		spec.QualityModel, err = parseQualityModel(doc, section, &ids)
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
	return spec, nil
}
