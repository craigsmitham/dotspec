// Package qualitymodel defines reusable quality attribute starters.
package qualitymodel

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

const Suffix = ".quality-model.spec.md"
const Collection = "quality-models"

var Layout = specfile.Layout{Kind: "quality-model", Collection: Collection, Suffix: Suffix}

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

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// QualityModel owns reusable definitions independently of any system.
type QualityModel struct {
	References               *specfile.ReferenceDocument    `json:"-"`
	ID                       string                         `json:"id,omitempty"`
	Slug                     string                         `json:"slug"`
	Name                     string                         `json:"name"`
	Description              string                         `json:"description"`
	Edition                  string                         `json:"edition"`
	Sources                  string                         `json:"sources"`
	QualityAttributes        []system.Attribute             `json:"quality_attributes"`
	QualityAttributesContent *string                        `json:"quality_attributes_content,omitempty"`
	Markdown                 string                         `json:"markdown"`
	Path                     string                         `json:"path"`
	Builtin                  bool                           `json:"builtin"`
	Location                 diagnostic.Location            `json:"location"`
	SectionLocations         map[string]diagnostic.Location `json:"-"`
	// Data retains the full source for copying and reference relocation.
	Data []byte `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid quality-model slug %q: use lowercase letters or digits separated by single hyphens (for example, organization-quality)", slug)
	}
	return nil
}

func Read(path string) (*QualityModel, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.quality-model-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.quality-model-name", err.Error(), diagnostic.Location{Path: path})
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
	spec.Slug, spec.Path = slug, path
	for i := range spec.QualityAttributes {
		spec.QualityAttributes[i].Location.Path = path
	}
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*QualityModel, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &QualityModel{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location), Data: append([]byte(nil), data...)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "quality-model", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	fields := map[string]*string{"Description": &spec.Description, "Edition": &spec.Edition, "Sources": &spec.Sources}
	for _, section := range doc.Sections {
		target, known := fields[section.Title]
		if !known && section.Title != "Quality attributes" {
			continue
		}
		location := doc.Location(section.Start, section.Title)
		if section.Level != 2 {
			return nil, diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-two heading", section.Title), location)
		}
		if previous, exists := spec.SectionLocations[section.Title]; exists {
			return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", section.Title), location, previous)
		}
		spec.SectionLocations[section.Title] = location
		if known {
			*target = doc.Body(section)
			continue
		}
		content := doc.Body(section)
		spec.QualityAttributesContent = &content
		spec.QualityAttributes, err = system.ParseQualityAttributes(doc, section, &ids)
		if err != nil {
			return nil, err
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
