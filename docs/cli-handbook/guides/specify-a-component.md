# Specify a component

Identify the responsibility a caller relies on and the part that can discharge it. Choose its parent system or component in the same project. If the caller uses an independently operated system outside that boundary, describe the service the caller requires under **Required interfaces** in the caller’s specification. [Components and their responsibilities](../explanations/components.md).

For a new part, start in your project root with the parent specification and a directory for the part. Create that directory if needed. In the example below the parent is `specs/personal-task-list.system.spec.md` and `app/` has no declaration of its own. Substitute your name and paths. If the part already has a declaration, edit that file.

```sh
dotspec component init "Task service" --path app --part-of specs/personal-task-list.system.spec.md
```

The starter writes **Part of** as a relative link from the new file to its parent. Under **Description** and **Responsibilities**, identify the part and say what callers depend on it to do. Where another part takes over, explain the handoff under **Boundary**. Under **Contracts**, link to an agreed contract if one exists. For a promise with conditions or an ordering that matters, link to its requirement and identify the component under that requirement’s **Applies to**. For example, To-do’s browser can acknowledge a successful save to the person because Task service must commit the change before sending its reply. [Task service](../../../examples/todo/app/specs/task-service.component.spec.md#responsibilities) · [Requirement governing confirmation](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md#applies-to).

Run `dotspec model overview` from the project root. Find the part and its **Explicit parent** in the system structure table. If the parent is wrong, correct the **Part of** link in the component file and regenerate the view. Follow the component’s specification and try a disputed case with a caller. Can they identify who owes the response and find the contract or requirement that settles it?
