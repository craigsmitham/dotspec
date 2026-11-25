package businessrule

import (
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System        *system.System  `json:"system"`
	BusinessRules []*BusinessRule `json:"business_rules"`
}

func Init(path, slug string) (string, error) {
	if err := Layout.ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	return Layout.Create(dir, slug, Template)
}

var collection = speccollection.Collection[BusinessRule]{
	Layout: Layout, Plural: "business rules", Read: Read,
	SetSystem: func(member *BusinessRule, owner *system.System) {
		member.System = owner
	},
}

func Show(path, selector string) (*BusinessRule, error) {
	return collection.Show(path, selector)
}

func List(path string) ([]*Group, error) {
	return speccollection.List(path, collection, func(owner *system.System, members []*BusinessRule) *Group {
		return &Group{System: owner, BusinessRules: members}
	})
}
