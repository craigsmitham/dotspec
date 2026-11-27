package system

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var ComponentLayout = specfile.ComponentLayout

type Component struct {
	Name       string
	ID         string
	Parent     specfile.Reference
	References *specfile.ReferenceDocument
}

func ParseComponent(data []byte) (*Component, error) {
	return ParseComponentMarkdown(specfile.NewMarkdown(data))
}

func ParseComponentMarkdown(input *specfile.Markdown) (*Component, error) {
	var meta metadata
	doc, err := input.TypedOutline("component", &meta)
	if err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		if section.Title == "Specification system" {
			return nil, diagnostic.Error("structure.part-of", "a component selects its governing system through Part of", doc.Location(section.Start, section.Title))
		}
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "component", input.MetadataLocation("id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Required interfaces", "Provided interfaces", "Concept uses", "Synchronizations"); err != nil {
		return nil, err
	}
	ids.AddAnchors(doc.Anchors)
	parent, present, err := doc.SingleSubjectLink("Part of", "structure.part-of")
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, diagnostic.Error("structure.part-of", "component requires a Part of link to its parent system or component", doc.NameLocation)
	}
	if err := doc.ValidateReferences(); err != nil {
		return nil, err
	}
	return &Component{Name: doc.Name, ID: meta.ID, Parent: parent, References: doc.References()}, nil
}

func readComponent(files *specfile.SourceSet, path string) (*Component, error) {
	owner := ComponentLayout.ContextDir(path)
	if _, err := specfile.LocalDeclaration(owner); err != nil {
		return nil, err
	}
	input, err := files.Markdown(path)
	if err != nil {
		return nil, err
	}
	c, err := ParseComponentMarkdown(input)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	if err := ComponentLayout.CheckName(path, c.Name); err != nil {
		return nil, err
	}
	return c, nil
}

func ComponentCheckType() check.Type {
	kind := check.ReadType(ComponentLayout, readComponent, func(c *Component) *specfile.ReferenceDocument { return c.References }, nil)
	kind.Context = CheckContext
	return kind
}

// SubjectChain returns the component declarations from the local owner to its
// enclosing system. Directory discovery selects only the first declaration.
func SubjectChain(owner string) ([]string, error) {
	return subjectChain(&specfile.SourceSet{}, owner)
}

func subjectChain(files *specfile.SourceSet, owner string) ([]string, error) {
	root, err := specfile.ContextRoot(owner)
	if err != nil {
		return nil, err
	}
	path, err := specfile.NearestDeclaration(owner, root)
	if err != nil || path == "" {
		return nil, err
	}
	if ComponentLayout.Match(path) {
		return parentChain(files, path, root)
	}
	_, err = ReadFrom(files, path)
	return []string{path}, err
}

func parentChain(files *specfile.SourceSet, path, boundary string) ([]string, error) {
	var paths []string
	seen := map[string]bool{}
	for {
		if seen[path] {
			return paths, fmt.Errorf("cycle in Part of links at %s", path)
		}
		seen[path] = true
		paths = append(paths, path)
		c, err := readComponent(files, path)
		if err != nil {
			return paths, err
		}
		ref := c.Parent
		ref.Location.Path = path
		u, err := url.Parse(ref.Destination)
		if err != nil || u.IsAbs() || u.Host != "" || u.Path == "" || filepath.IsAbs(u.Path) || strings.Contains(u.Path, "\\") || u.RawQuery != "" {
			return paths, diagnostic.Error("structure.part-of", "Part of must link to a local subject declaration", ref.Location)
		}
		path = filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(u.Path)))
		info, ok, err := check.ReferenceTarget(path, boundary)
		if err != nil || !ok || !info.Mode().IsRegular() {
			return paths, diagnostic.Error("structure.part-of", "unavailable parent declaration: "+path, ref.Location)
		}
		parentBoundary, err := specfile.ContextRoot(specfile.DeclarationDir(path))
		if err != nil {
			return paths, err
		}
		if parentBoundary != boundary {
			return paths, diagnostic.Error("structure.part-of", "parent belongs to another project", ref.Location)
		}
		var parentID string
		switch {
		case Layout.Match(path):
			parent, err := ReadFrom(files, path)
			if err != nil {
				return paths, err
			}
			parentID = parent.ID
		case ComponentLayout.Match(path):
			parent, err := readComponent(files, path)
			if err != nil {
				return paths, err
			}
			parentID = parent.ID
		default:
			return paths, diagnostic.Error("structure.part-of", "expected a system or component declaration", ref.Location)
		}
		if u.Fragment != "" && u.Fragment != parentID {
			md, err := files.Markdown(path)
			if err != nil {
				return paths, err
			}
			if u.Fragment != md.TitleAnchor() {
				return paths, diagnostic.Error("structure.part-of", "Part of must identify the parent declaration itself", ref.Location)
			}
		}
		if Layout.Match(path) {
			return append(paths, path), nil
		}
	}
}
