# Write a job from a concrete situation

Start with an occasion someone remembers. Use Outcome-Driven Innovation (ODI) to work out whose job it is, where it begins and ends, and how that person judges success. You can start a draft with `dotspec job init personal-commitments --preset odi` and edit the resulting Markdown as you learn.

## Reconstruct one occasion

Ask “Take me back to the last time this happened. What had you agreed to do? What happened before you opened the list? What happened after you closed it?” Record the person's words, the occasion and anything they consulted. If you start from an imagined situation, label it as a hypothesis under **Evidence**.

Suppose someone calls the plumber, marks the call complete, and moves on to collect some keys. The plumber asks for a photo the next day. “Reopen a task” describes the interface action. “Recover the next action after a break” describes a difficulty that could also arise with a paper list. Follow the errand back to the original promise and forward to the result the person owes.

For a job story, put the occasion under **Situation** (the plumber asked for a photo the day after the call), what the person wants to do under **Motivation** (recover what they agreed), and what that enables under **Intended progress** (get the visit booked in time). Check all three against the account.

## Identify the performer

Ask who took responsibility and who can say the work is done. In the example, the person arranging the repair has made a commitment to their guest. The plumber performs a different job. Under **Performer**, describe the person whose decisions the map will follow. If a different person chose or bought the tool, seek their account of that choice separately.

## Identify the steps

Try a first boundary, such as “resume an interrupted errand.” Walk through Bettencourt and Ulwick's eight universal job steps with the performer. For each, ask what they need to accomplish. Listen for work outside the boundary you guessed, including work done by hand or in conversation. Use the person's terms for the resulting step names. [Bettencourt and Ulwick, *The Customer-Centered Innovation Map*](https://hbr.org/2008/05/the-customer-centered-innovation-map) · [Strategyn, *Build Your Job Map*](https://strategyn.com/jobs-to-be-done/jobs-to-be-done-playbook/build-your-job-map/).

| Universal step | Ask the performer | What the repair errand brings into view |
| --- | --- | --- |
| Define | What result did you take responsibility for? | Establish what is owed. A working sink by the guest's visit is a different promise from making a call. |
| Locate | Where could you find what you needed to act? | Recall outstanding work. The plumber's reply may be the only record of what they need. |
| Prepare | What had to be put in order? | Arrange the information and access needed for a visit. |
| Confirm | What told you it was practical to go ahead? | Decide which errand can go ahead while another awaits a reply. |
| Execute | What did you do that advanced the commitment? | Arrange the visit or send the requested photo. |
| Monitor | How did you know whether things were on track? | Check what remains, including a reply that is now overdue. |
| Modify | What changed, and what did you reconsider? | Revise the plan when the plumber needs more information or can no longer come in time. |
| Conclude | What allowed you to stop following up? | Settle the commitment by checking the repair or confirming that someone accepted the handoff. |

Pursue an unexpected answer. If “done” means the call ended, ask who is still responsible for getting the sink fixed. If someone else has taken over, ask what counted as acceptance. Either answer can move the end of the job. Steps can repeat, several can belong to one stage, and some stages may have no distinct work here. Write down gaps to investigate. Keep only meaningful accomplishments under **Steps**. See Jim Kalbach, *The Jobs to Be Done Playbook* (2020), ch. 3.

## State the progress

Read the map back as a whole. Which objective explains why the same person would do these things? Use a verb and its object, adding context that fixes the boundary. In this case the map suggests this **Job** statement.

> See personal commitments through from accepting responsibility for an errand to establishing that it is settled or that someone else has agreed to take it on.

The earlier draft about resuming work accounts for a difficult moment. The broader statement also explains why the person clarifies a promise, waits for a reply, or checks a handoff. Test the statement against an uninterrupted errand, a different errand and a different way of keeping track. Narrow it if the performer or the test of completion changes. Record occasions that challenge it. [Strategyn, *Define Your Market*](https://strategyn.com/jobs-to-be-done/jobs-to-be-done-playbook/define-your-market/).

## Describe success

At each step ask “What takes too long, goes wrong or has to be done again? How can you tell?” For an ODI desired outcome, name the direction of improvement, one measure, the thing being measured and the circumstance that matters. “Make follow-up easier” might mean either of these:

- **Notice a stalled errand.** Minimize the time a commitment goes unattended after an expected reply is overdue.
- **Less reconstruction.** Minimize the time it takes to identify the next outstanding action when returning to an errand after an interruption.

The first delay can last days even when the person immediately remembers what to do. Keeping the criteria distinct lets people rate them independently. Ask about whole-job success too. Someone might remember every action and still miss the deadline. Keep importance and current satisfaction ratings with the dated study, its participants and method. [Ulwick, *Customer Needs Through a Jobs-to-be-Done Lens*](https://strategyn.com/customer-needs-through-a-jobs-to-be-done-lens/).

Declare an outcome once under the job or a step. In the worked example, the job owns **Less reconstruction**, and its **Applies to** section links both steps where it matters. A restart requirement supports that outcome. A participant can still recover every saved entry and have no idea whom to call next. Ask them to identify their actual next action and explain what helped. [Worked job](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md) · [Job spec format](../reference/formats/job.md#job-spec-format).

## Review the draft with a performer

Run `dotspec job map personal-commitments --outcomes` to read the map in order, then `dotspec job check personal-commitments --preset odi` for structural checks and wording suggestions. Take the map to someone who has recently done the work. Have them walk it against their experience, point out a missing piece and explain how they judged the result. Revise **Job**, **Steps** and **Outcomes** in light of that account and record its source and limits under **Evidence**.

If your question is why a person changed how they handle errands, reconstruct that choice in a switch interview. If you need to identify which abilities support a particular outcome, take it into a User Needs Map. [Switching and progress](../explanations/jobs-to-be-done.md#understand-a-switch) · [From a job map to a user needs map](../explanations/jobs-to-be-done.md#from-a-job-map-to-a-user-needs-map).
