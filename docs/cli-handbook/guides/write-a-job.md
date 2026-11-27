# Write a job from a concrete situation

Start with a person's account of a particular occasion, or a situation you can label as a hypothesis. Describe what the performer is trying to accomplish in terms that hold across different ways of doing the work. In the To-do example, someone crosses “Call the plumber” off their list. Ask **“When the call is finished, what does the person still owe?”**

## Start a job specification

From the directory that will own the job, run `dotspec new job personal-commitments`, substituting a new slug for your job. For the example slug this creates `specs/jobs/personal-commitments.job.spec.md`. If the job already has a draft, open that file. Replace the starter prompts with your account. The filename keeps its slug even if you change the title; inspect the draft by its source path. [Project layout](../../reference/layout.md#project-layout).

For a hand-written draft, start with `type: job-spec` in YAML frontmatter and a single H1 naming the job. You can add an optional `id` in frontmatter. Add sections as you learn more. Use each conventional heading at most once and match its capitalization.

## Who owns the result?

Name the **Performer** whose progress you are describing. The person who promised their guest a working sink owns the errand. The plumber owns the repair work. Each could finish their part at a different time.

## Where does the job end?

State what starts the work and how the performer knows it is complete. Name the circumstances that bound it, such as *when* the need arises or *until* what point the work continues. Try the statement against a change of tool or circumstance. “Reopen a task” assumes a list. “Resume an interrupted errand” misses the promise that would still matter if the person remembered every detail.

The To-do **Job** runs from accepting an errand to settling it or confirming that someone else has taken responsibility. Asking a housemate to let the plumber in settles one action. Who still owes the guest a working sink? [Worked job statement](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#job) · [Strategyn on defining a functional job](https://strategyn.com/jobs-to-be-done/jobs-to-be-done-playbook/define-your-market/).

## What work is easy to overlook?

Include **Steps** when they help locate a difficulty. Name each step as something the performer needs to accomplish, whichever solution they use. If the draft only covers the call, ask how the person notices an overdue reply. “Check what remains” names work that takes place even when the person uses no list. The universal job map can help you spot other accomplishments. Record the ones you understand and identify gaps in the evidence. [The eight universal job steps](../explanations/jobs-to-be-done.md#the-universal-job-map).

## How would the performer recognize success?

Give each desired outcome a name and a statement under the job’s or step’s **Outcomes** heading. If several steps share the criterion, declare it once and link those steps under **Applies to**. [Record steps and shared outcomes](#record-steps-and-shared-outcomes).

For a functional job, name what the performer wants to improve, the direction of improvement, how they would judge it, and when it matters. Ulwick gives this form for desired outcomes in *Jobs to Be Done: Theory to Practice* (2016, revised 2023), ch. 4, pp. 93–94. “Make follow-up easier” could mean either of these things. [Ulwick on desired outcomes](https://strategyn.com/customer-needs-through-a-jobs-to-be-done-lens/).

- **Notice a stalled errand.** Minimize the time a commitment goes unattended after an expected reply is overdue.
- **Less reconstruction.** Minimize the time it takes to identify the next outstanding action when returning to an errand after an interruption.

Someone could notice the delay days late and still know immediately whom to call. That makes the two criteria worth judging separately. A whole-job outcome might concern meeting the agreed deadline. For emotional or social progress, the performer's account of what changed may be more informative. Keep a quantified product commitment in a requirement and research ratings with the dated study under **Evidence**. [Choose a useful perspective](../explanations/jobs-to-be-done.md#choose-a-useful-perspective).

## What supports the account?

Use **Evidence** to distinguish observations from hypotheses. The repair errand is fictional. In a real case, a message requesting a photo could show that work continued after the call. Ask the person to walk through what happened next to learn how they recovered the action and whether they needed help. [Evidence and open questions in the To-do job](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#evidence).

For a job story, use **Situation**, **Motivation**, and **Intended progress** to express what prompted someone to act, what they wanted, and why it mattered. The **Job** can then state the story in one sentence: “When [situation], I want to [motivation], so I can [intended progress].” To investigate why they changed approaches on a particular occasion, ask about the switch. [Switching and job stories](../explanations/jobs-to-be-done.md#understand-a-switch).

## Record steps and shared outcomes

Give each step a level-three heading under `## Steps`. Order the steps as you want readers to see the map. Work may repeat or branch. If a stage label helps, add `**Stage:**` beneath the step heading. Labels are free text and may repeat.

This partial draft records one outcome under *Check what remains* and a shared outcome under the job. The links under `Applies to` identify both steps where the shared outcome matters.

```markdown
---
type: job-spec
---
# See personal commitments through

## Steps

### Recall outstanding work

**ID:** recall
**Stage:** locate

Recover the commitments that may need attention.

### Check what remains

Compare what has happened with what was promised.

#### Outcomes

##### Notice a stalled errand

Minimize the time a commitment goes unattended after an expected reply is overdue.

## Outcomes

### Less reconstruction

Minimize the time it takes to identify the next outstanding action after an interruption.

#### Applies to

[Recall outstanding work](#recall-outstanding-work) and [Check what remains](#check-what-remains).
```

Use a heading one level below the outcome name for fields such as `Applies to`, `Supporting content`, `Supports`, or `Derived from`. A job can declare outcomes before you know its steps. [Connect needs to capabilities](map-needs-to-capabilities.md).

Use an explicit ID when you need a link that survives a change to the heading. Put `**ID:**` beneath the step or outcome heading. Job, step, and outcome IDs are case-sensitive and unique within the file. Start with an ASCII letter or digit and use letters, digits, hyphens, or underscores thereafter. Links within the file can also use heading anchors, as in the example. [Reference checking](../../reference/inspection.md#references).

To read the map, run `dotspec build` from your project directory and open the job’s report in `specification/jobs/`. It presents the steps in your order, each with the outcomes it declares. Each step also lists the shared outcomes that apply to it under **Applies here**. For the partial draft above, *Recall outstanding work* declares no outcome but lists *Less reconstruction*. A job without steps says so. If a link is missing, check the heading level of `Applies to` and the target fragment in the job. [Worked applicability links](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#applies-to).

## Check the draft

In your project directory, run `dotspec check specs/jobs/personal-commitments.job.spec.md` with your job’s path. It looks for the job statement, performer, and statements for any declared outcomes. Qualitative criteria can be useful. The check cannot tell whether a functional job states its boundary, steps and measures of success, or whether a job story conveys its situation, motivation and intended progress. Judge those with the questions above.

Fix structural errors such as duplicate headings for the same field, conflicting IDs, and unresolved local links. Wording suggestions are advisory. Read the saved statement, performer and outcomes together in the job's file, or on its generated page after `dotspec build`. Ask the performer to describe the last time they considered the work settled. If asking a housemate to take over still left them chasing the repair, revise the endpoint to say what confirms the transfer. Record that occasion under **Evidence**, or keep the proposed change as a hypothesis until someone can speak to it.
