package concept

import (
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System   *system.System `json:"system"`
	Concepts []*Concept     `json:"concepts"`
}

func Init(path, slug string) (string, error) {
	return Layout.Create(path, slug, Template)
}

var collection = speccollection.Collection[Concept]{
	Layout: Layout, Plural: "concepts", Read: Read,
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
