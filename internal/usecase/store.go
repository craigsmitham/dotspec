package usecase

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System   *system.System `json:"system"`
	UseCases []*UseCase     `json:"use_cases"`
}

func Init(path, slug string) (string, error) {
	return InitTemplate(path, slug, "brief")
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
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	return Layout.Create(dir, slug, content)
}

// Show resolves a source path or unique slug within the selected owner.
func Show(path, slug string) (*UseCase, error) {
	if !Layout.Match(slug) {
		if err := ValidateSlug(slug); err != nil {
			return nil, err
		}
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, err
	}
	file, err := Layout.LocalPath(dir, slug)
	if err != nil {
		return nil, err
	}
	spec, err := Read(file)
	if err != nil {
		return nil, err
	}
	spec.System, err = associatedSystem(filepath.Dir(spec.Path), spec.Path)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// List groups use cases by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated use cases last.
// Use cases within a group are ordered by source path.
func List(path string) ([]*Group, error) {
	dir, err := specfile.Directory(path)
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
		spec, err := Read(path)
		if err != nil {
			return nil, err
		}
		dir := Layout.ContextDir(path)
		associated, err := associatedSystem(dir, path)
		if err != nil {
			return nil, err
		}
		spec.System = associated
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
		group.UseCases = append(group.UseCases, spec)
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
		return nil, fmt.Errorf("resolve associated system for use cases in %s: %w", dir, err)
	}
	return associated, nil
}
