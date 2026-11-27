package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/mermaid"
)

// NeedsMarkdown projects authored needs, contributions and capability dependencies.
// It does not turn dependency into containment or contribution into achievement.
func NeedsMarkdown(s *Snapshot, relativeTo string) (string, error) {
	s = s.navigation()
	l, err := newViewLinks(s, relativeTo)
	if err != nil {
		return "", err
	}
	d := l.d
	elements := []Record{}
	for _, r := range s.Elements {
		if designKind(String(r, "kind"), "Job", "JobStep", "DesiredOutcome", "Capability", "Feature", "System", "Component", "Requirement", "UseCase", "Concept") {
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
	l.anchoredPassages = true
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
		if designKind(String(r, "kind"), "JobStep", "DesiredOutcome", "Capability", "Component", "Feature") {
			if owner := String(r, "owner"); ids[owner] != "" {
				edges = append(edges, edge{owner, String(r, "id"), "declares"})
			}
		}
	}
	var out strings.Builder
	out.WriteString("# Needs and specification map\n\nGenerated from authored specifications. Regenerate this view instead of editing it.\n\n")
	fmt.Fprintf(&out, "Model `%s`; mapping `%s`. Only acquired records and explicit relationships are shown. Missing links do not establish missing needs, capabilities or implementation.\n\n", s.ModelVersion, s.MappingProfile)
	out.WriteString("`declares` locates a record; `applies to` identifies its explicit scope. For an outcome, that scope identifies where it is a success criterion; for a requirement, where the commitment applies. `supports` records intended contribution, not achievement. Capability dependencies are distinct from structural `part of` links. Providers are optional and non-exhaustive; human and procedural enablers remain in the descriptions. Step order does not imply a dependency. The diagram focuses on semantic relationships; declaration ownership appears in the table and definitions below.\n\n")
	if len(elements) > 0 {
		out.WriteString("```mermaid\nflowchart LR\n")
		for _, r := range elements {
			fmt.Fprintf(&out, "  %s[\"%s: %s\"]\n", ids[String(r, "id")], String(r, "kind"), mermaid.Label(d.name(r)))
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
		fmt.Fprintf(&out, "| %s | %s | %s |\n", l.link(d.index[e.from]), e.relation, l.link(d.index[e.to]))
	}
	out.WriteString("\n## Definitions and evidence\n\nShort diagram labels do not replace full statements. Evidence describes its source and limits; importance, satisfaction and opportunity assessments belong to dated research about a population, not to the outcome itself. Requirements express commitments, separately from desired user results.\n\n")
	for _, r := range elements {
		fmt.Fprintf(&out, "### %s: %s\n\nSource: %s\n\n", String(r, "kind"), l.label(r), l.link(r))
		for _, field := range []string{"statement", "description", "purpose", "performer", "situation", "motivation", "intended_progress", "enablers", "supporting_content", "supports_content", "applies_to_content", "depends_on_content", "provided_by_content", "evidence"} {
			value := r[field]
			text := strings.TrimSpace(d.text(value))
			if text == "" {
				continue
			}
			title := strings.ReplaceAll(strings.TrimSuffix(field, "_content"), "_", " ")
			if field == "supporting_content" {
				title = "supporting content"
			}
			title = strings.ToUpper(title[:1]) + title[1:]
			fmt.Fprintf(&out, "#### %s\n\n%s\n\n", title, l.passage(r, text, 5))
		}
	}
	if len(s.Diagnostics) > 0 {
		out.WriteString("## Inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
	}
	return projectReading(s, strings.TrimRight(out.String(), "\n")+"\n"), nil
}
