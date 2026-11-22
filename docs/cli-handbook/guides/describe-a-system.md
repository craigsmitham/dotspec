# Describe a system and its boundary

Choose a situation the system is meant to handle and one close enough that a stakeholder might expect it too. Find out who depends on the result and which part of it the team has agreed to deliver.

In the directory that will own the system, edit its existing `*.system.spec.md` file. For a new system, run `dotspec system init "Personal task list"`, substituting your system’s name. The name in the top-level heading determines the filename.

Fill in **Description**, **Purpose** and **Scope**. Say what the system is, whom it helps and which work it covers. For the personal task list, “pick up unfinished work” explains the benefit. “One person’s open and completed tasks on one computer” bounds the work. If the team has agreed to exclude a plausible case, put that decision under **Out of scope**. If the team still needs to decide, record the question under **Framing**. [Personal task list](../../../examples/todo/specs/personal-task-list.system.spec.md#scope).

For conditions outside the team’s control, say how they affect the boundary. The To-do service must accommodate a restart, but its recovery promise depends on the original database remaining available. Record a condition the system counts on as a named level-three heading under **Assumptions**. Meyer distinguishes the system decisions a team makes from conditions in its environment. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §§1.1.2 and 1.7.3](https://doi.org/10.1007/978-3-031-06739-6) · [Same local storage](../../../examples/todo/specs/personal-task-list.system.spec.md#same-local-storage).

From the system directory, run `dotspec system show` and check the displayed purpose, scope and source path. Read the exclusions and assumptions in that file. Can a stakeholder explain what to expect in both of your situations? If their answer turns on an unstated condition, add it to the declaration and try the cases again. [Systems and their boundaries](../explanations/systems.md).
