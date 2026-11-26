package system

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var FeatureLayout = specfile.Layout{Kind: "feature", Collection: "features", Suffix: ".feature.spec.md"}
var CapabilityLayout = specfile.Layout{Kind: "capability", Collection: "capabilities", Suffix: ".capability.spec.md"}

func rejectInlineSubject(doc *specfile.Outline, section *specfile.Section) error {
	if len(section.Children) == 0 {
		return nil
	}
	var kind string
	switch section.Title {
	case "Components":
		kind = "component"
	case "Features":
		kind = "feature"
	case "Capabilities":
		kind = "capability"
	default:
		return nil
	}
	return diagnostic.Error("structure."+kind+"-source", fmt.Sprintf("declare each %s in its own <slug>.%s.spec.md file", kind, kind), doc.Location(section.Start, section.Title))
}

type ProductSpec struct {
	Name        string                      `json:"name"`
	ID          string                      `json:"id,omitempty"`
	Description string                      `json:"description,omitempty"`
	Path        string                      `json:"path"`
	Markdown    string                      `json:"markdown"`
	References  *specfile.ReferenceDocument `json:"-"`
}

func ParseProductSpec(kind string, data []byte) (*ProductSpec, error) {
	return ParseProductSpecMarkdown(kind, specfile.NewMarkdown(data))
}

func ParseProductSpecMarkdown(kind string, input *specfile.Markdown) (*ProductSpec, error) {
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
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	spec := &ProductSpec{Name: doc.Name, ID: meta.ID, Markdown: strings.TrimSpace(string(doc.Source)), References: doc.References()}
	seen := map[string]diagnostic.Location{}
	for _, s := range doc.Sections {
		if s.Title == "Overview" && kind != FeatureLayout.Kind {
			continue
		}
		switch s.Title {
		case "Description", "Purpose", "Overview", "Groups", "Enablers", "Supports", "Depends on", "Provided by", "Specification system":
			location := doc.Location(s.Start, s.Title)
			if s.Level != 2 {
				return nil, diagnostic.Error("structure.section-level", s.Title+" requires a level-two heading", location)
			}
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
	slug, ok := strings.CutSuffix(filepath.Base(path), layout.Suffix)
	if !ok {
		return nil, fmt.Errorf("expected a %s source path", layout.Kind)
	}
	if err := layout.ValidateSlug(slug); err != nil {
		return nil, diagnostic.Error("layout."+layout.Kind+"-name", err.Error(), diagnostic.Location{Path: path})
	}
	input, err := files.Markdown(path)
	if err != nil {
		return nil, err
	}
	spec, err := ParseProductSpecMarkdown(layout.Kind, input)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	spec.Path = path
	return spec, nil
}

func ProductCheckType(layout specfile.Layout) check.Type {
	return check.Type{Layout: layout, Context: CheckContext, Inspect: func(files *specfile.SourceSet, path string) ([]diagnostic.Diagnostic, error) {
		_, err := readProduct(files, layout, path)
		return nil, err
	}, References: func(files *specfile.SourceSet, path string) (*specfile.ReferenceDocument, error) {
		spec, err := readProduct(files, layout, path)
		if err != nil {
			return nil, err
		}
		return spec.References, nil
	}}
}
