# Use cases

A reviewer declines a refund. Did the interaction succeed? If the goal was *Decide a refund request*, it succeeded once the reviewer knew the finance application had recorded the refusal. For a customer trying to *Obtain a refund*, the same refusal is an unsuccessful outcome. Naming the actor and the goal changes which story you need to tell.

A **scenario** follows one path through an interaction. A **use case** gathers the paths by which a particular actor pursues a goal with a particular system, including alternative successes and failures. The agreed account gives the people requesting, building, testing, and operating the system a way to spot different expectations about its behavior. Cockburn puts that agreement first in *The Mini-Book on Use Cases*, v1.1b, §§2.1–2.3. [Jacobson and Cockburn, *Use-Case Foundation*, v1.1, pp. 1–3](https://alistaircockburn.com/Use%20Case%20Foundation.pdf).

## Whose goal belongs at this boundary?

A customer might seek to recover the cost of an incorrect charge through a merchant, a bank, or another route. A job can describe that progress across those alternatives. A use case selects a system and asks what someone can accomplish with it. In *Decide a refund request*, the primary actor is the independent reviewer, the system is the finance application, and the task ends when the reviewer can establish which decision was recorded. Settling the payment belongs to the wider refund process. Verifying the reviewer’s identity might be one step within the task. Cockburn uses goal levels to distinguish a wider undertaking, a user task, and a subfunction. [JTBD perspectives](jobs-to-be-done.md#choose-a-useful-perspective) · [Cockburn, *The Mini-Book on Use Cases*, v1.1b, §§1.1 and 2.4](https://alistaircockburn.com/Books).

The boundary determines whose behavior the account describes. An external identity service would be a supporting actor if the application called on it. The customer has a stake in the decision even if they never interact with the finance application.

## Which paths belong together?

Approval and refusal both belong in the reviewer’s use case. In either case, the reviewer needs to know what the application recorded against the request. Losing confirmation after submitting approval exposes a different problem. Finance may already have an approval on record while the reviewer has no way to tell. Recording the decision and letting the reviewer establish what was recorded are two parts of the same success guarantee. Until the reviewer can recover that information, their goal remains open. Whatever happens to the confirmation, a refund subject to independent review must stay unauthorized until approval is recorded. This *minimal guarantee* protects Finance at every ending. [Work through the branches](../guides/write-a-use-case.md#explore-what-can-change-the-outcome) · [Review their guarantees](../guides/write-a-use-case.md#review-guarantees-and-ordering).

## Where does each decision live?

The use case owns the interaction and its outcomes. Finance’s rule owns the conditions under which payment needs independent approval. A requirement can carry a distinct software commitment, such as making the recorded decision retrievable after an interruption. In dotspec a branch can cite a rule or a requirement at the point where it matters. A later change to the threshold or the recovery promise then has one place to be made. [Refund rule and its software implications](business-rules.md#how-does-the-rule-reach-the-software).

A delivery story might cover recording a refusal in the first release. The use case still gives the team a place to discuss what happens when a reviewer loses confirmation. Cockburn develops this distinction between the behavioral account and increments of delivery in *Unifying User Stories, Use Cases, and Story Maps*, 2nd ed., §3.5, and *Slice the Problem, Grow the Solution*, v0.95b, ch. 8. [Cockburn’s books](https://alistaircockburn.com/Books) · [Use-Case 3.0, May 2024, pp. 20–23](https://www.ivarjacobson.com/files/use-case_3.0_v1.0.pdf).
