# Define a quality attribute for a system

Choose a concern whose name will help people find and discuss requirements for this system. If a reusable model supplies the term, preview the selection from your project directory. For example, with a system declared at `specs/new-system.system.spec.md`, run this preview.

```sh
dotspec quality-model apply iso-25010-2023 --system specs/new-system.system.spec.md --attributes recoverability --dry-run
```

Substitute your system's path and the model and attribute IDs you intend to use. The preview includes the selected attribute and its ancestors. Repeat the command without `--dry-run` to create the files. For a term you originate, put a `*.quality-attribute.spec.md` file under the owning system's quality directory with `type: quality-attribute-spec`, the attribute name as its title, and a **Definition**. [To-do durability file](../../../examples/todo/specs/quality/durability/durability.quality-attribute.spec.md).

Define what the quality means here. For recoverability, for example, identify the work or state someone needs to recover. Cite the source definition if you adapt one. To group a narrower concern under a broader one, give its `subattribute_of` field the parent IDs from the same system. Trace the chain back to its broadest terms and remove any link that leads back to itself. [What a quality means for one system](../explanations/system-quality-attributes.md).

To classify a requirement, add the system's attribute ID to its `quality_attributes` list. Read the requirement on its own. Its statement, applicability and criterion should let someone tell whether a particular result is acceptable. From the project directory, run `dotspec check` to inspect the declarations and references. [Classified restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md).
