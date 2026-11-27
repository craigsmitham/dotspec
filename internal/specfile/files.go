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
	// Create only real directories beneath the selected scope. In particular,
	// init must not write through a symlinked specs/ or collection directory.
	relative, err := filepath.Rel(dir, filepath.Dir(target))
	if err != nil {
		return "", err
	}
	parent := dir
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		parent = filepath.Join(parent, component)
		if err := os.Mkdir(parent, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
			return "", err
		}
		info, err := os.Lstat(parent)
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			return "", fmt.Errorf("%s must be a directory, not a file or symlink", parent)
		}
		if parent != dir {
			if err := CheckProjectBoundary(parent); err != nil {
				return "", err
			}
		}

	}
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

// CheckProjectBoundary prevents an initializer from writing through a nested project or repository root.
func CheckProjectBoundary(dir string) error {
	project, err := hasManifest(dir)
	if err != nil {
		return err
	}
	if project {
		return fmt.Errorf("destination belongs to another project: %s", dir)
	}
	_, err = os.Lstat(filepath.Join(dir, ".git"))
	if err == nil {
		return fmt.Errorf("destination belongs to another repository: %s", dir)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

// Discover finds matching files using the same ignore and symlink rules for all types.
func Discover(path string, match func(string) bool) ([]string, error) {
	return discover(path, ".", match)
}

func discover(path, pathspec string, match func(string) bool) ([]string, error) {
	dir, err := Directory(path)
	if err != nil {
		return nil, err
	}
	root, err := RepositoryRoot(dir)
	if err != nil {
		return nil, err
	}
	// Every nested manifest starts another project, even if that project cannot
	// currently be parsed. A parent scan only needs to know its boundary.
	boundaries := map[string]bool{}
	skip := func(child string) (bool, error) {
		if value, ok := boundaries[child]; ok {
			return value, nil
		}
		found, err := hasManifest(child)
		boundaries[child] = found
		return found, err
	}
	var paths []string
	if root != "" {
		paths, err = gitPaths(dir, root, pathspec, match, skip)
	} else {
		paths, err = walkSelectedPaths(dir, match, skip)
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
func gitPaths(dir, root, pathspec string, match func(string) bool, skip func(string) (bool, error)) ([]string, error) {
	cmd := exec.Command("git", "-C", dir, "ls-files", "--cached", "--others", "--exclude-standard", "--full-name", "-z", "--", pathspec)
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
		usable, err := regularPath(path, dir, skip)
		if err != nil {
			return paths, err
		}
		if usable {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

// Check each component: a tracked directory may have been replaced by a symlink
// since Git's index was written. Also skip deleted tracked files and nested repos.
func regularPath(path, boundary string, skip func(string) (bool, error)) (bool, error) {
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
			if excluded, err := skip(current); err != nil || excluded {
				return false, err
			}
			if _, err := os.Lstat(filepath.Join(current, ".git")); err == nil {
				return false, nil
			} else if !errors.Is(err, fs.ErrNotExist) {
				return false, err
			}
		}
	}
	return true, nil
}

type referenceBoundary struct{ reason string }

func (e *referenceBoundary) Error() string { return e.reason }

// ReferenceTarget applies the repository boundary policy for explicit links
// without reading a target. Check, snapshot extraction and reports share it.
// allowed is false only for a policy boundary; missing and unreadable targets
// remain distinct from prohibited traversal.
func ReferenceTarget(path, root string) (info os.FileInfo, allowed bool, err error) {
	info, err = referencePath(path, root)
	var boundary *referenceBoundary
	return info, !errors.As(err, &boundary), err
}

// Explicit links may address ignored files and ancestors within a repository,
// but must not traverse symlinks, Git internals, or another repository. Outside
// Git the filesystem root is the boundary, matching system context lookup.
func referencePath(path, root string) (os.FileInfo, error) {
	if root != "" {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, &referenceBoundary{"relative reference leaves the repository; use an external URL for external navigation"}
		}
	}
	var target os.FileInfo
	// Walk from root to leaf to avoid touching a path beneath a symlink first.
	var components []string
	for current := path; current != root && filepath.Dir(current) != current; current = filepath.Dir(current) {
		components = append(components, current)
	}
	for i := len(components) - 1; i >= 0; i-- {
		current := components[i]
		if filepath.Base(current) == ".git" {
			return nil, &referenceBoundary{"relative reference targets Git metadata"}
		}
		info, err := os.Lstat(current)
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, &referenceBoundary{"relative reference traverses a symlink: " + current}
		}
		if !info.IsDir() && i > 0 {
			return nil, fs.ErrNotExist
		}
		if info.IsDir() && current != root {
			if _, err := os.Lstat(filepath.Join(current, ".git")); err == nil {
				return nil, &referenceBoundary{"relative reference enters another repository: " + current}
			} else if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
		}
		target = info
	}
	if target == nil {
		return os.Lstat(path)
	}
	if !target.IsDir() && !target.Mode().IsRegular() {
		return nil, &referenceBoundary{"relative reference must target a regular file or directory"}
	}
	return target, nil
}

func walkSelectedPaths(dir string, match func(string) bool, skip func(string) (bool, error)) ([]string, error) {
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
				excluded, err := skip(path)
				if err != nil {
					return err
				}
				if excluded {
					return filepath.SkipDir
				}
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
