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
	"strings"

	"github.com/dotspec-md/dotspec/internal/model"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type writer struct {
	page                                         *Page
	file, titleMarkdown, titleAnchor, breadcrumb string
	out                                          strings.Builder
	sourceFiles                                  map[string]bool
	sourceLinks                                  []sourceRoute
	notePrefixes                                 map[*sourceFile]string // keep footnote labels local when combining sources
}

func (p *planner) newPage(file string, document Document, title string) *writer {
	page := &Page{Document: document, Title: title}
	p.plan.Pages[file] = page
	w := &writer{page: page, file: file, titleMarkdown: title, sourceFiles: map[string]bool{}}
	if document.Documents != "" {
		w.sourceFiles[sourceFileKey(file, document.Documents)] = true
	}
	return w
}

// A subject's primary file is represented by its footer. Other files get a
// source link at their first authored passage, in the page's reading order.
func (w *writer) passageSource(label, dest string) string {
	file := sourceFileKey(w.file, dest)
	if file == sourceFileKey(w.file, w.page.Documents) || w.sourceFiles[file] {
		return ""
	}
	w.sourceFiles[file] = true
	return w.sourceLink(label, dest)
}

func (w *writer) sourceLink(label, dest string) string {
	w.sourceLinks = append(w.sourceLinks, sourceRoute{label, dest})
	return link(label, dest)
}

func (w *writer) finish() {
	var header strings.Builder
	header.WriteString("---\ngenerated_by: dotspec\n")
	if w.page.Type != "" {
		fmt.Fprintf(&header, "type: %s\n", w.page.Type)
	}
	if source := w.page.Documents; source != "" {
		fmt.Fprintf(&header, "source: %q\n", source)
	}
	header.WriteString("---\n")
	if w.breadcrumb != "" {
		header.WriteString(w.breadcrumb + "\n\n")
	}
	if w.titleAnchor != "" {
		fmt.Fprintf(&header, "<a id=\"%s\"></a>\n\n", html.EscapeString(w.titleAnchor))
	}
	fmt.Fprintf(&header, "# %s\n", w.titleMarkdown)
	w.page.Markdown = strings.TrimSpace(header.String()+w.out.String()) + "\n"
	// Capture ordinary links with a source-like caption before validating the
	// delivered page. Authored citations are allowed to say "View source" too.
	w.page.otherSourceLabels = sourceLabelAllowances(w.page.Markdown, w.sourceLinks)
}
func (w *writer) text(value string) {
	if value != "" {
		w.out.WriteString("\n" + strings.TrimSpace(value) + "\n")
	}
}
func (w *writer) disclosure(label string, links []string) {
	if len(links) > 0 {
		w.disclosureBody(label, "- "+strings.Join(links, "\n- "))
	}
}
func (w *writer) disclosureBody(label, body string) {
	if body != "" {
		w.text("<details>\n<summary>" + label + "</summary>\n\n" + body + "\n\n</details>")
	}
}
func (w *writer) heading(value string) { w.text("## " + value) }
func (w *writer) link(title, target string) string {
	dest := relative(w.file, target)
	w.requireLink(dest)
	return link(title, dest)
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
	p.indexDetailReferences()
	entry := p.newPage("SPEC.md", Document{}, p.governing.name.title)
	entry.titleMarkdown = inlineText(p.governing.name.title)
	if p.governing.file != "SPEC.md" {
		entry.text(entry.link("Read the system specification", p.governing.file))
	} else {
		for _, title := range slices.Sorted(maps.Keys(p.governing.indexes)) {
			idx := p.governing.indexes[title]
			p.renderIndex(p.governing, title, idx)
			entry.text(entry.link(title, idx.file))
		}
	}
	entry.finish()
	for _, s := range p.subjects {
		p.renderSubject(s)
	}
	return p.plan
}

// The title supplies the current location. Each ancestor leads to its overview.
func (w *writer) setBreadcrumb(owner, member *subject) {
	var steps []string
	for ancestor := owner; ancestor != nil; ancestor = ancestor.navigationParent {
		steps = append(steps, w.link(inlineText(ancestor.name.title), ancestor.file))
	}
	slices.Reverse(steps)
	if owner != nil && member != nil {
		title := member.collection()
		if title == collectionTitle("Capability") && member.system != nil {
			owner = member.system
		}
		if idx := owner.indexes[title]; idx != nil {
			dest := owner.sectionIndex(title)
			steps = append(steps, w.link(title, dest))
		}
	}
	w.breadcrumb = strings.Join(steps, " › ")
}

func (p *planner) renderSubject(s *subject) {
	w := p.newPage(s.file, subjectDocument(s.file, s.source, s.name), s.name.title)
	w.titleMarkdown = p.headingMarkdown(w, s.source, s.name)
	w.titleAnchor = s.name.explicitAnchor
	w.setBreadcrumb(s.navigationParent, s)
	if s.isOverview() {
		p.renderOverview(w, s)
	} else {
		p.renderSubjectBody(w, s)
	}
	p.renderSpecifications(w, s)
	p.renderAssessment(w, s, s.r, "## Verification")
	if s.section("Related specifications") == nil {
		p.relationships(w, s.name, true)
	}
	p.appendFootnotes(w, s)
	for _, anchor := range s.emptyHeadingAnchors {
		w.text(`<a id="` + html.EscapeString(anchor) + `"></a>`)
	}
	w.text("---")
	footer := w.sourceLink("View source", w.page.Documents)
	if models := s.indexes["Quality models"]; s.kind() == "System" && models != nil {
		footer += " · " + w.link("Quality models", models.file)
	}
	w.text(footer)
	w.finish()
}

func (p *planner) renderSubjectBody(w *writer, s *subject) {
	headings := readingOrder(s)
	if hasDescriptionSection(s.kind()) && s.section("Description") != nil {
		end := slices.IndexFunc(headings, func(h *heading) bool {
			return h.level == s.name.level+1 && h.title != "Description"
		})
		if end < 0 {
			end = len(headings)
		}
		p.renderPassage(w, s, headings[:end], 0, false)
		headings = headings[end:]
	}
	var indexes []string
	for _, title := range slices.Sorted(maps.Keys(s.indexes)) {
		idx := s.indexes[title]
		p.renderIndex(s, title, idx)
		line := "- " + w.link(title, idx.file)
		indexes = append(indexes, line)
	}
	w.text(strings.Join(indexes, "\n"))
	if s.kind() == "Job" && !recordsSteps(s) {
		// The job page is the view of its map, so it reports absent steps
		// explicitly, ahead of any cross-references.
		end := slices.IndexFunc(headings, func(h *heading) bool {
			return h.level == s.name.level+1 && h.title == "Related specifications"
		})
		if end < 0 {
			end = len(headings)
		}
		p.renderPassage(w, s, headings[:end], 0, false)
		w.heading("Steps")
		w.text("No steps recorded.")
		headings = headings[end:]
	}
	p.renderPassage(w, s, headings, 0, false)
}

// recordsSteps reports whether a job declares a step or explains its Steps
// section in prose; an empty Steps heading records nothing.
func recordsSteps(s *subject) bool {
	if slices.ContainsFunc(s.source.headings, func(h *heading) bool { return h.subjectKind == "JobStep" }) {
		return true
	}
	steps := s.section("Steps")
	return steps != nil && strings.TrimSpace(s.source.data[steps.end:steps.bodyEnd]) != ""
}

// A system overview can introduce an authored passage under a shared section
// heading or at a deeper level. Its authored hierarchy stays together.
func (p *planner) renderPassage(w *writer, s *subject, headings []*heading, offset int, merged bool) {
	for i, h := range headings {
		if h.owner != s {
			continue
		}
		if h.detailEntry != nil && w.file == s.file {
			if h.detailEntry == h {
				body := sourceProse(s.source, h)
				preview := paragraphPreview(body)
				if h.subjectKind == "Assumption" {
					if preview == "" {
						preview = contentText(body)
					}
					preview = indexPreview(preview)
				}
				w.text(entryLine(w, entry{title: h.title, target: p.targets[h.node], preview: preview}))
			}
			continue
		}
		if h != s.name {
			// Preserve addressability for duplicated headings moved onto different pages.
			if h.explicitAnchor != "" {
				w.text(`<a id="` + html.EscapeString(h.explicitAnchor) + `"></a>`)
			}
			if !merged || h != headings[0] {
				caption := p.headingMarkdown(w, s.source, h)
				w.text(strings.Repeat("#", min(6, max(2, h.level-s.name.level+1+offset))) + " " + caption)
			}
			if h.subjectKind != "" {
				doc := subjectDocument(w.file, s.source, h)
				meta, _ := json.Marshal(doc)
				w.text("<!-- dotspec-doc " + string(meta) + " -->")
				w.text(w.passageSource("View section source", doc.Documents))
			}
		}
		if h == s.name && s.kind() == "Feature" {
			renderIndexMembers(w, "Use cases", p.featureGoals(s))
		}
		p.body(w, s.source, h)
		if h.level == s.name.level+1 && (h.title == "Extensions" || h.title == "Actions" || h.title == "Quality attributes") && !s.isOverview() {
			p.sectionNavigation(w, h, headings[i+1:])
		}
		if !s.isOverview() && h.level == s.name.level+1 {
			if dest := s.sectionIndex(h.title); dest != "" && !isDetailCollection(h.title) {
				w.text(w.link("Browse "+strings.ToLower(h.title), dest))
			}
		}
		if h.level == s.name.level+1 && h.title == "Related specifications" {
			p.relationships(w, s.name, false)
		}
		if h != s.name {
			p.relationships(w, h, false)
		}
	}
}

// Link to the planned fragments, including explicit anchors and duplicate suffixes.
func (p *planner) sectionNavigation(w *writer, parent *heading, rest []*heading) {
	var children []*heading
	for _, h := range rest {
		if h.level <= parent.level {
			break
		}
		if h.owner == parent.owner && h.level == parent.level+1 {
			children = append(children, h)
		}
	}
	if len(children) < 2 {
		return
	}
	var links []string
	for _, child := range children {
		links = append(links, w.link(inlineText(child.title), p.targets[child.node]))
	}
	if len(links) > 10 {
		w.disclosure("Browse "+strings.ToLower(parent.title), links)
	} else {
		w.text(strings.Join(links, " · "))
	}
}

// readingOrder promotes a few entire authored sections for the document's
// purpose. Descendants, scenario steps and all other siblings retain order.
func readingOrder(s *subject) []*heading {
	if s.overview != nil {
		var headings []*heading
		for _, section := range s.overview {
			for _, passage := range section.passages {
				headings = append(headings, passage...)
			}
		}
		return headings
	}
	groups := slices.Clone(s.passages)
	var first []string
	switch s.kind() {
	case "Requirement":
		first = []string{"statement", "applicability", "criterion"}
	case "UseCase":
		first = []string{"goal", "primary_actor", "trigger", "preconditions", "main_success_scenario", "extensions", "success_guarantees", "minimal_guarantees"}
	case "Feature":
		first = []string{"description", "purpose", "overview"}
	default:
		if hasDescriptionSection(s.kind()) {
			first = []string{"description"}
		}
	}
	priority := func(group []*heading) int {
		if group[0] == s.name {
			return 0
		}
		if strings.HasPrefix(group[0].field, "related_specifications") {
			return len(first) + 2 // Keep cross-references after the subject's own details.
		}
		if i := slices.Index(first, strings.Split(group[0].field, ".")[0]); i >= 0 {
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
	onPage := func(pos int) bool {
		h := s.source.headingAt(pos)
		if h == nil || h.owner != s || h.inGlossary {
			return false
		}
		file, _, _ := strings.Cut(p.targets[h.node], "#")
		return file == w.file
	}
	p.appendSourceFootnotes(w, s.source, onPage)
}

func (p *planner) appendSourceFootnotes(w *writer, sf *sourceFile, onPage func(int) bool) {
	definitions := map[string]model.Record{}
	needed := map[string]bool{}
	for _, n := range sf.footnotes {
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
			w.text(p.rebasedText(w, sf, position(n, "start"), position(n, "end")))
		}
	}
}
func (h *heading) reportTitle() string {
	if h.owner != nil && h.owner.kind() == "ConceptUse" && h.level == h.owner.name.level+1 && h.title == "Definition" {
		return relations["definition"].forward.label
	}
	return h.title
}

func (p *planner) headingMarkdown(w *writer, source *sourceFile, h *heading) string {
	if title := h.reportTitle(); title != h.title {
		return title
	}
	raw := p.rebasedText(w, source, h.start, h.end)
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

// kindName names a kind in file suffixes: the authoring form of a root kind,
// otherwise the kind in kebab case.
func kindName(kind string) string {
	if name, ok := specfile.RootKind(kind); ok {
		return name
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

// A report file replaces its source suffix's ".md" with "-report.md".
func reportSuffix(kind string) string {
	return strings.TrimSuffix(specfile.AuthoredSuffix(kindName(kind)), ".md") + "-report.md"
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
	edits = append(edits, sf.footnoteEdits(start, end, w.notePrefixes[sf])...)
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
	return sf.editedText(start, end, edits)
}

func (sf *sourceFile) editedText(start, end int, edits []replacement) string {
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
func sourceProse(sf *sourceFile, h *heading, edits ...replacement) string {
	return sf.passageMarkdown(h.end, h.bodyEnd, edits...)
}

func (sf *sourceFile) passageMarkdown(from, to int, edits ...replacement) string {
	for _, n := range sf.links {
		start, end := position(n, "start"), position(n, "end")
		if start < from || end > to || attr(n, "label") == "" || attr(n, "node_kind") == "link_definition" {
			continue
		}
		edits = append(edits, replacement{start, end, rewriteDestination(sf.data[start:end], n, attr(n, "destination"))})
	}
	return sf.editedText(from, to, edits)
}

func rewriteDestination(raw string, n model.Record, target string) string {
	// Reference definitions may be shared between sections. Expand reference
	// uses inline so each generated page works independently.
	if attr(n, "node_kind") == "link_definition" {
		at := strings.Index(raw, "]:")
		if at < 0 {
			return raw
		}
		return raw[:at+2] + " " + specfile.FormatDestination(target) + linkTitle(n) + "\n"
	}
	if close := specfile.LinkLabelEnd([]byte(raw), strings.IndexByte(raw, '[')); close >= 0 {
		return raw[:close+1] + "(" + specfile.FormatDestination(target) + linkTitle(n) + ")"
	}
	return raw
}
func linkTitle(n model.Record) string {
	if title := attr(n, "title"); title != "" {
		return " " + specfile.FormatTitle(title)
	}
	return ""
}

func (p *planner) indexRelationships() {
	p.routes = map[string][]route{}
	for _, a := range p.snapshot.Relationships() {
		projection := attr(a, "projection_path")
		rel, ok := reportRelation(projection)
		if !ok {
			continue
		}
		left, right := p.elements[attr(a, "subject")], p.elements[attr(a, "object")]
		add := func(from, to model.Record, caption relationCaption, forward bool) {
			dest := p.targets[id(to)]
			h := nameNode(from)
			title := attr(p.nodes[nameNode(to)], "decoded_text")
			if h != "" && title != "" && dest != "" {
				p.routes[h] = append(p.routes[h], route{relationCaption: caption, title: title, target: dest, forward: forward})
			}
		}
		add(left, right, rel.forward, true)
		add(right, left, rel.reverse, false)
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
			field = section.reportTitle()
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
		if h == h.owner.name && r.directoryHome(h.owner, r.target) != "" {
			continue
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
		w.text("**" + label + "**")
		w.text("- " + strings.Join(links, "\n- "))
	}
}
