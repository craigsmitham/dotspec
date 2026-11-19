# Project conventions

Give an agent a way to find the decisions your team has already made. Add a
short paragraph to the project's existing contributor or agent instructions.
Use paths to actual files and say what authority each carries. For shared
authoring methods, point to `dotspec guide`.

For example, a contributor working on the To-do app could use this entry in its `AGENTS.md`.

```markdown
## Specifications

Start at [the task list boundary](specs/personal-task-list.system.spec.md). The
[completion policy](specs/business-rules/completion.business-rule.spec.md) explains
why finished entries stay available. The [return scenario](specs/use-cases/resume-list.use-case.spec.md#keep-recorded-work-available)
owns the promise. The [restart requirement](specs/requirements/restart.requirement.spec.md)
refines it and links to its executable test.

Use `dotspec guide` to choose an authoring form. Run `dotspec check` in this
project after editing. Record an unresolved target with the evidence needed
to choose it.
```

Adapt this to your team's decisions about policies, executable specifications,
approvals and verification evidence. Start with one representative change and
check whether another contributor can find its governing source.
