# Define a quality model

Identify which systems could reuse the vocabulary and what questions their authors need it to raise. From your project directory, run `dotspec quality-model init team-quality`, substituting a new slug for your model. State its scope under **Description**. If you adapt a published vocabulary, record its version under **Edition** and cite the relevant passages under **Sources**. Say where your definitions depart and why. [What a quality model helps a team decide](../explanations/quality-models.md).

Under **Quality attributes**, give each attribute a level-three heading, a unique `**ID:**` and a short definition. If a more specific term belongs under another, add `**Sub-attribute of:**` with its parent ID. For each term, try a case from two systems. Can both teams recognize the concern even if they require different outcomes? Put each product's conditions and acceptable result in its requirement. [A reusable definition of durability](../../../examples/todo/specs/quality-models/personal-data.quality-model.spec.md#durability).

From the same directory, run `dotspec quality-model check team-quality` with the model's slug. When a system needs one of the terms, use [Define a quality attribute for a system](define-a-system-quality-attribute.md) to select and adapt it.
