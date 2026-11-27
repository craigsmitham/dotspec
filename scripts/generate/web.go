package main

import (
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

var publishedGuides = []string{"describe-a-system", "write-an-assessable-requirement", "list-a-companion"}

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
	for _, name := range []string{"init", "new requirement", "build", "check"} {
		cmd, remaining, err := commands.Find(strings.Fields(name))
		if err != nil || len(remaining) != 0 || cmd.CommandPath() != "dotspec "+name {
			return nil, fmt.Errorf("missing command %q in the CLI reference", name)
		}
		var body strings.Builder
		writeCommandBody(&body, cmd)
		body.WriteString("See the [full CLI reference](" + githubSource + "docs/reference/cli.md#dotspec) for global flags.\n")
		pages = append(pages, webPage("cli/"+strings.ReplaceAll(name, " ", "-"), cmd.CommandPath(), "Cobra command tree", body.String()))
	}
	return pages, nil
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
		if file == handbookDir+slug+".md" && slices.Contains(publishedGuides, slug) {
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
