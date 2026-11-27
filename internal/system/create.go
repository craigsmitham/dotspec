package system

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/rogpeppe/go-internal/lockedfile"
)

func InitNamed(path, name string) (string, error) {
	return createSubject(path, name, Layout, func(_ string) (string, string, error) {
		source := strings.Replace(Template, "# New system", "# "+name, 1)
		spec, err := Parse([]byte(source))
		if err != nil {
			return "", "", err
		}
		return source, spec.Name, nil
	})
}

func InitComponent(path, name, parent string) (string, error) {
	return createSubject(path, name, ComponentLayout, func(target string) (string, string, error) {
		// Resolve the working directory the same way as --path. A shell can
		// report a symlinked spelling for its current repository.
		if !filepath.IsAbs(parent) {
			cwd, err := specfile.Directory(".")
			if err != nil {
				return "", "", err
			}
			parent = filepath.Join(cwd, parent)
		}
		owner := specfile.DeclarationDir(target)
		boundary, err := specfile.ContextRoot(owner)
		if err != nil {
			return "", "", err
		}
		info, ok, err := check.ReferenceTarget(parent, boundary)
		if err != nil {
			return "", "", err
		}
		if !ok || !info.Mode().IsRegular() {
			return "", "", fmt.Errorf("unavailable parent declaration: %s", parent)
		}
		parentBoundary, err := specfile.ContextRoot(specfile.DeclarationDir(parent))
		if err != nil {
			return "", "", err
		}
		if boundary != parentBoundary {
			return "", "", fmt.Errorf("parent must belong to the same project: %s", parent)
		}
		switch {
		case Layout.Match(parent):
			if _, err := Read(parent); err != nil {
				return "", "", err
			}
		case ComponentLayout.Match(parent):
			if _, err := parentChain(&specfile.SourceSet{}, parent, boundary); err != nil {
				return "", "", err
			}
		default:
			return "", "", fmt.Errorf("expected a system or component declaration: %s", parent)
		}
		rel, err := filepath.Rel(filepath.Dir(target), parent)
		if err != nil {
			return "", "", err
		}
		dest := url.URL{Path: filepath.ToSlash(rel)}
		source := fmt.Sprintf("---\ntype: component-spec\n---\n\n# %s\n\n## Description\n\nDescribe this component’s role.\n\n## Part of\n\n[Parent specification](<%s>)\n", name, dest.EscapedPath())
		spec, err := ParseComponent([]byte(source))
		if err != nil {
			return "", "", err
		}
		return source, spec.Name, nil
	})
}

func createSubject(path, name string, layout specfile.Layout, template func(target string) (source, title string, err error)) (string, error) {
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", err
	}
	slug := specfile.NameSlug(name)
	if slug == "" || strings.ContainsAny(name, "\r\n") {
		return "", fmt.Errorf("invalid %s name %q", layout.Kind, name)
	}
	target := layout.Path(dir, slug)
	if owner := specfile.DeclarationDir(target); owner != dir && filepath.Dir(target) != dir {
		return "", fmt.Errorf("destination belongs to another project or repository: %s", owner)
	}
	dir = specfile.DeclarationDir(target)
	source, title, err := template(target)
	if err != nil {
		return "", err
	}
	if err := layout.CheckName(target, title); err != nil {
		return "", err
	}
	locks := filepath.Join(os.TempDir(), fmt.Sprintf("dotspec-subject-locks-%d", os.Getuid()))
	if err := os.MkdirAll(locks, 0700); err != nil {
		return "", err
	}
	unlock, err := lockedfile.MutexAt(filepath.Join(locks, fmt.Sprintf("%x.lock", sha256.Sum256([]byte(dir))))).Lock()
	if err != nil {
		return "", err
	}
	defer unlock()
	existing, err := specfile.LocalDeclaration(dir)
	if err != nil {
		return "", err
	}
	if existing != "" {
		return "", fmt.Errorf("%s already exists; refusing to overwrite it", existing)
	}
	rel, err := filepath.Rel(dir, target)
	if err != nil {
		return "", err
	}
	return specfile.Create(dir, rel, source)
}
