# Review a set of requirements

Choose a stakeholder goal and an occasion when several promises must hold together. Give a second reader the goal, the system boundary and the same occasion. Can you each work out what the person should experience, who owes each part of it and whether that outcome serves the goal?

## Follow a goal across the promises

Start at your project's generated `SPEC.md` if it has one. Follow the feature or use case to the requirements for the ordinary path. Include component commitments and environmental assumptions on which that path depends. In a project with authored specifications, run `dotspec model trace` from its root to inspect the recorded relationships. A missing promise may have no link to follow, so walk through an interruption or alternative from the actor's perspective too.

For To-do, start at the system overview, open **Pick up where I left off**, then follow **Decide what to do about a waiting errand**. On Thursday the owner records a reply to chase on Friday. The list acknowledges the change, the service stops and the owner opens Waiting on Friday. Before looking at the test, write down what should be on the screen. Try the same sequence with an entry the owner deleted before stopping the service. [Worked specification](../../../examples/todo/SPEC.md) · [Feature and its grouped specifications](../../../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md#groups).

## Resolve a disagreement at its source

Suppose a draft requires the service to restore the Waiting view exactly as it appeared. Thursday's screen said **Upcoming**. The follow-up requirement calls for **Follow up now** on Friday. Each instruction gives a tester a definite answer. Applied together, they give the owner two answers for the same reply.

Ask what the owner decided on Thursday and what the list determines when it serves Friday's page. The authored restart promise preserves the entry, expected reply and chosen day. The follow-up requirement evaluates that day when Waiting loads. The service owes a durable save before acknowledging the change. For the deleted entry, the availability promise keeps it out of every view. A correction to the recovery promise belongs in its source, where every use of it can follow the same decision. [Restart](../../../examples/todo/specs/requirements/restart.requirement.spec.md) · [Follow-up](../../../examples/todo/specs/requirements/follow-up.requirement.spec.md) · [Acknowledgement](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md) · [Availability](../../../examples/todo/specs/requirements/availability.requirement.spec.md).

## Judge the set within its scope

ISO/IEC/IEEE 29148 gives five characteristics for reviewing a set. Use the cases you just walked through to decide which judgments you can support. [ISO/IEC/IEEE 29148:2018, §5.2.6](https://www.iso.org/standard/72089.html).

| Characteristic | Question to try in your project |
| --- | --- |
| Complete | Can you account for the needs and exceptions inside the chosen boundary? Which outcome still depends on an unresolved decision? |
| Consistent | Can two obligations apply to the same event and require incompatible results? Check duplicated promises, terms, units and time boundaries. |
| Feasible | Can the team meet the promises together within its cost, time and operating constraints? If recovery assumes the original database, try a proposed promise to recover after disk loss. What would supply the missing data? |
| Comprehensible | Can another reader use the specifications to predict the outcome and explain which part of the system owes it? Ask where they had to guess. |
| Able to be validated | If the system met every promise, how could the people who depend on it judge whether it served their purpose under real conditions? |

Record the conflict, gap or open decision beside its owning specification, along with the evidence or person needed to resolve it. Repeat the case after the decision changes. From the project root, `dotspec model review` collects recorded criteria, open questions and inspection findings to bring to that discussion.

## Decide what the evidence supports

In Essence, the **Coherent** state of Requirements concerns a team's shared understanding of scenarios, priorities and conflicts. **Acceptable** calls for stakeholders to judge whether the proposed solution would meet their needs. The later **Addressed** and **Fulfilled** states consider what the software actually delivers. Use those distinctions to ask whose judgment or observation a decision needs. [OMG Essence 1.2, §8.3.2.1, Table 8.4](https://www.omg.org/spec/Essence/1.2/PDF).

Suppose every To-do acceptance check passes. The changed-day test establishes that the saved date survives restart and the cue changes. Can the owner now decide whom to call and what to ask? Invite someone returning to a real errand to show how they decide. Record what the person could recover and what they still had to look up. Until then, the benefit to the owner remains a hypothesis. [Acceptance checks](../../../examples/todo/tests/acceptance.py) · [Participant exercise](../../../examples/todo/README.md#try-the-workflow-with-a-participant).
