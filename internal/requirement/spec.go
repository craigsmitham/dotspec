// Package requirement reads and discovers general obligations.
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
type: requirement
---

# New requirement

## Applies to

Identify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.

## Statement

State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.

## Rationale

Explain why this obligation is needed and reference its source or governing goal.

## Verification

Describe how to establish that the obligation holds, including success criteria and evidence to retain.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Requirement keeps obligations as authored Markdown. Classification resolves
// against the associated system without establishing applicability or satisfaction.
type Requirement struct {
	Type string `json:"type"`
	*obligation.Content
	Statement    string `json:"statement"`
	Verification string `json:"verification"`
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

// Parse interprets the general requirement form without filesystem access.
func Parse(data []byte) (*Requirement, error) {
	content, statement, verification, err := obligation.Parse(data, Layout.Kind, "Statement", "Verification")
	if err != nil {
		return nil, err
	}
	return &Requirement{Type: Layout.Kind, Content: content, Statement: statement, Verification: verification}, nil
}
