package specfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const AttributeSuffix = ".quality-attribute.spec.md"

func AttributePath(dir string) string { return filepath.Join(dir, filepath.Base(dir)+AttributeSuffix) }

// OwnedPaths discovers matching sources, excluding the domains of child subjects.
func OwnedPaths(dir, pattern string, match func(string) bool) ([]string, error) {
	dir, err := Directory(dir)
	if err != nil {
		return nil, err
	}
	declaration, err := LocalDeclaration(dir)
	if err != nil {
		return nil, err
	}
	paths, err := discover(dir, ":(glob)**/"+pattern, match)
	if err != nil {
		return nil, err
	}
	owned := paths[:0]
	for _, p := range paths {
		subject, err := NearestDeclaration(filepath.Dir(p), dir)
		if err != nil {
			var conflict *diagnostic.Diagnostic
			if errors.As(err, &conflict) && conflict.Rule == "layout.subject-ambiguous" {
				continue
			} // The conflicting child scope owns its own files.
			return nil, err
		}
		if subject == declaration {
			owned = append(owned, p)
		}
	}
	return owned, nil
}

// SourcePath resolves a typed source path relative to the selected directory.
func (l Layout) SourcePath(dir, selector string) (string, error) {
	if !l.Match(selector) {
		return "", fmt.Errorf("expected a %s source path", l.Kind)
	}
	path := selector
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}
	path = filepath.Clean(path)
	boundary, err := ContextRoot(dir)
	if err != nil {
		return "", err
	}
	if boundary == "" {
		boundary = dir
	}
	rel, err := filepath.Rel(boundary, path)
	if err != nil || !filepath.IsLocal(rel) {
		return "", fmt.Errorf("source must be inside %s", boundary)
	}
	ok, err := regularPath(path, boundary, hasManifest)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("unavailable %s source: %s", l.Kind, path)
	}
	return path, nil
}

// LocalPath resolves a source path, or a slug owned by the nearest subject.
// A collection slug must follow the layout's slug grammar.
func (l Layout) LocalPath(dir, selector string) (string, error) {
	if l.Match(selector) {
		return l.SourcePath(dir, selector)
	}
	if err := l.validateName(selector); err != nil {
		return "", err
	}
	boundary, err := ContextRoot(dir)
	if err != nil {
		return "", err
	}
	declaration, err := NearestDeclaration(dir, boundary)
	if err != nil {
		return "", err
	}
	if declaration != "" {
		dir = DeclarationDir(declaration)
	}
	name := l.filename(selector)
	paths, err := OwnedPaths(dir, name, func(file string) bool { return file == name })
	if err != nil {
		return "", err
	}
	switch len(paths) {
	case 0:
		return "", &os.PathError{Op: "find " + l.Kind, Path: l.Path(dir, selector), Err: os.ErrNotExist}
	case 1:
		return paths[0], nil
	default:
		return "", fmt.Errorf("ambiguous %s %q; select a source path: %s", l.Kind, selector, strings.Join(paths, ", "))
	}
}
