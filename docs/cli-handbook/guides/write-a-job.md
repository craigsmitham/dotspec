# Write a job from a concrete situation

A useful job description lets a reader explain who wants to accomplish what, where that work ends, and how the person judges success. In the To-do example, someone crosses “Call the plumber” off their list. The question that reveals the job is **“When the call is finished, what does the person still owe?”**

## Who owns the result?

Name the **Performer** precisely enough that a reader knows whose decisions count. The person who promised their guest a working sink owns the errand. The plumber owns the repair work. Both have a reason to care about the same visit, but each could finish their part at a different time.

## Where does the job end?

Write the **Job** in terms that would still make sense if the performer used a notebook, exchanged messages, or kept the agreement in their head. “Reopen a task” depends on a particular tool. “Resume an interrupted errand” captures a real difficulty, but someone who remembers the errand perfectly still owes their guest a working sink. The worked example uses this boundary.

> See personal commitments through from accepting responsibility for an errand to establishing that it is settled or that someone else has agreed to take it on.

That statement explains why the person checks for a reply, sends a photo, or confirms a handoff. It also gives a reviewer something to challenge. If asking a housemate to meet the plumber leaves the original person responsible for the repair, that handoff settles only one action. [Worked To-do job](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md) · [Strategyn on defining a functional job](https://strategyn.com/jobs-to-be-done/jobs-to-be-done-playbook/define-your-market/).

## What work is easy to overlook?

When a draft starts and ends inside the product, the universal job steps from Outcome-Driven Innovation (ODI) can reveal a larger job. In the repair example, **Define** brings the original promise into view. **Monitor** reveals an overdue reply while there is still time to act. **Conclude** raises the question of who can say the commitment is settled. Those discoveries change the job statement itself. [The eight universal job steps](../explanations/jobs-to-be-done.md#the-universal-job-map).

Include **Steps** if they help a reader understand the work or locate a difficulty. Give each step a name the performer could recognize, such as “Check what remains.” A partial map can be useful when only part of the work is understood.

## How would the performer recognize success?

A useful **Outcome** distinguishes two ways the job could turn out. “Make follow-up easier” leaves open at least two quite different improvements.

| Possible desired outcome | What a performer could judge |
| --- | --- |
| Minimize the time a commitment goes unattended after an expected reply is overdue. | How long did the stalled errand escape notice? |
| Minimize the time it takes to identify the next outstanding action when returning to an errand after an interruption. | Once I returned to the errand, how long did it take to work out what to do? |

Someone could notice the delay days late and still know immediately whom to call. ODI expresses each success criterion with a direction and a measure so that performers can assess it independently. A whole-job criterion might concern meeting the agreed deadline. A job about reassurance or trust may call for the person’s own account of what changed. [Ulwick on desired outcomes](https://strategyn.com/customer-needs-through-a-jobs-to-be-done-lens/) · [JTBD perspectives](../explanations/jobs-to-be-done.md#choose-a-useful-perspective).

When the decision is which capabilities could help, a desired outcome can anchor a User Needs Map. In To-do, the need to recover the next action explains why keeping an entry available across a restart matters. [From a job map to a user needs map](../explanations/jobs-to-be-done.md#from-a-job-map-to-a-user-needs-map).

## What supports the account?

Under **Evidence**, give readers enough context to tell what was observed, what was inferred, and what remains open. The repair errand is fictional. A real message requesting a photo could establish that the errand continued after the call. It would take more evidence to say that the person struggled to remember their next action or that a saved entry helped.

If a particular occasion explains the progress sought, a job story can express it in **Situation**, **Motivation**, and **Intended progress**. If the question is why someone changed how they handle errands, a switch account can explain the choice. [Switching and job stories](../explanations/jobs-to-be-done.md#understand-a-switch).

To start an ODI draft, use `dotspec job init personal-commitments --preset odi`. To inspect its structure and wording, use `dotspec job check personal-commitments --preset odi`. The [job spec format](../reference/formats/job.md#job-spec-format) shows how to declare steps, outcomes and the links between them.
