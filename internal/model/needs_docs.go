package model

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// NeedsMarkdown projects authored needs, contributions and capability dependencies.
// It does not turn dependency into containment or contribution into achievement.
func NeedsMarkdown(s *Snapshot, relativeTo string) (string, error) {
	base, err := designOutputBase(relativeTo)
	if err != nil {
		return "", err
	}
	d := newDesignResolver(s)
	elements := []Record{}
	for _, r := range s.Elements {
		if designKind(String(r, "kind"), "Job or JobStep or DesiredOutcome or ProgressCase or Capability or Feature or System or Component or Obligation or UseCase or Concept") {
			elements = append(elements, r)
		}
	}
	sort.SliceStable(elements, func(i, j int) bool {
		a, b := elements[i], elements[j]
		pa, pb := String(d.index[String(a, "source_document")], "locator"), String(d.index[String(b, "source_document")], "locator")
		if pa != pb {
			return pa < pb
		}
		return String(a, "id") < String(b, "id")
	})
	label := func(r Record) string {
		if name := d.name(r); name != "" {
			return name
		}
		if id := String(r, "authored_id"); id != "" {
			return id
		}
		text := strings.Join(strings.Fields(d.text(r["statement"])), " ")
		runes := []rune(text)
		if len(runes) > 80 {
			text = string(runes[:77]) + "..."
		}
		if text != "" {
			return text
		}
		return String(r, "kind")
	}
	sourceURL := func(r Record) string {
		locator := String(d.index[String(r, "source_document")], "locator")
		rel, e := filepath.Rel(base, filepath.FromSlash(locator))
		if e != nil {
			rel = locator
		}
		u := url.URL{Path: filepath.ToSlash(rel)}
		for _, f := range s.Fragments {
			if String(f, "source") == String(r, "source_document") && slices.Contains(Strings(Object(r["name"])["content_nodes"]), String(f, "content_node")) && String(f, "generated_anchor") != "" {
				u.Fragment = String(f, "generated_anchor")
				break
			}
		}
		return u.String()
	}
	escape := func(s string) string {
		return strings.NewReplacer("|", "\\|", "[", "\\[", "]", "\\]", "\n", " ").Replace(s)
	}
	link := func(r Record) string { return "[" + escape(label(r)) + "](<" + sourceURL(r) + ">)" }
	ids := map[string]string{}
	for _, r := range elements {
		ids[String(r, "id")] = fmt.Sprintf("n%d", len(ids)+1)
	}
	type edge struct{ from, to, relation string }
	edges := []edge{}
	for _, r := range elements {
		for _, field := range []string{"supports", "applies_to", "groups", "depends_on", "provided_by", "part_of"} {
			targets := Strings(r[field])
			if target := String(r, field); target != "" {
				targets = []string{target}
			}
			for _, target := range targets {
				if ids[target] != "" {
					edges = append(edges, edge{String(r, "id"), target, strings.ReplaceAll(field, "_", " ")})
				}
			}
		}
		if designKind(String(r, "kind"), "JobStep or DesiredOutcome or ProgressCase or Capability or Component or Feature") {
			if owner := String(r, "owner"); ids[owner] != "" {
				edges = append(edges, edge{owner, String(r, "id"), "declares"})
			}
		}
	}
	connected := map[string]bool{}
	for _, e := range edges {
		if e.relation != "declares" {
			connected[e.from], connected[e.to] = true, true
		}
	}
	var out strings.Builder
	out.WriteString("# Needs and specification map\n\nGenerated from authored specifications. Regenerate this view instead of editing it.\n\n")
	fmt.Fprintf(&out, "Model `%s`; mapping `%s`. Only acquired records and explicit relationships are shown. Missing links do not establish missing needs, capabilities or implementation.\n\n", s.ModelVersion, s.MappingProfile)
	out.WriteString("`declares` locates a record; `applies to` identifies its explicit scope. For an outcome, that scope identifies where it is a success criterion; for an obligation, where the commitment applies. `supports` records intended contribution, not achievement. Capability dependencies are distinct from structural `part of` links. Providers are optional and non-exhaustive; human and procedural enablers remain in the descriptions. Step order and switching timelines are not dependency edges. The diagram focuses on semantic relationships; declaration ownership and progress cases appear in the table and definitions below.\n\n")
	if len(elements) > 0 {
		out.WriteString("```mermaid\nflowchart LR\n")
		for _, r := range elements {
			if String(r, "kind") == "ProgressCase" && !connected[String(r, "id")] {
				continue
			}
			fmt.Fprintf(&out, "  %s[\"%s: %s\"]\n", ids[String(r, "id")], String(r, "kind"), designDiagramLabel(label(r)))
		}
		for _, e := range edges {
			if e.relation == "declares" {
				continue
			}
			fmt.Fprintf(&out, "  %s -->|%s| %s\n", ids[e.from], e.relation, ids[e.to])
		}
		out.WriteString("```\n\n")
	} else {
		out.WriteString("No needs or specification subjects were acquired.\n\n")
	}
	out.WriteString("## Relationships\n\n| From | Relationship | To |\n| --- | --- | --- |\n")
	for _, e := range edges {
		fmt.Fprintf(&out, "| %s | %s | %s |\n", link(d.index[e.from]), e.relation, link(d.index[e.to]))
	}
	out.WriteString("\n## Definitions and evidence\n\nShort diagram labels do not replace full statements. Evidence describes its source and limits; importance, satisfaction and opportunity assessments belong to dated research about a population, not to the outcome itself. Requirements express commitments, separately from desired user results.\n\n")
	for _, r := range elements {
		fmt.Fprintf(&out, "### %s: %s\n\nSource: %s\n\n", String(r, "kind"), escape(label(r)), link(r))
		definitions := []string{}
		for _, n := range s.ContentNodes {
			if String(n, "node_kind") == "link_definition" && String(Object(n["span"]), "source") == String(r, "source_document") {
				definitions = append(definitions, d.text(Record{"content_nodes": []string{String(n, "id")}}))
			}
		}
		for _, field := range []string{"statement", "description", "purpose", "performer", "situation", "motivation", "intended_progress", "current_approach", "desired_progress", "circumstances", "struggling_moment", "participants", "alternatives_considered", "tradeoffs_and_selection_criteria", "enablers", "supporting_content", "supports_content", "applies_to_content", "depends_on_content", "provided_by_content", "forces", "timeline", "evidence"} {
			value := r[field]
			if field == "forces" || field == "timeline" || field == "evidence" {
				value = Object(value)["content"]
			}
			text := strings.TrimSpace(d.text(value))
			if text == "" {
				continue
			}
			title := strings.ReplaceAll(strings.TrimSuffix(field, "_content"), "_", " ")
			if field == "supporting_content" {
				title = "supporting content"
			}
			title = strings.ToUpper(title[:1]) + title[1:]
			fmt.Fprintf(&out, "#### %s\n\n%s\n\n", title, designHeadings(rebaseDesignMarkdown(text, sourceURL(r), strings.Join(definitions, "\n")), 5))
		}
	}
	if len(s.Diagnostics) > 0 {
		out.WriteString("## Inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
	}
	return strings.TrimRight(out.String(), "\n") + "\n", nil
}
