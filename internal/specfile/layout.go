package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
)

const Folder = "specs"

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

func Path(owner, collection, filename string) string {
	return filepath.Join(owner, Folder, collection, filename)
}

// Owner validates a spec's placement and returns the directory it describes.
// An empty collection denotes the singleton system spec directly under specs/.
func Owner(path, collection string) (string, error) {
	dir := filepath.Dir(path)
	metadata, err := metadataDirectory(dir)
	if err != nil {
		return "", err
	}
	if metadata != "" && dir == filepath.Join(metadata, collection) {
		return filepath.Dir(metadata), nil
	}
	expected := Path("<project-dir>", collection, filepath.Base(path))
	if metadata != "" {
		expected = Path(filepath.Dir(metadata), collection, filepath.Base(path))
	}
	return "", misplaced(path, expected, collection)
}

func misplaced(path, expected, collection string) error {
	kind := "system"
	if collection == "jobs" {
		kind = "job"
	}
	return diagnostic.Error("layout."+kind+"-location", fmt.Sprintf("misplaced spec: expected %s", expected), diagnostic.Location{Path: path})
}

// ValidateLocal rejects misplaced specs in this scope, including adjacent files
// and singular collection folders. Child project scopes are checked separately.
func ValidateLocal(owner, collection string, match func(string) bool) error {
	entries, err := os.ReadDir(owner)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && match(entry.Name()) {
			return misplaced(filepath.Join(owner, entry.Name()), Path(owner, collection, entry.Name()), collection)
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
	if collection != "" {
		dir := filepath.Join(metadata, collection)
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
	paths, err := walkPaths(metadata, match)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, err := Owner(path, collection); err != nil {
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
