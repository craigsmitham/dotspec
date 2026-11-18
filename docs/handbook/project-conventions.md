# Project conventions

Give an agent a way to find the decisions your team has already made. Add a
short paragraph to the project's existing contributor or agent instructions.
Use paths to actual files and say what authority each carries. For shared
authoring methods, point to `dotspec guide`.

For example, the document recovery project could use this entry in its
`AGENTS.md`:

```markdown
## Specifications

Start at [the workspace boundary](specs/system.spec.md). The
[recovery policy](specs/business-rules/retention.business-rule.spec.md) owns the
30-day window. The [restore use case](specs/use-cases/restore.use-case.spec.md#restore-content)
owns the software promise. Link to those statements when introducing a refinement.
The clock API contract is in `contracts/time.openapi.yaml`.

Use `dotspec guide` to choose an authoring form. Run `dotspec check` in this
project after editing. Record an unresolved target with the evidence needed
to choose it.
```

Adapt this to your team's decisions about policies, executable specifications,
approvals and verification evidence. Start with one representative change and
check whether another contributor can find its governing source.
