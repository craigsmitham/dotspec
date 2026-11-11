package requirement

import (
	"fmt"
	"sort"

	"github.com/craigsmitham/dotspec/internal/quality"
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

// Show selects exactly one requirement spec in the requested directory, without searching
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
	spec.System, err = associatedSystem(dir, spec.Path)
	if err != nil {
		return nil, err
	}
	spec.Classification = Classify(spec.QualityAttributes, spec.System)
	return spec, nil
}

// List inventories general and quality requirements, grouping them by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated requirement specs last.
// Requirements within a group are ordered by source path.
func List(path string) ([]*Group, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, func(name string) bool { return Layout.Match(name) || quality.Layout.Match(name) })
	if err != nil {
		return nil, err
	}
	bySystem := make(map[string]*Group)
	var groups []*Group
	for _, path := range paths {
		spec, err := readInventory(path)
		if err != nil {
			return nil, err
		}
		layout := Layout
		if spec.Type == quality.Layout.Kind {
			layout = quality.Layout
		}
		dir, err := layout.Owner(path)
		if err != nil {
			return nil, err
		}
		associated, err := associatedSystem(dir, path)
		if err != nil {
			return nil, err
		}
		spec.System = associated
		spec.Classification = Classify(spec.QualityAttributes, associated)
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

func associatedSystem(dir, source string) (*system.System, error) {
	associated, err := system.Associated(dir, source)
	if err != nil {
		return nil, fmt.Errorf("resolve associated system for requirement specs in %s: %w", dir, err)
	}
	return associated, nil
}

// Quality has specialized authoring fields but is the same kind of obligation.
// The inventory normalizes only its view; source Markdown and type stay intact.
func readInventory(path string) (*Requirement, error) {
	if Layout.Match(path) {
		return Read(path)
	}
	q, err := quality.Read(path)
	if err != nil {
		return nil, err
	}
	return &Requirement{Type: quality.Layout.Kind, Content: q.Content}, nil
}
