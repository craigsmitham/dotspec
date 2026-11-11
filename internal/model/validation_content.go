package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func natural(value any) (int, bool) {
	var text string
	switch value := value.(type) {
	case json.Number:
		text = string(value)
	case int:
		return value, value >= 0
	case float64:
		text = strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return 0, false
	}
	if n, err := strconv.ParseUint(text, 10, strconv.IntSize-1); err == nil {
		return int(n), true
	}
	number, ok := new(big.Rat).SetString(text)
	if !ok || !number.IsInt() || !number.Num().IsInt64() || number.Sign() < 0 {
		return 0, false
	}
	n := number.Num().Int64()
	if uint64(n) > uint64(^uint(0)>>1) {
		return 0, false
	}
	return int(n), true
}

func (v *graphValidator) checkValue(object Record, class, path string) {
	switch class {
	case "SourcePayload":
		if String(object, "storage") != "inline_base64" {
			return // External byte availability is not asserted by offline validation.
		}
		data, err := base64.StdEncoding.Strict().DecodeString(String(object, "bytes_base64"))
		if err != nil {
			v.add("payload.encoding", path+"/bytes_base64", "payload is not valid base64")
			return
		}
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != String(object, "sha256") {
			v.add("payload.digest", path+"/sha256", "digest does not match the original bytes")
		}
		length, valid := natural(object["byte_length"])
		if !valid || length != len(data) {
			v.add("payload.length", path+"/byte_length", "byte count does not match the original bytes")
		}
	case "SourceDocument":
		if String(object, "read_status") != "read" && object["payload"] != nil {
			v.add("payload.read-state", path+"/payload", "a failed or unattempted read cannot claim an acquired payload")
		}
	case "ByteSpan":
		start, validStart := natural(object["start"])
		end, validEnd := natural(object["end"])
		source := v.records[String(object, "source")].value
		length, validLength := natural(Object(source["payload"])["byte_length"])
		if !validStart || !validEnd || !validLength || start > end || end > length {
			v.add("source.span", path, "span must identify an ordered interval within acquired source bytes")
		}
	case "RichContent":
		nodes := Strings(object["content_nodes"])
		if (String(object, "presence") == "empty") != (len(nodes) == 0) {
			v.add("content.presence", path, "empty content has no body nodes; authored content has at least one")
		}
		lastEnd, source := -1, ""
		expectedSource := ""
		if parent, ok := v.containingRecord(path); ok {
			expectedSource = v.sourceOf(parent.value)
		}
		seen := map[string]bool{}
		for _, id := range nodes {
			node := v.records[id].value
			span := Object(node["span"])
			if node == nil {
				continue
			}
			if seen[id] || (source != "" && source != String(span, "source")) || validationNumber(span["start"]) < lastEnd {
				v.add("content.roots", path+"/content_nodes", "body roots must be ordered, distinct, non-overlapping occurrences in one source")
			}
			if expectedSource != "" && expectedSource != String(span, "source") {
				v.add("content.ownership", path+"/content_nodes", "domain content must be expressed in its record's source")
			}
			seen[id], lastEnd, source = true, validationNumber(span["end"]), String(span, "source")
		}
	case "Resolution":
		if String(object, "status") != "resolved" && object["target"] != nil {
			v.add("resolution.state", path+"/target", "only a resolved occurrence can claim a target")
		}
	case "SourceFragment", "Diagnostic", "SourceLocation":
		if span := Object(object["span"]); span != nil && String(object, "source") != "" && String(span, "source") != String(object, "source") {
			v.add("source.ownership", path+"/span", "source span disagrees with its location's source")
		}
		if class == "SourceFragment" {
			if node := v.records[String(object, "content_node")].value; node != nil && v.sourceOf(node) != String(object, "source") {
				v.add("source.ownership", path+"/content_node", "fragment node must belong to its source document")
			}
		}
	case "ContentNode":
		if String(object, "node_kind") == "opaque" && String(object, "opaque_reason") == "" {
			v.add("content.opaque-reason", path, "opaque content must state why it was not interpreted")
		}
	case "ContentBinding":
		owner := v.records[String(object, "owner")].value
		value, exists := fieldValue(owner, String(object, "field_path"))
		if owner != nil && !exists {
			v.add("binding.field", path+"/field_path", "binding field does not exist on its owner")
		}
		node := v.records[String(object, "content_node")].value
		if owner != nil && node != nil && v.sourceOf(owner) != v.sourceOf(node) {
			v.add("binding.source", path, "binding node and owner must have the same source")
		}
		role := String(object, "role")
		if (role == "fallback" || role == "extension") && String(object, "reason") == "" {
			v.add("binding.reason", path, "a fallback or extension binding must identify its interpretation boundary")
		}
		if String(owner, "kind") == "Document" && String(object, "field_path") == "root" && role != "structural_wrapper" && role != "fallback" {
			v.add("binding.root-role", path, "document root only provides structural or explicit fallback coverage")
		}
		if String(owner, "kind") == "Document" && String(object, "field_path") == "root" && role == "structural_wrapper" && String(node, "node_kind") != "document" && String(node, "node_kind") != "trivia" {
			v.add("binding.root-role", path, "document-root structural bindings may describe the root or retained trivia, not body content")
		}
		if exists && role != "structural_wrapper" && role != "fallback" {
			roots := v.representedNodes(value, map[string]bool{})
			if len(roots) > 0 {
				matched := false
				for _, root := range roots {
					if v.descendsFrom(String(object, "content_node"), root) {
						matched = true
					}
				}
				if !matched {
					v.add("binding.content", path, "binding node is outside the source content represented by its field")
				}
			}
		}
	case "CoverageReport":
		if String(object, "status") == "complete_for_profile" && (len(Strings(object["unbound_content_nodes"])) > 0 || len(records(object["unmapped_known_roles"])) > 0) {
			v.add("coverage.incomplete", path, "complete coverage cannot contain unbound nodes or unmapped known roles")
		}
	}
}

func fieldValue(owner Record, path string) (any, bool) {
	if path == "" {
		return nil, false
	}
	path = strings.TrimPrefix(path, "/")
	separator := "."
	if strings.Contains(path, "/") {
		separator = "/"
	}
	path = strings.ReplaceAll(strings.ReplaceAll(path, "[", separator), "]", "")
	var value any = owner
	for _, part := range strings.Split(path, separator) {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		if object := Object(value); object != nil {
			var exists bool
			value, exists = object[part]
			if !exists {
				return nil, false
			}
		} else if array, ok := value.([]any); ok {
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(array) {
				return nil, false
			}
			value = array[i]
		} else {
			return nil, false
		}
	}
	return value, true
}

// representedNodes follows embedded values and a directly bound record, but
// not that record's owner/related references. Source tracing must not spread
// through arbitrary graph edges until an unrelated field appears to match.
func (v *graphValidator) representedNodes(value any, seen map[string]bool) []string {
	if id, ok := value.(string); ok {
		if seen[id] {
			return nil
		}
		seen[id] = true
		r := v.records[id].value
		switch String(r, "kind") {
		case "ContentNode":
			return []string{id}
		case "SourceReference", "LabelOccurrence":
			return []string{String(r, "content_node")}
		}
		if r != nil {
			result := []string{}
			for name, slot := range v.catalog.Classes[String(r, "kind")].Slots {
				if slot.Inlined {
					result = append(result, v.representedNodes(r[name], seen)...)
				}
			}
			return result
		}
		return nil
	}
	if array, ok := value.([]any); ok {
		result := []string{}
		for _, item := range array {
			result = append(result, v.representedNodes(item, seen)...)
		}
		return result
	}
	if object := Object(value); object != nil {
		if roots, ok := object["content_nodes"]; ok {
			return Strings(roots)
		}
		result := []string{}
		for _, item := range object {
			result = append(result, v.representedNodes(item, seen)...)
		}
		return result
	}
	return nil
}

func (v *graphValidator) descendsFrom(id, root string) bool {
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(current string) bool {
		if current == id {
			return true
		}
		if seen[current] {
			return false
		}
		seen[current] = true
		for _, child := range Strings(v.records[current].value["children"]) {
			if visit(child) {
				return true
			}
		}
		return false
	}
	return visit(root)
}

func (v *graphValidator) checkContent() {
	parents := map[string]string{}
	sources := map[string]string{}
	for _, location := range v.sortedRecords() {
		node := location.value
		if String(node, "kind") != "ContentNode" {
			continue
		}
		document := v.records[String(node, "document")].value
		span := Object(node["span"])
		if String(span, "source") != String(document, "source") {
			v.add("content.source", location.path, "node and document must refer to the same source")
		}
		end := validationNumber(span["start"])
		for _, id := range Strings(node["children"]) {
			child := v.records[id].value
			if child == nil {
				continue
			}
			if prior := parents[id]; prior != "" {
				v.add("content.parent", location.path+"/children", fmt.Sprintf("node %q already has parent %q", id, prior))
			}
			parents[id] = String(node, "id")
			childSpan := Object(child["span"])
			if validationNumber(childSpan["start"]) > end {
				v.add("content.gap", location.path+"/children", "child content omits original source bytes; retain the gap as syntax, trivia or opaque content")
			}
			if String(child, "document") != String(node, "document") || validationNumber(childSpan["start"]) < end || validationNumber(childSpan["end"]) > validationNumber(span["end"]) {
				v.add("content.order", location.path+"/children", "children must be ordered, non-overlapping and contained in their parent's document/span")
			}
			end = validationNumber(childSpan["end"])
		}
		if (len(Strings(node["children"])) > 0 || String(node, "node_kind") == "document") && end != validationNumber(span["end"]) {
			v.add("content.gap", location.path+"/children", "children must account for the complete parent source span")
		}
	}
	v.checkCycles("content.cycle", "ContentNode", func(record Record) []string { return Strings(record["children"]) })
	for _, location := range v.sortedRecords() {
		document := location.value
		if String(document, "kind") != "Document" {
			continue
		}
		id, root := String(document, "id"), String(document, "root")
		if parents[root] != "" || String(v.records[root].value, "document") != id {
			v.add("content.root", location.path+"/root", "document root must be its own parentless content node")
		}
		source := v.records[String(document, "source")].value
		sourceID := String(document, "source")
		if prior := sources[sourceID]; prior != "" {
			v.add("document.duplicate-source", location.path+"/source", fmt.Sprintf("source already has parsed document %q", prior))
		}
		sources[sourceID] = id
		rootRecord := v.records[root].value
		rootSpan := Object(rootRecord["span"])
		if String(rootRecord, "node_kind") != "document" || validationNumber(rootSpan["start"]) != 0 || validationNumber(rootSpan["end"]) != validationNumber(Object(source["payload"])["byte_length"]) {
			v.add("content.root-span", location.path+"/root", "document root must cover every original source byte")
		}
		for _, front := range Strings(document["frontmatter"]) {
			n := v.records[front].value
			if String(n, "document") != id || String(n, "node_kind") != "frontmatter" {
				v.add("content.frontmatter", location.path+"/frontmatter", "frontmatter roots must be frontmatter nodes in this document")
			}
		}
		if metadata := v.records[String(document, "metadata_root")].value; metadata != nil && (String(metadata, "document") != id || !strings.HasPrefix(String(metadata, "node_kind"), "yaml_")) {
			v.add("content.metadata", location.path+"/metadata_root", "metadata root must be parsed YAML in this document")
		}
		if String(source, "read_status") != "read" {
			v.add("document.read-state", location.path, "parsed document requires acquired source bytes")
		}
		reachable := map[string]bool{}
		var visit func(string)
		visit = func(node string) {
			if reachable[node] {
				return
			}
			reachable[node] = true
			for _, child := range Strings(v.records[node].value["children"]) {
				visit(child)
			}
		}
		visit(root)
		for _, node := range v.sortedRecords() {
			if String(node.value, "kind") == "ContentNode" && String(node.value, "document") == id && !reachable[String(node.value, "id")] {
				v.add("content.orphan", node.path, "node is not reachable from its document root")
			}
		}
		v.checkDocumentCoverage(document, location.path, reachable)
		v.checkNormalization(document, location.path, source)
	}
	v.checkSnapshotCoverage()
}

func (v *graphValidator) checkDocumentCoverage(document Record, path string, nodes map[string]bool) {
	coverage := Object(document["coverage"])
	if String(coverage, "mapping_profile") != String(v.root, "mapping_profile") {
		v.add("coverage.profile", path+"/coverage/mapping_profile", "document coverage must name the snapshot mapping profile")
	}
	if String(coverage, "status") == "not_checked" {
		return
	}
	covered := map[string]bool{}
	uninterpreted, extensions := []string{}, []string{}
	for id := range nodes {
		if String(v.records[id].value, "node_kind") == "opaque" {
			uninterpreted = append(uninterpreted, id)
		}
	}
	for _, element := range v.sortedRecords() {
		if String(element.value, "kind") == "ExtensionSection" && v.sourceOf(element.value) == String(document, "source") {
			extensions = append(extensions, String(element.value, "id"))
		}
	}
	for _, binding := range records(v.root["content_bindings"]) {
		id := String(binding, "content_node")
		if !nodes[id] {
			continue
		}
		role := String(binding, "role")
		if role == "fallback" {
			uninterpreted = append(uninterpreted, id)
		}
		if known := String(binding, "known_role"); known != "" && (role == "fallback" || role == "extension") {
			reported := false
			for _, gap := range records(coverage["unmapped_known_roles"]) {
				if String(gap, "role") == known && slices.Contains(Strings(gap["content_nodes"]), id) {
					reported = true
				}
			}
			if !reported {
				v.add("coverage.known-role", path+"/coverage", "a known role mapped as fallback must be reported as unmapped")
			}
			if String(coverage, "status") == "complete_for_profile" {
				v.add("coverage.incomplete", path+"/coverage", "known roles retained only as fallback prevent complete domain mapping")
			}
		}
		var cover func(string)
		visited := map[string]bool{}
		cover = func(current string) {
			if visited[current] {
				return
			}
			visited[current] = true
			node := v.records[current].value
			// A heading's inline syntax can share its explicit heading binding.
			// Structural wrappers elsewhere cannot stand in for domain body mapping.
			if current == id || role != "structural_wrapper" || String(v.records[id].value, "node_kind") == "heading" || String(node, "node_kind") == "trivia" {
				covered[current] = true
			}
			if String(binding, "scope") == "subtree" {
				for _, child := range Strings(node["children"]) {
					cover(child)
				}
			}
		}
		cover(id)
	}
	actualUnbound := []string{}
	for id := range nodes {
		if !covered[id] {
			actualUnbound = append(actualUnbound, id)
		}
	}
	for field, actual := range map[string][]string{"unbound_content_nodes": actualUnbound, "uninterpreted_nodes": uninterpreted, "extension_sections": extensions} {
		reported := Strings(coverage[field])
		if len(reported) != len(uniqueStrings(reported)) || !sameStrings(reported, actual) {
			v.add("coverage.accounting", path+"/coverage/"+field, "coverage inventory disagrees with observed bindings and content records")
		}
	}
	if String(coverage, "status") == "complete_for_profile" && len(actualUnbound) > 0 {
		v.add("coverage.incomplete", path+"/coverage", "complete coverage must account for every content occurrence")
	}
	for _, gap := range records(coverage["unmapped_known_roles"]) {
		for _, id := range Strings(gap["content_nodes"]) {
			if !nodes[id] {
				v.add("coverage.source", path+"/coverage/unmapped_known_roles", "coverage gap refers outside this document")
			}
		}
	}
}

func (v *graphValidator) checkSnapshotCoverage() {
	coverage := Object(v.root["coverage"])
	if String(coverage, "mapping_profile") != String(v.root, "mapping_profile") {
		v.add("coverage.profile", "/coverage/mapping_profile", "coverage must name the snapshot mapping profile")
	}
	if String(coverage, "status") == "not_checked" {
		return
	}
	allComplete := true
	expected := map[string][]string{"unbound_content_nodes": {}, "uninterpreted_nodes": {}, "extension_sections": {}}
	gaps := []string{}
	for _, document := range records(v.root["documents"]) {
		child := Object(document["coverage"])
		if String(child, "status") != "complete_for_profile" {
			allComplete = false
		}
		for field := range expected {
			expected[field] = append(expected[field], Strings(child[field])...)
		}
		for _, gap := range records(child["unmapped_known_roles"]) {
			for _, id := range Strings(gap["content_nodes"]) {
				gaps = append(gaps, String(gap, "role")+"\x00"+id)
			}
		}
	}
	for field, actual := range expected {
		if !sameStrings(Strings(coverage[field]), actual) {
			v.add("coverage.aggregate", "/coverage/"+field, "snapshot coverage must aggregate document coverage inventories")
		}
	}
	reported := []string{}
	for _, gap := range records(coverage["unmapped_known_roles"]) {
		for _, id := range Strings(gap["content_nodes"]) {
			reported = append(reported, String(gap, "role")+"\x00"+id)
		}
	}
	if !sameStrings(reported, gaps) {
		v.add("coverage.aggregate", "/coverage/unmapped_known_roles", "snapshot known-role gaps must aggregate document gaps")
	}
	if String(coverage, "status") == "complete_for_profile" {
		if !allComplete {
			v.add("coverage.incomplete", "/coverage", "snapshot cannot claim complete mapping when a document has partial or unchecked coverage")
		}
		for _, source := range records(v.root["sources"]) {
			if String(source, "kind") == "SourceDocument" && String(source, "read_status") != "read" {
				v.add("coverage.incomplete", "/coverage", "unread source documents prevent a complete mapping claim")
			}
		}
	}
}

func (v *graphValidator) checkNormalization(document Record, path string, source Record) {
	segments := records(document["normalization"])
	payload := Object(source["payload"])
	data, inline := sourceBytes(source)
	if len(segments) == 0 {
		if inline && strings.HasPrefix(String(document, "parser"), "goldmark/") && (bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) || bytes.Contains(data, []byte("\r\n"))) {
			v.add("source.offset-map", path+"/normalization", "normalized parser input requires its original-byte offset map")
		}
		return
	}
	originalEnd, normalizedEnd := 0, 0
	for _, segment := range segments {
		os, oe := validationNumber(segment["original_start"]), validationNumber(segment["original_end"])
		ns, ne := validationNumber(segment["normalized_start"]), validationNumber(segment["normalized_end"])
		if os != originalEnd || ns != normalizedEnd || oe < os || ne < ns || oe > validationNumber(payload["byte_length"]) {
			v.add("source.offset-map", path+"/normalization", "normalization segments must form monotone contiguous in-bounds intervals")
		}
		transformation := String(segment, "transformation")
		if transformation == "" {
			v.add("source.offset-rule", path+"/normalization", "normalization must identify its transformation")
		}
		if os >= 0 && oe >= os && oe <= len(data) && inline {
			switch transformation {
			case "identity":
				if oe-os != ne-ns {
					v.add("source.offset-rule", path+"/normalization", "identity mapping must preserve byte counts")
				}
			case "BOM removal":
				if os != 0 || ne != ns || !bytes.Equal(data[os:oe], []byte{0xef, 0xbb, 0xbf}) {
					v.add("source.offset-rule", path+"/normalization", "BOM removal must identify the original leading UTF-8 BOM")
				}
			case "CRLF to LF":
				if ne-ns != 1 || !bytes.Equal(data[os:oe], []byte("\r\n")) {
					v.add("source.offset-rule", path+"/normalization", "CRLF normalization must identify exactly the original CRLF bytes")
				}
			}
		}
		originalEnd, normalizedEnd = oe, ne
	}
	if originalEnd != validationNumber(payload["byte_length"]) {
		v.add("source.offset-map", path+"/normalization", "normalization map must cover all original bytes")
	}
}

func sourceBytes(source Record) ([]byte, bool) {
	payload := Object(source["payload"])
	if String(payload, "storage") != "inline_base64" {
		return nil, false
	}
	data, err := base64.StdEncoding.Strict().DecodeString(String(payload, "bytes_base64"))
	return data, err == nil
}

func (v *graphValidator) sortedRecords() []locatedRecord {
	if v.orderedRecords != nil {
		return v.orderedRecords
	}
	result := make([]locatedRecord, 0, len(v.records))
	for _, record := range v.records {
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].path < result[j].path })
	v.orderedRecords = result
	return result
}

func (v *graphValidator) checkCycles(code, kind string, edges func(Record) []string) {
	state := map[string]int{}
	var visit func(string)
	visit = func(id string) {
		r := v.records[id]
		if state[id] == 1 {
			v.add(code, r.path, "relationship hierarchy contains a cycle")
			return
		}
		if state[id] == 2 || !v.isA(String(r.value, "kind"), kind) {
			return
		}
		state[id] = 1
		for _, target := range edges(r.value) {
			visit(target)
		}
		state[id] = 2
	}
	for _, record := range v.sortedRecords() {
		visit(String(record.value, "id"))
	}
}

func validationNumber(value any) int { number, _ := natural(value); return number }
