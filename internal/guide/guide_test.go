package guide

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

func TestPublishedRoutes(t *testing.T) {
	for topic := range topics {
		if _, err := Read(topic); err != nil {
			t.Errorf("topic %q: %v", topic, err)
		}
	}
	for rule := range ruleTopics {
		if _, err := Read(rule); err != nil {
			t.Errorf("diagnostic %q: %v", rule, err)
		}
	}
}

func TestDocumentationLinks(t *testing.T) {
	anchors := map[string]map[string]bool{}
	sources := []string{"../../README.md", "../../CONTRIBUTING.md", "../../examples/todo/README.md"}
	for _, root := range []string{"../../docs", "../../examples/todo/views"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err == nil && !entry.IsDir() && filepath.Ext(path) == ".md" {
				sources = append(sources, path)
			}
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, source := range sources {
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		for _, link := range specfile.MarkdownReferences(body).Links {
			u, err := url.Parse(link.Destination)
			if err != nil {
				t.Fatal(err)
			}
			if u.IsAbs() || u.Host != "" {
				continue
			}
			target := source
			if u.Path != "" {
				target = filepath.Join(filepath.Dir(source), filepath.FromSlash(u.Path))
			}
			if _, err := os.Stat(target); err != nil {
				t.Errorf("%s: %v", source, err)
				continue
			}
			if u.Fragment == "" || filepath.Ext(target) != ".md" {
				continue
			}
			if anchors[target] == nil {
				data, err := os.ReadFile(target)
				if err != nil {
					t.Fatal(err)
				}
				anchors[target] = specfile.MarkdownReferences(data).Anchors
			}
			if !anchors[target][u.Fragment] {
				t.Errorf("%s: missing heading in %s", source, link.Destination)
			}
		}
	}
}
