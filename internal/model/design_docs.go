package model

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// DesignMarkdown renders the authored concept descriptions and the explicitly
// resolved composition. It neither executes synchronizations nor infers edges.
// relativeTo is the directory where the returned Markdown will be read.
func DesignMarkdown(s *Snapshot, relativeTo string) (string, error) {
	base, err := designOutputBase(relativeTo)
	if err != nil {
		return "", err
	}
	d := newDesignResolver(s)
	elements := append([]Record{}, s.Elements...)
	sort.SliceStable(elements, func(i, j int) bool {
		a, b := elements[i], elements[j]
		pa, pb := String(d.index[String(a, "source_document")], "locator"), String(d.index[String(b, "source_document")], "locator")
		if pa != pb {
			return pa < pb
		}
		return d.name(a) < d.name(b)
	})
	var out strings.Builder
	out.WriteString("# Concept design reference\n\nGenerated from authored specifications. Regenerate this view instead of editing it.\n\n")
	fmt.Fprintf(&out, "Model `%s`; mapping `%s`. Relationships reflect explicit declarations. These views do not establish behavioral completeness or implementation guarantees.\n\n", s.ModelVersion, s.MappingProfile)
	if len(s.Diagnostics) > 0 {
		out.WriteString("## Inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
		out.WriteString("\n")
	}
	sourceLink := func(r Record) string {
		locator := String(d.index[String(r, "source_document")], "locator")
		rel, e := filepath.Rel(base, filepath.FromSlash(locator))
		if e != nil {
			rel = locator
		}
		return (&url.URL{Path: filepath.ToSlash(rel)}).String()
	}
	section := func(title string, r Record, value any) {
		text := d.text(value)
		// Reference-style definitions may live outside this section. Supply them to
		// the parser, then render references inline so labels cannot collide across concepts.
		definitions := []string{}
		for _, n := range s.ContentNodes {
			if String(n, "node_kind") == "link_definition" && String(Object(n["span"]), "source") == String(r, "source_document") {
				definitions = append(definitions, d.text(Record{"content_nodes": []string{String(n, "id")}}))
			}
		}
		rendered := rebaseDesignMarkdown(text, sourceLink(r), strings.Join(definitions, "\n"))
		if strings.TrimSpace(text) == "" {
			return
		}
		fmt.Fprintf(&out, "#### %s\n\n%s\n\n", title, designHeadings(rendered, 5))
	}
	out.WriteString("## Concepts\n\n")
	concepts := 0
	for _, r := range elements {
		if String(r, "kind") != "Concept" {
			continue
		}
		concepts++
		fmt.Fprintf(&out, "### %s\n\n[Source](%s).\n\n", d.name(r), sourceLink(r))
		section("Purpose", r, r["purpose"])
		section("Validation plan (needs and intended use)", r, Object(r["validation_plan"])["content"])
		section("Operational principles", r, r["operational_principles"])
		section("Type parameters", r, r["type_parameters_content"])
		section("Introduced types", r, r["introduced_types_content"])
		section("State", r, Object(r["abstract_state"])["content"])
		section("Actions", r, Object(r["actions"])["content"])
		if r["effects_table"] != nil {
			section("Effects table", r, Object(r["effects_table"])["content"])
			out.WriteString("Blank effects cells are unspecified, not unchanged. This authored matrix is not a proof of completeness.\n\n")
		}
		section("Explanatory notes", r, r["explanatory_notes"])
		section("Implementation suggestions (nonbinding)", r, r["implementation_suggestions"])
		section("Illustrative examples", r, r["illustrative_examples"])
	}
	if concepts == 0 {
		out.WriteString("No concept definitions were acquired in this scope.\n\n")
	}
	out.WriteString("## System compositions\n\n")
	systems := 0
	for _, sys := range elements {
		if String(sys, "kind") != "System" {
			continue
		}
		_, usesPresent := sys["concept_uses"]
		_, syncsPresent := sys["synchronizations"]
		if !usesPresent && !syncsPresent {
			continue
		}
		systems++
		fmt.Fprintf(&out, "### %s\n\n[Source](%s).\n\n", d.name(sys), sourceLink(sys))
		uses := Strings(sys["concept_uses"])
		if len(uses) > 0 {
			out.WriteString("| Use | Definition | Type bindings |\n| --- | --- | --- |\n")
			for _, id := range uses {
				use := d.index[id]
				definition := d.index[String(use, "definition")]
				definitionText := "Unresolved or undeclared"
				if definition != nil {
					definitionText = fmt.Sprintf("[%s](%s)", designCell(d.name(definition)), sourceLink(definition))
				}
				var bindings []string
				for _, bid := range Strings(use["bindings"]) {
					b := d.index[bid]
					parameter := d.name(d.index[String(b, "parameter")])
					if parameter == "" {
						parameter = "Unresolved parameter"
					}
					target := String(b, "primitive")
					if target == "" {
						target = d.name(d.index[String(b, "supplying_concept_use")]) + "." + d.name(d.index[String(b, "bound_type")])
						if target == "." {
							target = "Unresolved type"
						}
					}
					bindings = append(bindings, designCell(parameter+" = "+target))
				}
				fmt.Fprintf(&out, "| %s | %s | %s |\n", designCell(d.name(use)), definitionText, strings.Join(bindings, "; "))
			}
			out.WriteString("\n")
		}
		out.WriteString("#### Coordination\n\nOnly resolved action references appear as edges. A rule node does not imply a transaction.\n\n```mermaid\nflowchart LR\n")
		ids := map[string]string{}
		for i, id := range uses {
			ids[id] = fmt.Sprintf("u%d", i)
			fmt.Fprintf(&out, "    %s[\"%s\"]\n", ids[id], designDiagramLabel(d.name(d.index[id])))
		}
		for i, id := range Strings(sys["synchronizations"]) {
			sync := d.index[id]
			sid := fmt.Sprintf("s%d", i)
			fmt.Fprintf(&out, "    %s{\"%s\"}\n", sid, designDiagramLabel(d.name(sync)))
			for _, aid := range Strings(sync["when_actions"]) {
				a := d.index[aid]
				if uid := ids[String(a, "concept_use")]; uid != "" && a["action"] != nil {
					fmt.Fprintf(&out, "    %s -->|\"%s\"| %s\n", uid, designDiagramLabel(d.name(d.index[String(a, "action")])), sid)
				}
			}
			for _, aid := range Strings(sync["then_actions"]) {
				a := d.index[aid]
				if uid := ids[String(a, "concept_use")]; uid != "" && a["action"] != nil {
					fmt.Fprintf(&out, "    %s -->|\"%s\"| %s\n", sid, designDiagramLabel(d.name(d.index[String(a, "action")])), uid)
				}
			}
		}
		out.WriteString("```\n\n")
		out.WriteString("#### Application dependencies\n\nAn arrow means the source use requires the target use for its role in this system. Cycles are permitted. Bindings and action calls do not imply these edges.\n\n```mermaid\nflowchart LR\n")
		for _, id := range uses {
			fmt.Fprintf(&out, "    %s[\"%s\"]\n", ids[id], designDiagramLabel(d.name(d.index[id])))
		}
		for _, id := range uses {
			for _, dep := range Strings(d.index[id]["dependencies"]) {
				if target := ids[dep]; target != "" {
					fmt.Fprintf(&out, "    %s --> %s\n", ids[id], target)
				}
			}
		}
		out.WriteString("```\n\n")
		section("Concept use declarations", sys, sys["concept_uses_content"])
		section("Synchronization rules", sys, sys["synchronizations_content"])
	}
	if systems == 0 {
		out.WriteString("No explicit system compositions were acquired. Directory association alone does not declare a concept use.\n\n")
	}
	return strings.TrimSpace(out.String()) + "\n", nil
}
func designCell(s string) string { return strings.NewReplacer("|", "\\|", "\n", " ").Replace(s) }
func designDiagramLabel(s string) string {
	return strings.NewReplacer("&", "&amp;", "\"", "&quot;", "<", "&lt;", ">", "&gt;", "\n", " ", "|", "&#124;").Replace(s)
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
		if a, b := designLinkLabel(string(doc.Data[start:end])); a >= 0 {
			label = apply(start+a, start+b)
		}
		target := u.String()
		if strings.ContainsAny(target, "() ") {
			target = "<" + target + ">"
		}
		title := ""
		if String(n, "title") != "" {
			title = " " + strconv.Quote(String(n, "title"))
		}
		prefix := ""
		if form == "image" {
			prefix = "!"
		}
		edits = append(edits, designEdit{start, end, prefix + "[" + label + "](" + target + title + ")"})
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

// Return the authored label range, excluding link/image wrappers. Content-node
// spans may deliberately include enclosing syntax for source coverage.
func designLinkLabel(raw string) (int, int) {
	start := strings.IndexByte(raw, '[')
	if start < 0 {
		return -1, -1
	}
	depth := 1
	for i := start + 1; i < len(raw); i++ {
		switch raw[i] {
		case '\\':
			i++
		case '`':
			end := i
			for end < len(raw) && raw[end] == '`' {
				end++
			}
			delimiter := raw[i:end]
			if next := strings.Index(raw[end:], delimiter); next >= 0 {
				i = end + next + len(delimiter) - 1
			}
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return start + 1, i
			}
		}
	}
	return -1, -1
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
