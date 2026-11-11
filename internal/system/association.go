package system

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

// Associated selects the explicitly linked system, falling back to directory
// discovery only when no selection was authored. It never fetches external URLs.
func Associated(owner, source string) (*System, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	ref, explicit, err := specfile.SpecificationSystem(data)
	if err != nil {
		return nil, diagnostic.WithPath(err, source)
	}
	if !explicit {
		return Enclosing(owner)
	}
	ref.Location.Path = source
	invalid := func(message string) (*System, error) {
		return nil, diagnostic.Error("relationship.specification_system", message, ref.Location)
	}
	u, err := url.Parse(ref.Destination)
	if err != nil || u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") || strings.Contains(u.Path, "\\") || u.RawQuery != "" {
		return invalid("Specification system must link to a repository-local system document.")
	}
	target := filepath.Clean(filepath.Join(filepath.Dir(source), filepath.FromSlash(u.Path)))
	if filepath.Base(target) != Filename {
		return invalid("Specification system must identify a system.spec.md document.")
	}
	root, err := specfile.RepositoryRoot(owner)
	if err != nil {
		return nil, err
	}
	info, allowed, err := check.ReferenceTarget(target, root)
	if !allowed {
		return invalid("Specification system is outside the permitted repository boundary.")
	}
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return invalid("Specification system must identify a regular file.")
	}
	spec, err := Read(target)
	if err != nil {
		return nil, err
	}
	if u.Fragment != "" && u.Fragment != spec.ID {
		targetData, err := os.ReadFile(target)
		if err != nil {
			return nil, err
		}
		if u.Fragment != specfile.MarkdownTitleAnchor(targetData) {
			return invalid("Specification system must identify the system itself, not an entry within it.")
		}
	}
	return spec, nil
}
