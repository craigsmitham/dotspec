package check

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var referenceScheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

// inspectReferences follows one edge from each successfully parsed source and
// its system context. Referenced documents supply targets, not another scan.
func (r *Report) inspectReferences(sources []string, discovered map[string]bool, types []Type) {
	sources = append(append([]string{}, sources...), r.ContextFiles...)
	// Parsing and context lookup have already reported failed sources. Their
	// identities and links are unavailable; do not repeat those failures here.
	failed := make(map[string]bool)
	for _, finding := range r.Diagnostics {
		if finding.Severity == "error" {
			failed[finding.Location.Path] = true
		}
	}
	counted := make(map[string]bool)
	for path := range discovered {
		counted[path] = true
	}
	for _, path := range r.ContextFiles {
		counted[path] = true
	}
	count := func(path string) {
		if !counted[path] {
			counted[path] = true
			r.ContextFiles = append(r.ContextFiles, path)
			r.Summary.ContextFilesChecked++
		}
	}
	type loaded struct {
		doc *specfile.ReferenceDocument
		err error
	}
	cache := make(map[string]loaded)
	read := func(path string) (*specfile.ReferenceDocument, error) {
		if result, ok := cache[path]; ok {
			return result.doc, result.err
		}
		var result loaded
		if kind := findType(path, types); kind != nil && kind.References != nil {
			result.doc, result.err = kind.References(path)
		} else {
			var data []byte
			data, result.err = os.ReadFile(path)
			if result.err == nil {
				result.doc = specfile.MarkdownReferences(data)
			}
		}
		cache[path] = result
		return result.doc, result.err
	}
	visited := make(map[string]bool)
	for _, source := range sources {
		if visited[source] || failed[source] {
			continue
		}
		visited[source] = true
		kind := findType(source, types)
		if kind == nil || kind.References == nil {
			continue
		}
		doc, err := read(source)
		if err != nil {
			r.record(err, "scan.reference", source)
			continue
		}
		root, err := specfile.RepositoryRoot(filepath.Dir(source))
		if err != nil {
			r.record(err, "scan.reference", source)
			continue
		}
		for _, link := range doc.Links {
			destination := link.Destination
			if destination == "" || strings.HasPrefix(destination, "#") || strings.HasPrefix(destination, "/") || referenceScheme.MatchString(destination) {
				continue
			}
			location := link.Location
			location.Path = source
			u, err := url.Parse(destination)
			if err != nil {
				r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.invalid", fmt.Sprintf("invalid reference %q", destination), location))
				continue
			}
			// URLs and absolute paths retain their external navigation semantics.
			if u.IsAbs() || u.Host != "" || u.Path == "" || strings.HasPrefix(u.Path, "/") {
				continue
			}
			if strings.Contains(u.Path, "\\") {
				r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.invalid", "relative references must use forward slashes", location))
				continue
			}
			target := filepath.Clean(filepath.Join(filepath.Dir(source), filepath.FromSlash(u.Path)))
			info, err := referencePath(target, root)
			if err != nil {
				var boundary *referenceBoundary
				if errors.As(err, &boundary) {
					r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.boundary", boundary.Error(), location))
				} else if errors.Is(err, fs.ErrNotExist) {
					count(target)
					r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.unresolved", fmt.Sprintf("reference %q targets a missing file or directory", destination), location, diagnostic.Location{Path: target}))
				} else {
					count(target)
					r.record(err, "scan.reference", target)
				}
				continue
			}
			if info.IsDir() {
				if u.Fragment != "" {
					r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.unresolved", fmt.Sprintf("directory reference %q cannot resolve a fragment", destination), location))
				}
				continue
			}
			count(target)
			r.collectFile(target, &location)
			targetKind := findType(target, types)
			markdown := strings.EqualFold(filepath.Ext(target), ".md") || strings.EqualFold(filepath.Ext(target), ".markdown")
			// Native non-Markdown formats own their fragment semantics.
			if targetKind == nil && (!markdown || u.Fragment == "") {
				continue
			}
			if failed[target] {
				continue
			}
			targetDoc, err := read(target)
			if err != nil {
				r.record(err, "scan.reference", target)
				continue
			}
			if u.Fragment != "" && !targetDoc.Anchors[u.Fragment] {
				r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.unresolved", fmt.Sprintf("reference %q has no matching declared ID or Markdown heading anchor", destination), location, diagnostic.Location{Path: target}))
			}
		}
	}
}

type referenceBoundary struct{ reason string }

func (e *referenceBoundary) Error() string { return e.reason }

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
