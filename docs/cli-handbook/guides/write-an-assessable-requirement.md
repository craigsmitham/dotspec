# Write an assessable requirement

Start with a decision about what the system owes someone. Aim for a requirement that lets a reviewer judge a result, including a case just outside its scope. Have the decision or a reference to its source at hand. [What makes a requirement assessable?](../explanations/requirements.md).

## Give the commitment a home

From the owning system or component directory, run `dotspec requirement init preserve-entries`, choosing a new slug for your requirement. Open the path it prints. If a requirement already owns the decision, edit that file. Replace the starter prompts and give the file a title that names the promised result.

Use **Statement** for one required behavior or property. Use **Rationale** for the need or decision that makes it necessary. Cite a governing policy or use case where it explains the commitment. If you found the requirement written inline in a scenario, move the independent commitment into the requirement file and link to it from the scenario. Keep the actor's outcome in the use case. [ISO/IEC/IEEE 29148:2018, §§5.2.4–5.2.5](https://www.iso.org/standard/72089.html) · [Where the promise lives](../explanations/requirements.md#where-does-the-promise-live).

## Decide when it applies

Under **Applies to**, identify the subject and the circumstances in which the promise holds. Change one condition and judge the result again. The To-do restart promise covers acknowledged changes on functioning local storage. A missing task after acknowledgement can violate the promise. A request still awaiting acknowledgement falls outside its scope. Qualify the workload, population or exclusions when they change what the system owes. [Restart conditions](../../../examples/todo/specs/requirements/restart.requirement.spec.md#applies-to).

## State how to judge a result

Try one case that would meet the statement and one that would violate it. A qualitative rule such as “a deleted entry stays absent from every view” already distinguishes them. If a comparison needs explanation, add `## Criterion`. For a measured requirement, use `### Measure` for what to observe and how to count or calculate it, and `### Target` for the acceptable value or range. The To-do example counts changed, missing and unexpected entries. Without an acceptable count, two people could agree on the observation and disagree about whether to accept it. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §4.12](https://doi.org/10.1007/978-3-031-06739-6) · [Worked measure and target](../../../examples/todo/specs/requirements/restart.requirement.spec.md#criterion).

Give each bound one authoritative statement. If it is set elsewhere, reference it under **Criterion**. If a measured draft has no agreed bound yet, record the open decision and who can resolve it. Measurement advice is a prompt to inspect the criterion and its operating conditions. A separately written target is useful once the project has a value to defend.

## Check the classification and the promise

If a quality helps readers find this requirement, put its system attribute ID in the frontmatter list `quality_attributes`. For an undefined label, check its spelling and its definition in the requirement's specification system. If the concern belongs there and needs a definition, follow [Define a quality attribute for a system](define-a-system-quality-attribute.md).

From the same directory, run `dotspec requirement check preserve-entries` and `dotspec requirement show preserve-entries` with your slug. Follow the reported source if you selected the wrong draft. Read the promise and its conditions together, then ask someone to judge the two cases. Different answers identify terms or conditions to revisit. Link to any test or example that exercises the promise and explain which observations it makes. [Evidence in the restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md#evidence).
