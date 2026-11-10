package system

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	"github.com/craigsmitham/dotspec/internal/specfile"
	"github.com/yuin/goldmark/ast"
)

// QualityModel is a system-local vocabulary. A nil model means no vocabulary
// was declared; an empty model is an explicit, unfinished vocabulary.
type QualityModel struct {
	Markdown   string      `json:"markdown"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Definition      string              `json:"definition"`
	Broader         []string            `json:"broader"`
	Location        diagnostic.Location `json:"location"`
	broaderLocation diagnostic.Location
}

var attributeField = regexp.MustCompile(`^\*\*(ID|Broader):\*\*\s*(.*?)\s*$`)

func parseQualityModel(doc *specfile.Outline, section *specfile.Section, ids *specfile.Identifiers) (*QualityModel, error) {
	model := &QualityModel{Markdown: doc.Body(section), Attributes: []Attribute{}}
	for _, child := range section.Children {
		if child.Level != 3 || child.Title == "" {
			return nil, diagnostic.Error("structure.attribute-heading", "each quality attribute must have a nonempty level-three heading (### Attribute name)", doc.Location(child.Start, child.Title))
		}
		attribute, err := parseAttribute(doc, child, ids)
		if err != nil {
			return nil, err
		}
		model.Attributes = append(model.Attributes, attribute)
	}
	if err := model.validate(); err != nil {
		return nil, err
	}
	return model, nil
}

func parseAttribute(doc *specfile.Outline, section *specfile.Section, ids *specfile.Identifiers) (Attribute, error) {
	a := Attribute{Name: section.Title, Broader: []string{}, Location: doc.Location(section.Start, section.Title)}
	fields := make(map[string]string)
	locations := make(map[string]diagnostic.Location)
	fieldLines := make(map[int]bool)
	for _, node := range doc.Nodes(section) {
		if _, ok := node.(*ast.Paragraph); !ok {
			continue
		}
		for i := 0; i < node.Lines().Len(); i++ {
			segment := node.Lines().At(i)
			match := attributeField.FindStringSubmatch(strings.TrimSpace(string(segment.Value(doc.Source))))
			if match == nil {
				continue
			}
			location := doc.Location(segment.Start, a.Name+" / "+match[1])
			if previous, exists := locations[match[1]]; exists {
				return a, diagnostic.Error("structure.duplicate-field", fmt.Sprintf("attribute %s: duplicate %s field", a.Name, match[1]), location, previous)
			}
			fields[match[1]], locations[match[1]] = match[2], location
			fieldLines[bytes.LastIndexByte(doc.Source[:segment.Start], '\n')+1] = true
		}
	}
	a.ID = fields["ID"]
	if a.ID == "" {
		return a, diagnostic.Error("identity.required", "each quality attribute needs a nonempty **ID:** field", a.Location)
	}
	if err := ids.Register(a.ID, "attribute "+a.Name, locations["ID"]); err != nil {
		return a, err
	}
	if value, exists := fields["Broader"]; exists {
		a.broaderLocation = locations["Broader"]
		var references specfile.Identifiers
		for _, id := range strings.Split(value, ",") {
			id = strings.TrimSpace(id)
			if id == "" {
				return a, diagnostic.Error("quality.broader", "Broader must contain a comma-separated list of nonempty attribute IDs; omit the field for a root attribute", a.broaderLocation)
			}
			if err := references.Register(id, "broader attribute", a.broaderLocation); err != nil {
				return a, err
			}
			a.Broader = append(a.Broader, id)
		}
	}
	// Remove metadata only from recognized paragraph lines. Preserve examples,
	// blockquotes, and all deeper Markdown sections as part of the definition.
	var definition []string
	offset := section.ContentStart
	for _, line := range strings.Split(string(doc.Source[section.ContentStart:section.End]), "\n") {
		if !fieldLines[offset] {
			definition = append(definition, line)
		}
		offset += len(line) + 1
	}
	a.Definition = strings.TrimSpace(strings.Join(definition, "\n"))
	return a, nil
}

func (m *QualityModel) validate() error {
	byID := make(map[string]*Attribute, len(m.Attributes))
	for i := range m.Attributes {
		byID[m.Attributes[i].ID] = &m.Attributes[i]
	}
	// Iterative topological traversal avoids recursion limits on authored graphs.
	remaining := make(map[string]int, len(byID))
	children := make(map[string][]string)
	var roots []string
	for _, a := range m.Attributes {
		remaining[a.ID] = len(a.Broader)
		if len(a.Broader) == 0 {
			roots = append(roots, a.ID)
		}
		for _, id := range a.Broader {
			if byID[id] == nil {
				return diagnostic.Error("quality.broader", fmt.Sprintf("attribute %q refers to undefined broader attribute %q", a.ID, id), a.broaderLocation)
			}
			children[id] = append(children[id], a.ID)
		}
	}
	for i := 0; i < len(roots); i++ {
		for _, child := range children[roots[i]] {
			remaining[child]--
			if remaining[child] == 0 {
				roots = append(roots, child)
			}
		}
	}
	if len(roots) != len(m.Attributes) {
		for _, a := range m.Attributes {
			if remaining[a.ID] != 0 {
				return diagnostic.Error("quality.cycle", fmt.Sprintf("attribute %q has a cyclic broader relationship or depends on one", a.ID), a.broaderLocation)
			}
		}
	}
	return nil
}

func (m *QualityModel) Find(id string) *Attribute {
	if m != nil {
		for i := range m.Attributes {
			if m.Attributes[i].ID == id {
				return &m.Attributes[i]
			}
		}
	}
	return nil
}

// Ancestors returns distinct broader IDs in declaration order. It does not
// turn those classifications into authored labels or inherited obligations.
func (m *QualityModel) Ancestors(ids []string) []string {
	seen := make(map[string]bool)
	queue := append([]string(nil), ids...)
	for _, id := range ids {
		seen[id] = true
	}
	ancestors := make(map[string]bool)
	for i := 0; i < len(queue); i++ {
		if a := m.Find(queue[i]); a != nil {
			for _, broader := range a.Broader {
				ancestors[broader] = true
				if !seen[broader] {
					seen[broader] = true
					queue = append(queue, broader)
				}
			}
		}
	}
	result := []string{}
	if m != nil {
		for _, a := range m.Attributes {
			if ancestors[a.ID] {
				result = append(result, a.ID)
			}
		}
	}
	return result
}

func Check(spec *System) []diagnostic.Diagnostic {
	var findings []diagnostic.Diagnostic
	if spec.QualityModel != nil {
		for _, a := range spec.QualityModel.Attributes {
			if a.Definition == "" {
				location := a.Location
				location.Path = spec.Path
				findings = append(findings, diagnostic.Advice("quality.definition", "Define what this quality attribute means for the system.", location))
			}
		}
	}
	return findings
}
