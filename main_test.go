package main

import (
	"io/fs"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/docs"
	"github.com/craigsmitham/dotspec/internal/guide"
	"github.com/craigsmitham/dotspec/internal/specfile"
)

func TestHandbookPublicationBoundary(t *testing.T) {
	for _, id := range []string{"docs/README.md", "maintainers/releasing", "docs/maintainers/handbook.md", "AGENTS.md", "CONTRIBUTING.md", "model/README.md", "examples/todo/COVERAGE.md"} {
		if _, err := guide.Read(id); err == nil {
			t.Errorf("repository-only address opened: %s", id)
		}
	}
	// A page or section can be cited hundreds of times in the model reference.
	opened := map[string]error{}
	open := func(id string) error {
		if err, known := opened[id]; known {
			return err
		}
		_, err := guide.Read(id)
		opened[id] = err
		return err
	}
	err := fs.WalkDir(docs.Files, ".", func(file string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		if !strings.HasPrefix(file, "cli-handbook/") || path.Base(file) == "AGENTS.md" {
			t.Errorf("repository file embedded: %s", file)
		}
		if path.Ext(file) != ".md" {
			return nil
		}
		source := "docs/" + file
		id := (guide.Destination{Page: source}).ID()
		rendered, err := guide.Read(id)
		if err != nil {
			return err
		}
		raw, err := docs.Files.ReadFile(file)
		if err != nil {
			return err
		}
		for _, ref := range specfile.MarkdownReferences(raw).Links {
			u, err := url.Parse(ref.Destination)
			if err != nil {
				return err
			}
			if u.IsAbs() || u.Host != "" {
				continue
			}
			target := source
			if u.Path != "" {
				target = path.Join(path.Dir(source), u.Path)
			}
			if _, err := os.Stat(target); err != nil {
				return err
			}
			id := (guide.Destination{Page: target}).ID()
			if u.Fragment != "" {
				id += "#" + u.Fragment
			}
			command := "`dotspec docs " + id + "`"
			// An embedded file's fragment must open too. Repository citations
			// carry the source path without promising an offline command.
			fileErr := open((guide.Destination{Page: target}).ID())
			if fileErr != nil {
				if strings.Contains(rendered, command) || !strings.Contains(rendered, "`"+id+"`") {
					t.Errorf("%s: repository citation %s", source, id)
				}
			} else if err := open(id); err != nil || !strings.Contains(rendered, command) {
				t.Errorf("%s: offline citation %s: %v", source, id, err)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
