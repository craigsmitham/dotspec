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
	*obligation.Content
	Criterion   string `json:"criterion"`
	Measurement string `json:"measurement"`
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
	content, criterion, measurement, err := obligation.Parse(data, Layout.Kind, "Criterion", "Measurement")
	if err != nil {
		return nil, err
	}
	return &Requirement{Content: content, Criterion: criterion, Measurement: measurement}, nil
}
