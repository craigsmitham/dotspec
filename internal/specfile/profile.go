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
	"Description": "description", "Part of": "part_of", "Purpose": "purpose", "Scope": "scope", "Out of scope": "out_of_scope",
	"Project context": "project_context", "Environment": "environment", "Stakeholders": "stakeholders", "Constraints": "constraints",
	"Responsibilities": "responsibilities", "Boundary": "boundary_description", "Contracts": "contracts", "Architecture": "architecture", "Implementation": "implementation",
	"Interacts with": "interacts_with", "Specification system": "specification_system", "Enablers": "enablers",
	"Groups": "groups", "Depends on": "depends_on", "Provided by": "provided_by", "Supports": "supports", "Involves": "involves",
	"Derived from": "derived_from", "Relies on": "relies_on", "Refines": "refines", "Applies to": "applies_to",
	"Evidence": "evidence", "Authority": "authority_and_evidence", "Related specifications": "related_specifications", "Open questions": "open_questions",
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
	"Glossary": "glossary", "Obligations": "obligations", "Extensions": "extensions", "Concept uses": "concept_uses", "Synchronizations": "synchronizations",
}

var collectionKinds = map[string]string{
	"Goals": "Goal", "Assumptions": "Assumption", "Steps": "JobStep", "Outcomes": "DesiredOutcome",
	"Actions": "ConceptAction", "Required interfaces": "Interface", "Provided interfaces": "Interface",
	"Quality attributes": "QualityAttribute", "Glossary": "GlossaryTerm", "Obligations": "ObligationClause",
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

func specificationClass(kind string) string {
	if kind == "quality-attribute" {
		return "SystemQualityAttribute"
	}
	var class string
	for part := range strings.SplitSeq(kind, "-") {
		if part != "" {
			class += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return class
}

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
		case "Statement":
			known = class == "Requirement" || class == "BusinessRule"
		case "Boundary":
			known = class == "Component"
		case "System boundary":
			known = class == "UseCase"
		case "Steps":
			known = class == "Job"
		case "Applicability":
			known = class == "Assumption"
		case "Quality attributes":
			known = class == "QualityModel" || class == "System"
		}
		loc := o.Location(s.Start, s.Title)
		if !known && s.Title == "Obligations" {
			return diagnostic.Error("obligation.container", "This element cannot declare obligation clauses.", loc)
		}
		if !known {
			return diagnostic.Error("structure.unknown-section", fmt.Sprintf("Unrecognized %s section %q. Move its prose into the field it qualifies or use a declared section.", class, s.Title), loc)
		}
		if s.Level != level {
			return diagnostic.Error("structure.section-level", fmt.Sprintf("%s must use a level-%d heading", s.Title, level), loc)
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
			// Definitions and goal statements own their subordinate Markdown.
			if slices.Contains([]string{"QualityAttribute", "GlossaryTerm", "Goal"}, member) {
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
