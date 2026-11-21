# Write a job from a concrete situation

Use this guide when a draft could describe almost anyone doing anything. Start
with a recent occasion the performer can recount. If you have only a design
scenario, label it as a hypothesis and record the question you would ask.

## Reconstruct one occasion

Ask “What happened the last time you picked this up again?” Write down what
prompted the return, what the person had already done, and what they consulted.
For the To-do example, suppose the plumber asked for a photo the day after a
call. The call was finished. The commitment now has another step.

## Identify the performer

Name the person whose work is interrupted. Ask who must decide what to do next
and who can judge whether the errand is settled. Here it is the person arranging
the repair. The plumber's request supplies the occasion. Put that role under
**Performer**.

## State the progress

The draft “When I put work down, I want to remember what remains so I can pick
it up again” leaves the reader to guess what someone has forgotten. Try this
revision for the example:

> When an errand needs attention again after a break, I want to recover what
> I already did and what I still owe, so I can take the next action without
> repeating a finished step or overlooking a follow-up.

The revision identifies the decision at the moment of return. It still makes
sense if the person keeps a paper list. Put the statement under **Job**. If a
job story helps you distinguish situation, motivation and intended progress,
use the corresponding sections supplied by the job-story preset.

## Identify the steps

If the work has meaningful parts, name what the performer accomplishes in each.
For example, they identify outstanding commitments and reconsider an errand
after a new request. Use **Steps** when those distinctions help you locate a
difficulty or a desired outcome.

## Describe success

Ask “What would tell you that you can carry on?” In this example, the person
can identify the next action and explain which earlier step is settled. If
time matters, an ODI-style outcome could be “Minimize the time spent working
out which commitments still need attention after an interruption.” Record the
criterion under **Outcomes**. A numerical target needs a basis.

## Record what you know

Under **Evidence**, identify who supplied the account, or label the situation as
an assumption and give the next question. For the hypothetical repair, ask a
participant what they used to recover the next action and what, if anything,
they had to check again. Keep observations and interview detail with their source.
Compare your draft with the [worked job](../../../examples/todo/specs/jobs/resume.job.spec.md)
and the [return scenario](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md).

Read the statement back to the performer. Can they recall an occasion that fits,
describe how they knew enough to act, and point out a detail you guessed wrong?
Use their account to revise the draft. The [job format](../reference/formats/job.md#job-spec-format)
and [JTBD explanation](../explanations/jobs-to-be-done.md) cover syntax and method choices.

## Declare a shared outcome once

An outcome has exactly one declaration owner, a Job or JobStep. Its optional
`applies_to` links identify jobs or steps for which it is a success criterion.
Use a named heading for each outcome and write the full statement underneath.
Add an introductory `**ID:**` when the outcome needs a stable identifier. Put
applicability links under `Applies to`. A job-level outcome uses H3 for its
name and H4 for its fields. A step-level outcome uses H5 and H6. Ordinary links in supporting prose are citations.

In [Resume interrupted work](../../../examples/todo/specs/jobs/resume.job.spec.md),
the job owns “Less reconstruction”. Its Applies to link identifies the particular
step where the person recalls outstanding work.

`job outcome list --step ID` selects declaration ownership. `job map --outcomes`
shows the authored map and outcome declarations; it does not resolve shared
applicability. Use `model needs` for the relationship view across documents.
