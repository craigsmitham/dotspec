package specdoc

import "strings"

// Shared section definitions bind authored content and directory membership.
// Each layout takes value copies and supplies their order, nesting and navigation.
var (
	overviewLead      = overviewSection{sources: []string{"description"}}
	framing           = overviewSection{title: "Framing", level: 2, sources: []string{"framing"}}
	purpose           = overviewSection{title: "Purpose", level: 2, sources: []string{"purpose"}}
	goals             = overviewSection{title: "Goals", level: 2, sources: []string{"goals"}}
	stakeholders      = overviewSection{title: "Stakeholders", level: 2, sources: []string{"stakeholder_roles"}}
	features          = overviewDirectory("Features")
	capabilities      = overviewDirectory("Capabilities", "capabilities")
	concepts          = overviewDirectory("Concepts")
	conceptUses       = overviewDirectory("Concept uses", "concept_uses")
	synchronizations  = overviewDirectory("Synchronizations", "synchronizations")
	qualityAttributes = overviewSection{title: "Quality Attributes", collection: "Quality attributes", level: 2}
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

func overviewDirectory(title string, sources ...string) overviewSection {
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

var systemContextFields = []string{"scope", "out_of_scope", "environment", "assumptions", "constraints"}

func systemOutline(architecture []overviewSection) []overviewSection {
	sections := []overviewSection{overviewLead}
	sections = append(sections, navGroup("Intent", framing, purpose, goals, stakeholders, overviewDirectory("Jobs"))...)
	var context []overviewSection
	for _, field := range systemContextFields {
		level := 2
		if field == "out_of_scope" {
			level = 3
		}
		title := sentenceCase(strings.ReplaceAll(field, "_", " "))
		context = append(context, overviewSection{title: title, level: level, sources: []string{field}})
	}
	context = append(context, overviewDirectory("Business rules"))
	sections = append(sections, navGroup("Context", context...)...)
	sections = append(sections, navGroup("Behavior", features, capabilities, overviewDirectory("Use cases"))...)
	sections = append(sections, qualityAttributes, overviewDirectory("Requirements"))
	sections = append(sections, navGroup("Design", concepts, conceptUses, synchronizations)...)
	return append(sections, navGroup("Design", architecture...)...)
}

func componentOutline(architecture []overviewSection) []overviewSection {
	sections := []overviewSection{
		overviewLead, framing, purpose, goals, stakeholders,
		{title: "Responsibilities", level: 2, sources: []string{"responsibilities"}},
		{title: "Scope", level: 2, sources: []string{"boundary_description", "out_of_scope"}},
	}
	sections = append(sections, architecture...)
	provided := capabilities
	provided.level = 3
	return append(sections, provided, features, concepts, conceptUses, synchronizations, qualityAttributes)
}
