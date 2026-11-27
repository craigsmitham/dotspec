package specfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/craigsmitham/dotspec/internal/diagnostic"
	schema "github.com/craigsmitham/dotspec/model"
)

// A level heading introduces the domain slot it names in sentence case: lowercase
// the title and join its words with underscores, so "Out of scope" introduces
// out_of_scope and "Out Of Scope" introduces nothing.
// Only the slots listed here are authored as Markdown sections; the others are
// set from frontmatter, derived, or owned by a structured parent. Nested
// headings inside ordinary Markdown fields are local to that prose.
var sectionSlots = setOf(
	"description", "part_of", "purpose", "overview", "scope", "out_of_scope",
	"project_references", "environment", "constraints", "responsibilities", "contracts",
	"architecture", "implementation", "interacts_with", "enablers", "groups", "depends_on",
	"provided_by", "supports", "involves", "derived_from", "relies_on", "refines", "applies_to",
	"adopts", "evidence", "related_specifications", "open_questions",
	"performer", "situation", "motivation", "intended_progress",
	"goal", "primary_actor", "goal_level", "trigger", "supporting_actors", "stakeholder_interests",
	"main_success_scenario", "sequence", "preconditions", "minimal_guarantees", "success_guarantees",
	"operational_principles", "effects_table", "rationale", "criterion", "exceptions", "definition",
	"sources", "edition", "basis", "supporting_content", "configuration", "dependencies", "counterpart",
	"when", "where", "then", "semantics", "inputs", "outputs", "initiation", "effects", "rejections",
	"guards", "caller_assumptions", "condition", "handling",
	"goals", "assumptions", "outcomes", "actions", "required_interfaces", "provided_interfaces",
	"glossary", "extensions", "concept_uses", "synchronizations",
)

// irregularHeadings name the slot of a title that departs from the convention.
var irregularHeadings = map[string]string{
	"Stakeholders": "stakeholder_roles", "Authority": "authority_and_evidence", "State": "abstract_state",
	"Boundary": "boundary_description", "System boundary": "boundary_description", "Job": "statement",
	"Specification system": "specification_system_content",
}

// headingOwners restrict a title to the classes that author it as a section.
// Other classes declaring the same slot set it some other way.
var headingOwners = map[string][]string{
	"Job": {"Job"}, "Statement": {"Requirement", "BusinessRule"}, "Boundary": {"Component"},
	"System boundary": {"UseCase"}, "Steps": {"Job"}, "Applicability": {"Assumption", "AdoptedProvision"},
	"Quality attributes": {"QualityModel"},
}

func setOf(items ...string) map[string]bool {
	set := map[string]bool{}
	for _, item := range items {
		set[item] = true
	}
	return set
}

// sentenceCase spells a slot as its heading: "out_of_scope" is "Out of scope".
func sentenceCase(slot string) string {
	if slot == "" {
		return ""
	}
	return strings.ToUpper(slot[:1]) + strings.ReplaceAll(slot[1:], "_", " ")
}

// SectionSlot names the slot a level heading titled title introduces in class,
// and reports whether class accepts that heading as one of its sections.
func SectionSlot(class, title string) (string, bool) {
	schema, ok := classSlots()[class]
	if !ok {
		return "", false
	}
	slot, irregular := irregularHeadings[title]
	if !irregular {
		slot = strings.ReplaceAll(strings.ToLower(title), " ", "_")
		if title != sentenceCase(slot) {
			return "", false
		}
	}
	owners, restricted := headingOwners[title]
	if restricted && !slices.Contains(owners, class) || !restricted && !irregular && !sectionSlots[slot] {
		return "", false
	}
	if _, ok := schema.Slots[slot]; !ok {
		return "", false
	}
	return slot, true
}

const verificationProcedureAdvice = "List each verification procedure as a companion with role: verifies in the frontmatter, and keep any account of what it exercises under Evidence."

var collectionKinds = map[string]string{
	"Adopts": "AdoptedProvision", "Stakeholders": "SystemStakeholderRole",
	"Goals": "Goal", "Assumptions": "Assumption", "Steps": "JobStep", "Outcomes": "DesiredOutcome",
	"Actions": "ConceptAction", "Required interfaces": "Interface", "Provided interfaces": "Interface",
	"Quality attributes": "QualityModelAttribute", "Glossary": "GlossaryTerm",
	"Extensions": "UseCaseExtension", "Concept uses": "ConceptUse", "Synchronizations": "Synchronization",
}

// IsCollectionHeading reports whether a source heading introduces named entries.
func IsCollectionHeading(title string) bool {
	return collectionKinds[title] != ""
}

var classSlots = sync.OnceValue(func() map[string]struct{ Slots map[string]json.RawMessage } {
	var catalog struct {
		Classes map[string]struct{ Slots map[string]json.RawMessage }
	}
	if err := json.Unmarshal(schema.CatalogJSON, &catalog); err != nil {
		panic(err)
	}
	return catalog.Classes
})

var documentRootClasses = map[string]string{
	"business-rule":     "BusinessRule",
	"capability":        "Capability",
	"component":         "Component",
	"concept":           "Concept",
	"feature":           "Feature",
	"job":               "Job",
	"quality-attribute": "SystemQualityAttribute",
	"quality-model":     "QualityModel",
	"requirement":       "Requirement",
	"system":            "System",
	"use-case":          "UseCase",
}

// DocumentRootClass identifies the primary subject of a built-in authoring form.
func DocumentRootClass(kind string) string { return documentRootClasses[kind] }

// RootKind names the built-in authoring form whose primary subject is class.
func RootKind(class string) (string, bool) {
	for kind, root := range documentRootClasses {
		if root == class {
			return kind, true
		}
	}
	return "", false
}

var reportTypes = sync.OnceValue(func() map[string]string {
	type slot struct {
		Fixed  string   `json:"equals_string"`
		Ranges []string `json:"reference_ranges"`
	}
	types := map[string]string{}
	for _, class := range classSlots() {
		var documents, typ slot
		if class.Slots["documents"] == nil || class.Slots["type"] == nil {
			continue
		}
		if err := json.Unmarshal(class.Slots["documents"], &documents); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(class.Slots["type"], &typ); err != nil {
			panic(err)
		}
		if len(documents.Ranges) == 1 && typ.Fixed != "" {
			types[documents.Ranges[0]] = typ.Fixed
		}
	}
	return types
})

// ReportType returns the fixed type of the report class that documents kind.
func ReportType(kind string) string { return reportTypes()[kind] }

func (o *Outline) checkFields(class string, sections []*Section, level int) error {
	// Independently registered types own their authoring contracts.
	if _, builtIn := classSlots()[class]; !builtIn {
		return nil
	}
	for _, s := range sections {
		_, known := SectionSlot(class, s.Title)
		loc := o.Location(s.Start, s.Title)
		if s.Title == "Verification procedures" {
			return diagnostic.Error("structure.verification-procedure", verificationProcedureAdvice, loc)
		}
		if s.Title == "Obligations" || s.Title == "Requirements" {
			return diagnostic.Error("structure.inline-requirement", "Declare each requirement in a *.requirement.spec.md file and link to it here. Keep scenario outcomes in the use case.", loc)
		}
		if !known {
			return diagnostic.Error("structure.unknown-section", fmt.Sprintf("Unrecognized %s section %q. Move its prose into the field it qualifies or use a declared section.", class, s.Title), loc)
		}
		if s.Level != level {
			return diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-%d heading", s.Title, level), loc)
		}
		if class == "Requirement" && s.Title == "Criterion" {
			for _, sub := range s.Children {
				if sub.Title == "Verification procedures" {
					return diagnostic.Error("structure.verification-procedure", verificationProcedureAdvice, o.Location(sub.Start, sub.Title))
				}
			}
		}
		member := collectionKinds[s.Title]
		if member == "" {
			continue
		}
		if err := o.collectionIntro(s); err != nil {
			return err
		}
		for _, child := range s.Children {
			if child.Level != level+1 || strings.TrimSpace(child.Title) == "" {
				return diagnostic.Error("structure.member-heading", "Use a named heading immediately below its collection.", o.Location(child.Start, s.Title))
			}
			if member == "AdoptedProvision" {
				if len(o.sectionLinks(child)) != 1 {
					return diagnostic.Error("structure.adoption", "Identify exactly one provision with a link under the named entry.", o.Location(child.Start, child.Title))
				}
				fields := map[string]int{}
				for _, sub := range child.Children {
					if o.Body(sub) != "" {
						fields[sub.Title]++
					}
				}
				if len(child.Children) != 2 || fields["Edition"] != 1 || fields["Applicability"] != 1 {
					return diagnostic.Error("structure.adoption", "Give the entry one Edition and one Applicability subsection, each with its value.", o.Location(child.Start, child.Title))
				}
			}
			// Definitions and goal statements own their subordinate Markdown.
			if slices.Contains([]string{"QualityModelAttribute", "GlossaryTerm", "Goal", "SystemStakeholderRole"}, member) {
				continue
			}
			if err := o.checkFields(member, child.Children, level+2); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Outline) collectionIntro(s *Section) error {
	intro := o.Source[s.ContentStart:o.IntroEnd(s)]
	trimmed := bytes.TrimLeft(intro, " \t\r\n")
	if len(trimmed) > 0 {
		return diagnostic.Error("structure.collection-content", fmt.Sprintf("%s contains named entries. Move this passage into its entry or the field it qualifies.", s.Title), o.Location(s.ContentStart+len(intro)-len(trimmed), s.Title))
	}
	return nil
}

func (o *Outline) sectionLinks(s *Section) []Reference {
	var refs ReferenceDocument
	for _, n := range o.Nodes(s) {
		outlineLinks(o, n, &refs)
	}
	return refs.Links
}
