package system

import (
	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var AttributeLayout = specfile.Layout{Kind: "quality-attribute", Suffix: specfile.AttributeSuffix}

type AttributeDocument struct {
	Attribute
	References *specfile.ReferenceDocument
}

func ParseQualityAttribute(path string, data []byte) (*AttributeDocument, error) {
	return ParseQualityAttributeMarkdown(path, specfile.NewMarkdown(data))
}

func ParseQualityAttributeMarkdown(path string, input *specfile.Markdown) (*AttributeDocument, error) {
	var meta struct {
		specfile.Metadata `yaml:",inline"`
		ID                string   `yaml:"id"`
		Parents           []string `yaml:"subattribute_of"`
	}
	doc, err := input.TypedOutline(AttributeLayout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	if err := AttributeLayout.CheckName(path, doc.Name); err != nil {
		return nil, err
	}
	if meta.ID == "" {
		meta.ID = specfile.NameSlug(doc.Name)
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "quality attribute", doc.NameLocation); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	a := Attribute{ID: meta.ID, Name: doc.Name, Location: doc.NameLocation, SubattributeOf: meta.Parents}
	a.Location.Path = path
	a.subattributeLocation = a.Location
	var parents specfile.Identifiers
	for _, id := range meta.Parents {
		if err := parents.Register(id, "parent attribute", a.Location); err != nil {
			return nil, err
		}
	}
	seenDefinition := false
	for _, s := range doc.Sections {
		if s.Title == "Definition" {
			if seenDefinition {
				return nil, diagnostic.Error("structure.attribute-definition", "use one Definition section", doc.Location(s.Start, s.Title))
			}
			seenDefinition = true
			a.Definition = doc.Body(s)
		}
	}
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	return &AttributeDocument{Attribute: a, References: doc.References()}, nil
}

func readAttribute(files *specfile.SourceSet, path string) (*AttributeDocument, error) {
	input, err := files.Markdown(path)
	if err != nil {
		return nil, err
	}
	doc, err := ParseQualityAttributeMarkdown(path, input)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	return doc, nil
}

func (s *System) loadQualityTree(files *specfile.SourceSet) error {
	dir := specfile.DeclarationDir(s.Path)
	paths, err := specfile.OwnedPaths(dir, "*"+AttributeLayout.Suffix, AttributeLayout.Match)
	if err != nil {
		return err
	}
	ids := map[string]bool{}
	for _, path := range paths {
		doc, err := readAttribute(files, path)
		if err != nil {
			return err
		}
		a := doc.Attribute
		if ids[a.ID] {
			return diagnostic.Error("quality.identity", "duplicate quality attribute ID "+a.ID, a.Location)
		}
		ids[a.ID] = true
		s.QualityAttributes = append(s.QualityAttributes, a)
	}
	return ValidateQualityAttributes(s.QualityAttributes)
}

var attributeContext = &check.ContextProvider{
	Resolve: func(_ *specfile.SourceSet, owner string) ([]string, error) {
		boundary, err := specfile.ContextRoot(owner)
		if err != nil {
			return nil, err
		}
		file, err := specfile.NearestDeclaration(owner, boundary)
		if err != nil {
			return nil, err
		}
		if !Layout.Match(file) {
			return nil, diagnostic.Error("quality.scope", "quality attributes require an owning system", diagnostic.Location{Path: owner})
		}
		return []string{file}, nil
	},
	Inspect: inspect,
}

func AttributeCheckType() check.Type {
	kind := check.ReadType(AttributeLayout, readAttribute, func(a *AttributeDocument) *specfile.ReferenceDocument { return a.References }, nil)
	kind.Context = attributeContext
	return kind
}
