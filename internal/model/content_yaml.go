package model

import (
	"bytes"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

func (e *contentExtractor) frontmatter() int {
	md := e.p.Syntax
	end := len(md.Source) - len(md.Body)
	if md.MetadataError != nil {
		limit := end
		if limit == 0 {
			limit = e.lineEnd(0)
		}
		e.diagnostic("content.frontmatter", md.MetadataError.Error(), 0, limit)
		e.p.Document["parse_status"] = "partial"
	}
	for _, key := range []string{"type", "id"} {
		if n := yamlScalarTarget(md.Value(key)); n != nil && n.Kind == yaml.ScalarNode && n.ShortTag() == "!!str" {
			field := "declared_type"
			if key == "id" {
				field = "authored_id"
			}
			e.p.Document[field] = n.Value
		}
	}
	return end
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
