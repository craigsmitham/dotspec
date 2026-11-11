package job

import (
	"fmt"
	"sort"

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
	dir, err := specfile.Scope(path)
	if err != nil {
		return "", err
	}
	if err := Layout.ValidateLocal(dir); err != nil {
		return "", err
	}
	return Layout.Create(dir, slug, template)
}

// Show selects exactly one job in the requested directory, without searching
// ancestor directories for a matching slug.
func Show(path, slug string) (*Job, error) {
	if err := ValidateSlug(slug); err != nil {
		return nil, err
	}
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	if err := Layout.ValidateLocal(dir); err != nil {
		return nil, err
	}
	job, err := Read(Layout.Path(dir, slug))
	if err != nil {
		return nil, err
	}
	job.System, err = associatedSystem(dir, job.Path)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// List groups jobs by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated jobs last.
// Jobs within a group are ordered by source path.
func List(path string) ([]*Group, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, Layout.Match)
	if err != nil {
		return nil, err
	}
	bySystem := make(map[string]*Group)
	var groups []*Group
	for _, path := range paths {
		job, err := Read(path)
		if err != nil {
			return nil, err
		}
		dir, err := Layout.Owner(path)
		if err != nil {
			return nil, err
		}
		associated, err := associatedSystem(dir, path)
		if err != nil {
			return nil, err
		}
		job.System = associated
		key := ""
		if associated != nil {
			key = associated.Path
		}
		group, found := bySystem[key]
		if !found {
			group = &Group{System: associated}
			bySystem[key] = group
			groups = append(groups, group)
		}
		group.Jobs = append(group.Jobs, job)
	}
	sort.Slice(groups, func(i, j int) bool {
		a, b := groups[i].System, groups[j].System
		if a == nil {
			return false
		}
		if b == nil {
			return true
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.Path < b.Path
	})
	return groups, nil
}

func associatedSystem(dir, source string) (*system.System, error) {
	associated, err := system.Associated(dir, source)
	if err != nil {
		return nil, fmt.Errorf("resolve associated system for jobs in %s: %w", dir, err)
	}
	return associated, nil
}
