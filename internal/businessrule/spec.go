// Package businessrule reads business constraints independently of enforcement.
package businessrule

import (
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

var Layout = specfile.Layout{Kind: "business-rule", Collection: "business-rules", Suffix: ".business-rule.spec.md", SlugExample: "purchase-approval"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
}

const Template = `---
type: business-rule-spec
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

// BusinessRule expresses a business constraint independently of enforcement.
type BusinessRule struct {
	References            *specfile.ReferenceDocument    `json:"-"`
	ID                    string                         `json:"id,omitempty"`
	Slug                  string                         `json:"slug"`
	Name                  string                         `json:"name"`
	Statement             string                         `json:"statement"`
	Applicability         string                         `json:"applicability"`
	Rationale             string                         `json:"rationale"`
	Authority             string                         `json:"authority"`
	Exceptions            string                         `json:"exceptions"`
	RelatedSpecifications string                         `json:"related_specifications"`
	Markdown              string                         `json:"markdown"`
	Path                  string                         `json:"path"`
	Location              diagnostic.Location            `json:"location"`
	SectionLocations      map[string]diagnostic.Location `json:"-"`
}

func Read(path string) (*BusinessRule, error) {
	return readFrom(&specfile.SourceSet{}, path)
}

func readFrom(files *specfile.SourceSet, path string) (*BusinessRule, error) {
	spec, slug, err := specfile.ReadMember(files, Layout, path, ParseMarkdown)
	if err != nil {
		return nil, err
	}
	spec.Slug, spec.Path = slug, path
	spec.Location.Path = path
	return spec, nil
}

// Parse accepts unfinished content, but rejects ambiguous document structure
// and unresolved same-document references. Check commands resolve file links.
func Parse(data []byte) (*BusinessRule, error) {
	return ParseMarkdown(specfile.NewMarkdown(data))
}

func ParseMarkdown(input *specfile.Markdown) (*BusinessRule, error) {
	var meta metadata
	doc, err := input.TypedOutline(Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	spec := &BusinessRule{ID: meta.ID, Name: doc.Name, Markdown: strings.TrimSpace(string(doc.Source)), Location: doc.NameLocation, SectionLocations: make(map[string]diagnostic.Location)}
	var ids specfile.Identifiers
	if err := ids.Register(spec.ID, "business-rule", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	fields := map[string]*string{
		"Statement": &spec.Statement, "Applies to": &spec.Applicability,
		"Rationale": &spec.Rationale, "Authority": &spec.Authority,
		"Exceptions": &spec.Exceptions, "Related specifications": &spec.RelatedSpecifications,
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
