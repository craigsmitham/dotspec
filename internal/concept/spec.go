// Package concept reads and discovers descriptions of purposeful product behavior.
package concept

import (
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

var Layout = specfile.Layout{Kind: "concept", Collection: "concepts", Suffix: ".concept.spec.md", SlugExample: "job-map"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: concept-spec
---

# New concept

## Purpose

Describe the benefit this concept provides and why it exists.

## State

Describe the abstract information needed to explain the behavior, independently of its storage representation.

## Actions

### Perform an action

Describe the action, its inputs, outputs, conditions and effects.
Distinguish caller assumptions from enforced guards, and state rejection behavior where it matters.

## Operational principles

Tell a characteristic story showing how the actions and state fulfill the purpose.
Use multiple stories when needed to explain the essential behavior.
The story explains the design; it does not establish exhaustive behavior or verification.
`

// Concept retains section bodies as Markdown. State and actions have no imposed
// language; system association does not imply an implementation boundary.
type Concept struct {
	EffectsTable          string                         `json:"effects_table,omitempty"`
	References            *specfile.ReferenceDocument    `json:"-"`
	ID                    string                         `json:"id,omitempty"`
	Slug                  string                         `json:"slug"`
	Name                  string                         `json:"name"`
	Purpose               string                         `json:"purpose"`
	State                 string                         `json:"state"`
	Actions               string                         `json:"actions"`
	OperationalPrinciples string                         `json:"operational_principles"`
	Markdown              string                         `json:"markdown"`
	Path                  string                         `json:"path"`
	Location              diagnostic.Location            `json:"location"`
	SectionLocations      map[string]diagnostic.Location `json:"-"`
}

func Read(path string) (*Concept, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*Concept, error) {
	spec, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	spec.Slug, spec.Path = slug, path
	spec.Location.Path = path
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*Concept, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*Concept, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &Concept{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "concept", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Effects table": &spec.EffectsTable,
		"Purpose":       &spec.Purpose, "State": &spec.State,
		"Actions": &spec.Actions, "Operational principles": &spec.OperationalPrinciples,
	}
	if err := doc.Capture(fields, spec.SectionLocations); err != nil {
		return nil, err
	}
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	spec.References = doc.References()
	return spec, nil
}
