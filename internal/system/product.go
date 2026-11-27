package system

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var FeatureLayout = specfile.Layout{Kind: "feature", Collection: "features", Suffix: ".feature.spec.md", SlugExample: "resume-work"}
var CapabilityLayout = specfile.Layout{Kind: "capability", Collection: "capabilities", Suffix: ".capability.spec.md", SlugExample: "local-storage"}

type ProductSpec struct {
	Name        string                      `json:"name"`
	ID          string                      `json:"id,omitempty"`
	Description string                      `json:"description,omitempty"`
	Path        string                      `json:"path"`
	Markdown    string                      `json:"markdown"`
	References  *specfile.ReferenceDocument `json:"-"`
}

func ParseProductSpec(kind string, input *specfile.Markdown) (*ProductSpec, error) {
	var meta metadata
	doc, err := input.TypedOutline(kind, &meta)
	if err != nil {
		return nil, err
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, kind, input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	spec := &ProductSpec{Name: doc.Name, ID: meta.ID, Markdown: strings.TrimSpace(string(doc.Source)), References: doc.References()}
	seen := map[string]diagnostic.Location{}
	for _, s := range doc.Sections {
		switch s.Title {
		case "Description", "Purpose", "Overview", "Groups", "Enablers", "Supports", "Depends on", "Provided by", "Specification system":
			location := doc.Location(s.Start, s.Title)
			if prev, ok := seen[s.Title]; ok {
				return nil, diagnostic.Error("structure.duplicate-section", "duplicate "+s.Title+" section", location, prev)
			}
			seen[s.Title] = location
			if s.Title == "Description" {
				spec.Description = doc.Body(s)
			}
		}
	}
	return spec, nil
}

func ReadProductSpec(layout specfile.Layout, path string) (*ProductSpec, error) {
	return readProduct(&specfile.SourceSet{}, layout, path)
}

func readProduct(files *specfile.SourceSet, layout specfile.Layout, path string) (*ProductSpec, error) {
	spec, _, err := specfile.ReadMember(files, layout, path, func(input *specfile.Markdown) (*ProductSpec, error) { return ParseProductSpec(layout.Kind, input) })
	if err != nil {
		return nil, err
	}
	spec.Path = path
	return spec, nil
}

func ProductCheckType(layout specfile.Layout) check.Type {
	read := func(files *specfile.SourceSet, path string) (*ProductSpec, error) {
		return readProduct(files, layout, path)
	}
	kind := check.ReadType(layout, read, func(spec *ProductSpec) *specfile.ReferenceDocument { return spec.References }, nil)
	kind.Context = CheckContext
	return kind
}

// InitProduct creates a feature or capability draft named by slug.
func InitProduct(layout specfile.Layout, path, slug string) (string, error) {
	sections := "## Description\n\nDescribe what a stakeholder can do.\n\n## Groups\n\nLink to the specifications that define this behavior."
	if layout == CapabilityLayout {
		sections = "## Description\n\nDescribe the ability and what it enables.\n\n## Provided by\n\nLink to the systems or components that provide this ability."
	}
	source := fmt.Sprintf("---\ntype: %s\n---\n\n# %s\n\n%s\n", specfile.AuthoredType(layout.Kind), specfile.SlugTitle(slug), sections)
	return layout.Create(path, slug, source)
}
