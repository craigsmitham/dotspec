# Describe a system and its boundary

Choose a situation the system is meant to handle and one close enough that a stakeholder might expect it too. Find out who depends on the result and which part of it the team has agreed to deliver.

In the directory that will own the system, edit its existing `*.system.spec.md` file. For a new system, run `dotspec system init "Personal task list"`, substituting your system’s name. The name in the top-level heading determines the filename.

Use **Description** to introduce the system in a sentence and **Purpose** to explain why it exists. For the personal task list, the benefit is being able to pick up unfinished work after an interruption. If a reader new to the system needs to see how its parts fit together before the details make sense, explain that under **Overview**. The task list’s Overview introduces its Open, Waiting and Completed views through a repair the owner still owes after sending the plumber a photo. An Overview explains without committing. When the team must be able to check a statement, write it as a requirement. [Personal task list overview (dotspec checkout)](../../../examples/todo/specs/personal-task-list.system.spec.md#overview)

Then choose the heading that fits the claim you want to make.

| Heading | Question to answer |
| --- | --- |
| **Scope** | What work can people expect this system to handle, for whom and in which situations? |
| **Out of scope** | What might someone reasonably expect that the team has agreed to exclude? |
| **Environment** | What conditions in the surrounding world help explain a requirement? |
| **Assumptions** | What premise are we relying on, and which promise would we revisit if it failed? |
| **Constraints** | What choice has already been limited, and who or what imposes that limit? |

For the personal task list, one person keeping commitments on one computer sets the scope. Sending reminders is an agreed exclusion. When a team still needs to decide whether to cover a case, list the question under **Open questions** with the person or evidence that can settle it. [Personal task list boundary (dotspec checkout)](../../../examples/todo/specs/personal-task-list.system.spec.md#scope)

Under **Stakeholders**, give each role a level-three heading in the project’s own terms. Describe what people in that role depend on or decide. To find perspectives you may have missed, ask who operates, funds, maintains or regulates the system, who benefits from its results, and who else experiences its effects. One person can fill several roles. In To-do, *List owner* identifies who decides what to discard. Link to that heading from a job’s **Performer** or a use case’s **Primary actor** or **Stakeholder interests** when it identifies the same role. [Meyer, §3.4, G.7](https://se.inf.ethz.ch/requirements/plan.pdf#page=5) · [List owner](../../../examples/todo/specs/personal-task-list.system.spec.md#list-owner) · [Resume list primary actor](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md#primary-actor).

Under **Goals**, give each benefit stakeholders expect its own level-three heading and say what should improve. In Meyer’s Standard Plan, expected benefits are the core of the Goals book and help keep a project focused. A job that serves the goal links to it under a level-two **Supports** heading, and its outcomes say what to observe. In To-do, the commitments job supports *Resume unfinished work*, and its outcome *Less reconstruction* concerns the time to find the next action after an interruption. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §3.4, G.3](https://se.inf.ethz.ch/requirements/plan.pdf) · [Personal task list goals](../../../examples/todo/specs/personal-task-list.system.spec.md#goals) · [Less reconstruction](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#less-reconstruction).

Meyer’s Environment book covers the application domain as well as the physical and technical setting. In the task list, the owner can stop the service and return on a different local day. Record that circumstance under **Environment**. What the owner should see on the new day belongs in the follow-up requirement. [Meyer, *Handbook of Requirements and Business Analysis* (2022), §3.5](https://se.inf.ethz.ch/requirements/plan.pdf) · [Task list environment (dotspec checkout)](../../../examples/todo/specs/personal-task-list.system.spec.md#environment) · [Follow-up requirement (dotspec checkout)](../../../examples/todo/specs/requirements/follow-up.requirement.spec.md#statement)

Under **Assumptions**, give each premise a level-three heading. The restart requirement relies on the owner returning to the same writable database. Link to that named assumption from the requirement’s **Relies on** section. Under the assumption, level-four **Basis** and **Applicability** headings explain why it is reasonable to rely on the premise and where it holds. From the project root, `dotspec model pegs` flags an assumption that lacks either. [Same local storage (dotspec checkout)](../../../examples/todo/specs/personal-task-list.system.spec.md#same-local-storage) · [Restart requirement (dotspec checkout)](../../../examples/todo/specs/requirements/restart.requirement.spec.md#relies-on)

For **Constraints**, identify the source of the restriction. Suppose a customer’s hosting agreement requires all task data to stay on the owner’s computer. That term limits which storage solutions the team can choose. Cite the agreement under **Constraints**. If a requirement states the storage obligation and its acceptance criterion, link to that requirement there.

From the system directory, run `dotspec system show` and check the displayed purpose, scope and source path. Read the exclusions and assumptions in that file. Can a stakeholder explain what to expect in both of your situations? If their answer turns on an unstated condition, add it to the declaration and try the cases again. [Systems and their boundaries](../explanations/systems.md).
