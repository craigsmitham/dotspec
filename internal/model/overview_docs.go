package model

import (
	"fmt"
	"net/url"
	"path/filepath"

	"sort"
	"strings"
)

// OverviewMarkdown describes the acquired specification subjects and features.
// Only authored part_of edges form the structural diagram; obligations are not
// inherited, and feature involvement is derived solely from grouped obligations.
func OverviewMarkdown(s *Snapshot, relativeTo string) (string, error) {
	s = s.navigation()
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
	sourceURL := func(r Record) string {
		locator := String(d.index[String(r, "source_document")], "locator")
		rel, e := filepath.Rel(base, filepath.FromSlash(locator))
		if e != nil {
			rel = locator
		}
		return (&url.URL{Path: filepath.ToSlash(rel)}).String()
	}
	link := func(id string) string {
		r := d.index[id]
		if r == nil {
			return "Unresolved or undeclared"
		}
		target, _ := url.Parse(sourceURL(r))
		for _, f := range s.Fragments {
			if String(f, "source") == String(r, "source_document") && namesFragment(r, f) && String(f, "generated_anchor") != "" {
				target.Fragment = String(f, "generated_anchor")
				break
			}
		}
		label := strings.NewReplacer("[", "\\[", "]", "\\]", "|", "\\|", "\n", " ").Replace(d.name(r))
		return "[" + label + "](<" + target.String() + ">)"
	}
	links := func(ids []string) string {
		values := []string{}
		for _, id := range ids {
			values = append(values, link(id))
		}
		if len(values) == 0 {
			return "None resolved"
		}
		return strings.Join(values, ", ")
	}
	var out strings.Builder
	out.WriteString("# System specification overview\n\nGenerated from authored specifications. Regenerate this view instead of editing it.\n\n")
	fmt.Fprintf(&out, "Model `%s`; mapping `%s`. This view covers only acquired records within the export's discovery and direct-reference boundary. Missing links do not establish missing behavior or missing implementation.\n\n", s.ModelVersion, s.MappingProfile)
	out.WriteString("## System structure\n\nArrows read **part of**, from a component or subsystem to its explicit structural parent. Declaration location and directory nesting do not create these edges. No requirements are inherited down the structure.\n\n")
	ids := map[string]string{}
	for _, r := range elements {
		if designKind(String(r, "kind"), "System or Component") {
			ids[String(r, "id")] = fmt.Sprintf("n%d", len(ids)+1)
		}
	}
	if len(ids) > 0 {
		out.WriteString("```mermaid\nflowchart BT\n")
		for _, r := range elements {
			if id := ids[String(r, "id")]; id != "" {
				fmt.Fprintf(&out, "  %s[\"%s: %s\"]\n", id, String(r, "kind"), designDiagramLabel(d.name(r)))
			}
		}
		for _, r := range elements {
			if source, target := ids[String(r, "id")], ids[String(r, "part_of")]; source != "" && target != "" {
				fmt.Fprintf(&out, "  %s -->|part of| %s\n", source, target)
			}
		}
		out.WriteString("```\n\n")
	} else {
		out.WriteString("No systems or components were acquired.\n\n")
	}
	out.WriteString("| Subject | Explicit parent | Declared in |\n| --- | --- | --- |\n")
	for _, r := range elements {
		if ids[String(r, "id")] != "" {
			declared := "System document"
			if owner := String(r, "owner"); owner != "" {
				declared = link(owner)
			}
			fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "id")), link(String(r, "part_of")), declared)
		}
	}
	out.WriteString("\n## Interfaces\n\nRequired and provided are relative to the declaring subject. Each entry is an addressable element; native contract references do not prove implementation or compatibility.\n\n| Subject | Direction | Interface |\n| --- | --- | --- |\n")
	for _, r := range elements {
		for _, field := range []string{"required_interfaces", "provided_interfaces"} {
			for _, id := range Strings(r[field]) {
				fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "id")), strings.TrimSuffix(field, "_interfaces"), link(id))
			}
		}
	}
	out.WriteString("\n## Features\n\nFeature membership is many-to-many and does not own the grouped records. Component involvement below is derived only from the direct `applies_to` targets of grouped requirements or obligation clauses. It is not an implementation or allocation claim.\n\n| Feature | Owner | Grouped specifications | Components involved (derived) |\n| --- | --- | --- | --- |\n")
	for _, r := range elements {
		if String(r, "kind") != "Feature" {
			continue
		}
		components := []string{}
		for _, member := range Strings(r["groups"]) {
			grouped := d.index[member]
			if !designKind(String(grouped, "kind"), "Obligation") {
				continue
			}
			for _, subject := range Strings(grouped["applies_to"]) {
				if String(d.index[subject], "kind") == "Component" {
					components = appendUnique(components, subject)
				}
			}
		}
		fmt.Fprintf(&out, "| %s | %s | %s | %s |\n", link(String(r, "id")), link(String(r, "specification_owner")), links(Strings(r["groups"])), links(components))
	}
	out.WriteString("\n## Applicable obligations\n\nOnly explicit applicability to acquired systems and components is listed. Declaration inside a component does not implicitly supply an applicability target.\n\n| Obligation | Subjects | Criterion |\n| --- | --- | --- |\n")
	guidance := func(recorded bool) string {
		if recorded {
			return "Recorded in source"
		}
		return "Not recorded"
	}
	for _, r := range elements {
		if !designKind(String(r, "kind"), "Obligation") {
			continue
		}
		subjects := []string{}
		for _, id := range Strings(r["applies_to"]) {
			if ids[id] != "" {
				subjects = append(subjects, id)
			}
		}
		if len(subjects) > 0 {
			fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "id")), links(subjects), guidance(strings.TrimSpace(d.text(Object(r["criterion"])["content"])) != ""))
		}
	}
	out.WriteString("\n## Specification association and discovery context\n\nAn explicit specification system overrides the directory default. These associations do not establish structural parentage.\n\n| Specification | System | Origin |\n| --- | --- | --- |\n")
	for _, r := range elements {
		if designKind(String(r, "kind"), "SystemContextSpecification") {
			origin := String(r, "specification_system_origin")
			if origin == "" {
				origin = "None"
			}
			fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "id")), link(String(r, "specification_system")), origin)
		}
		if String(r, "kind") == "System" && String(r, "enclosing_system") != "" {
			fmt.Fprintf(&out, "| %s | %s | enclosing system by directory discovery |\n", link(String(r, "id")), link(String(r, "enclosing_system")))
		}
	}
	out.WriteString("\n## Concept uses\n\n| System | Named use | Concept definition |\n| --- | --- | --- |\n")
	for _, r := range elements {
		if String(r, "kind") == "ConceptUse" {
			fmt.Fprintf(&out, "| %s | %s | %s |\n", link(String(r, "owner")), link(String(r, "id")), link(String(r, "definition")))
		}
	}
	out.WriteString("\n## Descriptions and native references\n\nLinked contracts, architecture decisions, and implementation artifacts retain their own authority. Normative commitments belong in requirements or obligation clauses.\n\n")
	for _, r := range elements {
		included := designKind(String(r, "kind"), "System or Component or Feature or Interface")
		if designKind(String(r, "kind"), "Obligation") {
			for _, target := range Strings(r["applies_to"]) {
				included = included || ids[target] != ""
			}
		}
		if !included {
			continue
		}
		fmt.Fprintf(&out, "### %s: %s\n\n%s\n\n", String(r, "kind"), d.name(r), link(String(r, "id")))
		for _, field := range []string{"description", "purpose", "scope", "responsibilities", "boundary_description", "contracts", "counterpart", "architecture", "implementation", "part_of_content", "groups_content", "statement", "criterion"} {
			value := r[field]
			if field == "counterpart" || field == "contracts" || field == "architecture" || field == "implementation" || field == "criterion" {
				value = Object(value)["content"]
			}
			text := d.text(value)
			if strings.TrimSpace(text) == "" {
				continue
			}
			title := strings.ReplaceAll(strings.TrimSuffix(field, "_content"), "_", " ")
			title = strings.ToUpper(title[:1]) + title[1:]
			fmt.Fprintf(&out, "#### %s\n\n%s\n\n", title, designHeadings(rebaseDesignMarkdown(text, sourceURL(r), d.definitions(r)), 5))
		}
	}
	if len(s.Diagnostics) > 0 {
		out.WriteString("## Inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
	}
	review, err := newReviewView(s, relativeTo)
	if err != nil {
		return "", err
	}
	review.boundaries()
	review.assumptions()
	out.WriteString(review.out.String())
	return projectReading(s, strings.TrimRight(out.String(), "\n")+"\n"), nil
}
