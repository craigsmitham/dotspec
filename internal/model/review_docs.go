package model

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// ReviewMarkdown assembles advisory questions from recorded content and links.
// Absence in the acquired scope is not absence of behavior or evidence elsewhere.
func ReviewMarkdown(s *Snapshot, relativeTo string) (string, error) {
	v, err := newReviewView(s, relativeTo)
	if err != nil {
		return "", err
	}
	v.out.WriteString("# Requirements review\n\nGenerated from acquired specifications. Findings below are advisory review prompts, not completeness, acceptance or implementation claims. Missing links mean **no relationship recorded in this scope**. Drafts remain valid.\n\n")
	fmt.Fprintf(&v.out, "Model `%s`; mapping `%s`.\n\n", s.ModelVersion, s.MappingProfile)
	v.out.WriteString("## Goal contribution paths\n\nPaths follow explicit supports, derived-from and refinement links from obligations, concept actions or use cases. A recorded path identifies intent; reviewers must assess whether the behavior is sufficient.\n\n")
	reached := map[string]bool{}
	var reach func(string, map[string]bool)
	reach = func(id string, seen map[string]bool) {
		if seen[id] {
			return
		}
		seen[id] = true
		r := v.d.index[id]
		if String(r, "kind") == "SystemGoal" {
			reached[id] = true
		}
		for _, field := range []string{"supports", "derived_from", "refines"} {
			for _, target := range Strings(r[field]) {
				reach(target, seen)
			}
		}
	}
	for _, r := range v.records {
		if designKind(String(r, "kind"), "Obligation or ConceptAction or UseCase") {
			reach(String(r, "id"), map[string]bool{})
		}
	}
	for _, r := range v.records {
		if String(r, "kind") == "SystemGoal" {
			status := "No contribution path recorded from specified behavior; review required."
			if reached[String(r, "id")] {
				status = "Contribution path recorded; sufficiency requires review."
			}
			fmt.Fprintf(&v.out, "- %s: %s\n", v.link(r), status)
		}
	}
	v.out.WriteString("\n## Scenario grounding\n\nAn Illustrates link identifies a rule or action behind a step. Participation links and ordinary citations do not establish this relationship; not every conversational step needs a separate obligation.\n\n")
	for _, r := range v.records {
		switch String(r, "kind") {
		case "ScenarioStep":
			if ids := Strings(r["illustrates"]); len(ids) > 0 {
				fmt.Fprintf(&v.out, "- %s: illustrates %s.\n", v.link(r), v.links(ids))
			} else {
				fmt.Fprintf(&v.out, "- %s: no rule or action relationship recorded; review required.\n", v.link(r))
			}
		case "Scenario":
			if len(Strings(r["steps"])) == 0 {
				fmt.Fprintf(&v.out, "- %s: no steps recorded; review required.\n", v.link(r))
			}
		}
	}
	v.boundaries()
	v.out.WriteString("## Obligation verification\n\nPlans are guidance for verification or measurement, not evidence of execution.\n\n")
	for _, r := range v.records {
		if designKind(String(r, "kind"), "Obligation") {
			fmt.Fprintf(&v.out, "### %s\n\n", v.link(r))
			present := false
			for _, f := range []string{"verification_plan", "measurement_plan"} {
				if v.body(r, reviewLabel(f), Object(r[f])["content"]) {
					present = true
				}
			}
			v.body(r, "Criterion", Object(r["criterion"])["content"])
			if !present {
				v.out.WriteString("No verification or measurement guidance recorded; review required.\n\n")
			}
		}
	}
	v.questions()
	v.glossary()
	v.out.WriteString("## Explanations and suggestions\n\nThese roles distinguish explanation, nonbinding implementation advice and examples from specified obligations. They do not record stakeholder endorsement.\n\n")
	for _, r := range v.records {
		for _, f := range []string{"explanatory_notes", "implementation_suggestions", "illustrative_examples"} {
			if strings.TrimSpace(v.d.text(r[f])) != "" {
				fmt.Fprintf(&v.out, "### %s — %s\n\n", v.link(r), reviewLabel(f))
				v.body(r, "", r[f])
			}
		}
	}
	if len(s.Diagnostics) > 0 {
		v.out.WriteString("## Inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&v.out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
	}
	return strings.TrimRight(v.out.String(), "\n") + "\n", nil
}

// PEGSMarkdown is a reader perspective, independent of model concern subsets.
// It links project artifacts instead of duplicating their workflow or lifecycle.
func PEGSMarkdown(s *Snapshot, relativeTo string) (string, error) {
	v, err := newReviewView(s, relativeTo)
	if err != nil {
		return "", err
	}
	v.out.WriteString("# Project, Environment, Goals, System\n\nGenerated PEGS reading perspective. This view complements business, user and software requirements subsets; it does not replace them or impose a dependency order. It covers only acquired specifications. Referenced project plans and native contracts remain authoritative.\n\n")
	fmt.Fprintf(&v.out, "Model `%s`; mapping `%s`.\n\n", s.ModelVersion, s.MappingProfile)
	v.out.WriteString("## Project\n\n")
	for _, r := range v.records {
		if String(r, "kind") == "System" {
			fmt.Fprintf(&v.out, "### %s\n\n", v.link(r))
			if !v.body(r, "Project context", Object(r["project_context"])["content"]) {
				v.out.WriteString("No project context recorded. Link schedules, responsibilities, risks and decisions where they are maintained.\n\n")
			}
		}
	}
	v.boundaries()
	v.out.WriteString("## Goals\n\n")
	for _, r := range v.records {
		if !designKind(String(r, "kind"), "SystemGoal or Job or JobStep or DesiredOutcome") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s: %s\n\n", String(r, "kind"), v.link(r))
		for _, f := range []string{"statement", "description", "intended_progress", "situation"} {
			v.body(r, reviewLabel(f), r[f])
		}
		if ids := Strings(r["supports"]); len(ids) > 0 {
			fmt.Fprintf(&v.out, "Supports: %s.\n\n", v.links(ids))
		}
	}
	v.out.WriteString("## System\n\nSpecifications below describe the solution's subjects, concepts, interactions and obligations. Declaration, structural parentage, applicability and behavioral relationships retain their distinct meanings.\n\n")
	for _, r := range v.records {
		if !designKind(String(r, "kind"), "System or Component or Feature or Capability or Concept or UseCase or Obligation") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s: %s\n\n", String(r, "kind"), v.link(r))
		for _, f := range []string{"purpose", "scope", "responsibilities", "statement", "goal"} {
			v.body(r, reviewLabel(f), r[f])
		}
		for _, f := range []string{"supports", "applies_to", "involves", "groups", "interacts_with"} {
			if ids := Strings(r[f]); len(ids) > 0 {
				fmt.Fprintf(&v.out, "%s: %s.\n\n", reviewLabel(f), v.links(ids))
			}
		}
	}
	v.questions()
	return strings.TrimRight(v.out.String(), "\n") + "\n", nil
}

type reviewView struct {
	d        *designResolver
	snapshot *Snapshot
	records  []Record
	base     string
	out      strings.Builder
}

func newReviewView(s *Snapshot, relativeTo string) (*reviewView, error) {
	base, err := designOutputBase(relativeTo)
	if err != nil {
		return nil, err
	}
	v := &reviewView{d: newDesignResolver(s), snapshot: s, base: base, records: append([]Record{}, s.Elements...)}
	sort.SliceStable(v.records, func(i, j int) bool {
		a, b := v.records[i], v.records[j]
		key := func(r Record) string {
			return String(v.d.index[String(r, "source_document")], "locator") + "/" + String(r, "id")
		}
		return key(a) < key(b)
	})
	return v, nil
}
func (v *reviewView) source(r Record) string {
	loc := String(v.d.index[String(r, "source_document")], "locator")
	rel, err := filepath.Rel(v.base, filepath.FromSlash(loc))
	if err != nil {
		rel = loc
	}
	return (&url.URL{Path: filepath.ToSlash(rel)}).String()
}
func (v *reviewView) link(r Record) string {
	name := v.d.name(r)
	if name == "" {
		name = String(r, "kind")
		if label := String(r, "authored_label"); label != "" {
			name += " " + label
		}
	}
	target, _ := url.Parse(v.source(r))
	for _, f := range v.snapshot.Fragments {
		if String(f, "source") == String(r, "source_document") && slices.Contains(Strings(Object(r["name"])["content_nodes"]), String(f, "content_node")) {
			target.Fragment = String(f, "generated_anchor")
			break
		}
	}
	if String(r, "kind") == "ScenarioStep" {
		owner := v.d.index[String(r, "owner")]
		if owner != nil {
			name = v.d.name(owner) + " / " + name
		}
	}
	name = strings.NewReplacer("[", "\\[", "]", "\\]", "\n", " ").Replace(name)
	return "[" + name + "](<" + target.String() + ">)"
}
func (v *reviewView) links(ids []string) string {
	links := []string{}
	for _, id := range ids {
		if r := v.d.index[id]; r != nil {
			links = append(links, v.link(r))
		}
	}
	return strings.Join(links, ", ")
}
func reviewLabel(field string) string {
	s := strings.ReplaceAll(field, "_", " ")
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
func (v *reviewView) body(r Record, label string, value any) bool {
	text := v.d.text(value)
	if strings.TrimSpace(text) == "" {
		return false
	}
	defs := []string{}
	for _, n := range v.snapshot.ContentNodes {
		if String(n, "node_kind") == "link_definition" && String(Object(n["span"]), "source") == String(r, "source_document") {
			defs = append(defs, v.d.text(Record{"content_nodes": []string{String(n, "id")}}))
		}
	}
	if label != "" {
		fmt.Fprintf(&v.out, "**%s**\n\n", label)
	}
	fmt.Fprintf(&v.out, "%s\n\n", designHeadings(rebaseDesignMarkdown(text, v.source(r), strings.Join(defs, "\n")), 4))
	return true
}
func (v *reviewView) boundaries() {
	v.out.WriteString("## Environment and boundaries\n\nOperating conditions and constraints describe what the solution must accommodate. Assumptions are premises relied upon and require a justified basis and validation. Explicit exclusions delimit scope. External partners are separate from structural parents.\n\n")
	for _, r := range v.records {
		kind := String(r, "kind")
		if !designKind(kind, "System or Component or Assumption or Constraint or BusinessRule") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s: %s\n\n", kind, v.link(r))
		v.body(r, "Environment", Object(r["environment"])["content"])
		for _, f := range []string{"statement", "out_of_scope", "basis", "applicability"} {
			value := r[f]
			if kind == "BusinessRule" && f == "applicability" {
				value = Object(value)["content"]
			}
			v.body(r, reviewLabel(f), value)
		}
		if kind == "Assumption" {
			v.body(r, "Validation plan", Object(r["validation_plan"])["content"])
			missing := []string{}
			for _, f := range []string{"basis", "applicability", "validation_plan"} {
				value := r[f]
				if f == "validation_plan" {
					value = Object(value)["content"]
				}
				if strings.TrimSpace(v.d.text(value)) == "" {
					missing = append(missing, reviewLabel(f))
				}
			}
			if len(missing) > 0 {
				fmt.Fprintf(&v.out, "Review required: no %s recorded.\n\n", strings.Join(missing, ", "))
			}
		}
		if ids := Strings(r["interacts_with"]); len(ids) > 0 {
			fmt.Fprintf(&v.out, "Interacts with: %s.\n\n", v.links(ids))
		}
		for _, field := range []string{"required_interfaces", "provided_interfaces"} {
			collection := Object(r[field])
			if collection == nil {
				continue
			}
			v.body(r, reviewLabel(field), collection["content"])
			items := reviewObjects(collection["interfaces"])
			if len(items) == 0 {
				v.out.WriteString("Review required: no named interface descriptions recorded.\n\n")
			}
			for _, item := range items {
				if len(Strings(Object(item["contracts"])["references"])) == 0 {
					fmt.Fprintf(&v.out, "Review required: %s %s has no native contract reference recorded.\n\n", reviewLabel(field), v.d.name(item))
				}
			}
		}
	}
}
func reviewObjects(value any) []Record {
	switch values := value.(type) {
	case []Record:
		return values
	case []any:
		result := []Record{}
		for _, value := range values {
			if r := Object(value); r != nil {
				result = append(result, r)
			}
		}
		return result
	default:
		return nil
	}
}
func (v *reviewView) questions() {
	v.out.WriteString("## Open questions\n\nConsolidated authored questions. External tracking references own resolution workflow. Omitted specification fields are not automatically questions; missing metadata below is advisory.\n\n")
	count := 0
	for _, r := range v.records {
		for _, q := range reviewObjects(Object(r["open_questions"])["questions"]) {
			count++
			fmt.Fprintf(&v.out, "### Question %d — %s\n\n", count, v.link(r))
			v.body(r, "Question", q["question"])
			missing := []string{}
			for _, f := range []string{"raised_by", "raised_on", "resolve_by", "importance", "evidence_needed"} {
				if !v.body(r, reviewLabel(f), q[f]) {
					missing = append(missing, reviewLabel(f))
				}
			}
			v.body(r, "Tracking reference", Object(q["tracking_reference"])["content"])
			if len(missing) > 0 {
				fmt.Fprintf(&v.out, "Review required: no %s recorded.\n\n", strings.Join(missing, ", "))
			}
		}
	}
	if count == 0 {
		v.out.WriteString("No explicit open questions recorded in this scope.\n\n")
	}
}
func (v *reviewView) glossary() {
	v.out.WriteString("## Terminology review\n\nUse one preferred term consistently in specification prose; declared alternative names and synonym entries remain useful. Collisions below are scoped to the declaring System and require human review. This report does not scan prose for undeclared synonyms.\n\n")
	type designation struct {
		record Record
		role   string
	}
	terms := map[string][]designation{}
	for _, r := range v.records {
		if String(r, "kind") == "GlossaryTerm" {
			names := []struct{ text, role string }{{v.d.name(r), "preferred term"}}
			for _, alias := range reviewObjects(r["alternative_names"]) {
				names = append(names, struct{ text, role string }{strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(v.d.text(alias)), "-*+ ")), "alternative name"})
			}
			for _, name := range names {
				key := String(r, "owner") + "\x00" + strings.ToLower(strings.TrimSpace(name.text))
				terms[key] = append(terms[key], designation{r, name.role})
			}
		}
	}
	count := 0
	for _, key := range sortedKeys(terms) {
		entries := terms[key]
		owners := map[string]bool{}
		for _, e := range entries {
			owners[String(e.record, "id")] = true
		}
		if len(owners) < 2 {
			continue
		}
		count++
		labels := []string{}
		for _, e := range entries {
			labels = append(labels, v.link(e.record)+" ("+e.role+")")
		}
		fmt.Fprintf(&v.out, "- Review required for **%s**: %s.\n", strings.SplitN(key, "\x00", 2)[1], strings.Join(labels, ", "))
	}
	if count == 0 {
		v.out.WriteString("No collisions found among recorded preferred terms and alternative names in each System. This does not establish terminology consistency.\n")
	}
	v.out.WriteString("\n")
}
