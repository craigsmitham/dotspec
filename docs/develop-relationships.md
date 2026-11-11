# Declare relationships between specifications

Use explicit relationships when a reader or tool needs to distinguish an actual
support, applicability or boundary claim from a related document. Ordinary
Markdown links remain useful context. They do not automatically make those claims.
The [relationship reference](reference/model/relationships.md) defines the model's
edge categories, cardinalities, scopes and projections.

## Identify a use case's system boundary

Keep the boundary explanation in the use case's System boundary section. When the
boundary has a system specification, identify it with an explicit System field:

```markdown
## System boundary

**System:** [Delivery service](../system.spec.md)

The actor submits an address through this service; the recipient's mail client
is outside the boundary.
```

The narrative maps to `system_boundary`. The field retains its source reference
in `boundary_references`; a resolved System populates `boundary`. A plain link
elsewhere in the narrative is only a citation. The boundary may differ from the
use case's directory-derived `system_context`.

## State which user progress a use case supports

Use a Supports jobs section to declare a support relationship:

```markdown
## Supports jobs

[Share information](../jobs/share-information.job.spec.md) and its
[Send a link step](../jobs/share-information.job.spec.md#send-a-link).
```

Every link in this section declares a target: a Job, JobStep or Outcome. A file
link identifies the job; a fragment identifies the explicitly mapped step or
outcome. The section is preserved as `supported_jobs_content`, its occurrences as
`supported_jobs_references`, and resolved targets as `supported_jobs`.

This is an authored design claim, not proof that the use case achieves the job's
outcomes. Put research citations and loosely related jobs in Related jobs instead.
Cross-system support relationships are permitted.

## Identify the subjects of a requirement or business rule

Use an explicit Targets field within Applies to:

```markdown
## Applies to

**Targets:** [Share a link](../use-cases/share-link.use-case.spec.md)

The recipient must have consented to delivery.

### Conditions

Applies when the address is sent to a recipient outside the organization.
```

Both general and quality requirements, and business rules, use this form. Targets
may identify any acquired domain record, including a system, concept use, use case
or explicitly named step. Cross-system applicability is permitted.

The complete section remains `applies_to: Applicability`. Source occurrences and
resolved domain targets are retained on the root record as
`applicability_targets_references` and `applicability_targets`. A Subject or
Subjects subsection remains contextual; its links do not silently become Targets.
Conditions, exclusions, workloads and populations remain authored explanations.
A target link does not establish enforcement, satisfaction or compliance.

## Keep ownership and composition distinct

`owner` means declaration containment: a step belongs to a scenario, an action to
a concept, or a named use to a system declaration. It is checked against the
containing record's member list and source document. It does not describe runtime
object lifetime or business ownership.

`system_context` associates a non-system specification with the enclosing system
found through directory discovery. Systems use `parent_system` for their hierarchy.
Standalone QualityModels have neither association. Neither relationship supplies
inherited requirements or establishes a use case boundary.

For explicit concept inclusion, type bindings, dependencies and coordination, use
[Develop concepts](develop-concepts.md). Mutual application dependencies are valid;
containment, system hierarchy and quality sub-attribute hierarchy must be acyclic.

## Check and consume the result

```sh
dotspec check
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
```

Model export additionally checks relationship target kinds and scope. A broken or
incompatible local reference is an authoring error; an external or unacquired
reference remains unresolved with advice. If a concept definition or supplying
use is unknown, dependent scope checks are deferred. Unknown is not incompatible.
The snapshot retains the authored references and diagnostics in either case.

Consumers should read the typed fields as authoritative. The optional `assertions`
collection provides derived views such as BoundaryIs, SupportsJob and AppliesTo.
Every assertion has both endpoints, source-backed inputs, and a `projection_path`
identifying its canonical field. That path also distinguishes concept parameters
from introduced types when both project as HasConceptType. Assertions are a partial
view of the full graph; their absence does not remove a typed relationship.

## Trace motivation, derivation and assumptions

Use these headings on a specification or on a named traceable element such as a
system goal, job step, outcome, concept use, concept action or synchronization:

```markdown
## Supports

[Reduce lost work](../system.spec.md#reduce-lost-work).

## Derived from

[Recovery policy](../business-rules/retention.business-rule.spec.md) supplies the
retention interval.

## Relies on

[Trusted time](../system.spec.md#trusted-time).
```

Use the appropriate deeper heading level inside a named element. Supports accepts
SystemGoal, Job, JobStep and Outcome targets. The existing Supports jobs form stays
available on use cases with its narrower Job/JobStep/Outcome range. These are
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

## Identify an obligation without creating another document

Use an Obligations section within a specification, concept action or other
traceable element when a rule needs its own identity:

```markdown
## Obligations

### Restore retained content

**ID:** restore-content

The workspace must restore retained document content and metadata when an
 authorized member requests restoration within the recovery window.

#### Derived from

[Recovery policy](../business-rules/retention.business-rule.spec.md).

#### Applies to

**Targets:** [DocumentTrash](../system.spec.md#documenttrash)

#### Verification

Compare content and metadata before deletion and after restoration.
```

The introductory prose becomes the clause's statement. Applies to, Rationale and
Verification use the same forms as standalone general requirements. Relationship
headings can appear beneath the clause. An optional ID must be nonempty and unique
among declared IDs in the document; otherwise link to the generated heading anchor.
IDs and anchors work in same-file and cross-file links. Generated reader links use
Markdown heading anchors so they also work in ordinary Markdown viewers.

This creates an owned `ObligationClause`, not a new file. Do not copy the statement
into a standalone requirement just to link it. If a narrative already owns the
obligation, refer to that clause from scenarios or other views. The declaration
does not itself record acceptance; retain the consuming project's authority policy.

## Refine an obligation

Use Refines on a standalone requirement, business rule or obligation clause:

```markdown
## Refines

[Restore retained content](../use-cases/restore.use-case.spec.md#restore-content).

This requirement identifies the metadata fields while preserving the parent's
recovery window and authorization condition.
```

Targets must be Requirement, BusinessRule or ObligationClause records. Multiple
parents are permitted. Explain the added precision in the section or Rationale.
Refinement preserves the parent commitment; it is not supersession, duplication
or proof that the statements are consistent. The tool checks endpoint kinds and
cycles, not logical implication from prose. A cyclic edge is withheld from the
resolved graph, with an error and its authored reference retained in the snapshot.

## Identify participating concept uses

In a use case, name the particular uses participating in the actor-goal interaction:

```markdown
## Involves

[DocumentTrash](../system.spec.md#documenttrash).
```

Each target must be a ConceptUse. Linking a reusable Concept definition is
insufficient when the system might use that definition more than once. Use cases
and concepts remain complementary, many-to-many views with no mandatory authoring
order. Applicability, application dependencies and participation are separate claims.

## Generate traceability views

```sh
dotspec model trace --path . --relative-to docs/reference > docs/reference/traceability.md
```

The result includes a layered Mermaid diagram, motivation links, behavioral
participation, direct and transitive change-review candidates, verification
guidance, and inspection findings. It derives reverse navigation from canonical
fields and emits only resolved relationships. It does not write files unless
output is redirected. A nonzero authoring exit still produces the review document.

Missing motivation or verification guidance is advisory. A verification plan is
not an execution result, and a path through the graph does not prove satisfaction.
The [generated recovery example](reference/requirements-traceability-example.md)
shows the complete flow; [Requirements layers](requirements-layers.md) explains
why these distinctions matter.
