package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ParseContent retains original bytes before interpreting syntax. It performs no
// I/O and imposes no dotspec authoring validity requirements on Markdown or YAML.
func ParseContent(locator string, data []byte) *ParsedDocument {
	data = bytes.Clone(data)
	key := strings.TrimPrefix(SourceID(locator), "source:")
	p := &ParsedDocument{Data: data, NodeByID: map[string]Record{}}
	p.Source = SourceRecord(locator, data)
	p.Source["media_type"] = "text/markdown"
	p.Document = Record{"id": "document:" + key, "kind": "Document", "source": p.Source["id"], "parse_status": "parsed", "parser": "goldmark/1.8.7 (CommonMark, GFM, footnotes); yaml.v3/3.0.6", "domain_decomposition": "not_attempted", "coverage": Record{"status": "not_checked", "mapping_profile": MappingProfile}}
	e := &contentExtractor{p: p, key: key, counts: map[string]int{}}
	e.normalize()
	root := e.node("document", 0, len(e.buffer))
	p.Document["root"] = root["id"]
	root["span"] = e.rawSpan(0, len(data))
	body, front := e.frontmatter()
	children := []string{}
	if front != "" {
		children = append(children, front)
	}
	e.base = body
	footnotes := &sourceFootnotes{}
	md := goldmark.New(goldmark.WithExtensions(extension.GFM, footnotes), goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	tree := md.Parser().Parse(text.NewReader(e.buffer[body:]))
	footnotes.restore(tree)
	for n := tree.FirstChild(); n != nil; n = n.NextSibling() {
		if _, ok := n.(*east.FootnoteList); ok {
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if id := e.markdown(c); id != "" {
					children = append(children, id)
				}
			}
			continue
		}
		if id := e.markdown(n); id != "" {
			children = append(children, id)
		}
	}
	children = e.fillGaps(root, children)
	root["children"] = children
	e.sections(children)
	e.referenceDefinitions()
	if !utf8.Valid(data) {
		e.diagnostic("content.encoding", "Source contains invalid UTF-8; exact bytes are retained.", 0, len(e.buffer))
		p.Document["parse_status"] = "partial"
	}
	if len(p.Diagnostics) > 0 {
		ids := []string{}
		for _, d := range p.Diagnostics {
			ids = append(ids, String(d, "id"))
		}
		p.Document["diagnostics"] = ids
	}
	return p
}

type contentExtractor struct {
	p       *ParsedDocument
	key     string
	buffer  []byte
	offsets []int
	base    int
	counts  map[string]int
}

func (e *contentExtractor) normalize() {
	data := e.p.Data
	i := 0
	segments := []Record{}
	if bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		i = 3
		segments = append(segments, Record{"normalized_start": 0, "normalized_end": 0, "original_start": 0, "original_end": 3, "transformation": "BOM removal"})
	}
	e.offsets = append(e.offsets, i)
	identityStart, identityNorm := i, 0
	flush := func(end int) {
		if end > identityStart {
			segments = append(segments, Record{"normalized_start": identityNorm, "normalized_end": len(e.buffer), "original_start": identityStart, "original_end": end, "transformation": "identity"})
		}
	}
	changed := i != 0
	for i < len(data) {
		if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
			flush(i)
			start := len(e.buffer)
			e.buffer = append(e.buffer, '\n')
			e.offsets = append(e.offsets, i+2)
			segments = append(segments, Record{"normalized_start": start, "normalized_end": start + 1, "original_start": i, "original_end": i + 2, "transformation": "CRLF to LF"})
			i += 2
			identityStart, identityNorm = i, len(e.buffer)
			changed = true
			continue
		}
		e.buffer = append(e.buffer, data[i])
		i++
		e.offsets = append(e.offsets, i)
	}
	flush(i)
	if changed {
		e.p.Document["normalization"] = segments
	}
}
func (e *contentExtractor) rawSpan(start, end int) Record {
	return Record{"source": e.p.Source["id"], "start": start, "end": end}
}
func (e *contentExtractor) node(kind string, start, end int) Record {
	start = max(0, min(start, len(e.buffer)))
	end = max(start, min(end, len(e.buffer)))
	key := fmt.Sprintf("%s:%d:%d", kind, e.offsets[start], e.offsets[end])
	count := e.counts[key]
	e.counts[key]++
	id := fmt.Sprintf("node:%s:%s:%d", e.key, key, count)
	r := Record{"id": id, "kind": "ContentNode", "document": e.p.Document["id"], "node_kind": kind, "span": e.rawSpan(e.offsets[start], e.offsets[end])}
	e.p.Nodes = append(e.p.Nodes, r)
	e.p.NodeByID[id] = r
	return r
}
func (e *contentExtractor) diagnostic(code, message string, start, end int) {
	r := Record{"id": fmt.Sprintf("diagnostic:%s:%d", e.key, len(e.p.Diagnostics)), "kind": "Diagnostic", "code": code, "severity": "error", "message": message, "source": e.p.Source["id"], "phase": "parse", "span": e.rawSpan(e.offsets[max(0, min(start, len(e.buffer)))], e.offsets[max(0, min(end, len(e.buffer)))])}
	e.p.Diagnostics = append(e.p.Diagnostics, r)
}
func (e *contentExtractor) lineEnd(pos int) int {
	pos = max(0, min(pos, len(e.buffer)))
	if n := bytes.IndexByte(e.buffer[pos:], '\n'); n >= 0 {
		return pos + n + 1
	}
	return len(e.buffer)
}
func (e *contentExtractor) lineStart(pos int) int {
	pos = max(0, min(pos, len(e.buffer)))
	if n := bytes.LastIndexByte(e.buffer[:pos], '\n'); n >= 0 {
		return n + 1
	}
	return 0
}
func (e *contentExtractor) spanBounds(r Record) (int, int) {
	s := Object(r["span"])
	return Number(s["start"]), Number(s["end"])
}
func (e *contentExtractor) normalizedOffset(raw int) int { return sort.SearchInts(e.offsets, raw) }

// fillGaps retains syntax delimiters and whitespace omitted by the parser. Child
// nodes keep their own spans; no content node is invented from a rendered view.
func (e *contentExtractor) fillGaps(parent Record, children []string) []string {
	sort.SliceStable(children, func(i, j int) bool {
		a, _ := e.spanBounds(e.p.NodeByID[children[i]])
		b, _ := e.spanBounds(e.p.NodeByID[children[j]])
		return a < b
	})
	begin, end := e.spanBounds(parent)
	out := []string{}
	at := begin
	add := func(a, b int) {
		if a >= b {
			return
		}
		start, stop := e.normalizedOffset(a), e.normalizedOffset(b)
		if start > len(e.buffer) {
			start = len(e.buffer)
		}
		if stop > len(e.buffer) {
			stop = len(e.buffer)
		}
		r := e.node("trivia", start, stop)
		r["span"] = e.rawSpan(a, b)
		if String(parent, "node_kind") == "document" && len(bytes.TrimSpace(bytes.TrimPrefix(e.p.Data[a:b], []byte{0xef, 0xbb, 0xbf}))) > 0 {
			r["node_kind"] = "opaque"
			r["opaque_reason"] = "Source region not represented by the configured Markdown parser"
		}
		if String(parent, "node_kind") == "text" {
			r["node_kind"] = "text"
			r["decoded_text"] = decodeMarkdownText(string(e.p.Data[a:b]))
		}
		if bytes.HasPrefix(bytes.TrimSpace(e.p.Data[a:b]), []byte("<!--")) {
			r["node_kind"] = "comment"
		}
		out = append(out, String(r, "id"))
	}
	for _, id := range children {
		a, b := e.spanBounds(e.p.NodeByID[id])
		if a > at {
			add(at, a)
		}
		out = append(out, id)
		at = max(at, b)
	}
	if at < end {
		add(at, end)
	}
	return out
}
func (e *contentExtractor) bounds(n ast.Node) (int, int) {
	start := n.Pos()
	end := start
	if start < 0 {
		start = len(e.buffer) - e.base
		end = 0
	}
	if n.Type() != ast.TypeInline {
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			s := lines.At(i)
			start = min(start, s.Start)
			end = max(end, s.Stop)
		}
	}
	switch x := n.(type) {
	case *ast.Text:
		return e.base + x.Segment.Start, e.base + x.Segment.Stop
	case *ast.RawHTML:
		for i := 0; i < x.Segments.Len(); i++ {
			s := x.Segments.At(i)
			start = min(start, s.Start)
			end = max(end, s.Stop)
		}
	case *ast.HTMLBlock:
		if x.HasClosure() {
			end = max(end, x.ClosureLine.Stop)
		}
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if _, ok := c.(*east.FootnoteBacklink); ok {
			continue
		}
		a, b := e.bounds(c)
		start = min(start, a-e.base)
		end = max(end, b-e.base)
	}
	start += e.base
	end += e.base
	start = max(e.base, min(start, len(e.buffer)))
	end = max(start, min(end, len(e.buffer)))
	switch x := n.(type) {
	case *ast.Heading:
		end = e.lineEnd(max(start, end-1))
		if start < len(e.buffer) && e.buffer[start] != '#' && end < len(e.buffer) {
			next := strings.TrimSpace(string(e.buffer[end:e.lineEnd(end)]))
			if next != "" && (strings.Trim(next, "=") == "" || strings.Trim(next, "-") == "") {
				end = e.lineEnd(end)
			}
		}
	case *ast.FencedCodeBlock:
		openEnd := e.lineEnd(start)
		end = max(end, openEnd)
		// Lines are Goldmark's actual code body. Only the line immediately after
		// that body can close this node; fence-looking body lines are literal data.
		fence := bytes.TrimLeft(e.buffer[start:openEnd], " \t")
		if len(fence) > 0 && end < len(e.buffer) {
			char := fence[0]
			width := 0
			for width < len(fence) && fence[width] == char {
				width++
			}
			lineEnd := e.lineEnd(end)
			line := e.buffer[end:lineEnd]
			quotes := 0
			listIndent := 0
			for parent := n.Parent(); parent != nil; parent = parent.Parent() {
				switch parent := parent.(type) {
				case *ast.Blockquote:
					quotes++
				case *ast.ListItem:
					listIndent += parent.Offset
				}
			}
			valid := true
			for range quotes {
				line = bytes.TrimLeft(line, " \t")
				if len(line) == 0 || line[0] != '>' {
					valid = false
					break
				}
				line = line[1:]
				if len(line) > 0 && line[0] == ' ' {
					line = line[1:]
				}
			}
			indent := len(line) - len(bytes.TrimLeft(line, " \t"))
			if quotes == 0 && indent < listIndent {
				valid = false
			}
			line = bytes.TrimSpace(line)
			run := 0
			for run < len(line) && line[run] == char {
				run++
			}
			if valid && run >= width && len(bytes.TrimSpace(line[run:])) == 0 {
				end = lineEnd
			}
		}
	case *ast.Emphasis:
		if first := n.FirstChild(); first != nil {
			a, _ := e.bounds(first)
			start = max(start, a-x.Level)
		}
		end = min(len(e.buffer), end+x.Level)
	case *east.Strikethrough:
		width := 2
		closing := 0
		for end+closing < len(e.buffer) && e.buffer[end+closing] == '~' {
			closing++
		}
		if first := n.FirstChild(); first != nil {
			a, _ := e.bounds(first)
			width = min(2, max(1, a-start), max(1, closing))
			start = max(start, a-width)
		}
		end = min(len(e.buffer), end+width)
	case *ast.CodeSpan:
		run := 0
		for start+run < len(e.buffer) && e.buffer[start+run] == '`' {
			run++
		}
		for i := start + run; i < len(e.buffer); {
			if e.buffer[i] != '`' {
				i++
				continue
			}
			j := i
			for j < len(e.buffer) && e.buffer[j] == '`' {
				j++
			}
			if j-i == run {
				end = j
				break
			}
			i = j
		}
	case *ast.Link, *ast.Image:
		end = e.linkEnd(start, end)
		if next := n.NextSibling(); next != nil {
			a, _ := e.bounds(next)
			if a >= start {
				end = min(end, a)
			}
		}
	case *ast.AutoLink:
		end = start + len(x.Label(e.buffer[e.base:]))
		if start < len(e.buffer) && e.buffer[start] == '<' {
			end += 2
		}
		end = min(end, len(e.buffer))
	case *east.TaskCheckBox:
		end = min(start+3, len(e.buffer))
	case *east.FootnoteLink:
		if start+1 < len(e.buffer) && e.buffer[start] == '!' && e.buffer[start+1] == '[' {
			start++
		}
		if rel := bytes.IndexByte(e.buffer[start:], ']'); rel >= 0 {
			end = start + rel + 1
		}
	case *ast.ThematicBreak:
		end = e.lineEnd(start)
	case *east.TableHeader, *east.TableRow, *east.Table:
		end = e.lineEnd(max(start, end-1))
	}
	return start, end
}
func (e *contentExtractor) linkEnd(start, end int) int {
	i := start
	if i < len(e.buffer) && e.buffer[i] == '!' {
		i++
	}
	if i >= len(e.buffer) || e.buffer[i] != '[' {
		return end
	}
	close := e.balanced(i, '[', ']')
	if close < 0 {
		return end
	}
	end = close + 1
	i = end
	if i < len(e.buffer) && (e.buffer[i] == '(' || e.buffer[i] == '[') {
		right := byte(')')
		if e.buffer[i] == '[' {
			right = ']'
		}
		if j := e.balanced(i, e.buffer[i], right); j >= 0 {
			return j + 1
		}
	}
	return end
}
func (e *contentExtractor) balanced(start int, left, right byte) int {
	depth := 0
	for i := start; i < len(e.buffer); i++ {
		c := e.buffer[i]
		if c == '\\' {
			i++
			continue
		}
		if left == '[' && c == '`' {
			run := 1
			for i+run < len(e.buffer) && e.buffer[i+run] == '`' {
				run++
			}
			if close := bytes.Index(e.buffer[i+run:], bytes.Repeat([]byte{'`'}, run)); close >= 0 {
				i += run + close + run - 1
				continue
			}
		}
		if left == '(' && ((c == '"' || c == '\'') && i > start && (e.buffer[i-1] == ' ' || e.buffer[i-1] == '\t') || c == '<') {
			quote := c
			if c == '<' {
				quote = '>'
			}
			for i++; i < len(e.buffer); i++ {
				if e.buffer[i] == '\\' {
					i++
					continue
				}
				if e.buffer[i] == quote {
					break
				}
			}
			continue
		}
		if c == left {
			depth++
		}
		if c == right {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func syntaxProperty(name, value string) Record {
	return Record{"name": name, "value": value, "value_type": "string"}
}
func (e *contentExtractor) markdown(n ast.Node) string {
	if _, ok := n.(*east.FootnoteBacklink); ok {
		return ""
	}
	start, end := e.bounds(n)
	kind := "opaque"
	switch x := n.(type) {
	case *ast.Heading:
		kind = "heading"
	case *ast.Paragraph, *ast.TextBlock:
		kind = "paragraph"
	case *ast.List:
		kind = "list"
	case *ast.ListItem:
		kind = "list_item"
	case *ast.Blockquote:
		kind = "block_quote"
	case *ast.FencedCodeBlock:
		kind = "fenced_code"
	case *ast.CodeBlock:
		kind = "indented_code"
	case *ast.ThematicBreak:
		kind = "thematic_break"
	case *ast.HTMLBlock:
		kind = "html_block"
		if bytes.HasPrefix(bytes.TrimSpace(e.buffer[start:end]), []byte("<!--")) {
			kind = "comment"
		}
	case *ast.LinkReferenceDefinition:
		kind = "link_definition"
	case *ast.Text, *ast.String:
		kind = "text"
	case *ast.Emphasis:
		kind = "emphasis"
		if x.Level == 2 {
			kind = "strong"
		}
	case *ast.Link:
		kind = "link"
	case *ast.Image:
		kind = "image"
	case *ast.AutoLink:
		kind = "autolink"
	case *ast.CodeSpan:
		kind = "code_span"
	case *ast.RawHTML:
		kind = "inline_html"
		if bytes.HasPrefix(bytes.TrimSpace(e.buffer[start:end]), []byte("<!--")) {
			kind = "comment"
		}
	case *east.Table:
		kind = "table"
	case *east.TableHeader, *east.TableRow:
		kind = "table_row"
	case *east.TableCell:
		kind = "table_cell"
	case *east.Strikethrough:
		kind = "strikethrough"
	case *east.Footnote:
		kind = "footnote_definition"
	case *east.FootnoteLink:
		kind = "footnote_reference"
	case *east.TaskCheckBox:
		kind = "text"
	}
	r := e.node(kind, start, end)
	id := String(r, "id")
	if kind == "opaque" {
		r["opaque_reason"] = "Unsupported Goldmark node: " + n.Kind().String()
	}
	switch x := n.(type) {
	case *ast.Heading:
		r["heading_level"] = x.Level
		r["decoded_text"] = e.inlineText(x)
		if anchor, ok := x.AttributeString("id"); ok {
			if value, ok := anchor.([]byte); ok {
				e.fragment(id, string(value), "goldmark/1.8.7 auto-heading-id")
			}
		}
	case *ast.List:
		r["list_ordered"] = x.IsOrdered()
		r["list_start"] = x.Start
		r["list_delimiter"] = string(x.Marker)
	case *ast.FencedCodeBlock:
		if x.Info != nil {
			r["code_info"] = string(x.Info.Segment.Value(e.buffer[e.base:]))
		}
		r["decoded_text"] = string(x.Lines().Value(e.buffer[e.base:]))
	case *ast.CodeBlock:
		r["decoded_text"] = string(x.Lines().Value(e.buffer[e.base:]))
	case *ast.LinkReferenceDefinition:
		r["label"] = string(x.Label)
		r["destination"] = string(x.Destination)
		if len(x.Title) > 0 {
			r["title"] = string(x.Title)
		}
	case *ast.Text:
		r["decoded_text"] = decodeMarkdownText(string(x.Segment.Value(e.buffer[e.base:])))
		if x.IsRaw() {
			r["decoded_text"] = string(x.Segment.Value(e.buffer[e.base:]))
		}
	case *ast.String:
		r["decoded_text"] = string(x.Value)
	case *ast.Link:
		r["destination"] = string(x.Destination)
		if len(x.Title) > 0 {
			r["title"] = string(x.Title)
		}
		e.linkReference(r, "inline_link", string(x.Destination), start, end)
	case *ast.Image:
		r["destination"] = string(x.Destination)
		if len(x.Title) > 0 {
			r["title"] = string(x.Title)
		}
		e.linkReference(r, "image", string(x.Destination), start, end)
	case *ast.AutoLink:
		r["destination"] = string(x.URL(e.buffer[e.base:]))
		r["decoded_text"] = string(x.Label(e.buffer[e.base:]))
		e.reference(r, "autolink", String(r, "destination"))
	case *east.TableHeader:
		r["syntax_properties"] = []Record{syntaxProperty("header", "true")}
	case *east.TableCell:
		r["syntax_properties"] = []Record{syntaxProperty("alignment", x.Alignment.String())}
	case *east.Footnote:
		r["label"] = string(x.Ref)
	case *east.FootnoteLink:
		target := string(e.buffer[start:end])
		target = strings.TrimSuffix(strings.TrimPrefix(target, "[^"), "]")
		r["label"] = target
		e.reference(r, "footnote", target)
	case *east.TaskCheckBox:
		r["task_state"] = "unchecked"
		if x.IsChecked {
			r["task_state"] = "checked"
		}
		r["decoded_text"] = string(e.buffer[start:end])
	}
	children := []string{}
	for c := n.FirstChild(); c != nil; {
		next := c.NextSibling()
		if current, ok := c.(*ast.Text); ok && len(children) > 0 && !current.IsRaw() {
			_, priorEnd := e.spanBounds(e.p.NodeByID[children[len(children)-1]])
			prior := e.normalizedOffset(priorEnd) - e.base
			if current.Segment.Start < prior {
				width := current.Segment.Stop - current.Segment.Start
				candidate := prior + e.base
				if width >= 0 && candidate+width <= len(e.buffer) && bytes.Equal(current.Segment.Value(e.buffer[e.base:]), e.buffer[candidate:candidate+width]) {
					fixed := ast.NewTextSegment(text.NewSegment(prior, prior+width))
					fixed.SetSoftLineBreak(current.SoftLineBreak())
					fixed.SetHardLineBreak(current.HardLineBreak())
					c = fixed
				}
			}
		}
		if first, ok := c.(*ast.Text); ok && !first.IsRaw() {
			last := first
			for next != nil && !last.SoftLineBreak() && !last.HardLineBreak() {
				candidate, ok := next.(*ast.Text)
				if !ok || candidate.IsRaw() || last.Segment.Stop != candidate.Segment.Start {
					break
				}
				last = candidate
				next = next.NextSibling()
			}
			if last != first {
				merged := ast.NewTextSegment(text.NewSegment(first.Segment.Start, last.Segment.Stop))
				merged.SetSoftLineBreak(last.SoftLineBreak())
				merged.SetHardLineBreak(last.HardLineBreak())
				c = merged
			}
		}
		if child := e.markdown(c); child != "" {
			children = append(children, child)
			if x, ok := c.(*ast.Text); ok && (x.SoftLineBreak() || x.HardLineBreak()) {
				a := e.base + x.Segment.Stop
				b := e.lineEnd(a)
				if a < len(e.buffer) && b <= end {
					kind := "soft_break"
					if x.HardLineBreak() {
						kind = "hard_break"
					}
					br := e.node(kind, a, b)
					children = append(children, String(br, "id"))
				}
			}
			if state := String(e.p.NodeByID[child], "task_state"); state != "" {
				r["task_state"] = state
			}
		}
		c = next
	}
	if x, ok := n.(*ast.Text); ok && !x.IsRaw() {
		children = e.textSyntax(start, end)
	}
	if len(children) > 0 {
		r["children"] = e.fillGaps(r, children)
	}
	return id
}
func (e *contentExtractor) fragment(node, anchor, rule string) {
	r := Record{"id": fmt.Sprintf("fragment:%s:%d", e.key, len(e.p.Fragments)), "kind": "SourceFragment", "source": e.p.Source["id"], "content_node": node, "span": e.p.NodeByID[node]["span"], "generated_anchor": anchor, "anchor_rule": rule}
	e.p.Fragments = append(e.p.Fragments, r)
}
func (e *contentExtractor) reference(node Record, form, target string) {
	e.p.References = append(e.p.References, Record{"id": fmt.Sprintf("reference:%s:%d", e.key, len(e.p.References)), "kind": "SourceReference", "source_document": e.p.Source["id"], "content_node": node["id"], "form": form, "authored_target": target})
}
func (e *contentExtractor) linkReference(node Record, form, target string, start, end int) {
	raw := e.buffer[start:end]
	i := 0
	if len(raw) > 0 && raw[0] == '!' {
		i++
	}
	close := e.balanced(start+i, '[', ']')
	if close >= 0 && close+1 < end && e.buffer[close+1] == '[' {
		label := string(e.buffer[close+2 : end-1])
		if label == "" {
			label = string(e.buffer[start+i+1 : close])
		}
		target = label
		if form != "image" {
			form = "reference_link"
		}
		node["label"] = label
	} else if close >= 0 && close+1 == end {
		target = string(e.buffer[start+i+1 : close])
		node["label"] = target
		if form != "image" {
			form = "reference_link"
		}
	}
	e.reference(node, form, target)
}
func (e *contentExtractor) referenceDefinitions() {
	defs := map[string]string{}
	for _, n := range e.p.Nodes {
		if String(n, "node_kind") == "link_definition" || String(n, "node_kind") == "footnote_definition" {
			label := strings.ToLower(strings.Join(strings.Fields(String(n, "label")), " "))
			if _, ok := defs[label]; !ok {
				defs[label] = String(n, "id")
			}
		}
	}
	for _, r := range e.p.References {
		n := e.p.NodeByID[String(r, "content_node")]
		if String(n, "label") != "" {
			label := strings.ToLower(strings.Join(strings.Fields(String(n, "label")), " "))
			if id := defs[label]; id != "" {
				r["definition_node"] = id
			}
		}
	}
}
func (e *contentExtractor) sections(ids []string) {
	stack := []*Section{}
	seenTitle := false
	for _, id := range ids {
		n := e.p.NodeByID[id]
		kind := String(n, "node_kind")
		if kind == "trivia" || kind == "frontmatter" {
			continue
		}
		if kind == "heading" {
			level := Number(n["heading_level"])
			if !seenTitle && level == 1 && len(stack) == 0 {
				e.p.Title = e.p.Text(id)
				e.p.TitleNode = id
				seenTitle = true
				continue
			}
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			s := &Section{ID: "section:" + id, HeadingID: id, Title: e.p.Text(id), Level: level, Content: []string{}, Intro: []string{}}
			if len(stack) == 0 {
				e.p.Sections = append(e.p.Sections, s)
			} else {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, s)
			}
			for _, parent := range stack {
				parent.Content = append(parent.Content, id)
			}
			stack = append(stack, s)
			continue
		}
		if len(stack) == 0 {
			e.p.Preamble = append(e.p.Preamble, id)
		}
		for _, s := range stack {
			s.Content = append(s.Content, id)
			if len(s.Children) == 0 {
				s.Intro = append(s.Intro, id)
			}
		}
	}
}

// Raw returns the exact bytes of a node's source span as a string.
func (p *ParsedDocument) Raw(id string) string {
	n := p.NodeByID[id]
	s := Object(n["span"])
	a, b := Number(s["start"]), Number(s["end"])
	if a < 0 || b < a || b > len(p.Data) {
		return ""
	}
	return string(p.Data[a:b])
}

// Text exposes decoded values where the syntax parser has supplied one, and
// otherwise the authored syntax. It never mutates the retained payload.
func (p *ParsedDocument) Text(id string) string {
	if text, ok := p.NodeByID[id]["decoded_text"].(string); ok {
		return text
	}
	return strings.ReplaceAll(p.Raw(id), "\r\n", "\n")
}
func (p *ParsedDocument) Rich(ids []string) Record {
	nodes := []string{}
	for _, id := range ids {
		if n := p.NodeByID[id]; n != nil && String(n, "node_kind") != "trivia" {
			nodes = append(nodes, id)
		}
	}
	presence := "empty"
	if len(nodes) > 0 {
		presence = "authored"
	}
	return Record{"presence": presence, "content_nodes": nodes, "decomposition": "unsegmented"}
}

var markdownEntity = regexp.MustCompile(`&(?:#[0-9]+|#x[0-9a-fA-F]+|[a-zA-Z][a-zA-Z0-9]+);`)

func decodeMarkdownText(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && strings.ContainsRune(`!"#$%&'()*+,-./:;<=>?@[\]^_`+"`"+`{|}~`, rune(s[i+1])) {
			i++
		}
		b.WriteByte(s[i])
	}
	return html.UnescapeString(b.String())
}
func (e *contentExtractor) textSyntax(start, end int) []string {
	children := []string{}
	for i := start; i < end; i++ {
		if e.buffer[i] == '\\' && i+1 < end && strings.ContainsRune(`!"#$%&'()*+,-./:;<=>?@[\]^_`+"`"+`{|}~`, rune(e.buffer[i+1])) {
			r := e.node("escape", i, i+2)
			r["decoded_text"] = string(e.buffer[i+1 : i+2])
			children = append(children, String(r, "id"))
			i++
			continue
		}
		if e.buffer[i] == '&' {
			loc := markdownEntity.FindIndex(e.buffer[i:end])
			if loc != nil && loc[0] == 0 {
				stop := i + loc[1]
				r := e.node("entity", i, stop)
				r["decoded_text"] = html.UnescapeString(string(e.buffer[i:stop]))
				children = append(children, String(r, "id"))
				i = stop - 1
				continue
			}
		}
		// Goldmark leaves undefined reference-style links as text. Retain the authored
		// notation as an occurrence without pretending it resolved or rendered as a link.
		if e.buffer[i] == '[' {
			close := e.balanced(i, '[', ']')
			if close < 0 || close >= end {
				continue
			}
			stop := close + 1
			label := ""
			kind, form := "link", "reference_link"
			if i+1 < close && e.buffer[i+1] == '^' {
				label = string(e.buffer[i+2 : close])
				kind, form = "footnote_reference", "footnote"
			} else if stop < end && e.buffer[stop] == '[' {
				last := e.balanced(stop, '[', ']')
				if last < 0 || last >= end {
					continue
				}
				label = string(e.buffer[stop+1 : last])
				if label == "" {
					label = string(e.buffer[i+1 : close])
				}
				stop = last + 1
			}
			if label != "" {
				r := e.node(kind, i, stop)
				r["label"] = label
				r["syntax_properties"] = []Record{syntaxProperty("interpretation", "undefined reference notation retained from text")}
				e.reference(r, form, label)
				children = append(children, String(r, "id"))
				i = stop - 1
			}
		}
	}
	return children
}

// Goldmark's HTML-oriented footnote transformer drops unused definitions and
// moves used ones to the end. Keep its syntax parsers but restore source parents
// after inline parsing, so this content model retains definitions and source order.
type footnotePlacement struct{ node, parent ast.Node }
type sourceFootnotes struct{ placements []footnotePlacement }
type sourceFootnoteBlockParser struct {
	parser.BlockParser
	owner *sourceFootnotes
}

func (p *sourceFootnoteBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	p.owner.placements = append(p.owner.placements, footnotePlacement{node, node.Parent()})
	p.BlockParser.Close(node, reader, pc)
}
func (f *sourceFootnotes) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(parser.WithBlockParsers(util.Prioritized(&sourceFootnoteBlockParser{extension.NewFootnoteBlockParser(), f}, 999)), parser.WithInlineParsers(util.Prioritized(sourceFootnoteInline{extension.NewFootnoteParser()}, 101)))
}
func (f *sourceFootnotes) restore(root ast.Node) {
	affected := map[ast.Node]bool{}
	for _, p := range f.placements {
		affected[p.parent] = true
		p.node.Parent().RemoveChild(p.node.Parent(), p.node)
		p.parent.AppendChild(p.parent, p.node)
	}
	var clean func(ast.Node)
	clean = func(n ast.Node) {
		for c := n.FirstChild(); c != nil; {
			next := c.NextSibling()
			if _, ok := c.(*east.FootnoteList); ok {
				n.RemoveChild(n, c)
			} else {
				clean(c)
			}
			c = next
		}
		if affected[n] {
			nodes := []ast.Node{}
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				nodes = append(nodes, c)
			}
			sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Pos() < nodes[j].Pos() })
			for _, c := range nodes {
				n.RemoveChild(n, c)
				n.AppendChild(n, c)
			}
		}
	}
	clean(root)
}

// SourceID identifies a locator within the snapshot independently of read state.
// Callers normalize their locator before extraction, so discovery and references
// share the same identity without imposing filesystem I/O on this package.
func SourceID(locator string) string {
	return fmt.Sprintf("source:%x", sha256.Sum256([]byte(locator)))[:27]
}

// SourceRecord retains acquired bytes, without assuming a Markdown media type.
func SourceRecord(locator string, data []byte) Record {
	sum := sha256.Sum256(data)
	payload := Record{"storage": "inline_base64", "bytes_base64": base64.StdEncoding.EncodeToString(data), "sha256": fmt.Sprintf("%x", sum), "byte_length": len(data)}
	if utf8.Valid(data) {
		payload["text_encoding"] = "UTF-8"
	}
	return Record{"id": SourceID(locator), "kind": "SourceDocument", "locator": locator, "read_status": "read", "payload": payload}
}

type sourceFootnoteInline struct{ parser.InlineParser }

func (p sourceFootnoteInline) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	line, _ := reader.PeekLine()
	if !bytes.HasPrefix(line, []byte("[^")) && !bytes.HasPrefix(line, []byte("![^")) {
		return nil
	}
	return p.InlineParser.Parse(parent, reader, pc)
}

// inlineText reads syntax leaves without the deprecated approximate Node.Text API.
func (e *contentExtractor) inlineText(n ast.Node) string {
	var b strings.Builder
	var visit func(ast.Node)
	visit = func(n ast.Node) {
		switch x := n.(type) {
		case *ast.Text:
			value := string(x.Segment.Value(e.buffer[e.base:]))
			if !x.IsRaw() {
				value = decodeMarkdownText(value)
			}
			b.WriteString(value)
			if x.SoftLineBreak() || x.HardLineBreak() {
				b.WriteByte('\n')
			}
		case *ast.String:
			b.Write(x.Value)
		case *ast.AutoLink:
			b.Write(x.Label(e.buffer[e.base:]))
		default:
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				visit(c)
			}
		}
	}
	visit(n)
	return b.String()
}
