# Develop use cases

Use this guide when you need to agree on what an actor can accomplish with a
system, investigate alternatives, or retain behavior across delivery increments.
Start with a readable sketch and expand the portions needed for the next decision.
The [format reference](../README.md#use-case-spec-format) defines the supported
sections and checks; this guide explains how to use them.

## Start with actors and goals

List the system's external participants, including other software, and the goals
for which each primary actor seeks its service. A supporting actor provides a
service or receives a notification. Roles depend on the boundary and interaction;
being important does not make a stakeholder the primary actor.

Create a sketch for an important goal:

```sh
dotspec use-case init reserve-equipment
```

Replace the title with an active verb phrase. Name the primary actor, boundary,
and goal level. The boundary can be an application or a business process; an
enclosing system file organizes the document but does not determine that boundary.
Keep the actor's goal distinct from the broader job or desired outcome it supports.

Use summary (kite) for a larger accomplishment, user task (sea) for a coherent
user goal, and subfunction (fish) for a complex or reused supporting interaction.
A button click rarely warrants its own use case. Link related use cases in prose
rather than forcing them into a strict hierarchy or assuming a code-module boundary.

Inspect the collection with:

```sh
dotspec use-case list --overview
```

The view displays existing prose, including missing fields. It cannot establish
that every actor or goal has been found. During an early scope discussion, a
simple actor-goal list in an ordinary Markdown document is also useful; each entry
does not yet need its own specification file.

## Sketch the interaction, then investigate alternatives

Connect important user tasks with one or a few summary narratives. The repository's
[Prepare specifications for review](../specs/use-cases/prepare-specifications.use-case.spec.md)
links creating a draft, discovering sources, and checking specifications.

For a selected user task, write a short typical success path. State actor intentions
and visible system responsibilities, including exchanges with external collaborators.
Use a paragraph or numbered steps. Keep data schemas, interface details, and internal
design in their appropriate sources. Explain whether sequence is representative;
link any rule that actually requires an order.

Brainstorm alternative success paths and failures before solving them all. A
condition-only extension with “handling unresolved” is useful evidence of a question.
Prioritize investigation by the value and risk of the upcoming decision. A filled
section is not evidence that people have agreed or that the behavior is adequate.

For each developed extension, explain:

- Where the condition arises in the interaction.
- What the actor, system, or collaborator does in response.
- Whether the interaction resumes, succeeds another way, or ends without the goal.

Use nested prose when a response itself can fail. If branches become hard to read,
link a decision table, state model, or substantial subfunction. Avoid enumerating
all interleavings merely to make the narrative appear complete.

For a nested branch that model export should distinguish, use a labelled heading
or list item, or group descriptive headings under an `Extensions` subsection:

```markdown
### 1a. Equipment unavailable

Offer an alternative.

#### 1a1a. Alternative declined

End the request without a reservation.

#### Notes

The equipment catalogue remains the source of availability information.
```

Here the labelled child is a branch and Notes is explanatory content. Labels such
as `1a2a` and relative `.1a` are optional reading aids. Markdown nesting determines
branch ownership. For casual prose, simply indent and keep writing; it remains
fully represented even when it is not decomposed into branch records.

Use `Condition`, `Location` and `Handling` subsections only when that separation
helps the reader. Handling can be a paragraph or steps, including where the
interaction resumes or ends. Write unresolved questions in prose. There is no
required handling status or actor/action/result form for each step.

## Add context where it changes the promise

For a new document needing more prompts, use:

```sh
dotspec use-case init reserve-equipment-detailed --template expanded
```

Both starters use the same format. For an existing document, add the relevant
headings directly; `init` never overwrites it. Remove unhelpful prompts rather
than maintaining empty ceremony.

Distinguish the initiating event from preconditions already expected true. If the
system must handle a condition being false, describe that behavior instead of
assuming it away. State the success criterion and, when consequential, minimum
guarantees that survive failure. Supporting actors and stakeholder interests can
expose promises that the primary actor's success alone misses.

[Check specifications](../specs/use-cases/check-specifications.use-case.spec.md)
illustrates this distinction: errors found in specifications can accompany a
successful complete report; unreadable input can prevent that goal. Source
preservation applies in either outcome.

## Link branches and decide behavioral authority

Use a descriptive heading when a branch needs references. For example:

```markdown
## Extensions

### Payment result unknown

During authorization, the provider's response is lost.
Handling unresolved: determine how to reconcile before confirming a reservation.
```

A work item, requirement, or test description can link to
`reserve-equipment.use-case.spec.md#payment-result-unknown`. Within the same file,
use `#payment-result-unknown`. Keep the heading stable while references rely on it;
if it changes, update links and run `dotspec check`. Labels such as `3a` can help
readers locate a branch, but they are not durable identities. Dotspec does not
parse step numbering into a control-flow graph. You do not need separate files
or authored IDs for branches; model export creates source-backed branch records.

Choose where each obligation is defined. Accepted use-case prose can carry a
behavioral obligation directly; it need not be duplicated in a requirement file.
Alternatively, a narrative can explain and link to a rule owned by another source.
Explain that role around the link. Keep general rules and complementary data or
quality specifications where they can state obligations more clearly.

Tests may define accepted examples or verify an existing rule. State which role
they play. Link resolution establishes that a target exists; it does not establish
execution, branch coverage, completeness, stakeholder agreement, or customer value.
The project's [source-authority principle](../README.md#design-principle-keep-meaning-in-the-source-keep-connections-simple)
governs these choices.

## Plan growth without weakening promises

A delivery story can address selected paths, data variations, or rules across
several use cases. Preserve the behavioral account while work items track changes.
For a rental service, an illustrative sequence might be:

| Increment | Behavior and evidence to consider |
| --- | --- |
| Technical walking skeleton | Exercise a fixed request through reservation and confirmation with simulated payment; evidence of connectivity only |
| Restricted live service | Support one depot and equipment category, including failures possible within that operating scope |
| Broader offering | Add substitute equipment or another payment method, together with its relevant failure handling |
| Learning investigation | Resolve uncertain payment outcomes before accepting a promise about confirmation |

These are planning examples, not product commitments. Moving payment failure to a
later story does not remove it from a live service that accepts payments. Narrow
the offering explicitly or handle the condition. A story map can connect slices
across use cases; it need not become another dotspec record type.

## Review and check

Review goals, boundaries, promises, and unresolved extensions with business and
technical participants. Run `dotspec use-case check <slug>` or `dotspec check` for
structural errors, broken links, and missing-content advice. Advice suggests next
writing steps; it does not require every document to reach the same precision.

This workflow adapts Alistair Cockburn's *The Mini-Book on Use Cases*, v1.1b,
“The use case body parts,” “Sketch-and-expand,” and “Manage precision,” and
*Unifying User Stories, Use Cases, and Story Maps*, second edition, Parts 3–6.
His public [Use Case Foundation](https://alistaircockburn.com/Use%20Case%20Foundation.pdf)
provides a short account. Template choices, link behavior, and authority rules
above are dotspec decisions, not requirements imposed by his method.
