// Package speccollection reads and checks specification families whose members are presented
// with their associated system. The family supplies parsing and enrichment.
package speccollection

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Collection[T any] struct {
	Layout    specfile.Layout
	Plural    string
	Read      func(path string) (*T, error)
	SetSystem func(*T, *system.System)
}

func (c Collection[T]) readMember(path string) (*T, *system.System, error) {
	member, err := c.Read(path)
	if err != nil {
		return nil, nil, err
	}
	dir := c.Layout.ContextDir(path)
	owner, err := system.Associated(dir, path)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve associated system for %s in %s: %w", c.Plural, dir, err)
	}
	c.SetSystem(member, owner)
	return member, owner, nil
}

// List groups discovered members by system source identity, preserving source
// order within each group. The caller constructs its family-specific result.
func List[T, G any](path string, c Collection[T], makeGroup func(*system.System, []*T) G) ([]G, error) {
	paths, err := specfile.Discover(path, c.Layout.Match)
	if err != nil {
		return nil, err
	}
	type bucket struct {
		owner   *system.System
		members []*T
	}
	byPath := make(map[string]*bucket)
	for _, file := range paths {
		member, owner, err := c.readMember(file)
		if err != nil {
			return nil, err
		}
		key := ""
		if owner != nil {
			key = owner.Path
		}
		group := byPath[key]
		if group == nil {
			group = &bucket{owner: owner}
			byPath[key] = group
		}
		group.members = append(group.members, member)
	}
	ordered := slices.SortedFunc(maps.Values(byPath), func(left, right *bucket) int {
		a, b := left.owner, right.owner
		switch {
		case a == nil && b == nil:
			return 0
		case a == nil:
			return 1
		case b == nil:
			return -1
		default:
			return cmp.Or(cmp.Compare(a.Name, b.Name), cmp.Compare(a.Path, b.Path))
		}
	})
	var groups []G
	for _, group := range ordered {
		groups = append(groups, makeGroup(group.owner, group.members))
	}
	return groups, nil
}

// Member is a discovered specification of any layout with its associated system.
type Member struct {
	Name   string
	Path   string
	System *system.System
}

// Members discovers the files of a layout below path and associates each with
// its system. The order follows List. The name function reads one file and
// returns its heading.
func Members(path string, layout specfile.Layout, plural string, name func(path string) (string, error)) ([]*Member, error) {
	c := Collection[Member]{Layout: layout, Plural: plural,
		Read: func(file string) (*Member, error) {
			heading, err := name(file)
			if err != nil {
				return nil, err
			}
			return &Member{Name: heading, Path: file}, nil
		},
		SetSystem: func(m *Member, owner *system.System) { m.System = owner },
	}
	groups, err := List(path, c, func(_ *system.System, members []*Member) []*Member { return members })
	if err != nil {
		return nil, err
	}
	members := []*Member{}
	for _, group := range groups {
		members = append(members, group...)
	}
	return members, nil
}
