---
generated_by: dotspec
type: job-spec-report
source: "../../specs/jobs/personal-commitments.job.spec.md#see-personal-commitments-through"
---
# See personal commitments through

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Job

See personal commitments through from accepting responsibility for an errand to establishing that it is settled or that someone else has agreed to take it on.

## Performer

The person who has taken responsibility for the errand and must decide what still needs attention alongside their other commitments.

## Situation

Suppose I agree to arrange a plumbing repair before a guest arrives. I also have keys to collect. I call the plumber between errands. The next day they ask for a photo before they can book a visit. The call is finished, but I still owe my guest a working sink.

## Supports

[Resume unfinished work](<../goals.md#resume-unfinished-work>).

## Steps

### Establish what is owed

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#establish-what-is-owed"} -->

**Stage:** define

Agree on the result, who is responsible and when it matters. For the repair, a working sink by the visit is the result I have agreed to arrange.

#### Outcomes

##### Clear agreement

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#clear-agreement"} -->

Minimize the likelihood of accepting an errand with a different understanding of completion from the person relying on it.

### Recall outstanding work

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#recall-outstanding-work"} -->

**ID:** recall
**Stage:** locate

Find the commitments that could need attention and recover the details needed to act. After a break I may have to consult a message, another person or my own notes.

**Applies here**

- [Less reconstruction](<#less-reconstruction>)

### Arrange what is needed

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#arrange-what-is-needed"} -->

**Stage:** prepare

Get the details or materials needed for the next action into usable form. To request a visit I may need the address, a description of the leak and times when someone can let the plumber in.

### Decide which errand can go ahead

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#decide-which-errand-can-go-ahead"} -->

**Stage:** confirm

Check which commitments are ready to act on given the time, access and people available. I can collect the keys before the shop closes while I wait for the plumber to reply.

#### Outcomes

##### Avoid wasted trips

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#avoid-wasted-trips"} -->

Minimize the likelihood of starting an errand before a necessary prerequisite is in place.

### Carry out the next action

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#carry-out-the-next-action"} -->

**Stage:** execute

Do the work that advances the commitment. Call the plumber to arrange a visit or collect the keys from the shop.

### Check what remains

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#check-what-remains"} -->

**Stage:** monitor

Compare what has happened with what was promised. If I am waiting on a reply, determine when I need to follow up for the repair to happen in time.

**Applies here**

- [Less reconstruction](<#less-reconstruction>)

#### Outcomes

##### Notice a stalled errand

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#notice-a-stalled-errand"} -->

Minimize the time a commitment goes unattended after an expected reply is overdue.

### Revise the plan

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#revise-the-plan"} -->

**ID:** reconsider
**Stage:** modify

Work out what a changed circumstance requires. When the plumber asks for a photo, sending it becomes the next action. If a visit before the deadline becomes impossible, agree a new arrangement with my guest.

### Settle the commitment

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#settle-the-commitment"} -->

**Stage:** conclude

Establish that the agreed result has been delivered, or confirm who has accepted the remaining responsibility. A completed repair settles the errand. A booked visit still needs someone to let the plumber in.

## Outcomes

### Commitments fulfilled

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#commitments-fulfilled"} -->

Minimize the likelihood of missing an agreed deadline for a personal commitment.

### Less reconstruction

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/personal-commitments.job.spec.md#less-reconstruction"} -->

**ID:** less-reconstruction

Minimize the time it takes to identify the next outstanding action when returning to an errand after an interruption.

**Supported by**

- [Local task storage](<../capabilities/local-task-storage.capability.spec-report.md>)
- [Preserve tasks across restart](<../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

#### Applies to

[Recall outstanding work](<#recall-outstanding-work>) and [Check what remains](<#check-what-remains>).

## Evidence

This map and its outcomes are hypotheses drawn from a fictional errand. Ask someone to walk through the last commitment they carried over several days. What counted as accepting it? What sign told them another action was due? Who could say it was settled? A handoff that left them still responsible would change the proposed endpoint.

The [return scenario](<../use-cases/return-to-my-list.use-case.spec-report.md>) illustrates the list's contribution when a finished action needs follow-up. The [participant exercise](<../../README.md#try-the-workflow-with-a-participant>) examines that interaction. Whether keeping the entry helps the person recover their actual next action is a separate question.

## Related specifications

**Supported by**

- [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
- [Return to my list](<../use-cases/return-to-my-list.use-case.spec-report.md>)
