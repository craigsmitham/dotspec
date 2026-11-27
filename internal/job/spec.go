// Package job reads and discovers jobs to be done.
package job

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var Layout = specfile.Layout{Kind: "job", Collection: "jobs", Suffix: ".job.spec.md", SlugExample: "develop-requirements"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: job-spec
---

# New job

## Job

Describe what someone seeks to accomplish or change, independent of the chosen solution.

## Performer

Describe who is trying to accomplish this job.

## Situation

Describe the circumstances in which this job arises.

## Motivation

Describe why someone wants to act in this situation.

## Intended progress

Describe the progress or result they want to achieve.
Use an optional Outcomes section for individual success criteria.
`

type DesiredOutcome struct {
	Name      string              `json:"name,omitempty"`
	AppliesTo string              `json:"applies_to,omitempty"`
	ID        string              `json:"id,omitempty"`
	Statement string              `json:"statement"`
	Location  diagnostic.Location `json:"location"`
}

type Step struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Stage       string              `json:"stage,omitempty"`
	Description string              `json:"description"`
	Location    diagnostic.Location `json:"location"`
	IDLocation  diagnostic.Location `json:"-"`
	Outcomes    []DesiredOutcome    `json:"outcomes"`
}

type Job struct {
	References       *specfile.ReferenceDocument    `json:"-"`
	ID               string                         `json:"id,omitempty"`
	Slug             string                         `json:"slug"`
	Name             string                         `json:"name"`
	Statement        string                         `json:"statement"`
	Performer        string                         `json:"performer"`
	Situation        string                         `json:"situation"`
	Motivation       string                         `json:"motivation"`
	IntendedProgress string                         `json:"intended_progress"`
	Evidence         string                         `json:"evidence"`
	Steps            []Step                         `json:"steps"`
	Outcomes         []DesiredOutcome               `json:"outcomes"`
	Markdown         string                         `json:"markdown"`
	Path             string                         `json:"path"`
	Location         diagnostic.Location            `json:"location"`
	SectionLocations map[string]diagnostic.Location `json:"-"`
}

func Read(path string) (*Job, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*Job, error) {
	job, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	job.Slug, job.Path = slug, path
	job.Location.Path = path
	for i := range job.Outcomes {
		job.Outcomes[i].Location.Path = path
	}
	for i := range job.Steps {
		job.Steps[i].Location.Path = path
		for j := range job.Steps[i].Outcomes {
			job.Steps[i].Outcomes[j].Location.Path = path
		}
	}
	return job, nil
}

// Parse recognizes structured sections but keeps all Markdown for display.
// Missing content is a draft, while conflicting identities and structure are errors.
func Parse(data []byte) (*Job, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*Job, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	job := &Job{Steps: []Step{}, Outcomes: []DesiredOutcome{}, ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(job.ID, "job", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	// Outcomes and Steps are parsed from their sections below. Capturing them
	// too rejects every duplicate section, in document order, before that parsing.
	fields := map[string]*string{
		"Job": &job.Statement, "Performer": &job.Performer, "Situation": &job.Situation,
		"Motivation": &job.Motivation, "Intended progress": &job.IntendedProgress,
		"Evidence": &job.Evidence, "Outcomes": new(string), "Steps": new(string),
	}
	if err := doc.Capture(fields, job.SectionLocations); err != nil {
		return nil, err
	}
	sections := make(map[string]*specfile.Section)
	for _, section := range doc.Sections {
		sections[section.Title] = section
	}
	if section := sections["Outcomes"]; section != nil {
		job.Outcomes, err = parseOutcomes(doc, section)
		if err != nil {
			return nil, err
		}
		for _, outcome := range job.Outcomes {
			if err := ids.Register(outcome.ID, "job outcomes", outcome.Location); err != nil {
				return nil, err
			}
		}
	}
	if section := sections["Steps"]; section != nil {
		for _, s := range section.Children {
			step, err := parseStep(doc, s)
			if err != nil {
				return nil, err
			}
			if err := ids.Register(step.ID, "step "+step.Name, step.IDLocation); err != nil {
				return nil, err
			}
			for _, outcome := range step.Outcomes {
				if err := ids.Register(outcome.ID, "outcomes for step "+step.Name, outcome.Location); err != nil {
					return nil, err
				}
			}
			job.Steps = append(job.Steps, step)
		}
	}
	ids.AddAnchors(doc.Anchors)
	err = doc.ValidateReferences()
	if err != nil {
		return nil, err
	}
	job.References = doc.References()
	return job, nil
}

func parseStep(doc *specfile.Outline, section *specfile.Section) (Step, error) {
	step := Step{Outcomes: []DesiredOutcome{}, Name: section.Title, Location: doc.Location(section.Start, section.Title)}
	fields, err := doc.ReadFields("step "+step.Name, section, "ID", "Stage")
	if err != nil {
		return step, err
	}
	step.ID, step.Stage = fields.Values["ID"], fields.Values["Stage"]
	step.IDLocation = fields.Locations["ID"]
	step.Description = fields.Strip(doc, section.ContentStart, doc.IntroEnd(section))
	seen := false
	var outcomesLocation diagnostic.Location
	for _, child := range section.Children {
		if child.Title != "Outcomes" {
			continue
		}
		if seen {
			return step, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("step %s: duplicate Outcomes section", step.Name), doc.Location(child.Start, step.Name+" / Outcomes"), outcomesLocation)
		}
		seen = true
		outcomesLocation = doc.Location(child.Start, step.Name+" / Outcomes")
		step.Outcomes, err = parseOutcomes(doc, child)
		if err != nil {
			return step, err
		}
	}
	return step, nil
}

func parseOutcomes(doc *specfile.Outline, section *specfile.Section) ([]DesiredOutcome, error) {
	outcomes := []DesiredOutcome{}
	for _, child := range section.Children {
		fields, err := doc.ReadFields("outcome "+child.Title, child, "ID")
		if err != nil {
			return nil, err
		}
		outcome := DesiredOutcome{Name: child.Title, ID: fields.Values["ID"], Statement: fields.Strip(doc, child.ContentStart, doc.IntroEnd(child)), Location: doc.Location(child.Start, child.Title)}
		seenAppliesTo := false
		for _, sub := range child.Children {
			if sub.Title == "Applies to" {
				if seenAppliesTo {
					return nil, diagnostic.Error("structure.duplicate-section", "duplicate outcome Applies to section", doc.Location(sub.Start, child.Title))
				}
				seenAppliesTo = true
				outcome.AppliesTo = doc.Body(sub)
			}
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}
