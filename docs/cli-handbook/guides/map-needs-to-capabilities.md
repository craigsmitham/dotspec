# Map needs to capabilities

Use this guide to trace which abilities support a job and which parts of the system provide them. Start with the [job and its desired outcomes](write-a-job.md).

## Connect needs to capabilities and specifications

Choose a user need at the level of the decision. A broad job can anchor a workshop about the whole service. For a question about delays after an interruption, use the specific desired outcome. Ask which abilities support it, what each ability depends on and who or what provides it. Preserve the user's full criterion behind the short map label. [JTBD and User Needs Mapping](../explanations/jobs-to-be-done.md#from-a-job-map-to-a-user-needs-map) · [Richard Allen's User Needs Mapping](https://www.userneedsmapping.com/).

| Type | Role |
| --- | --- |
| Capability | An ability needed to deliver value, possibly combining people, procedures and technology |
| Feature | A recognizable grouping of specified product behavior |
| System / Component | The subject or structural part being specified |

Create a capability with `dotspec capability init local-task-storage`. The
`*.capability.spec.md` file starts with `type: capability-spec`, an optional
frontmatter `id`, and an H1 naming the ability. Use level-two headings for
Description, Purpose, Supports, Depends on, Provided by and Enablers as needed.
The nearest enclosing system or component owns the file. Record providers
explicitly, including when a different component provides the ability.
See [Local task storage](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).

- `supports` links to goals, jobs, steps, outcomes or use cases; it expresses intended contribution.
- `depends_on` links to other capabilities; dependencies can be many-to-many and cyclic.
- `provided_by` links to explicitly identified systems or components; it is optional and non-exhaustive.
- `enablers` describes human, procedural or other enabling context and external references.

A dependency is not a structural `part_of` relationship, and a provider link is
not verified implementation. Capabilities can span systems. Human and procedural
providers do not need to be falsely represented as software components.

Job, JobStep and DesiredOutcome remain in the core
`user-requirements` subset. Capability participates in the core
`business-requirements` and `software-requirements` subsets; these memberships
organize model definitions, not the strategic importance of each capability.

## Generate a needs-to-specification view

```sh
dotspec model needs
```

The generated Markdown includes a Mermaid diagram, linked relationship table,
full statements and evidence. Short labels keep the graph readable; the
statements retain their meaning. Unresolved references remain in the authored
narrative and inspection findings rather than becoming invented edges.

From the repository root, run `dotspec model needs --project todo`. Find **Less reconstruction** and follow its applicability to **Recall outstanding work** and **Check what remains**. Follow the contributions in the other direction. Both **Local task storage** and **Preserve tasks across restart** support this outcome, and **Task service** provides the capability. The linked [outcome statement](../../../examples/todo/specification/jobs/see-personal-commitments-through.job.spec-report.md#less-reconstruction) tells you what a person would have to experience for that support to matter.
