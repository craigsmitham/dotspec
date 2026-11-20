package specfile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const AttributeSuffix = ".quality-attribute.spec.md"

// QualityFolders returns the attribute folders below the nearest specs directory.
func QualityFolders(path string) []string {
	var folders []string
	for dir := filepath.Dir(path); dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if filepath.Base(dir) == Folder {
			break
		}
		if filepath.Base(dir) == "quality" && filepath.Base(filepath.Dir(dir)) == Folder {
			slices.Reverse(folders)
			return folders
		}
		folders = append(folders, filepath.Base(dir))
	}
	return nil
}

func AttributePath(dir string) string { return filepath.Join(dir, filepath.Base(dir)+AttributeSuffix) }

// LocalPath resolves a slug across the ordinary collection and its quality folders.
func (l Layout) LocalPath(owner, slug string) (string, error) {
	ordinary := l.Path(owner, slug)
	if !l.QualityTree {
		return ordinary, nil
	}
	matches, err := walkPathsIfPresent(filepath.Join(owner, Folder, "quality"), func(name string) bool { return name == l.filename(slug) })
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return ordinary, nil
	}
	if _, err := os.Lstat(ordinary); err == nil {
		matches = append(matches, ordinary)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous %s %q: %s", l.Kind, slug, strings.Join(matches, ", "))
	}
	return matches[0], nil
}

func walkPathsIfPresent(dir string, match func(string) bool) ([]string, error) {
	info, err := os.Lstat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s must be a directory", dir)
	}
	return walkPaths(dir, match)
}
