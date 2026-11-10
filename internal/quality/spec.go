// Package quality reads and discovers quality obligations.
package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"go.yaml.in/yaml/v3"
)

const Suffix = ".quality.spec.md"
const Collection = "quality"

var Layout = specfile.Layout{Kind: "quality", Collection: Collection, Suffix: Suffix}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string    `yaml:"id"`
	Attributes        yaml.Node `yaml:"attributes"`
}

const Template = `---
type: quality
---

# New quality requirement

## Applies to

Identify the subject and the conditions, workload, or population covered by this requirement.

## Criterion

State one coherent obligation, including its measure, comparison, and target where applicable.
If a target is unknown, record the open decision and the evidence needed to choose it.

## Measurement

Describe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.

## Rationale

Explain why this quality matters and any relevant tradeoffs.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Requirement keeps obligations as authored Markdown. Classification resolves
// against the associated system without establishing applicability or satisfaction.
type Requirement struct {
	ID                 string                         `json:"id,omitempty"`
	Slug               string                         `json:"slug"`
	Name               string                         `json:"name"`
	AppliesTo          string                         `json:"applies_to"`
	Criterion          string                         `json:"criterion"`
	Measurement        string                         `json:"measurement"`
	Rationale          string                         `json:"rationale"`
	Attributes         []string                       `json:"attributes"`
	Classification     Classification                 `json:"classification"`
	Markdown           string                         `json:"markdown"`
	Path               string                         `json:"path"`
	System             *system.System                 `json:"system"`
	Location           diagnostic.Location            `json:"location"`
	SectionLocations   map[string]diagnostic.Location `json:"-"`
	AttributeLocations map[string]diagnostic.Location `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid quality slug %q: use lowercase letters or digits separated by single hyphens (for example, diagnostic-severity)", slug)
	}
	return nil
}

func Read(path string) (*Requirement, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.quality-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.quality-name", err.Error(), diagnostic.Location{Path: path})
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
	spec.Location.Path = path
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Other links remain navigation.
func Parse(data []byte) (*Requirement, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &Requirement{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location), Attributes: []string{}, AttributeLocations: make(map[string]diagnostic.Location)}
	if err := parseAttributes(spec, &meta.Attributes); err != nil {
		return nil, err
	}
	spec.Classification = Classify(spec.Attributes, nil)
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "quality", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Applies to": &spec.AppliesTo, "Criterion": &spec.Criterion,
		"Measurement": &spec.Measurement, "Rationale": &spec.Rationale,
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
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	return spec, nil
}

func parseAttributes(spec *Requirement, node *yaml.Node) error {
	if node.Kind == 0 {
		return nil
	}
	location := diagnostic.Location{Line: node.Line + 1, Column: node.Column, Element: "attributes"}
	if node.Kind != yaml.SequenceNode {
		return diagnostic.Error("quality.attributes", "attributes must be a list of unique, nonempty string IDs", location)
	}
	var ids specfile.Identifiers
	for _, value := range node.Content {
		location := diagnostic.Location{Line: value.Line + 1, Column: value.Column, Element: "attributes"}
		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || strings.TrimSpace(value.Value) == "" {
			return diagnostic.Error("quality.attributes", "each attribute must be a nonempty string ID", location)
		}
		if err := ids.Register(value.Value, "attribute", location); err != nil {
			return err
		}
		spec.Attributes = append(spec.Attributes, value.Value)
		spec.AttributeLocations[value.Value] = location
	}
	return nil
}
