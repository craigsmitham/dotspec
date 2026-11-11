package system

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

type Node struct {
	System   *System `json:"system"`
	Children []*Node `json:"children"`
}

func Init(path string) (string, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return "", err
	}
	if err := Layout.ValidateLocal(dir); err != nil {
		return "", err
	}
	return Layout.Create(dir, "", Template)
}

// Enclosing returns the nearest system, or nil when none exists. Malformed
// enclosing specs remain errors. Lookup stops at the repository boundary.
func Enclosing(path string) (*System, error) {
	specs, err := enclosing(path, 1)
	if err != nil || len(specs) == 0 {
		return nil, err
	}
	return specs[0], nil
}

// Show resolves the nearest system document and its enclosing directory context.
func Show(path string) (spec, parent *System, err error) {
	specs, err := enclosing(path, 2)
	if err != nil {
		return nil, nil, err
	}
	if len(specs) == 0 {
		return nil, nil, fmt.Errorf("no system spec found for %s; use dotspec system init to create one", path)
	}
	if len(specs) > 1 {
		parent = specs[1]
	}
	return specs[0], parent, nil
}

func enclosing(path string, limit int) ([]*System, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	root, err := specfile.ContextRoot(dir)
	if err != nil {
		return nil, err
	}
	var specs []*System
	for current := dir; ; current = filepath.Dir(current) {
		if err := Layout.ValidateLocal(current); err != nil {
			return nil, err
		}
		candidate, err := Read(Layout.Path(current, ""))
		if err == nil {
			specs = append(specs, candidate)
			if len(specs) == limit {
				return specs, nil
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if current == root || filepath.Dir(current) == current {
			break
		}
	}
	return specs, nil
}

// List returns directory nesting in source-path order, not declared part_of structure.
// Enclosing documents outside the selected subtree are absent from the forest.
func List(path string) ([]*Node, error) {
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, Layout.Match)
	if err != nil {
		return nil, err
	}
	byDir := make(map[string]*Node)
	var nodes []*Node
	for _, path := range paths {
		spec, err := Read(path)
		if err != nil {
			return nil, err
		}
		node := &Node{System: spec, Children: []*Node{}}
		nodes = append(nodes, node)
		owner, err := Layout.Owner(path)
		if err != nil {
			return nil, err
		}
		byDir[owner] = node
	}
	var roots []*Node
	for _, node := range nodes {
		var parent *Node
		owner, _ := Layout.Owner(node.System.Path) // Read already validated placement.
		for current := owner; current != dir; {
			current = filepath.Dir(current)
			if parent = byDir[current]; parent != nil {
				break
			}
		}
		if parent == nil {
			roots = append(roots, node)
		} else {
			parent.Children = append(parent.Children, node)
		}
	}
	return roots, nil
}
