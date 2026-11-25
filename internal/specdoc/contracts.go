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
	for _, page := range p.plan.Pages {
		page.SectionDocs = map[string]Document{}
		page.RequiredSectionContent = map[string][]string{}
	}
	for _, s := range p.subjects {
		page := p.plan.Pages[s.file]
		if page == nil {
			return fmt.Errorf("%s: no generated destination for %s", s.source.path, s.name.title)
		}
		page.Document = subjectDocument(s.file, s.source, s.name)
		page.Title = s.name.title
		requireRoute(page, page.Documents, "View source")
		page.RequiredContent = nil
		for _, section := range s.overview {
			if section.level == 2 {
				requireRoute(page, relativeURL(s.file, s.file, section.anchor), section.title)
			}
		}
		for _, h := range s.source.headings {
			if h.owner != s || h.inGlossary {
				continue
			}
			file, fragment, _ := strings.Cut(p.targets[h.node], "#")
			if h == h.detailEntry {
				requireRoute(page, relative(s.file, p.targets[h.node]), inlineText(h.title))
			}
			passagePage := p.plan.Pages[file]
			if passagePage == nil {
				return fmt.Errorf("%s: missing passage destination %s", s.source.path, file)
			}
			if fragment != "" {
				passagePage.RequiredAnchors = append(passagePage.RequiredAnchors, fragment)
				for _, r := range p.routes[h.node] {
					if r.label == "Referenced by" {
						requireSectionRoute(passagePage, fragment, relative(file, r.target))
					}
				}
				if h.subjectKind != "" {
					doc := subjectDocument(file, s.source, h)
					passagePage.SectionDocs[fragment] = doc
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
				link := p.documentLink(file, s.source.path, raw)
				switch attr(n, "node_kind") {
				case "link":
					requireRoute(passagePage, link, "")
				case "image":
					passagePage.RequiredResources = append(passagePage.RequiredResources, link)
				}
			}
			if h.end < h.bodyEnd {
				prose := sourceProse(s.source, h)

				if fragment != "" {
					passagePage.RequiredSectionContent[fragment] = []string{prose}
				}
			}
		}
		// Membership and section ownership are checked against the independent
		// source inventory above. Page sequence follows its reader-facing outline.
		for _, h := range readingOrder(s) {
			if h.owner != s || h.end >= h.bodyEnd {
				continue
			}
			file, _, _ := strings.Cut(p.targets[h.node], "#")
			destination := p.plan.Pages[file]
			destination.RequiredContent = append(destination.RequiredContent, sourceProse(s.source, h))
		}
		for title, idx := range s.indexes {
			home, _, _ := strings.Cut(s.sectionIndex(title), "#")
			indexPage := p.plan.Pages[home]
			if indexPage == nil {
				return fmt.Errorf("%s: missing collection index %s", s.source.path, idx.file)
			}
			if s.browseCollection(title, idx) || title == "Quality models" {
				requireRoute(page, relative(s.file, idx.file), collectionTitle(title))
			}
			if idx.overviewSections == nil {
				requireRoute(indexPage, relative(home, s.file), s.name.title)
				indexPage.DirectoryMembers = map[string][]string{}
			}
			for group, members := range idx.groups {
				anchor := specfile.NameSlug(group)
				if group == "" {
					anchor = specfile.NameSlug(collectionTitle(title))
				}
				for _, e := range members {
					dest := relative(home, e.target)
					if section := idx.overviewSections[group]; section != nil {
						if specific := s.conceptSections[e.target]; specific != nil {
							section = specific
						}
						requireSectionRoute(indexPage, section.anchor, dest)
					} else {
						indexPage.DirectoryMembers[anchor] = append(indexPage.DirectoryMembers[anchor], dest)
					}
					requireRoute(indexPage, dest, inlineText(e.title))
				}
			}
			// Full passages on an index introduce their file at its first
			// contribution. Linked directory previews lead to subject reports.
			files := map[string]bool{}
			requireSource := func(anchor, dest, label string) {
				key := sourceFileKey(home, dest)
				if !files[key] {
					requireRoute(indexPage, dest, label)
					requireSectionRoute(indexPage, anchor, dest)
					files[key] = true
				}
			}
			for _, passage := range idx.passages {
				h := passage[0]
				_, anchor, _ := strings.Cut(p.targets[h.node], "#")
				requireSource(anchor, subjectDocument(home, s.source, h).Documents, "View section source")
			}
			if title == "Glossary" {
				for _, term := range s.glossary {
					for i, h := range term.headings {
						if h == term.topicsHeading {
							continue
						}
						anchor := term.anchors[i]
						if h.subjectKind != "" {
							indexPage.SectionDocs[anchor] = subjectDocument(idx.file, term.source, h)
						}
						edits := term.source.footnoteEdits(h.end, h.bodyEnd, term.notePrefix)
						body := sourceProse(term.source, h, edits...)
						if i == 0 {
							definition, ref := glossaryDefinition(body)
							indexPage.RequiredSectionContent[anchor] = []string{definition, ref}
						} else {
							indexPage.RequiredSectionContent[anchor] = []string{body}
						}
					}
					if term.topicsHeading != nil {
						anchor := term.anchors[0]
						topics := strings.Join(model.Strings(term.record["topics"]), ", ")
						indexPage.RequiredSectionContent[anchor] = append(indexPage.RequiredSectionContent[anchor], inlineText(topics))
					}
					dest := relativeURL(idx.file, term.source.path, term.headings[0].anchor)
					requireSource(term.anchors[0], dest, "Term source")
				}
			}
			if title == "Quality attributes" {
				for _, id := range model.Strings(s.r["quality_attributes"]) {
					h := p.nodes[nameNode(p.elements[id])]
					dest := relative(home, p.targets[id])
					requireRoute(indexPage, dest, inlineText(attr(h, "decoded_text")))
					if section := idx.overviewSections[""]; section != nil {
						requireSectionRoute(indexPage, section.anchor, dest)
					}
					requireRoute(page, relative(s.file, p.targets[id]), inlineText(attr(h, "decoded_text")))
				}
			}
		}
		if s != p.governing {
			parent := s.navigationParent
			if parent == nil {
				parent = s.system
			}
			requireRoute(page, relative(s.file, parent.file), parent.name.title)
		}
	}
	// Read typed relation facts independently of the renderer's choice to
	// present a relation inline, beside a section or in a related-items list.
	for _, a := range p.snapshot.Relationships() {
		projection := attr(a, "projection_path")
		_, field, _ := strings.Cut(projection, ".")
		if _, ok := relations[field]; !ok || projection == "QualityModel.quality_attributes" || projection == "System.quality_attributes" {
			continue
		}
		left, right := attr(a, "subject"), attr(a, "object")
		for i, edge := range [][2]string{{left, right}, {right, left}} {
			from := p.ownerPage(p.elements[edge[0]])
			source, destination := p.targets[edge[0]], p.targets[edge[1]]
			if from == nil || source == "" || destination == "" || destination == source || p.nodes[nameNode(p.elements[edge[1]])] == nil {
				continue
			}
			file, fragment, _ := strings.Cut(source, "#")
			if destination == file {
				continue
			}
			if i == 1 && source == from.file {
				if idx := p.relationshipIndex(from, relations[field].reverse, destination); idx != nil {
					file = idx.file
					if section := idx.overviewSections[""]; section != nil {
						file, fragment = from.file, section.anchor
					}
				}
			}
			page := p.plan.Pages[file]
			dest := relative(file, destination)
			requireRoute(page, dest, "")
			if fragment != "" {
				requireSectionRoute(page, fragment, dest)
			}
		}
	}
	if p.governing.file != "SPEC.md" {
		requireRoute(p.plan.Pages["SPEC.md"], relative("SPEC.md", p.governing.file), "Read the system specification")
	} else {
		for title, idx := range p.governing.indexes {
			requireRoute(p.plan.Pages["SPEC.md"], relative("SPEC.md", idx.file), collectionTitle(title))
		}
	}
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
