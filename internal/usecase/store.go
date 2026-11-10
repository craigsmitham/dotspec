package usecase

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System   *system.System `json:"system"`
	UseCases []*UseCase     `json:"use_cases"`
}

func Init(path, slug string) (string, error) {
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	dir, err := specfile.Scope(path)
	if err != nil {
		return "", err
	}
	if err := specfile.ValidateLocal(dir, Collection, isSpec); err != nil {
		return "", err
	}
	return specfile.Create(dir, filepath.Join(specfile.Folder, Collection, slug+Suffix), Template)
}

func isSpec(name string) bool { return strings.HasSuffix(name, Suffix) }

// Show selects exactly one use case in the requested directory, without searching
// ancestor directories for a matching slug.
func Show(path, slug string) (*UseCase, error) {
	if err := ValidateSlug(slug); err != nil {
		return nil, err
	}
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	if err := specfile.ValidateLocal(dir, Collection, isSpec); err != nil {
		return nil, err
	}
	spec, err := Read(specfile.Path(dir, Collection, slug+Suffix))
	if err != nil {
		return nil, err
	}
	spec.System, err = associatedSystem(dir)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// List groups use cases by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated use cases last.
// Use cases within a group are ordered by source path.
func List(path string) ([]*Group, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, isSpec)
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
		dir, err := specfile.Owner(path, Collection)
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

func associatedSystem(dir string) (*system.System, error) {
	associated, err := system.Enclosing(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve associated system for use cases in %s: %w", dir, err)
	}
	return associated, nil
}
