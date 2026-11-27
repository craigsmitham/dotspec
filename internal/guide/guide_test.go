package guide

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/craigsmitham/dotspec/internal/specfile"
)

func TestPublishedRoutes(t *testing.T) {
	files, err := destinations()
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		id := pageID(file)
		if handbook.names[strings.ToLower(id)] != "" || HasRule(id) {
			t.Errorf("topic or diagnostic shadows page address %q", id)
		}
	}
	for topic := range handbook.topics {
		if _, err := Read(topic); err != nil {
			t.Errorf("topic %q: %v", topic, err)
		}
	}
	for rule := range handbook.rules {
		if _, err := Read(rule); err != nil {
			t.Errorf("diagnostic %q: %v", rule, err)
		}
	}
}

func TestDocumentationLinks(t *testing.T) {
	anchors := map[string]map[string]bool{}
	sources := []string{"../../README.md", "../../CONTRIBUTING.md", "../../examples/todo/README.md"}
	for _, root := range []string{"../../docs", "../../examples/todo/specification"} {
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

func TestExampleMatchesRepositoryBytes(t *testing.T) {
	const file = "examples/todo/specs/jobs/personal-commitments.job.spec.md"
	want, err := os.ReadFile("../../" + file)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Read(file)
	if err != nil || got != string(want) {
		t.Fatalf("example differs from %s: %v", file, err)
	}
}

func TestSearchRanksTopicAboveBody(t *testing.T) {
	results, err := Read("success criteria")
	if err != nil {
		t.Fatal(err)
	}
	guide := strings.Index(results, "dotspec docs guides/write-a-job\n")
	index := strings.Index(results, "dotspec docs README\n")
	if guide < 0 || index <= guide {
		t.Fatalf("expected the Jobs topic's how-to above the index page:\n%s", results)
	}
}

func TestTopicLookupNamesAreUnique(t *testing.T) {
	topics := []Topic{{Name: "job", Title: "Jobs", Description: "Write a job.", Synonyms: []string{"jtbd"}, Pages: []Destination{{Page: "docs/cli-handbook/guides/write-a-job.md"}}}}
	for _, collision := range []string{"job", "JTBD", "guides/write-a-job"} {
		conflict := topics[0]
		conflict.Name = collision
		conflict.Synonyms = nil
		if _, err := indexTopics(append(topics, conflict)); err == nil {
			t.Errorf("accepted duplicate lookup %q", collision)
		}
	}
}

func TestBundledQualityModelPages(t *testing.T) {
	page, err := Read("quality-models/iso-25010-2023")
	if err != nil || !strings.HasPrefix(page, "# ISO/IEC 25010:2023") || !strings.Contains(page, "**ID:** recoverability") {
		t.Fatalf("page: %v\n%s", err, page)
	}
	section, err := Read("quality-models/iso-25010-2023#reliability")
	if err != nil || !strings.HasPrefix(section, "### Reliability") || strings.Contains(section, "### Security") {
		t.Fatalf("section: %v\n%s", err, section)
	}
	found, err := Read("iso 25010")
	if err != nil || !strings.Contains(found, "dotspec docs quality-models/iso-25010-2023\n") {
		t.Fatalf("search: %v\n%s", err, found)
	}
	for _, missing := range []string{"quality-models/unknown", "quality-models/../store.go", "quality-models/iso-25010-2023.quality-model.md"} {
		if _, err := Read(missing); err == nil {
			t.Errorf("opened %q", missing)
		}
	}
}
