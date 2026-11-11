// Package quality reads and discovers quality obligations.
package quality

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

const Suffix = ".quality.spec.md"
const Collection = "quality"

var Layout = specfile.Layout{Kind: "quality", Collection: Collection, Suffix: Suffix}

const Template = `---
type: quality
---

# New quality requirement

## Applies to

Identify the subject and operating conditions, workload, population, and exclusions.
For outcomes of use, identify the users, their task, available assistance, and environment; link relevant actors, jobs, or use cases.

## Statement

State one coherent quality obligation. What observable result would violate it?

## Criterion

Define or reference what counts as satisfaction. For a measured property, identify the measure, required bound, units, aggregation, and observation window as needed. Keep each required bound in one place; reference it if already stated above or in an authoritative source.
If a bound is unknown, record an open question and the evidence needed to choose it. Qualitative criteria are valid; keep preferred targets distinct from required bounds.

## Rationale

Explain why this quality matters and any relevant tradeoffs.

## Verification plan

Describe how to obtain and compare evidence against the statement and criterion. Identify the realization, method, assessment conditions, sampling where needed, and evidence to retain. Explain how assessment conditions represent the operating context in Applies to. Reference required bounds rather than repeating them. This plan is not a completed verification result.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Requirement keeps obligations as authored Markdown. Classification resolves
// against the associated system without establishing applicability or satisfaction.
type Requirement struct {
	*obligation.Content
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

// Parse interprets the quality form without filesystem access.
func Parse(data []byte) (*Requirement, error) {
	content, err := obligation.Parse(data, Layout.Kind)
	if err != nil {
		return nil, err
	}
	return &Requirement{Content: content}, nil
}
