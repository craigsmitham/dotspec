package model

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

type traceEdge struct{ source, field, target string }

// TraceMarkdown presents explicit instance relationships and their inverses.
// Reachability only identifies review candidates, never satisfaction or evidence.
func TraceMarkdown(s *Snapshot, relativeTo string) (string, error) {
	base, err := designOutputBase(relativeTo)
	if err != nil {
		return "", err
	}
	d := newDesignResolver(s)
	elements := append([]Record{}, s.Elements...)
	key := func(r Record) string {
		source := String(r, "source_document")
		if source == "" {
			source = String(r, "source")
		}
		return String(d.index[source], "locator") + "/" + d.name(r) + "/" + String(r, "id")
	}
	sort.Slice(elements, func(i, j int) bool { return key(elements[i]) < key(elements[j]) })
	name := func(r Record) string {
		n := d.name(r)
		if n == "" {
			n = String(r, "authored_id")
		}
		if n == "" && String(r, "kind") == "SourceFragment" {
			n = String(r, "generated_anchor")
		}
		if n == "" {
			n = String(r, "kind")
		}
		return n
	}
	link := func(id string) string {
		r := d.index[id]
		source := String(r, "source_document")
		if source == "" {
			source = String(r, "source")
		}
		locator := String(d.index[source], "locator")
		rel, e := filepath.Rel(base, filepath.FromSlash(locator))
		if e != nil {
			rel = locator
		}
		anchor := ""
		for _, f := range s.Fragments {
			if String(f, "source") == source && (String(f, "id") == id || slices.Contains(Strings(Object(r["name"])["content_nodes"]), String(f, "content_node"))) && String(f, "generated_anchor") != "" {
				anchor = String(f, "generated_anchor")
				break
			}
		}
		target := (&url.URL{Path: filepath.ToSlash(rel), Fragment: anchor}).String()
		label := strings.NewReplacer("[", "\\[", "]", "\\]", "\n", " ", "|", "\\|").Replace(name(r))
		return "[" + label + "](<" + target + ">)"
	}
	edges := []traceEdge{}
	// These relationships have distinct, documented change-review meanings.
	fields := []string{"supports", "derived_from", "refines", "relies_on", "involves", "applies_to", "definition", "dependencies", "illustrates", "interacts_with"}
	nodes := map[string]bool{}
	for _, r := range elements {
		for _, field := range fields {
			targets := Strings(r[field])
			if field == "definition" {
				if id := String(r, field); id != "" {
					targets = []string{id}
				}
			}
			for _, target := range targets {
				edges = append(edges, traceEdge{String(r, "id"), field, target})
				nodes[String(r, "id")] = true
				nodes[target] = true
			}
		}
	}
	label := func(field string) string {
		switch field {
		case "supports":
			return "supports"
		case "applies_to":
			return "applies to"
		default:
			return strings.ReplaceAll(field, "_", " ")
		}
	}
	var out strings.Builder
	out.WriteString("# Requirements traceability\n\nGenerated from authored specifications. Regenerate this view instead of editing it.\n\n")
	fmt.Fprintf(&out, "Model `%s`; mapping `%s`. Business/user/software subsets describe concerns; core/supporting/generic classifies model subsets. Artifact kinds and relationships have separate meanings.\n\n", s.ModelVersion, s.MappingProfile)
	out.WriteString("This view covers only acquired records and resolved explicit relationships within the export's discovery and direct-reference boundary. Missing links do not establish missing behavior. Review candidates are not proven impacts; paths do not establish satisfaction, achieved outcomes, completeness or normative authority.\n\n")
	out.WriteString("## Why: motivation and basis\n\n| Element | Relationship | Basis or subject |\n| --- | --- | --- |\n")
	for _, e := range edges {
		fmt.Fprintf(&out, "| %s | %s | %s |\n", link(e.source), label(e.field), link(e.target))
	}
	if len(edges) == 0 {
		out.WriteString("| — | No resolved relationships acquired | — |\n")
	}
	out.WriteString("\n## Layered relationship diagram\n\nArrows read from the declaring element to its named target. They do not impose authoring order.\n\n")
	if len(nodes) > 0 {
		out.WriteString("```mermaid\nflowchart BT\n")
		ids := map[string]string{}
		records := []Record{}
		for _, r := range elements {
			if nodes[String(r, "id")] {
				records = append(records, r)
			}
		}
		// Source fragments can be derivation bases without being domain records.
		for _, r := range s.Fragments {
			if nodes[String(r, "id")] {
				records = append(records, r)
			}
		}
		sort.Slice(records, func(i, j int) bool { return key(records[i]) < key(records[j]) })
		for i, r := range records {
			ids[String(r, "id")] = fmt.Sprintf("n%d", i)
		}
		groups := []struct{ id, title string }{{"business-requirements", "Business requirements"}, {"user-requirements", "User requirements"}, {"software-requirements", "Software requirements"}, {"shared", "Shared concerns and source context"}}
		for gi, g := range groups {
			fmt.Fprintf(&out, "    subgraph g%d[\"%s\"]\n", gi, g.title)
			for _, r := range records {
				subsets := compiledModel().catalog.Classes[String(r, "kind")].Subsets
				group := "shared"
				if len(subsets) == 1 {
					group = subsets[0]
				}
				if group != "business-requirements" && group != "user-requirements" && group != "software-requirements" {
					group = "shared"
				}
				if group == g.id {
					fmt.Fprintf(&out, "        %s[\"%s: %s\"]\n", ids[String(r, "id")], String(r, "kind"), designDiagramLabel(name(r)))
				}
			}
			out.WriteString("    end\n")
		}
		for _, e := range edges {
			fmt.Fprintf(&out, "    %s -->|\"%s\"| %s\n", ids[e.source], label(e.field), ids[e.target])
		}
		out.WriteString("```\n\n")
	} else {
		out.WriteString("No resolved relationships acquired.\n\n")
	}
	out.WriteString("## Behavior: use cases and concept uses\n\nUse cases and concepts are complementary views. Participation is many-to-many and identifies a named use, not a reusable definition alone.\n\n| Use case | Participating concept uses | Explicit obligations applying to this use case |\n| --- | --- | --- |\n")
	count := 0
	for _, r := range elements {
		if String(r, "kind") != "UseCase" {
			continue
		}
		count++
		uses, obligations := []string{}, []string{}
		for _, id := range Strings(r["involves"]) {
			uses = append(uses, link(id))
		}
		for _, id := range Strings(r["obligations"]) {
			obligations = append(obligations, link(id)+" (declared here)")
		}
		for _, e := range edges {
			if e.target == String(r, "id") && e.field == "applies_to" {
				obligations = append(obligations, link(e.source))
			}
		}
		fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "id")), traceList(uses), traceList(obligations))
	}
	if count == 0 {
		out.WriteString("| — | No use cases acquired | — |\n")
	}
	out.WriteString("\n## Change review\n\nReverse navigation is generated from the same typed fields. Each direct candidate declares the displayed relationship to the changed record. Further candidates follow those links transitively; they require human review, not automatic changes. Containment and general citations are excluded.\n\n| Changed record | Direct review candidates and reasons | Further review candidates |\n| --- | --- | --- |\n")
	targets := sortedKeys(nodes)
	sort.Slice(targets, func(i, j int) bool { return key(d.index[targets[i]]) < key(d.index[targets[j]]) })
	for _, target := range targets {
		direct := []string{}
		seen := map[string]bool{target: true}
		queue := []string{target}
		for _, e := range edges {
			if e.target == target {
				direct = append(direct, link(e.source)+" ("+label(e.field)+")")
				seen[e.source] = true
				queue = append(queue, e.source)
			}
		}
		if len(direct) == 0 {
			continue
		}
		further := []string{}
		for i := 1; i < len(queue); i++ {
			for _, e := range edges {
				if e.target == queue[i] && !seen[e.source] {
					seen[e.source] = true
					queue = append(queue, e.source)
					further = append(further, link(e.source))
				}
			}
		}
		fmt.Fprintf(&out, "| %s | %s | %s |\n", link(target), traceList(direct), traceList(further))
	}
	out.WriteString("\n## System verification planning\n\nThis inventory covers standalone requirements and explicitly distinguished obligation clauses. Other normative prose is not automatically decomposed. Guidance plans assessment of realizations against obligations; it is not an executed check or requirements-quality review. Revision-bound executions, satisfaction claims and evidence verdicts are not represented by this view.\n\n| Obligation | Authored check guidance | Execution evidence |\n| --- | --- | --- |\n")
	count = 0
	for _, r := range elements {
		if !designKind(String(r, "kind"), "Requirement or ObligationClause") {
			continue
		}
		count++
		guidanceParts := []string{}
		for _, field := range []string{"verification_plan", "measurement_plan"} {
			if strings.TrimSpace(d.text(Object(r[field])["content"])) != "" {
				guidanceParts = append(guidanceParts, strings.ReplaceAll(field, "_", " "))
			}
		}
		guidance := "Not recorded"
		if len(guidanceParts) > 0 {
			guidance = "Recorded in source (" + strings.Join(guidanceParts, ", ") + ")"
		}

		fmt.Fprintf(&out, "| %s | %s | Not evaluated |\n", link(String(r, "id")), guidance)
	}
	if count == 0 {
		out.WriteString("| — | No addressable obligations acquired | Not evaluated |\n")
	}
	out.WriteString("\n## Review findings\n\nUnresolved or invalid declarations remain in the source-backed snapshot, even when omitted from the resolved graph. Missing justification or check guidance is advisory.\n\n")
	if len(s.Diagnostics) == 0 {
		out.WriteString("No inspection findings in this acquired scope.\n")
	}
	for _, f := range s.Diagnostics {
		prefix := ""
		if id := String(f, "record"); d.index[id] != nil {
			prefix = link(id) + ": "
		}
		fmt.Fprintf(&out, "- **%s** `%s`: %s%s\n", String(f, "severity"), String(f, "code"), prefix, designCell(String(f, "message")))
	}
	return projectReading(s, out.String()), nil
}
func traceList(values []string) string {
	if len(values) == 0 {
		return "None recorded in acquired scope"
	}
	return strings.Join(values, "; ")
}
