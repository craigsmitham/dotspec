# Requirements layers and relationships

Dotspec represents requirements as a layered traceability graph. This explanation
is for authors and model maintainers deciding where an element belongs and what
its relationships mean. The [canonical model](../model/dotspec.yaml) defines the
types; [Declare relationships](develop-relationships.md) gives the Markdown forms.

## Separate concern, kind and classification

| Axis | Meaning | Examples |
| --- | --- | --- |
| Concern subset | Whose concerns and which level of abstraction the model serves | business-requirements, user-requirements, software-requirements |
| Subset classification | The subset's role in dotspec's domain | core, supporting, generic |
| Artifact kind | What a record represents | SystemGoal, BusinessRule, Job, UseCase, Concept, Requirement |
| Relationship | Why or how two particular elements are connected | supports, derived_from, refines, relies_on, applies_to, involves |

Subset membership may overlap. It does not make every member an obligation,
assign priority, determine acceptance, or impose a creation order. A reusable
concept is not thereby a member of a generic DDD subdomain. Pace of change is
another concern: a business policy can change faster than a software protocol.
No fixed pace or authority is inferred from a type or directory.

| Concern | Main question | Typical content |
| --- | --- | --- |
| Business | Why should the system exist, and which organizational rules govern it? | Goals, business rules, organizational constraints |
| User | What must people accomplish, and which outcomes matter? | Jobs, steps, outcomes, stakeholder needs |
| Software | What behavior and qualities must the software provide? | Use cases, concepts and composition, behavioral obligations, quality requirements |

There is no requirement to populate every layer on every path. A software
requirement can derive directly from a business rule. A design choice can
introduce a new obligation. Environment assumptions and project constraints can
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

## Distinguish an obligation from its document

`Requirement` is the abstract family of standalone GeneralRequirement and
QualityRequirement specifications. It is not the superclass of every normative
expression. Accepted use-case prose or concept contracts can also express
obligations. Acceptance comes from the consuming project's authority policy;
neither a type tag nor a generated projection establishes it.

Use an optional `ObligationClause` when a particular rule within an existing
specification needs its own reference, refinement, applicability or verification
guidance. The clause retains the original statement and source identity. Other
views link to that expression rather than create a second authority. Prose need
not be split into clauses merely to satisfy the model, and scenarios are not
assumed to enumerate all required behavior.

## Give each relationship one meaning

| Relationship | Meaning | Does not establish |
| --- | --- | --- |
| supports | Intended contribution to a goal, job, step or outcome | Achievement or stakeholder validation |
| derived_from | Source or rationale for the element | Logical entailment or refinement |
| refines | More precise obligation preserving a parent commitment | Replacement, supersession or a proof of consistency |
| relies_on | Validity conditioned on a named assumption | That the assumption is true |
| applies_to | Subject governed by an obligation or business rule | Implementation satisfaction |
| involves | Named concept uses participating in a use case | Calls, service boundaries or required execution order |

Refinement may have several parents and must be acyclic. Children may collectively
cover a broader commitment; the presence of one refining child does not establish
that coverage. Alternative designs require an explicit decision, not a claim
that every alternative is simultaneously required. Support and application
relationships may cycle; the entire graph is not a DAG.

`ConceptUse.dependencies` remains an application-composition relationship. It is
not a substitute for requirement derivation. Containment, directory context,
ordinary citations and semantic associations also retain distinct meanings.

## Generate views from the same relationships

`dotspec model trace` produces an instance-level diagram and four reader views:

- Motivation and basis: follow explicit justification and applicability.
- Behavior: inspect use cases, participating concept uses and related obligations.
- Change review: inspect direct and transitive reverse links with their meanings.
- Verification guidance: find authored check plans and missing guidance.

These views are bounded by acquired sources and explicit links. Transitive
reachability identifies review candidates; it does not prove actual impact.
Missing justification or verification guidance is advisory so unfinished drafts
remain usable. Broken or incompatible local references and refinement cycles are
errors. Unacquired and external references retain their occurrences with advice.

Implementation realization, verification definitions and revision-bound execution
evidence need separate relationships when introduced. Current verification and
Evidence prose must not be promoted into satisfaction claims or execution results.

The [recovery example](reference/requirements-traceability-example.md) demonstrates
the relationships without requiring a duplicate requirement for every scenario step.

## Method basis

This is dotspec's synthesis, not a schema prescribed by any one source.

- Chad Fowler's [Pace Layers and AI Integration](https://chadfowler.com/regenerative-software/3maob46kbz22v/)
  distinguishes rates of change. [The Specification Is Not a Document](https://chadfowler.com/regenerative-software/3mtgs36dnq22o/)
  supports multiple queryable views of identifiable, related specification knowledge.
- Bertrand Meyer's [Handbook of Requirements and Business Analysis](https://doi.org/10.1007/978-3-031-06739-6),
  §§1.1–1.2 distinguishes PEGS dimensions and properties from statements; chapter 5
  distinguishes reference/refinement from duplication; §§7.5–7.6 discusses scenarios
  and general requirements. Allowing accepted normative scenarios remains dotspec's
  choice, not Meyer's endorsement of scenarios as complete specifications.
- [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html), §§5.3–5.4,
  addresses business, stakeholder, system and software information items, iterative
  refinement and design-derived requirements. §6.3.3.4 discusses using scenarios and
  use cases to identify requirements. Its operational concept is not synonymous
  with Jackson's reusable concept.
- Daniel Jackson's [Operational principles](https://essenceofsoftware.com/tutorials/concept-basics/principle/)
  distinguishes purpose-explaining stories, use cases and fuller state/action
  specifications. The sibling view and `involves` relationship are dotspec's synthesis.
