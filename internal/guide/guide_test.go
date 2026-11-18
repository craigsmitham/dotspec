package guide

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func TestHandbookLinks(t *testing.T) {
	for _, source := range []string{"README.md", "maintaining.md", "project-conventions.md"} {
		dir := "../../docs/handbook"
		body, err := os.ReadFile(filepath.Join(dir, source))
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
			target := filepath.Join(dir, filepath.FromSlash(u.Path))
			if u.Path == "" {
				target = filepath.Join(dir, source)
			}
			data, err := os.ReadFile(target)
			if err != nil {
				t.Errorf("%s: %v", source, err)
				continue
			}
			if u.Fragment != "" && strings.HasSuffix(target, ".md") && !specfile.MarkdownReferences(data).Anchors[u.Fragment] {
				t.Errorf("%s: missing heading in %s", source, link.Destination)
			}
		}
	}
}
