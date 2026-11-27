# Define a quality attribute for a system

Start with a system declaration and a concern whose name will help people find and discuss its requirements. Aim for a local definition that lets a second author recognize when the concern applies. A reusable model can supply the starting term. [What a quality means for one system](../explanations/system-quality-attributes.md).

## Select or define the term

From the directory of the system you are specifying, read the terms of a quality model. dotspec bundles an ISO/IEC 25010:2023 model that you can read offline; each term gives its ID and any **Sub-attribute of** parent. For example, suppose you need to decide which work a user can recover after a restart.

```sh
dotspec docs quality-models/iso-25010-2023
dotspec new quality-attribute recoverability --from iso-25010-2023 --dry-run
```

The preview identifies each file that adoption would create, including the broader Reliability definition. Run `dotspec new quality-attribute recoverability --from iso-25010-2023` to create the files. Open Recoverability at the path printed by the command and describe which work the user needs back. Later, adopting Fault tolerance reuses the same local Reliability definition. Adoption refuses an attribute the system already declares; edit that declaration instead.

Adoption writes to the nearest system that encloses the current directory. To adopt into another system, run the command from its directory or select that directory with `-C`. To adopt from a team vocabulary, pass its source path to `--from`, relative to the current directory. A slug such as `--from team-quality` finds only a model owned by the system or component that contains the current directory. [Define a quality model](define-a-quality-model.md).

For a term you originate, run `dotspec new quality-attribute durability` without `--from`, substituting the term's ID. It creates a blank declaration in the system's quality directory, such as `specs/quality/durability/durability.quality-attribute.spec.md` next to a declaration in `specs/`. To write one by hand, use the same location. Put `type: quality-attribute-spec` in YAML frontmatter and name the attribute in the H1. Use its title slug for the filename before `.quality-attribute.spec.md`. The title slug also supplies the ID unless you specify `id` in frontmatter. [To-do durability file](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md).

Under **Definition**, identify what someone needs to preserve, recover or otherwise assess in this system. For recoverability, name the work or state someone needs back after an interruption. Replace a copied prompt with the local definition and keep a citation to its source. Leave the accepted loss or delay with the requirement that sets it. [ISO/IEC 25030:2019, §7.4.2, steps 2–5](https://www.iso.org/standard/72116.html) · [Classified restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md).

## Group the attribute within its owner

If a broader concern would help someone find this term, give the narrower attribute a frontmatter list such as `subattribute_of: [reliability]`. Each ID must name an attribute of the same system. If a parent cannot be found, correct the ID or add the intended declaration. If the checker finds a cycle, trace the parent links and remove the one that points back to a narrower term. For a definition inside a reusable model, edit the model's [named attribute entries](define-a-quality-model.md#define-and-group-the-shared-terms).

## Check it from a requirement

For a requirement in this specification system, add the attribute ID to its frontmatter, for example `quality_attributes: [recoverability]`. From the project directory, run `dotspec check specs/requirements/preserve-entries.requirement.spec.md`, substituting the path of the requirement you classified. The check reports a label that the requirement's system does not declare. After `dotspec build`, the attribute's generated page lists the requirements it classifies under **Related specifications**; a requirement listed under an unexpected attribute calls for a review of the parent links. If the label is undefined, check which system supplies its declarations.

Read the requirement. Can someone tell which observed result would meet it? If the attribute label is all they have, state the promised result and applicable conditions in the [requirement](write-an-assessable-requirement.md). From the project directory, run `dotspec check` to inspect the declarations and references.

Check the subject under **Applies to**. In ISO/IEC 25030’s maturity example, the authors identify the software they are evaluating, select a measure of remaining defect density and give an acceptable range. In dotspec, a requirement can name the particular part or use of a concept that carries the promise. Suppose the owner has just one open task and restarts the service. Its loss would violate the durability promise even though the owner has never waited for a reply. The restart requirement therefore targets **My checklist** and compares entries across its three views. [ISO/IEC 25030:2019, Annex C, Figure C.1, p. 27](https://www.iso.org/standard/72116.html) · [Restart requirement and its criterion](../../../examples/todo/specs/requirements/restart.requirement.spec.md#applies-to).
