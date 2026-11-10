package specfile

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func Create(path, filename, template string) (string, error) {
	dir, err := Directory(path)
	if err != nil {
		return "", err
	}
	target := filepath.Join(dir, filename)
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return "", fmt.Errorf("%s already exists; refusing to overwrite it", target)
	}
	if err != nil {
		return "", err
	}
	_, writeErr := file.WriteString(template)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return "", err
	}
	return target, nil
}

// Discover finds matching files using the same ignore and symlink rules for all types.
func Discover(path string, match func(string) bool) ([]string, error) {
	dir, err := Directory(path)
	if err != nil {
		return nil, err
	}
	root, err := RepositoryRoot(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	if root != "" {
		paths, err = gitPaths(dir, root, match)
	} else {
		paths, err = walkPaths(dir, match)
	}
	sort.Strings(paths)
	return paths, err
}

func Directory(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory; --path selects a directory", path)
	}
	return abs, nil
}

func RepositoryRoot(dir string) (string, error) {
	for current := dir; ; current = filepath.Dir(current) {
		_, err := os.Lstat(filepath.Join(current, ".git"))
		if err == nil {
			return current, nil // .git may be a directory or a worktree metadata file.
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		if filepath.Dir(current) == current {
			return "", nil
		}
	}
}

// Let Git apply its own ignore rules, including nested .gitignore files, local
// exclusions, and negations. Tracked files remain visible even if now ignored.
func gitPaths(dir, root string, match func(string) bool) ([]string, error) {
	cmd := exec.Command("git", "-C", dir, "ls-files", "--cached", "--others", "--exclude-standard", "--full-name", "-z", "--", ".")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("discover specs with Git: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	seen := make(map[string]bool)
	var paths []string
	for _, entry := range bytes.Split(output, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(string(entry)))
		if !match(filepath.Base(path)) || seen[path] {
			continue
		}
		seen[path] = true
		usable, err := regularPath(path, dir)
		if err != nil {
			return nil, err
		}
		if usable {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

// Check each component: a tracked directory may have been replaced by a symlink
// since Git's index was written. Also skip deleted tracked files and nested repos.
func regularPath(path, boundary string) (bool, error) {
	for current := path; current != boundary; current = filepath.Dir(current) {
		if filepath.Base(current) == ".git" {
			return false, nil
		}
		info, err := os.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 || (current == path && !info.Mode().IsRegular()) {
			return false, nil
		}
		if current != path {
			if _, err := os.Lstat(filepath.Join(current, ".git")); err == nil {
				return false, nil
			} else if !errors.Is(err, fs.ErrNotExist) {
				return false, err
			}
		}
	}
	return true, nil
}

func walkPaths(dir string, match func(string) bool) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			if path != dir {
				if _, err := os.Lstat(filepath.Join(path, ".git")); err == nil {
					return filepath.SkipDir
				} else if !errors.Is(err, fs.ErrNotExist) {
					return err
				}
			}
			return nil
		}
		if match(entry.Name()) && entry.Type().IsRegular() {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
