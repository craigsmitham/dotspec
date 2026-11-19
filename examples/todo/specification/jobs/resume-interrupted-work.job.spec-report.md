---
generated_by: dotspec
type: job-spec-report
documents: "../../specs/jobs/resume.job.spec.md#resume-interrupted-work"
---
# Resume interrupted work

[Personal task list](<../system.spec-report.md>)

## Job

When I put work down, I want to remember what remains so I can pick it up again.

## Performer

A person managing their own errands or work.

## Supports

[Resume unfinished work](<../system.spec-report.md#resume-unfinished-work>).

## Steps

### Recall outstanding work

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/resume.job.spec.md#recall-outstanding-work"} -->

**ID:** recall
**Stage:** locate

Identify the commitments that still need attention.

**Applies here** · [Less reconstruction](<#less-reconstruction>)

### Reconsider finished work

<!-- dotspec-doc {"type":"job-step-spec-report","documents":"../../specs/jobs/resume.job.spec.md#reconsider-finished-work"} -->

**ID:** reconsider
**Stage:** modify

Recognize when a previously finished commitment needs another action.

## Outcomes

### Less reconstruction

<!-- dotspec-doc {"type":"desired-outcome-spec-report","documents":"../../specs/jobs/resume.job.spec.md#less-reconstruction"} -->

**ID:** less-reconstruction

Minimize the time spent reconstructing outstanding commitments after an interruption.

**Supported by** · [Preserve tasks across restart](<../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

#### Applies to

[Recall outstanding work](<#recall-outstanding-work>).

## Related specifications

**Supported by**

- [My checklist](<../concept-uses/my-checklist.concept-use.spec-report.md>)
- [Return to my list](<../use-cases/return-to-my-list.use-case.spec-report.md>)

[Edit source](<../../specs/jobs/resume.job.spec.md#resume-interrupted-work>)
