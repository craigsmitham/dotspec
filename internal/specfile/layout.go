package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const Folder = "specs"

var collectionSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Layout describes a collection of slug-named files. The owning type supplies
// its kind and optional example slug.
type Layout struct {
	NamedSubject bool
	Kind         string
	Collection   string
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

// ReadMember reads one slug-named collection file. It rejects a filename
// outside the layout before reading, and adds the path to parse failures.
func ReadMember[T any](files *SourceSet, l Layout, path string, parse func(*Markdown) (T, error)) (T, string, error) {
	var zero T
	rule := "layout." + l.Kind + "-name"
	slug, ok := strings.CutSuffix(filepath.Base(path), l.Suffix)
	if !ok {
		return zero, "", diagnostic.Error(rule, fmt.Sprintf("expected a <slug>%s filename", l.Suffix), diagnostic.Location{Path: path})
	}
	if err := l.ValidateSlug(slug); err != nil {
		return zero, "", diagnostic.Error(rule, err.Error(), diagnostic.Location{Path: path})
	}
	input, err := files.Markdown(path)
	if err != nil {
		return zero, "", err
	}
	spec, err := parse(input)
	if err != nil {
		return zero, "", diagnostic.WithPath(err, path)
	}
	return spec, slug, nil
}

func (l Layout) Match(name string) bool {
	name = filepath.Base(name)
	return strings.HasSuffix(name, l.Suffix) || (l.NamedSubject && name == strings.TrimPrefix(l.Suffix, "."))
}

// validateName applies the collection slug grammar; named subjects follow NameSlug.
func (l Layout) validateName(slug string) error {
	if l.NamedSubject {
		return nil
	}
	return l.ValidateSlug(slug)
}

func (l Layout) filename(slug string) string {
	return slug + l.Suffix
}

func (l Layout) Path(owner, slug string) string {
	return filepath.Join(defaultSpecDir(owner), l.Collection, l.filename(slug))
}

// Create writes a new member named by slug in the directory selected by path.
// A collection slug must follow the layout's slug grammar; it is checked before
// the filesystem is touched.
func (l Layout) Create(path, slug, template string) (string, error) {
	if err := l.validateName(slug); err != nil {
		return "", err
	}
	owner, err := Directory(path)
	if err != nil {
		return "", err
	}
	target := l.Path(owner, slug)
	rel, err := filepath.Rel(owner, target)
	if err != nil {
		return "", err
	}
	if !l.NamedSubject {
		template = titled(template, slug)
	}
	return Create(owner, rel, template)
}

// SlugTitle is the draft heading for a collection slug: fast-save becomes
// "Fast save". Authors rename the heading freely; the slug stays the filename.
func SlugTitle(slug string) string {
	title := strings.ReplaceAll(slug, "-", " ")
	if title == "" {
		return title
	}
	return strings.ToUpper(title[:1]) + title[1:]
}

// titled replaces a template's placeholder "# New ..." heading with the slug's
// title, so a fresh draft has a distinct heading and generated page.
func titled(template, slug string) string {
	for _, prefix := range []string{"# New ", "\n# New "} {
		i := strings.Index(template, prefix)
		if i < 0 || (prefix == "# New " && i != 0) {
			continue
		}
		start := i + len(prefix) - len("# New ")
		end := strings.IndexByte(template[start:], '\n')
		if end < 0 {
			end = len(template) - start
		}
		return template[:start] + "# " + SlugTitle(slug) + template[start+end:]
	}
	return template
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
