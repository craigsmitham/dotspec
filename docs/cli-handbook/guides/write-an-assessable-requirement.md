# Write an assessable requirement

Start with a decision about what the system owes someone. Aim for a requirement that lets a reviewer judge a result, including a case just outside its scope. Have the decision or a reference to its source at hand. [What makes a requirement assessable?](../explanations/requirements.md).

## Give the commitment a home

From the owning system or component directory, run `dotspec new requirement preserve-entries`, choosing a new slug for your requirement. Open the path it prints. If a requirement already owns the decision, edit that file. Replace the starter prompts and give the file a title that names the promised result.

Use **Statement** for one required behavior or property. Use **Rationale** for the need or decision that makes it necessary. Cite a governing policy or use case where it explains the commitment. If you found the requirement written inline in a scenario, move the independent commitment into the requirement file and link to it from the scenario. Keep the actor's outcome in the use case. [ISO/IEC/IEEE 29148:2018, §§5.2.4–5.2.5](https://www.iso.org/standard/72089.html) · [Where the promise lives](../explanations/requirements.md#where-does-the-promise-live).

## Check the promise against its source

Suppose a draft says a reply marked for Friday becomes due on Saturday. Two testers could agree exactly how to check it. Compare their answer with the owner's decision to chase the reply **on Friday**. The draft is verifiable and still gets the day wrong. Correct the requirement against the [follow-up policy](../../../examples/todo/specs/business-rules/follow-up.business-rule.spec.md#statement).

Use the source to review **correctness** and **necessity**. What did the person actually ask for? If you removed this obligation, which need would go unmet? Check that the detail is **appropriate** to the subject. The owner needs an acknowledged change to survive a restart. The ordering of a database commit and a browser response belongs with the service that controls it. Ask whether each condition qualifies the same promise. Preserving an acknowledged entry across restart is a **singular** obligation even though the details to preserve depend on whether the errand was open, waiting or completed. Split obligations when they can change independently. [ISO/IEC/IEEE 29148:2018, §5.2.5](https://www.iso.org/standard/72089.html) · [Task service commitment](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md).

Before choosing the item for development, ask the team what makes it feasible, what else it affects and why it takes priority over another item. For restart, access to the same functioning database matters to that judgment. Essence's *Described* state asks for this kind of justification and understanding. Its *Verified* state concerns evidence from an implementation. [OMG Essence 1.2, Annex A.3.2.1, Table A.5](https://www.omg.org/spec/Essence/1.2/PDF).

## Decide when it applies

Under **Applies to**, identify the subject and the circumstances in which the promise holds. Name the subject on a line that begins with `**Targets:**`, as the restart requirement does for **My checklist**. After `dotspec build`, the subject’s report lists the requirement under **Applies here**. A link elsewhere in the section remains an ordinary citation. Change one condition and judge the result again. The To-do restart promise covers acknowledged changes on functioning local storage. A missing task after acknowledgement can violate the promise. A request still awaiting acknowledgement falls outside its scope. Qualify the workload, population or exclusions when they change what the system owes. [Restart conditions](../../../examples/todo/specs/requirements/restart.requirement.spec.md#applies-to).

## State how to judge a result

Try one case that would meet the statement and one that would violate it. A qualitative rule such as “a deleted entry stays absent from every view” already distinguishes them. If a comparison needs explanation, add `## Criterion`. For a measured requirement, use `### Measure` for what to observe and how to count or calculate it, and `### Target` for the acceptable value or range. The To-do example counts changed, missing and unexpected entries. Without an acceptable count, two people could agree on the observation and disagree about whether to accept it. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §4.12](https://doi.org/10.1007/978-3-031-06739-6) · [Worked measure and target](../../../examples/todo/specs/requirements/restart.requirement.spec.md#criterion).

Give each bound one authoritative statement. If it is set elsewhere, reference it under **Criterion**. If a measured draft has no agreed bound yet, record the open decision and who can resolve it. Measurement advice is a prompt to inspect the criterion and its operating conditions. A separately written target is useful once the project has a value to defend.

## Check the classification and the promise

Durability describes a property of the product. Verifiability describes whether a reviewer can judge fulfillment of the written promise. If a product quality helps readers find this requirement, put its system attribute ID in the frontmatter list `quality_attributes`. For an undefined label, check its spelling and its definition in the requirement's specification system. If the concern belongs there and needs a definition, follow [Define a quality attribute for a system](define-a-system-quality-attribute.md).

From the same directory, run `dotspec check specs/requirements/preserve-entries.requirement.spec.md` with your file's path. Read the promise and its conditions together, then ask someone to judge the two cases. If readers disagree, look for ambiguous terms. If both have to supply a missing condition, complete the promise. Check it against any authoring conventions your project has adopted. List any test or example that exercises the promise as a companion, and explain under **Evidence** which observations it makes. [List a companion](list-a-companion.md#choose-a-role) · [Evidence in the restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md#evidence).

When several promises apply to the same user goal, [review the set together](review-a-set-of-requirements.md).
