// Package job reads and discovers jobs to be done.
package job

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

## Desired outcome

Describe the progress or result they want to achieve.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
var fieldPattern = regexp.MustCompile(`^\*\*(ID|Stage):\*\*\s*(.*?)\s*$`)
var outcomePattern = regexp.MustCompile(`^\*\*([^*]+):\*\*\s*`)

type Outcome struct {
	ID        string
	Statement string
}

type Step struct {
	ID, Name, Stage, Description string
	Outcomes                     []Outcome
}

type Job struct {
	ID, Slug, Name                                              string
	Statement, Performer, Situation, Motivation, DesiredOutcome string
	ProgressCases, Evidence                                     string
	Steps                                                       []Step
	Outcomes                                                    []Outcome
	Markdown                                                    string
	Path                                                        string
	System                                                      *system.System
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
		return nil, fmt.Errorf("%s: expected a <slug>%s filename", path, Suffix)
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	job, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
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
	job := &Job{ID: doc.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source))}
	ids := make(map[string]string)
	register := func(id, location string) error {
		if id == "" {
			return nil
		}
		if !identifierPattern.MatchString(id) {
			return fmt.Errorf("%s: invalid ID %q; use letters, digits, hyphens, or underscores", location, id)
		}
		if previous, exists := ids[id]; exists {
			return fmt.Errorf("duplicate ID %q in %s and %s", id, previous, location)
		}
		ids[id] = location
		return nil
	}
	if err := register(job.ID, "job"); err != nil {
		return nil, err
	}
	sections := make(map[string]*specfile.Section)
	for _, section := range doc.Sections {
		switch section.Title {
		case "Job", "Performer", "Situation", "Motivation", "Desired outcome", "Outcomes", "Steps", "Progress cases", "Evidence":
			if section.Level != 2 {
				return nil, fmt.Errorf("%s must use a level-two heading (## %s)", section.Title, section.Title)
			}
			if _, exists := sections[section.Title]; exists {
				return nil, fmt.Errorf("duplicate %s section", section.Title)
			}
			sections[section.Title] = section
		}
	}
	for name, target := range map[string]*string{
		"Job": &job.Statement, "Performer": &job.Performer, "Situation": &job.Situation,
		"Motivation": &job.Motivation, "Desired outcome": &job.DesiredOutcome,
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
			if err := register(outcome.ID, "job outcomes"); err != nil {
				return nil, err
			}
		}
	}
	if section := sections["Steps"]; section != nil {
		for _, s := range section.Children {
			if s.Level != 3 || s.Title == "" {
				return nil, fmt.Errorf("each step must have a nonempty level-three heading (### Step name)")
			}
			step, err := parseStep(doc, s)
			if err != nil {
				return nil, err
			}
			if err := register(step.ID, "step "+step.Name); err != nil {
				return nil, err
			}
			for _, outcome := range step.Outcomes {
				if err := register(outcome.ID, "outcomes for step "+step.Name); err != nil {
					return nil, err
				}
			}
			job.Steps = append(job.Steps, step)
		}
	}
	for id := range ids {
		doc.Anchors[id] = true
	}
	err = ast.Walk(doc.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}
		destination := string(link.Destination)
		if !strings.HasPrefix(destination, "#") || destination == "#" {
			return ast.WalkContinue, nil
		}
		fragment, err := url.PathUnescape(destination[1:])
		if err != nil {
			return ast.WalkStop, fmt.Errorf("invalid local reference %q", destination)
		}
		if !doc.Anchors[fragment] {
			return ast.WalkStop, fmt.Errorf("unresolved local reference %q; use a declared ID or a Markdown heading anchor", destination)
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func parseStep(doc *specfile.Outline, section *specfile.Section) (Step, error) {
	step := Step{Name: section.Title}
	fields := make(map[string]string)
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
				return step, fmt.Errorf("step %s: duplicate %s field", step.Name, match[1])
			}
			fields[match[1]] = match[2]
			fieldLines[bytes.LastIndexByte(doc.Source[:segment.Start], '\n')+1] = true
		}
	}
	step.ID, step.Stage = fields["ID"], fields["Stage"]
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
	for _, child := range section.Children {
		if child.Title != "Outcomes" {
			continue
		}
		if child.Level != 4 {
			return step, fmt.Errorf("step %s: outcomes must use #### Outcomes", step.Name)
		}
		if seen {
			return step, fmt.Errorf("step %s: duplicate Outcomes section", step.Name)
		}
		seen = true
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
			var paragraphs []string
			for block := item.FirstChild(); block != nil; block = block.NextSibling() {
				switch block.(type) {
				case *ast.Paragraph, *ast.TextBlock:
					paragraphs = append(paragraphs, strings.TrimSpace(string(block.Lines().Value(doc.Source))))
				}
			}
			statement := strings.TrimSpace(strings.Join(paragraphs, "\n\n"))
			outcome := Outcome{Statement: statement}
			if match := outcomePattern.FindStringSubmatch(statement); match != nil {
				outcome.ID = match[1]
				outcome.Statement = strings.TrimSpace(strings.TrimPrefix(statement, match[0]))
			}
			outcomes = append(outcomes, outcome)
		}
	}
	return outcomes, nil
}
