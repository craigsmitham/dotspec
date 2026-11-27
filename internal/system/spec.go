// Package system reads and discovers system specifications.
package system

import (
	"fmt"
	"path/filepath"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

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

type StakeholderRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type System struct {
	Description string `json:"description,omitempty"`

	OutOfScope         string                      `json:"out_of_scope,omitempty"`
	InteractsWith      string                      `json:"interacts_with,omitempty"`
	RequiredInterfaces string                      `json:"required_interfaces,omitempty"`
	ProvidedInterfaces string                      `json:"provided_interfaces,omitempty"`
	ProjectReferences  string                      `json:"project_references,omitempty"`
	ID                 string                      `json:"id,omitempty"`
	PartOf             string                      `json:"part_of,omitempty"`
	ConceptUses        string                      `json:"concept_uses,omitempty"`
	Synchronizations   string                      `json:"synchronizations,omitempty"`
	References         *specfile.ReferenceDocument `json:"-"`
	Name               string                      `json:"name"`
	Purpose            string                      `json:"purpose"`
	Overview           string                      `json:"overview,omitempty"`
	Scope              string                      `json:"scope"`
	Goals              string                      `json:"goals,omitempty"`
	Environment        string                      `json:"environment,omitempty"`
	StakeholderRoles   []StakeholderRole           `json:"stakeholder_roles,omitempty"`
	Assumptions        string                      `json:"assumptions,omitempty"`
	Constraints        string                      `json:"constraints,omitempty"`
	GlossarySources    []string                    `json:"glossary_sources,omitempty"`
	Glossary           string                      `json:"glossary,omitempty"`
	Path               string                      `json:"path"`
	QualityAttributes  []Attribute                 `json:"quality_attributes"`
	advice             []diagnostic.Diagnostic
}

func Read(path string) (*System, error) {
	return ReadFrom(&specfile.SourceSet{}, path)
}

func ReadFrom(files *specfile.SourceSet, path string) (*System, error) {
	if !Layout.Match(filepath.Base(path)) {
		return nil, diagnostic.Error("layout.system-location", fmt.Sprintf("expected a <name>%s filename", Layout.Suffix), diagnostic.Location{Path: path})
	}
	owner := Layout.ContextDir(path)
	if _, err := specfile.LocalDeclaration(owner); err != nil {
		return nil, err
	}
	input, err := files.Markdown(path)
	if err != nil {
		return nil, err
	}
	spec, err := ParseMarkdown(input)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	if err := Layout.CheckName(path, spec.Name); err != nil {
		return nil, err
	}
	spec.Path = path
	if err := spec.loadQualityTree(files); err != nil {
		return nil, err
	}
	return spec, nil
}

func Parse(data []byte) (*System, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*System, error) {

	var meta struct {
		metadata        `yaml:",inline"`
		GlossarySources []string `yaml:"glossary_sources,omitempty"`
	}
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &System{Name: doc.Name, ID: meta.ID, GlossarySources: meta.GlossarySources}
	fields := map[string]*string{"Description": &spec.Description, "Out of scope": &spec.OutOfScope, "Interacts with": &spec.InteractsWith, "Required interfaces": &spec.RequiredInterfaces, "Provided interfaces": &spec.ProvidedInterfaces, "Project references": &spec.ProjectReferences, "Part of": &spec.PartOf, "Purpose": &spec.Purpose, "Overview": &spec.Overview, "Scope": &spec.Scope, "Goals": &spec.Goals, "Environment": &spec.Environment, "Assumptions": &spec.Assumptions, "Constraints": &spec.Constraints, "Glossary": &spec.Glossary, "Concept uses": &spec.ConceptUses, "Synchronizations": &spec.Synchronizations}
	locations := make(map[string]diagnostic.Location)
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "system", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Required interfaces", "Provided interfaces", "Concept uses", "Synchronizations"); err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		target, known := fields[section.Title]
		if !known && section.Title != "Stakeholders" {
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
		for _, role := range section.Children {
			spec.StakeholderRoles = append(spec.StakeholderRoles, StakeholderRole{Name: role.Title, Description: doc.Body(role)})
		}
	}
	if spec.Purpose == "" {
		location, found := locations["Purpose"]
		if !found {
			location = doc.NameLocation
			location.Element = "Purpose"
		}
		return nil, diagnostic.Error("structure.required-section", "missing or empty Purpose section (use ## Purpose)", location)
	}
	ids.AddAnchors(doc.Anchors)
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	spec.advice = boundaryAdvice(doc)
	if err := spec.glossaryReferences(input); err != nil {
		return nil, err
	}
	return spec, nil
}
