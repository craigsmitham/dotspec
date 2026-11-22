# Specify a component

Identify the responsibility a caller relies on and the part that can discharge it. Choose its parent system or component in the same project. If the caller uses an independently operated system outside that boundary, describe the service the caller requires under **Required interfaces** in the caller’s specification. [Components and their responsibilities](../explanations/components.md).

For a new part, start in your project root with the parent specification and a directory for the part. Create that directory if needed. In the example below the parent is `specs/personal-task-list.system.spec.md` and `app/` has no declaration of its own. Substitute your name and paths. If the part already has a declaration, edit that file.

```sh
dotspec component init "Task service" --path app --part-of specs/personal-task-list.system.spec.md
```

The starter writes **Part of** as a relative link from the new file to its parent. Under **Description** and **Responsibilities**, identify the part and say what callers depend on it to do. Where another part takes over, explain the handoff under **Boundary**. Under **Contracts**, link to an agreed contract if one exists.

When a caller relies on two operations happening in order, ask what must already be true for the second operation to proceed. In Meyer's insurance claim example, the adjuster can negotiate once the insurer has evaluated the damage and set a reserve. In To-do, the browser can report a successful save once the task service has committed the change. Put that condition in a requirement for the part that controls it. Link to the requirement from the component’s responsibilities and name the component under the requirement’s **Applies to**. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §8.7.1, pp. 147–149](https://se.inf.ethz.ch/requirements/logical.pdf) · [Task service](../../../examples/todo/app/specs/task-service.component.spec.md#responsibilities) · [Requirement governing confirmation](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md#applies-to).

Try a case where the prerequisite fails partway through the caller’s action. Completing a waiting errand both settles the commitment and clears the expected reply. If storage accepts the first write but rejects the next, what would the owner see? Under **Evidence** in the component requirement, identify how to check the response and the resulting entry. A successful restart exercises the case where the service could save everything. The failed completion exercises the point where a success response would mislead the caller. [Task service evidence](../../../examples/todo/app/specs/requirements/acknowledge-after-commit.requirement.spec.md#evidence).

Run `dotspec model overview` from the project root. Find the part and its **Explicit parent** in the system structure table. If the parent is wrong, correct the **Part of** link in the component file and regenerate the view. Follow the component’s specification and try a disputed case with a caller. Can they identify who owes the response and find the contract or requirement that settles it?
