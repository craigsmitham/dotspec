// Package system reads and discovers system specifications.
package system

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.yaml.in/yaml/v3"
)

const Filename = "system.spec.md"

const Template = `---
type: system
---

# New system

## Purpose

Describe why this system exists and whom it serves.

## Scope

Describe the system's boundary and what it includes.
`

type System struct {
	Name    string
	Purpose string
	Scope   string
	Path    string
}

func Read(path string) (*System, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	spec.Path = path
	return spec, nil
}

// Parse accepts YAML frontmatter and Markdown with one document-level H1 and
// nonempty Purpose and Scope H2 sections. Section bodies retain their Markdown.
func Parse(data []byte) (*System, error) {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) == 0 || string(bytes.TrimSpace(lines[0])) != "---" {
		return nil, fmt.Errorf("expected YAML frontmatter beginning with ---")
	}
	end := 1
	for end < len(lines) && string(bytes.TrimSpace(lines[end])) != "---" {
		end++
	}
	if end == len(lines) {
		return nil, fmt.Errorf("YAML frontmatter must end with ---")
	}
	var metadata struct {
		Type string `yaml:"type"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(bytes.Join(lines[1:end], []byte("\n"))))
	decoder.KnownFields(true)
	if err := decoder.Decode(&metadata); err != nil {
		return nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}
	if metadata.Type != "system" {
		return nil, fmt.Errorf("frontmatter must declare type: system")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("frontmatter must contain exactly one YAML document")
	}

	body := bytes.Join(lines[end+1:], []byte("\n"))
	doc := goldmark.DefaultParser().Parse(text.NewReader(body))
	var headings []*ast.Heading
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if h, ok := node.(*ast.Heading); ok && h.Level <= 2 {
			headings = append(headings, h)
		}
	}
	spec := &System{}
	titles := 0
	sections := make(map[string]string)
	for i, h := range headings {
		name := strings.TrimSpace(string(h.Text(body)))
		if h.Level == 1 {
			titles++
			spec.Name = name
			continue
		}
		if name != "Purpose" && name != "Scope" {
			continue
		}
		if titles == 0 {
			return nil, fmt.Errorf("system name heading must precede %s", name)
		}
		if _, exists := sections[name]; exists {
			return nil, fmt.Errorf("duplicate %s section", name)
		}
		stop := len(body)
		if i+1 < len(headings) {
			stop = lineStart(body, headings[i+1].Pos())
		}
		sections[name] = strings.TrimSpace(string(body[headingEnd(body, h):stop]))
	}
	if titles != 1 || spec.Name == "" {
		return nil, fmt.Errorf("expected exactly one nonempty top-level heading naming the system")
	}
	for _, name := range []string{"Purpose", "Scope"} {
		if sections[name] == "" {
			return nil, fmt.Errorf("missing or empty %s section (use ## %s)", name, name)
		}
	}
	spec.Purpose = sections["Purpose"]
	spec.Scope = sections["Scope"]
	return spec, nil
}

func lineStart(source []byte, pos int) int {
	return bytes.LastIndexByte(source[:pos], '\n') + 1
}

func lineEnd(source []byte, pos int) int {
	if n := bytes.IndexByte(source[pos:], '\n'); n >= 0 {
		return pos + n + 1
	}
	return len(source)
}

func headingEnd(source []byte, h *ast.Heading) int {
	last := h.Pos()
	if h.Lines().Len() > 0 {
		last = h.Lines().At(h.Lines().Len() - 1).Start
	}
	end := lineEnd(source, last)
	// Setext headings have an underline following their text lines.
	if !bytes.HasPrefix(bytes.TrimSpace(source[lineStart(source, h.Pos()):]), []byte("#")) {
		end = lineEnd(source, end)
	}
	return end
}
