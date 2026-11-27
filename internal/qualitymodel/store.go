package qualitymodel

import (
	"embed"
	"fmt"
	"sort"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Default = "iso-25010-2023"

//go:embed specs/quality-models/*.quality-model.md
var starters embed.FS

func Builtins() ([]*QualityModel, error) {
	entries, err := starters.ReadDir("specs/quality-models")
	if err != nil {
		return nil, err
	}
	result := []*QualityModel{}
	for _, entry := range entries {
		data, err := starters.ReadFile("specs/quality-models/" + entry.Name())
		if err != nil {
			return nil, err
		}
		m, err := Parse(data)
		if err != nil {
			return nil, fmt.Errorf("built-in %s: %w", entry.Name(), err)
		}
		m.Slug = entry.Name()[:len(entry.Name())-len(Layout.Suffix)]
		m.Path = "builtin:" + m.Slug
		m.Builtin = true
		result = append(result, m)
	}
	return result, nil
}

func builtin(slug string) (*QualityModel, error) {
	models, err := Builtins()
	if err != nil {
		return nil, err
	}
	for _, m := range models {
		if m.Slug == slug {
			return m, nil
		}
	}
	return nil, nil
}

func Init(path, slug, from string) (string, error) {
	if m, err := builtin(slug); err != nil {
		return "", err
	} else if m != nil {
		return "", fmt.Errorf("%q is a reserved built-in name; choose a project-specific slug", slug)
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	content := Template
	if from != "" {
		m, err := builtin(from)
		if err != nil {
			return "", err
		}
		if m == nil {
			return "", fmt.Errorf("--from must name a bundled starter: %q", from)
		}
		content = "---\ntype: quality-model\nid: " + slug + "\n---\n\n" + m.Markdown + "\n"
	}
	return Layout.Create(dir, slug, content)
}

// Show resolves a source path or unique slug within the selected owner.
func Show(path, slug string) (*QualityModel, error) {
	if m, err := builtin(slug); err != nil || m != nil {
		return m, err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, err
	}
	file, err := Layout.LocalPath(dir, slug)
	if err != nil {
		return nil, err
	}
	return Read(file)
}

func List(path string) ([]*QualityModel, error) {
	result, err := Builtins()
	if err != nil {
		return nil, err
	}
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, Layout.Match)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		m, err := Read(path)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Builtin != result[j].Builtin {
			return result[i].Builtin
		}
		return result[i].Path < result[j].Path
	})
	return result, nil
}

// Selection describes the vocabulary used by this invocation.
type Selection struct {
	Basis  string `json:"basis"`
	ID     string `json:"id"`
	Source string `json:"source"`
}

func Resolve(dir string, project *specfile.Project, explicit string) (*QualityModel, Selection, error) {
	basis, selector := "builtin_default", Default
	var m *QualityModel
	var err error
	switch {
	case explicit != "":
		basis, selector = "explicit", explicit
		m, err = Show(dir, selector)
	case project != nil && project.QualityModel != "":
		basis, selector = "project", project.QualityModel
		if Layout.Match(selector) {
			m, err = Show(project.Root, selector)
		} else {
			m, err = builtin(selector)
			if m == nil && err == nil {
				err = diagnostic.Error("project.quality-model", fmt.Sprintf("unknown bundled quality model %q", selector), project.QualityLocation)
			}
		}
		if err != nil {
			err = fmt.Errorf("%s: quality.model %q: %w", project.Manifest(), selector, err)
		}
	default:
		m, err = builtin(selector)
	}
	if err != nil {
		return nil, Selection{}, err
	}
	id := m.ID
	if id == "" {
		id = m.Slug
	}
	return m, Selection{Basis: basis, ID: id, Source: m.Path}, nil
}

func UnknownAttribute(m *QualityModel, id string) error {
	return fmt.Errorf("unknown attribute ID %q in model %s; use dotspec quality-attribute browse to find a term, or author a system *.quality-attribute.spec.md file (dotspec docs quality-attribute)", id, m.Slug)
}
