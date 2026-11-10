// Package concept reads and discovers descriptions of purposeful product behavior.
package concept

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

const Suffix = ".concept.spec.md"
const Collection = "concepts"

var Layout = specfile.Layout{Kind: "concept", Collection: Collection, Suffix: Suffix}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: concept
---

# New concept

## Purpose

Describe the benefit this concept provides and why it exists.

## State

Describe the abstract information needed to explain the behavior, independently of its storage representation.

## Actions

Describe or reference the available actions, their inputs and outputs, conditions, and effects.
Distinguish caller assumptions from enforced guards, and state rejection behavior where it matters.
Use subheadings for individual actions when useful.

## Operational principle

Tell a characteristic story showing how the actions and state fulfill the purpose.
The story explains the design; it does not establish exhaustive behavior or verification.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Concept retains section bodies as Markdown. State and actions have no imposed
// language; system association does not imply an implementation boundary.
type Concept struct {
	ID                   string                         `json:"id,omitempty"`
	Slug                 string                         `json:"slug"`
	Name                 string                         `json:"name"`
	Purpose              string                         `json:"purpose"`
	State                string                         `json:"state"`
	Actions              string                         `json:"actions"`
	OperationalPrinciple string                         `json:"operational_principle"`
	Markdown             string                         `json:"markdown"`
	Path                 string                         `json:"path"`
	System               *system.System                 `json:"system"`
	Location             diagnostic.Location            `json:"location"`
	SectionLocations     map[string]diagnostic.Location `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid concept slug %q: use lowercase letters or digits separated by single hyphens (for example, job-map)", slug)
	}
	return nil
}

func Read(path string) (*Concept, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.concept-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.concept-name", err.Error(), diagnostic.Location{Path: path})
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
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Other links remain navigation.
func Parse(data []byte) (*Concept, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &Concept{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "concept", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Purpose": &spec.Purpose, "State": &spec.State,
		"Actions": &spec.Actions, "Operational principle": &spec.OperationalPrinciple,
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
