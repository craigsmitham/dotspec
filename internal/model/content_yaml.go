package model

import (
	"bytes"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

func (e *contentExtractor) frontmatter() (int, string) {
	firstEnd := e.lineEnd(0)
	if strings.TrimRight(string(e.buffer[:firstEnd]), " \t\r\n") != "---" {
		return 0, ""
	}
	closeStart, closeEnd := -1, -1
	for at := firstEnd; at < len(e.buffer); {
		end := e.lineEnd(at)
		line := strings.TrimRight(string(e.buffer[at:end]), " \t\r\n")
		if line == "---" || line == "..." {
			closeStart, closeEnd = at, end
			break
		}
		at = end
	}
	if closeStart < 0 {
		r := e.node("frontmatter", 0, len(e.buffer))
		opaque := e.node("opaque", firstEnd, len(e.buffer))
		opaque["opaque_reason"] = "Frontmatter closing delimiter is missing"
		r["children"] = e.fillGaps(r, []string{String(opaque, "id")})
		e.p.Document["frontmatter"] = []string{String(r, "id")}
		e.p.Document["frontmatter_parse_status"] = "failed"
		e.p.Document["parse_status"] = "partial"
		e.diagnostic("content.frontmatter", "Frontmatter closing delimiter is missing.", 0, len(e.buffer))
		return len(e.buffer), String(r, "id")
	}
	r := e.node("frontmatter", 0, closeEnd)
	e.p.Document["frontmatter"] = []string{String(r, "id")}
	var metadata yaml.Node
	err := yaml.Unmarshal(e.buffer[firstEnd:closeStart], &metadata)
	children := []string{}
	if err != nil {
		opaque := e.node("opaque", firstEnd, closeStart)
		opaque["opaque_reason"] = "Malformed YAML: " + err.Error()
		children = append(children, String(opaque, "id"))
		e.p.Document["frontmatter_parse_status"] = "failed"
		e.p.Document["parse_status"] = "partial"
		e.diagnostic("content.frontmatter", err.Error(), firstEnd, closeStart)
	} else {
		e.p.Metadata = &metadata
		e.p.Document["frontmatter_parse_status"] = "parsed"
		if len(metadata.Content) > 0 {
			yn := metadata.Content[0]
			node := e.yamlNode(yn, firstEnd, closeStart, false)
			children = append(children, String(node, "id"))
			e.p.Document["metadata_root"] = node["id"]
			if yn.Kind == yaml.MappingNode {
				for i := 0; i+1 < len(yn.Content); i += 2 {
					key, value := yn.Content[i], yamlScalarTarget(yn.Content[i+1])
					if value == nil || value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
						continue
					}
					switch key.Value {
					case "type":
						e.p.Document["declared_type"] = value.Value
					case "id":
						e.p.Document["authored_id"] = value.Value
					}
				}
			}
		}
	}
	r["children"] = e.yamlGaps(r, children)
	return closeEnd, String(r, "id")
}
func (e *contentExtractor) yamlPosition(n *yaml.Node, base int) int {
	at := base
	for line := 1; line < n.Line && at < len(e.buffer); line++ {
		at = e.lineEnd(at)
	}
	for col := 1; col < n.Column && at < len(e.buffer); col++ {
		_, size := utf8.DecodeRune(e.buffer[at:])
		at += size
	}
	return at
}
func (e *contentExtractor) yamlNode(n *yaml.Node, base, limit int, flow bool) Record {
	start := e.yamlPosition(n, base)
	start = min(start, limit)
	kind := "yaml_scalar"
	switch n.Kind {
	case yaml.MappingNode:
		kind = "yaml_mapping"
	case yaml.SequenceNode:
		kind = "yaml_sequence"
	case yaml.AliasNode:
		kind = "yaml_alias"
	}
	flow = flow || n.Style&yaml.FlowStyle != 0
	end := e.yamlEnd(n, start, limit, flow)
	r := e.node(kind, start, end)
	if n.Tag != "" {
		r["scalar_tag"] = n.Tag
	}
	if n.Anchor != "" {
		r["anchor_name"] = n.Anchor
	}
	if n.Kind == yaml.AliasNode {
		r["alias_name"] = n.Value
		return r
	}
	children := []string{}
	switch n.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, value := n.Content[i], n.Content[i+1]
			entryLimit := end
			if i+2 < len(n.Content) {
				entryLimit = min(entryLimit, e.yamlPosition(n.Content[i+2], base))
			}
			k := e.yamlNode(key, base, e.yamlPosition(value, base), flow)
			v := e.yamlNode(value, base, entryLimit, flow)
			ks, _ := e.spanBounds(k)
			_, ve := e.spanBounds(v)
			entry := e.node("yaml_entry", e.normalizedOffset(ks), e.normalizedOffset(ve))
			entry["children"] = e.yamlGaps(entry, []string{String(k, "id"), String(v, "id")})
			children = append(children, String(entry, "id"))
		}
	case yaml.SequenceNode:
		for i, child := range n.Content {
			childLimit := end
			if i+1 < len(n.Content) {
				childLimit = min(childLimit, e.yamlPosition(n.Content[i+1], base))
			}
			v := e.yamlNode(child, base, childLimit, flow)
			children = append(children, String(v, "id"))
		}
	default:
		r["decoded_text"] = n.Value
		var value any
		_ = n.Decode(&value) // Lexical value and tag remain available if conversion is unavailable.
		switch n.Tag {
		case "!!str":
			r["scalar_kind"] = "string"
			r["scalar_string"] = n.Value
		case "!!bool":
			r["scalar_kind"] = "boolean"
			if v, ok := value.(bool); ok {
				r["scalar_boolean"] = v
			}
		case "!!int":
			r["scalar_kind"] = "integer"
			switch v := value.(type) {
			case int:
				r["scalar_integer"] = v
			case int64:
				r["scalar_integer"] = v
			case uint64:
				r["scalar_integer"] = v
			}
		case "!!float":
			r["scalar_kind"] = "float"
			if v, ok := value.(float64); ok && !math.IsInf(v, 0) && !math.IsNaN(v) {
				r["scalar_float"] = v
			}
		case "!!null":
			r["scalar_kind"] = "null"
		case "!!timestamp":
			r["scalar_kind"] = "timestamp"
			if v, ok := value.(time.Time); ok {
				r["scalar_timestamp"] = v.Format(time.RFC3339Nano)
			}
		case "!!binary":
			r["scalar_kind"] = "binary"
			r["scalar_binary_base64"] = strings.Join(strings.Fields(n.Value), "")
		default:
			r["scalar_kind"] = "custom"
			r["scalar_string"] = n.Value
		}
	}
	if n.Kind == yaml.ScalarNode && n.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 && n.LineComment != "" {
		headerEnd := min(e.lineEnd(start), end)
		if offset := bytes.Index(e.buffer[start:headerEnd], []byte(n.LineComment)); offset >= 0 {
			commentStart := start + offset
			if commentStart > start {
				prefix := e.node("trivia", start, commentStart)
				children = append(children, String(prefix, "id"))
			}
			comment := e.node("comment", commentStart, headerEnd)
			children = append(children, String(comment, "id"))
			if headerEnd < end {
				body := e.node("text", headerEnd, end)
				body["decoded_text"] = n.Value
				children = append(children, String(body, "id"))
			}
		}
	}
	if len(children) > 0 {
		r["children"] = e.yamlGaps(r, children)
	}
	return r
}
func (e *contentExtractor) yamlEnd(n *yaml.Node, start, limit int, flow bool) int {
	limit = max(start, min(limit, len(e.buffer)))
	at := start
	// Tags and anchors precede a scalar or flow collection, and are part of its span.
	for at < limit && (e.buffer[at] == '!' || e.buffer[at] == '&') {
		for at < limit && !bytes.ContainsRune([]byte(" \t\r\n"), rune(e.buffer[at])) {
			at++
		}
		for at < limit && (e.buffer[at] == ' ' || e.buffer[at] == '\t') {
			at++
		}
	}
	if at < limit && (e.buffer[at] == '\'' || e.buffer[at] == '"') {
		quote := e.buffer[at]
		for i := at + 1; i < limit; i++ {
			if quote == '"' && e.buffer[i] == '\\' {
				i++
				continue
			}
			if e.buffer[i] == quote {
				if quote == '\'' && i+1 < limit && e.buffer[i+1] == '\'' {
					i++
					continue
				}
				return i + 1
			}
		}
	}
	if at < limit && (e.buffer[at] == '[' || e.buffer[at] == '{') {
		left := e.buffer[at]
		right := byte(']')
		if left == '{' {
			right = '}'
		}
		depth := 0
		quote := byte(0)
		for i := at; i < limit; i++ {
			c := e.buffer[i]
			if quote != 0 {
				if c == '\\' && quote == '"' {
					i++
					continue
				}
				if c == quote {
					if quote == '\'' && i+1 < limit && e.buffer[i+1] == '\'' {
						i++
						continue
					}
					quote = 0
				}
				continue
			}
			if c == '#' && (i == at || e.buffer[i-1] == ' ' || e.buffer[i-1] == '\t' || e.buffer[i-1] == '\n') {
				i = min(e.lineEnd(i), limit) - 1
				continue
			}
			if c == '\'' || c == '"' {
				quote = c
				continue
			}
			if c == left {
				depth++
			}
			if c == right {
				depth--
				if depth == 0 {
					return i + 1
				}
			}
		}
	}
	if n.Kind == yaml.ScalarNode || n.Kind == yaml.AliasNode {
		if at < limit && (e.buffer[at] == '|' || e.buffer[at] == '>') {
			return e.yamlBlockEnd(start, at, limit)
		}
		return e.yamlPlainEnd(n, start, at, limit, flow)
	}
	// A collection includes its intervening whitespace and comments. Keeping the
	// untrimmed limit lets the final scalar retain its final line and chomp bytes.
	return limit
}

func (e *contentExtractor) yamlPlainEnd(n *yaml.Node, start, at, limit int, flow bool) int {
	scan := func(a, b int) int {
		for i := a; i < b; i++ {
			c := e.buffer[i]
			if c == '#' && (i == a || e.buffer[i-1] == ' ' || e.buffer[i-1] == '\t') {
				return i
			}
			if flow && (c == ',' || c == ']' || c == '}') {
				return i
			}
			if c == ':' && (i+1 == b || e.buffer[i+1] == ' ' || e.buffer[i+1] == '\n') {
				return i
			}
		}
		return b
	}
	firstLineEnd := min(e.lineEnd(at), limit)
	end := scan(at, firstLineEnd)
	if n.Kind != yaml.ScalarNode || n.Style != 0 || n.Tag != "!!str" || end < firstLineEnd {
		return trimYAMLWhitespace(e.buffer, start, end)
	}
	lineStart := e.lineStart(start)
	prefix := e.buffer[lineStart:start]
	indentation := len(prefix) - len(bytes.TrimLeft(prefix, " "))
	minimum := indentation
	if len(bytes.TrimSpace(prefix)) > 0 {
		minimum++
	}
	for next := e.lineEnd(at); next < limit; {
		stop := min(e.lineEnd(next), limit)
		line := e.buffer[next:stop]
		trim := bytes.TrimSpace(line)
		if len(trim) == 0 {
			next = stop
			continue
		}
		indent := len(line) - len(bytes.TrimLeft(line, " "))
		if trim[0] == '#' || (!flow && indent < minimum) {
			break
		}
		end = scan(next, stop)
		if end < stop {
			break
		}
		next = stop
	}
	return trimYAMLWhitespace(e.buffer, start, end)
}

func (e *contentExtractor) yamlBlockEnd(start, indicator, limit int) int {
	headerEnd := min(e.lineEnd(indicator), limit)
	line := e.buffer[e.lineStart(start):start]
	baseIndent := len(line) - len(bytes.TrimLeft(line, " "))
	contentIndent := -1
	for _, c := range e.buffer[indicator:headerEnd] {
		if c >= '1' && c <= '9' {
			contentIndent = baseIndent + int(c-'0')
			break
		}
		if c == '#' {
			break
		}
	}
	end := headerEnd
	for at := headerEnd; at < limit; {
		stop := min(e.lineEnd(at), limit)
		line := e.buffer[at:stop]
		trim := bytes.TrimSpace(line)
		if len(trim) > 0 {
			indent := len(line) - len(bytes.TrimLeft(line, " "))
			if contentIndent < 0 {
				if indent <= baseIndent {
					break
				}
				contentIndent = indent
			}
			if indent < contentIndent {
				break
			}
		}
		end = stop
		at = stop
	}
	return end
}

func trimYAMLWhitespace(data []byte, start, end int) int {
	for end > start && (data[end-1] == ' ' || data[end-1] == '\t' || data[end-1] == '\n' || data[end-1] == '\r') {
		end--
	}
	return end
}
func (e *contentExtractor) yamlGaps(parent Record, children []string) []string {
	all := e.fillGaps(parent, children)
	out := []string{}
	for _, id := range all {
		r := e.p.NodeByID[id]
		if String(r, "node_kind") != "trivia" {
			out = append(out, id)
			continue
		}
		a, b := e.spanBounds(r)
		start, end := e.normalizedOffset(a), e.normalizedOffset(b)
		at := start
		pieces := []string{}
		for i := start; i < end; i++ {
			if e.buffer[i] != '#' {
				continue
			}
			if at < i {
				piece := e.node("trivia", at, i)
				pieces = append(pieces, String(piece, "id"))
			}
			stop := min(e.lineEnd(i), end)
			comment := e.node("comment", i, stop)
			pieces = append(pieces, String(comment, "id"))
			at = stop
			i = stop - 1
		}
		if len(pieces) == 0 {
			out = append(out, id)
			continue
		}
		if at < end {
			piece := e.node("trivia", at, end)
			pieces = append(pieces, String(piece, "id"))
		}
		out = append(out, pieces...)
		// The original gap was replaced, so remove it from the snapshot inventory.
		delete(e.p.NodeByID, id)
		for i, n := range e.p.Nodes {
			if String(n, "id") == id {
				e.p.Nodes = append(e.p.Nodes[:i], e.p.Nodes[i+1:]...)
				break
			}
		}
	}
	return out
}

// Metadata interpretation may follow YAML aliases, while syntax nodes retain
// the authored alias occurrence and its own span instead of copying the target.
func yamlScalarTarget(n *yaml.Node) *yaml.Node {
	seen := map[*yaml.Node]bool{}
	for n != nil && n.Kind == yaml.AliasNode {
		if seen[n] {
			return nil
		}
		seen[n] = true
		n = n.Alias
	}
	return n
}
