package specfile

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const Folder = "specs"

// Layout describes either a singleton filename or a collection of slug-named
// files. The owning type supplies its kind; layout never infers domain identity
// from a directory name. Slug validation remains the type's responsibility.
type Layout struct {
	NamedSubject bool
	Kind         string
	Collection   string
	Filename     string
	Suffix       string
}

func (l Layout) Match(name string) bool {
	name = filepath.Base(name)
	if l.Filename != "" {
		return name == l.Filename
	}
	return strings.HasSuffix(name, l.Suffix) || (l.NamedSubject && name == strings.TrimPrefix(l.Suffix, "."))
}

func (l Layout) filename(slug string) string {
	if l.Filename != "" {
		return l.Filename
	}
	return slug + l.Suffix
}

func (l Layout) Path(owner, slug string) string {
	return l.path(owner, l.filename(slug))
}

func (l Layout) path(owner, filename string) string {
	return filepath.Join(defaultSpecDir(owner), l.Collection, filename)
}

func (l Layout) Create(owner, slug, template string) (string, error) {
	target := l.Path(owner, slug)
	rel, err := filepath.Rel(owner, target)
	if err != nil {
		return "", err
	}
	return Create(owner, rel, template)
}

// DeclarationDir is the directory claimed by a subject declaration.
func DeclarationDir(path string) string {
	dir := filepath.Dir(path)
	if specsAlias(dir) {
		return filepath.Dir(dir)
	}
	return dir
}

// ContextDir supplies the starting directory for local subject resolution.
func (l Layout) ContextDir(path string) string {
	if l.NamedSubject {
		return DeclarationDir(path)
	}
	return filepath.Dir(path)
}

func defaultSpecDir(owner string) string {
	if specsAlias(owner) {
		return owner
	}
	return filepath.Join(owner, Folder)
}

// A project or repository root owns its declarations even when named specs.
func specsAlias(dir string) bool {
	if filepath.Base(dir) != Folder {
		return false
	}
	if _, err := os.Lstat(filepath.Join(dir, ".git")); !errors.Is(err, fs.ErrNotExist) {
		return false
	}
	project, err := hasManifest(dir)
	return err == nil && !project
}
