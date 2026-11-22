# Define a quality model

Use a reusable model when authors of several systems need a shared vocabulary. Have two systems or proposed uses in mind. If the term matters only to one system, give it a [system quality attribute](define-a-system-quality-attribute.md).

## Set the scope and basis

From your project directory, run `dotspec quality-model init team-quality`, substituting a new slug for your model. If the team already uses a bundled vocabulary, `dotspec quality-model init team-quality --from iso-25010-2023` starts an editable copy instead. Use one of these routes. If the model already exists, edit its source.

Give the model a title and replace the prompts in **Description** with its purpose and intended uses. Identify any edition and sources you drew on under **Edition** and **Sources**. If you adapt an external term, tell readers which part of the source meaning you kept and why. For example, ISO/IEC 25010 treats recovery of affected data and restoration of the system's state together. A model specifically about personal records might isolate retention of acknowledged changes for discussion across products. [ISO/IEC 25010:2023, §3.5.4](https://www.iso.org/standard/78176.html) · [Personal data quality](../../../examples/todo/specs/quality-models/personal-data.quality-model.spec.md).

## Define and group the shared terms

Under `## Quality attributes`, give each attribute a level-three heading, a unique ID line such as `**ID:** durability`, and a definition that another team can apply. To group a narrower term under a broader one, put the parent ID on a `**Sub-attribute of:**` line. List multiple parents with commas when several groupings help readers find the term. Omit the line for a root. If the checker flags an older `Broader` or `Parents` field, replace it with `**Sub-attribute of:**` and retain the intended parent IDs. [What a quality model helps a team decide](../explanations/quality-models.md).

Ask what would count in each system. A task list and an expense tracker can both consider whether acknowledged work survives a restart. Their authors can apply that definition to different records. If your shared definition specifies how quickly the expense tracker must recover receipts, move that commitment into the tracker's requirement. [From a shared meaning to a specific promise](../explanations/system-quality-attributes.md#from-a-shared-meaning-to-a-specific-promise).

From the same directory, run `dotspec quality-model check team-quality` and `dotspec quality-model show team-quality` with the model's slug. Inspect the definitions as the second team would read them. Correct any missing parent ID or circular grouping in this model, then run the check again. Once each team can select the relevant terms, use [Define a quality attribute for a system](define-a-system-quality-attribute.md) to tailor them.
