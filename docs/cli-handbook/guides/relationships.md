# Declare relationships between specifications

Use explicit relationships when a reader or tool needs to distinguish an actual
support, applicability or boundary claim from a related document. Ordinary
Markdown links remain useful context. They do not automatically make those claims.
The [relationship reference](../reference/model/relationships.md) defines the model's
edge categories, cardinalities, scopes and projections.

## Identify a use case's system boundary

Keep the boundary explanation in the use case's System boundary section. When the
boundary has a system specification, identify it with an explicit Subject system field:

```markdown
## System boundary

**Subject system:** [Delivery service](../delivery-service.system.spec.md)

The actor submits an address through this service; the recipient's mail client
is outside the boundary.
```

The narrative maps to `boundary_description`. The field retains its source reference
in `subject_system_references`; a resolved System populates `subject_system`. A plain link
elsewhere in the narrative is only a citation. The boundary may differ from the
use case's `specification_system`.

## State which user progress a use case supports

Use a Supports section to declare a support relationship:

```markdown
## Supports

[Share information](../jobs/share-information.job.spec.md) and its
[Send a link step](../jobs/share-information.job.spec.md#send-a-link).
```

Every link in this section declares a target: a Goal, Job, JobStep, DesiredOutcome
or higher-level UseCase. A file link identifies the job or use case; a fragment
identifies the explicitly mapped goal, step or outcome. The section is preserved
as `supports_content`, its occurrences as `supports_references`, and resolved
targets as `supports`.

This is an authored design claim, not proof that the use case achieves the job's
outcomes. Put research citations and loosely related jobs in Related
specifications instead. Cross-system support relationships are permitted.

## Identify the subjects of a requirement or business rule

Use an explicit Targets field within Applies to:

```markdown
## Applies to

**Targets:** [Share a link](../use-cases/share-link.use-case.spec.md)

The recipient must have consented to delivery.

### Conditions

Applies when the address is sent to a recipient outside the organization.
```

Requirements and business rules use this form. Targets
may identify any acquired domain record, including a system, concept use, use case
or explicitly named step. Cross-system applicability is permitted.

The complete section remains `applicability: Applicability`. Source occurrences and
resolved domain targets are retained on the requirement or business-rule record as
`applies_to_references` and `applies_to`. Only the Targets field declares targets;
links elsewhere in the section do not silently become Targets. Conditions,
Workload, Population and Exclusions subsections remain authored explanations.
A target link does not establish enforcement, satisfaction or compliance.

## Keep ownership and composition distinct

`owner` means declaration containment: a step belongs to a scenario, an action to
a concept, or a named use to a system declaration. It is checked against the
containing record's member list and source document. It does not describe runtime
object lifetime or business ownership.

`specification_system` associates a non-system specification with a system: the
one selected by a Specification system section, or by default the nearest enclosing
system found through directory discovery. `specification_system_origin` records
which, as `explicit` or `directory`. Systems retain directory context in
`enclosing_system`. Standalone QualityModels have neither association. These
relationships supply no inherited requirements or use-case boundary.

For explicit concept inclusion, dependencies and coordination, use
[Develop concepts](concepts.md). Mutual application dependencies are valid;
containment, explicit `part_of` structure and quality sub-attribute hierarchy must be acyclic.

## Check and consume the result

```sh
dotspec check
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
```

Model export additionally checks relationship target kinds and scope. A broken or
incompatible local reference is an authoring error; an external or unacquired
reference remains unresolved with advice.
The snapshot retains the authored references and diagnostics in either case.

Consumers should read the typed fields as authoritative. The optional `assertions`
collection provides derived views of those fields. Every assertion has both
endpoints, source-backed inputs, and a `projection_path` naming its canonical
field, such as `UseCase.subject_system`, `TraceableElement.supports` or
`Requirement.applies_to`. That path also distinguishes roles, such as
`System.required_interfaces` from `System.provided_interfaces`. Assertions are a view
of the typed graph; their absence does not remove a typed relationship.

## Trace motivation, derivation and assumptions

Use these headings on a specification or on a named traceable element such as a
goal, job step, outcome, concept use, concept action or synchronization:

```markdown
## Supports

[Reduce lost work](../document-workspace.system.spec.md#reduce-lost-work).

## Derived from

[Recovery policy](../business-rules/retention.business-rule.spec.md) supplies the
retention interval.

## Relies on

[Trusted time](../document-workspace.system.spec.md#trusted-time).
```

Use the appropriate deeper heading level inside a named element. Supports accepts
Goal, Job, JobStep, DesiredOutcome and UseCase targets. These are
intended contributions, not achieved outcomes.

Derived from accepts an acquired domain record or a source fragment, including an
explicit decision or a heading in an ordinary referenced Markdown document. A file
link to a specification promotes to its root domain record. For an ordinary
Markdown decision or source, link to the relevant heading. External resources
remain unresolved contextual occurrences with advice; dotspec does not fetch them.
Derived from records a basis, not entailment. Relies on accepts Assumption records.

Each section retains its full prose, including rationale, in `<field>_content`,
reference occurrences in `<field>_references`, and resolved targets in `<field>`.
Every link in one of these sections declares a target. Keep incidental citations
in Notes or Related specifications so they do not declare unintended relationships.
Empty sections are valid drafts. No links are inferred from prose.

## Refine a requirement or business rule

Use Refines on a requirement or business rule:

```markdown
## Refines

[Restore retained content](restore-content.requirement.spec.md).

This requirement identifies the metadata fields while preserving the parent's
recovery window and authorization condition.
```

A requirement can refine Requirement or BusinessRule
records. A business rule refines only other business rules. Multiple parents are
permitted. Explain the added precision in the section or Rationale.
Refinement preserves the parent commitment; it is not supersession, duplication
or proof that the statements are consistent. The tool checks endpoint kinds and
cycles, not logical implication from prose. A cyclic edge is withheld from the
resolved graph, with an error and its authored reference retained in the snapshot.

## Identify participating concept uses

In a use case, name the particular uses participating in the actor-goal interaction:

```markdown
## Involves

[DocumentTrash](../document-workspace.system.spec.md#documenttrash).
```

Each target must be a ConceptUse. Linking a reusable Concept definition is
insufficient when the system might use that definition more than once. Use cases
and concepts remain complementary, many-to-many views with no mandatory authoring
order. Applicability, application dependencies and participation are separate claims.

## Generate traceability views

```sh
dotspec model trace
```

The result includes a layered Mermaid diagram, motivation links, behavioral
participation, direct and transitive change-review candidates, and inspection findings. It derives reverse navigation from canonical
fields and emits only resolved relationships. It does not write files unless
output is redirected. A nonzero authoring exit still produces the review document.

Missing motivation prompts review of the requirement’s justification.
The [To-do restart requirement](../../../examples/todo/specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md)
links to the availability requirement it refines, the storage assumption it relies on
and the desired outcome it supports. See [Requirements layers](../explanations/requirements-layers.md)
for the meaning of those relationships.


## Describe systems, components and features

Use a System for the whole being specified, a Component for a named structural
part, and a Feature for a recognizable grouping of specified product behavior. These are
partial specification subjects. Declare each system and component in its own named
source. For example, `components/checkout/specs/checkout.component.spec.md`
can identify the system that contains it:

```markdown
---
type: component-spec
id: checkout
---
# Checkout

## Part of
[Online Store](../../../specs/online-store.system.spec.md)

## Responsibilities
Calculate totals and coordinate submission.

## Boundary
Checkout owns order submission. Inventory owns stock availability.
```

Group stakeholder-visible behavior in a standalone file such as
`specs/features/place-order.feature.spec.md`. Create it with
`dotspec feature init place-order`.

```markdown
---
type: feature-spec
id: place-order
---
# Place an order

## Description
A customer places an order for available products.

## Purpose
Let a customer buy available products.

## Groups
- [Reserve stock](../requirements/reserve-stock.requirement.spec.md)
```

Each component declares one parent system or component. A nested component links
to its immediate parent. Subsystems use a single system link in their own
`Part of` section. The root system omits that section.

`Groups` can reference requirements, use cases, concepts,
named concept uses, synchronizations and interfaces. Several features can group the same
record. The nearest enclosing system or component owns the feature source. Its overview
lists the feature automatically, using the first paragraph of its Description.

In a requirement, explicitly target a component:

```markdown
## Applies to

**Targets:** [Checkout](../../components/checkout/specs/checkout.component.spec.md)
```

Purpose, Responsibilities and Boundary explain a component. A requirement targeting that component owns the required behavior, rationale and criterion.
Contracts, Architecture and Implementation are ordinary narrative and references
to authoritative native artifacts, not separately modeled architecture definitions.

### Select a specification system

In a business rule, job, use case, concept or requirement:

```markdown
## Specification system

[Online Store](../../store/specs/online-store.system.spec.md#online-store)
```

Use one repository-local link identifying a system document, its title heading,
or its authored System ID. An omitted section uses the nearest directory system;
a broken, external, non-system, empty or ambiguous explicit selection does not
fall back. Relative links obey the existing repository and direct-reference
boundaries; external systems are not fetched. The association governs grouping
and quality classification in show, list, check and model export. It does not set
`subject_system` for a use case or change a requirement's `applies_to` targets.

### Review the overview

```sh
dotspec model overview
```

The command reports explicit part relationships, feature groups, directly
applicable requirements, criteria, concept uses,
descriptions and native references. To read the same subjects in context, open
the [To-do system specification](../../../examples/todo/specification/personal-task-list.system.spec-report.md) and follow the
feature, component or interface that interests you.

The command derives feature involvement from the direct applicability targets
of grouped requirements.
The overview distinguishes declaration location, selected system and directory
context, and retains unresolved authored references in the descriptions and
inspection findings. It makes no completeness, allocation or satisfaction claim.

This scope follows Meyer's distinction between identifying components and their
responsibilities (*Handbook of Requirements and Business Analysis*, §§1.8, 3.6,
8.1) and the separate PEGS treatment of components, functionality, interfaces and
scenarios. It also supports Fowler's emphasis on durable specification and
boundary intent while allowing implementation structure to change:
[Phoenix Primitives](https://chadfowler.com/regenerative-software/3mjfruwwuck2d/),
[Compile to Architecture](https://chadfowler.com/regenerative-software/3mgfsrk75ac2l/).


For enabling abilities that may combine people, procedures and technology, use
the optional [Capability model](map-needs-to-capabilities.md#connect-needs-to-capabilities-and-specifications).
Its dependencies and providers are distinct from feature membership and structural
parentage. The [To-do job section](../../../examples/todo/specification/jobs/resume-interrupted-work.job.spec-report.md) shows the job,
its steps and the outcome supported by the restart requirement.

## Review boundaries and scenario grounding

Use Interacts with on Systems or Components for explicit boundary partners, separate
from Part of. Use an Illustrates bullet inside a scenario step to identify its
underlying requirement, business rule or concept action. See [Review requirements](review.md)
for these forms, exclusions, required/provided interfaces, assumption validation,
open questions and the generated review and PEGS perspectives.
