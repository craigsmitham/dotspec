# Define the scope of a feature

Name the behavior a stakeholder would recognize and explain what they gain from it. From the owning system or component directory, run `dotspec feature init resume-work`, using a slug for your own feature. In the resulting `*.feature.spec.md` file, describe the behavior under **Description** and its benefit under **Purpose**.

Under **Groups**, link to the specifications that together describe this behavior. For example, *Pick up where I left off* links to the return interaction and the restart guarantee. Keep the conditions and criteria with the requirement so every feature that needs the guarantee can refer to the same text. [Worked feature](../../../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md#groups).

Try removing one proposed member from the description. What part of the stakeholder’s expected experience would now be unexplained? If the only connection you can name is the team or screen involved, revisit the scope. [What holds a feature together?](../explanations/features.md).
