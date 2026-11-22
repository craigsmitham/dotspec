# Specify a component

Choose the system or component that contains this part. From your project root, supply the path to that existing declaration and an existing directory for the new part. For example, with a system declared at `specs/personal-task-list.system.spec.md` and an `app/` directory available for its service, use the following command. Substitute your component name and paths.

```sh
dotspec component init "Task service" --path app --part-of specs/personal-task-list.system.spec.md
```

Keep the **Part of** link the starter creates. Under **Responsibilities**, say what another part can rely on this one to do. Under **Boundary**, explain what it includes or delegates when the division needs clarification. Link to the authoritative contract under **Contracts**. If the timing or result of an operation matters, state the promise in a requirement and link back to the component under **Applies to**. [Task service](../../../examples/todo/app/specs/task-service.component.spec.md) · [Requirement governing confirmation](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md#applies-to).

Ask a reader where they would look to resolve a disagreement between the caller and this part. They should be able to name the responsible part and find the contract or requirement governing the exchange. [Components and their responsibilities](../explanations/components.md).
