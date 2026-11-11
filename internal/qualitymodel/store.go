package qualitymodel

import (
	"embed"
	"fmt"
	"sort"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

const Default = "iso-25010-2023"

//go:embed specs/quality-models/*.quality-model.spec.md
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
		m.Slug = entry.Name()[:len(entry.Name())-len(Suffix)]
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
	if err := ValidateSlug(slug); err != nil {
		return "", err
	}
	if m, err := builtin(slug); err != nil {
		return "", err
	} else if m != nil {
		return "", fmt.Errorf("%q is a reserved built-in name; choose a project-specific slug", slug)
	}
	dir, err := specfile.Scope(path)
	if err != nil {
		return "", err
	}
	if err := Layout.ValidateLocal(dir); err != nil {
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

// Show uses reserved built-in slugs or exactly one model in the selected project directory.
func Show(path, slug string) (*QualityModel, error) {
	if err := ValidateSlug(slug); err != nil {
		return nil, err
	}
	if m, err := builtin(slug); err != nil || m != nil {
		return m, err
	}
	dir, err := specfile.Scope(path)
	if err != nil {
		return nil, err
	}
	if err := Layout.ValidateLocal(dir); err != nil {
		return nil, err
	}
	return Read(Layout.Path(dir, slug))
}

func List(path string) ([]*QualityModel, error) {
	result, err := Builtins()
	if err != nil {
		return nil, err
	}
	dir, err := specfile.Scope(path)
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
