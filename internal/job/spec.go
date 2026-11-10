// Package job reads and discovers jobs to be done.
package job

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/yuin/goldmark/ast"
)

const Suffix = ".job.spec.md"

const Template = `---
type: job
---

# New job

## Job

Describe what someone wants to accomplish, including where this job begins and ends.

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

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
var fieldPattern = regexp.MustCompile(`^\*\*(ID|Stage):\*\*\s*(.*?)\s*$`)
var outcomePattern = regexp.MustCompile(`^\*\*([^*]+):\*\*\s*`)

type Outcome struct {
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
	Outcomes    []Outcome           `json:"outcomes"`
}

type Job struct {
	ID               string                         `json:"id,omitempty"`
	Slug             string                         `json:"slug"`
	Name             string                         `json:"name"`
	Statement        string                         `json:"statement"`
	Performer        string                         `json:"performer"`
	Situation        string                         `json:"situation"`
	Motivation       string                         `json:"motivation"`
	IntendedProgress string                         `json:"intended_progress"`
	ProgressCases    string                         `json:"progress_cases"`
	Evidence         string                         `json:"evidence"`
	Steps            []Step                         `json:"steps"`
	Outcomes         []Outcome                      `json:"outcomes"`
	Markdown         string                         `json:"markdown"`
	Path             string                         `json:"path"`
	System           *system.System                 `json:"system"`
	Location         diagnostic.Location            `json:"location"`
	SectionLocations map[string]diagnostic.Location `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid job slug %q: use lowercase letters or digits separated by single hyphens (for example, develop-requirements)", slug)
	}
	return nil
}

func Read(path string) (*Job, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.job-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.job-name", err.Error(), diagnostic.Location{Path: path})
	}
	if _, err := specfile.Owner(path, "jobs"); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	job, err := Parse(data)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	job.Slug, job.Path = slug, path
	return job, nil
}

// Parse recognizes structured sections but keeps all Markdown for display.
// Missing content is a draft, while conflicting identities and structure are errors.
func Parse(data []byte) (*Job, error) {
	doc, err := specfile.ParseOutline(data, "job")
	if err != nil {
		return nil, err
	}
	job := &Job{Steps: []Step{}, Outcomes: []Outcome{}, ID: doc.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	ids := make(map[string]diagnostic.Location)
	register := func(id, element string, location diagnostic.Location) error {
		if id == "" {
			return nil
		}
		location.Element = element
		if !identifierPattern.MatchString(id) {
			return diagnostic.Error("identity.invalid", fmt.Sprintf("%s: invalid ID %q; use letters, digits, hyphens, or underscores", element, id), location)
		}
		if previous, exists := ids[id]; exists {
			return diagnostic.Error("identity.duplicate", fmt.Sprintf("duplicate ID %q in %s and %s; assign distinct IDs", id, previous.Element, element), location, previous)
		}
		ids[id] = location
		return nil
	}
	if err := register(job.ID, "job", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	sections := make(map[string]*specfile.Section)
	for _, section := range doc.Sections {
		name := section.Title
		if name == "Desired outcome" {
			name = "Intended progress"
		}
		switch name {
		case "Job", "Performer", "Situation", "Motivation", "Intended progress", "Outcomes", "Steps", "Progress cases", "Evidence":
			location := doc.Location(section.Start, name)
			if section.Level != 2 {
				return nil, diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-two heading (## %s)", section.Title, section.Title), location)
			}
			if _, exists := sections[name]; exists {
				return nil, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("duplicate %s section", name), location, job.SectionLocations[name])
			}
			sections[name] = section
			job.SectionLocations[name] = location
		}
	}
	for name, target := range map[string]*string{
		"Job": &job.Statement, "Performer": &job.Performer, "Situation": &job.Situation,
		"Motivation": &job.Motivation, "Intended progress": &job.IntendedProgress,
		"Progress cases": &job.ProgressCases, "Evidence": &job.Evidence,
	} {
		if s := sections[name]; s != nil {
			*target = doc.Body(s)
		}
	}
	if section := sections["Outcomes"]; section != nil {
		job.Outcomes, err = parseOutcomes(doc, section)
		if err != nil {
			return nil, err
		}
		for _, outcome := range job.Outcomes {
			if err := register(outcome.ID, "job outcomes", outcome.Location); err != nil {
				return nil, err
			}
		}
	}
	if section := sections["Steps"]; section != nil {
		for _, s := range section.Children {
			if s.Level != 3 || s.Title == "" {
				return nil, diagnostic.Error("structure.step-heading", "each step must have a nonempty level-three heading (### Step name)", doc.Location(s.Start, s.Title))
			}
			step, err := parseStep(doc, s)
			if err != nil {
				return nil, err
			}
			if err := register(step.ID, "step "+step.Name, step.IDLocation); err != nil {
				return nil, err
			}
			for _, outcome := range step.Outcomes {
				if err := register(outcome.ID, "outcomes for step "+step.Name, outcome.Location); err != nil {
					return nil, err
				}
			}
			job.Steps = append(job.Steps, step)
		}
	}
	for id := range ids {
		doc.Anchors[id] = true
	}
	err = specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors)
	if err != nil {
		return nil, err
	}
	if job.Outcomes == nil {
		job.Outcomes = []Outcome{}
	}
	for i := range job.Steps {
		if job.Steps[i].Outcomes == nil {
			job.Steps[i].Outcomes = []Outcome{}
		}
	}
	return job, nil
}

func parseStep(doc *specfile.Outline, section *specfile.Section) (Step, error) {
	step := Step{Outcomes: []Outcome{}, Name: section.Title, Location: doc.Location(section.Start, section.Title)}
	fields := make(map[string]string)
	fieldLocations := make(map[string]diagnostic.Location)
	fieldLines := make(map[int]bool)
	for _, node := range doc.Nodes(section) {
		if _, ok := node.(*ast.Paragraph); !ok {
			continue
		}
		for i := 0; i < node.Lines().Len(); i++ {
			segment := node.Lines().At(i)
			match := fieldPattern.FindStringSubmatch(strings.TrimSpace(string(segment.Value(doc.Source))))
			if match == nil {
				continue
			}
			if _, exists := fields[match[1]]; exists {
				return step, diagnostic.Error("structure.duplicate-field", fmt.Sprintf("step %s: duplicate %s field", step.Name, match[1]), doc.Location(segment.Start, step.Name+" / "+match[1]), fieldLocations[match[1]])
			}
			fields[match[1]] = match[2]
			fieldLocations[match[1]] = doc.Location(segment.Start, step.Name+" / "+match[1])
			fieldLines[bytes.LastIndexByte(doc.Source[:segment.Start], '\n')+1] = true
		}
	}
	step.ID, step.Stage = fields["ID"], fields["Stage"]
	step.IDLocation = fieldLocations["ID"]
	intro := string(doc.Source[section.ContentStart:doc.IntroEnd(section)])
	// Remove only metadata lines that were recognized in actual paragraphs.
	// Code fences and other examples remain in the original Markdown.
	var description []string
	offset := section.ContentStart
	for _, line := range strings.Split(intro, "\n") {
		if !fieldLines[offset] {
			description = append(description, line)
		}
		offset += len(line) + 1
	}
	step.Description = strings.TrimSpace(strings.Join(description, "\n"))
	seen := false
	var outcomesLocation diagnostic.Location
	for _, child := range section.Children {
		if child.Title != "Outcomes" {
			continue
		}
		if child.Level != 4 {
			return step, diagnostic.Error("structure.section-level", fmt.Sprintf("step %s: outcomes must use #### Outcomes", step.Name), doc.Location(child.Start, step.Name+" / Outcomes"))
		}
		if seen {
			return step, diagnostic.Error("structure.duplicate-section", fmt.Sprintf("step %s: duplicate Outcomes section", step.Name), doc.Location(child.Start, step.Name+" / Outcomes"), outcomesLocation)
		}
		seen = true
		outcomesLocation = doc.Location(child.Start, step.Name+" / Outcomes")
		var err error
		step.Outcomes, err = parseOutcomes(doc, child)
		if err != nil {
			return step, err
		}
	}
	return step, nil
}

func parseOutcomes(doc *specfile.Outline, section *specfile.Section) ([]Outcome, error) {
	var outcomes []Outcome
	for _, node := range doc.Nodes(section) {
		list, ok := node.(*ast.List)
		if !ok {
			continue
		}
		for item := list.FirstChild(); item != nil; item = item.NextSibling() {
			location := doc.Location(item.Pos(), section.Title)
			var paragraphs []string
			for block := item.FirstChild(); block != nil; block = block.NextSibling() {
				switch block.(type) {
				case *ast.Paragraph, *ast.TextBlock:
					paragraphs = append(paragraphs, strings.TrimSpace(string(block.Lines().Value(doc.Source))))
				}
			}
			statement := strings.TrimSpace(strings.Join(paragraphs, "\n\n"))
			outcome := Outcome{Statement: statement, Location: location}
			if match := outcomePattern.FindStringSubmatch(statement); match != nil {
				outcome.ID = match[1]
				outcome.Statement = strings.TrimSpace(strings.TrimPrefix(statement, match[0]))
			}
			outcomes = append(outcomes, outcome)
		}
	}
	return outcomes, nil
}
