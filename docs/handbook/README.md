# The dotspec handbook

dotspec helps people express what matters about a system in a form they can
read, question, verify and change. Start with the decision you need to make.
For the reasoning behind these choices and where to change them, see
[Maintaining the handbook](maintaining.md).

From an installed CLI, run `dotspec guide` for this index, `dotspec guide use-case`
for a topic, or `dotspec guide examples` to browse the bundled examples. Adopt
[project conventions](project-conventions.md) to point contributors at local sources of authority.

## Choose what to write

Consider a workspace that lets members recover accidentally deleted documents.
The same feature raises several questions. Each guide below helps answer one.

| Question | How the recovery example answers it | Authoring guide |
| --- | --- | --- |
| What is in scope? | The [workspace system](../../examples/requirements-review/specs/system.spec.md) names its boundary, the clock it depends on and the recovery goal. | [Boundaries and interfaces](../review-requirements.md#delimit-the-boundary-and-its-interfaces) |
| Which policy governs the service? | The [business rule](../../examples/requirements-review/specs/business-rules/retention.business-rule.spec.md) sets a 30-day recovery window. Authors can discuss that policy independently of storage design. | [Business rules](../business-rules.md) |
| What progress does a person seek? | The [job](../../examples/requirements-review/specs/jobs/recover.job.spec.md) is to resume work after losing it. Less time spent recovering is a desired outcome. | [Jobs and user needs](../jobs-and-user-needs.md) |
| What can an actor accomplish with this system? | The [use case](../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md) follows a member restoring a document. It also accounts for an expired window or missing authorization. | [Develop use cases](../develop-use-cases.md) |
| What makes the behavior intelligible? | [Trash](../../examples/requirements-review/specs/concepts/trash.concept.spec.md) explains how deletion can be reversible. Its operational principle tells the story, and its actions explain what changes. A glossary entry suffices for a term that only needs a definition. | [Develop concepts](../develop-concepts.md) |
| Where should a promise live? | The use case owns [Restore retained content](../../examples/requirements-review/specs/use-cases/restore.use-case.spec.md#restore-retained-content). A [requirement](../../examples/requirements-review/specs/requirements/restoration.requirement.spec.md) refines which metadata must survive. Link to the governing statement when describing its consequences elsewhere. | [Choose an authoring form](../develop-requirements.md#choose-an-authoring-form) · [Declare relationships](../develop-relationships.md) |
| Which qualities matter, and how good is good enough? | A quality model provides vocabulary for selecting concerns. The [feedback requirement](../../examples/requirements-review/specs/requirements/recovery-feedback.requirement.spec.md) identifies elapsed time as a measure. The bound and applicable operating conditions are still open. | [Quality models](../quality-models.md) · [Define a criterion](../develop-requirements.md#define-or-reference-the-criterion-once) |

## Review what the specification says

Read the promise back to someone who has a stake in it. What would violate it?
What happens at the boundary? For the recovery window, “before the deadline”
leaves a reviewer a concrete question about a request arriving at the deadline.
Review the set as well as the individual statement. A clear recovery rule can
still conflict with a deletion policy. [Review requirements](../review-requirements.md).

Choose the assessment by the question it answers. Comparing restored metadata
with the original can verify a build against the promise. Watching members try
to resume their work can help validate the solution against their need. Record
what was assessed, against what basis, and under which conditions.
[Verification and validation](../verification-and-validation.md).

## Find the corresponding tool or definition

- To create, inspect or check a spec, look up its command and starter. [CLI reference](../cli.md) · [Command examples](../examples.md).
- To understand what a type or relationship means across the domain, start with the explanation and follow it to the model. [Requirements layers](../requirements-layers.md) · [Core domain reference](../reference/model/core.md).
- To determine how particular Markdown becomes a model record, consult the mapping and its source-preservation rules. [Markdown mapping](../../model/markdown-mapping.md) · [Specification model](../specification-model.md).
