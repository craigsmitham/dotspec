package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/dotspec-md/dotspec/internal/model"
	"github.com/dotspec-md/dotspec/internal/specfile"
	"github.com/spf13/cobra"
)

const handbookDir = "docs/cli-handbook/guides/"
const githubSource = "https://github.com/dotspec-md/dotspec/blob/main/"

var publishedGuides = []string{"describe-a-system", "write-an-assessable-requirement", "catalog-a-specification"}

func webPages(commands *cobra.Command) ([]generatedFile, error) {
	var pages []generatedFile
	for _, slug := range publishedGuides {
		source := handbookDir + slug + ".md"
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}
		body, err := webGuide(source, specfile.Normalize(data))
		if err != nil {
			return nil, fmt.Errorf("publish %s: %w", source, err)
		}
		heading, body, ok := strings.Cut(body, "\n")
		if !ok || !strings.HasPrefix(heading, "# ") {
			return nil, fmt.Errorf("missing guide title in %s", source)
		}
		pages = append(pages, webPage("guides/"+slug, heading[2:], source, body))
	}
	for _, file := range []string{restartReport, assumptionReport} {
		page, err := publishReport(file)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	walkthrough, err := checkWalkthrough()
	if err != nil {
		return nil, err
	}
	var visit func(*cobra.Command) any
	visit = func(cmd *cobra.Command) any {
		route := commandRoute(cmd)
		var body strings.Builder
		writeCommandBody(&body, cmd)
		switch route {
		case "cli/check":
			body.WriteString(walkthrough)
		case "cli/build":
			body.WriteString("[Repair a stale specification](/troubleshooting/spec-stale).\n")
		}
		if cmd.HasParent() {
			body.WriteString("See [global flags](/cli/overview) for options available to every command.\n")
		}
		var children []any
		for _, child := range cmd.Commands() {
			if child.Hidden {
				continue
			}
			if len(children) == 0 {
				body.WriteString("\n## Choose a command\n\n")
			}
			fmt.Fprintf(&body, "- [%s](/%s) — %s\n", child.CommandPath(), commandRoute(child), child.Short)
			children = append(children, visit(child))
		}
		pages = append(pages, webPage(route, cmd.CommandPath(), "Cobra command tree", body.String()))
		if len(children) == 0 {
			return route
		}
		label := cmd.Name()
		if !cmd.HasParent() {
			label = "CLI reference"
		}
		return map[string]any{"group": label, "root": route, "pages": children}
	}
	group := visit(commands)
	const config = "mintlify/docs.json"
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	var site map[string]json.RawMessage
	if err = json.Unmarshal(data, &site); err != nil {
		return nil, err
	}
	var nav map[string]json.RawMessage
	if err = json.Unmarshal(site["navigation"], &nav); err != nil {
		return nil, err
	}
	var groups []json.RawMessage
	if err = json.Unmarshal(nav["groups"], &groups); err != nil {
		return nil, err
	}
	found := false
	for i, raw := range groups {
		var meta struct{ Group string }
		if err = json.Unmarshal(raw, &meta); err != nil {
			return nil, err
		}
		if meta.Group == "CLI reference" {
			groups[i], err = json.Marshal(group)
			if err != nil {
				return nil, err
			}
			found = true
		}
	}
	if !found {
		return nil, errors.New("missing CLI reference navigation group")
	}
	nav["groups"], err = json.Marshal(groups)
	if err != nil {
		return nil, err
	}
	site["navigation"], err = json.Marshal(nav)
	if err != nil {
		return nil, err
	}
	data, err = json.MarshalIndent(site, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(pages, generatedFile{config, string(data) + "\n"}), nil
}

func commandRoute(cmd *cobra.Command) string {
	if !cmd.HasParent() {
		return "cli/overview"
	}
	return "cli/" + strings.ReplaceAll(strings.TrimPrefix(cmd.CommandPath(), cmd.Root().Name()+" "), " ", "-")
}

func webPage(route, title, source, body string) generatedFile {
	content := fmt.Sprintf("---\ntitle: %q\n---\n\n{/* Generated from %s by mise run generate. */}\n\n%s\n", title, source, strings.TrimSpace(body))
	return generatedFile{"mintlify/" + route + ".mdx", content}
}

// Change link destinations without touching examples, labels or heading anchors.
// The content parser resolves reference-style links and gives byte ranges even
// when a link contains an image.
func webGuide(source string, data []byte) (string, error) {
	type replacement struct {
		start, end int
		suffix     string
	}
	var replacements []replacement
	for _, node := range model.ParseContent(source, data).Nodes {
		kind := model.String(node, "node_kind")
		if kind != "link" && kind != "image" {
			continue
		}
		dest, err := url.Parse(model.String(node, "destination"))
		if err != nil {
			return "", err
		}
		if dest.IsAbs() || dest.Host != "" || dest.Path == "" {
			continue
		}
		file := path.Join(path.Dir(source), dest.Path)
		if _, err := os.Stat(file); err != nil {
			return "", err
		}
		slug := strings.TrimSuffix(path.Base(file), ".md")
		if route, ok := reportRoutes[file]; ok {
			dest.Path = "/" + route
		} else if file == handbookDir+slug+".md" && slices.Contains(publishedGuides, slug) {
			dest.Path = "/guides/" + slug
		} else {
			dest.Scheme, dest.Host, dest.Path = "https", "github.com", "/dotspec-md/dotspec/blob/main/"+file
		}
		span := model.Object(node["span"])
		start, end := model.Number(span["start"]), model.Number(span["end"])
		if kind == "image" {
			start++
		}
		bracket := specfile.LinkLabelEnd(data[:end], start)
		if bracket < 0 {
			return "", errors.New("missing link label")
		}
		suffix := "(" + specfile.FormatDestination(dest.String())
		if title := model.String(node, "title"); title != "" {
			suffix += " " + specfile.FormatTitle(title)
		}
		replacements = append(replacements, replacement{bracket + 1, end, suffix + ")"})
	}
	slices.SortFunc(replacements, func(a, b replacement) int { return a.start - b.start })
	var out strings.Builder
	var offset int
	for _, repl := range replacements {
		out.Write(data[offset:repl.start])
		out.WriteString(repl.suffix)
		offset = repl.end
	}
	out.Write(data[offset:])
	return out.String(), nil
}
