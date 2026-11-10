package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const Folder = "specs"

// Layout describes either a singleton filename or a collection of slug-named
// files. The owning type supplies its kind; layout never infers domain identity
// from a directory name. Slug validation remains the type's responsibility.
type Layout struct {
	Kind       string
	Collection string
	Filename   string
	Suffix     string
}

func (l Layout) Match(name string) bool {
	name = filepath.Base(name)
	if l.Filename != "" {
		return name == l.Filename
	}
	return strings.HasSuffix(name, l.Suffix)
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
	return filepath.Join(owner, Folder, l.Collection, filename)
}

func (l Layout) Create(owner, slug, template string) (string, error) {
	return Create(owner, filepath.Join(Folder, l.Collection, l.filename(slug)), template)
}

// Scope resolves a project directory. The reserved specs tree contains metadata,
// so it cannot introduce another project scope.
func Scope(path string) (string, error) {
	dir, err := Directory(path)
	if err != nil {
		return "", err
	}
	metadata, err := metadataDirectory(dir)
	if err != nil {
		return "", err
	}
	if metadata != "" {
		return "", fmt.Errorf("%s is inside a specs directory; --path selects a project directory, use %s", path, filepath.Dir(metadata))
	}
	return dir, nil
}

// Owner validates a spec's placement and returns the directory it describes.
// An empty collection places a spec directly under specs/.
func (l Layout) Owner(path string) (string, error) {
	dir := filepath.Dir(path)
	metadata, err := metadataDirectory(dir)
	if err != nil {
		return "", err
	}
	if metadata != "" && dir == filepath.Join(metadata, l.Collection) {
		return filepath.Dir(metadata), nil
	}
	expected := l.path("<project-dir>", filepath.Base(path))
	if metadata != "" {
		expected = l.path(filepath.Dir(metadata), filepath.Base(path))
	}
	return "", l.misplaced(path, expected)
}

func (l Layout) misplaced(path, expected string) error {
	return diagnostic.Error("layout."+l.Kind+"-location", fmt.Sprintf("misplaced spec: expected %s", expected), diagnostic.Location{Path: path})
}

// ValidateLocal rejects misplaced specs in this scope, including adjacent files
// and singular collection folders. Child project scopes are checked separately.
func (l Layout) ValidateLocal(owner string) error {
	entries, err := os.ReadDir(owner)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && l.Match(entry.Name()) {
			return l.misplaced(filepath.Join(owner, entry.Name()), l.path(owner, entry.Name()))
		}
	}
	metadata := filepath.Join(owner, Folder)
	info, err := os.Lstat(metadata)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return diagnostic.Error("layout.directory", "must be a directory, not a file or symlink", diagnostic.Location{Path: metadata})
	}
	if l.Collection != "" {
		dir := filepath.Join(metadata, l.Collection)
		info, err := os.Lstat(dir)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if err == nil && !info.IsDir() {
			return diagnostic.Error("layout.directory", "must be a directory, not a file or symlink", diagnostic.Location{Path: dir})
		}
	}
	// Like direct reads, local validation does not need Git or apply ignore
	// rules. Discovery retains its separate Git-aware visibility rules.
	paths, err := walkPaths(metadata, l.Match)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, err := l.Owner(path); err != nil {
			return err
		}
	}
	return nil
}

// Return the outermost specs directory inside this repository (or up to the
// filesystem root outside Git), so nested specs/ trees cannot create scopes.
func metadataDirectory(dir string) (string, error) {
	root, err := RepositoryRoot(dir)
	if err != nil {
		return "", err
	}
	var metadata string
	for current := dir; current != root; current = filepath.Dir(current) {
		if filepath.Base(current) == Folder {
			metadata = current
		}
		if filepath.Dir(current) == current {
			break
		}
	}
	return metadata, nil
}
