package model

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"go.yaml.in/yaml/v3"
)

// MapDomain applies the versioned Markdown mapping profile. Existing parsers
// remain authoritative for authoring errors. Mapping recognizes explicit
// headings, lists and labelled fields; it never infers assertions from prose.
func MapDomain(doc *ParsedDocument, kind string) Mapping {
	m := &domainMapper{doc: doc, result: Mapping{Elements: []Record{}, Labels: []Record{}, Diagnostics: []Record{}}}
	if kind == "glossary" {
		context, err := system.ParseGlossaryMarkdown(doc.Syntax)
		if err != nil {
			m.parseFailure(err)
			return m.result
		}
		if context != "" {
			doc.Document["glossary_context"] = strings.TrimSpace(context)
		}
		doc.Document["glossary_terms"] = []string{}
		m.glossaryTerms(doc.Document, "glossary_terms", doc.Sections)
		return m.result
	}
	if kind == "" {
		return m.result
	}
	var err error
	var authoredID string
	var glossarySources []string
	var labels []string
	switch kind {
	case "feature", "capability":
		var v *system.ProductSpec
		v, err = system.ParseProductSpecMarkdown(kind, doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "quality-attribute":
		var a *system.AttributeDocument
		a, err = system.ParseQualityAttributeMarkdown(String(doc.Source, "locator"), doc.Syntax)
		if err == nil {
			authoredID = a.ID
			labels = a.SubattributeOf
		}
	case "quality-model":
		var v *qualitymodel.QualityModel
		v, err = qualitymodel.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "business-rule":
		var v *businessrule.BusinessRule
		v, err = businessrule.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "component":
		var v *system.Component
		v, err = system.ParseComponentMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
			err = system.ComponentLayout.CheckName(String(doc.Source, "locator"), v.Name)
		}
	case "system":
		var v *system.System
		v, err = system.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
			glossarySources = v.GlossarySources
			err = system.Layout.CheckName(String(doc.Source, "locator"), v.Name)
		}
	case "job":
		var v *job.Job
		v, err = job.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "concept":
		var v *concept.Concept
		v, err = concept.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "use-case":
		var v *usecase.UseCase
		v, err = usecase.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID = v.ID
		}
	case "requirement":
		var v *requirement.Requirement
		v, err = requirement.ParseMarkdown(doc.Syntax)
		if err == nil {
			authoredID, labels = v.ID, v.QualityAttributes
		}
	default:
		err = fmt.Errorf("no domain mapping for specification type %q", kind)
	}
	if err != nil {
		m.parseFailure(err)
		return m.result
	}
	r := m.record(specfile.DocumentRootClass(kind), "")
	m.result.RootID = String(r, "id")
	if authoredID != "" {
		r["authored_id"] = authoredID
	}
	if kind != "system" {
		name := filepath.Base(String(doc.Source, "locator"))
		r["slug"] = strings.TrimSuffix(name, specfile.AuthoredSuffix(kind))
	}
	if doc.TitleNode != "" {
		r["name"] = doc.Rich([]string{doc.TitleNode})
	}
	if len(doc.Preamble) > 0 {
		r["preamble"] = doc.Rich(doc.Preamble)
		m.markField(r, "preamble", doc.Preamble)
	}
	m.metadata(r)
	if designKind(String(r, "kind"), "Specification") {
		m.companions(r)
	}
	if glossarySources != nil {
		r["glossary_source_paths"] = glossarySources
	}
	if kind == "quality-attribute" {
		for _, parent := range labels {
			node := m.metadataNode("subattribute_of", parent)
			m.label(r, parent, node, "subattribute_of")
		}
	}
	if kind == "requirement" {
		m.classification(r, labels)
	}
	for _, section := range doc.Sections {
		if section.Level != 2 {
			m.subjectSection(r, section)
			continue
		}
		if m.common(r, section) {
			continue
		}
		mapped := false
		switch kind {
		case "feature", "capability":
			switch {
			case section.Title == "Description" || section.Title == "Purpose" || (kind == "feature" && section.Title == "Overview") || (kind == "capability" && section.Title == "Enablers"):
				m.section(r, strings.ToLower(section.Title), section)
				mapped = true
			case kind == "feature" && section.Title == "Groups":
				m.structureReference(r, "groups", section)
				mapped = true
			}
		case "quality-attribute":
			if section.Title == "Definition" {
				m.section(r, "definition", section)
				mapped = true
			}
		case "quality-model":
			mapped = m.qualityModel(r, section)
		case "business-rule":
			mapped = m.businessRule(r, section)
		case "component":
			mapped = m.componentSection(r, section)
		case "system":
			mapped = m.system(r, section)
		case "job":
			mapped = m.job(r, section)
		case "concept":
			mapped = m.concept(r, section)
		case "use-case":
			mapped = m.useCase(r, section)
		case "requirement":
			mapped = m.requirement(r, section)
		}
		if !mapped {
			m.subjectSection(r, section)
		}
	}
	m.nestedTraceability()
	return m.result
}

type domainMapper struct {
	doc         *ParsedDocument
	result      Mapping
	next        int
	traceMapped map[string]bool
}

func (m *domainMapper) id(suffix string) string {
	m.next++
	return fmt.Sprintf("%s:%s:%d", String(m.doc.Document, "id"), suffix, m.next)
}

func (m *domainMapper) record(kind, owner string) Record {
	r := Record{"id": m.id("domain"), "kind": kind, "source_document": String(m.doc.Source, "id")}
	if owner != "" {
		r["owner"] = owner
	}
	m.result.Elements = append(m.result.Elements, r)
	return r
}

func (m *domainMapper) markField(owner Record, field string, nodes []string) {
	for _, id := range nodes {
		node := m.doc.NodeByID[id]
		node["field_owner"], node["field_path"] = owner["id"], field
	}
}

func (m *domainMapper) markHeading(field, id string) {
	m.doc.NodeByID[id]["domain_field"] = field
}

func (m *domainMapper) authoredFragment(owner Record, nodes []string) {
	if len(nodes) == 0 || owner["source_fragment"] != nil || String(owner, "authored_id") == "" {
		return
	}
	id := nodes[0]
	fragment := Record{"id": m.id("fragment"), "kind": "SourceFragment", "source": m.doc.Source["id"], "content_node": id, "authored_anchor": owner["authored_id"], "span": m.doc.NodeByID[id]["span"]}
	m.doc.Fragments = append(m.doc.Fragments, fragment)
	owner["source_fragment"] = fragment["id"]
}

func (m *domainMapper) section(owner Record, field string, s *Section) {
	m.setRich(owner, field, s.Content)
	m.markSection(owner, field, s)
}

func (m *domainMapper) setRich(owner Record, field string, nodes []string) {
	previous := Object(owner[field])
	if previous != nil {
		nodes = append(append([]string{}, Strings(previous["content_nodes"])...), nodes...)
	}
	owner[field] = m.doc.Rich(nodes)
}

func (m *domainMapper) markSection(owner Record, field string, s *Section) {
	m.markField(owner, field, s.Content)
	m.markHeading(field, s.HeadingID)
}

// A declaration collection contains named members. Its heading records presence.
func (m *domainMapper) declared(owner Record, field string, s *Section) {
	m.markHeading(field, s.HeadingID)
	for _, id := range s.Intro {
		if String(m.doc.NodeByID[id], "node_kind") != "trivia" {
			m.authoringFinding(owner, id, "structure.collection-content", fmt.Sprintf("%s contains named entries. Put this passage under its entry or move it to the field it qualifies.", s.Title))
			break
		}
	}
	if owner[field] == nil {
		owner[field] = []string{}
	}
}

func (m *domainMapper) named(kind string, owner Record, s *Section) Record {
	r := m.record(kind, String(owner, "id"))
	r["name"] = m.doc.Rich([]string{s.HeadingID})
	if declaredWithID[kind] {
		if id := m.fieldText(s, "ID"); id != "" {
			r["authored_id"] = id
			m.authoredFragment(r, m.fieldNodes(s, "ID"))
		}
	}
	return r
}

// Elements declared by name inside a specification may carry an explicit ID.
var declaredWithID = map[string]bool{"Component": true, "ConceptUse": true, "Synchronization": true, "Interface": true}

func appendID(r Record, field, id string) { r[field] = append(Strings(r[field]), id) }

func (m *domainMapper) parseFailure(err error) {
	d := Record{"id": m.id("diagnostic"), "kind": "Diagnostic", "code": "model.domain-parse", "severity": "error", "message": err.Error(), "source": String(m.doc.Source, "id"), "phase": "parse"}
	var finding *diagnostic.Diagnostic
	if errors.As(err, &finding) {
		d["code"], d["severity"], d["message"] = finding.Rule, finding.Severity, finding.Message
		if finding.Location.Line > 0 {
			d["line"] = finding.Location.Line
		}
		if finding.Location.Column > 0 {
			d["column"] = finding.Location.Column
		}
		if finding.Location.Element != "" {
			d["element"] = finding.Location.Element
		}
		if len(finding.Related) > 0 {
			related := []Record{}
			for _, location := range finding.Related {
				source := String(m.doc.Source, "id")
				if location.Path != "" && filepath.Clean(location.Path) != filepath.Clean(String(m.doc.Source, "locator")) {
					source = ""
				}
				related = append(related, diagnosticSourceLocation(location, source))
			}
			d["related"] = related
		}
	}
	m.result.Diagnostics = append(m.result.Diagnostics, d)

}

// diagnosticSourceLocation preserves reported coordinates. A missing line or
// column stays absent; a locator is retained only when the producer supplied it.
func diagnosticSourceLocation(location diagnostic.Location, source string) Record {
	r := Record{}
	if source != "" {
		r["source"] = source
	}
	if location.Path != "" {
		r["locator"] = location.Path
	}
	if location.Line > 0 {
		r["line"] = location.Line
	}
	if location.Column > 0 {
		r["column"] = location.Column
	}
	if location.Element != "" {
		r["element"] = location.Element
	}
	return r
}

func (m *domainMapper) subjectSection(owner Record, s *Section) {
	if designKind(String(owner, "kind"), "Specification") && m.common(owner, s) {
		return
	}
	if m.traceSection(owner, s) || m.boundarySection(owner, s) {
		return
	}
	// Any domain element can carry explicitly headed open questions.
	if kind := String(owner, "kind"); openQuestionHeading(s.Title) && kind != "Document" {
		m.openQuestions(owner, []*Section{s})
		return
	}
	m.authoringFinding(owner, s.HeadingID, "structure.unknown-section", fmt.Sprintf("Unrecognized section %q. Move its prose into the field it qualifies or use a declared section.", s.Title))
}

func (m *domainMapper) authoringFinding(owner Record, node, code, message string) {
	span := Object(m.doc.NodeByID[node]["span"])
	start := Number(span["start"])
	m.result.Diagnostics = append(m.result.Diagnostics, Record{
		"id": m.id("diagnostic"), "kind": "Diagnostic", "severity": "error", "code": code,
		"source": m.doc.Source["id"], "record": owner["id"], "span": span, "phase": "parse",
		"line": 1 + bytes.Count(m.doc.Data[:start], []byte("\n")), "message": message,
	})
}

func (m *domainMapper) common(r Record, s *Section) bool {
	if m.traceSection(r, s) || m.boundarySection(r, s) {
		return true
	}
	switch s.Title {
	case "Adopts":
		if !designKind(String(r, "kind"), "Specification") {
			return false
		}
		m.adoptions(r, s)
	case "Specification system":
		if _, ok := compiledModel().catalog.Classes[String(r, "kind")].Slots["specification_system_content"]; !ok {
			return false
		}
		m.structureReference(r, "specification_system", s)
	case "Related specifications":
		m.contextField(r, "related_specifications", s.Content)
		m.markSection(r, "related_specifications.content", s)
	case "Authority":
		m.section(r, "authority_and_evidence", s)
	case "Evidence":
		m.section(r, "evidence", s)
	case "Open questions":
		m.openQuestions(r, []*Section{s})
	default:
		return false
	}
	return true
}

func (m *domainMapper) contextField(owner Record, field string, nodes []string) {
	if previous := Object(owner[field]); previous != nil {
		nodes = append(append([]string{}, Strings(Object(previous["content"])["content_nodes"])...), nodes...)
	}
	owner[field] = m.context(nodes)
}

func (m *domainMapper) context(nodes []string) Record {
	return Record{"content": m.doc.Rich(nodes), "references": m.references(nodes)}
}

func (m *domainMapper) descendants(nodes []string) map[string]bool {
	seen := map[string]bool{}
	var visit func(string)
	visit = func(id string) {
		if seen[id] {
			return
		}
		seen[id] = true
		for _, c := range Strings(m.doc.NodeByID[id]["children"]) {
			visit(c)
		}
	}
	for _, id := range nodes {
		visit(id)
	}
	return seen
}

func (m *domainMapper) references(nodes []string) []string {
	set := m.descendants(nodes)
	ids := []string{}
	for _, ref := range m.doc.References {
		if set[String(ref, "content_node")] {
			ids = append(ids, String(ref, "id"))
		}
	}
	return ids
}

// listItems selects direct items of lists at the field's top level. Nested
// bullet prose belongs to its item; it is not another outcome or interaction.
func (m *domainMapper) listItems(nodes []string, orderedOnly bool) []string {
	ids := []string{}
	for _, id := range nodes {
		n := m.doc.NodeByID[id]
		if String(n, "node_kind") != "list" {
			continue
		}
		ordered, _ := n["list_ordered"].(bool)
		if orderedOnly && !ordered {
			continue
		}
		for _, child := range Strings(n["children"]) {
			if String(m.doc.NodeByID[child], "node_kind") == "list_item" {
				ids = append(ids, child)
			}
		}
	}
	return ids
}

func (m *domainMapper) fieldNodes(s *Section, label string) []string {
	needle := "**" + label + ":**"
	for _, id := range s.Intro {
		if String(m.doc.NodeByID[id], "node_kind") != "paragraph" || !strings.Contains(m.doc.Text(id), needle) {
			continue
		}
		children := Strings(m.doc.NodeByID[id]["children"])
		for i, child := range children {
			if String(m.doc.NodeByID[child], "node_kind") != "strong" || strings.TrimSpace(m.doc.Text(child)) != needle {
				continue
			}
			values := []string{}
			valueStart := Number(Object(m.doc.NodeByID[child]["span"])["end"])
			lineEnd := len(m.doc.Data)
			if offset := bytes.IndexByte(m.doc.Data[valueStart:], '\n'); offset >= 0 {
				lineEnd = valueStart + offset
			}
			for _, value := range children[i+1:] {
				n := m.doc.NodeByID[value]
				span := Object(n["span"])
				kind := String(n, "node_kind")
				if Number(span["start"]) >= lineEnd || kind == "soft_break" || kind == "hard_break" {
					break
				}
				if Number(span["end"]) > lineEnd {
					continue
				}
				if (kind == "text" || kind == "trivia") && strings.TrimSpace(m.doc.Text(value)) == "" {
					continue
				}
				values = append(values, value)
			}
			return values
		}
		return []string{id}
	}
	return nil
}

func (m *domainMapper) fieldText(s *Section, label string) string {
	for _, id := range s.Intro {
		if String(m.doc.NodeByID[id], "node_kind") != "paragraph" {
			continue
		}
		for _, line := range strings.Split(m.doc.Text(id), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "**"+label+":**") {
				return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "**"+label+":**"))
			}
		}
	}
	return ""
}

func (m *domainMapper) withoutFields(s *Section, nodes []string, labels ...string) []string {
	intro := map[string]bool{}
	for _, id := range s.Intro {
		intro[id] = true
	}
	result := []string{}
	for _, id := range nodes {
		n := m.doc.NodeByID[id]
		if !intro[id] || String(n, "node_kind") != "paragraph" {
			result = append(result, id)
			continue
		}
		ranges := [][2]int{}
		offset := Number(Object(n["span"])["start"])
		for _, line := range strings.SplitAfter(m.doc.Raw(id), "\n") {
			trimmed := strings.TrimSpace(line)
			for _, label := range labels {
				if strings.HasPrefix(trimmed, "**"+label+":**") {
					ranges = append(ranges, [2]int{offset, offset + len(line)})
					break
				}
			}
			offset += len(line)
		}
		if len(ranges) == 0 {
			result = append(result, id)
			continue
		}
		for _, child := range Strings(n["children"]) {
			span := Object(m.doc.NodeByID[child]["span"])
			a, b := Number(span["start"]), Number(span["end"])
			remove := false
			for _, line := range ranges {
				if a >= line[0] && b <= line[1] {
					remove = true
					break
				}
			}
			if !remove {
				result = append(result, child)
			}
		}
	}
	return result
}

func (m *domainMapper) metadata(r Record) {
	if String(r, "authored_id") != "" {
		if node := m.metadataNode("id", String(r, "authored_id")); node != "" {
			m.authoredFragment(r, []string{node})
		}
	}
}

func (m *domainMapper) metadataNode(key, value string) string {
	root := m.doc.Syntax.Value(key)
	if root == nil {
		return ""
	}
	target := yamlScalarTarget(root)
	if target == nil {
		return ""
	}
	nodes := []*yaml.Node{root}
	if target.Kind == yaml.SequenceNode {
		nodes = target.Content
	}
	for _, n := range nodes {
		if scalar := yamlScalarTarget(n); scalar == nil || scalar.Value != value {
			continue
		}
		if root.Kind == yaml.AliasNode {
			n = root
		}
		return m.yamlContent(n, target.Style&yaml.FlowStyle != 0, value)
	}
	return ""
}

// yamlContent records a frontmatter value as a content node spanning its source.
func (m *domainMapper) yamlContent(n *yaml.Node, flow bool, value string) string {
	e := contentExtractor{p: m.doc, buffer: m.doc.Data}
	start := e.yamlPosition(n, e.lineEnd(0))
	end := e.yamlEnd(n, start, len(m.doc.Data)-len(m.doc.Syntax.Body), flow)
	id := fmt.Sprintf("%s:metadata:%d", String(m.doc.Source, "id"), start)
	if m.doc.NodeByID[id] == nil {
		node := Record{"id": id, "span": e.sourceSpan(start, end), "node_kind": "metadata", "decoded_text": value}
		m.doc.NodeByID[id] = node
		m.doc.Nodes = append(m.doc.Nodes, node)
	}
	return id
}

// companions maps each validated frontmatter entry to an owned Companion whose
// target is an ordinary source reference, resolved like any authored link.
func (m *domainMapper) companions(r Record) {
	list := m.doc.Syntax.Value(specfile.CompanionsField)
	if list == nil || list.Kind != yaml.SequenceNode {
		return
	}
	r["companions"] = []string{}
	for _, item := range list.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}
		flow := item.Style&yaml.FlowStyle != 0 || list.Style&yaml.FlowStyle != 0
		values := map[string]*yaml.Node{}
		for i := 0; i+1 < len(item.Content); i += 2 {
			values[item.Content[i].Value] = item.Content[i+1]
		}
		target, role := values["target"], values["role"]
		if target == nil || role == nil {
			continue
		}
		companion := m.record("Companion", String(r, "id"))
		companion["role"] = role.Value
		if note := values["note"]; note != nil && note.Value != "" {
			companion["note"] = note.Value
		}
		node := m.yamlContent(target, flow, target.Value)
		ref := Record{"id": fmt.Sprintf("reference:%s:companion:%d", String(m.doc.Source, "id"), len(m.doc.References)), "kind": "SourceReference", "source_document": m.doc.Source["id"], "content_node": node, "span": m.doc.NodeByID[node]["span"], "form": "frontmatter", "authored_target": target.Value}
		m.doc.References = append(m.doc.References, ref)
		companion["target"] = ref["id"]
		m.bindReferenceOccurrences(companion, "target", []string{String(ref, "id")})
		appendID(r, "companions", String(companion, "id"))
	}
}

func (m *domainMapper) classification(r Record, labels []string) {
	if m.doc.Syntax.Value("quality_attributes") == nil {
		return
	}

	for _, label := range labels {
		node := m.metadataNode("quality_attributes", label)
		if node == "" {
			continue
		}
		m.label(r, label, node, "quality_classification")
	}
}

func (m *domainMapper) label(owner Record, label, node, role string) string {
	r := Record{"id": m.id("label"), "kind": "LabelOccurrence", "owner": String(owner, "id"), "role": role, "label": label, "span": m.doc.NodeByID[node]["span"], "content_node": node}
	m.result.Labels = append(m.result.Labels, r)
	return String(r, "id")
}
