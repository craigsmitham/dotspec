package job

import (
	"github.com/craigsmitham/dotspec/internal/speccollection"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System *system.System `json:"system"`
	Jobs   []*Job         `json:"jobs"`
}

func Init(path, slug string) (string, error) { return InitWithPreset(path, slug, "general") }

func InitWithPreset(path, slug, preset string) (string, error) {
	template, err := TemplateFor(preset)
	if err != nil {
		return "", err
	}
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	return Layout.Create(dir, slug, template)
}

var collection = speccollection.Collection[Job]{
	Layout: Layout, Plural: "jobs", ValidateSlug: ValidateSlug, Read: Read,
	SetSystem: func(member *Job, owner *system.System) {
		member.System = owner
	},
}

func Show(path, selector string) (*Job, error) {
	return collection.Show(path, selector)
}

func List(path string) ([]*Group, error) {
	return speccollection.List(path, collection, func(owner *system.System, members []*Job) *Group {
		return &Group{System: owner, Jobs: members}
	})
}
