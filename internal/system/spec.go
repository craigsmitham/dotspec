// Package system reads and discovers system specifications.
package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Suffix = ".system.spec.md"

var Layout = specfile.SystemLayout

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id,omitempty"`
}

const Template = `---
type: system-spec
---

# New system

## Description

Describe what this system is in one short paragraph.

## Purpose

Describe why this system exists and whom it serves.

## Scope

Describe the system's boundary and what it includes.
`

type System struct {
	AssessmentPlans []specfile.AssessmentPlan `json:"assessment_plans,omitempty"`
	Description     string                    `json:"description,omitempty"`

	OutOfScope               string                      `json:"out_of_scope,omitempty"`
	InteractsWith            string                      `json:"interacts_with,omitempty"`
	RequiredInterfaces       string                      `json:"required_interfaces,omitempty"`
	ProvidedInterfaces       string                      `json:"provided_interfaces,omitempty"`
	ProjectContext           string                      `json:"project_context,omitempty"`
	ID                       string                      `json:"id,omitempty"`
	Components               string                      `json:"components,omitempty"`
	PartOf                   string                      `json:"part_of,omitempty"`
	ConceptUses              string                      `json:"concept_uses,omitempty"`
	Synchronizations         string                      `json:"synchronizations,omitempty"`
	References               *specfile.ReferenceDocument `json:"-"`
	Name                     string                      `json:"name"`
	Purpose                  string                      `json:"purpose"`
	Scope                    string                      `json:"scope"`
	Goals                    string                      `json:"goals,omitempty"`
	Environment              string                      `json:"environment,omitempty"`
	Stakeholders             string                      `json:"stakeholders,omitempty"`
	Assumptions              string                      `json:"assumptions,omitempty"`
	Constraints              string                      `json:"constraints,omitempty"`
	GlossarySources          []string                    `json:"glossary_sources,omitempty"`
	Glossary                 string                      `json:"glossary,omitempty"`
	Path                     string                      `json:"path"`
	QualityAttributes        []Attribute                 `json:"quality_attributes"`
	QualityAttributesContent *string                     `json:"quality_attributes_content,omitempty"`
}

func Read(path string) (*System, error) {
	if !Layout.Match(filepath.Base(path)) {
		return nil, diagnostic.Error("layout.system-location", fmt.Sprintf("expected a <name>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	owner := Layout.ContextDir(path)
	if _, err := specfile.LocalDeclaration(owner); err != nil {
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
	if err := Layout.CheckName(path, spec.Name); err != nil {
		return nil, err
	}
	spec.Path = path
	for i := range spec.QualityAttributes {
		spec.QualityAttributes[i].Location.Path = path
	}
	if err := spec.loadQualityTree(); err != nil {
		return nil, err
	}
	return spec, nil
}

func Parse(data []byte) (*System, error) {
	var meta struct {
		metadata        `yaml:",inline"`
		GlossarySources []string `yaml:"glossary_sources,omitempty"`
	}
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &System{AssessmentPlans: doc.AssessmentPlans, Name: doc.Name, ID: meta.ID, GlossarySources: meta.GlossarySources}
	fields := map[string]*string{"Description": &spec.Description, "Out of scope": &spec.OutOfScope, "Interacts with": &spec.InteractsWith, "Required interfaces": &spec.RequiredInterfaces, "Provided interfaces": &spec.ProvidedInterfaces, "Project context": &spec.ProjectContext, "Components": &spec.Components, "Part of": &spec.PartOf, "Purpose": &spec.Purpose, "Scope": &spec.Scope, "Goals": &spec.Goals, "Environment": &spec.Environment, "Stakeholders": &spec.Stakeholders, "Assumptions": &spec.Assumptions, "Constraints": &spec.Constraints, "Glossary": &spec.Glossary, "Concept uses": &spec.ConceptUses, "Synchronizations": &spec.Synchronizations}
	locations := make(map[string]diagnostic.Location)
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "system", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Components", "Concept uses", "Synchronizations"); err != nil {
		return nil, err
	}
	if err := ids.RegisterInterfaces(doc); err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		if err := rejectInlineSubject(doc, section); err != nil {
			return nil, err
		}
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
		if len(section.Children) > 0 {
			return nil, diagnostic.Error("quality.attribute-source", "declare each system attribute in its own <name>.quality-attribute.spec.md file", location)
		}
		spec.QualityAttributes = []Attribute{}
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
	if err := spec.glossaryReferences(data); err != nil {
		return nil, err
	}
	return spec, nil
}
