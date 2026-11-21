# Map needs to capabilities

Use this guide to trace which abilities support a job and which parts of the system provide them. Start with the [job and its desired outcomes](write-a-job.md).

## Connect needs to capabilities and specifications

Richard Allen's [User Needs Mapping](https://leanpub.com/userneedsmapping) connects
user needs to enabling capabilities and their dependencies, then considers
technical systems and organizational design. dotspec adopts the needs and
capability relationships while leaving team design and Wardley evolution in
external artifacts. The distinction is informed by *User Needs Mapping*,
printed pp. 39–49; the mapping below is a dotspec design choice.

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

In the [To-do job](../../../examples/todo/specification/jobs/resume-interrupted-work.job.spec-report.md#less-reconstruction),
Less reconstruction links to the restart requirement that supports it and the job
step where it applies. To inspect contributions together in the diagram and table,
run `dotspec model needs --project todo` from the dotspec repository.
