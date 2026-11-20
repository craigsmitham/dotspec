package requirement

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
)

type Group struct {
	System       *system.System `json:"system"`
	Requirements []*Requirement `json:"requirements"`
}

func Init(path, slug string) (string, error) { return InitWithPreset(path, slug, "general") }

// InitWithPreset creates a requirement from the selected authoring template.
func InitWithPreset(path, slug, preset string) (string, error) {
	return InitWithAttribute(path, slug, preset, "")
}

func InitWithAttribute(path, slug, preset, attribute string) (string, error) {
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
	if attribute == "" {
		return Layout.Create(dir, slug, template)
	}
	chain, err := system.SubjectChain(dir)
	if err != nil {
		return "", err
	}
	if len(chain) == 0 {
		return "", fmt.Errorf("no governing system for quality attribute %q in %s", attribute, dir)
	}
	sys, err := system.Read(chain[len(chain)-1])
	if err != nil {
		return "", err
	}
	if sys.FindQualityAttribute(attribute) == nil {
		return "", fmt.Errorf("unknown quality attribute %q in %s", attribute, sys.Path)
	}
	template = strings.Replace(template, "type: requirement-spec\n", fmt.Sprintf("type: requirement-spec\nquality_attributes: [%q]\n", attribute), 1)
	return Layout.Create(dir, slug, template)
}

// Show resolves a source path or unique slug within the selected owner.
func Show(path, slug string) (*Requirement, error) {
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
	spec.Classification = Classify(spec.QualityAttributes, spec.System)
	return spec, nil
}

// List inventories requirements, grouping them by the identity of their enclosing system (its source path).
// Groups are ordered by system name then path, with unassociated requirement specs last.
// Requirements within a group are ordered by source path.
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
