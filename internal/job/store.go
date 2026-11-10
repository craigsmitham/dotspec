package job

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System *system.System
	Jobs   []*Job
}

func Init(path, slug string) (string, error) {
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	return specfile.Create(path, slug+Suffix, Template)
}

// Show selects exactly one job in the requested directory, without searching
// ancestor directories for a matching slug.
func Show(path, slug string) (*Job, error) {
	if err := ValidateSlug(slug); err != nil {
		return nil, err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, err
	}
	job, err := Read(filepath.Join(dir, slug+Suffix))
	if err != nil {
		return nil, err
	}
	job.System, err = associatedSystem(dir)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// List groups jobs by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated jobs last.
// Jobs within a group are ordered by source path.
func List(path string) ([]*Group, error) {
	paths, err := specfile.Discover(path, func(name string) bool { return strings.HasSuffix(name, Suffix) })
	if err != nil {
		return nil, err
	}
	bySystem := make(map[string]*Group)
	byDirectory := make(map[string]*system.System)
	var groups []*Group
	for _, path := range paths {
		job, err := Read(path)
		if err != nil {
			return nil, err
		}
		dir := filepath.Dir(path)
		associated, found := byDirectory[dir]
		if !found {
			associated, err = associatedSystem(dir)
			if err != nil {
				return nil, err
			}
			byDirectory[dir] = associated
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

func associatedSystem(dir string) (*system.System, error) {
	associated, err := system.Enclosing(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve associated system for jobs in %s: %w", dir, err)
	}
	return associated, nil
}
