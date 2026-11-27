package model

import (
	"regexp"
	"strconv"
)

func (m *domainMapper) useCase(r Record, s *Section) bool {
	switch s.Title {
	case "System boundary":
		m.section(r, "boundary_description", s)
		m.designField(r, "subject_system", s, "Subject system")
	case "Main success scenario":
		x := m.named("Scenario", r, s)
		m.section(x, "content", s)
		x["steps"] = m.scenarioSteps(x, s.Intro)
		r["main_success_scenario"] = String(x, "id")
		m.markHeading("main_success_scenario", s.HeadingID)
	case "Extensions":
		m.declared(r, "extensions", s)
		for _, sub := range s.Children {
			x := m.branch(r, sub)
			appendID(r, "extensions", String(x, "id"))
		}
	default:
		return false
	}
	return true
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
			label := strconv.Itoa(index)
			if match := orderedMarker.FindStringSubmatch(m.doc.Raw(item)); len(match) > 1 {
				label = match[1]
			}
			index++
			x := m.record("ScenarioStep", String(owner, "id"))
			x["authored_label"] = label
			x["content"] = m.doc.Rich([]string{item})
			m.stepIllustrates(x, item)
			m.markField(x, "content", []string{item})
			steps = append(steps, String(x, "id"))
		}
	}
	return steps
}

var orderedMarker = regexp.MustCompile(`^\s*([0-9]+)[.)]\s+`)

func (m *domainMapper) branch(owner Record, s *Section) Record {
	x := m.named("UseCaseExtension", owner, s)
	m.section(x, "content", s)
	for _, sub := range s.Children {
		if sub.Title == "Handling" {
			m.section(x, "handling", sub)
			x["steps"] = m.scenarioSteps(x, sub.Intro)
		} else {
			m.subjectSection(x, sub)
		}
	}
	return x
}
