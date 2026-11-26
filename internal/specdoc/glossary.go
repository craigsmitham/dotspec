package specdoc

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type glossaryPassage struct {
	record             model.Record
	context, qualifier string
	incoming           []entry
	source             *sourceFile
	headings           []*heading
	topicsHeading      *heading
	notePrefix         string
	anchors            []string // local to this system's combined glossary
}

func (term glossaryPassage) title() string {
	return term.headings[0].title + term.qualifier
}

// Fold a plain topic list into the term. Rich lists and qualified passages
// retain their Markdown and their own heading.
func (term glossaryPassage) compactTopics() *heading {
	topics := model.Strings(term.record["topics"])
	if len(topics) == 0 {
		return nil
	}
headings:
	for i, h := range term.headings {
		if h.level != term.headings[0].level+1 || h.title != "Topics" || (i+1 < len(term.headings) && term.headings[i+1].level > h.level) {
			continue
		}
		var labels []string
		for line := range strings.SplitSeq(strings.TrimSpace(term.source.data[h.end:h.bodyEnd]), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if len(line) < 3 || !strings.ContainsRune("-*+", rune(line[0])) || line[1] != ' ' {
				continue headings
			}
			labels = append(labels, strings.TrimSpace(line[2:]))
		}
		if slices.Equal(labels, topics) {
			return h
		}
	}
	return nil
}

func (term glossaryPassage) caption(i int) string {
	if i == 0 {
		return inlineText(term.title())
	}
	h := term.headings[i]
	if h.level == term.headings[0].level+1 {
		switch h.title {
		case "Relationships", "Example":
			return inlineText(term.title() + " " + strings.ToLower(h.title))
		}
	}
	return inlineText(h.title)
}

func (p *planner) planGlossary(s *subject) {
	names := map[string]int{}
	for _, term := range model.SystemGlossary(s.r, p.elements) {
		r := p.elements[term]
		src := p.sources[attr(r, "source_document")]
		node := nameNode(r)
		start := slices.IndexFunc(src.headings, func(h *heading) bool { return h.node == node })
		if start < 0 {
			continue
		}
		end := start + 1
		for end < len(src.headings) && src.headings[end].level > src.headings[start].level {
			end++
		}
		names[strings.ToLower(src.headings[start].title)]++
		context := attr(p.elements[attr(r, "owner")], "glossary_context")
		s.glossary = append(s.glossary, glossaryPassage{record: r, context: context, source: src, headings: src.headings[start:end]})
	}
	if len(s.glossary) == 0 {
		return
	}
	idx := s.index("Glossary")
	local := map[string]string{}
	p.glossaryTargets[idx.file] = local
	p.glossaryTargets[s.file] = local
	slices.SortFunc(s.glossary, func(a, b glossaryPassage) int {
		left, right := a.headings[0], b.headings[0]
		return cmp.Or(strings.Compare(strings.ToLower(left.title), strings.ToLower(right.title)), strings.Compare(strings.ToLower(a.context), strings.ToLower(b.context)), strings.Compare(a.source.path, b.source.path), cmp.Compare(left.start, right.start))
	})
	prefixes := map[*sourceFile]string{}
	ids := parser.NewContext().IDs()
	ids.Generate([]byte("Glossary"), ast.KindHeading)
	for i := range s.glossary {
		term := &s.glossary[i]
		if prefixes[term.source] == "" {
			prefixes[term.source] = fmt.Sprintf("source%d-", len(prefixes)+1)
		}
		term.notePrefix = prefixes[term.source]
		if names[strings.ToLower(term.headings[0].title)] > 1 {
			context := term.context
			if context == "" {
				context = s.name.title
			}
			term.qualifier = " (" + context + ")"
		}
		term.topicsHeading = term.compactTopics()
		for j, h := range term.headings {
			var anchor string
			if h == term.topicsHeading {
				anchor = term.anchors[0]
			} else {
				anchor = string(ids.Generate([]byte(term.caption(j)), ast.KindHeading))
			}
			term.anchors = append(term.anchors, anchor)
			local[sourceKey(term.source.path, h.anchor)] = idx.file + "#" + anchor
			if term.source == s.source {
				h.inGlossary = true
				p.targets[h.node] = idx.file + "#" + anchor
			}
		}
	}
	p.glossaryBacklinks(s, idx.file)
}

// Only links in a term's authored Relationships contribute incoming links.
// Each backlink returns to the passage that explains the direction and meaning.
func (p *planner) glossaryBacklinks(s *subject, file string) {
	// Join resolved reference occurrences to the terms included by this system.
	headings := map[string]*glossaryPassage{}
	for i := range s.glossary {
		term := &s.glossary[i]
		for _, h := range term.headings {
			headings[h.node] = term
		}
	}
	fragments := map[string]*glossaryPassage{}
	for _, fragment := range p.snapshot.Fragments {
		if term := headings[attr(fragment, "content_node")]; term != nil {
			fragments[id(fragment)] = term
		}
	}
	targets := map[string]*glossaryPassage{}
	for _, occurrence := range p.snapshot.References {
		resolution := model.Object(occurrence["resolution"])
		if target := fragments[attr(resolution, "target")]; target != nil {
			targets[id(occurrence)] = target
		}
	}
	links := map[string]model.Record{}
	for _, ref := range p.snapshot.References {
		if targets[id(ref)] != nil {
			links[id(ref)] = p.nodes[attr(ref, "content_node")]
		}
	}
	for _, term := range s.glossary {
		for _, ref := range model.Strings(model.Object(term.record["relationships"])["references"]) {
			node := links[ref]
			if attr(node, "node_kind") != "link" {
				continue
			}
			target := targets[ref]
			pos := position(node, "start")
			i := slices.IndexFunc(term.headings, func(h *heading) bool { return pos >= h.end && pos < h.bodyEnd })
			if i < 0 {
				continue
			}
			from := entry{title: term.title(), target: file + "#" + term.anchors[i]}
			from.preview = term.relationshipPreview(term.headings[i], pos)
			if !slices.ContainsFunc(target.incoming, func(e entry) bool {
				return e.title == from.title && e.target == from.target && e.preview == from.preview
			}) {
				target.incoming = append(target.incoming, from)
			}
		}
	}
}

// Quote the paragraph containing the cited term without inferring its inverse.
func (term glossaryPassage) relationshipPreview(h *heading, pos int) string {
	var preview string
	body := []byte(term.source.data[h.end:h.bodyEnd])
	offset := pos - h.end
	_ = ast.Walk(markdownParser.Parse(text.NewReader(body)), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && (n.Kind() == ast.KindParagraph || n.Kind() == ast.KindTextBlock) && n.Lines().Len() > 0 {
			first, last := n.Lines().At(0), n.Lines().At(n.Lines().Len()-1)
			if first.Start <= offset && offset < last.Stop {
				start := h.end
				preview = indexPreview(contentText(term.source.passageMarkdown(start+first.Start, start+last.Stop)))
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return preview
}

// A model reference can share the metadata line when it occupies the final
// paragraph by itself. Keep all other authored blocks and inline links intact.
func glossaryDefinition(body string) (string, string) {
	data := []byte(body)
	last := markdownParser.Parse(text.NewReader(data)).LastChild()
	if last == nil || last.Kind() != ast.KindParagraph || last.ChildCount() != 1 {
		return body, ""
	}
	ref, ok := last.FirstChild().(*ast.Link)
	if !ok || visibleText(ref, data) != "Model reference" {
		return body, ""
	}
	lines := last.Lines()
	start, end := lines.At(0).Start, lines.At(lines.Len()-1).Stop
	if strings.TrimSpace(body[end:]) != "" {
		return body, ""
	}
	return strings.TrimSpace(body[:start]), strings.TrimSpace(body[start:])
}

func renderGlossaryTopics(w *writer, terms []glossaryPassage) {
	topics := map[string][]string{}
	for _, term := range terms {
		for _, topic := range model.Strings(term.record["topics"]) {
			dest := w.link(term.title(), w.file+"#"+term.anchors[0])
			if !slices.Contains(topics[topic], dest) {
				topics[topic] = append(topics[topic], dest)
			}
		}
	}
	if len(topics) == 0 {
		return
	}
	w.text("**Browse by topic**")
	rows := []string{}
	for _, topic := range slices.Sorted(maps.Keys(topics)) {
		rows = append(rows, "- **"+inlineText(topic)+"** · "+strings.Join(topics[topic], ", "))
	}
	w.text(strings.Join(rows, "\n"))
}

func (p *planner) renderGlossary(w *writer, s *subject) {
	renderGlossaryTopics(w, s.glossary)
	w.notePrefixes = map[*sourceFile]string{}
	ranges := map[*sourceFile][]glossaryPassage{}
	for _, term := range s.glossary {
		w.notePrefixes[term.source] = term.notePrefix
		ranges[term.source] = append(ranges[term.source], term)
		first := term.headings[0]
		for i, h := range term.headings {
			if h == term.topicsHeading {
				continue
			}
			w.text(strings.Repeat("#", min(6, 2+h.level-first.level)) + " " + term.caption(i))
			// Native fragments use the plain caption. Keep references authored in a
			// heading accessible even when their destination changes on this page.
			if slices.ContainsFunc(term.source.links, func(n model.Record) bool { return position(n, "start") >= h.start && position(n, "end") <= h.end }) {
				w.text(p.headingMarkdown(w, term.source, h))
			}
			if h.subjectKind != "" {
				doc := subjectDocument(w.file, term.source, h)
				meta, _ := json.Marshal(doc)
				w.text("<!-- dotspec-doc " + string(meta) + " -->")
			}
			if i == 0 && term.qualifier == "" && term.context != "" {
				w.text("Context · " + inlineText(term.context))
			}
			if i != 0 {
				p.body(w, term.source, h)
				continue
			}
			body, ref := glossaryDefinition(p.rebasedText(w, term.source, h.end, h.bodyEnd))
			var metadata []string
			w.text(body)
			if term.topicsHeading != nil {
				var labels []string
				for _, topic := range model.Strings(term.record["topics"]) {
					labels = append(labels, inlineText(topic))
				}
				metadata = append(metadata, strings.Join(labels, ", "))
			}
			if ref != "" {
				metadata = append(metadata, ref)
			}
			if source := w.passageSource("Term source", relativeURL(w.file, term.source.path, first.anchor)); source != "" {
				metadata = append(metadata, source)
			}
			w.text(strings.Join(metadata, " · "))
			for _, origin := range term.incoming {
				line := "From " + w.link(inlineText(origin.title), origin.target)
				if origin.preview != "" {
					line += " · " + inlineText(origin.preview)
				}
				w.text(line)
			}
		}
	}
	for _, term := range s.glossary {
		passages := ranges[term.source]
		if len(passages) == 0 {
			continue
		}
		p.appendSourceFootnotes(w, term.source, func(pos int) bool {
			return slices.ContainsFunc(passages, func(passage glossaryPassage) bool {
				return pos >= passage.headings[0].start && pos < passage.headings[len(passage.headings)-1].bodyEnd
			})
		})
		delete(ranges, term.source)
	}
}

// Footnote labels belong to their defining source. Qualify the labels at
// their parsed byte spans so identically named notes can share one page.
func (src *sourceFile) footnoteEdits(start, end int, prefix string) []replacement {
	var edits []replacement
	if prefix == "" {
		return nil
	}
	for _, n := range src.footnotes {
		a, b := position(n, "start"), position(n, "end")
		if a < start || b > end {
			continue
		}
		label := "[^" + attr(n, "label") + "]"
		offset := strings.Index(src.data[a:b], label)
		if offset >= 0 {
			at := a + offset + 2
			edits = append(edits, replacement{at, at, prefix})
		}
	}
	return edits

}
