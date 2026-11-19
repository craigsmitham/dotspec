package specdoc

import (
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
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
		page.Document = subjectDocument(s.file, s.source, s.name)
		page.SectionDocs = map[string]Document{}
		page.DirectoryMembers = map[string][]string{}
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
				if h.subjectKind != "" {
					page.SectionDocs[dest.Fragment] = subjectDocument(s.file, s.source, h)
				}
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
		if s.attributes != nil {
			for _, id := range model.Strings(s.r["quality_attributes"]) {
				h := p.nodes[nameNode(p.elements[id])]
				requireRoute(page, relative(s.file, p.targets[id]), attr(h, "decoded_text"))
			}
		}
		for _, d := range s.directories {
			for _, e := range d.entries {
				dest := relative(s.file, e.target)
				page.DirectoryMembers[d.anchor] = append(page.DirectoryMembers[d.anchor], dest)
				requireRoute(page, dest, e.title)
			}
		}
		if s != p.governing {
			parent := s.parent
			if parent == nil {
				parent = s.system
			}
			requireRoute(page, relative(s.file, parent.file), parent.name.title)
		}
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
	if p.plan.Pages[reviewPath] == nil {
		return fmt.Errorf("%s: review report is missing", reviewPath)
	}
	p.plan.Pages[reviewPath].Document = Document{Type: RequirementsReviewReport}
	requireRoute(p.plan.Pages[p.governing.file], relative(p.governing.file, reviewPath), "")
	requireRoute(p.plan.Pages["SPEC.md"], relative("SPEC.md", p.governing.file), "Read the system specification")
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
