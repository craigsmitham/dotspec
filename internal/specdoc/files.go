package specdoc

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/rogpeppe/go-internal/lockedfile"
)

// Changes lists generated paths relative to the selected project.
type Changes struct {
	Written []string
	Removed []string
}

type inventory struct {
	files      map[string][]byte
	changes    Changes
	collisions []diagnostic.Diagnostic
}

func projectDirectory(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}

func outputPath(root, rel string) (string, error) {
	local, err := filepath.Localize(rel)
	if err != nil || !filepath.IsLocal(local) || (rel != "SPEC.md" && (!strings.HasPrefix(rel, "specification/") || !strings.HasSuffix(rel, ".md"))) {
		return "", fmt.Errorf("invalid generated spec path %q", rel)
	}
	// Inspect every ancestor before touching its descendants. Paths outside the
	// generated tree (including symlink targets) are never build destinations.
	for p := root; ; {
		part, rest, _ := strings.Cut(local, string(filepath.Separator))
		p = filepath.Join(p, part)
		info, err := os.Lstat(p)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		if err == nil && ((rest != "" && !info.IsDir()) || (rest == "" && !info.Mode().IsRegular())) {
			return "", diagnostic.Error("spec.collision", "generated destination requires real directories and a regular file", diagnostic.Location{Path: p})
		}
		if info != nil && info.IsDir() {
			boundary, err := nestedBoundary(p)
			if err != nil {
				return "", err
			}
			if boundary {
				return "", diagnostic.Error("spec.collision", "generated path enters a separate project or repository", diagnostic.Location{Path: p})
			}
		}
		if rest == "" {
			return p, nil
		}
		local = rest
	}
}

func nestedBoundary(dir string) (bool, error) {
	for _, name := range []string{"dotspec.yaml", ".git"} {
		if _, err := os.Lstat(filepath.Join(dir, name)); err == nil {
			return true, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
	}
	return false, nil
}

func inspectFiles(root string, plan *Plan) (*inventory, error) {
	result := &inventory{files: map[string][]byte{}}
	for _, rel := range slices.Sorted(maps.Keys(plan.Pages)) {
		path, err := outputPath(root, rel)
		var d *diagnostic.Diagnostic
		if errors.As(err, &d) {
			result.collisions = append(result.collisions, *d)
			continue
		}
		if err != nil {
			return nil, err
		}
		actual, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if err == nil {
			result.files[rel] = actual
			header, _ := splitHeader(actual)
			if header.GeneratedBy != "dotspec" {
				result.collisions = append(result.collisions, *diagnostic.Error("spec.collision", "authored file occupies a generated destination", diagnostic.Location{Path: path}))
				continue
			}
		}
		// A checkout may materialize Git text files using native line endings.
		if !bytes.Equal(bytes.ReplaceAll(actual, []byte("\r\n"), []byte("\n")), []byte(plan.Pages[rel].Markdown)) {
			result.changes.Written = append(result.changes.Written, rel)
		}
	}
	tree := filepath.Join(root, "specification")
	err := filepath.WalkDir(tree, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if path == tree && errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if entry.IsDir() {
			boundary, err := nestedBoundary(path)
			if err != nil {
				return err
			}
			if boundary {
				return fs.SkipDir
			}
		}
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if _, ok := plan.Pages[rel]; ok {
			return nil
		}
		actual, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		header, _ := splitHeader(actual)
		if header.GeneratedBy == "dotspec" {
			result.files[rel] = actual
			result.changes.Removed = append(result.changes.Removed, rel)
		}
		return nil
	})
	return result, err
}

// Check reads the actual output, compares the complete generated file set, and
// runs semantic checks over the Markdown a reader would receive.
func Check(root string, plan *Plan) ([]diagnostic.Diagnostic, error) {
	root, err := projectDirectory(root)
	if err != nil {
		return nil, err
	}
	i, err := inspectFiles(root, plan)
	if err != nil {
		return nil, err
	}
	ds := i.collisions
	for _, rel := range i.changes.Written {
		message := "generated page is stale"
		if _, ok := i.files[rel]; !ok {
			message = "generated page is missing"
		}
		ds = append(ds, *diagnostic.Error("spec.stale", message+"; run dotspec build in this project", diagnostic.Location{Path: filepath.Join(root, filepath.FromSlash(rel))}))
	}
	for _, rel := range i.changes.Removed {
		ds = append(ds, *diagnostic.Error("spec.obsolete", "obsolete generated page; run dotspec build in this project", diagnostic.Location{Path: filepath.Join(root, filepath.FromSlash(rel))}))
	}
	more, err := validateRendered(root, plan, i.files)
	return append(ds, more...), err
}

// Write validates the complete candidate, preflights every destination, and
// installs only changed pages. An unsuccessful install restores previous files.
func Write(root string, plan *Plan) (Changes, error) {
	root, err := projectDirectory(root)
	if err != nil {
		return Changes{}, err
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return Changes{}, err
	}
	lockDir := filepath.Join(cache, "dotspec", "build-locks")
	if err := os.MkdirAll(lockDir, 0700); err != nil {
		return Changes{}, err
	}
	lock := lockedfile.MutexAt(filepath.Join(lockDir, fmt.Sprintf("%x.lock", sha256.Sum256([]byte(root)))))
	unlock, err := lock.Lock()
	if err != nil {
		return Changes{}, err
	}
	defer unlock()
	candidate := make(map[string][]byte, len(plan.Pages))
	for rel, p := range plan.Pages {
		candidate[rel] = []byte(p.Markdown)
	}
	i, err := inspectFiles(root, plan)
	if err != nil {
		return Changes{}, err
	}
	if len(i.collisions) > 0 {
		return Changes{}, &i.collisions[0]
	}
	ds, err := validateRendered(root, plan, candidate)
	if err != nil {
		return Changes{}, err
	}
	if len(ds) > 0 {
		return Changes{}, &ds[0]
	}
	changes := i.changes
	if len(changes.Written)+len(changes.Removed) == 0 {
		return changes, nil
	}
	if err := installPages(root, candidate, changes); err != nil {
		return Changes{}, err
	}
	return changes, nil
}

func installPages(root string, files map[string][]byte, changes Changes) (err error) {
	stage, err := os.MkdirTemp(root, ".dotspec-build-")
	if err != nil {
		return err
	}
	// If rollback itself encounters a filesystem failure, keep the backups and
	// include their path in the returned error.
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.RemoveAll(stage)
		}
	}()
	stageName := func(group, rel string) string { return filepath.Join(stage, group, filepath.FromSlash(rel)) }
	for _, rel := range changes.Written {
		tmp := stageName("new", rel)
		if err := os.MkdirAll(filepath.Dir(tmp), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(tmp, files[rel], 0644); err != nil {
			return err
		}
	}
	type applied struct {
		dest, backup string
		installed    bool
	}
	var edits []applied
	var madeDirs []string
	defer func() {
		if err == nil {
			return
		}
		var rollback error
		for _, edit := range slices.Backward(edits) {
			if edit.installed {
				rollback = errors.Join(rollback, os.Remove(edit.dest))
			}
			if edit.backup != "" {
				rollback = errors.Join(rollback, os.Rename(edit.backup, edit.dest))
			}
		}
		for _, dir := range slices.Backward(madeDirs) {
			_ = os.Remove(dir)
		}
		if rollback != nil {
			cleanup = false
			err = errors.Join(err, fmt.Errorf("restoring generated pages from %s: %w", stage, rollback))
		}
	}()
	for _, rel := range slices.Concat(changes.Written, changes.Removed) {
		dest, err := outputPath(root, rel)
		if err != nil {
			return err
		}
		if _, err := os.Lstat(dest); err == nil {
			// Recheck ownership immediately before moving an existing file.
			data, err := os.ReadFile(dest)
			if err != nil {
				return err
			}
			h, _ := splitHeader(data)
			if h.GeneratedBy != "dotspec" {
				return diagnostic.Error("spec.collision", "generated destination changed during build", diagnostic.Location{Path: dest})
			}
			backup := stageName("old", rel)
			if err := os.MkdirAll(filepath.Dir(backup), 0755); err != nil {
				return err
			}
			if err := os.Rename(dest, backup); err != nil {
				return err
			}
			edits = append(edits, applied{dest: dest, backup: backup})
		} else if errors.Is(err, fs.ErrNotExist) {
			edits = append(edits, applied{dest: dest})
		} else {
			return err
		}
		if _, ok := files[rel]; ok {
			// Create parents top down and remember which ones were introduced.
			var missing []string
			for d := filepath.Dir(dest); d != root; d = filepath.Dir(d) {
				if _, err := os.Lstat(d); err == nil {
					break
				} else if !errors.Is(err, fs.ErrNotExist) {
					return err
				}
				missing = append(missing, d)
			}
			for _, d := range slices.Backward(missing) {
				if err := os.Mkdir(d, 0755); err != nil {
					return err
				}
				madeDirs = append(madeDirs, d)
			}
			if err := os.Rename(stageName("new", rel), dest); err != nil {
				return err
			}
			edits[len(edits)-1].installed = true
		}
	}
	// Remove empty former collections. os.Remove preserves any authored files
	// or other siblings added to these directories.
	for _, rel := range changes.Removed {
		for d := filepath.Dir(filepath.Join(root, filepath.FromSlash(rel))); d != root; d = filepath.Dir(d) {
			if err := os.Remove(d); err != nil {
				break
			}
		}
	}
	return nil
}
