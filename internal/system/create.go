package system

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotspec-md/dotspec/internal/specfile"
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

// InitProject declares the directory at path as a specification project and
// creates its root system specification. Nothing is written when either file
// already exists.
func InitProject(path, name, project string) (manifest, root string, err error) {
	dir, err := specfile.Directory(path)
	if err != nil {
		return "", "", err
	}
	if specfile.NameSlug(name) == "" || strings.ContainsAny(name, "\r\n") {
		return "", "", fmt.Errorf("invalid %s name %q", Layout.Kind, name)
	}
	if !specfile.ValidProjectName(project) {
		return "", "", fmt.Errorf("invalid project name %q", project)
	}
	manifest = filepath.Join(dir, specfile.ProjectFilename)
	if _, err := os.Lstat(manifest); err == nil {
		return "", "", fmt.Errorf("%s already exists; refusing to overwrite it", specfile.ProjectFilename)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", "", err
	}
	root, err = InitNamed(dir, name)
	if err != nil {
		return "", "", err
	}
	if _, err := specfile.Create(dir, specfile.ProjectFilename, "name: "+project+"\n"); err != nil {
		return "", "", errors.Join(err, os.Remove(root))
	}
	return manifest, root, nil
}

// InitComponent creates a component whose Part of link names parent, a system
// or component declaration in the same project.
func InitComponent(path, name, parent string) (string, error) {
	return initPart(path, name, parent, ComponentLayout, func(name, link string) string {
		return fmt.Sprintf("---\ntype: component-spec\n---\n\n# %s\n\n## Description\n\nDescribe this component’s role.\n\n## Part of\n\n[Parent specification](<%s>)\n", name, link)
	})
}

// InitSubsystem creates a system whose Part of link names parent, a system
// declaration in the same project. A project's root system comes from
// InitProject.
func InitSubsystem(path, name, parent string) (string, error) {
	return initPart(path, name, parent, Layout, func(name, link string) string {
		return strings.Replace(Template, "# New system", "# "+name, 1) + fmt.Sprintf("\n## Part of\n\n[Parent specification](<%s>)\n", link)
	})
}

// initPart writes a named subject with an explicit Part of link. The parent
// path is relative to the working directory.
func initPart(path, name, parent string, layout specfile.Layout, template func(name, link string) string) (string, error) {
	return createSubject(path, name, layout, func(target string) (string, string, error) {
		// Resolve the working directory the same way as command scopes. A shell can
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
		info, ok, err := specfile.ReferenceTarget(parent, boundary)
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
		case layout == ComponentLayout && ComponentLayout.Match(parent):
			if _, err := parentChain(&specfile.SourceSet{}, parent, boundary); err != nil {
				return "", "", err
			}
		case layout == ComponentLayout:
			return "", "", fmt.Errorf("expected a system or component declaration: %s", parent)
		default:
			return "", "", fmt.Errorf("expected a system declaration: %s", parent)
		}
		rel, err := filepath.Rel(filepath.Dir(target), parent)
		if err != nil {
			return "", "", err
		}
		dest := url.URL{Path: filepath.ToSlash(rel)}
		source := template(name, dest.EscapedPath())
		title, err := subjectName(layout, source)
		if err != nil {
			return "", "", err
		}
		return source, title, nil
	})
}

func subjectName(layout specfile.Layout, source string) (string, error) {
	if layout == ComponentLayout {
		spec, err := ParseComponent([]byte(source))
		if err != nil {
			return "", err
		}
		return spec.Name, nil
	}
	spec, err := Parse([]byte(source))
	if err != nil {
		return "", err
	}
	return spec.Name, nil
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
	if existing == target {
		return "", fmt.Errorf("%s already exists; refusing to overwrite it", existing)
	}
	if existing != "" {
		return "", fmt.Errorf("%s already declares this directory's subject; create the %s in a directory of its own", existing, layout.Kind)
	}
	rel, err := filepath.Rel(dir, target)
	if err != nil {
		return "", err
	}
	return specfile.Create(dir, rel, source)
}
