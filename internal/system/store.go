package system

import (
	"fmt"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

type Node struct {
	System   *System `json:"system"`
	Children []*Node `json:"children"`
}

// Show resolves the nearest system document and its enclosing directory context.
func Show(path string) (spec, parent *System, err error) {
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, nil, err
	}
	root, err := specfile.ContextRoot(dir)
	if err != nil {
		return nil, nil, err
	}
	for current := dir; ; current = filepath.Dir(current) {
		file, err := Layout.LocalSubject(current)
		if err != nil {
			return nil, nil, err
		}
		if file != "" {
			candidate, err := Read(file)
			if err != nil {
				return nil, nil, err
			}
			if spec != nil {
				return spec, candidate, nil
			}
			spec = candidate
		}
		if current == root || filepath.Dir(current) == current {
			break
		}
	}
	if spec == nil {
		return nil, nil, fmt.Errorf("no system spec found for %s; start a project with dotspec init <system-name>", path)
	}
	return spec, nil, nil
}

// List returns directory nesting in source-path order, not declared part_of structure.
// Enclosing documents outside the selected subtree are absent from the forest.
func List(path string) ([]*Node, error) {
	dir, err := specfile.Directory(path)
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
		owner := Layout.ContextDir(path)
		byDir[owner] = node
	}
	var roots []*Node
	for _, node := range nodes {
		var parent *Node
		owner := Layout.ContextDir(node.System.Path)
		for current := owner; current != dir && filepath.Dir(current) != current; {
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
