package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const Folder = "specs"

var collectionSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Layout describes either a singleton filename or a collection of slug-named
// files. The owning type supplies its kind and optional example slug.
type Layout struct {
	NamedSubject bool
	Kind         string
	Collection   string
	Filename     string
	Suffix       string
	SlugExample  string
}

// ValidateSlug checks a user-supplied collection slug. Named subject filenames
// follow the Unicode heading conversion in NameSlug.
func (l Layout) ValidateSlug(slug string) error {
	if collectionSlug.MatchString(slug) {
		return nil
	}
	advice := "use lowercase letters or digits separated by single hyphens"
	if l.SlugExample != "" {
		advice += " (for example, " + l.SlugExample + ")"
	}
	return fmt.Errorf("invalid %s slug %q: %s", l.Kind, slug, advice)
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
