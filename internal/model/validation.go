package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	modelschema "github.com/craigsmitham/dotspec/model"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Finding struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
	Source  string `json:"source,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}

type slotDefinition struct {
	Range           string   `json:"range"`
	AnyOf           []string `json:"any_of"`
	ReferenceRanges []string `json:"reference_ranges"`
	InlineRanges    []string `json:"inline_ranges"`
	Required        bool     `json:"required"`
	Multivalued     bool     `json:"multivalued"`
	Inlined         bool     `json:"inlined"`
	InlinedAsList   bool     `json:"inlined_as_list"`
	Ordered         bool     `json:"ordered"`
}

type classDefinition struct {
	Abstract    bool                      `json:"abstract"`
	Ancestors   []string                  `json:"ancestors"`
	Slots       map[string]slotDefinition `json:"slots"`
	Annotations map[string]string         `json:"annotations"`
}

type modelCatalog struct {
	Classes map[string]classDefinition `json:"classes"`
}

type modelRuntime struct {
	classes map[string]*jsonschema.Schema
	schema  *jsonschema.Schema
	catalog modelCatalog
	err     error
}

var compiledModel = sync.OnceValue(func() modelRuntime {
	var catalog modelCatalog
	if err := json.Unmarshal(modelschema.CatalogJSON, &catalog); err != nil {
		return modelRuntime{err: fmt.Errorf("read model catalog: %w", err)}
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(modelschema.SchemaJSON))
	if err != nil {
		return modelRuntime{err: err}
	}
	compiler := jsonschema.NewCompiler()
	const address = "https://dotspec.invalid/snapshot.schema.json"
	if err := compiler.AddResource(address, doc); err != nil {
		return modelRuntime{err: err}
	}
	schema, err := compiler.Compile(address)
	if err != nil {
		return modelRuntime{err: err}
	}
	classes := map[string]*jsonschema.Schema{}
	for name := range catalog.Classes {
		concrete, err := compiler.Compile(address + "#/$defs/" + escapePointer(name))
		if err != nil {
			return modelRuntime{err: fmt.Errorf("compile model class %s: %w", name, err)}
		}
		classes[name] = concrete
	}
	return modelRuntime{schema: schema, catalog: catalog, classes: classes}
})

// ValidateJSON checks the complete interchange shape and graph contract. It
// never fetches references, mutates source files, or treats source diagnostics
// as defects in a faithful snapshot of those diagnostics.
func ValidateJSON(data []byte) ([]Finding, error) {
	runtime := compiledModel()
	if runtime.err != nil {
		return nil, runtime.err
	}
	value, err := decodeSnapshotJSON(data)
	if err != nil {
		return []Finding{{Code: "shape.json", Path: "/", Message: err.Error()}}, nil
	}
	if err := runtime.schema.Validate(value); err != nil {
		return shapeFindings(runtime, value, err), nil
	}
	v := &graphValidator{catalog: runtime.catalog, root: Object(value), records: map[string]locatedRecord{}, findings: []Finding{}}
	v.walk(Object(value), "Snapshot", "", true)
	v.walk(Object(value), "Snapshot", "", false)
	v.checkContent()
	v.checkDomain()
	sort.Slice(v.findings, func(i, j int) bool {
		a, b := v.findings[i], v.findings[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Code != b.Code {
			return a.Code < b.Code
		}
		return a.Message < b.Message
	})
	return v.findings, nil
}

func Validate(snapshot *Snapshot) ([]Finding, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}
	return ValidateJSON(data)
}

type locatedRecord struct {
	value Record
	path  string
}

type graphValidator struct {
	catalog        modelCatalog
	root           Record
	records        map[string]locatedRecord
	findings       []Finding
	orderedRecords []locatedRecord
}

func (v *graphValidator) add(code, path, message string) {
	f := Finding{Code: code, Path: path, Message: message}
	// The most specific containing record makes diagnostics deterministic.
	if record, ok := v.containingRecord(path); ok {
		if source, exists := v.records[v.sourceOf(record.value)]; exists {
			f.Source = String(source.value, "locator")
		}
	}
	v.findings = append(v.findings, f)
}

func (v *graphValidator) isA(kind, parent string) bool {
	return kind == parent || slices.Contains(v.catalog.Classes[kind].Ancestors, parent)
}

func (v *graphValidator) walk(object Record, class, path string, index bool) {
	if object == nil {
		return
	}
	if concrete := String(object, "kind"); v.isA(concrete, class) {
		class = concrete
	}
	definition := v.catalog.Classes[class]
	if index && v.isA(class, "Record") {
		id := String(object, "id")
		if earlier, found := v.records[id]; found {
			v.add("identity.duplicate", path+"/id", fmt.Sprintf("ID %q also occurs at %s", id, earlier.path))
		} else {
			v.records[id] = locatedRecord{object, path}
		}
	}
	if !index {
		v.checkValue(object, class, path)
	}
	names := make([]string, 0, len(definition.Slots))
	for name := range definition.Slots {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		slot := definition.Slots[name]
		value, present := object[name]
		if !present {
			continue
		}
		values := []any{value}
		if slot.Multivalued {
			values, _ = value.([]any)
		}
		for i, value := range values {
			location := path + "/" + escapePointer(name)
			if slot.Multivalued {
				location += "/" + strconv.Itoa(i)
			}
			if len(slot.ReferenceRanges) > 0 {
				if index {
					continue
				}
				id, _ := value.(string)
				target, exists := v.records[id]
				if !exists {
					v.add("reference.missing", location, fmt.Sprintf("record %q is absent from this snapshot", id))
					continue
				}
				allowed := false
				for _, parent := range slot.ReferenceRanges {
					allowed = allowed || v.isA(String(target.value, "kind"), parent)
				}
				if !allowed {
					v.add("reference.kind", location, fmt.Sprintf("%s is not an allowed target (%s)", String(target.value, "kind"), strings.Join(slot.ReferenceRanges, ", ")))
				}
			} else if _, exists := v.catalog.Classes[slot.Range]; exists {
				v.walk(Object(value), slot.Range, location, index)
			}
		}
	}
}

func escapePointer(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "~", "~0"), "/", "~1")
}

func records(value any) []Record {
	var result []Record
	switch values := value.(type) {
	case []Record:
		return values
	case []any:
		for _, value := range values {
			if object := Object(value); object != nil {
				result = append(result, object)
			}
		}
	}
	return result
}

func valuesAt(object Record, path string) []string {
	parts := strings.Split(path, ".")
	var value any = object
	for _, part := range parts {
		value = Object(value)[part]
	}
	if text, ok := value.(string); ok {
		return []string{text}
	}
	return Strings(value)
}

func (v *graphValidator) containingRecord(path string) (locatedRecord, bool) {
	var result locatedRecord
	found := false
	for _, record := range v.records {
		if (path == record.path || strings.HasPrefix(path, record.path+"/")) && (!found || len(record.path) > len(result.path)) {
			result = record
			found = true
		}
	}
	return result, found
}

// decodeSnapshotJSON rejects transports that encoding/json would otherwise
// normalize by replacing invalid UTF-8 or silently overwriting duplicate keys.
func decodeSnapshotJSON(data []byte) (any, error) {
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("snapshot JSON is not valid UTF-8")
	}
	if err := checkJSONSurrogates(data); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var read func(string, int) (any, error)
	read = func(path string, depth int) (any, error) {
		if depth > 1000 {
			return nil, fmt.Errorf("snapshot JSON nesting exceeds the supported model depth at %s", path)
		}
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("invalid snapshot JSON at %s: %w", path, err)
		}
		delimiter, isContainer := token.(json.Delim)
		if !isContainer {
			return token, nil
		}
		switch delimiter {
		case '{':
			object := map[string]any{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, fmt.Errorf("expected object key at %s", path)
				}
				if _, exists := object[key]; exists {
					return nil, fmt.Errorf("duplicate JSON key at %s/%s", path, escapePointer(key))
				}
				value, err := read(path+"/"+escapePointer(key), depth+1)
				if err != nil {
					return nil, err
				}
				object[key] = value
			}
			close, err := decoder.Token()
			if err != nil || close != json.Delim('}') {
				return nil, fmt.Errorf("invalid JSON object ending at %s", path)
			}
			return object, nil
		case '[':
			array := []any{}
			for decoder.More() {
				value, err := read(path+"/"+strconv.Itoa(len(array)), depth+1)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			close, err := decoder.Token()
			if err != nil || close != json.Delim(']') {
				return nil, fmt.Errorf("invalid JSON array ending at %s", path)
			}
			return array, nil
		default:
			return nil, fmt.Errorf("unexpected JSON delimiter at %s", path)
		}
	}
	value, err := read("", 0)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("snapshot JSON must contain exactly one complete value")
	}
	return value, nil
}

// JSON's Unicode escape syntax must not rely on the decoder's replacement of
// unpaired UTF-16 surrogates; that would silently change an authored value.
func checkJSONSurrogates(data []byte) error {
	quoted := false
	for i := 0; i < len(data); i++ {
		if data[i] == '"' {
			quoted = !quoted
			continue
		}
		if !quoted || data[i] != '\\' {
			continue
		}
		if i+1 >= len(data) {
			break
		}
		if data[i+1] != 'u' {
			i++
			continue
		}
		if i+6 > len(data) {
			break
		}
		code, err := strconv.ParseUint(string(data[i+2:i+6]), 16, 16)
		if err != nil {
			return fmt.Errorf("invalid JSON Unicode escape at byte %d", i)
		}
		if code >= 0xdc00 && code <= 0xdfff {
			return fmt.Errorf("unpaired JSON Unicode surrogate at byte %d", i)
		}
		if code >= 0xd800 && code <= 0xdbff {
			if i+12 > len(data) || data[i+6] != '\\' || data[i+7] != 'u' {
				return fmt.Errorf("unpaired JSON Unicode surrogate at byte %d", i)
			}
			low, err := strconv.ParseUint(string(data[i+8:i+12]), 16, 16)
			if err != nil || low < 0xdc00 || low > 0xdfff {
				return fmt.Errorf("unpaired JSON Unicode surrogate at byte %d", i)
			}
			i += 11
		} else {
			i += 5
		}
	}
	return nil
}
