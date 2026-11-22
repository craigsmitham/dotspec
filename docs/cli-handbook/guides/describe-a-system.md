# Describe a system and its boundary

In the directory that will own the system, start a declaration with `dotspec system init "Personal task list"`, substituting your system’s name. If that directory already declares your system, edit its existing `*.system.spec.md` file.

Under **Purpose**, name the people the system serves and the benefit they seek. Under **Scope**, identify the work and situations it covers. Put a plausible but deliberately excluded case under **Out of scope**. If a promise depends on the surroundings, record the relevant condition as a named entry under **Assumptions**. [Worked system description](../../../examples/todo/specs/personal-task-list.system.spec.md).

Ask someone to try two nearby cases. Can they tell from the text why reopening the same database calls for the saved tasks, and why losing the database raises a different question? If they have to guess, name the condition or exclusion that settles it. [Systems and their boundaries](../explanations/systems.md).
