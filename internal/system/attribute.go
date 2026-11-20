package system

import (
	"os"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var AttributeLayout = specfile.Layout{Kind: "quality-attribute", Suffix: specfile.AttributeSuffix, Collection: "quality"}

type AttributeDocument struct {
	AssessmentPlans []specfile.AssessmentPlan `json:"assessment_plans,omitempty"`
	Attribute
	References *specfile.ReferenceDocument
}

func ParseQualityAttribute(path string, data []byte) (*AttributeDocument, error) {
	var meta struct {
		specfile.Metadata `yaml:",inline"`
		ID                string   `yaml:"id"`
		Parents           []string `yaml:"subattribute_of"`
	}
	doc, err := specfile.ParseOutline(data, AttributeLayout.Kind, &meta)
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
		if s.Title == "Specification system" {
			return nil, diagnostic.Error("quality.scope", "the nearest subject declares this quality attribute", doc.Location(s.Start, s.Title))
		}
		if s.Title == "Definition" {
			if seenDefinition || s.Level != 2 {
				return nil, diagnostic.Error("structure.attribute-definition", "use one level-two Definition section", doc.Location(s.Start, s.Title))
			}
			seenDefinition = true
			a.Definition = doc.Body(s)
		}
	}
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	return &AttributeDocument{AssessmentPlans: doc.AssessmentPlans, Attribute: a, References: doc.References()}, nil
}

func ReadQualityAttribute(path string) (*AttributeDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := ParseQualityAttribute(path, data)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	return doc, nil
}

func (s *System) loadQualityTree() error {
	dir := specfile.DeclarationDir(s.Path)
	paths, err := specfile.OwnedPaths(dir, "*"+AttributeLayout.Suffix, AttributeLayout.Match)
	if err != nil {
		return err
	}
	ids := map[string]bool{}
	for _, a := range s.QualityAttributes {
		ids[a.ID] = true
	}
	for _, path := range paths {
		doc, err := ReadQualityAttribute(path)
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
	return validateQualityAttributes(s.QualityAttributes)
}

var attributeContext = &check.ContextProvider{
	Resolve: func(owner string) ([]string, error) {
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
	return check.Type{Layout: AttributeLayout, Context: attributeContext, Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
		_, err := ReadQualityAttribute(path)
		return nil, err
	}, References: func(path string) (*specfile.ReferenceDocument, error) {
		a, err := ReadQualityAttribute(path)
		if err != nil {
			return nil, err
		}
		return a.References, nil
	}}
}
