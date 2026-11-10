// Package quality is an illustrative test integration, not a supported public
// specification format. It deliberately depends only on the shared machinery.
package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var Layout = specfile.Layout{Kind: "quality-requirement", Collection: "quality-requirements", Suffix: ".quality-requirement.spec.md"}

type metadata struct {
	specfile.Metadata `yaml:",inline"`
	ID                string `yaml:"id"`
	Unit              string `yaml:"unit"`
}

type Requirement struct {
	ID, Unit, Markdown string
	Location           diagnostic.Location
}

func Parse(data []byte) (*Requirement, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, Layout.Kind, &meta)
	if err != nil {
		return nil, err
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "requirement", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	return &Requirement{ID: meta.ID, Unit: meta.Unit, Markdown: string(doc.Source), Location: doc.NameLocation}, nil
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func validateSlug(slug string) error {
	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("invalid quality-requirement slug %q", slug)
	}
	return nil
}

func CheckType() check.Type {
	return check.Type{Layout: Layout, ValidateSlug: validateSlug, Inspect: func(path string) ([]diagnostic.Diagnostic, error) {
		if err := validateSlug(strings.TrimSuffix(filepath.Base(path), Layout.Suffix)); err != nil {
			return nil, diagnostic.Error("layout.quality-requirement-name", err.Error(), diagnostic.Location{Path: path})
		}
		if _, err := Layout.Owner(path); err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		requirement, err := Parse(data)
		if err != nil {
			return nil, diagnostic.WithPath(err, path)
		}
		if requirement.Unit == "" {
			location := requirement.Location
			location.Path = path
			return []diagnostic.Diagnostic{diagnostic.Advice("quality-requirement.unit", "Identify the unit of measurement.", location)}, nil
		}
		return nil, nil
	}}
}
