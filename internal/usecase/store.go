package usecase

import (
	"fmt"

	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System   *system.System `json:"system"`
	UseCases []*UseCase     `json:"use_cases"`
}

// InitTemplate selects writing prompts, not a separate format or maturity state.
func InitTemplate(path, slug, template string) (string, error) {
	content := Template
	switch template {
	case "brief":
	case "expanded":
		content = ExpandedTemplate
	default:
		return "", fmt.Errorf("unknown use-case template %q; use brief or expanded", template)
	}
	return Layout.Create(path, slug, content)
}

var collection = speccollection.Collection[UseCase]{
	Layout: Layout, Plural: "use cases", Read: Read,
	SetSystem: func(member *UseCase, owner *system.System) {
		member.System = owner
	},
}

func Show(path, selector string) (*UseCase, error) {
	return collection.Show(path, selector)
}

func List(path string) ([]*Group, error) {
	return speccollection.List(path, collection, func(owner *system.System, members []*UseCase) *Group {
		return &Group{System: owner, UseCases: members}
	})
}
