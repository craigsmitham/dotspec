// Package guide resolves topics and reads documentation bundled with the CLI.
package guide

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/craigsmitham/dotspec/docs"
)

// A Topic groups destinations in editorial order. The first is the starting
// point for a contributor who has just created a draft of this type.
type Topic struct {
	Name        string `json:"topic"`
	Title       string
	Description string
	Synonyms    []string
	Pages       []Destination
}

// Rules map stable diagnostic IDs to heading fragments within this page.
// An empty fragment uses the destination as written, including its own fragment.
type Destination struct {
	Page  string
	Rules map[string]string
}

func (d Destination) ID() string { return pageID(d.Page) }

func pageID(source string) string {
	file, anchor, section := strings.Cut(source, "#")
	if strings.HasPrefix(file, docs.Root+"/") {
		file = strings.TrimSuffix(strings.TrimPrefix(file, docs.Root+"/"), ".md")
	}
	if section {
		return file + "#" + anchor
	}
	return file
}

func (d Destination) Form() string {
	switch {
	case strings.HasPrefix(d.Page, docs.Root+"/explanations/"):
		return "Explanation"
	case strings.HasPrefix(d.Page, docs.Root+"/tutorials/"):
		return "Tutorial"
	case strings.HasPrefix(d.Page, docs.Root+"/guides/"):
		return "How-to"
	case strings.HasPrefix(d.Page, "examples/"), strings.HasPrefix(d.Page, "specs/"):
		return "Example"
	case strings.HasPrefix(d.Page, docs.Root+"/reference/"), strings.HasPrefix(d.Page, "model/"):
		return "Reference"
	default:
		return "Handbook"
	}
}

type catalog struct {
	topics map[string]Topic
	names  map[string]string
	rules  map[string]string
}

var handbook = loadCatalog()

func loadCatalog() catalog {
	data, err := docs.Files.ReadFile("cli-handbook/topics.json")
	if err != nil {
		panic(err)
	}
	var entries []Topic
	if err := json.Unmarshal(data, &entries); err != nil {
		panic(err)
	}
	c, err := indexTopics(entries)
	if err != nil {
		panic(err)
	}
	return c
}

func indexTopics(entries []Topic) (catalog, error) {
	c := catalog{map[string]Topic{}, map[string]string{}, map[string]string{}}
	used := make(map[string]bool)
	// Pages can be shared across topics, but a topic or rule must never
	// intercept the command printed for a page.
	for _, topic := range entries {
		for _, page := range topic.Pages {
			used[strings.ToLower(page.ID())] = true
		}
	}
	for _, t := range entries {
		if t.Name == "" || t.Title == "" || t.Description == "" || len(t.Pages) == 0 {
			return c, fmt.Errorf("incomplete documentation topic %q", t.Name)
		}
		for _, name := range append([]string{t.Name}, t.Synonyms...) {
			key := strings.ToLower(name)
			if key == "" || used[key] {
				return c, fmt.Errorf("invalid or duplicate documentation lookup %q", name)
			}
			used[key] = true
			c.names[key] = t.Name
		}
		c.topics[t.Name] = t
		for _, d := range t.Pages {
			for rule, fragment := range d.Rules {
				key := strings.ToLower(rule)
				if key == "" || used[key] {
					return c, fmt.Errorf("duplicate documentation lookup %q", rule)
				}
				used[key] = true
				target := d.Page
				if fragment != "" {
					file, _, _ := strings.Cut(target, "#")
					target = file + "#" + fragment
				}
				c.rules[rule] = target
			}
		}
	}
	return c, nil
}

func Topics() []Topic {
	var result []Topic
	for _, name := range slices.Sorted(maps.Keys(handbook.topics)) {
		result = append(result, handbook.topics[name])
	}
	return result
}

func HasRule(rule string) bool { return handbook.rules[rule] != "" }

func StartingPage(topic string) string { return handbook.topics[topic].Pages[0].ID() }

// PageTitle uses the authored heading for both whole pages and section links.
func (d Destination) PageTitle() (string, error) {
	data, err := readDestination(d.Page)
	if err != nil {
		return "", err
	}
	return headingTitle(data, d.Page), nil
}
