# Define a quality attribute for a system

Start with a system declaration and a concern whose name will help people find and discuss its requirements. Aim for a local definition that lets a second author recognize when the concern applies. A reusable model can supply the starting term. [What a quality means for one system](../explanations/system-quality-attributes.md).

## Select or define the term

From the directory of the system you are specifying, browse the available terms. With a new project, dotspec uses the bundled ISO/IEC 25010:2023 starter offline. For example, suppose you need to decide which work a user can recover after a restart.

```sh
dotspec quality-attribute browse
dotspec quality-attribute browse reliability
dotspec quality-attribute add recoverability --dry-run
```

The preview identifies each file that adoption would create, including the broader Reliability definition. Run `dotspec quality-attribute add recoverability` to create the files. Open Recoverability at the path printed by the command and describe which work the user needs back. Later, adopting Fault tolerance reuses the same local Reliability definition. Repeating a selection shows where the system already declares it.

When several systems could own the term, supply `--system` with the intended source path. To explore another vocabulary for this decision, pass the same `--model` to browse and add. A team can also [choose a project default](define-a-quality-model.md#choose-the-project-default).

For a term you originate, create a file in the system's quality directory. For example, next to a declaration in `specs/`, use `specs/quality/durability/durability.quality-attribute.spec.md`. Put `type: quality-attribute-spec` in YAML frontmatter and name the attribute in the H1. Use its title slug for the filename before `.quality-attribute.spec.md`. The title slug also supplies the ID unless you specify `id` in frontmatter. [To-do durability file](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md).

Under **Definition**, identify what someone needs to preserve, recover or otherwise assess in this system. For recoverability, name the work or state someone needs back after an interruption. Replace a copied prompt with the local definition and keep a citation to its source. Leave the accepted loss or delay with the requirement that sets it. [ISO/IEC 25030:2019, §7.4.2, steps 2–5](https://www.iso.org/standard/72116.html) · [Classified restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md).

## Group the attribute within its owner

If a broader concern would help someone find this term, give the narrower attribute a frontmatter list such as `subattribute_of: [reliability]`. Each ID must name an attribute of the same system. If a parent cannot be found, correct the ID or add the intended declaration. If the checker finds a cycle, trace the parent links and remove the one that points back to a narrower term. For a definition inside a reusable model, edit the model's [named attribute entries](define-a-quality-model.md#define-and-group-the-shared-terms).

## Check it from a requirement

For a requirement in this specification system, add the attribute ID to its frontmatter, for example `quality_attributes: [recoverability]`. From the requirement's owning directory, run `dotspec requirement show preserve-entries`, substituting the slug of the requirement you classified. Inspect its reported system, direct labels and ancestor attributes. An unexpected ancestor calls for a review of the attribute's parent links. If the label is undefined, check which system supplies its declarations.

Read the requirement. Can someone tell which observed result would meet it? If the attribute label is all they have, state the promised result and applicable conditions in the [requirement](write-an-assessable-requirement.md). From the project directory, run `dotspec check` to inspect the declarations and references.

Check the subject under **Applies to**. In ISO/IEC 25030’s maturity example, the authors identify the software they are evaluating, select a measure of remaining defect density and give an acceptable range. In dotspec, a requirement can name the particular part or use of a concept that carries the promise. Suppose the owner has just one open task and restarts the service. Its loss would violate the durability promise even though the owner has never waited for a reply. The restart requirement therefore targets **My checklist** and compares entries across its three views. [ISO/IEC 25030:2019, Annex C, Figure C.1, p. 27](https://www.iso.org/standard/72116.html) · [Restart requirement and its criterion](../../../examples/todo/specs/requirements/restart.requirement.spec.md#applies-to).
