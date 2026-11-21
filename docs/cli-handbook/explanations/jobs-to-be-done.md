# What makes something a job to be done?

A person calls a plumber, crosses the call off a list and leaves to collect some keys. What would count as finishing the job? If they agreed to get the sink fixed before a guest arrives, finishing the call leaves a commitment open. They might use a notebook, rely on a text exchange or remember the arrangement unaided. The commitment gives each of those approaches its purpose.

## Choose a useful perspective

The handbook starts with Anthony Ulwick's **Outcome-Driven Innovation (ODI)** for requirements discovery. Its anchor is a functional job a performer wants to accomplish across different solutions. A job map breaks the work into accomplishments and gives researchers specific places to ask how performers judge success. In the repair example, checking what remains is an accomplishment. Time spent identifying the next action after a break is one measure of success. ODI research asks performers how important that criterion is and how satisfied they are with current approaches. Those ratings belong with the study because their interpretation depends on who took part and when. See Ulwick, *Jobs to Be Done: Theory to Practice*, chs. 2 and 4. [Strategyn on customer needs](https://strategyn.com/customer-needs-through-a-jobs-to-be-done-lens/).

Christensen's theory asks what **progress in particular circumstances** explains a choice. The same person may use a notebook for their own errands and agree to exchange texts with a guest about the repair. A written exchange might help them coordinate access (functional), feel reassured (emotional), or be seen as someone who keeps a promise (social). Each is a possible explanation to investigate with the person. Their choice can also involve carrying on as before. [Christensen Institute, *Jobs to Be Done Theory*](https://www.christenseninstitute.org/theory/jobs-to-be-done/).

These inquiries can inform each other. A switch account might reveal how much work it took to establish who owed the next action. That is a reason to investigate the **Monitor** step with other performers. A job map might reveal that people have several ways to confirm a handoff. A switch interview can explain why one of those ways became worth trying on a particular occasion. Evidence about one choice has the scope of that case. [Jim Kalbach on functional jobs and aspirations](https://experiencinginformation.com/2020/10/23/drills-and-milkshakes-say-yes-to-the-mess/).

## The universal job map

Bettencourt and Ulwick describe eight kinds of accomplishment that can make up a functional job. They are useful when an obvious activity, such as making a phone call, hides work on either side of it. The stage names give authors a way to look for that work. The step names in a particular map describe what matters to its performer. [Bettencourt and Ulwick, *The Customer-Centered Innovation Map*](https://hbr.org/2008/05/the-customer-centered-innovation-map) · [Strategyn, *Build Your Job Map*](https://strategyn.com/jobs-to-be-done/jobs-to-be-done-playbook/build-your-job-map/).

| Universal step | What it brings into view | In the repair errand |
| --- | --- | --- |
| Define | The result to achieve and the constraints that matter | Agree that the sink needs to work before the guest arrives. |
| Locate | The people, materials or information the work depends on | Find the plumber’s reply and recover what they need. |
| Prepare | The work of getting those inputs into usable condition | Assemble the address, a description of the leak and possible visit times. |
| Confirm | How the performer knows they can proceed | Establish that someone can let the plumber in at the proposed time. |
| Execute | The action that advances the commitment | Arrange the visit. |
| Monitor | Signs that the work is on track or has stalled | Notice that the plumber has yet to confirm. |
| Modify | Adjustments when circumstances change | Supply the requested photo or arrange a different visit. |
| Conclude | What allows the performer to consider the work settled | Check that the sink works or confirm who has taken responsibility. |

A performer may check readiness several times or finish an errand without changing the plan. Some stages may have no distinct step worth recording. The value of the map is the work it reveals and the criteria it helps explain. [Worked job and its outcomes](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#steps). See also Jim Kalbach, *The Jobs to Be Done Playbook* (2020), ch. 3.

## Understand a switch

Bob Moesta and the Re-Wired Group's switching model addresses “Why change now?” A useful account identifies the previous approach, what made someone reconsider it, the alternatives they took seriously, and the trade-off they accepted. A timeline can tie the explanation to actual events, from first doubts to trying the new approach. If a buyer chose for someone else, the account explains the buyer’s choice. See Moesta, *Learning to Build*, ch. 3. [Re-Wired Group on JTBD interviews](https://therewiredgroup.com/learn/complete-guide-jobs-to-be-done/) · [The buyer’s timeline](https://therewiredgroup.com/learn/the-six-stages-of-the-buyers-timeline/).

Suppose someone starts keeping commitments in a notebook. Losing an agreement in a long message thread might supply the **push** of the old situation. Seeing all their errands in one place might supply the **pull** of the new approach. Worry about leaving the notebook at home could create **anxiety**. The **habit** of replying in the same thread where a request arrived might explain why messages remained attractive. These are candidate explanations. A switch interview can establish which forces the person actually experienced and which mattered to their choice.

The switch timeline belongs with the evidence for that decision. The functional job map describes how someone sees commitments through, whether they switch tools or carry on as before.

A **job story** is a writing form for carrying a finding into design discussion. For example, “When someone asks me for an update after a few days, I want to recover what we agreed, so I can give them an answer they can rely on.” Use `dotspec job init give-an-update --preset job-story` to draft those three fields. Keep the source of the story under **Evidence**, including whether it is still a hypothesis. [Paul Adams, *How we accidentally invented Job Stories*](https://www.intercom.com/blog/accidentally-invented-job-stories).

## From a job map to a user needs map

Richard Allen's **User Needs Mapping** starts with who needs what and follows the dependencies that make it possible to meet that need. A broad job can anchor the map. For a closer look at one part of the job, an ODI desired outcome can supply a user need. Allen uses short labels on the diagram and permits starting with task-like language. Preserve the full criterion so a reader can recover what success means. For example, “identify the next action” alone loses the time cost and the circumstance of returning after an interruption. See *User Needs Mapping*, first edition, ch. 4, step 3 (pp. 39–45) and steps 4–7. [Allen's companion site](https://www.userneedsmapping.com/) · [Worked moviegoer map](https://www.userneedsmapping.com/moviegoer-walkthrough).

The To-do example follows one need from the person doing the work to the software that supports it.

- **Who needs what?** The person responsible for the errand wants to spend less time identifying the next outstanding action after a break. The map can label that need “Less reconstruction.” [Full desired outcome](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#less-reconstruction).
- **What ability helps?** Retaining accepted changes allows the person to consult the same entry tomorrow. Their own practice of recording what was agreed also matters. A stored title such as “Plumber” may tell them very little. [Local task storage capability](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).
- **What provides the ability?** The task service provides local storage. Its restart requirement states what entry data must survive. Allen's map can go on to show which teams own or depend on the capabilities, exposing coordination work. [Task service](../../../examples/todo/app/specs/task-service.component.spec.md) · [Restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md).

In dotspec, `job map personal-commitments --outcomes` shows steps in authored order and outcomes where they were declared. `model needs` follows declared relationships across needs, capabilities and specifications. It can show the task service's contribution to this outcome. Finding out whether a returning person can act sooner takes evidence from doing the job. Keep workshop maps and team overlays with the study or organizational decision that used them. [Map needs to capabilities](../guides/map-needs-to-capabilities.md#connect-needs-to-capabilities-and-specifications).
