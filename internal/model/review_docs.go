package model

import (
	"fmt"
	"net/url"
	"path/filepath"

	"sort"
	"strings"

	"github.com/craigsmitham/dotspec/internal/requirement"
)

// ReviewMarkdown assembles advisory questions from recorded content and links.
// Absence in the acquired scope is not absence of behavior or evidence elsewhere.
func ReviewMarkdown(s *Snapshot, relativeTo string) (string, error) {
	v, err := newReviewView(s, relativeTo)
	if err != nil {
		return "", err
	}
	v.out.WriteString("# Requirements review\n\nGenerated from acquired specifications. Findings below are advisory review prompts, not completeness, acceptance or implementation claims. Missing links mean **no relationship recorded in this scope**. Drafts may conform to the authoring format while requirements verification and validation remain unfinished.\n\n")
	fmt.Fprintf(&v.out, "Model `%s`; mapping `%s`.\n\n", s.ModelVersion, s.MappingProfile)
	v.out.WriteString("## Needs and intent\n\nRequirements validation asks whether the specification expresses the right needs. System validation asks whether a proposed or realized solution serves intended use. The relationships below identify review questions, not either conclusion.\n\n### Goal contribution paths\n\nPaths follow explicit supports, derived-from and refinement links from requirements, concept actions or use cases. A recorded path identifies intent; reviewers must assess whether the behavior is sufficient.\n\n")
	reached := map[string]bool{}
	var reach func(string, map[string]bool)
	reach = func(id string, seen map[string]bool) {
		if seen[id] {
			return
		}
		seen[id] = true
		r := v.d.index[id]
		if String(r, "kind") == "Goal" {
			reached[id] = true
		}
		for _, field := range []string{"supports", "derived_from", "refines"} {
			for _, target := range Strings(r[field]) {
				reach(target, seen)
			}
		}
	}
	for _, r := range v.records {
		if designKind(String(r, "kind"), "Requirement or ConceptAction or UseCase") {
			reach(String(r, "id"), map[string]bool{})
		}
	}
	for _, r := range v.records {
		if String(r, "kind") == "Goal" {
			status := "No contribution path recorded from specified behavior; review required."
			if reached[String(r, "id")] {
				status = "Contribution path recorded; sufficiency requires review."
			}
			fmt.Fprintf(&v.out, "- %s: %s\n", v.link(r), status)
		}
	}
	v.out.WriteString("\n### Scenario grounding\n\nWalk the main path, extensions and guarantees. Could a system meet the written specifications and still prevent the actor from completing a step? An Illustrates link identifies a requirement, business rule or concept action exercised by that step.\n\n")
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
	v.assumptions()
	v.out.WriteString("## Requirement review\n\nReview the criterion recorded on each requirement.\n\n")
	for _, r := range v.records {
		if !designKind(String(r, "kind"), "Requirement") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s\n\n", v.link(r))
		v.body(r, "Criterion", Object(r["criterion"])["content"])
		if missing := v.d.measurementGuidance(r).Missing(); len(missing) > 0 {
			fmt.Fprintf(&v.out, "Review required: %s\n\n", requirement.MeasurementMessage(missing))
		}
	}
	v.questions()
	v.out.WriteString("## Requirements quality\n\nRequirements verification examines individual statements and the set for clarity, consistency and assessability against declared quality criteria. Terminology checks and inspection findings below supply selected review prompts; passing them does not establish complete requirements verification.\n\n")
	v.glossary()
	if len(s.Diagnostics) > 0 {
		v.out.WriteString("## Automated inspection findings\n\n")
		for _, f := range s.Diagnostics {
			fmt.Fprintf(&v.out, "- **%s** `%s`: %s\n", String(f, "severity"), String(f, "code"), strings.ReplaceAll(String(f, "message"), "\n", " "))
		}
	}
	return projectReading(s, strings.TrimRight(v.out.String(), "\n")+"\n"), nil
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
			if !v.body(r, "Project references", Object(r["project_references"])["content"]) {
				v.out.WriteString("No project references recorded. Link schedules, responsibilities, risks and decisions where they are maintained.\n\n")
			}
		}
	}
	v.boundaries()
	v.assumptions()
	v.out.WriteString("## Goals\n\n")
	for _, r := range v.records {
		if !designKind(String(r, "kind"), "Goal or Job or JobStep or DesiredOutcome") {
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
	v.out.WriteString("## System\n\nSpecifications below describe the solution's subjects, concepts, interactions and requirements. Declaration, structural parentage, applicability and behavioral relationships retain their distinct meanings.\n\n")
	for _, r := range v.records {
		if !designKind(String(r, "kind"), "System or Component or Feature or Capability or Concept or ConceptUse or Synchronization or Interface or UseCase or Requirement") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s: %s\n\n", String(r, "kind"), v.link(r))
		for _, f := range []string{"description", "framing", "purpose", "scope", "responsibilities", "statement", "goal", "configuration", "when", "where", "then", "semantics"} {
			v.body(r, reviewLabel(f), r[f])
		}
		if definition := String(r, "definition"); definition != "" {
			fmt.Fprintf(&v.out, "Definition: %s.\n\n", v.link(v.d.index[definition]))
		}
		for _, f := range []string{"contracts", "counterpart"} {
			v.body(r, reviewLabel(f), Object(r[f])["content"])
		}
		for _, f := range []string{"supports", "applies_to", "involves", "groups", "interacts_with"} {
			if ids := Strings(r[f]); len(ids) > 0 {
				fmt.Fprintf(&v.out, "%s: %s.\n\n", reviewLabel(f), v.links(ids))
			}
		}
	}
	v.questions()
	return projectReading(s, strings.TrimRight(v.out.String(), "\n")+"\n"), nil
}

type reviewView struct {
	d        *designResolver
	snapshot *Snapshot
	records  []Record
	base     string
	out      strings.Builder
}

func newReviewView(s *Snapshot, relativeTo string) (*reviewView, error) {
	s = s.navigation()
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
		if String(f, "source") == String(r, "source_document") && namesFragment(r, f) {
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
	if label != "" {
		fmt.Fprintf(&v.out, "**%s**\n\n", label)
	}
	fmt.Fprintf(&v.out, "%s\n\n", designHeadings(rebaseDesignMarkdown(text, v.source(r), v.d.definitions(r)), 4))
	return true
}
func (v *reviewView) boundaries() {
	v.out.WriteString("## Environment and boundaries\n\nOperating conditions and constraints describe what the solution must accommodate. Assumptions are assessed separately for evidence and justified reliance. Explicit exclusions delimit scope. External partners are separate from structural parents.\n\n")
	for _, r := range v.records {
		kind := String(r, "kind")
		if !designKind(kind, "System or Component or BusinessRule") {
			continue
		}
		fmt.Fprintf(&v.out, "### %s: %s\n\n", kind, v.link(r))
		v.body(r, "Environment", r["environment"])
		v.body(r, "Constraints", r["constraints"])
		for _, f := range []string{"statement", "out_of_scope", "basis", "applicability"} {
			value := r[f]
			if kind == "BusinessRule" && f == "applicability" {
				value = Object(value)["content"]
			}
			v.body(r, reviewLabel(f), value)
		}

		if ids := Strings(r["interacts_with"]); len(ids) > 0 {
			fmt.Fprintf(&v.out, "Interacts with: %s.\n\n", v.links(ids))
		}
		for _, field := range []string{"required_interfaces", "provided_interfaces"} {
			if r[field] == nil {
				continue
			}

			items := Strings(r[field])
			if len(items) == 0 {
				v.out.WriteString("Review required: no named interfaces recorded.\n\n")
			}
			for _, id := range items {
				item := v.d.index[id]
				if len(Strings(Object(item["contracts"])["references"])) == 0 {
					fmt.Fprintf(&v.out, "Review required: %s %s has no native contract reference recorded.\n\n", reviewLabel(field), v.d.name(item))
				}
			}
		}
	}
}
func (v *reviewView) assumptions() {
	v.out.WriteString("## Assumptions: evidence and justified reliance\n\nAssess both empirical support in the stated conditions and whether relying on the premise is appropriate for intended use. An observation supports a bounded claim, not a guarantee that the premise always holds.\n\n")
	for _, r := range v.records {
		if String(r, "kind") != "Assumption" {
			continue
		}
		fmt.Fprintf(&v.out, "### %s\n\n", v.link(r))
		for _, field := range []string{"statement", "basis", "applicability"} {
			v.body(r, reviewLabel(field), r[field])
		}
		missing := []string{}
		for _, field := range []string{"basis", "applicability"} {
			value := r[field]
			if strings.TrimSpace(v.d.text(value)) == "" {
				missing = append(missing, reviewLabel(field))
			}
		}
		if len(missing) > 0 {
			fmt.Fprintf(&v.out, "Review required: no %s recorded.\n\n", strings.Join(missing, ", "))
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

// OpenQuestions returns the authored question records belonging to an element.
func OpenQuestions(r Record) []Record {
	return records(r["open_questions"])
}

func (v *reviewView) questions() {
	v.out.WriteString("## Open questions\n\nConsolidated authored questions. External tracking references own resolution workflow. Omitted specification fields are not automatically questions.\n\n")
	count := 0
	for _, r := range v.records {
		for _, q := range OpenQuestions(r) {
			count++
			fmt.Fprintf(&v.out, "### Question %d — %s\n\n", count, v.link(r))
			v.body(r, "Question", q["question"])

		}
	}
	if count == 0 {
		v.out.WriteString("No explicit open questions recorded in this scope.\n\n")
	}
}
func (v *reviewView) glossary() {
	v.out.WriteString("### Terminology review\n\nUse one preferred term consistently in specification prose; declared alternative names and synonym entries remain useful. Review preferred terms and aliases within each context of a system. Inline terms and vocabularies without a named context share the system default. This report does not scan prose for undeclared synonyms.\n\n")
	type designation struct {
		record Record
		role   string
	}
	terms := map[string][]designation{}
	for _, sys := range v.records {
		if String(sys, "kind") != "System" {
			continue
		}
		for _, term := range SystemGlossary(sys, v.d.index) {
			r := v.d.index[term]
			names := []struct{ text, role string }{{v.d.name(r), "preferred term"}}
			for _, alias := range reviewObjects(r["alternative_names"]) {
				names = append(names, struct{ text, role string }{strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(v.d.text(alias)), "-*+ ")), "alternative name"})
			}
			context := String(v.d.index[String(r, "owner")], "glossary_context")
			scope := String(sys, "id") + "\x00" + strings.ToLower(context)
			for _, name := range names {
				key := scope + "\x00" + strings.ToLower(strings.TrimSpace(name.text))
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
		systemID, rest, _ := strings.Cut(key, "\x00")
		_, name, _ := strings.Cut(rest, "\x00")
		scope := v.d.name(v.d.index[systemID])
		if context := String(v.d.index[String(entries[0].record, "owner")], "glossary_context"); context != "" {
			scope = context + " (" + scope + ")"
		}
		fmt.Fprintf(&v.out, "- Review required for **%s** in %s: %s.\n", name, scope, strings.Join(labels, ", "))
	}
	if count == 0 {
		v.out.WriteString("No preferred term or alias is shared by different terms within the same system and context. This does not establish terminology consistency.\n")
	}
	v.out.WriteString("\n")
}
