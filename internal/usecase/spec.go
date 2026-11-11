// Package usecase reads and discovers actor-goal use cases.
package usecase

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

const Suffix = ".use-case.spec.md"
const Collection = "use-cases"

var Layout = specfile.Layout{Kind: "use-case", Collection: Collection, Suffix: Suffix}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

// Template starts a readable sketch; expanded prompts are optional elaboration.
const Template = `---
type: use-case
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

List alternative success paths and failure conditions before resolving their handling.
Mark unknown responses as unresolved; expand the branches needed for the next decision.
`

const ExpandedTemplate = `---
type: use-case
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

## Minimum guarantees

State or link to what remains protected even when the actor's goal is not achieved.

## Success criterion

State or link to the observable condition that completes the actor's goal.

## Main success scenario

1. Describe the actor's action.
2. Describe the system's response and how the actor reaches the goal.

## Extensions

Describe alternative success paths and failure conditions. For each developed branch,
identify where its condition arises, the response, and whether the interaction resumes,
succeeds another way, or ends without achieving the goal. Link to governing rules.
Condition-only entries with explicitly unresolved handling are useful drafts.
Use descriptive subheadings when other sources need to link to a branch.

## Sequence

State whether the scenario's order is representative or required; link to the governing rule for required order.

## Related jobs

Link to the jobs, steps, or desired outcomes this use case supports.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// UseCase retains section bodies as Markdown for the authoring commands. Model
// export separately exposes source-backed scenario and extension records. File
// ownership does not set the system boundary.
type UseCase struct {
	ValidationPlan       string                         `json:"validation_plan,omitempty"`
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
	SuccessCriterion     string                         `json:"success_criterion"`
	MainSuccessScenario  string                         `json:"main_success_scenario"`
	Extensions           string                         `json:"extensions"`
	Sequence             string                         `json:"sequence"`
	SupportingActors     string                         `json:"supporting_actors"`
	StakeholderInterests string                         `json:"stakeholder_interests"`
	MinimumGuarantees    string                         `json:"minimum_guarantees"`
	RelatedJobs          string                         `json:"related_jobs"`
	Supports             string                         `json:"supports,omitempty"`
	Markdown             string                         `json:"markdown"`
	Path                 string                         `json:"path"`
	System               *system.System                 `json:"system"`
	Location             diagnostic.Location            `json:"location"`
	SectionLocations     map[string]diagnostic.Location `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid use-case slug %q: use lowercase letters or digits separated by single hyphens (for example, check-specifications)", slug)
	}
	return nil
}

func Read(path string) (*UseCase, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.use-case-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.use-case-name", err.Error(), diagnostic.Location{Path: path})
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
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*UseCase, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &UseCase{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "use case", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterObligations(doc); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{"Validation plan": &spec.ValidationPlan,
		"Goal": &spec.Goal, "System boundary": &spec.BoundaryDescription,
		"Primary actor": &spec.PrimaryActor, "Goal level": &spec.GoalLevel,
		"Trigger": &spec.Trigger, "Preconditions": &spec.Preconditions,
		"Supporting actors": &spec.SupportingActors, "Stakeholder interests": &spec.StakeholderInterests,
		"Minimum guarantees": &spec.MinimumGuarantees,
		"Success criterion":  &spec.SuccessCriterion, "Main success scenario": &spec.MainSuccessScenario,
		"Extensions": &spec.Extensions, "Sequence": &spec.Sequence, "Related jobs": &spec.RelatedJobs, "Supports": &spec.Supports,
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
	spec.References = doc.References()
	return spec, nil
}
