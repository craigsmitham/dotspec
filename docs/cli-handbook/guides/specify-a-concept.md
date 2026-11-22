# Specify a concept

Describe something a person can do, the benefit they get from doing it, and the information the behavior depends on. Start with a case you can walk through with that person. The [Checklist explanation](../explanations/concepts.md) follows an errand that someone finishes and later needs to reopen.

## Connect the actions to a purpose

From your project directory, run `dotspec concept init checklist` to create `specs/concepts/checklist.concept.spec.md`. Use a slug that names your concept. For an existing draft, work in its concept file. Replace the starter prompts with your account.

Under **Purpose**, describe what a person gains. Under **Operational principles**, tell a short story that ends where the benefit becomes visible. “An entry is marked complete” leaves open why the change matters. “When I discover a follow-up, I can return the original errand to my open list” gives you a result to check. Jackson calls this characteristic sequence an operational principle. [Jackson, “Concept criteria” (2023), “Three parts”](https://essenceofsoftware.com/tutorials/concept-basics/criteria/).

Under **State**, describe what the concept must remember for that story to work. Under **Actions**, give each action a level-three heading and describe its effect. In the worked example, Complete changes an entry’s status and preserves its identity. Reopen depends on that choice. If the caller must supply an existing entry, describe that obligation under **Caller assumptions** within the action. When the action itself checks a condition, use **Guards** and specify what happens on failure under **Rejections**. For Rename, the example specifies that a blank title is rejected and the entry keeps its earlier title. [Checklist state and actions](../../../examples/todo/specs/concepts/checklist.concept.spec.md#state).

When a detail belongs to entries in one state, ask two questions. Could an entry in that state lack the detail? Could it keep the detail after leaving that state? In Jackson’s Invitation example, every pending invitation has a reminder date. In Checklist, putting an errand in Waiting requires both an expected reply and a follow-up day. Completing it clears both, so reopening it cannot revive the old wait. An ordinary sentence under **Invariants** can express which entries must carry the details. [Jackson, “Concept state” (2023), “Defining State,” Invitation example](https://essenceofsoftware.com/tutorials/concept-basics/state/#defining-state) · [Checklist invariants](../../../examples/todo/specs/concepts/checklist.concept.spec.md#invariants).

## Try a different history

Give a reviewer two errands with the same title. Have them complete the first and rename the second. Can they predict which entry can be reopened and what its title will be? If the reviewer cannot distinguish the errands, record how state identifies an entry. If they can identify the entry but cannot predict its title, clarify what the action preserves. When a system adopts the concept, add a level-three heading naming its use under **Concept uses** in the system or component specification. Under the use’s level-four **Definition** heading, link to the concept file. [Worked declaration](../../../examples/todo/specs/personal-task-list.system.spec.md#my-checklist).

From the project directory containing your draft, run `dotspec concept check checklist` with your own slug. It checks structure and references and suggests missing content. A reviewer’s account of what happens in the changed story tells you whether the description explains the behavior.
