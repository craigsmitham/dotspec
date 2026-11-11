// Package model represents complete, source-backed specification snapshots.
// Heterogeneous records retain all fields; the canonical schema supplies their
// closed shapes and typed references instead of a lossy base-struct projection.
package model

import (
	"encoding/json"

	"go.yaml.in/yaml/v3"
)

const ModelVersion = "0.2.0-draft"
const RepresentationVersion = "0.1.0"
const MappingProfile = "dotspec-markdown/1"

// Record is a schema-validated tagged record or an embedded model value.
// Values remain JSON-compatible and are never decoded through a base class
// that would discard concrete subtype fields.
type Record map[string]any

type Snapshot struct {
	ModelVersion          string   `json:"model_version"`
	RepresentationVersion string   `json:"representation_version"`
	MappingProfile        string   `json:"mapping_profile"`
	Scope                 Record   `json:"scope"`
	Sources               []Record `json:"sources"`
	Documents             []Record `json:"documents"`
	ContentNodes          []Record `json:"content_nodes"`
	Fragments             []Record `json:"fragments,omitempty"`
	Elements              []Record `json:"elements"`
	References            []Record `json:"references,omitempty"`
	LabelOccurrences      []Record `json:"label_occurrences,omitempty"`
	Resolutions           []Record `json:"resolutions,omitempty"`
	Assertions            []Record `json:"assertions,omitempty"`
	ContentBindings       []Record `json:"content_bindings"`
	Diagnostics           []Record `json:"diagnostics,omitempty"`
	ContextSources        []Record `json:"context_sources,omitempty"`
	Coverage              Record   `json:"coverage"`
}

func NewSnapshot(scope Record) *Snapshot {
	return &Snapshot{
		ModelVersion: ModelVersion, RepresentationVersion: RepresentationVersion,
		MappingProfile: MappingProfile, Scope: scope,
		Sources: []Record{}, Documents: []Record{}, ContentNodes: []Record{},
		Elements: []Record{}, ContentBindings: []Record{},
		Coverage: Record{"status": "not_checked", "mapping_profile": MappingProfile},
	}
}

// ParsedDocument is the extraction work buffer, not an interchange type.
// Its records use original-byte spans; sections expose the structure needed by
// the explicit domain mapping without reparsing narrative as natural language.
type ParsedDocument struct {
	Data        []byte
	Source      Record
	Document    Record
	Nodes       []Record
	NodeByID    map[string]Record
	Fragments   []Record
	References  []Record
	Diagnostics []Record
	Metadata    *yaml.Node
	Title       string
	TitleNode   string
	Preamble    []string
	Sections    []*Section
}

type Section struct {
	ID        string
	HeadingID string
	Title     string
	Level     int
	Content   []string // Ordered body roots, including child sections.
	Intro     []string // Body roots preceding the first child section.
	Children  []*Section
}

type Mapping struct {
	RootID      string
	Elements    []Record
	Bindings    []Record
	Labels      []Record
	Diagnostics []Record
}

func String(r Record, key string) string {
	v, _ := r[key].(string)
	return v
}

func Strings(value any) []string {
	switch values := value.(type) {
	case []string:
		return values
	case []any:
		result := make([]string, 0, len(values))
		for _, value := range values {
			if text, ok := value.(string); ok {
				result = append(result, text)
			}
		}
		return result
	}
	return nil
}

func Object(value any) Record {
	switch value := value.(type) {
	case Record:
		return value
	case map[string]any:
		return Record(value)
	}
	return nil
}

func Number(value any) int {
	switch value := value.(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	case json.Number:
		n, _ := value.Int64()
		return int(n)
	}
	return 0
}
