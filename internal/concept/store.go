package concept

import (
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System   *system.System `json:"system"`
	Concepts []*Concept     `json:"concepts"`
}

func Init(path, slug string) (string, error) {
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	return Layout.Create(dir, slug, Template)
}

var collection = speccollection.Collection[Concept]{
	Layout: Layout, Plural: "concepts", ValidateSlug: ValidateSlug, Read: Read,
	SetSystem: func(member *Concept, owner *system.System) {
		member.System = owner
	},
}

func Show(path, selector string) (*Concept, error) {
	return collection.Show(path, selector)
}

func List(path string) ([]*Group, error) {
	return speccollection.List(path, collection, func(owner *system.System, members []*Concept) *Group {
		return &Group{System: owner, Concepts: members}
	})
}
