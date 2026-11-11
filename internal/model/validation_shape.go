package model

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

// shapeFindings narrows tagged-union errors to the declared concrete class.
// A wrong UseCase field should identify that field, not print every unrelated
// domain subtype's failed discriminator and required-property list.
func shapeFindings(runtime modelRuntime, value any, err error) []Finding {
	var invalid *jsonschema.ValidationError
	if !errors.As(err, &invalid) {
		return []Finding{{Code: "shape.invalid", Path: "/", Message: err.Error()}}
	}
	leaves := validationLeaves(invalid)
	findings := []Finding{}
	handled := []string{}
	root := Object(value)
	index := shapeRecordIndex(root)
	failureBelow := func(path string) bool {
		for _, leaf := range leaves {
			p := instancePointer(leaf.InstanceLocation)
			if p == path || strings.HasPrefix(p, path+"/") {
				return true
			}
		}
		return false
	}
	var visit func(any, string, string)
	visit = func(value any, class, path string) {
		object := Object(value)
		if object == nil {
			return
		}
		definition := runtime.catalog.Classes[class]
		if slices.Contains(definition.Ancestors, "Record") && failureBelow(path) {
			declared := String(object, "kind")
			concrete, known := runtime.catalog.Classes[declared]
			if !known || concrete.Abstract || !slices.Contains(concrete.Ancestors, class) {
				finding := Finding{Code: "shape.invalid", Path: path + "/kind", Message: fmt.Sprintf("kind %q must identify a concrete %s record", declared, class)}
				locateShapeFinding(&finding, root, index, object, path)
				findings = append(findings, finding)
				handled = append(handled, path)
				return
			}
			if classErr := runtime.classes[declared].Validate(map[string]any(object)); classErr != nil {
				var classInvalid *jsonschema.ValidationError
				if errors.As(classErr, &classInvalid) {
					for _, leaf := range validationLeaves(classInvalid) {
						for _, finding := range leafFindings(leaf, path) {
							locateShapeFinding(&finding, root, index, object, path)
							findings = append(findings, finding)
						}
					}
					handled = append(handled, path)
					return
				}
			}
		}
		if declared := String(object, "kind"); slices.Contains(runtime.catalog.Classes[declared].Ancestors, class) {
			definition = runtime.catalog.Classes[declared]
		}
		names := make([]string, 0, len(definition.Slots))
		for name := range definition.Slots {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			slot := definition.Slots[name]
			if !slot.Inlined || len(slot.InlineRanges) == 0 {
				continue
			}
			if slot.Multivalued {
				array, _ := object[name].([]any)
				for i, item := range array {
					visit(item, slot.InlineRanges[0], path+"/"+escapePointer(name)+"/"+strconv.Itoa(i))
				}
			} else {
				visit(object[name], slot.InlineRanges[0], path+"/"+escapePointer(name))
			}
		}
	}
	visit(value, "Snapshot", "")
	for _, leaf := range leaves {
		path := instancePointer(leaf.InstanceLocation)
		skip := false
		for _, recordPath := range handled {
			if path == recordPath || strings.HasPrefix(path, recordPath+"/") {
				skip = true
				break
			}
		}
		if !skip {
			findings = append(findings, leafFindings(leaf, "")...)
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		return findings[i].Message < findings[j].Message
	})
	findings = slices.CompactFunc(findings, func(a, b Finding) bool { return a.Path == b.Path && a.Message == b.Message })
	if len(findings) > 20 {
		findings = findings[:20]
	}
	return findings
}

func validationLeaves(err *jsonschema.ValidationError) []*jsonschema.ValidationError {
	if len(err.Causes) == 0 {
		return []*jsonschema.ValidationError{err}
	}
	var leaves []*jsonschema.ValidationError
	for _, cause := range err.Causes {
		leaves = append(leaves, validationLeaves(cause)...)
	}
	return leaves
}
func instancePointer(parts []string) string {
	var path strings.Builder
	for _, part := range parts {
		path.WriteByte('/')
		path.WriteString(escapePointer(part))
	}
	return path.String()
}
func leafFindings(leaf *jsonschema.ValidationError, prefix string) []Finding {
	path := prefix + instancePointer(leaf.InstanceLocation)
	message := leaf.Error()
	if output := leaf.BasicOutput(); output.Error != nil {
		message = output.Error.String()
	}
	if runes := []rune(message); len(runes) > 1000 {
		message = string(runes[:1000]) + "…"
	}
	fields := []string{}
	switch problem := leaf.ErrorKind.(type) {
	case *kind.Required:
		fields = problem.Missing
	case *kind.AdditionalProperties:
		fields = problem.Properties
	}
	if len(fields) == 0 {
		if path == "" {
			path = "/"
		}
		return []Finding{{Code: "shape.invalid", Path: path, Message: message}}
	}
	findings := []Finding{}
	for _, field := range fields {
		if len(findings) >= 20 {
			break
		}
		findings = append(findings, Finding{Code: "shape.invalid", Path: path + "/" + escapePointer(field), Message: message})
	}
	return findings
}

func shapeRecordIndex(root Record) map[string]Record {
	index := map[string]Record{}
	for _, name := range []string{"sources", "documents", "content_nodes", "elements", "references", "label_occurrences", "fragments"} {
		for _, r := range records(root[name]) {
			if id := String(r, "id"); id != "" {
				index[id] = r
			}
		}
	}
	return index
}

func locateShapeFinding(f *Finding, root Record, index map[string]Record, record Record, recordPath string) {
	source := String(record, "source_document")
	if source == "" {
		source = String(record, "source")
	}
	if String(record, "kind") == "SourceDocument" {
		source = String(record, "id")
	}
	if source == "" && String(record, "kind") == "ContentNode" {
		source = String(index[String(record, "document")], "source")
	}
	if source == "" {
		source = String(index[String(record, "owner")], "source_document")
	}
	f.Source = String(index[source], "locator")
	node := index[String(record, "content_node")]
	relative := strings.TrimPrefix(f.Path, recordPath+"/")
	best := -1
	for _, binding := range records(root["content_bindings"]) {
		if String(binding, "owner") != String(record, "id") {
			continue
		}
		field := strings.ReplaceAll(String(binding, "field_path"), ".", "/")
		if relative != field && !strings.HasPrefix(relative, field+"/") {
			continue
		}
		if len(field) > best {
			node = index[String(binding, "content_node")]
			best = len(field)
		}
	}
	span := Object(node["span"])
	if span == nil {
		span = Object(record["span"])
	}
	if span == nil {
		roots := Strings(Object(record["name"])["content_nodes"])
		if len(roots) > 0 {
			span = Object(index[roots[0]]["span"])
		}
	}
	if span != nil {
		if spanSource := String(span, "source"); spanSource != "" {
			source = spanSource
			f.Source = String(index[source], "locator")
		}
		start, ok := natural(span["start"])
		data, available := sourceBytes(index[source])
		if ok && available && start <= len(data) {
			f.Line = bytes.Count(data[:start], []byte{'\n'}) + 1
			f.Column = start - bytes.LastIndexByte(data[:start], '\n')
		}
	}
}
