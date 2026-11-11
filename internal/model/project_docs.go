package model

import (
	"fmt"
	"sort"
	"strings"
)

func projectReading(s *Snapshot, text string) string {
	project := Object(s.Scope["project"])
	if project == nil {
		return text
	}
	// Domain context is distinct from merely acquiring the workspace manifest.
	sourceProjects := map[string]string{}
	for _, source := range s.Sources {
		sourceProjects[String(source, "id")] = String(Object(source["project"]), "name")
	}
	names := map[string]bool{}
	for _, element := range s.Elements {
		if name := sourceProjects[String(element, "source_document")]; name != "" && name != String(project, "name") {
			names[name] = true
		}
	}
	external := []string{}
	for name := range names {
		external = append(external, name)
	}
	sort.Strings(external)
	note := fmt.Sprintf("Specification project: **%s**. Project membership does not imply a System relationship.\n\n", String(project, "name"))
	if len(external) > 0 {
		note += "Explicitly acquired domain context from other projects: " + strings.Join(external, ", ") + ". Those projects are not recursively scanned.\n\n"
	}
	heading, rest, ok := strings.Cut(text, "\n\n")
	if !ok {
		return note + text
	}
	return heading + "\n\n" + note + rest
}
