package model

import (
	"regexp"
	"strconv"
	"strings"
)

func (m *domainMapper) useCase(r Record, s *Section) bool {
	fields := map[string]string{"Goal": "goal", "System boundary": "system_boundary", "Trigger": "trigger", "Success criterion": "success_criterion"}
	if field, ok := fields[s.Title]; ok {
		m.section(r, field, s)
		return true
	}
	switch s.Title {
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
		for _, sub := range s.Children {
			x := m.branch(r, sub)
			appendID(v, "branches", String(x, "id"))
		}
		for _, id := range m.listItems(s.Intro, false) {
			if explicitBranch.MatchString(m.doc.Text(id)) {
				x := m.listBranch(r, id)
				appendID(v, "branches", String(x, "id"))
			}
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
		if String(n, "node_kind") != "list" || !ordered {
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
			x := m.record("ScenarioStep", String(owner, "id"))
			x["authored_label"] = strconv.Itoa(index)
			if match := orderedMarker.FindStringSubmatch(m.doc.Raw(item)); len(match) > 1 {
				x["authored_label"] = match[1]
			}
			index++
			x["content"] = m.doc.Rich([]string{item})
			x["governing_references"] = m.references([]string{item})
			m.bind(x, "content", []string{item}, "structured_item", "ordered-scenario-step")
			// Labels are explicit Markdown fields inside this step, not guesses
			// about the grammatical actor/action in its prose sentence.
			s := &Section{Intro: Strings(m.doc.NodeByID[item]["children"])}
			for _, pair := range [][2]string{{"Action", "action"}, {"Response", "response"}, {"Outcome", "outcome"}} {
				if body := m.fieldNodes(s, pair[0]); body != nil {
					x[pair[1]] = m.doc.Rich(body)
					m.bind(x, pair[1], body, "named_field", "step-field/"+pair[0])
					m.bind(x, pair[1], m.fieldLabelNodes(s, pair[0]), "structural_wrapper", "step-field-label/"+pair[0])
				}
			}
			if body := m.fieldNodes(s, "Actor"); body != nil {
				x["actor"] = m.actor(body)
				m.bind(x, "actor.content", body, "named_field", "step-actor")
				m.bind(x, "actor.content", m.fieldLabelNodes(s, "Actor"), "structural_wrapper", "step-actor-label")
			}
			steps = append(steps, String(x, "id"))
		}
	}
	return steps
}

var orderedMarker = regexp.MustCompile(`^\s*([0-9]+)[.)]\s+`)
var explicitBranch = regexp.MustCompile(`^\s*(?:[-+*]|[0-9]+[.)])?\s*(?:\*\*)?(?:[0-9]+[a-z]|\*[a-z])(?:[.:)]|\s|\*\*)`)

func (m *domainMapper) branch(owner Record, s *Section) Record {
	x := m.named("Extension", owner, s)
	m.section(x, "content", s)
	x["governing_references"] = m.references(s.Content)
	if flow := m.scenarioSteps(x, s.Intro); len(flow) > 0 {
		x["flow"] = flow
	}
	for _, sub := range s.Children {
		field := map[string]string{"Condition": "condition", "Location": "location", "At step": "location", "Point of occurrence": "location", "Response": "response", "Continuation": "continuation"}[sub.Title]
		if field != "" {
			m.section(x, field, sub)
			continue
		}
		switch sub.Title {
		case "Unresolved handling":
			x["handling_state"] = "unresolved"
			m.sectionBindings(x, "handling_state", sub)
		case "Handling":
			value := strings.ToLower(strings.TrimSpace(m.contentText(sub.Content)))
			if value == "specified" || value == "unresolved" {
				x["handling_state"] = value
				m.sectionBindings(x, "handling_state", sub)
			}
		case "Flow":
			x["flow"] = append(Strings(x["flow"]), m.scenarioSteps(x, sub.Intro)...)
			for _, detail := range sub.Children {
				m.extension(x, detail)
			}
		case "Governing references":
			x["governing_references"] = appendUniqueReferences(Strings(x["governing_references"]), m.references(sub.Content))
		default:
			b := m.branch(x, sub)
			appendID(x, "branches", String(b, "id"))
		}
	}
	for _, pair := range [][2]string{{"Condition", "condition"}, {"At step", "location"}, {"Response", "response"}, {"Continuation", "continuation"}} {
		if nodes := m.fieldNodes(s, pair[0]); nodes != nil {
			m.setRich(x, pair[1], nodes)
			m.bind(x, pair[1], nodes, "named_field", "extension-field/"+pair[0])
			m.bind(x, pair[1], m.fieldLabelNodes(s, pair[0]), "structural_wrapper", "extension-field-label/"+pair[0])
		}
	}
	if handling := strings.ToLower(m.fieldText(s, "Handling")); handling == "specified" || handling == "unresolved" {
		x["handling_state"] = handling
		m.bind(x, "handling_state", m.fieldNodes(s, "Handling"), "named_field", "extension-handling")
	}
	for _, id := range m.listItems(s.Intro, false) {
		if explicitBranch.MatchString(m.doc.Text(id)) {
			b := m.listBranch(x, id)
			appendID(x, "branches", String(b, "id"))
		}
	}
	return x
}

func appendUniqueReferences(existing, more []string) []string {
	for _, value := range more {
		found := false
		for _, old := range existing {
			if old == value {
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, value)
		}
	}
	return existing
}

func (m *domainMapper) listBranch(owner Record, id string) Record {
	x := m.record("Extension", String(owner, "id"))
	x["content"] = m.doc.Rich([]string{id})
	x["governing_references"] = m.references([]string{id})
	m.bind(x, "content", []string{id}, "structured_item", "labelled-extension-item")
	children := Strings(m.doc.NodeByID[id]["children"])
	if flow := m.scenarioSteps(x, children); len(flow) > 0 {
		x["flow"] = flow
	}
	for _, sub := range m.listItems(children, false) {
		if explicitBranch.MatchString(m.doc.Text(sub)) {
			b := m.listBranch(x, sub)
			appendID(x, "branches", String(b, "id"))
		}
	}
	return x
}
