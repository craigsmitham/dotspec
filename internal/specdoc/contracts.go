package specdoc

import (
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
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
		page.Title = s.name.title
		page.RequiredContent = nil
		page.RequiredSectionContent = map[string]string{}
		for _, section := range s.overview {
			if section.level == 2 && section.hasBody {
				requireRoute(page, relativeURL(s.file, s.file, section.anchor), section.title)
			}
		}
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
		for title, idx := range s.indexes {
			indexPage := p.plan.Pages[idx.file]
			if indexPage == nil {
				return fmt.Errorf("%s: missing collection index %s", s.source.path, idx.file)
			}
			label := title
			if s.kind() == "System" {
				label = overviewIndexLabel(title)
			}
			requireRoute(page, relative(s.file, idx.file), label)
			requireRoute(indexPage, relative(idx.file, s.file), s.name.title)
			indexPage.DirectoryMembers = map[string][]string{}
			for group, headings := range idx.sources {
				for _, h := range headings {
					dest := relativeURL(idx.file, s.source.path, h.anchor)
					requireRoute(indexPage, dest, h.sourceLabel())
					if group != "" {
						requireSectionRoute(indexPage, specfile.NameSlug(group), dest)
					}
				}
			}
			for group, members := range idx.groups {
				anchor := specfile.NameSlug(group)
				if group == "" {
					anchor = specfile.NameSlug(title)
				}
				for _, e := range members {
					dest := relative(idx.file, e.target)
					indexPage.DirectoryMembers[anchor] = append(indexPage.DirectoryMembers[anchor], dest)
					requireRoute(indexPage, dest, inlineText(e.title))
				}
			}
			if title == "Glossary" {
				indexPage.SectionDocs = map[string]Document{}
				indexPage.RequiredSectionContent = map[string]string{}
				for _, term := range s.glossary {
					for i, h := range term.headings {
						anchor := term.anchors[i]
						if h.subjectKind != "" {
							indexPage.SectionDocs[anchor] = subjectDocument(idx.file, term.source, h)
						}
						edits := term.source.footnoteEdits(h.end, h.bodyEnd, term.notePrefix)
						indexPage.RequiredSectionContent[anchor] = sourceProse(term.source, h, edits...)
					}
					dest := relativeURL(idx.file, term.source.path, term.headings[0].anchor)
					requireRoute(indexPage, dest, "View term source")
					requireSectionRoute(indexPage, term.anchors[0], dest)
				}
			}
			if title == "Quality attributes" {
				for _, id := range model.Strings(s.r["quality_attributes"]) {
					h := p.nodes[nameNode(p.elements[id])]
					requireRoute(indexPage, relative(idx.file, p.targets[id]), attr(h, "decoded_text"))
					requireRoute(page, relative(s.file, p.targets[id]), attr(h, "decoded_text"))
				}
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
		for i, edge := range [][2]string{{left, right}, {right, left}} {
			from := p.ownerPage(p.elements[edge[0]])
			source, destination := p.targets[edge[0]], p.targets[edge[1]]
			if from == nil || source == "" || destination == "" || destination == from.file || destination == source || p.nodes[nameNode(p.elements[edge[1]])] == nil {
				continue
			}
			file := from.file
			if i == 1 && source == file {
				if idx := relationshipIndex(from, relations[field].reverse, destination); idx != nil {
					file = idx.file
				}
			}
			page := p.plan.Pages[file]
			dest := relative(file, destination)
			requireRoute(page, dest, "")
			if _, fragment, ok := strings.Cut(source, "#"); ok {
				requireSectionRoute(page, fragment, dest)
			}
		}
	}
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
	key := sourceKey(target, u.Fragment)
	canonical := p.glossaryTargets[pagePath][key]
	if canonical == "" {
		canonical = p.targets[key]
	}
	if canonical != "" {
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
