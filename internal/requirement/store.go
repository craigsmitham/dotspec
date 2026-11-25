package requirement

import (
	"fmt"
	"strings"

	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System       *system.System `json:"system"`
	Requirements []*Requirement `json:"requirements"`
}

func Init(path, slug string) (string, error) { return InitWithPreset(path, slug, "general") }

// InitWithPreset creates a requirement from the selected authoring template.
func InitWithPreset(path, slug, preset string) (string, error) {
	return InitWithAttribute(path, slug, preset, "")
}

func InitWithAttribute(path, slug, preset, attribute string) (string, error) {
	template, err := TemplateFor(preset)
	if err != nil {
		return "", err
	}
	if err := Layout.ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	if attribute == "" {
		return Layout.Create(dir, slug, template)
	}
	chain, err := system.SubjectChain(dir)
	if err != nil {
		return "", err
	}
	if len(chain) == 0 {
		return "", fmt.Errorf("no governing system for quality attribute %q in %s", attribute, dir)
	}
	sys, err := system.Read(chain[len(chain)-1])
	if err != nil {
		return "", err
	}
	if sys.FindQualityAttribute(attribute) == nil {
		return "", fmt.Errorf("unknown quality attribute %q in %s", attribute, sys.Path)
	}
	template = strings.Replace(template, "type: requirement-spec\n", fmt.Sprintf("type: requirement-spec\nquality_attributes: [%q]\n", attribute), 1)
	return Layout.Create(dir, slug, template)
}

var collection = speccollection.Collection[Requirement]{
	Layout: Layout, Plural: "requirement specs", Read: Read,
	SetSystem: func(member *Requirement, owner *system.System) {
		member.System = owner
		member.Classification = Classify(member.QualityAttributes, owner)
	},
}

func Show(path, selector string) (*Requirement, error) {
	return collection.Show(path, selector)
}

func List(path string) ([]*Group, error) {
	return speccollection.List(path, collection, func(owner *system.System, members []*Requirement) *Group {
		return &Group{System: owner, Requirements: members}
	})
}
