package quality

import (
	"fmt"
	"sort"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System       *system.System `json:"system"`
	Requirements []*Requirement `json:"requirements"`
}

func Init(path, slug string) (string, error) {
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
	return Layout.Create(dir, slug, Template)
}

// Show selects exactly one quality spec in the requested directory, without searching
// ancestor directories for a matching slug.
func Show(path, slug string) (*Requirement, error) {
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
	spec, err := Read(Layout.Path(dir, slug))
	if err != nil {
		return nil, err
	}
	spec.System, err = associatedSystem(dir)
	if err != nil {
		return nil, err
	}
	spec.Classification = Classify(spec.Attributes, spec.System)
	return spec, nil
}

// List groups quality specs by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated quality specs last.
// Requirements within a group are ordered by source path.
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
	byDirectory := make(map[string]*system.System)
	var groups []*Group
	for _, path := range paths {
		spec, err := Read(path)
		if err != nil {
			return nil, err
		}
		dir, err := Layout.Owner(path)
		if err != nil {
			return nil, err
		}
		associated, found := byDirectory[dir]
		if !found {
			associated, err = associatedSystem(dir)
			if err != nil {
				return nil, err
			}
			byDirectory[dir] = associated
		}
		spec.System = associated
		spec.Classification = Classify(spec.Attributes, associated)
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
		group.Requirements = append(group.Requirements, spec)
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

func associatedSystem(dir string) (*system.System, error) {
	associated, err := system.Enclosing(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve associated system for quality specs in %s: %w", dir, err)
	}
	return associated, nil
}
