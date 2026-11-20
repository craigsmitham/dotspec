# Requirements layers and relationships

dotspec represents requirements as a layered traceability graph. This explanation
is for authors and model maintainers deciding where an element belongs and what
its relationships mean. The [canonical model](../../model/metamodel.yaml) defines the
types; [Declare relationships](relationships.md) gives the Markdown forms.

## Separate concern, kind and classification

| Axis | Meaning | Examples |
| --- | --- | --- |
| Concern subset | Whose concerns and which level of abstraction the model serves | business-requirements, user-requirements, software-requirements |
| Subset classification | The subset's role in dotspec's domain | core, supporting, generic |
| Record kind | What a record represents | Goal, BusinessRule, Job, UseCase, Concept, Requirement |
| Relationship | Why or how two particular elements are connected | supports, derived_from, refines, relies_on, applies_to, involves |

Subset membership may overlap. It does not make every member a requirement,
assign priority, determine acceptance, or impose a creation order. A reusable
concept is not thereby a member of a generic DDD subdomain. Pace of change is
another concern: a business policy can change faster than a software protocol.
No fixed pace or authority is inferred from a type or directory.

| Concern | Main question | Typical content |
| --- | --- | --- |
| Business | Why should the system exist, and which organizational rules govern it? | Goals, business rules, organizational constraints |
| User | What must people accomplish, and which outcomes matter? | Jobs, steps, outcomes, stakeholder needs |
| Software | What behavior and qualities must the software provide? | Use cases, concepts and composition, behavioral requirements, quality requirements |

There is no requirement to populate every layer on every path. A software
requirement can derive directly from a business rule. A design choice can
introduce a new requirement. Environment assumptions and project constraints can
matter at several levels. Revisiting an earlier decision is compatible with an
acyclic refinement relation: iteration in the authoring process is not a semantic
refinement cycle in the current model.

## Concepts and use cases are complementary

A use case describes an actor pursuing a goal across a system boundary. A Jackson
concept describes purposeful behavior through state, actions and operational
principles. A concept use includes that definition under a name in a particular
system; synchronizations coordinate uses.

Use cases and concepts have a many-to-many relationship. `involves` points from a
use case to the particular concept uses participating in it. A use case may reveal
a missing concept, and a reusable concept may suggest useful interactions. A
concept overview can precede use cases in a reader view without becoming a
mandatory dependency or authoring order.

An operational principle explains a concept's characteristic purpose. It does
not replace a complete use case or state/action specification. A business domain
term in a glossary also need not become a concept with actions and state.

## Choose the authoritative specification

Each specification owns one subject. Edit an interaction and its guarantees in the use case, an action and its effects in the concept, and a policy in the business rule. Put each product requirement in a `*.requirement.spec.md` file. Its conditions, criterion, rationale and quality classification travel with the statement. A use case can exercise several requirements and several use cases can exercise the same requirement.

| Specification kind | Authoritative content |
| --- | --- |
| System | Purpose, goals, scope, environment, assumptions and system boundary |
| Component | The component’s responsibility, boundary, interfaces and declared place in the system |
| Job | The progress sought, steps, and desired outcomes |
| Use case | An actor’s goal and interaction, including preconditions, success path, extensions and guarantees |
| Concept | The reusable behavior, including state, actions, effects and invariants |
| Feature | A recognizable unit of product behavior and the subjects explicitly grouped into it |
| Capability | What a provider enables and its declared providers |
| Business rule | The governing business policy, its applicability, exceptions and basis |
| Requirement | One required property or behavior, its conditions, criterion, rationale and declared relationships |
| Quality model | A reusable vocabulary of quality attributes and their definitions |
| System quality attribute | The meaning of a quality concern for this system, including its relationship to a reusable definition |

A named part such as a scenario step or a concept action belongs to its declaring file. Link to the part when the relationship concerns that precise passage. Generated reports provide other reading paths to the same source.

For an older `## Obligations` declaration, move a general requirement and its relationships to its own file, then update incoming links. If the text describes a scenario-specific result, keep it in the success guarantees or extension handling. [Requirement format and migration](../reference/formats/requirement.md#requirement-format)

## Give each relationship one meaning

| Relationship | Meaning | Does not establish |
| --- | --- | --- |
| supports | Intended contribution to a goal, job, step, outcome or higher-level use case | Achievement or stakeholder validation |
| derived_from | Source or rationale for the element | Logical entailment or refinement |
| refines | More precise requirement or business rule preserving a parent commitment | Replacement, supersession or a proof of consistency |
| relies_on | Validity conditioned on a named assumption | That the assumption is true |
| applies_to | Subject governed by a requirement or business rule | Implementation satisfaction |
| illustrates | A scenario step demonstrates a requirement, business rule or concept action | A new requirement or proof of scenario completeness |
| involves | Named concept uses participating in a use case | Calls, service boundaries or required execution order |

A requirement or business rule declares `applies_to` with the Targets field of its
Applies to section. The narrative operating context in that section is its
`applicability`, which creates no relationship.

Refinement may have several parents and must be acyclic. A requirement can refine
a requirement or a business rule; a business rule refines
only other business rules. Children may collectively cover a broader commitment;
the presence of one refining child does not establish that coverage. Alternative designs require an explicit decision, not a claim
that every alternative is simultaneously required. Support and application
relationships may cycle; the entire graph is not a DAG.

`ConceptUse.dependencies` remains an application-composition relationship. It is
not a substitute for requirement derivation. Containment, directory context,
ordinary citations and semantic associations also retain distinct meanings.

## Generate views from the same relationships

`dotspec model trace` produces an instance-level diagram and three reader views:

- Motivation and basis: follow explicit justification and applicability.
- Behavior: inspect use cases, participating concept uses and related requirements.
- Change review: inspect direct and transitive reverse links with their meanings.

These views are bounded by acquired sources and explicit links. Transitive
reachability identifies review candidates; it does not prove actual impact.
Missing justification is advisory so unfinished drafts
remain usable. Broken or incompatible local references and refinement cycles are
errors. Unacquired and external references retain their occurrences with advice.

Link to tests and evidence at their source. Following a relationship identifies a
commitment affected by a change. Assessing whether the product meets it requires evidence.

The [To-do example](../../examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md) demonstrates
the relationships without requiring a duplicate requirement for every scenario step.

## Method basis

This is dotspec's synthesis, not a schema prescribed by any one source.

- Chad Fowler's [Pace Layers and AI Integration](https://chadfowler.com/regenerative-software/3maob46kbz22v/)
  distinguishes rates of change. [The Specification Is Not a Document](https://chadfowler.com/regenerative-software/3mtgs36dnq22o/)
  supports multiple queryable views of identifiable, related specification knowledge.
- Bertrand Meyer's [Handbook of Requirements and Business Analysis](https://doi.org/10.1007/978-3-031-06739-6),
  §§1.1–1.2 distinguishes PEGS dimensions and properties from statements; chapter 5
  distinguishes reference/refinement from duplication; §§7.5–7.6 discusses scenarios
  and general requirements. Section 11.3 uses identified scenarios to check whether the specified functionality is sufficient.
- [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html), §§5.3–5.4,
  addresses business, stakeholder, system and software information items, iterative
  refinement and design-derived requirements. §6.3.3.4 discusses using scenarios and
  use cases to identify requirements. Its operational concept is not synonymous
  with Jackson's reusable concept.
- Daniel Jackson's [Operational principles](https://essenceofsoftware.com/tutorials/concept-basics/principle/)
  distinguishes purpose-explaining stories, use cases and fuller state/action
  specifications. The sibling view and `involves` relationship are dotspec's synthesis.


## Verification and validation across layers

Requirements verification reviews the quality of statements and sets; requirements
validation reviews their faithfulness to needs and intended use. System verification
assesses a realization against requirements; system validation assesses its intended-use
fitness. These distinctions apply across all concern layers, and methods can support
more than one purpose. See [Verification and validation](verification-and-validation.md)
for the governing vocabulary and [Review requirements](review.md) for review prompts.
