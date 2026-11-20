package specdoc

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/internal/model"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

type glossaryPassage struct {
	record             model.Record
	context, qualifier string
	incoming           []entry
	source             *sourceFile
	headings           []*heading
	notePrefix         string
	anchors            []string // local to this system's combined glossary
}

func (term glossaryPassage) title() string {
	return term.headings[0].title + term.qualifier
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
	used := map[string]bool{}
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
		for j, h := range term.headings {
			base := "term." + specfile.NameSlug(h.title)
			if j > 0 {
				base = term.anchors[0] + "." + specfile.NameSlug(h.title)
			} else if term.context != "" {
				base = "term." + specfile.NameSlug(term.context) + "." + specfile.NameSlug(h.title)
			}
			anchor := base
			for n := 2; used[anchor]; n++ {
				anchor = fmt.Sprintf("%s-%d", base, n)
			}
			used[anchor] = true
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
	for _, resolution := range p.snapshot.Resolutions {
		if target := fragments[attr(resolution, "target")]; target != nil {
			targets[attr(resolution, "occurrence")] = target
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
			if !slices.Contains(target.incoming, from) {
				target.incoming = append(target.incoming, from)
			}
		}
	}
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
	w.heading("Browse by topic")
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
			anchor := term.anchors[i]
			w.text(`<a id="` + anchor + `"></a>`)
			title := p.headingMarkdown(w, term.source, h)
			if i == 0 {
				title += inlineText(term.qualifier)
			}
			w.text(strings.Repeat("#", min(6, 2+h.level-first.level)) + " " + title)
			if h.subjectKind != "" {
				doc := subjectDocument(w.file, term.source, h)
				meta, _ := json.Marshal(doc)
				w.text("<!-- dotspec-doc " + string(meta) + " -->")
			}
			if i == 0 && term.qualifier == "" && term.context != "" {
				w.text("Context · " + inlineText(term.context))
			}
			p.body(w, term.source, h)
			if i == 0 && len(term.incoming) > 0 {
				var from []string
				for _, origin := range term.incoming {
					from = append(from, w.link(origin.title, origin.target))
				}
				w.text("Relationships from " + strings.Join(from, ", "))
			}
		}
		dest := relativeURL(w.file, term.source.path, first.anchor)
		w.text(link("View term source", dest))
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
