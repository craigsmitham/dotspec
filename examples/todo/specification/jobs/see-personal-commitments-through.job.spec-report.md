---
generated_by: dotspec
type: job-spec-report
source: "../../specs/jobs/personal-commitments.job.spec.md#see-personal-commitments-through"
---
# See personal commitments through

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Job

See personal commitments through from accepting responsibility for an errand to establishing that the promised result has been delivered or someone else has agreed to take responsibility for it.

## Performer

The person who has taken responsibility for the errand and must decide what still needs attention alongside their other commitments.

## Situation

Suppose I promise a guest a working sink by Friday. On Monday I call the plumber for a visit time, then leave to collect keys. I decide to chase the reply on Wednesday, while I still have time to arrange another visit. I might pencil a note in my diary or ask my housemate to remind me. Either way a request for a photo puts the next action back with me.

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

**Supported by**

- [Follow-up review](<../capabilities/follow-up-review.capability.spec-report.md>)

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

This map and its outcomes are hypotheses drawn from a fictional errand. Ask someone to walk through the last commitment they carried over several days. What counted as accepting it? What sign told them another action was due? Who could say it was settled? If they still had to chase the result after asking someone else to do the work, ask what would have confirmed a transfer of responsibility.

The [return scenario](<../use-cases/decide-what-to-do-about-a-waiting-errand.use-case.spec-report.md>) illustrates how a person might keep track of a reply between calls. The [participant exercise](<../../README.md#try-the-workflow-with-a-participant>) examines that interaction. Whether keeping the entry helps the person recover their actual next action is a separate question.

## Related specifications

**Supported by**

- [Decide what to do about a waiting errand](<../use-cases/decide-what-to-do-about-a-waiting-errand.use-case.spec-report.md>)
- [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
