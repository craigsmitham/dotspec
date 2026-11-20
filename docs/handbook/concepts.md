# Develop concepts

Use this guide when you need to explain a coherent piece of product behavior,
reuse that behavior in another setting, or describe how independent concepts work
together. Start with purpose and a characteristic story; add structure where it
helps people make a design decision. The [format reference](../reference/formats/concept.md#concept-spec-format)
defines the basic file contract.

A Concept is a specification with its own document. ConceptUse and Synchronization
are elements declared inside a System specification, not specifications. They are
traceable, addressable through heading anchors and optional introductory `**ID:**`
values. Use the declared fields for each element and place qualifications next to the behavior they describe. See
[specification types and declared elements](../reference/specification-model.md#specification-types-and-declared-elements).

## Find a purposeful boundary

A concept is a unit of behavior people can understand and use: Shorthand, Trash,
Reservation or Authentication. It can involve several kinds of objects. A domain
noun, class, UI widget, endpoint or delivery story does not by itself define a
concept. For an API, programmers are users too.

Ask what benefit the concept provides on its own and what history of interactions
shows that benefit. Registering credentials alone does not show authentication;
registration followed by a successful or unsuccessful attempt does. If explaining
the benefit requires several independent purposes, consider separate concepts and
a composition rule.

Create a draft with `dotspec concept init shorthand`. Purpose, State, Actions and
Operational principles are conventional sections; unfinished descriptions remain
valid. Use subheadings inside State or Operational principles to organize longer descriptions.

## Write the operational principles

Tell one or a few short stories explaining why the design works. Use actions and
observable state, without prescribing buttons or screen layouts. Include enough
history to reveal the benefit, and any conditions on which the story depends.

Use `## Operational principles` once per document. Named subsections distinguish stories without
requiring scenario records. For Trash, recovery and permanent removal can be two
stories. The stories need not all be successful actor-goal tasks.

An operational principle explains the essential design. State and action behavior
answer the remaining questions. Use cases can describe broader interactions across
concepts; they do not become the concept's operational principle automatically.

## Describe state and actions

Include information the concept must remember to support its actions. Sets and
relations often express that information more clearly than implementation classes.
Do not infer a database layout, UI representation or service boundary from the
abstract state.

Under `## State`, optional `### Initial state` and `### Invariants` describe the
starting state and constraints. Other H3 headings identify components. A component
can have `#### Type` (also Value or Values) and `#### Derivation` subsections.
Derived state need not be stored separately.

Use H3 headings under `## Actions` for individual actions. Inputs, Outputs,
Effects, Rejections, Caller assumptions, Guards and Initiation are optional H4
subsections. Caller assumptions and enforced guards make different promises;
rejection behavior states what happens when a request cannot succeed. Initiation
can explain system behavior such as expiration. None of these headings selects an
execution engine or scheduler.

Describe inputs and outputs in the language of the concept. For example, renaming
a checklist entry requires its identity and new title. State what happens if the
title is blank.

## Separate a definition from a use

A system can declare a named use under **Concept uses**, link its **Definition**,
and describe its **Configuration** and **Dependencies**. Two uses of the same
definition may have different configuration and identities.

## State coordination explicitly

Under **Synchronizations**, give each rule a name. Describe when the rule applies,
what follows, and any timing, retry or ordering guarantees the system promises.
Link to the relevant actions in the prose.

## Inspect and generate views

```sh
dotspec concept check shorthand
dotspec check
dotspec model export --path . > snapshot.json
dotspec model validate snapshot.json
dotspec model docs
```

Authoring checks handle
structure and ordinary links. Model export additionally diagnoses composition
names, reference kinds and declaration scope. It returns status 1 for
composition errors while retaining the snapshot. External or unacquired targets
remain unresolved and produce advice; export does not recursively acquire a
concept package. Select a scope containing the definitions needed for resolution.
A saved snapshot can be validated without reading the working tree.

The generated reference presents purpose and operational principles before state
and actions, then the system’s named uses, dependencies and coordination rules. When saving the output to a file,
use `--relative-to DIR` to resolve links from its destination directory.

In the [To-do spec](../../examples/todo/specification/concepts/checklist.concept.spec-report.md), the
Checklist section keeps the state, actions and effects table together and links to
the system's named use, My checklist. Generate that reading path for an opted-in
project with `dotspec spec build`.

Review purpose, independence, adequate state, end-to-end behavior and composition
promises with people who understand the intended use. Structural validation cannot
establish those qualities, execution guarantees or requirement satisfaction.

This guide adapts Daniel Jackson's [concept criteria](https://essenceofsoftware.com/tutorials/concept-basics/criteria/),
[state](https://essenceofsoftware.com/tutorials/concept-basics/state/),
[operational principles](https://essenceofsoftware.com/tutorials/concept-basics/principle/)
and [dependencies](https://essenceofsoftware.com/tutorials/concept-basics/dependency/),
and Meng and Jackson's [WYSIWID paper](https://arxiv.org/html/2508.14511v2).
The Markdown conventions are dotspec choices. See the
[model rationale](../reference/specification-model.md#concept-definitions-and-composition) for
how these ideas map to the core domain.
