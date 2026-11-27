// Package usecase reads and discovers actor-goal use cases.
package usecase

import (
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

var Layout = specfile.Layout{Kind: "use-case", Collection: "use-cases", Suffix: ".use-case.spec.md", SlugExample: "check-specifications"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

// Template starts a readable sketch; expanded prompts are optional elaboration.
const Template = `---
type: use-case-spec
---

# New use case

## Goal

Name what the primary actor wants to accomplish, using an active verb phrase.

## System boundary

Name the system whose external behavior this use case describes.

## Primary actor

Identify the role pursuing the goal.

## Goal level

Identify the level: summary (kite), user task (sea), or subfunction (fish).

## Main success scenario

Describe the typical successful interaction in a short paragraph or a few steps.
Keep internal implementation and interface details in their own sources.

## Extensions

### An alternative path

#### Condition

Describe when this alternative arises.

#### Handling

Describe the response and where the interaction continues or ends.

`

const ExpandedTemplate = `---
type: use-case-spec
---

# New use case

## Goal

Describe what the primary actor wants to accomplish through this interaction.

## System boundary

Name the system whose behavior this use case describes.

## Primary actor

Describe the role pursuing the goal.

## Goal level

Identify the level: summary, user task, or subfunction.

## Trigger

Describe the action or time event that starts the interaction.

## Preconditions

Describe the state expected to be true before the interaction, or state that none is assumed.

## Supporting actors

Name external collaborators and the services or notifications exchanged with them.

## Stakeholder interests

Identify interests that affect the promised behavior, including those of nonparticipants.

## Minimal guarantees

State or link to what remains protected even when the actor's goal is not achieved.

## Success guarantees

State or link to what holds when the actor's goal is achieved.

## Main success scenario

1. Describe the actor's action.
2. Describe the system's response and how the actor reaches the goal.

## Extensions

### An alternative path

#### Condition

Describe when this alternative arises.

#### Handling

Describe the response and where the interaction continues or ends.

## Sequence

State whether the scenario's order is representative or required; link to the governing rule for required order.

## Supports

Link to the jobs, steps, desired outcomes, or higher-level use cases this use case supports.
`

// UseCase retains section bodies as Markdown for the authoring commands. Model
// export separately exposes source-backed scenario and extension records. File
// ownership does not set the system boundary.
type UseCase struct {
	References           *specfile.ReferenceDocument    `json:"-"`
	ID                   string                         `json:"id,omitempty"`
	Slug                 string                         `json:"slug"`
	Name                 string                         `json:"name"`
	Goal                 string                         `json:"goal"`
	BoundaryDescription  string                         `json:"boundary_description"`
	PrimaryActor         string                         `json:"primary_actor"`
	GoalLevel            string                         `json:"goal_level"`
	Trigger              string                         `json:"trigger"`
	Preconditions        string                         `json:"preconditions"`
	SuccessGuarantees    string                         `json:"success_guarantees"`
	MainSuccessScenario  string                         `json:"main_success_scenario"`
	Extensions           string                         `json:"extensions"`
	Sequence             string                         `json:"sequence"`
	SupportingActors     string                         `json:"supporting_actors"`
	StakeholderInterests string                         `json:"stakeholder_interests"`
	MinimalGuarantees    string                         `json:"minimal_guarantees"`
	Supports             string                         `json:"supports,omitempty"`
	Markdown             string                         `json:"markdown"`
	Path                 string                         `json:"path"`
	System               *system.System                 `json:"system"`
	Location             diagnostic.Location            `json:"location"`
	SectionLocations     map[string]diagnostic.Location `json:"-"`
}

func Read(path string) (*UseCase, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*UseCase, error) {
	spec, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	spec.Slug, spec.Path = slug, path
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*UseCase, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*UseCase, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &UseCase{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "use case", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Goal": &spec.Goal, "System boundary": &spec.BoundaryDescription,
		"Primary actor": &spec.PrimaryActor, "Goal level": &spec.GoalLevel,
		"Trigger": &spec.Trigger, "Preconditions": &spec.Preconditions,
		"Supporting actors": &spec.SupportingActors, "Stakeholder interests": &spec.StakeholderInterests,
		"Minimal guarantees": &spec.MinimalGuarantees,
		"Success guarantees": &spec.SuccessGuarantees, "Main success scenario": &spec.MainSuccessScenario,
		"Extensions": &spec.Extensions, "Sequence": &spec.Sequence, "Supports": &spec.Supports,
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
