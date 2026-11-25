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

// A heading introduces a field only when its owner declares that domain slot.
// Nested headings inside ordinary Markdown fields are local to that prose.
var headingFields = map[string]string{
	"Description": "description", "Framing": "framing", "Part of": "part_of", "Purpose": "purpose", "Scope": "scope", "Out of scope": "out_of_scope",
	"Project references": "project_references", "Environment": "environment", "Stakeholders": "stakeholders", "Constraints": "constraints",
	"Responsibilities": "responsibilities", "Boundary": "boundary_description", "Contracts": "contracts", "Architecture": "architecture", "Implementation": "implementation",
	"Interacts with": "interacts_with", "Specification system": "specification_system", "Enablers": "enablers",
	"Groups": "groups", "Depends on": "depends_on", "Provided by": "provided_by", "Supports": "supports", "Involves": "involves",
	"Derived from": "derived_from", "Relies on": "relies_on", "Refines": "refines", "Applies to": "applies_to",
	"Adopts": "adopts", "Verification procedures": "verification_procedures", "Evidence": "evidence", "Authority": "authority_and_evidence", "Related specifications": "related_specifications", "Open questions": "open_questions",
	"Job": "statement", "Performer": "performer", "Situation": "situation", "Motivation": "motivation", "Intended progress": "intended_progress",
	"Goal": "goal", "System boundary": "boundary_description", "Primary actor": "primary_actor", "Goal level": "goal_level", "Trigger": "trigger",
	"Supporting actors": "supporting_actors", "Stakeholder interests": "stakeholder_interests", "Main success scenario": "main_success_scenario",
	"Sequence": "sequence", "Preconditions": "preconditions", "Minimal guarantees": "minimal_guarantees", "Success guarantees": "success_guarantees",
	"State": "abstract_state", "Operational principles": "operational_principles", "Effects table": "effects_table",
	"Statement": "statement", "Rationale": "rationale", "Criterion": "criterion", "Exceptions": "exceptions", "Definition": "definition",
	"Sources": "sources", "Edition": "edition", "Basis": "basis", "Applicability": "applicability", "Supporting content": "supporting_content",
	"Configuration": "configuration", "Dependencies": "dependencies", "Counterpart": "counterpart",
	"When": "when", "Where": "where", "Then": "then", "Semantics": "semantics", "Inputs": "inputs", "Outputs": "outputs",
	"Initiation": "initiation", "Effects": "effects", "Rejections": "rejections", "Guards": "guards", "Caller assumptions": "caller_assumptions",
	"Condition": "condition", "Handling": "handling", "Alternative names": "alternative_names", "Topics": "topics",
	"Relationships": "relationships", "Related terms": "related_terms",
	"Goals": "goals", "Assumptions": "assumptions", "Steps": "steps", "Outcomes": "outcomes", "Actions": "actions",
	"Required interfaces": "required_interfaces", "Provided interfaces": "provided_interfaces", "Quality attributes": "quality_attributes",
	"Glossary": "glossary", "Extensions": "extensions", "Concept uses": "concept_uses", "Synchronizations": "synchronizations",
}

var collectionKinds = map[string]string{
	"Adopts": "AdoptedProvision",
	"Goals":  "Goal", "Assumptions": "Assumption", "Steps": "JobStep", "Outcomes": "DesiredOutcome",
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

func (o *Outline) checkFields(class string, sections []*Section, level int) error {
	schema, builtIn := classSlots()[class]
	// Independently registered types own their authoring contracts.
	if !builtIn {
		return nil
	}
	slots := schema.Slots
	for _, s := range sections {
		field := headingFields[s.Title]
		_, known := slots[field]
		switch s.Title {
		case "Job":
			known = class == "Job"
		case "Specification system":
			_, known = slots["specification_system_content"]
		case "Statement":
			known = class == "Requirement" || class == "BusinessRule"
		case "Boundary":
			known = class == "Component"
		case "System boundary":
			known = class == "UseCase"
		case "Steps":
			known = class == "Job"
		case "Applicability":
			known = class == "Assumption" || class == "AdoptedProvision"
		case "Quality attributes":
			known = class == "QualityModel" || class == "System"
		}
		loc := o.Location(s.Start, s.Title)
		if s.Title == "Obligations" || s.Title == "Requirements" {
			return diagnostic.Error("structure.inline-requirement", "Declare each requirement in a *.requirement.spec.md file and link to it here. Keep scenario outcomes in the use case.", loc)
		}
		if !known {
			return diagnostic.Error("structure.unknown-section", fmt.Sprintf("Unrecognized %s section %q. Move its prose into the field it qualifies or use a declared section.", class, s.Title), loc)
		}
		if s.Level != level {
			return diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-%d heading", s.Title, level), loc)
		}
		if s.Title == "Verification procedures" && len(o.sectionLinks(s)) == 0 {
			return diagnostic.Error("structure.verification-procedure", "Link to at least one verification procedure.", loc)
		}
		if class == "Requirement" && s.Title == "Criterion" {
			for _, sub := range s.Children {
				if sub.Title == "Verification procedures" && (sub.Level != level+1 || len(o.sectionLinks(sub)) == 0) {
					return diagnostic.Error("structure.verification-procedure", "Use a direct criterion subsection with a link to a procedure.", o.Location(sub.Start, sub.Title))
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
			if slices.Contains([]string{"QualityModelAttribute", "GlossaryTerm", "Goal"}, member) {
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
