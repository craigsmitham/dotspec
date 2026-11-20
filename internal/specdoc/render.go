package specdoc

import (
	"cmp"
	"encoding/json"
	"fmt"
	"html"
	"maps"
	"net/url"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const reviewPath = "specification/requirements-review.report.md"

type writer struct {
	page                *Page
	file, titleMarkdown string
	out                 strings.Builder
}

func (p *planner) newPage(file string, document Document, title string) *writer {
	page := &Page{Document: document, Title: title}
	p.plan.Pages[file] = page
	return &writer{page: page, file: file, titleMarkdown: title}
}
func (w *writer) finish() {
	var header strings.Builder
	header.WriteString("---\ngenerated_by: dotspec\n")
	if w.page.Type != "" {
		fmt.Fprintf(&header, "type: %s\n", w.page.Type)
	}
	if w.page.Documents != "" {
		fmt.Fprintf(&header, "documents: %s\n", strconv.Quote(w.page.Documents))
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
	entry := p.newPage("SPEC.md", Document{}, p.governing.name.title)
	entry.titleMarkdown = inlineText(p.governing.name.title)
	entry.text(entry.link("Read the system specification", p.governing.file))
	entry.finish()
	for _, s := range p.subjects {
		p.renderSubject(s)
	}
	p.renderReview()
	return p.plan
}

func (p *planner) renderSubject(s *subject) {
	w := p.newPage(s.file, subjectDocument(s.file, s.source, s.name), s.name.title)
	w.titleMarkdown = p.headingMarkdown(w, s, s.name)
	headings := readingOrder(s)
	var opening []string
	if s.kind() == "System" {
		opening = []string{"Description", "Purpose", "Scope"}
	} else if hasDescriptionSection(s.kind()) && s.section("Description") != nil {
		opening = []string{"Description"}
	}
	if len(opening) > 0 {
		// Keep each opening section and its descendants ahead of navigation.
		end := slices.IndexFunc(headings, func(h *heading) bool {
			return h.level == s.name.level+1 && !slices.Contains(opening, h.title)
		})
		if end < 0 {
			end = len(headings)
		}
		p.renderHeadings(w, s, headings[:end])
		headings = headings[end:]
	}
	if s.kind() == "System" {
		var nav []string
		if s != p.governing {
			nav = append(nav, w.link(s.system.name.title, s.system.file))
		}
		if d := s.directories["Requirements"]; d != nil {
			nav = append(nav, w.link("Requirements", s.file+"#"+d.anchor))
		}
		nav = append(nav, w.link("Requirements review", reviewPath))
		w.text("---")
		w.text(strings.Join(nav, " · "))
	} else {
		parent := s.parent
		if parent == nil {
			parent = s.system
		}
		w.text(w.link(parent.name.title, parent.file))
	}
	var earlyDirectory *directory
	if features := s.directories["Features"]; s.kind() == "System" && features != nil && features.section == nil {
		earlyDirectory = features
		w.heading("Features")
		features.render(w)
	}
	p.renderHeadings(w, s, headings)
	if s.kind() == "System" && len(model.Strings(s.r["quality_attributes"])) > 0 && s.section("Quality attributes") == nil {
		w.heading("Quality attributes")
		p.attributeTable(w, s)
	}
	for _, label := range slices.Sorted(maps.Keys(s.directories)) {
		if d := s.directories[label]; d.section == nil && d != earlyDirectory {
			w.heading(label)
			d.render(w)
		}
	}
	if s.section("Related specifications") == nil {
		p.relationships(w, s.name, true)
	}
	p.appendFootnotes(w, s)
	w.text(w.sourceLink("Edit source", s.source.path, s.name.anchor))
	w.finish()
}

func (p *planner) renderHeadings(w *writer, s *subject, headings []*heading) {
	for _, h := range headings {
		if h.owner != s {
			continue
		}
		if h != s.name {
			// Preserve addressability for duplicated headings moved onto different pages.
			if h.explicitAnchor != "" {
				w.text(`<a id="` + html.EscapeString(h.explicitAnchor) + `"></a>`)
			}
			w.text(strings.Repeat("#", max(2, h.level-s.name.level+1)) + " " + p.headingMarkdown(w, s, h))
			if h.subjectKind != "" {
				meta, _ := json.Marshal(subjectDocument(w.file, s.source, h))
				w.text("<!-- dotspec-doc " + string(meta) + " -->")
			}
			if d := s.directories[h.title]; d != nil && d.section == h {
				d.render(w)
				if strings.TrimSpace(s.source.data[h.end:h.bodyEnd]) != "" {
					// A following authored Markdown list must remain a separate list.
					w.text("<!-- end directory -->")
				}
			}
		}
		if h == s.name && s.kind() == "Feature" {
			p.featureGoals(w, s)
		}
		p.body(w, s.source, h)
		if s.kind() == "System" && h.title == "Quality attributes" && h.level == 2 && len(model.Strings(s.r["quality_attributes"])) > 0 {
			p.attributeTable(w, s)
		}
		if h.level == s.name.level+1 && h.title == "Related specifications" {
			p.relationships(w, s.name, false)
		}
		if h != s.name {
			p.relationships(w, h, false)
			if h.obligation {
				w.text(w.sourceLink("Edit obligation", s.source.path, h.anchor))
			}
		}
	}
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
	var first []string
	switch s.kind() {
	case "Requirement":
		first = []string{"Statement", "Applies to", "Conditions", "Criterion", "Verification plan"}
	case "UseCase":
		first = []string{"Goal", "Primary actor", "Trigger", "Preconditions", "Main success scenario", "Extensions", "Success guarantees", "Minimal guarantees"}
	case "System":
		first = []string{"Description", "Purpose", "Scope", "Features"}
	default:
		if hasDescriptionSection(s.kind()) {
			first = []string{"Description"}
		}
	}
	priority := func(group []*heading) int {
		if group[0] == s.name {
			return 0
		}
		if group[0].title == "Related specifications" {
			return len(first) + 2 // Keep cross-references after the subject's own details.
		}
		if i := slices.Index(first, group[0].title); i >= 0 {
			return i + 1
		}
		return len(first) + 1
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
	onPage := func(pos int) bool {
		for _, h := range s.source.headings {
			if h.start <= pos && pos < h.bodyEnd {
				return h.owner == s
			}
		}
		return false
	}
	for _, n := range s.source.footnotes {
		key := strings.ToLower(attr(n, "label"))
		if attr(n, "node_kind") == "footnote_definition" {
			definitions[key] = n
			continue
		}
		if onPage(position(n, "start")) {
			needed[key] = true
		}
	}
	for _, key := range slices.Sorted(maps.Keys(needed)) {
		if n := definitions[key]; n != nil && !onPage(position(n, "start")) {
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
	if kind == "SystemQualityAttribute" {
		return "quality-attribute"
	}
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
	var breaks []replacement
	if h.qualityAttribute {
		breaks = attributeFieldBreaks(sf.data, h.end, h.bodyEnd)
	}
	w.text(p.rebasedText(w, sf, h.end, h.bodyEnd, breaks...))
}

// The authoring syntax recognizes fields only on lines of direct paragraphs.
// Leave matching text in examples, block quotes and nested lists untouched.
func attributeFieldBreaks(source string, start, end int) []replacement {
	body := []byte(source[start:end])
	doc := markdownParser.Parse(text.NewReader(body))
	var breaks []replacement
	for block := doc.FirstChild(); block != nil; block = block.NextSibling() {
		if block.Kind() != ast.KindParagraph {
			continue
		}
		previousFieldEnd := -1
		for i := 0; i < block.Lines().Len(); i++ {
			segment := block.Lines().At(i)
			line := strings.TrimRight(source[start+segment.Start:start+segment.Stop], "\r\n")
			trimmed := strings.TrimSpace(line)
			if !attributeField(trimmed) {
				previousFieldEnd = -1
				continue
			}
			if previousFieldEnd >= 0 {
				breaks = append(breaks, replacement{previousFieldEnd, previousFieldEnd, "\\"})
			}
			previousFieldEnd = start + segment.Start + len(line)
			if strings.HasSuffix(line, "  ") || strings.HasSuffix(line, "\\") {
				previousFieldEnd = -1 // already a Markdown hard break
			}
		}
	}
	return breaks
}

type replacement struct {
	start, end int
	text       string
}

// Rebase link destinations alongside any requested field breaks. Link edits
// are disjoint even when an image is the label of a link.
func (p *planner) rebasedText(w *writer, sf *sourceFile, start, end int, edits ...replacement) string {
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
	"subattribute_of": {"Parent attribute", "Subattributes"}, "groups": {"Groups", "Feature"}, "refines": {"Refines", "Refined by"}, "supports": {"Supports", "Supported by"},
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

// Only a link in the corresponding authored relationship section fills that
// role. Navigation and links to the same passage for other reasons can coexist.
func (p *planner) relationships(w *writer, h *heading, showHeading bool) {
	if len(p.routes[h.node]) == 0 {
		return
	}
	sf := h.owner.source
	type labeledLink struct{ label, dest string }
	authored := map[labeledLink]bool{}
	field := ""
	for _, section := range sf.headings {
		if section.start <= h.start {
			continue
		}
		if section.level <= h.level {
			break
		}
		if section.level == h.level+1 {
			field = section.title
			if field == "Quality attributes" {
				field = "Quality attribute"
			}
		}
		if field == "" || section.owner != h.owner {
			continue
		}
		for _, n := range sf.links {
			pos := position(n, "start")
			if pos >= section.end && pos < section.bodyEnd && attr(n, "node_kind") == "link" {
				dest := p.documentLink(w.file, sf.path, attr(n, "destination"))
				authored[labeledLink{field, dest}] = true
			}
		}
	}
	byLabel := map[string][]string{}
	for _, r := range p.routes[h.node] {
		if r.target == w.file || (r.forward && authored[labeledLink{r.label, relative(w.file, r.target)}]) {
			continue
		}
		if h == h.owner.name {
			dirs := h.owner.directories
			if r.label == "Parts" && (dirs["Child systems"].contains(r.target) || dirs["Components"].contains(r.target)) {
				continue
			}
			if r.label == "Applies here" && dirs["Applicable requirements"].contains(r.target) {
				continue
			}
		}
		byLabel[r.label] = append(byLabel[r.label], w.link(r.title, r.target))
	}
	if showHeading && len(byLabel) != 0 {
		w.heading("Related specifications")
	}
	for _, label := range slices.Sorted(maps.Keys(byLabel)) {
		links := byLabel[label]
		slices.Sort(links)
		links = slices.Compact(links)
		if len(links) == 1 {
			w.text("**" + label + "** · " + links[0])
		} else {
			w.text("**" + label + "**")
			w.text("- " + strings.Join(links, "\n- "))
		}
	}
}
func (p *planner) renderReview() {
	w := p.newPage(reviewPath, Document{Type: RequirementsReviewReport}, "Requirements review")
	w.text(w.link(p.governing.name.title, p.governing.file))
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
				labels := questionPassages(sourceProse(s.source, h))
				for _, label := range labels {
					questions = append(questions, w.link(inlineText(label), p.targets[h.node])+" · "+inlineText(subject))
				}
				if len(labels) == 0 && h.title == "Open question" {
					questions = append(questions, w.link(subject+" — "+h.title, p.targets[h.node]))
				}
			case "Verification plan":
				verification = append(verification, w.link(subject, p.targets[h.node]))
			case "Validation plan":
				validation = append(validation, w.link(subject, p.targets[h.node]))
			default:
				if len(ancestors) > 0 && ancestors[len(ancestors)-1].title == "Open questions" {
					labels := questionPassages(sourceProse(s.source, h))
					if len(labels) == 0 {
						labels = []string{h.title}
					}
					for _, label := range labels {
						questions = append(questions, w.link(inlineText(label), p.targets[h.node])+" · "+inlineText(s.name.title))
					}
				}
			}
			ancestors = append(ancestors, h)
		}
	}
	byFinding := map[string][]string{}
	for _, d := range p.snapshot.Diagnostics {
		msg := strings.ReplaceAll(attr(d, "message"), filepath.ToSlash(p.root)+"/", "")
		dest := p.targets[attr(d, "source")]
		if dest != "" {
			src := p.sources[attr(d, "source")]
			label := p.plan.Pages[dest].Title
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
					label = h.owner.name.title
					if h != h.owner.name {
						label += " · " + h.title
					}
					break
				}
			}
			byFinding[msg] = append(byFinding[msg], w.link(label, dest))
		} else {
			findings = append(findings, inlineText(msg))
		}
	}
	for msg, passages := range byFinding {
		slices.Sort(passages)
		if len(passages) == 1 {
			findings = append(findings, inlineText(msg)+" · "+passages[0])
		} else {
			findings = append(findings, inlineText(msg)+"\n  - "+strings.Join(passages, "\n  - "))
		}
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

func questionPassages(md string) []string {
	source := []byte(md)
	doc := markdownParser.Parse(text.NewReader(source))
	var passages []string
	add := func(n ast.Node) {
		if label := visibleText(n, source); label != "" {
			passages = append(passages, label)
		}
	}
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if n.Kind() == ast.KindList {
			for item := n.FirstChild(); item != nil; item = item.NextSibling() {
				add(item)
			}
		} else {
			add(n)
		}
	}
	return passages
}
