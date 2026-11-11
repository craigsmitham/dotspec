package model

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func (m *domainMapper) useCase(r Record, s *Section) bool {
	fields := map[string]string{"Goal": "goal", "Trigger": "trigger", "Success criterion": "success_criterion"}
	if field, ok := fields[s.Title]; ok {
		m.section(r, field, s)
		return true
	}
	switch s.Title {
	case "System boundary":
		m.section(r, "boundary_description", s)
		m.designField(r, "subject_system", s, "Subject system")
	case "Primary actor":
		r["primary_actor"] = m.actor(s.Content)
		m.sectionBindings(r, "primary_actor.content", s)
	case "Goal level":
		v := Record{"content": m.doc.Rich(s.Content)}
		text := strings.ToLower(strings.TrimSpace(m.contentText(s.Content)))
		text = strings.TrimSuffix(text, ".")
		switch text {
		case "summary", "kite", "summary (kite)":
			v["level"] = "summary"
		case "user task", "user-task", "sea", "sea level", "user task (sea level)":
			v["level"] = "user_task"
		case "subfunction", "fish", "subfunction (fish)":
			v["level"] = "subfunction"
		}
		r["goal_level"] = v
		m.sectionBindings(r, "goal_level.content", s)
	case "Preconditions":
		r["preconditions"] = m.conditions(s)
		m.sectionBindings(r, "preconditions.content", s)
	case "Minimum guarantees":
		r["minimum_guarantees"] = m.conditions(s)
		m.sectionBindings(r, "minimum_guarantees.content", s)
	case "Supporting actors":
		actors := []Record{}
		for _, id := range m.listItems(s.Intro, false) {
			actors = append(actors, m.actor([]string{id}))
		}
		for _, sub := range s.Children {
			v := m.actor(sub.Content)
			v["role_name"] = m.doc.Rich([]string{sub.HeadingID})
			for _, detail := range sub.Children {
				if detail.Title == "Services" {
					m.setRich(v, "services", detail.Content)
				} else {
					m.extension(r, detail)
				}
			}
			actors = append(actors, v)
		}
		r["supporting_actors"] = Record{"content": m.doc.Rich(s.Content), "actors": actors}
		m.sectionBindings(r, "supporting_actors.content", s)
	case "Stakeholder interests":
		items := []Record{}
		for _, id := range m.listItems(s.Intro, false) {
			items = append(items, Record{"content": m.doc.Rich([]string{id}), "interest": m.doc.Rich([]string{id})})
		}
		for _, sub := range s.Children {
			items = append(items, Record{"content": m.doc.Rich(append([]string{sub.HeadingID}, sub.Content...)), "stakeholder": m.actor([]string{sub.HeadingID}), "interest": m.doc.Rich(sub.Content)})
		}
		r["stakeholder_interests"] = Record{"content": m.doc.Rich(s.Content), "interests": items}
		m.sectionBindings(r, "stakeholder_interests.content", s)
	case "Main success scenario":
		x := m.named("Scenario", r, s)
		m.section(x, "content", s)
		x["steps"] = m.scenarioSteps(x, s.Intro)
		for _, sub := range s.Children {
			m.extension(x, sub)
		}
		r["main_success_scenario"] = String(x, "id")
		m.bind(r, "main_success_scenario", []string{s.HeadingID}, "structured_item", "main-scenario")
	case "Extensions":
		v := m.collection(r, "extensions", "branches", s)
		v["branches"] = append(Strings(v["branches"]), m.listBranches(r, s.Intro)...)
		for _, sub := range s.Children {
			x := m.branch(r, sub)
			appendID(v, "branches", String(x, "id"))
		}
	case "Sequence":
		v := Record{"content": m.doc.Rich(s.Content)}
		mode := strings.ToLower(m.fieldText(s, "Mode"))
		if mode == "representative" || mode == "required" || mode == "partially_constrained" {
			v["mode"] = mode
		}
		items := []Record{}
		for _, id := range m.listItems(s.Intro, false) {
			items = append(items, Record{"statement": m.doc.Rich([]string{id}), "governing_references": m.references([]string{id})})
		}
		if len(items) > 0 {
			v["constraints"] = items
		}
		r["sequence"] = v
		m.sectionBindings(r, "sequence.content", s)
	case "Related jobs":
		r["related_jobs"] = m.context(s.Content)
		m.sectionBindings(r, "related_jobs.content", s)
	default:
		return false
	}
	return true
}

func (m *domainMapper) actor(nodes []string) Record {
	return Record{"content": m.doc.Rich(nodes), "references": m.references(nodes)}
}

func (m *domainMapper) conditions(s *Section) Record {
	items := []Record{}
	for _, id := range m.listItems(s.Intro, false) {
		items = append(items, Record{"statement": m.doc.Rich([]string{id}), "references": m.references([]string{id})})
	}
	for _, child := range s.Children {
		items = append(items, Record{"statement": m.doc.Rich(append([]string{child.HeadingID}, child.Content...)), "references": m.references(child.Content)})
	}
	return Record{"content": m.doc.Rich(s.Content), "conditions": items}
}

func (m *domainMapper) conditionField(owner Record, field string, s *Section) {
	v := m.conditions(s)
	if previous := Object(owner[field]); previous != nil {
		nodes := append(append([]string{}, Strings(Object(previous["content"])["content_nodes"])...), s.Content...)
		v["content"] = m.doc.Rich(nodes)
		old, _ := previous["conditions"].([]Record)
		items, _ := v["conditions"].([]Record)
		v["conditions"] = append(old, items...)
	}
	owner[field] = v
	m.sectionBindings(owner, field+".content", s)
}

func (m *domainMapper) contentText(nodes []string) string {
	var out []string
	for _, id := range nodes {
		if String(m.doc.NodeByID[id], "node_kind") != "trivia" {
			out = append(out, m.doc.Text(id))
		}
	}
	return strings.Join(out, "\n")
}

func (m *domainMapper) scenarioSteps(owner Record, nodes []string) []string {
	steps := []string{}
	for _, id := range nodes {
		n := m.doc.NodeByID[id]
		ordered, _ := n["list_ordered"].(bool)
		if String(n, "node_kind") != "list" {
			continue
		}
		index := Number(n["list_start"])
		if index == 0 {
			index = 1
		}
		for _, item := range Strings(n["children"]) {
			if String(m.doc.NodeByID[item], "node_kind") != "list_item" {
				continue
			}
			text := m.doc.Text(item)
			label := authoredLabel(explicitStep, text)
			if authoredLabel(explicitBranch, text) != "" || (!ordered && label == "") {
				index++
				continue
			}
			if label == "" {
				label = strconv.Itoa(index)
				if match := orderedMarker.FindStringSubmatch(m.doc.Raw(item)); len(match) > 1 {
					label = match[1]
				}
			}
			index++
			x := m.record("ScenarioStep", String(owner, "id"))
			x["authored_label"] = label
			x["content"] = m.doc.Rich([]string{item})
			x["references"] = m.references([]string{item})
			m.bind(x, "content", []string{item}, "structured_item", "labelled-scenario-step")
			steps = append(steps, String(x, "id"))
		}
	}
	return steps
}

var orderedMarker = regexp.MustCompile(`^\s*([0-9]+)[.)]\s+`)

// Labels are cues for explicitly distinguished items, not a control-flow grammar.
// Markdown containment determines ownership; label digits do not resolve steps.
var explicitBranch = regexp.MustCompile(`^\s*(?:[-+*]|[0-9]+[.)])?\s*(?:\*\*)?((?:[0-9]+[a-z])+|\.[0-9]+[a-z]|\*[a-z])(?:[.:)]|\s|\*\*|$)`)
var explicitStep = regexp.MustCompile(`^\s*(?:[-+*]|[0-9]+[.)])?\s*(?:\*\*)?((?:[0-9]+[a-z])+[0-9]+|\.[0-9]+)(?:[.:)]|\s|\*\*|$)`)

func authoredLabel(pattern *regexp.Regexp, text string) string {
	if match := pattern.FindStringSubmatch(text); len(match) > 1 {
		return match[1]
	}
	return ""
}

func (m *domainMapper) branch(owner Record, s *Section) Record {
	x := m.named("UseCaseExtension", owner, s)
	m.section(x, "content", s)
	if label := authoredLabel(explicitBranch, s.Title); label != "" {
		x["authored_label"] = label
	}
	x["references"] = m.references(s.Content)
	x["steps"] = m.scenarioSteps(x, s.Intro)
	x["branches"] = m.listBranches(x, s.Intro)
	for _, pair := range [][2]string{{"Condition", "condition"}, {"Location", "location"}, {"At step", "location"}, {"Point of occurrence", "location"}, {"Handling", "handling"}, {"Response", "handling"}, {"Continuation", "handling"}} {
		if nodes := m.fieldNodes(s, pair[0]); nodes != nil {
			m.setRich(x, pair[1], nodes)
			m.bind(x, pair[1], nodes, "named_field", "extension-field/"+pair[0])
			m.bind(x, pair[1], m.fieldLabelNodes(s, pair[0]), "structural_wrapper", "extension-field-label/"+pair[0])
		}
	}
	for _, sub := range s.Children {
		field := map[string]string{"Condition": "condition", "Location": "location", "At step": "location", "Point of occurrence": "location", "Handling": "handling", "Response": "handling", "Flow": "handling", "Continuation": "handling", "Unresolved handling": "handling"}[sub.Title]
		if field != "" {
			m.section(x, field, sub)
			if field == "handling" {
				// Retain the explicit qualification instead of inventing a status enum.
				if sub.Title == "Unresolved handling" {
					m.setRich(x, field, []string{sub.HeadingID})
				}
				x["steps"] = append(Strings(x["steps"]), m.scenarioSteps(x, sub.Intro)...)
				x["branches"] = append(Strings(x["branches"]), m.listBranches(x, sub.Intro)...)
			}
			for _, detail := range sub.Children {
				m.branchDetail(x, detail)
			}
			continue
		}
		m.branchDetail(x, sub)
	}
	// Explicit fields may be authored in any order. Preserve source order when
	// several old Response/Flow/Continuation fields contribute to one narrative.
	for _, field := range []string{"condition", "location", "handling"} {
		if value := Object(x[field]); value != nil {
			nodes := Strings(value["content_nodes"])
			sort.SliceStable(nodes, func(i, j int) bool {
				return Number(Object(m.doc.NodeByID[nodes[i]]["span"])["start"]) < Number(Object(m.doc.NodeByID[nodes[j]]["span"])["start"])
			})
			x[field] = m.doc.Rich(nodes)
		}
	}
	return x
}

func (m *domainMapper) branchDetail(owner Record, s *Section) {
	switch {
	case s.Title == "Extensions":
		owner["branches"] = append(Strings(owner["branches"]), m.listBranches(owner, s.Intro)...)
		for _, sub := range s.Children {
			x := m.branch(owner, sub)
			appendID(owner, "branches", String(x, "id"))
		}
	case authoredLabel(explicitBranch, s.Title) != "":
		x := m.branch(owner, s)
		appendID(owner, "branches", String(x, "id"))
	default:
		m.extension(owner, s)
	}
}

func (m *domainMapper) listBranches(owner Record, nodes []string) []string {
	branches := []string{}
	for _, id := range m.listItems(nodes, false) {
		if authoredLabel(explicitBranch, m.doc.Text(id)) != "" {
			x := m.listBranch(owner, id)
			branches = append(branches, String(x, "id"))
		} else {
			// An error within a handling step belongs to the enclosing branch.
			branches = append(branches, m.listBranches(owner, Strings(m.doc.NodeByID[id]["children"]))...)
		}
	}
	return branches
}

func (m *domainMapper) listBranch(owner Record, id string) Record {
	x := m.record("UseCaseExtension", String(owner, "id"))
	x["authored_label"] = authoredLabel(explicitBranch, m.doc.Text(id))
	x["content"] = m.doc.Rich([]string{id})
	x["references"] = m.references([]string{id})
	m.bind(x, "content", []string{id}, "structured_item", "labelled-extension-item")
	children := Strings(m.doc.NodeByID[id]["children"])
	x["steps"] = m.scenarioSteps(x, children)
	x["branches"] = m.listBranches(x, children)
	return x
}
