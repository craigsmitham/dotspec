// Package requirement reads and discovers requirement specifications.
package requirement

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/obligation"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Suffix = ".requirement.spec.md"
const Collection = "requirements"

var Layout = specfile.Layout{Kind: "requirement", Collection: Collection, Suffix: Suffix}

const Template = `---
type: requirement-spec
---

# New requirement

## Applies to

Identify the subject and operating conditions, including any exclusions. Distinguish assumptions from conditions the system must accommodate. For outcomes of use, identify the users, task, and environment.

## Statement

State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.
What observable result would violate this obligation? Add a Criterion section if an explicit assessment rule helps, defining or referencing each required bound once.

## Rationale

Explain why this obligation is needed and reference its source or governing goal.

## Verification plan

Describe how to check a system or process realization against the statement and any criterion. Identify the method, assessment conditions, and evidence to retain; explain how the conditions represent the operating context in Applies to. Reference required bounds rather than repeating them. This plan is not a completed verification result.
`

// QualityTemplate is the quality preset. It prompts for an explicit criterion;
// the result is an ordinary requirement, classified by its quality attributes.
const QualityTemplate = `---
type: requirement-spec
---

# New quality requirement

## Applies to

Identify the subject and operating conditions, workload, population, and exclusions.
For outcomes of use, identify the users, their task, available assistance, and environment; link relevant actors, jobs, or use cases.

## Statement

State one coherent quality obligation. What observable result would violate it?

## Criterion

Define or reference what counts as satisfaction. For a measured property, state the measure (including units, aggregation, and observation window as needed) and the required bound. Keep each required bound in one place; reference it if already stated above or in an authoritative source.
If a bound is unknown, record an open question and the evidence needed to choose it. Qualitative criteria are valid; keep preferred targets distinct from required bounds.

## Rationale

Explain why this quality matters and any relevant tradeoffs.

## Verification plan

Describe how to obtain and compare evidence against the statement and criterion. Identify the realization, method (including measurement procedure and sampling where needed), assessment conditions, and evidence to retain. Explain how assessment conditions represent the operating context in Applies to. Reference required bounds rather than repeating them. This plan is not a completed verification result.
`

// Presets guide authoring; they do not create another specification type.
func ValidatePreset(preset string) error {
	switch preset {
	case "general", "quality":
		return nil
	}
	return fmt.Errorf("unknown requirement preset %q: choose general or quality", preset)
}

func TemplateFor(preset string) (string, error) {
	if err := ValidatePreset(preset); err != nil {
		return "", err
	}
	if preset == "quality" {
		return QualityTemplate, nil
	}
	return Template, nil
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Requirement keeps obligations as authored Markdown. Classification resolves
// against the associated system without establishing applicability or satisfaction.
type Requirement struct {
	*obligation.Content
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid requirement slug %q: use lowercase letters or digits separated by single hyphens (for example, reject-existing-target)", slug)
	}
	return nil
}

func Read(path string) (*Requirement, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.requirement-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.requirement-name", err.Error(), diagnostic.Location{Path: path})
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

// Parse interprets the requirement form without filesystem access.
func Parse(data []byte) (*Requirement, error) {
	content, err := obligation.Parse(data, Layout.Kind)
	if err != nil {
		return nil, err
	}
	return &Requirement{Content: content}, nil
}
