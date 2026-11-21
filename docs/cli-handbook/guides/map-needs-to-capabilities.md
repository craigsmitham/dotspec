# Map needs to capabilities

## Connect needs to capabilities and specifications

Choose the need that explains why an ability matters. A job can anchor a map of the whole service. To investigate a delay, choose the desired outcome that identifies where the performer loses time. In To-do, **Less reconstruction** means spending less time identifying the next outstanding action after an interruption. Keep that full criterion in the job and link to it from the capability. [Desired outcome](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#less-reconstruction) · [How User Needs Mapping uses a job or outcome](../explanations/jobs-to-be-done.md#from-a-job-map-to-a-user-needs-map).

Name the capability as an ability. **Local task storage** retains accepted changes between sessions. **Task service** is the component that provides that ability. A feature can group the product behavior a person recognizes, such as **Pick up where I left off**. [Worked capability](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md) · [Feature in the To-do specification](../../../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md).

From the owning system or component directory, substitute your ability’s name for `local-task-storage`. Use `dotspec capability init local-task-storage` to start a `*.capability.spec.md` file under its owning system or component. The starter sets `type: capability-spec` in frontmatter. Name the ability in the H1 and describe it under `Description`. Choose level-two headings for the relationships you can explain.

| Heading | What belongs there |
| --- | --- |
| `Supports` | Link to the goal, job, step, desired outcome or use case the ability helps. |
| `Depends on` | Link to other capabilities this ability needs. Several capabilities can depend on each other. |
| `Provided by` | Link to the systems or components known to provide the ability. A provider can differ from the file's owner. |
| `Enablers` | Describe people, procedures or other conditions the ability depends on. Cite external material where helpful. |

## Generate a needs-to-specification view

Run `dotspec model needs` in your project, or `dotspec model needs --project todo` from the dotspec repository root. In the relationship table, check that each link answers the question its label asks. The To-do example has three routes worth checking.

- **Where does the criterion apply?** Less reconstruction applies to both *Recall outstanding work* and *Check what remains*. [Authored applicability](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#applies-to).
- **What contributes to it?** Local task storage and the requirement *Preserve tasks across restart* each support the same outcome. [Capability’s contribution](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md#supports) · [Requirement’s contribution](../../../examples/todo/specs/requirements/restart.requirement.spec.md#supports).
- **What provides the ability?** Local task storage names Task service as its provider. [Provider declaration](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md#provided-by).

If an expected link is missing, inspect its heading and target in the source. Even when all three routes resolve, an entry titled “Plumber” may give the person too little to act on. To judge whether storage helped, observe what they recorded and how they found their next action after the break. [To-do participant exercise](../../../examples/todo/README.md#try-the-workflow-with-a-participant).
