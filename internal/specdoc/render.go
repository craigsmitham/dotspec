package specdoc

import (
	"cmp"
	"fmt"
	"html"
	"maps"
	"net/url"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
)

const contentsPath = "specification/index.md"
const reviewPath = "specification/review.md"

type writer struct {
	page                *Page
	file, titleMarkdown string
	out                 strings.Builder
}

func (p *planner) newPage(file, typ, subjectKind, title string) *writer {
	page := &Page{Type: typ, SubjectKind: subjectKind, Title: title}
	p.plan.Pages[file] = page
	return &writer{page: page, file: file, titleMarkdown: title}
}
func (w *writer) finish() {
	var header strings.Builder
	fmt.Fprintf(&header, "---\ntype: %s\ngenerated_by: dotspec\n", w.page.Type)
	if w.page.SubjectKind != "" {
		fmt.Fprintf(&header, "subject_kind: %s\n", w.page.SubjectKind)
	}
	fmt.Fprintf(&header, "---\n# %s\n", w.titleMarkdown)
	w.page.Markdown = strings.TrimSpace(header.String()+w.out.String()) + "\n"
}
func (w *writer) text(value string) {
	if value != "" {
		w.out.WriteString("\n" + strings.TrimSpace(value) + "\n")
	}
}
func (w *writer) heading(value string) { w.text("## " + value) }
func (w *writer) link(title, target string) string {
	dest := relative(w.file, target)
	w.requireLink(dest)
	return link(title, dest)
}
func (w *writer) sourceLink(label, file, fragment string) string {
	dest := relativeURL(w.file, file, fragment)
	w.requireLink(dest)
	return link(label, dest)
}
func (w *writer) requireLink(dest string) {
	if !slices.Contains(w.page.RequiredLinks, dest) {
		w.page.RequiredLinks = append(w.page.RequiredLinks, dest)
	}
}
func link(title, dest string) string {
	title = strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]").Replace(title)
	return "[" + title + "](<" + dest + ">)"
}
func relative(from, to string) string {
	file, frag, _ := strings.Cut(to, "#")
	return relativeURL(from, file, frag)
}
func relativeURL(from, file, frag string) string {
	rel, _ := filepath.Rel(filepath.FromSlash(path.Dir(from)), filepath.FromSlash(file))
	u := url.URL{Path: filepath.ToSlash(rel)}
	if frag != "" {
		u.Fragment = frag
		if from == file {
			u.Path = ""
		}
	}
	return u.String()
}
func (p *planner) render() *Plan {
	p.indexRelationships()
	for _, s := range p.subjects {
		p.renderSubject(s)
	}
	for _, file := range slices.Sorted(maps.Keys(p.collections)) {
		p.renderIndex(p.collections[file])
	}
	p.renderContents()
	p.renderReview()
	return p.plan
}

func (p *planner) renderSubject(s *subject) {
	typ := SpecificationSection
	if s.kind() == "System" {
		typ = RequirementsSpecification
	}
	w := p.newPage(s.file, typ, kindName(s.kind()), s.name.title)
	w.titleMarkdown = p.headingMarkdown(w, s, s.name)
	if s.kind() != "System" {
		w.text(w.link(s.system.name.title, s.system.file) + " · " + w.link("Contents", contentsPath))
	}
	navigationShown := false
	for _, h := range readingOrder(s) {
		if s.kind() == "System" && !navigationShown && h != s.name && h.title != "Purpose" && h.title != "Scope" {
			var nav []string
			if s != p.governing {
				nav = append(nav, w.link(s.system.name.title, s.system.file))
			}
			nav = append(nav, w.link("Contents", contentsPath), w.link("Review", reviewPath))
			w.text(strings.Join(nav, " · "))
			navigationShown = true
		}
		if h.owner != s {
			if h.owner != nil && h.owner.name == h && p.ownerPage(p.elements[attr(h.owner.r, "owner")]) == s {
				w.text("- " + w.link(h.title, h.owner.file))
			}
			continue
		}
		if h != s.name {
			// Preserve addressability for duplicated headings moved onto different pages.
			if h.explicitAnchor != "" {
				w.text(`<a id="` + html.EscapeString(h.explicitAnchor) + `"></a>`)
			}
			w.text(strings.Repeat("#", max(2, h.level-s.name.level+1)) + " " + p.headingMarkdown(w, s, h))
		}
		p.body(w, s.source, h)
		if h != s.name {
			p.embeddedRelationships(w, h)
			if h.obligation {
				w.text(w.sourceLink("Edit obligation", s.source.path, h.anchor))
			}
		}
	}
	if s.kind() == "System" {
		w.heading("Explore this system")
		w.out.WriteByte('\n')
		for _, file := range slices.Sorted(maps.Keys(p.collections)) {
			c := p.collections[file]
			if c.system == s {
				fmt.Fprintln(&w.out, "- "+w.link(sentenceCase(c.label), c.file))
			}
		}
		fmt.Fprintln(&w.out, "- "+w.link("Full specification contents", contentsPath))
		fmt.Fprintln(&w.out, "- "+w.link("Requirements review", reviewPath))
	}
	p.relationships(w, s)
	p.appendFootnotes(w, s)
	w.text(w.sourceLink("Edit source", s.source.path, s.name.anchor))
	w.finish()
}

// readingOrder promotes a few entire authored sections for the document's
// purpose. Descendants, scenario steps and all other siblings retain order.
func readingOrder(s *subject) []*heading {
	var groups [][]*heading
	for _, h := range s.source.headings {
		if h.start < s.name.start {
			continue
		}
		if h != s.name && h.level <= s.name.level {
			break
		}
		if h == s.name || h.level == s.name.level+1 || len(groups) == 0 {
			groups = append(groups, []*heading{h})
		} else {
			i := len(groups) - 1
			groups[i] = append(groups[i], h)
		}
	}
	priority := func(group []*heading) int {
		if group[0] == s.name {
			return 0
		}
		if s.kind() == "Requirement" && group[0].title == "Statement" {
			return 1
		}
		if s.kind() == "System" {
			switch group[0].title {
			case "Purpose":
				return 1
			case "Scope":
				return 2
			case "Features":
				return 3
			}
		}
		return 4
	}
	slices.SortStableFunc(groups, func(a, b []*heading) int { return cmp.Compare(priority(a), priority(b)) })
	var headings []*heading
	for _, group := range groups {
		headings = append(headings, group...)
	}
	return headings
}

// A footnote can be defined anywhere in its authored file. Keep its
// definition available locally on every generated page that uses it.
func (p *planner) appendFootnotes(w *writer, s *subject) {
	definitions := map[string]model.Record{}
	needed := map[string]bool{}
	owner := func(pos int) *subject {
		for _, h := range s.source.headings {
			if h.start <= pos && pos < h.bodyEnd {
				return h.owner
			}
		}
		return nil
	}
	for _, n := range s.source.footnotes {
		key := strings.ToLower(attr(n, "label"))
		if attr(n, "node_kind") == "footnote_definition" {
			definitions[key] = n
			continue
		}
		if owner(position(n, "start")) == s {
			needed[key] = true
		}
	}
	for _, key := range slices.Sorted(maps.Keys(needed)) {
		if n := definitions[key]; n != nil && owner(position(n, "start")) != s {
			w.text(p.rebasedText(w, s.source, position(n, "start"), position(n, "end")))
		}
	}
}
func (p *planner) headingMarkdown(w *writer, s *subject, h *heading) string {
	raw := p.rebasedText(w, s.source, h.start, h.end)
	if strings.HasPrefix(raw, "#") {
		// CommonMark permits a closing run of hashes preceded by whitespace.
		raw = strings.TrimSpace(strings.TrimLeft(raw, "#"))
		before := strings.TrimRight(raw, "#")
		if len(before) < len(raw) && len(before) > 0 && (before[len(before)-1] == ' ' || before[len(before)-1] == '\t') {
			raw = strings.TrimSpace(before)
		}
		return raw
	}
	// A Setext heading's final line is its underline.
	if end := strings.LastIndexByte(raw, '\n'); end > 0 {
		return strings.Join(strings.Fields(raw[:end]), " ")
	}
	return h.title
}
func kindName(kind string) string {
	var b strings.Builder
	for i, r := range kind {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
		}
		b.WriteString(strings.ToLower(string(r)))
	}
	return b.String()
}
func sentenceCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (p *planner) body(w *writer, sf *sourceFile, h *heading) {
	w.text(p.rebasedText(w, sf, h.end, h.bodyEnd))
}

type replacement struct {
	start, end int
	text       string
}

// Change only link destinations. Their intervals are disjoint even when an
// image is the label of a link. All surrounding source Markdown stays intact.
func (p *planner) rebasedText(w *writer, sf *sourceFile, start, end int) string {
	var edits []replacement
	for _, n := range sf.links {
		a, b := position(n, "start"), position(n, "end")
		raw := attr(n, "destination")
		if raw == "" || a < start || b > end {
			continue
		}
		target := p.documentLink(w.file, sf.path, raw)
		if attr(n, "node_kind") == "link" {
			w.requireLink(target)
		}
		original := sf.data[a:b]
		changed := rewriteDestination(original, n, target)
		common := 0
		for common < len(original) && common < len(changed) && original[common] == changed[common] {
			common++
		}
		if common < len(original) || common < len(changed) {
			edits = append(edits, replacement{a + common, b, changed[common:]})
		}
	}
	slices.SortFunc(edits, func(a, b replacement) int { return cmp.Compare(a.start, b.start) })
	var out strings.Builder
	pos := start
	for _, edit := range edits {
		if edit.start < pos {
			continue
		}
		out.WriteString(sf.data[pos:edit.start])
		out.WriteString(edit.text)
		pos = edit.end
	}
	out.WriteString(sf.data[pos:end])
	return strings.TrimSpace(strings.ReplaceAll(out.String(), "\r\n", "\n"))
}
func sourceProse(sf *sourceFile, h *heading) string {
	var b strings.Builder
	pos := h.end
	for _, n := range sf.links {
		start, end := position(n, "start"), position(n, "end")
		if start < pos || end > h.bodyEnd || attr(n, "label") == "" || attr(n, "node_kind") == "link_definition" {
			continue
		}
		b.WriteString(sf.data[pos:start])
		b.WriteString(rewriteDestination(sf.data[start:end], n, attr(n, "destination")))
		pos = end
	}
	b.WriteString(sf.data[pos:h.bodyEnd])
	return strings.TrimSpace(b.String())
}
func rewriteDestination(raw string, n model.Record, target string) string {
	// Reference definitions may be shared between sections. Expand reference
	// uses inline so each generated page works independently.
	if attr(n, "node_kind") == "link_definition" {
		at := strings.Index(raw, "]:")
		if at < 0 {
			return raw
		}
		return raw[:at+2] + " <" + target + ">" + linkTitle(n) + "\n"
	}
	level := 0
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '\\':
			i++
		case '[':
			level++
		case ']':
			level--
			if level == 0 {
				return raw[:i+1] + "(<" + target + ">" + linkTitle(n) + ")"
			}
		}
	}
	return raw
}
func linkTitle(n model.Record) string {
	if title := attr(n, "title"); title != "" {
		return " \"" + strings.ReplaceAll(title, "\"", "\\\"") + "\""
	}
	return ""
}

type relation struct{ forward, reverse string }

var relations = map[string]relation{
	"groups": {"Groups", "Feature"}, "refines": {"Refines", "Refined by"}, "supports": {"Supports", "Supported by"},
	"applies_to": {"Applies to", "Applies here"}, "quality_attributes": {"Quality attribute", "Classifies"},
	"derived_from": {"Derived from", "Basis for"}, "relies_on": {"Relies on", "Relied on by"},
	"involves": {"Involves", "Used in"}, "illustrates": {"Illustrates", "Illustrated by"},
	"definition": {"Definition", "Used by"}, "part_of": {"Part of", "Parts"},
}

type route struct {
	label, title, target string
	forward              bool
}

func (p *planner) indexRelationships() {
	p.routes = map[string][]route{}
	for _, a := range p.snapshot.Assertions {
		projection := attr(a, "projection_path")
		_, field, _ := strings.Cut(projection, ".")
		words, ok := relations[field]
		if !ok || projection == "System.quality_attributes" || projection == "QualityModel.quality_attributes" {
			continue
		}
		left, right := p.elements[attr(a, "subject")], p.elements[attr(a, "object")]
		add := func(from, to model.Record, label string, forward bool) {
			dest := p.targets[id(to)]
			h := nameNode(from)
			title := attr(p.nodes[nameNode(to)], "decoded_text")
			if h != "" && title != "" && dest != "" {
				p.routes[h] = append(p.routes[h], route{label, title, dest, forward})
			}
		}
		add(left, right, words.forward, true)
		add(right, left, words.reverse, false)
	}
}

// Place inverse links beside an embedded obligation or outcome. An authored
// forward link anywhere within that heading's child sections already serves
// the reader, so materialize only additional routes.
func (p *planner) embeddedRelationships(w *writer, h *heading) {
	if len(p.routes[h.node]) == 0 {
		return
	}
	sf := h.owner.source
	limit := len(sf.data)
	for _, next := range sf.headings {
		if next.start > h.start && next.level <= h.level {
			limit = next.start
			break
		}
	}
	var links []string
	for _, r := range p.routes[h.node] {
		if r.target == w.file {
			continue
		}
		found := false
		if r.forward {
			dest := relative(w.file, r.target)
			for _, n := range sf.links {
				pos := position(n, "start")
				if pos >= h.start && pos < limit && attr(n, "node_kind") == "link" && p.documentLink(w.file, sf.path, attr(n, "destination")) == dest {
					found = true
					break
				}
			}
		}
		if !found {
			links = append(links, "**"+r.label+"** · "+w.link(r.title, r.target))
		}
	}
	slices.Sort(links)
	for _, line := range slices.Compact(links) {
		w.text(line)
	}
}
func (p *planner) relationships(w *writer, s *subject) {
	var lines []string
	for _, r := range p.routes[s.name.node] {
		if r.target == w.file || slices.Contains(w.page.RequiredLinks, relative(w.file, r.target)) {
			continue
		}
		lines = append(lines, r.label+" — "+w.link(r.title, r.target))
	}
	if len(lines) > 0 {
		w.heading("Related specifications")
		w.out.WriteByte('\n')
		slices.Sort(lines)
		for _, line := range slices.Compact(lines) {
			fmt.Fprintln(&w.out, "- "+line)
		}
	}
}
func (p *planner) renderIndex(c *collection) {
	title := sentenceCase(c.label) + " for " + c.system.name.title
	w := p.newPage(c.file, CollectionIndex, "", title)
	w.text(w.link(c.system.name.title, c.system.file) + " · " + w.link("Contents", contentsPath))
	slices.SortFunc(c.entries, func(a, b entry) int {
		return cmp.Or(strings.Compare(a.title, b.title), strings.Compare(a.target, b.target))
	})
	w.out.WriteByte('\n')
	for _, e := range c.entries {
		fmt.Fprintln(&w.out, "- "+w.link(e.title, e.target))
	}
	w.finish()
}
func (p *planner) renderContents() {
	w := p.newPage(contentsPath, TableOfContents, "", "Specification contents")
	systems := []*subject{p.governing}
	for _, s := range p.subjects {
		if s != p.governing && s.kind() == "System" {
			systems = append(systems, s)
		}
	}
	slices.SortFunc(systems[1:], func(a, b *subject) int { return strings.Compare(a.file, b.file) })
	w.heading("By system")
	w.out.WriteByte('\n')
	for _, s := range systems {
		fmt.Fprintln(&w.out, strings.Repeat("  ", strings.Count(s.file, "/systems/"))+"- "+w.link(s.name.title, s.file))
	}
	byKind := map[string][]*collection{}
	for _, file := range slices.Sorted(maps.Keys(p.collections)) {
		c := p.collections[file]
		byKind[c.label] = append(byKind[c.label], c)
	}
	w.heading("By kind")
	w.out.WriteByte('\n')
	for _, label := range slices.Sorted(maps.Keys(byKind)) {
		var links []string
		for _, c := range byKind[label] {
			links = append(links, w.link(sentenceCase(label)+" for "+c.system.name.title, c.file))
		}
		fmt.Fprintln(&w.out, "- "+strings.Join(links, " · "))
	}
	w.text(w.link("Requirements review", reviewPath))
	w.finish()
}
func (p *planner) renderReview() {
	w := p.newPage(reviewPath, RequirementsReviewReport, "", "Requirements review")
	w.text(w.link(p.governing.name.title, p.governing.file) + " · " + w.link("Contents", contentsPath))
	var questions, findings, verification, validation []string
	for _, s := range p.subjects {
		var ancestors []*heading
		for _, h := range s.source.headings {
			if h.owner != s {
				continue
			}
			for len(ancestors) > 0 && ancestors[len(ancestors)-1].level >= h.level {
				ancestors = ancestors[:len(ancestors)-1]
			}
			subject := s.name.title
			if len(ancestors) > 1 {
				subject = ancestors[len(ancestors)-1].title + " · " + subject
			}
			switch h.title {
			case "Open questions", "Open question":
				// Lists and prose remain legible at the review destination;
				// nested question headings also receive an individual route.
				label := contentText(sourceProse(s.source, h))
				if label != "" {
					questions = append(questions, w.link(label, p.targets[h.node])+" · "+s.name.title)
				} else if h.title == "Open question" {
					questions = append(questions, w.link(subject+" — "+h.title, p.targets[h.node]))
				}
			case "Verification plan":
				verification = append(verification, w.link(subject, p.targets[h.node]))
			case "Validation plan":
				validation = append(validation, w.link(subject, p.targets[h.node]))
			default:
				if len(ancestors) > 0 && ancestors[len(ancestors)-1].title == "Open questions" {
					questions = append(questions, w.link(h.title, p.targets[h.node]))
				}
			}
			ancestors = append(ancestors, h)
		}
	}
	for _, d := range p.snapshot.Diagnostics {
		msg := strings.ReplaceAll(attr(d, "message"), filepath.ToSlash(p.root)+"/", "")
		dest := p.targets[attr(d, "source")]
		if dest != "" {
			src := p.sources[attr(d, "source")]
			start := 0
			for line := 1; line < model.Number(d["line"]); line++ {
				next := strings.IndexByte(src.data[start:], '\n')
				if next < 0 {
					break
				}
				start += next + 1
			}
			for _, h := range src.headings {
				if h.start <= start && start < h.bodyEnd && p.targets[h.node] != "" {
					dest = p.targets[h.node]
					break
				}
			}
			msg += " · " + w.link(src.path, dest)
		}
		findings = append(findings, msg)
	}
	for _, section := range []struct {
		name, empty string
		items       []string
	}{
		{"Open questions", "No open questions recorded.", questions},
		{"Verification plans", "No verification plans recorded.", verification},
		{"Validation plans", "No validation plans recorded.", validation},
		{"Inspection findings", "No inspection findings.", findings},
	} {
		w.heading(section.name)
		if len(section.items) == 0 {
			w.text(section.empty)
			continue
		}
		slices.Sort(section.items)
		w.out.WriteByte('\n')
		for _, item := range section.items {
			fmt.Fprintln(&w.out, "- "+item)
		}
	}
	w.finish()
}
