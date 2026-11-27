package system

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"go.yaml.in/yaml/v3"
)

func (s *System) glossaryReferences(input *specfile.Markdown) error {
	if len(s.GlossarySources) == 0 {
		return nil
	}
	location := input.MetadataLocation(specfile.GlossarySourcesField)
	seen := map[string]bool{}
	for _, name := range s.GlossarySources {
		ext := strings.ToLower(filepath.Ext(name))
		if filepath.IsAbs(name) || strings.ContainsAny(name, "\\:") || (ext != ".md" && ext != ".markdown") || strings.HasSuffix(name, ".spec.md") || strings.HasSuffix(name, specfile.QualityModelSuffix) {
			return diagnostic.Error("glossary.source", fmt.Sprintf("expected a relative Markdown glossary path, found %q", name), location)
		}
		key := filepath.Clean(name)
		if seen[key] {
			return diagnostic.Error("glossary.source", "duplicate glossary source "+name, location)
		}
		seen[key] = true
		s.References.Links = append(s.References.Links, specfile.Reference{Destination: (&url.URL{Path: name}).String(), Location: location})
	}
	return nil
}

func (s *System) checkGlossarySources(files *specfile.SourceSet) error {
	if len(s.GlossarySources) == 0 {
		return nil
	}
	dir := filepath.Dir(s.Path)
	root, err := specfile.RepositoryRoot(dir)
	if err != nil {
		return err
	}
	for _, name := range s.GlossarySources {
		path := filepath.Join(dir, filepath.FromSlash(name))
		loc := diagnostic.Location{Path: s.Path, Element: specfile.GlossarySourcesField}
		info, allowed, err := check.ReferenceTarget(path, root)
		if !allowed || os.IsNotExist(err) || (err == nil && !info.Mode().IsRegular()) {
			return diagnostic.Error("glossary.source", "glossary source must be an accessible Markdown file within the repository: "+name, loc)
		}
		if err != nil {
			return err
		}
		data, err := files.Markdown(path)
		if err != nil {
			return err
		}
		if _, err := ParseGlossary(data); err != nil {
			return diagnostic.WithPath(err, path)
		}
	}
	return nil
}

// ParseGlossary validates the vocabulary and reads its optional context name.
func ParseGlossary(input *specfile.Markdown) (string, error) {
	doc, err := input.Outline("glossary")
	if err != nil {
		return "", err
	}
	for _, term := range doc.Sections {
		if term.Level != 2 || term.Title == "" || doc.Body(term) == "" {
			return "", diagnostic.Error("glossary.term", "use a named level-two heading and a definition for each glossary entry", doc.Location(term.Start, term.Title))
		}
	}
	if err := doc.ValidateReferences(); err != nil {
		return "", err
	}
	if input.MetadataError != nil {
		return "", input.MetadataError
	}
	location := input.MetadataLocation("context")
	node := input.Value("context")
	if node == nil {
		return "", nil
	}
	if node.Kind != yaml.ScalarNode || node.ShortTag() != "!!str" || strings.TrimSpace(node.Value) == "" || strings.ContainsAny(strings.TrimSpace(node.Value), "\r\n") {
		return "", diagnostic.Error("glossary.context", "expected a nonempty, single-line context name", location)
	}
	return node.Value, nil
}
