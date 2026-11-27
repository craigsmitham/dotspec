package model

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

var (
	// viewLabelEscape makes link text safe in a Markdown link or table cell.
	viewLabelEscape = strings.NewReplacer("[", "\\[", "]", "\\]", "|", "\\|", "\n", " ")
	// viewCellEscape makes plain text safe in a table cell or list item.
	viewCellEscape = strings.NewReplacer("|", "\\|", "\n", " ")
)

// viewLinks links the records of a reading view to their declarations, with
// paths relative to the directory the view is written to.
type viewLinks struct {
	d        *designResolver
	snapshot *Snapshot
	base     string
	// name gives a record's link text before escaping.
	name func(Record) string
	// escape makes link text safe where the view places it.
	escape *strings.Replacer
	// fragment reports whether a source fragment addresses a record. The first
	// fragment that does supplies the link's anchor.
	fragment func(r, f Record) bool
	// ownSource links a record without a source document, such as a source
	// fragment, to the source that contains it.
	ownSource bool
	// anchoredPassages links footnotes in a rebased passage to the record's
	// anchor rather than to the top of its document.
	anchoredPassages bool
	// definitions caches each source document's reference definitions.
	definitions map[string]string
}

func newViewLinks(s *Snapshot, relativeTo string) (*viewLinks, error) {
	base, err := designOutputBase(relativeTo)
	if err != nil {
		return nil, err
	}
	d := newDesignResolver(s)
	return &viewLinks{d: d, snapshot: s, base: base, name: d.name, escape: viewLabelEscape, fragment: anchoredFragment, definitions: map[string]string{}}, nil
}

func anchoredFragment(r, f Record) bool {
	return namesFragment(r, f) && String(f, "generated_anchor") != ""
}

func (l *viewLinks) document(r Record) string {
	source := String(r, "source_document")
	if source == "" && l.ownSource {
		source = String(r, "source")
	}
	return source
}

func (l *viewLinks) sourceURL(r Record) url.URL {
	locator := String(l.d.index[l.document(r)], "locator")
	rel, err := filepath.Rel(l.base, filepath.FromSlash(locator))
	if err != nil {
		rel = locator
	}
	return url.URL{Path: filepath.ToSlash(rel)}
}

// source returns the path of the document that declares r.
func (l *viewLinks) source(r Record) string {
	u := l.sourceURL(r)
	return u.String()
}

// target returns the path of r's document with the anchor of r's fragment.
func (l *viewLinks) target(r Record) string {
	u := l.sourceURL(r)
	document := l.document(r)
	for _, f := range l.snapshot.Fragments {
		if String(f, "source") == document && l.fragment(r, f) {
			u.Fragment = String(f, "generated_anchor")
			break
		}
	}
	return u.String()
}

func (l *viewLinks) label(r Record) string { return l.escape.Replace(l.name(r)) }

func (l *viewLinks) link(r Record) string {
	return "[" + l.label(r) + "](<" + l.target(r) + ">)"
}

// passage renders authored Markdown from r with links rebased onto the view
// and headings starting at level.
func (l *viewLinks) passage(r Record, text string, level int) string {
	source := l.source(r)
	if l.anchoredPassages {
		source = l.target(r)
	}
	return designHeadings(rebaseDesignMarkdown(text, source, l.documentDefinitions(r)), level)
}

type designEdit struct {
	start, end int
	text       string
}

func rebaseDesignMarkdown(text, source, definitions string) string {
	base, err := url.Parse(source)
	if err != nil {
		return text
	}
	doc := ParseContent("render.md", []byte(text+"\n\n"+definitions))
	refs := append([]Record{}, doc.References...)
	sort.SliceStable(refs, func(i, j int) bool {
		a := Object(doc.NodeByID[String(refs[i], "content_node")]["span"])
		b := Object(doc.NodeByID[String(refs[j], "content_node")]["span"])
		return Number(a["end"])-Number(a["start"]) < Number(b["end"])-Number(b["start"])
	})
	edits := []designEdit{}
	apply := func(start, end int) string {
		ordered := append([]designEdit{}, edits...)
		sort.SliceStable(ordered, func(i, j int) bool {
			if ordered[i].start == ordered[j].start {
				return ordered[i].end > ordered[j].end
			}
			return ordered[i].start < ordered[j].start
		})
		var out strings.Builder
		cursor := start
		for _, e := range ordered {
			if e.start < cursor || e.end > end {
				continue
			}
			out.Write(doc.Data[cursor:e.start])
			out.WriteString(e.text)
			cursor = e.end
		}
		out.Write(doc.Data[cursor:end])
		return out.String()
	}
	for _, ref := range refs {
		n := doc.NodeByID[String(ref, "content_node")]
		span := Object(n["span"])
		start, end := Number(span["start"]), Number(span["end"])
		if start < 0 || end > len(text) || end <= start {
			continue
		}
		form := String(ref, "form")
		if form == "footnote" {
			edits = append(edits, designEdit{start, end, "[" + strings.Trim(string(doc.Data[start:end]), "[]") + "](" + source + ")"})
			continue
		}
		if form != "inline_link" && form != "reference_link" && form != "image" {
			continue
		}
		destination := String(n, "destination")
		if destination == "" {
			continue
		}
		u, err := url.Parse(destination)
		if err != nil {
			continue
		}
		if !u.IsAbs() && u.Host == "" && !strings.HasPrefix(u.Path, "/") {
			if u.Path == "" {
				u.Path = base.Path
			} else {
				u.Path = filepath.ToSlash(filepath.Join(filepath.Dir(base.Path), filepath.FromSlash(u.Path)))
			}
		}
		label := ""
		// Content-node spans may include enclosing syntax for source coverage;
		// keep only the authored label inside the link or image wrapper.
		raw := doc.Data[start:end]
		if open := bytes.IndexByte(raw, '['); open >= 0 {
			if close := specfile.LinkLabelEnd(raw, open); close >= 0 {
				label = apply(start+open+1, start+close)
			}
		}
		title := ""
		if String(n, "title") != "" {
			title = " " + specfile.FormatTitle(String(n, "title"))
		}
		prefix := ""
		if form == "image" {
			prefix = "!"
		}
		edits = append(edits, designEdit{start, end, prefix + "[" + label + "](" + specfile.FormatDestination(u.String()) + title + ")"})
	}
	return apply(0, len(text))
}
func designHeadings(text string, level int) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	minLevel := 7
	fence := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			if fence == "" {
				fence = trimmed[:3]
			} else if strings.HasPrefix(trimmed, fence) {
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, "#"))
		if n > 0 && n <= 6 && len(line) > n && line[n] == ' ' && n < minLevel {
			minLevel = n
		}
	}
	fence = ""
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			if fence == "" {
				fence = trimmed[:3]
			} else if strings.HasPrefix(trimmed, fence) {
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, "#"))
		if n > 0 && n <= 6 && len(line) > n && line[n] == ' ' {
			depth := n - minLevel + level
			if depth <= 6 {
				lines[i] = strings.Repeat("#", depth) + line[n:]
			} else {
				lines[i] = "**" + strings.TrimSpace(line[n:]) + "**"
			}
		}
	}
	return strings.Join(lines, "\n")
}

func designOutputBase(relativeTo string) (string, error) {
	base, err := filepath.Abs(relativeTo)
	if err != nil {
		return "", err
	}
	// Source discovery resolves symlinks. Resolve the existing ancestor of an
	// output directory too, including when the output directory is not created yet.
	ancestor := base
	for {
		if _, e := os.Stat(ancestor); e == nil {
			break
		}
		parent := filepath.Dir(ancestor)
		if parent == ancestor {
			break
		}
		ancestor = parent
	}
	if physical, e := filepath.EvalSymlinks(ancestor); e == nil {
		suffix, _ := filepath.Rel(ancestor, base)
		base = filepath.Join(physical, suffix)
	}
	return base, nil
}

// Markdown views can be generated directly from an exported snapshot. Extract
// reference definitions once per source, including definitions outside the field.
func (l *viewLinks) documentDefinitions(owner Record) string {
	source := String(owner, "source_document")
	if defs, ok := l.definitions[source]; ok {
		return defs
	}
	text := String(l.d.index[source], "markdown")
	var definitions []string
	if strings.Contains(text, "]:") {
		doc := ParseContent(String(l.d.index[source], "locator"), []byte(text))
		for _, n := range doc.Nodes {
			if String(n, "node_kind") == "link_definition" {
				definitions = append(definitions, doc.Raw(String(n, "id")))
			}
		}
	}
	l.definitions[source] = strings.Join(definitions, "\n")
	return l.definitions[source]
}
