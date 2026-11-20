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

var ComponentLayout = specfile.ComponentLayout

type Component struct {
	AssessmentPlans []specfile.AssessmentPlan `json:"assessment_plans,omitempty"`
	Name            string
	ID              string
	Parent          specfile.Reference
	References      *specfile.ReferenceDocument
}

func ParseComponent(data []byte) (*Component, error) {
	var meta metadata
	doc, err := specfile.ParseOutline(data, "component", &meta)
	if err != nil {
		return nil, err
	}
	for _, section := range doc.Sections {
		if err := rejectInlineSubject(doc, section); err != nil {
			return nil, err
		}
		if section.Title == "Specification system" {
			return nil, diagnostic.Error("structure.part-of", "a component selects its governing system through Part of", doc.Location(section.Start, section.Title))
		}
	}
	var ids specfile.Identifiers
	if err := ids.Register(meta.ID, "component", specfile.MetadataLocation(data, "id")); err != nil {
		return nil, err
	}
	if err := ids.RegisterDeclarations(doc, "Concept uses", "Synchronizations"); err != nil {
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
	return &Component{AssessmentPlans: doc.AssessmentPlans, Name: doc.Name, ID: meta.ID, Parent: parent, References: doc.References()}, nil
}

func ReadComponent(path string) (*Component, error) {
	owner := ComponentLayout.ContextDir(path)
	if _, err := specfile.LocalDeclaration(owner); err != nil {
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
	path, err := specfile.NearestDeclaration(owner, root)
	if err != nil || path == "" {
		return nil, err
	}
	if ComponentLayout.Match(path) {
		return parentChain(path, root)
	}
	_, err = Read(path)
	return []string{path}, err
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
