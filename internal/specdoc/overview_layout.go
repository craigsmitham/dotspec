package specdoc

import "strings"

// Shared section definitions bind authored content and directory membership.
// Each layout takes value copies and supplies their order, nesting and navigation.
var (
	overviewLead      = overviewSection{sources: []string{"description"}}
	purpose           = overviewSection{title: "Purpose", level: 2, sources: []string{"purpose"}}
	overview          = overviewSection{title: "Overview", level: 2, sources: []string{"overview"}}
	goals             = overviewSection{title: "Goals", level: 2, sources: []string{"goals"}}
	stakeholders      = overviewSection{title: "Stakeholders", level: 2, sources: []string{"stakeholder_roles"}}
	features          = overviewDirectory("Feature")
	capabilities      = overviewDirectory("Capability")
	concepts          = overviewDirectory("Concept", "concept_uses")
	synchronizations  = overviewDirectory("Synchronization", "synchronizations")
	qualityAttributes = overviewSection{title: "Quality Attributes", collection: collectionTitle("SystemQualityAttribute"), level: 2}
)

var overviewReferences = []overviewSection{
	{title: "Adopts", level: 2, sources: []string{"adopts"}},
	{title: "Contracts", level: 2, sources: []string{"contracts"}},
	{title: "Implementation", level: 2, sources: []string{"implementation"}},
	{title: "Open questions", level: 2, sources: []string{"open_questions"}},
	{title: "Glossary", level: 2, sources: []string{"glossary"}},
	{title: "Project references", level: 2, sources: []string{"project_references"}},
	{title: "Reference", level: 2, sources: []string{"related_specifications", "specification_system", "authority_and_evidence", "evidence", "supports", "derived_from", "relies_on"}},
}

// A directory section takes its title from the collection that lists kind, so
// the outline and the collection cannot drift apart.
func overviewDirectory(kind string, sources ...string) overviewSection {
	title := collectionTitle(kind)
	return overviewSection{title: title, collection: title, level: 2, sources: sources}
}

// Only top-level sections have entries in the navigation map.
func navGroup(label string, sections ...overviewSection) []overviewSection {
	for i := range sections {
		if sections[i].level == 2 {
			sections[i].group = label
		}
	}
	return sections
}

func systemOutline(architecture []overviewSection) []overviewSection {
	sections := []overviewSection{overviewLead}
	sections = append(sections, navGroup("Intent", purpose, overview, goals, stakeholders, overviewDirectory("Job"))...)
	var context []overviewSection
	// A top-level field travels with its subsections when the reading order changes.
	contextFields := [][]string{
		{"scope", "out_of_scope"},
		{"environment"},
		{"assumptions"},
		{"constraints"},
	}
	for _, fields := range contextFields {
		for i, field := range fields {
			level := 2
			if i > 0 {
				level = 3
			}
			title := sentenceCase(strings.ReplaceAll(field, "_", " "))
			context = append(context, overviewSection{title: title, level: level, sources: []string{field}, omitEmptyPassages: true})
		}
	}
	context = append(context, overviewDirectory("BusinessRule"))
	sections = append(sections, navGroup("Context", context...)...)
	sections = append(sections, navGroup("Behavior", features, capabilities, overviewDirectory("UseCase"))...)
	sections = append(sections, qualityAttributes, overviewDirectory("Requirement"))
	sections = append(sections, navGroup("Design", concepts, synchronizations)...)
	return append(sections, navGroup("Design", architecture...)...)
}

func componentOutline(architecture []overviewSection) []overviewSection {
	sections := []overviewSection{
		overviewLead, purpose, overview,
		{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
		{title: "Scope", level: 2, sources: []string{"boundary_description", "out_of_scope"}},
	}
	sections = append(sections, architecture...)
	provided := capabilities
	provided.level = 3
	return append(sections, provided, features, concepts, synchronizations, qualityAttributes)
}
