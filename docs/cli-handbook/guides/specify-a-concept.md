# Specify a concept

Describe something a person can do, the benefit they get from doing it, and the information the behavior depends on. Start with a case you can walk through with that person. The [Checklist explanation](../explanations/concepts.md) follows an errand that someone finishes and later needs to reopen.

## Connect the actions to a purpose

From your project directory, run `dotspec concept init checklist` to create `specs/concepts/checklist.concept.spec.md`. Use a slug that names your concept. Under **Purpose**, describe what a person gains. Under **Operational principles**, tell a short story that ends where the benefit becomes visible. “An entry is marked complete” leaves open why the change matters. “When I discover a follow-up, I can return the original errand to my open list” gives you a result to check.

Under **State**, describe what the concept must remember for that story to work. Under **Actions**, give each action a level-three heading and describe its effect. In the worked example, Complete changes an entry’s status and preserves its identity. Reopen depends on that choice. When a restriction matters, distinguish a condition the caller assumes from a guard the action enforces. For Rename, the example specifies that a blank title is rejected and the entry keeps its earlier title. [Checklist state and actions](../../../examples/todo/specs/concepts/checklist.concept.spec.md#state).

## Try a different history

Give a reviewer two errands with the same title. Have them complete the first and rename the second. Can they predict which entry can be reopened and what its title will be? If the answer depends on a detail that appears only in the interface mockup, decide whether the state or an action needs that detail. When a system adopts the concept, name the use and link its definition under **Concept uses** in the system or component specification. [Worked declaration](../../../examples/todo/specs/personal-task-list.system.spec.md#my-checklist).

From the project directory containing your draft, run `dotspec concept check checklist` with your own slug. It checks structure and references and suggests missing content. A reviewer’s account of what happens in the changed story tells you whether the description explains the behavior.
