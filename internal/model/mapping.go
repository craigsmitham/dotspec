package model

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/businessrule"
	"github.com/craigsmitham/dotspec/internal/concept"
	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/job"
	"github.com/craigsmitham/dotspec/internal/quality"
	"github.com/craigsmitham/dotspec/internal/qualitymodel"
	"github.com/craigsmitham/dotspec/internal/requirement"
	"github.com/craigsmitham/dotspec/internal/system"
	"github.com/craigsmitham/dotspec/internal/usecase"
	"go.yaml.in/yaml/v3"
)

// MapDomain applies the versioned Markdown mapping profile. Existing parsers
// remain authoritative for authoring errors. Mapping recognizes explicit
// headings, lists and labelled fields; it never infers assertions from prose.
func MapDomain(doc *ParsedDocument, kind string) Mapping {
	m := &domainMapper{doc: doc, result: Mapping{Elements: []Record{}, Bindings: []Record{}, Labels: []Record{}, Diagnostics: []Record{}}}
	if kind == "" {
		// Ordinary referenced Markdown has a content model, not a fabricated
		// specification type. Unknown sections remain visible domain gaps.
		for _, s := range doc.Sections {
			m.extension(doc.Document, s)
		}
		doc.Document["domain_decomposition"] = "not_attempted"
		if root := String(doc.Document, "root"); root != "" {
			m.result.Bindings = append(m.result.Bindings, Record{"id": m.id("binding"), "kind": "ContentBinding", "content_node": root, "scope": "subtree", "owner": String(doc.Document, "id"), "field_path": "root", "role": "fallback", "mapping_rule": MappingProfile + "/ordinary-markdown", "reason": "Referenced Markdown is preserved without applying a dotspec domain authoring profile."})
		}
		return m.result
	}
	var err error
	var authoredID string
	var labels []string
	switch kind {
	case "quality-model":
		var v *qualitymodel.QualityModel
		v, err = qualitymodel.Parse(doc.Data)
		if err == nil {
			authoredID = v.ID
		}
	case "business-rule":
		var v *businessrule.BusinessRule
		v, err = businessrule.Parse(doc.Data)
		if err == nil {
			authoredID = v.ID
		}
	case "system":
		_, err = system.Parse(doc.Data)
	case "job":
		var v *job.Job
		v, err = job.Parse(doc.Data)
		if err == nil {
			authoredID = v.ID
		}
	case "concept":
		var v *concept.Concept
		v, err = concept.Parse(doc.Data)
		if err == nil {
			authoredID = v.ID
		}
	case "use-case":
		var v *usecase.UseCase
		v, err = usecase.Parse(doc.Data)
		if err == nil {
			authoredID = v.ID
		}
	case "requirement":
		var v *requirement.Requirement
		v, err = requirement.Parse(doc.Data)
		if err == nil {
			authoredID, labels = v.ID, v.QualityAttributes
		}
	case "quality":
		var v *quality.Requirement
		v, err = quality.Parse(doc.Data)
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
	class := map[string]string{"quality-model": "QualityModel", "business-rule": "BusinessRule", "system": "System", "job": "Job", "concept": "Concept", "use-case": "UseCase", "requirement": "GeneralRequirement", "quality": "QualityRequirement"}[kind]
	r := m.record(class, "")
	m.result.RootID = String(r, "id")
	if authoredID != "" {
		r["authored_id"] = authoredID
	}
	if kind != "system" {
		name := filepath.Base(String(doc.Source, "locator"))
		r["slug"] = strings.TrimSuffix(name, "."+kind+".spec.md")
	}
	if doc.TitleNode != "" {
		r["name"] = doc.Rich([]string{doc.TitleNode})
		m.bind(r, "name", []string{doc.TitleNode}, "record_identity", "title")
	}
	if len(doc.Preamble) > 0 {
		r["preamble"] = doc.Rich(doc.Preamble)
		m.bind(r, "preamble", doc.Preamble, "named_field", "preamble")
	}
	m.metadata(r)
	if kind == "quality" || kind == "requirement" {
		m.classification(r, labels)
	}
	for _, section := range doc.Sections {
		if section.Level != 2 {
			m.extension(r, section)
			continue
		}
		if m.common(r, section) {
			continue
		}
		mapped := false
		switch kind {
		case "quality-model":
			mapped = m.qualityModel(r, section)
		case "business-rule":
			mapped = m.businessRule(r, section)
		case "system":
			mapped = m.system(r, section)
		case "job":
			mapped = m.job(r, section)
		case "concept":
			mapped = m.concept(r, section)
		case "use-case":
			mapped = m.useCase(r, section)
		case "requirement", "quality":
			mapped = m.requirement(r, section)
		}
		if !mapped {
			m.extension(r, section)
		}
	}
	m.nestedTraceability()
	doc.Document["domain_decomposition"] = "complete"
	m.authoredFragments()
	m.bind(doc.Document, "root", []string{String(doc.Document, "root")}, "structural_wrapper", "document-root")
	for _, node := range doc.Nodes {
		if String(node, "node_kind") == "trivia" {
			m.bind(doc.Document, "root", []string{String(node, "id")}, "structural_wrapper", "source-trivia")
		}
	}
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
	r := Record{"id": m.id("domain"), "kind": kind, "source_document": String(m.doc.Source, "id"), "decomposition": "unsegmented"}
	if owner != "" {
		r["owner"] = owner
	}
	m.result.Elements = append(m.result.Elements, r)
	return r
}

func (m *domainMapper) bind(owner Record, field string, nodes []string, role, rule string) {
	for _, node := range nodes {
		if node == "" {
			continue
		}
		scope := "subtree"
		if role == "structural_wrapper" && String(m.doc.NodeByID[node], "node_kind") != "heading" {
			scope = "node"
		}
		m.result.Bindings = append(m.result.Bindings, Record{"id": m.id("binding"), "kind": "ContentBinding", "content_node": node, "scope": scope, "owner": String(owner, "id"), "field_path": field, "role": role, "mapping_rule": MappingProfile + "/" + rule})
	}
}

func (m *domainMapper) section(owner Record, field string, s *Section) {
	m.setRich(owner, field, s.Content)
	m.sectionBindings(owner, field, s)
}

func (m *domainMapper) setRich(owner Record, field string, nodes []string) {
	previous := Object(owner[field])
	if previous != nil {
		nodes = append(append([]string{}, Strings(previous["content_nodes"])...), nodes...)
	}
	owner[field] = m.doc.Rich(nodes)
}

func (m *domainMapper) sectionBindings(owner Record, field string, s *Section) {
	m.bind(owner, field, s.Content, "named_field", "section/"+s.Title)
	m.bind(owner, field, []string{s.HeadingID}, "structural_wrapper", "section-heading/"+s.Title)
}

func (m *domainMapper) collection(owner Record, field, items string, s *Section) Record {
	v := Object(owner[field])
	if v == nil {
		v = Record{items: []string{}}
	}
	m.setRich(v, "content", s.Content)
	owner[field] = v
	m.sectionBindings(owner, field+".content", s)
	return v
}

func (m *domainMapper) named(kind string, owner Record, s *Section) Record {
	r := m.record(kind, String(owner, "id"))
	r["name"] = m.doc.Rich([]string{s.HeadingID})
	m.bind(r, "name", []string{s.HeadingID}, "record_identity", "item-heading")
	return r
}

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
	m.doc.Document["domain_decomposition"] = "failed"
	root := String(m.doc.Document, "root")
	if root != "" {
		b := Record{"id": m.id("binding"), "kind": "ContentBinding", "content_node": root, "scope": "subtree", "owner": String(m.doc.Document, "id"), "field_path": "root", "role": "fallback", "mapping_rule": MappingProfile + "/invalid-domain", "reason": "Authoring parser rejected the document; source and syntax are retained without a valid domain object."}
		m.result.Bindings = append(m.result.Bindings, b)
	}
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

func (m *domainMapper) extension(owner Record, s *Section) Record {
	if m.traceSection(owner, s) {
		return nil
	}
	r := m.record("ExtensionSection", String(owner, "id"))
	r["heading"] = m.doc.Rich([]string{s.HeadingID})
	r["content"] = m.doc.Rich(s.Content)
	r["reason"] = "No named domain role is defined for this section in " + MappingProfile
	m.bind(r, "heading", []string{s.HeadingID}, "extension", "unknown-section")
	m.bind(r, "content", s.Content, "extension", "unknown-section")
	for _, b := range m.result.Bindings {
		if String(b, "owner") == String(r, "id") {
			b["reason"] = r["reason"]
		}
	}
	if String(owner, "kind") == "Document" {
		appendID(owner, "elements", String(r, "id"))
	} else if String(owner, "kind") == "ExtensionSection" {
		appendID(owner, "children", String(r, "id"))
	} else {
		appendID(owner, "extension_sections", String(r, "id"))
	}
	for _, child := range s.Children {
		m.extension(r, child)
	}
	return r
}

func (m *domainMapper) authoredFragments() {
	for _, element := range m.result.Elements {
		anchor := String(element, "authored_id")
		if anchor == "" {
			continue
		}
		var node string
		for _, binding := range m.result.Bindings {
			if String(binding, "owner") == String(element, "id") && String(binding, "field_path") == "authored_id" {
				node = String(binding, "content_node")
				break
			}
		}
		if node == "" {
			continue
		}
		fragment := Record{"id": m.id("fragment"), "kind": "SourceFragment", "source": String(m.doc.Source, "id"), "content_node": node, "authored_anchor": anchor, "anchor_rule": MappingProfile + "/authored-id", "span": m.doc.NodeByID[node]["span"]}
		m.doc.Fragments = append(m.doc.Fragments, fragment)
		element["source_fragment"] = String(fragment, "id")
	}
}

func (m *domainMapper) common(r Record, s *Section) bool {
	if m.traceSection(r, s) {
		return true
	}
	switch s.Title {
	case "Business rules":
		m.contextField(r, "business_rules", s.Content)
		m.sectionBindings(r, "business_rules.content", s)
	case "Related specifications":
		m.contextField(r, "related_specifications", s.Content)
		m.sectionBindings(r, "related_specifications.content", s)
	case "Related use cases":
		m.contextField(r, "related_use_cases", s.Content)
		m.sectionBindings(r, "related_use_cases.content", s)
	case "Related concepts":
		m.contextField(r, "related_concepts", s.Content)
		m.sectionBindings(r, "related_concepts.content", s)
	case "Authority", "Authority and evidence":
		m.section(r, "authority_and_evidence", s)
	case "Notes":
		m.section(r, "notes", s)
	case "Evidence":
		m.evidence(r, "evidence", s)
	case "Open question", "Open questions":
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
	v := Record{"content": m.doc.Rich(nodes), "references": m.references(nodes)}
	items := []Record{}
	for _, id := range m.listItems(nodes, false) {
		items = append(items, Record{"description": m.doc.Rich([]string{id}), "references": m.references([]string{id})})
	}
	if len(items) > 0 {
		v["items"] = items
	}
	return v
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

func (m *domainMapper) fieldLabelNodes(s *Section, label string) []string {
	needle := "**" + label + ":**"
	for _, id := range s.Intro {
		if String(m.doc.NodeByID[id], "node_kind") != "paragraph" {
			continue
		}
		for _, child := range Strings(m.doc.NodeByID[id]["children"]) {
			if String(m.doc.NodeByID[child], "node_kind") == "strong" && strings.TrimSpace(m.doc.Text(child)) == needle {
				return []string{child}
			}
		}
	}
	return nil
}

// withoutFields removes only explicit introductory metadata lines. Inline
// nodes remain rooted in the document tree; no rewritten prose or fabricated
// source ranges are introduced. Deeper examples with the same text are kept.
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
	if root := String(m.doc.Document, "metadata_root"); root != "" {
		m.bind(m.doc.Document, "metadata_root", []string{root}, "metadata", "frontmatter")
	}
	for _, id := range Strings(m.doc.Document["frontmatter"]) {
		m.bind(m.doc.Document, "frontmatter", []string{id}, "metadata", "frontmatter")
	}
	if String(r, "authored_id") != "" {
		if node := m.metadataNode("id", String(r, "authored_id")); node != "" {
			m.bind(r, "authored_id", []string{node}, "record_identity", "authored-id")
		}
	}
}

func (m *domainMapper) metadataNode(key, value string) string {
	// Prefer the YAML subtree under the named key so equal strings elsewhere do
	// not accidentally become the source of a semantic label.
	for _, n := range m.doc.Nodes {
		if String(n, "node_kind") != "yaml_entry" {
			continue
		}
		children := Strings(n["children"])
		if len(children) < 2 || String(m.doc.NodeByID[children[0]], "scalar_string") != key {
			continue
		}
		set := m.descendants(children[1:])
		for _, candidate := range m.doc.Nodes {
			if set[String(candidate, "id")] && String(candidate, "node_kind") == "yaml_scalar" && String(candidate, "scalar_string") == value {
				return String(candidate, "id")
			}
		}
		if authored := yamlValue(m.doc.Metadata, key); authored != nil && authored.Kind == yaml.AliasNode {
			for _, candidate := range m.doc.Nodes {
				if set[String(candidate, "id")] && String(candidate, "node_kind") == "yaml_alias" && String(candidate, "alias_name") == authored.Value {
					return String(candidate, "id")
				}
			}
		}
	}
	return ""
}

func (m *domainMapper) classification(r Record, labels []string) {
	if yamlValue(m.doc.Metadata, "quality_attributes") == nil {
		return
	}
	r["quality_attribute_labels"] = []string{}
	for _, label := range labels {
		node := m.metadataNode("quality_attributes", label)
		if node == "" {
			continue
		}
		id := m.label(r, label, node, "quality_classification")
		appendID(r, "quality_attribute_labels", id)
		m.bind(r, "quality_attribute_labels", []string{node}, "named_field", "quality-label")
	}
}

func (m *domainMapper) label(owner Record, label, node, role string) string {
	r := Record{"id": m.id("label"), "kind": "LabelOccurrence", "owner": String(owner, "id"), "role": role, "label": label, "content_node": node}
	m.result.Labels = append(m.result.Labels, r)
	return String(r, "id")
}

func yamlValue(root *yaml.Node, key string) *yaml.Node {
	if root == nil {
		return nil
	}
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		root = root.Content[0]
	}
	if root.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == key {
			return root.Content[i+1]
		}
	}
	return nil
}

var outcomeID = regexp.MustCompile(`^\s*(?:[-+*]|[0-9]+[.)])?\s*\*\*([^*]+):\*\*`)
