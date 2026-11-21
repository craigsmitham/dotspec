# The dotspec handbook

dotspec helps people express what matters about a system in a form they can
read, question, verify and change. Start with the decision you need to make.

From an installed CLI, run `dotspec docs` to browse, `dotspec docs job` to choose
a Jobs to be Done page, or `dotspec docs success criteria` to search. Run
`dotspec docs examples` to browse the bundled To-do files. Adopt
[project conventions](guides/project-conventions.md) to point contributors at local sources of authority.

## Start with the To-do app

Add an errand, mark it complete, stop the app and return to it tomorrow. What
would you expect to find? The [runnable walkthrough](../../examples/todo/README.md)
takes that question from user intent to a real process restart.

| Question | Follow it in the app | Authoring guide |
| --- | --- | --- |
| What progress does someone seek? | The [job](../../examples/todo/specs/jobs/resume.job.spec.md) describes the effort of reconstructing unfinished work. | [Write a job from a concrete situation](guides/write-a-job.md) |
| What is in scope? | The [system](../../examples/todo/specs/personal-task-list.system.spec.md) serves one person returning to the same computer. | [Boundaries and interfaces](guides/review.md#delimit-the-boundary-and-its-interfaces) |
| Which policy governs completion? | The [business rule](../../examples/todo/specs/business-rules/completion.business-rule.spec.md) explains why a completed errand is still available for follow-up. | [Business rules](guides/business-rules.md) |
| How does a person act on that policy? | The [use case](../../examples/todo/specs/use-cases/resume-list.use-case.spec.md) follows a person revisiting a finished errand. The [checklist concept](../../examples/todo/specs/concepts/checklist.concept.spec.md) defines how reopening changes its state. | [Use cases](guides/use-cases.md) · [Concepts](guides/concepts.md) |
| Where should a promise live? | The return scenario illustrates [Keep recorded work available](../../examples/todo/specs/requirements/availability.requirement.spec.md). The [restart requirement](../../examples/todo/specs/requirements/restart.requirement.spec.md) refines what must survive and links to its executable verification. | [Choose an authoring form](guides/requirements.md#choose-an-authoring-form) · [Declare relationships](guides/relationships.md) |
| What does durability mean here? | The [quality model](../../examples/todo/specs/quality-models/personal-data.quality-model.spec.md) names the concern. The requirement supplies the measure and the acceptable result. | [Quality models](guides/quality-models.md) · [Define a criterion](guides/requirements.md#define-or-reference-the-criterion-once) |

## Review what the specification says

What would violate the restart requirement? Losing a title, bringing a deleted
task back or quietly giving a reopened task a different identity would each
matter. Run the walkthrough to see how the acceptance test compares what the
person recorded with what the new process serves. [Review requirements](guides/review.md).

To assess whether the workflow helps a person resume work, follow
the [participant exercise](../../examples/todo/README.md#try-the-workflow-with-a-participant).
Record what the participant tried and where they needed help.
[Verification and validation](explanations/verification-and-validation.md).

## Find the corresponding tool or definition

The [reference index](reference/README.md) lists commands, file formats, workspaces and model fields.

- To create, inspect or check a spec, look up its command and starter. [CLI reference](reference/cli.md) · [Command examples](reference/command-examples.md).
- To understand what a type or relationship means across the domain, start with the explanation and follow it to the model. [Requirements layers](explanations/requirements-layers.md) · [Core domain reference](reference/model/core.md).
- To determine how particular Markdown becomes a model record, consult the mapping and its source-preservation rules. [Markdown mapping](../../model/markdown-mapping.md) · [Specification model](reference/specification-model.md).

## Browse by topic

<!-- Generated topic navigation -->

### Specification authority

Find the editable specification for a subject or its named parts. Run `dotspec docs authority`.

- [Choose the authoritative specification](explanations/requirements-layers.md#choose-the-authoritative-specification) (Explanation)

### Business rules

State business constraints and link their enforcement. Run `dotspec docs business-rule`.

- [Business rules](guides/business-rules.md) (How-to)

### Capabilities

Describe an ability, its dependencies and its providers. Run `dotspec docs capability`.

- [Connect needs to capabilities and specifications](guides/map-needs-to-capabilities.md#connect-needs-to-capabilities-and-specifications) (How-to)

### Components

Name a component, declare its parent and organize its local specifications. Run `dotspec docs component`.

- [Component spec format](reference/formats/system.md#component-spec-format) (Reference)

### Concepts

Explain purposeful behavior through state and actions. Run `dotspec docs concept`.

- [Develop concepts](guides/concepts.md) (How-to)

### Project conventions

Identify a project’s local specification authorities. Run `dotspec docs conventions`.

- [Project conventions](guides/project-conventions.md) (How-to)

### Core domain

Inspect the core domain types, relationships and diagrams. Run `dotspec docs core-domain`.

- [Core domain model](reference/model/core.md) (Reference)

### Worked examples

Browse the bundled To-do project and its specifications. Run `dotspec docs examples`.

- [To-do](../../examples/todo/README.md) (Example)

### Features

Group stakeholder-visible behavior across specifications. Run `dotspec docs feature`.

- [Describe systems, components and features](guides/relationships.md#describe-systems-components-and-features) (How-to)

### Authoring handbook

Choose what to write using the runnable To-do example. Run `dotspec docs handbook`.

- [Start with the To-do app](README.md#start-with-the-to-do-app) (Handbook)

### Inspection

Interpret check results, source discovery and references. Run `dotspec docs inspection`.

- [Checking and discovering specifications](reference/inspection.md) (Reference)

### Jobs to be done

Describe the progress someone seeks and the success criteria they use to judge it. Run `dotspec docs job`.

- [Write a job from a concrete situation](guides/write-a-job.md) (How-to)
- [What makes something a job to be done?](explanations/jobs-to-be-done.md) (Explanation)
- [Job spec format](reference/formats/job.md#job-spec-format) (Reference)
- [Resume interrupted work](../../examples/todo/specs/jobs/resume.job.spec.md) (Example)

### Project layout

Find a specification by source path or owning directory. Run `dotspec docs layout`.

- [Specification layout](reference/layout.md) (Reference)

### Markdown mapping

Look up how authored Markdown maps to model elements. Run `dotspec docs mapping`.

- [Markdown mapping profile](../../model/markdown-mapping.md) (Reference)

### Quality criteria

Define a criterion and distinguish a measure from its required bound. Run `dotspec docs quality`.

- [Define or reference the criterion once](guides/requirements.md#define-or-reference-the-criterion-once) (How-to)

### Quality models

Select quality attributes and apply them to a system. Run `dotspec docs quality-model`.

- [Populate system quality attributes from a model](guides/quality-models.md) (How-to)

### Relationships

Connect specifications through explicit relationships. Run `dotspec docs relationships`.

- [Declare relationships between specifications](guides/relationships.md) (How-to)

### Requirements

Author assessable requirements and satisfaction criteria. Run `dotspec docs requirement`.

- [Develop requirements](guides/requirements.md) (How-to)

### Requirements review

Review a requirement set and surface unresolved decisions. Run `dotspec docs review`.

- [Review requirements and unresolved questions](guides/review.md) (How-to)

### Systems

Describe a system, its boundary, goals and quality attributes. Run `dotspec docs system`.

- [Systems](reference/formats/system.md) (Reference)

### Use cases

Describe actor goals, boundaries and alternative paths. Run `dotspec docs use-case`.

- [Develop use cases](guides/use-cases.md) (How-to)

### Verification and validation

Distinguish the four review questions and the evidence each needs. Run `dotspec docs verification`.

- [Verification and validation](explanations/verification-and-validation.md) (Explanation)
