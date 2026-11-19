package system

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/check"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

var ComponentLayout = specfile.Layout{Kind: "component", Suffix: ".component.spec.md", NamedSubject: true}

type Component struct {
	Name       string
	ID         string
	Parent     specfile.Reference
	References *specfile.ReferenceDocument
}

func ParseComponent(data []byte) (*Component, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, "component", &meta)
	if err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		if section.Title == "Components" && len(section.Children) > 0 {
			return nil, diagnostic.Error("structure.component-source", "declare each component in its own <name>.component.spec.md file", doc.Location(section.Start, section.Title))
		}
		if section.Title == "Specification system" {
			return nil, diagnostic.Error("structure.part-of", "a component selects its governing system through Part of", doc.Location(section.Start, section.Title))
		}
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "component", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Features", "Capabilities", "Concept uses", "Synchronizations"); err != nil {
		return nil, err
	}
	if err := ids.RegisterInterfaces(doc); err != nil {
		return nil, err
	}
	if err := ids.RegisterObligations(doc); err != nil {
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
	if err := specfile.ValidateReferences(doc.Original, doc.Source, doc.AST, doc.Anchors); err != nil {
		return nil, err
	}
	return &Component{Name: doc.Name, ID: meta.ID, Parent: parent, References: doc.References()}, nil
}

func ReadComponent(path string) (*Component, error) {
	owner, err := ComponentLayout.Owner(path)
	if err != nil {
		return nil, err
	}
	if _, err := localSubject(owner); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := ParseComponent(data)
	if err != nil {
		return nil, diagnostic.WithPath(err, path)
	}
	if err := ComponentLayout.CheckName(path, c.Name); err != nil {
		return nil, err
	}
	return c, nil
}

func ComponentCheckType() check.Type {
	return check.Type{Layout: ComponentLayout, Context: CheckContext, Inspect: func(path string) ([]diagnostic.Diagnostic, error) { _, err := ReadComponent(path); return nil, err }, References: func(path string) (*specfile.ReferenceDocument, error) {
		c, err := ReadComponent(path)
		if err != nil {
			return nil, err
		}
		return c.References, nil
	}}
}

// SubjectChain returns the component declarations from the local owner to its
// enclosing system. Directory discovery selects only the first declaration.
func SubjectChain(owner string) ([]string, error) {
	root, err := specfile.ContextRoot(owner)
	if err != nil {
		return nil, err
	}
	path, err := nearestSubject(owner, root)
	if err != nil || path == "" {
		return nil, err
	}
	if ComponentLayout.Match(path) {
		return parentChain(path, root)
	}
	_, err = Read(path)
	return []string{path}, err
}

func nearestSubject(owner, boundary string) (string, error) {
	for dir := owner; ; dir = filepath.Dir(dir) {
		for _, layout := range []specfile.Layout{Layout, ComponentLayout} {
			if err := layout.ValidateLocal(dir); err != nil {
				return "", err
			}
		}
		path, err := localSubject(dir)
		if path != "" || err != nil || dir == boundary || filepath.Dir(dir) == dir {
			return path, err
		}
	}
}

func localSubject(dir string) (string, error) {
	system, err := Layout.LocalSubject(dir)
	if err != nil {
		return "", err
	}
	component, err := ComponentLayout.LocalSubject(dir)
	if err != nil {
		return "", err
	}
	if system != "" && component != "" {
		return "", diagnostic.Error("layout.subject-ambiguous", "directory contains both a system and a component declaration", diagnostic.Location{Path: component}, diagnostic.Location{Path: system})
	}
	if component != "" {
		return component, nil
	}
	return system, nil
}

func parentChain(path, boundary string) ([]string, error) {
	var paths []string
	seen := map[string]bool{}
	for {
		if seen[path] {
			return paths, fmt.Errorf("cycle in Part of links at %s", path)
		}
		seen[path] = true
		paths = append(paths, path)
		c, err := ReadComponent(path)
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
		parentBoundary, err := specfile.ContextRoot(filepath.Dir(filepath.Dir(path)))
		if err != nil {
			return paths, err
		}
		if parentBoundary != boundary {
			return paths, diagnostic.Error("structure.part-of", "parent belongs to another project", ref.Location)
		}
		var parentID string
		switch {
		case Layout.Match(path):
			parent, err := Read(path)
			if err != nil {
				return paths, err
			}
			parentID = parent.ID
		case ComponentLayout.Match(path):
			parent, err := ReadComponent(path)
			if err != nil {
				return paths, err
			}
			parentID = parent.ID
		default:
			return paths, diagnostic.Error("structure.part-of", "expected a system or component declaration", ref.Location)
		}
		if u.Fragment != "" && u.Fragment != parentID {
			data, err := os.ReadFile(path)
			if err != nil {
				return paths, err
			}
			if u.Fragment != specfile.MarkdownTitleAnchor(data) {
				return paths, diagnostic.Error("structure.part-of", "Part of must identify the parent declaration itself", ref.Location)
			}
		}
		if Layout.Match(path) {
			return append(paths, path), nil
		}
	}
}
