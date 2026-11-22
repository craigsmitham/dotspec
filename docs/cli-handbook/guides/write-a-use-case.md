# Write a use case

Ask someone who does the work to describe the last time they finished it. What did they want from the system, and what told them they were done? Sketch that interaction first. Add detail where another person could reasonably expect a different result. Cockburn calls this *sketch-and-expand* in *The Mini-Book on Use Cases*, v1.1b, Part 3, “Sketch-and-expand, the writing strategy”. [Cockburn’s books](https://alistaircockburn.com/Books).

## Name the actor, boundary, and finish

If the story spans several visits, ask what brought the person to the system each time and what they had accomplished by the time they left. The errand owner records a wait on one visit and decides what to do about the repair on another. Cockburn uses what a person can accomplish in one sitting as a clue to a user task’s scope. A summary use case can follow a goal across several visits. [Cockburn, *The Mini-Book on Use Cases*, v1.1b, §4.2](https://alistaircockburn.com/Books).

Vary one detail of the story to check a precondition. If the owner chose Wednesday to chase the plumber, could they act on a reply on Tuesday? The decision requires a recorded wait. The day the owner returns belongs in the scenario. Cockburn makes this distinction in *Request benefits summary*, where someone can ask about benefits even before filing a claim. [Cockburn, *Writing Effective Use Cases*, 1999 draft, ch. 8, “Precondition”](https://jpaulgibson.synology.me/~jpaulgibson/TSP/Teaching/Teaching-ReadingMaterial/Cockburn00.pdf#page=54).

Name the use case for the result the primary actor wants. In the refund example below, *Decide a refund request* is a user task. *Obtain a refund* would switch to the customer’s goal. Selecting a request is one step in the reviewer’s work. [Actors and goal levels](../explanations/use-cases.md#whose-goal-belongs-at-this-boundary).

Under **System boundary**, describe what the narrative covers in the actor’s terms. If the subject system has a specification, a **Subject system:** link there can identify it, as in the To-do example. File location determines which system the CLI groups the use case under by default. To choose that association explicitly, put one relative Markdown link to the system’s `*.system.spec.md` file under **Specification system**. [To-do system boundary (dotspec checkout)](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#system-boundary).

From your project directory, run `dotspec use-case init decide-refund` to create `specs/use-cases/decide-refund.use-case.spec.md`. Use your own slug for a different task. For a new draft whose trigger, participants or guarantees need discussion, add `--template expanded` to that command. For an existing draft, edit its current sections and add headings where you have more to say. The following passages make up one worked draft.

```markdown
---
type: use-case-spec
---
# Decide a refund request

## Primary actor

An independent refund reviewer.

## Goal

Decide a pending request and establish which decision was recorded against it.

## System boundary

The finance application, from review of the pending request through
confirmation of the recorded decision. Payment settlement happens elsewhere.

## Goal level

User task.
```

## Describe an ordinary success

Describe what the actor and system do that the other can observe. A paragraph can be sufficient. Use an ordered Markdown list when readers need to discuss individual steps. This passage continues the same draft.

```markdown
## Main success scenario

1. The reviewer identifies a pending request.
2. The application supplies the amount, requester, and reason for the refund.
3. The reviewer approves the request.
4. The application records the approval and confirms it to the reviewer
   against the selected request.
```

Check the verbs with the reviewer. “The reviewer approves” describes a decision. “The application confirms” tells the reviewer what became of it. If confirmation merely acknowledges receipt, the last step promises more than the team has agreed.

Walk the last step as the actor. What observation tells them they reached their goal? Jacobson and Cockburn’s *Browse and shop* lists a missing purchase confirmation among the paths to investigate. In the To-do example the repair disappears from Waiting after either Resume or Complete. The owner can establish what happened by opening the destination list and finding the same errand there. Include that observation when it matters to the goal. [*Use-Case Foundation*, v1.1, “Browse and shop,” p. 5](https://alistaircockburn.com/Use%20Case%20Foundation.pdf#page=5) · [To-do return scenario](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#main-success-scenario).

## Explore what can change the outcome

List conditions worth investigating while the successful path is fresh. Start with one alternative way to finish and one circumstance that could leave the actor short of the goal. Give each extension a descriptive heading. Its **Condition** identifies when the path diverges. Its **Handling** explains what happens next, where the interaction resumes or ends, or which decision remains open.

```markdown
## Extensions

### Reviewer declines the request

#### Condition

After examining the request, the reviewer decides to decline it.

#### Handling

The application records the refusal and confirms the declined request to the
reviewer. The requester can see the refusal. The request remains unauthorized
for payment, and the review ends here.

### Confirmation is lost

#### Condition

The reviewer submits approval but loses the connection before receiving
confirmation.

#### Handling

The reviewer cannot tell whether approval was recorded. Open question for
Finance: how can the reviewer establish the decision before trying again?
```

The refusal is an alternative success for this actor. In the second branch, loss of confirmation alone leaves the recorded state unknown. Discuss the handling with Finance before promising a retry or telling the requester that payment can proceed. A precondition such as “the connection stays available” would hide the very situation the team needs to investigate.

## Review guarantees and ordering

At the end of each successful path, what can the actor establish? Read the proposed success guarantee against each ending. Here the reviewer must be able to establish the recorded decision after either approval or refusal. Ask Finance which protection must hold even if confirmation is lost. Suppose they agree on the following guarantees. The recovery procedure can remain open while the protection is explicit.

```markdown
## Success guarantees

The decision is recorded against the intended request, and the reviewer can
establish which decision was recorded.

## Minimal guarantees

A request that requires independent approval stays unauthorized for payment
unless that approval has been recorded.

## Sequence

The reviewer may examine the amount and the reason in either order.
Where independent approval is required, it must precede payment.
```

Numbering the main scenario identifies one path. Use **Sequence** when changing the order would change a promise, and cite the rule that imposes it. If the team agrees to a separate recovery commitment, give it a requirement file and link to it from the lost-confirmation branch. [How the refund rule reaches the software](../explanations/business-rules.md#how-does-the-rule-reach-the-software).

From the project directory containing your draft, run `dotspec use-case check decide-refund` to check its structure and references and get suggestions about missing content. Then walk the paths with the reviewer. For each identified path, inspect the specified behavior supporting each step. If Finance needs to retrieve the decision after a lost confirmation, find the specified behavior that lets the reviewer recover it. If that behavior is still undecided, keep the question in the extension for Finance to resolve. Meyer calls this check of the functionality against identified paths *scenario completeness*. A passing CLI check leaves those judgments to the people reviewing the specification. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §11.3](https://doi.org/10.1007/978-3-031-06739-6).
