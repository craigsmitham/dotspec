// Package qualitymodel defines reusable quality models, including the bundled ISO/IEC 25010 model.
package qualitymodel

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

var Layout = specfile.Layout{Kind: "quality-model", Collection: "quality-models", Suffix: specfile.QualityModelSuffix, SlugExample: "organization-quality"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: quality-model
---

# New quality model

## Description

Describe the purpose and intended scope of this reusable set of quality attributes.

## Edition

Identify the edition or revision, if applicable.

## Sources

Cite the sources informing these definitions, if any.

## Quality attributes
`

// QualityModel owns reusable definitions independently of any system.
type QualityModel struct {
	References        *specfile.ReferenceDocument `json:"-"`
	ID                string                      `json:"id,omitempty"`
	Slug              string                      `json:"slug"`
	Name              string                      `json:"name"`
	Description       string                      `json:"description"`
	Edition           string                      `json:"edition"`
	Sources           string                      `json:"sources"`
	QualityAttributes []system.Attribute          `json:"quality_attributes"`

	Markdown string              `json:"markdown"`
	Path     string              `json:"path"`
	Builtin  bool                `json:"builtin"`
	Location diagnostic.Location `json:"location"`
	// Data retains the full source for copying and reference relocation.
	Data []byte `json:"-"`
}

func Read(path string) (*QualityModel, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*QualityModel, error) {
	spec, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	spec.Slug, spec.Path = slug, path
	spec.Location.Path = path
	for i := range spec.QualityAttributes {
		spec.QualityAttributes[i].Location.Path = path
	}
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*QualityModel, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*QualityModel, error) {
	data := input.Source
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &QualityModel{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, Data: append([]byte(nil), data...)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "quality-model", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	fields := map[string]*string{"Description": &spec.Description, "Edition": &spec.Edition, "Sources": &spec.Sources}
	sections := make(map[string]diagnostic.Location)
	for _, section := range doc.Sections {
		target, known := fields[section.Title]
		if !known && section.Title != "Quality attributes" {
			continue
		}
		location := doc.Location(section.Start, section.Title)
		if previous, exists := sections[section.Title]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		sections[section.Title] = location
		if known {
			*target = doc.Body(section)
			continue
		}
		spec.QualityAttributes, err = system.ParseQualityAttributes(doc, section, &ids)
		if err != nil {
			return nil, err
		}
	}
	ids.AddAnchors(doc.Anchors)
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}
