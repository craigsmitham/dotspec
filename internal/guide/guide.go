package guide

import (
	"bytes"
	"fmt"
	"io/fs"
	"maps"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/docs"
	"github.com/craigsmitham/dotspec/examples"
	"github.com/craigsmitham/dotspec/internal/specfile"
	modelschema "github.com/craigsmitham/dotspec/model"
	"github.com/craigsmitham/dotspec/specs"
	"github.com/yuin/goldmark/ast"
)

var roots = map[string]fs.FS{"docs": docs.Files, "model": modelschema.GuideFiles, "examples": examples.Files, "specs": specs.Files}

func readFile(name string) ([]byte, error) {
	root, file, _ := strings.Cut(name, "/")
	if root == "examples" {
		if content, ok := examples.AppSpecs[file]; ok {
			return []byte(content), nil
		}
	}
	if root == "docs" && path.Ext(file) != ".md" && file != "cli-handbook/topics.json" {
		return nil, fs.ErrNotExist
	}
	if files := roots[root]; files != nil {
		return fs.ReadFile(files, file)
	}
	return nil, fs.ErrNotExist
}

func readDestination(source string) ([]byte, error) {
	file, anchor, section := strings.Cut(source, "#")
	data, err := readFile(file)
	if err != nil {
		return nil, err
	}
	if !section {
		return data, nil
	}
	md := specfile.NewMarkdown(data)
	start, level := 0, 0
	for n := md.Tree.FirstChild(); n != nil; n = n.NextSibling() {
		h, ok := n.(*ast.Heading)
		if !ok {
			continue
		}
		line := bytes.LastIndexByte(md.Body[:h.Pos()], '\n') + 1
		if level > 0 && h.Level <= level {
			return md.Body[start:line], nil
		}
		id, ok := h.AttributeString("id")
		if ok && string(id.([]byte)) == anchor {
			start, level = line, h.Level
		}
	}
	if level > 0 {
		return md.Body[start:], nil
	}
	return nil, fmt.Errorf("unknown section %q in %s", anchor, file)
}

func headingTitle(data []byte, source string) string {
	_, _, section := strings.Cut(source, "#")
	md := specfile.NewMarkdown(data)
	for n := md.Tree.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && (h.Level == 1 || section) {
			return strings.TrimSpace(string(h.Lines().Value(md.Body)))
		}
	}
	return path.Base(source)
}

func sourceFor(id string) string {
	file, anchor, section := strings.Cut(id, "#")
	// Check the documentation file first so IDs can contain dots.
	if _, err := fs.Stat(docs.Files, "cli-handbook/"+file+".md"); err == nil {
		file = docs.Root + "/" + file + ".md"
	} else if file == "topics.json" {
		file = docs.Root + "/" + file
	}
	if section {
		return file + "#" + anchor
	}
	return file
}

// Read uses only versioned, embedded sources. Explicit addresses never fall
// through to search, so a broken link is visible even if its words occur elsewhere.
func Read(query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		var out strings.Builder
		fmt.Fprintf(&out, "# Authoring handbook\n\nStart with the To-do walkthrough: dotspec docs %s\n\n", StartingPage("handbook"))
		for _, t := range Topics() {
			fmt.Fprintf(&out, "%s — %s\n  dotspec docs %s\n", t.Title, t.Description, t.Name)
		}
		return out.String(), nil
	}
	if name := handbook.names[strings.ToLower(query)]; name != "" {
		t := handbook.topics[name]
		var out strings.Builder
		fmt.Fprintf(&out, "# %s\n\n%s\n\n", t.Title, t.Description)
		for _, d := range t.Pages {
			title, err := d.PageTitle()
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&out, "%s (%s)\n  dotspec docs %s\n", title, d.Form(), d.ID())
		}
		if name == "examples" {
			files, err := bundledFiles("examples")
			if err != nil {
				return "", err
			}
			out.WriteString("\n## Bundled files\n")
			for _, file := range files {
				fmt.Fprintf(&out, "\n  dotspec docs %s", file)
			}
			out.WriteByte('\n')
		}
		return out.String(), nil
	}
	source := handbook.rules[query]
	if source == "" {
		source = sourceFor(query)
	}
	data, err := readDestination(source)
	if err == nil {
		if strings.HasPrefix(source, "examples/") || strings.HasPrefix(source, "specs/") {
			return string(data), nil
		}
		file, _, _ := strings.Cut(source, "#")
		if path.Ext(file) != ".md" {
			return string(data), nil
		}
		return linkedPage(source, data), nil
	}
	if strings.ContainsAny(query, "/#.\\") || handbook.rules[query] != "" {
		return "", fmt.Errorf("cannot open %q: %w; run dotspec docs to browse topics", query, err)
	}
	return search(query)
}

func linkedPage(source string, data []byte) string {
	var out strings.Builder
	out.Write(data)
	name, _, _ := strings.Cut(source, "#")
	seen := map[string]bool{}
	var offline, repository strings.Builder
	for _, link := range specfile.MarkdownReferences(data).Links {
		u, err := url.Parse(link.Destination)
		if err != nil || u.IsAbs() || u.Host != "" || path.IsAbs(u.Path) {
			continue
		}
		file := name
		if u.Path != "" {
			file = path.Join(path.Dir(name), u.Path)
		}
		id := pageID(file)
		if u.Fragment != "" {
			id += "#" + u.Fragment
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		if _, err := readFile(file); err == nil {
			fmt.Fprintf(&offline, "- `dotspec docs %s`\n", id)
		} else {
			fmt.Fprintf(&repository, "- `%s`\n", id)
		}
	}
	if offline.Len() > 0 {
		out.WriteString("\n## Read linked pages offline\n\n" + offline.String())
	}
	if repository.Len() > 0 {
		out.WriteString("\n## Repository references\n\n" + repository.String())
	}
	return out.String()
}

func bundledFiles(root string) ([]string, error) {
	var files []string
	err := fs.WalkDir(roots[root], ".", func(file string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (path.Ext(file) == ".md" || root == "examples") {
			files = append(files, root+"/"+file)
		}
		return err
	})
	if root == "examples" {
		for file := range examples.AppSpecs {
			files = append(files, "examples/"+file)
		}
	}
	slices.Sort(files)
	return files, err
}

// destinations includes searchable pages outside the curated topic collections.
func destinations() ([]string, error) {
	set := map[string]bool{}
	for _, root := range []string{"docs", "examples", "model", "specs"} {
		files, err := bundledFiles(root)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			set[file] = true
		}
	}
	for _, t := range handbook.topics {
		for _, d := range t.Pages {
			set[d.Page] = true
		}
	}
	return slices.Sorted(maps.Keys(set)), nil
}

func Complete(prefix string) []string {
	set := map[string]bool{}
	for name := range handbook.topics {
		set[name] = true
	}
	files, _ := destinations() // Only embedded, build-time filesystems.
	for _, file := range files {
		set[pageID(file)] = true
	}
	var result []string
	for _, id := range slices.Sorted(maps.Keys(set)) {
		if strings.HasPrefix(id, prefix) {
			result = append(result, id)
		}
	}
	return result
}

func search(query string) (string, error) {
	files, err := destinations()
	if err != nil {
		return "", err
	}
	subject := map[string]string{}
	for _, t := range handbook.topics {
		for _, d := range t.Pages {
			subject[d.Page] += " " + strings.Join(append([]string{t.Name, t.Title, t.Description}, t.Synonyms...), " ")
		}
	}
	terms := strings.Fields(strings.ToLower(query))
	matches := func(s string) bool {
		s = strings.ToLower(s)
		for _, term := range terms {
			if !strings.Contains(s, term) {
				return false
			}
		}
		return true
	}
	var primary, secondary strings.Builder
	for _, file := range files {
		data, err := readDestination(file)
		if err != nil {
			return "", err
		}
		d := Destination{Page: file}
		title := headingTitle(data, file)
		out := &primary
		if !matches(title + " " + subject[file]) {
			if !matches(string(data)) {
				continue
			}
			out = &secondary
		}
		fmt.Fprintf(out, "%s (%s)\n  dotspec docs %s\n", title, d.Form(), d.ID())
	}
	if primary.Len()+secondary.Len() == 0 {
		return "", fmt.Errorf("no documentation matches %q; run dotspec docs to browse topics", query)
	}
	return primary.String() + secondary.String(), nil
}
