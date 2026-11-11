// Package businessrule reads business constraints independently of enforcement.
package businessrule

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

const Suffix = ".business-rule.spec.md"
const Collection = "business-rules"

var Layout = specfile.Layout{Kind: "business-rule", Collection: Collection, Suffix: Suffix}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: business-rule
---

# New business rule

## Statement

State the business constraint independently of how it is enforced.

## Applies to

Describe the business activity, subjects, and applicable conditions.

## Rationale

Explain why this rule exists and the business objective it serves.

## Authority

Identify the policy, decision, agreement, or other source establishing the rule.

## Exceptions

Describe explicit qualifications or link to exception rules, if any are known.

## Related specifications

Link to software requirements or use cases that address this rule. Describe their role;
a link alone does not establish enforcement or compliance.
`

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// BusinessRule expresses a business constraint independently of enforcement.
type BusinessRule struct {
	References            *specfile.ReferenceDocument    `json:"-"`
	ID                    string                         `json:"id,omitempty"`
	Slug                  string                         `json:"slug"`
	Name                  string                         `json:"name"`
	Statement             string                         `json:"statement"`
	AppliesTo             string                         `json:"applies_to"`
	Rationale             string                         `json:"rationale"`
	Authority             string                         `json:"authority"`
	Exceptions            string                         `json:"exceptions"`
	RelatedSpecifications string                         `json:"related_specifications"`
	Markdown              string                         `json:"markdown"`
	Path                  string                         `json:"path"`
	System                *system.System                 `json:"system"`
	Location              diagnostic.Location            `json:"location"`
	SectionLocations      map[string]diagnostic.Location `json:"-"`
}

func ValidateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid business-rule slug %q: use lowercase letters or digits separated by single hyphens (for example, job-map)", slug)
	}
	return nil
}

func Read(path string) (*BusinessRule, error) {
	name := filepath.Base(path)
	if !strings.HasSuffix(name, Suffix) {
		return nil, diagnostic.Error("layout.business-rule-name", fmt.Sprintf("expected a <slug>%s filename", Suffix), diagnostic.Location{Path: path})
	}
	slug := strings.TrimSuffix(name, Suffix)
	if err := ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout.business-rule-name", err.Error(), diagnostic.Location{Path: path})
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
func Parse(data []byte) (*BusinessRule, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &BusinessRule{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "business-rule", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Statement": &spec.Statement, "Applies to": &spec.AppliesTo,
		"Rationale": &spec.Rationale, "Authority": &spec.Authority,
		"Exceptions": &spec.Exceptions, "Related specifications": &spec.RelatedSpecifications,
	}

	for _, section := range doc.Sections {
		if section.Title == "Authority and evidence" {
			section.Title = "Authority"
		}
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
