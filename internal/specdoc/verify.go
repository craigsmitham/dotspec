// Package specdoc generates and checks the reader-facing requirements specification.
package specdoc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/diagnostic"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.yaml.in/yaml/v3"
)

type pageHeader struct {
	Source      yaml.Node `yaml:"source"`
	Type        string    `yaml:"type"`
	GeneratedBy string    `yaml:"generated_by"`
}

func splitHeader(data []byte) (pageHeader, []byte) {
	var header pageHeader
	data = specfile.Normalize(data)
	if bytes.HasPrefix(data, []byte("---\n")) {
		if end := bytes.Index(data[4:], []byte("\n---\n")); end >= 0 {
			if yaml.Unmarshal(data[4:4+end], &header) == nil {
				return header, data[end+9:]
			}
		}
	}
	return pageHeader{}, data
}

var markdownParser = goldmark.New(goldmark.WithExtensions(extension.GFM), goldmark.WithParserOptions(parser.WithAutoHeadingID())).Parser()

type sourceRoute struct{ label, dest string }

// Relative destinations on the same report have the same base directory.
// URL parsing decodes the path and discards query and fragment for file identity.
func sourceFileKey(page, dest string) string {
	u, err := url.Parse(dest)
	if err != nil || u.IsAbs() || u.Host != "" || u.Path == "" {
		return ""
	}
	return path.Join(path.Dir(page), u.Path)
}

func isSourceLabel(label string) bool {
	return label == "View source" || label == "View section source" || label == "Term source"
}

func sourceLabelAllowances(markdown string, generated []sourceRoute) map[sourceRoute]int {
	_, body := splitHeader([]byte(markdown))
	counts := map[sourceRoute]int{}
	doc := markdownParser.Parse(text.NewReader(body))
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := n.(*ast.Link); ok && entering {
			if label := visibleText(l, body); isSourceLabel(label) {
				counts[sourceRoute{label, string(util.URLEscape(l.Destination, true))}]++
			}
		}
		return ast.WalkContinue, nil
	})
	for _, r := range generated {
		counts[sourceRoute{r.label, string(util.URLEscape([]byte(r.dest), true))}]--
	}
	return counts
}

type parsedPage struct {
	footerSource   string
	sectionDocs    map[string][]Document
	sectionLists   map[string]int
	firstList      map[string][]directoryRow
	header         pageHeader
	titles         []string
	visible        string
	links          map[string]int
	labels         map[string][]string
	sectionLinks   map[string]map[string]bool
	sectionContent map[string][]string
	resources      map[string]bool
	anchors        map[string]bool
}

func parseGenerated(data []byte) parsedPage {
	header, body := splitHeader(data)
	doc := markdownParser.Parse(text.NewReader(body))
	p := parsedPage{sectionDocs: map[string][]Document{}, sectionLists: map[string]int{}, firstList: map[string][]directoryRow{}, sectionContent: map[string][]string{}, sectionLinks: map[string]map[string]bool{}, header: header, links: make(map[string]int), labels: make(map[string][]string), resources: make(map[string]bool), anchors: make(map[string]bool), visible: visibleText(doc, body)}
	type section struct {
		level int
		ids   []string
	}
	var overviewIDs specfile.RenderedHeadingIDs
	if header.Type == "system-spec-report" || header.Type == "component-spec-report" {
		overviewIDs = specfile.RenderedHeadingIDs{}
	}
	var ancestors []section
	var aliases []string
	addAnchor := func(id string) { p.anchors[id] = true; aliases = append(aliases, id) }
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if n.Parent() == doc && n.Kind() != ast.KindHeading && len(ancestors) > 0 {
			bodyText := visibleText(n, body)
			for _, id := range ancestors[len(ancestors)-1].ids {
				p.sectionContent[id] = append(p.sectionContent[id], bodyText)
			}
		}
		// The directory is the first direct list in its section. Later authored
		// prose and links can refer to other subjects for different reasons.
		if list, ok := n.(*ast.List); ok && n.Parent() == doc && len(ancestors) > 0 {
			for _, ancestor := range ancestors {
				for _, id := range ancestor.ids {
					p.sectionLists[id]++
				}
			}
			ids := ancestors[len(ancestors)-1].ids
			if _, seen := p.firstList[ids[0]]; !seen {
				rows := parseDirectoryRows(list, body)
				for _, id := range ids {
					p.firstList[id] = rows
				}
			}
		}
		if anchor := specfile.HTMLAnchor(n, body); anchor != "" {
			addAnchor(anchor)
		}
		switch n := n.(type) {
		case *ast.Heading:
			if attr, ok := n.AttributeString("id"); ok {
				name := string(attr.([]byte))
				if overviewIDs != nil {
					name = overviewIDs.Next(specfile.RenderedText(n, body))
				}
				p.anchors[name] = true
				if n.Parent() == doc {
					for len(ancestors) > 0 && ancestors[len(ancestors)-1].level >= n.Level {
						ancestors = ancestors[:len(ancestors)-1]
					}
					ids := append(aliases, name)
					if n.Level == 1 {
						ids = append(ids, "") // the title's introductory section
					}
					ancestors = append(ancestors, section{n.Level, ids})
					aliases = nil
				}
			}
			if n.Level == 1 {
				p.titles = append(p.titles, visibleText(n, body))
			}
		case *ast.Link:
			if label := visibleText(n, body); label != "" {
				dest := string(util.URLEscape(n.Destination, true))
				if label == "View source" && n.Parent() == doc.LastChild() && n == n.Parent().FirstChild() &&
					n.Parent().Kind() == ast.KindParagraph && n.Parent().PreviousSibling() != nil && n.Parent().PreviousSibling().Kind() == ast.KindThematicBreak {
					p.footerSource = dest
				}
				p.links[dest]++
				p.labels[dest] = append(p.labels[dest], label)
				for _, section := range ancestors {
					for _, id := range section.ids {
						if p.sectionLinks[id] == nil {
							p.sectionLinks[id] = map[string]bool{}
						}
						p.sectionLinks[id][dest] = true
					}
				}
			}
		case *ast.Image:
			p.resources[string(util.URLEscape(n.Destination, true))] = true
		case *ast.HTMLBlock:
			raw := strings.TrimSpace(string(n.Lines().Value(body)))
			if data, ok := strings.CutPrefix(raw, "<!-- dotspec-doc "); ok && n.Parent() == doc {
				anchor := ""
				if len(ancestors) > 0 {
					anchor = ancestors[len(ancestors)-1].ids[0]
				}
				var document Document
				_ = json.Unmarshal([]byte(strings.TrimSuffix(data, " -->")), &document)
				p.sectionDocs[anchor] = append(p.sectionDocs[anchor], document)
			}
		}
		return ast.WalkContinue, nil
	})
	return p
}

// Read visible text from the AST. A link label survives destination rebasing,
// and processor comments and frontmatter cannot satisfy a visible-content rule.
func visibleText(node ast.Node, source []byte) string {
	return strings.Join(strings.Fields(specfile.RenderedText(node, source)), " ")
}

func contentText(md string) string {
	source := []byte(md)
	return visibleText(markdownParser.Parse(text.NewReader(source)), source)
}

// A directory item is a bullet starting with one title link. Any following
// preview is plain text, optionally with an italic note or label. A nested list holds the item's own members, one level
// deeper. Keep malformed items in the count for diagnostics.
func parseDirectoryRows(list *ast.List, body []byte) []directoryRow {
	return appendDirectoryRows(nil, list, body, 0)
}

func appendDirectoryRows(rows []directoryRow, list *ast.List, body []byte, depth int) []directoryRow {
	if list.IsOrdered() {
		return rows
	}
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		row := directoryRow{text: visibleText(item, body), depth: depth}
		block := item.FirstChild()
		var nested *ast.List
		if block != nil && block.NextSibling() != nil && block.NextSibling().NextSibling() == nil {
			nested, _ = block.NextSibling().(*ast.List)
		}
		if nested != nil {
			row.text = visibleText(block, body)
		}
		if block != nil && (block.NextSibling() == nil || nested != nil) {
			if link, ok := block.FirstChild().(*ast.Link); ok {
				row.destination = string(util.URLEscape(link.Destination, true))
				row.title = visibleText(link, body)
				for sibling := link.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
					if sibling.Kind() != ast.KindText && sibling.Kind() != ast.KindString && sibling.Kind() != ast.KindEmphasis {
						row.destination = ""
						break
					}
				}
			}
		}
		rows = append(rows, row)
		if nested != nil {
			rows = appendDirectoryRows(rows, nested, body, depth+1)
		}
	}
	return rows
}

func checkSourceLinks(file string, expected *Page, p parsedPage, report func(string, string)) {
	if expected.Documents != "" && p.footerSource != string(util.URLEscape([]byte(expected.Documents), true)) {
		report("spec.navigation", "expected View source in the footer to open "+expected.Documents)
	}
	ordinary := maps.Clone(expected.otherSourceLabels)
	sources := map[string]int{}
	for dest, labels := range p.labels {
		for _, label := range labels {
			if !isSourceLabel(label) {
				continue
			}
			route := sourceRoute{label, dest}
			if ordinary[route] > 0 {
				ordinary[route]--
				continue
			}
			if key := sourceFileKey(file, dest); key != "" {
				sources[key]++
			}
		}
	}
	required := map[string]bool{}
	if expected.Documents != "" {
		required[sourceFileKey(file, expected.Documents)] = true
	}
	for _, declaration := range expected.SectionDocs {
		required[sourceFileKey(file, declaration.Documents)] = true
	}
	for _, key := range slices.Sorted(maps.Keys(required)) {
		if sources[key] == 0 {
			report("spec.navigation", "missing generated source link to "+key)
		}
	}
	for _, key := range slices.Sorted(maps.Keys(sources)) {
		if sources[key] > 1 {
			report("spec.navigation", "multiple generated source links to "+key)
		} else if !required[key] {
			report("spec.navigation", "unexpected generated source link to "+key)
		}
	}
}

// Check the rendered document against the plan's source-derived contract.
func checkPage(file string, expected *Page, p parsedPage, report func(string, string)) {
	if p.header.GeneratedBy != "dotspec" || p.header.Type != expected.Type {
		report("spec.document-type", fmt.Sprintf("expected a generated %s document for %q", expected.Type, expected.Title))
	}
	if expected.Documents != "" && (p.header.Source.ShortTag() != "!!str" || p.header.Source.Value != expected.Documents) {
		report("spec.document-type", fmt.Sprintf("expected authored subject %q for %q", expected.Documents, expected.Title))
	}
	if len(p.titles) != 1 || p.titles[0] != contentText(expected.Title) {
		report("spec.title", "expected one title naming "+expected.Title)
	}
	checkSourceLinks(file, expected, p, report)
	remaining := p.visible
	for _, md := range expected.RequiredContent {
		needle := contentText(md)
		if needle == "" {
			continue
		}
		if _, rest, ok := strings.Cut(remaining, needle); ok {
			remaining = rest
			continue
		}
		excerpt := []rune(needle)
		if len(excerpt) > 100 {
			excerpt = append(excerpt[:100], '…')
		}
		report("spec.content", fmt.Sprintf("missing, altered, or out-of-order authored content: %q", string(excerpt)))
	}
	for _, anchor := range slices.Sorted(maps.Keys(expected.SectionDocs)) {
		if !slices.Equal(p.sectionDocs[anchor], []Document{expected.SectionDocs[anchor]}) {
			report("spec.document-type", "expected one typed document for the authored subject at #"+anchor)
		}
	}
	for _, anchor := range slices.Sorted(maps.Keys(p.sectionDocs)) {
		if _, ok := expected.SectionDocs[anchor]; !ok {
			report("spec.document-type", "unexpected subject document at #"+anchor)
		}
	}
	for _, link := range expected.RequiredLinks {
		if p.links[string(util.URLEscape([]byte(link), true))] == 0 {
			report("spec.navigation", fmt.Sprintf("missing visible link to %q", link))
		}
	}
	for _, link := range slices.Sorted(maps.Keys(expected.RequiredLabels)) {
		labels := p.labels[string(util.URLEscape([]byte(link), true))]
		for _, label := range expected.RequiredLabels[link] {
			if !slices.Contains(labels, contentText(label)) {
				report("spec.navigation", fmt.Sprintf("missing link labeled %q to %q", label, link))
			}
		}
	}
	for _, anchor := range slices.Sorted(maps.Keys(expected.RequiredSectionContent)) {
		body := strings.Join(p.sectionContent[anchor], " ")
		for _, part := range expected.RequiredSectionContent[anchor] {
			needle := contentText(part)
			if needle != "" && !strings.Contains(body, needle) {
				report("spec.content", fmt.Sprintf("source content missing or altered in section %q", anchor))
				break
			}
		}
	}
	for _, resource := range expected.RequiredResources {
		if !p.resources[string(util.URLEscape([]byte(resource), true))] {
			report("spec.content", fmt.Sprintf("missing image %q", resource))
		}
	}
	for _, anchor := range slices.Sorted(maps.Keys(expected.RequiredSectionLinks)) {
		for _, dest := range expected.RequiredSectionLinks[anchor] {
			if !p.sectionLinks[anchor][string(util.URLEscape([]byte(dest), true))] {
				report("spec.navigation", fmt.Sprintf("section %q is missing link to %q", anchor, dest))
			}
		}
	}
	for _, anchor := range expected.RequiredAnchors {
		if !p.anchors[anchor] {
			report("spec.structure", fmt.Sprintf("missing addressable section %q", anchor))
		}
	}
	for _, anchor := range slices.Sorted(maps.Keys(expected.DirectoryMembers)) {
		directory := expected.DirectoryMembers[anchor]
		members := slices.Clone(directory.rows)
		for i := range members {
			row := &members[i]
			row.destination = string(util.URLEscape([]byte(row.destination), true))
			row.text = strings.Join(strings.Fields(row.text), " ")
		}
		if !slices.Equal(members, p.firstList[anchor]) || (directory.wholeSection && p.sectionLists[anchor] != 1) {
			report("spec.membership", "directory at #"+anchor+" must list each scoped title and preview once in alphabetical order")
		}
	}
}

// validateRendered compares source-derived page contracts against the visible
// files, independently of the formatters that generated their Markdown.
func validateRendered(root string, plan *Plan, files map[string][]byte) ([]diagnostic.Diagnostic, error) {
	repo, err := specfile.RepositoryRoot(root)
	if err != nil {
		return nil, err
	}
	pages := make(map[string]parsedPage, len(files))
	for rel, data := range files {
		pages[rel] = parseGenerated(data)
	}
	cache := map[string]*specfile.ReferenceDocument{}
	edges := map[string][]string{}
	var ds []diagnostic.Diagnostic
	claims := map[string]string{}
	for _, rel := range slices.Sorted(maps.Keys(pages)) {
		p := pages[rel]
		if p.header.GeneratedBy == "dotspec" && p.header.Source.ShortTag() == "!!str" && p.header.Source.Value != "" {
			source, err := url.Parse(p.header.Source.Value)
			if err == nil && !source.IsAbs() {
				key := (&url.URL{Path: path.Join(path.Dir(rel), source.Path), Fragment: source.Fragment}).String()
				if previous, ok := claims[key]; ok {
					ds = append(ds, *diagnostic.Error("spec.duplicate-report", "two report pages document "+key,
						diagnostic.Location{Path: filepath.Join(root, rel)}, diagnostic.Location{Path: filepath.Join(root, previous)}))
				} else {
					claims[key] = rel
				}
			}
		}
		expected := plan.Pages[rel]
		if expected == nil {
			continue // Obsolete pages can duplicate a claim but have no current content contract.
		}
		report := func(rule, message string) {
			ds = append(ds, *diagnostic.Error(rule, message, diagnostic.Location{Path: filepath.Join(root, filepath.FromSlash(rel)), Element: expected.Title}))
		}
		checkPage(rel, expected, p, report)
		destinations := maps.Clone(p.resources)
		if p.header.Source.ShortTag() == "!!str" {
			destinations[p.header.Source.Value] = true
		}
		for dest := range p.links {
			destinations[dest] = true
		}
		for _, dest := range slices.Sorted(maps.Keys(destinations)) {
			u, err := url.Parse(dest)
			if err != nil {
				report("spec.link", "invalid URL "+dest)
				continue
			}
			if u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") {
				continue
			}
			target := rel
			if u.Path != "" {
				target = path.Join(path.Dir(rel), u.Path)
			}
			if _, generated := plan.Pages[target]; generated {
				t, ok := pages[target]
				if !ok {
					report("spec.link", "missing generated destination "+dest)
					continue
				}
				if u.Fragment != "" && !t.anchors[u.Fragment] {
					report("spec.link", "missing generated anchor "+dest)
					continue
				}
				if p.links[dest] > 0 {
					edges[rel] = append(edges[rel], target)
				}
				continue
			}
			targetPath := filepath.Join(root, filepath.FromSlash(target))
			info, allowed, err := specfile.ReferenceTarget(targetPath, repo)
			if !allowed || errors.Is(err, fs.ErrNotExist) {
				report("spec.link", fmt.Sprintf("unavailable local destination %q", dest))
				continue
			}
			if err != nil {
				return ds, err
			}
			if u.Fragment == "" {
				continue
			}
			if info.IsDir() {
				report("spec.link", "directory has no anchor "+dest)
				continue
			}
			if !strings.EqualFold(path.Ext(target), ".md") && !strings.EqualFold(path.Ext(target), ".markdown") {
				continue
			}
			ref := cache[targetPath]
			if ref == nil {
				data, err := os.ReadFile(targetPath)
				if err != nil {
					return ds, err
				}
				ref = specfile.MarkdownReferences(data)
				cache[targetPath] = ref
			}
			if !ref.Anchors[u.Fragment] {
				report("spec.link", "missing source heading "+dest)
			}
		}
	}
	// Missing pages already have actionable findings. Verify that every existing
	// page can be reached by a reader starting at the root document.
	if _, ok := pages["SPEC.md"]; ok {
		seen := map[string]bool{"SPEC.md": true}
		queue := []string{"SPEC.md"}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			for _, next := range edges[current] {
				if !seen[next] {
					seen[next] = true
					queue = append(queue, next)
				}
			}
		}
		for _, rel := range slices.Sorted(maps.Keys(pages)) {
			if plan.Pages[rel] != nil && !seen[rel] {
				ds = append(ds, *diagnostic.Error("spec.unreachable", "no visible route from SPEC.md", diagnostic.Location{Path: filepath.Join(root, filepath.FromSlash(rel))}))
			}
		}
	}
	return ds, nil
}
