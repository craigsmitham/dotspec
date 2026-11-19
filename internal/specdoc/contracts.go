package specdoc

import (
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"
)

// addSourceContracts runs against the complete subject, heading, and collection
// inventory. A renderer that skips a body, membership list, or relationship
// cannot make that expectation disappear by omitting its formatting step.
func addSourceContracts(p *planner) error {
	for _, s := range p.subjects {
		page := p.plan.Pages[s.file]
		if page == nil {
			return fmt.Errorf("%s: no generated destination for %s", s.source.path, s.name.title)
		}
		page.Type = SpecificationSection
		if s.kind() == "System" {
			page.Type = RequirementsSpecification
		}
		page.SubjectKind = kindName(s.kind())
		page.Title = s.name.title
		page.RequiredContent = nil
		page.RequiredSectionContent = map[string]string{}
		for _, h := range readingOrder(s) {
			if h.owner != s {
				continue
			}
			dest, _ := url.Parse(p.targets[h.node])
			if dest != nil && dest.Fragment != "" {
				page.RequiredAnchors = append(page.RequiredAnchors, dest.Fragment)
				if h.obligation {
					source := relativeURL(s.file, s.source.path, h.anchor)
					requireRoute(page, source, "Edit obligation")
					requireSectionRoute(page, dest.Fragment, source)
				}
			}
			for _, n := range s.source.links {
				if position(n, "start") < h.end || position(n, "end") > h.bodyEnd {
					continue
				}
				raw := attr(n, "destination")
				if raw == "" {
					continue
				}
				link := p.documentLink(s.file, s.source.path, raw)
				switch attr(n, "node_kind") {
				case "link":
					requireRoute(page, link, "")
				case "image":
					page.RequiredResources = append(page.RequiredResources, link)
				}
			}
			if h.end < h.bodyEnd {
				prose := sourceProse(s.source, h)
				page.RequiredContent = append(page.RequiredContent, prose)
				if dest != nil && dest.Fragment != "" {
					page.RequiredSectionContent[dest.Fragment] = prose
				}
			}
		}
		requireRoute(page, relativeURL(s.file, s.source.path, s.name.anchor), "Edit source")
		if s != p.governing {
			requireRoute(page, relative(s.file, s.system.file), s.system.name.title)
		}
	}
	for _, coll := range p.collections {
		index := p.plan.Pages[coll.file]
		if index == nil {
			return fmt.Errorf("%s: collection index is missing", coll.file)
		}
		index.Type = CollectionIndex
		index.SubjectKind = ""
		index.Members = nil
		// Only system and contents links are navigation. Any other visible
		// destination must correspond to an expected collection member.
		index.RequiredLinks = nil
		requireRoute(index, relative(coll.file, coll.system.file), coll.system.name.title)
		requireRoute(index, relative(coll.file, contentsPath), "Contents")
		for _, e := range coll.entries {
			dest := relative(coll.file, e.target)
			index.Members = append(index.Members, dest)
			requireRoute(index, dest, e.title)
		}
		requireRoute(p.plan.Pages[coll.system.file], relative(coll.system.file, coll.file), "")
		contents := p.plan.Pages[contentsPath]
		if contents == nil {
			return fmt.Errorf("%s: specification contents is missing", contentsPath)
		}
		requireRoute(contents, relative(contentsPath, coll.file), "")
	}
	// Read typed relation facts independently of the renderer's choice to
	// present a relation inline, beside a section or in a related-items list.
	for _, a := range p.snapshot.Assertions {
		projection := attr(a, "projection_path")
		_, field, _ := strings.Cut(projection, ".")
		if _, ok := relations[field]; !ok || projection == "QualityModel.quality_attributes" || projection == "System.quality_attributes" {
			continue
		}
		left, right := attr(a, "subject"), attr(a, "object")
		for _, edge := range [][2]string{{left, right}, {right, left}} {
			from := p.ownerPage(p.elements[edge[0]])
			source, destination := p.targets[edge[0]], p.targets[edge[1]]
			if from == nil || source == "" || destination == "" || destination == from.file || destination == source || p.nodes[nameNode(p.elements[edge[1]])] == nil {
				continue
			}
			page := p.plan.Pages[from.file]
			dest := relative(from.file, destination)
			requireRoute(page, dest, "")
			if _, fragment, ok := strings.Cut(source, "#"); ok {
				requireSectionRoute(page, fragment, dest)
			}
		}
	}
	for required, role := range map[string]string{"SPEC.md": RequirementsSpecification, contentsPath: TableOfContents, reviewPath: RequirementsReviewReport} {
		page := p.plan.Pages[required]
		if page == nil {
			return fmt.Errorf("%s: required document is missing", required)
		}
		page.Type = role
		if required != "SPEC.md" {
			page.SubjectKind = ""
		}
	}
	root := p.plan.Pages["SPEC.md"]
	requireRoute(root, relative("SPEC.md", contentsPath), "")
	requireRoute(root, relative("SPEC.md", reviewPath), "")
	return nil
}

func requireSectionRoute(page *Page, fragment, dest string) {
	if page.RequiredSectionLinks == nil {
		page.RequiredSectionLinks = map[string][]string{}
	}
	if !slices.Contains(page.RequiredSectionLinks[fragment], dest) {
		page.RequiredSectionLinks[fragment] = append(page.RequiredSectionLinks[fragment], dest)
	}
}

func requireRoute(page *Page, dest, label string) {
	if !slices.Contains(page.RequiredLinks, dest) {
		page.RequiredLinks = append(page.RequiredLinks, dest)
	}
	if label == "" {
		return
	}
	if page.RequiredLabels == nil {
		page.RequiredLabels = map[string][]string{}
	}
	if !slices.Contains(page.RequiredLabels[dest], label) {
		page.RequiredLabels[dest] = append(page.RequiredLabels[dest], label)
	}
}

// documentLink preserves external references and resolves an authored relative
// URL through the canonical target map before rebasing it for a generated page.
func (p *planner) documentLink(pagePath, sourcePath, raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") {
		return raw
	}
	target := sourcePath
	if u.Path != "" {
		target = path.Join(path.Dir(sourcePath), u.Path)
	}
	dest := relativeURL(pagePath, target, u.Fragment)
	if canonical := p.targets[sourceKey(target, u.Fragment)]; canonical != "" {
		dest = relative(pagePath, canonical)
	}
	if u.RawQuery != "" || u.ForceQuery {
		out, err := url.Parse(dest)
		if err == nil {
			out.RawQuery = u.RawQuery
			out.ForceQuery = u.ForceQuery
			return out.String()
		}
	}
	return dest
}
