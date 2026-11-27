package system

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/dotspec-md/dotspec/internal/check"
	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
)

var ComponentLayout = specfile.ComponentLayout

type Component struct {
	Name       string
	ID         string
	Parent     specfile.Reference
	References *specfile.ReferenceDocument
	advice     []diagnostic.Diagnostic
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
	return &Component{Name: doc.Name, ID: meta.ID, Parent: parent, References: doc.References(), advice: boundaryAdvice(doc)}, nil
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
	c.advice = withPath(c.advice, path)
	return c, nil
}

func ComponentCheckType() check.Type {
	kind := check.ReadType(ComponentLayout, readComponent, func(c *Component) *specfile.ReferenceDocument { return c.References }, func(c *Component) []diagnostic.Diagnostic { return c.advice })
	kind.Context = CheckContext
	return kind
}

// ComponentMember is a discovered component with the declaration its Part of
// link names and the system at the end of that chain.
type ComponentMember struct {
	Name   string
	Path   string
	Parent string
	System *System
}

// Components discovers the components below path in source order, each with
// its explicit Part of parent and governing system.
func Components(path string) ([]ComponentMember, error) {
	dir, err := specfile.Directory(path)
	if err != nil {
		return nil, err
	}
	paths, err := specfile.Discover(dir, ComponentLayout.Match)
	if err != nil {
		return nil, err
	}
	files := &specfile.SourceSet{}
	members := []ComponentMember{}
	for _, file := range paths {
		boundary, err := specfile.ContextRoot(ComponentLayout.ContextDir(file))
		if err != nil {
			return nil, err
		}
		chain, err := parentChain(files, file, boundary)
		if err != nil {
			return nil, err
		}
		c, err := readComponent(files, file)
		if err != nil {
			return nil, err
		}
		governing, err := ReadFrom(files, chain[len(chain)-1])
		if err != nil {
			return nil, err
		}
		members = append(members, ComponentMember{Name: c.Name, Path: file, Parent: chain[1], System: governing})
	}
	return members, nil
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
		info, ok, err := specfile.ReferenceTarget(path, boundary)
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
