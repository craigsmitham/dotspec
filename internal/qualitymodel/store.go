package qualitymodel

import (
	"embed"
	"fmt"
	"strings"

	"github.com/dotspec-md/dotspec/internal/specfile"
)

//go:embed specs/quality-models/*.quality-model.md
var starters embed.FS

// Bundled returns the source of a bundled model by slug, so the offline
// handbook can publish it as a reference page.
func Bundled(slug string) ([]byte, error) {
	return starters.ReadFile("specs/quality-models/" + slug + Layout.Suffix)
}

// BundledSlugs lists the bundled models in source order.
func BundledSlugs() ([]string, error) {
	entries, err := starters.ReadDir("specs/quality-models")
	if err != nil {
		return nil, err
	}
	slugs := make([]string, 0, len(entries))
	for _, entry := range entries {
		slugs = append(slugs, strings.TrimSuffix(entry.Name(), Layout.Suffix))
	}
	return slugs, nil
}

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
		// A bundled model has no file; its locations name the model instead.
		m.Location.Path = m.Path
		for i := range m.QualityAttributes {
			m.QualityAttributes[i].Location.Path = m.Path
		}
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
	if err := Layout.ValidateSlug(slug); err != nil {
		return "", err
	}
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
			return "", fmt.Errorf("--from must name a bundled model: %q", from)
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

func UnknownAttribute(m *QualityModel, id string) error {
	where := "read its terms in " + m.Path
	if m.Builtin {
		where = "read its terms with dotspec docs " + Layout.Collection + "/" + m.Slug
	}
	return fmt.Errorf("unknown attribute ID %q in model %s; %s, or start a blank system attribute with dotspec new quality-attribute %s (dotspec docs quality-attribute)", id, m.Slug, where, id)
}
