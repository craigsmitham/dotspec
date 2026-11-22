# Write a use case

Ask someone who does the work to describe the last time they finished it. What did they want from the system, and what told them they were done? Sketch that interaction first. Add detail where another person could reasonably expect a different result. Cockburn calls this *sketch-and-expand* in *The Mini-Book on Use Cases*, v1.1b, Part 3. [Cockburn’s books](https://alistaircockburn.com/Books).

## Name the actor, boundary, and finish

Give the use case a verb phrase the primary actor would recognize. In this fictional draft, *Decide a refund request* is a user task. *Obtain a refund* would switch to the customer’s goal. Selecting a request is one step in the reviewer’s work. State the system boundary in the actor’s terms. File location determines default system association, while **System boundary** describes what the narrative covers. If you need to select a system explicitly, put one relative Markdown link to its `*.system.spec.md` file under **Specification system**. [Actors and goal levels](../explanations/use-cases.md#whose-goal-belongs-at-this-boundary).

From your project directory, run `dotspec use-case init decide-refund` to create `specs/use-cases/decide-refund.use-case.spec.md`. Use your own slug for a different task. The brief starter and `--template expanded` use the same format. Choose the extra prompts when the trigger, other participants, or guarantees need discussion.

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

## Explore what can change the outcome

List conditions worth investigating while the successful path is fresh. Start with one alternative way to finish and one circumstance that could leave the actor short of the goal. Give each extension a descriptive heading. Its **Condition** identifies where the path diverges. Its **Handling** explains what happens next, where the interaction resumes or ends, or which decision remains open.

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

Ask the reviewer what must be true at either successful ending, then ask Finance which protection must hold even if the connection fails. Suppose they agree on the following guarantees. The recovery procedure can remain open while the protection is explicit.

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

From the project directory containing your draft, run `dotspec use-case check decide-refund` to check its structure and references and get suggestions about missing content. Then walk the paths with the reviewer. Could a system meet the stated commitments and still strand them after a lost confirmation? What would establish that a refusal was recorded? These questions apply Meyer’s scenario-completeness test to the identified paths. A passing CLI check leaves those judgments to the people reviewing the specification. [Meyer, *Handbook of Requirements and Business Analysis*, §11.3](https://doi.org/10.1007/978-3-031-06739-6).
