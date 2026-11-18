# The dotspec handbook

dotspec helps people express what matters about a system in a form they can
read, question, verify and change. Start with the decision you need to make.
For the reasoning behind these choices and where to change them, see
[Maintaining the handbook](maintaining.md).

From an installed CLI, run `dotspec guide` for this index, `dotspec guide use-case`
for a topic, or `dotspec guide examples` to browse the bundled examples. Adopt
[project conventions](project-conventions.md) to point contributors at local sources of authority.

## Start with the To-do app

Add an errand, mark it complete, stop the app and return to it tomorrow. What
would you expect to find? The [runnable walkthrough](../../examples/todo/README.md)
takes that question from user intent to a real process restart.

| Question | Follow it in the app | Authoring guide |
| --- | --- | --- |
| What progress does someone seek? | The [job](../../examples/todo/specs/jobs/resume.job.spec.md) describes the effort of reconstructing unfinished work. | [Jobs and user needs](../jobs-and-user-needs.md) |
| What is in scope? | The [system](../../examples/todo/specs/system.spec.md) serves one person returning to the same computer. | [Boundaries and interfaces](../review-requirements.md#delimit-the-boundary-and-its-interfaces) |
| Which policy governs completion? | The [business rule](../../examples/todo/specs/business-rules/completion.business-rule.spec.md) explains why a completed errand is still available for follow-up. | [Business rules](../business-rules.md) |
| How does a person act on that policy? | The [use case](../../examples/todo/specs/use-cases/resume-list.use-case.spec.md) follows a person revisiting a finished errand. The [checklist concept](../../examples/todo/specs/concepts/checklist.concept.spec.md) defines how reopening changes its state. | [Use cases](../develop-use-cases.md) · [Concepts](../develop-concepts.md) |
| Where should a promise live? | The use case owns [Keep recorded work available](../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available). The [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md) refines what must survive and links to its executable verification. | [Choose an authoring form](../develop-requirements.md#choose-an-authoring-form) · [Declare relationships](../develop-relationships.md) |
| What does durability mean here? | The [quality model](../../examples/todo/specs/quality-models/personal-data.quality-model.spec.md) names the concern. The requirement supplies the measure and the acceptable result. | [Quality models](../quality-models.md) · [Define a criterion](../develop-requirements.md#define-or-reference-the-criterion-once) |

## Review what the specification says

What would violate the restart requirement? Losing a title, bringing a deleted
task back or quietly giving a reopened task a different identity would each
matter. Run the walkthrough to see how the acceptance test compares what the
person recorded with what the new process serves. [Review requirements](../review-requirements.md).

To assess whether the workflow helps a person resume work, try the exercise in
the [system validation plan](../../examples/todo/specs/system.spec.md#validation-plan).
Record what the participant tried and where they needed help.
[Verification and validation](../verification-and-validation.md).

## Find the corresponding tool or definition

- To create, inspect or check a spec, look up its command and starter. [CLI reference](../cli.md) · [Command examples](../examples.md).
- To understand what a type or relationship means across the domain, start with the explanation and follow it to the model. [Requirements layers](../requirements-layers.md) · [Core domain reference](../reference/model/core.md).
- To determine how particular Markdown becomes a model record, consult the mapping and its source-preservation rules. [Markdown mapping](../../model/markdown-mapping.md) · [Specification model](../specification-model.md).
