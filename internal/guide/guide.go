// Package guide resolves handbook topics and worked examples in the installed binary.
package guide

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/docs"
	"github.com/craigsmitham/dotspec/examples"
	"github.com/craigsmitham/dotspec/internal/specfile"
	modelschema "github.com/craigsmitham/dotspec/model"
	"github.com/craigsmitham/dotspec/specs"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// The CLI and AXM distributions use the same maintainer-owned catalog.
var topics, ruleTopics = loadTopics()

func loadTopics() (map[string]string, map[string]string) {
	data, err := docs.Files.ReadFile("handbook/topics.json")
	if err != nil {
		panic(err)
	}
	var catalog []struct {
		Topic string
		Page  string
		Rules []string
	}
	if err := json.Unmarshal(data, &catalog); err != nil {
		panic(err)
	}
	topics := map[string]string{}
	rules := map[string]string{}
	for _, entry := range catalog {
		if _, exists := topics[entry.Topic]; exists {
			panic("duplicate guide topic: " + entry.Topic)
		}
		topics[entry.Topic] = entry.Page
		for _, rule := range entry.Rules {
			if _, exists := rules[rule]; exists {
				panic("duplicate guide rule: " + rule)
			}
			rules[rule] = entry.Topic
		}
	}
	return topics, rules
}

// HasRule reports whether a diagnostic links to an authored explanation.
func HasRule(rule string) bool { return ruleTopics[rule] != "" }

func readPage(name string) ([]byte, error) {
	root, file, _ := strings.Cut(name, "/")
	switch root {
	case "docs":
		if path.Ext(file) == ".md" || file == "handbook/topics.json" {
			return docs.Files.ReadFile(file)
		}
	case "model":
		return modelschema.GuideFiles.ReadFile(file)
	case "specs":
		return specs.Files.ReadFile(file)
	case "examples":
		return examples.Files.ReadFile(file)
	}
	return nil, fs.ErrNotExist
}

// section returns an authored section, retaining nested headings and fenced examples.
func section(data []byte, id string) ([]byte, error) {
	md := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	doc := md.Parser().Parse(text.NewReader(data))
	start, level := 0, 0
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok {
			continue
		}
		line := bytes.LastIndexByte(data[:h.Pos()], '\n') + 1
		if level > 0 && h.Level <= level {
			return data[start:line], nil
		}
		anchor, ok := h.AttributeString("id")
		if ok && string(anchor.([]byte)) == id {
			start, level = line, h.Level
		}
	}
	if level > 0 {
		return data[start:], nil
	}
	return nil, fmt.Errorf("unknown guide section %q", id)
}

func names() []string {
	result := []string{"examples"}
	for name := range topics {
		result = append(result, name)
	}
	slices.Sort(result)
	return result
}

// Read returns Markdown and never consults the current project or the network.
func Read(query string) (string, error) {
	if query == "examples" {
		var list strings.Builder
		list.WriteString("# Worked examples\n\nOpen a file with `dotspec guide PATH`.\n\n")
		err := fs.WalkDir(examples.Files, ".", func(p string, entry fs.DirEntry, err error) error {
			if err == nil && !entry.IsDir() {
				fmt.Fprintf(&list, "- examples/%s\n", p)
			}
			return err
		})
		return list.String(), err
	}
	key := query
	if topic := ruleTopics[key]; topic != "" {
		key = topic
	}
	target := key
	if key == "" {
		target = topics["handbook"]
	} else if file := topics[key]; file != "" {
		target = file
	}
	name, fragment, _ := strings.Cut(target, "#")
	data, err := readPage(name)
	if err != nil {
		return "", fmt.Errorf("unknown guide topic or file %q; run dotspec guide for topics", query)
	}
	if fragment != "" {
		data, err = section(data, fragment)
		if err != nil {
			return "", err
		}
	}
	// Return example files verbatim so the reader can also save them as fixtures.
	if strings.HasPrefix(name, "examples/") || strings.HasPrefix(name, "specs/") || path.Ext(name) != ".md" {
		return string(data), nil
	}

	var out strings.Builder
	if query == "" {
		fmt.Fprintf(&out, "Topics: %s\n\n", strings.Join(names(), ", "))
	}
	out.Write(data)
	seen := map[string]bool{}
	for _, link := range specfile.MarkdownReferences(data).Links {
		u, err := url.Parse(link.Destination)
		if err != nil || u.IsAbs() || u.Host != "" || u.Path == "" || path.IsAbs(u.Path) {
			continue
		}
		file := path.Join(path.Dir(name), u.Path)
		if seen[file] {
			continue
		}
		if _, err := readPage(file); err != nil {
			continue
		}
		if len(seen) == 0 {
			out.WriteString("\n## Read linked pages offline\n\n")
		}
		seen[file] = true
		fmt.Fprintf(&out, "- `dotspec guide %s`\n", file)
	}
	return out.String(), nil
}
