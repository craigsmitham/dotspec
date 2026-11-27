package check

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

var referenceScheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

// inspectReferences follows one edge from each successfully parsed source and
// its system context, including modeled passages of explicitly included glossaries.
func (r *Report) inspectReferences(sources []string, discovered map[string]bool, types []Type) {
	sources = append([]string{}, sources...)
	for _, path := range r.ContextFiles {
		if r.Scope.Project == nil || path != r.Scope.Project.Manifest() {
			sources = append(sources, path)
		}
	}
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
		if kind := MatchType(path, types); kind != nil && kind.References != nil {
			result.doc, result.err = kind.References(r.sources, path)
		} else {
			var md *specfile.Markdown
			md, result.err = r.sources.Markdown(path)
			if result.err == nil {
				result.doc = md.AnchorReferences()
			}
		}
		cache[path] = result
		return result.doc, result.err
	}
	visited := make(map[string]bool)
	glossaries := map[string]bool{}
	for len(sources) > 0 {
		source := sources[0]
		sources = sources[1:]
		if visited[source] || failed[source] {
			continue
		}
		kind := MatchType(source, types)
		if !glossaries[source] && (kind == nil || kind.References == nil) {
			continue
		}
		visited[source] = true
		var doc *specfile.ReferenceDocument
		var err error
		if glossaries[source] {
			var md *specfile.Markdown
			md, err = r.sources.Markdown(source)
			if err == nil {
				doc, err = md.GlossaryDetailReferences()
			}
		} else {
			doc, err = read(source)
		}
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
			info, allowed, err := specfile.ReferenceTarget(target, root)
			if err != nil {
				if !allowed {
					r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.boundary", err.Error(), location))
				} else if errors.Is(err, fs.ErrNotExist) {
					count(target)
					if r.missing == nil {
						r.missing = make(map[string]bool)
					}
					r.missing[target] = true
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
			r.collectFile(target, &specfile.Reference{Destination: destination, Location: location, Role: link.Role})
			if location.Element == specfile.GlossarySourcesField && kind != nil && kind.Layout.Kind == "system" && !glossaries[target] {
				glossaries[target] = true
				sources = append(sources, target)
			}
			targetKind := MatchType(target, types)
			markdown := strings.EqualFold(filepath.Ext(target), ".md") || strings.EqualFold(filepath.Ext(target), ".markdown")
			// A companion fragment names a definition in a format dotspec can index.
			// Other native formats keep their own fragment semantics.
			if location.Element == specfile.CompanionsField && u.Fragment != "" && targetKind == nil && !markdown && specfile.NativeDefinitionFormat(target) {
				line, err := specfile.NativeDefinitionAt(target, u.Fragment)
				if err != nil {
					r.record(err, "scan.reference", target)
				} else if line == 0 {
					r.Diagnostics = append(r.Diagnostics, *diagnostic.Error("reference.unresolved", fmt.Sprintf("companion %q names no definition %q in %s", destination, u.Fragment, filepath.Base(target)), location, diagnostic.Location{Path: target}))
				}
				continue
			}
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
