package system

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var AttributeLayout = specfile.Layout{Kind: "quality-attribute", Suffix: specfile.AttributeSuffix, Collection: "quality", QualityTree: true}

type AttributeDocument struct {
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
	if filepath.Base(path) != filepath.Base(filepath.Dir(path))+AttributeLayout.Suffix {
		return nil, fmt.Errorf("%s: attribute directory must match the filename", path)
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
			return nil, diagnostic.Error("quality.scope", "the quality folder determines the declaring system", doc.Location(s.Start, s.Title))
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
	return &AttributeDocument{Attribute: a, References: doc.References()}, nil
}

func ReadQualityAttribute(path string) (*AttributeDocument, error) {
	if _, err := AttributeLayout.Owner(path); err != nil {
		return nil, err
	}
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
	root := filepath.Join(filepath.Dir(s.Path), "quality")
	info, err := os.Lstat(root)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s must be a directory", root)
	}
	byDir := map[string]string{}
	ids := map[string]bool{}
	for _, a := range s.QualityAttributes {
		ids[a.ID] = true
	}
	err = filepath.WalkDir(root, func(dir string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() || dir == root {
			return nil
		}
		path := specfile.AttributePath(dir)
		info, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("quality folder requires its named declaration: %w", err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("expected a regular attribute source: %s", path)
		}
		doc, err := ReadQualityAttribute(path)
		if err != nil {
			return err
		}
		a := doc.Attribute
		if ids[a.ID] {
			return diagnostic.Error("quality.identity", "duplicate quality attribute ID "+a.ID, a.Location)
		}
		ids[a.ID] = true
		byDir[dir] = a.ID
		if parent := byDir[filepath.Dir(dir)]; parent != "" {
			additional := slices.DeleteFunc(a.SubattributeOf, func(id string) bool { return id == parent })
			a.SubattributeOf = append([]string{parent}, additional...)
		}
		s.QualityAttributes = append(s.QualityAttributes, a)
		return nil
	})
	if err != nil {
		return err
	}
	return validateQualityAttributes(s.QualityAttributes)
}

var attributeContext = &check.ContextProvider{
	Resolve: func(owner string) ([]string, error) {
		file, err := Layout.LocalSubject(owner)
		if err != nil {
			return nil, err
		}
		if file == "" {
			return nil, diagnostic.Error("quality.scope", "quality attributes require a system in the owning specs directory", diagnostic.Location{Path: owner})
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
