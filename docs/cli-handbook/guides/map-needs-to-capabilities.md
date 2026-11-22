# Map needs to capabilities

## Connect needs to capabilities and specifications

Start with one performer and the need you want to investigate. A job can anchor a map of the whole service. For a narrower question, choose the desired outcome that identifies the difficulty. In To-do, **Less reconstruction** means spending less time identifying the next outstanding action after an interruption. If the need has yet to be recorded, capture it in the job so the map can link back to the performer's full criterion. This follows Richard Allen’s practice of tracing what a particular user depends on. [Write a job](write-a-job.md) · [Worked desired outcome](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#less-reconstruction) · [Allen, *User Needs Mapping* (1st ed., 2025), ch. 4, steps 3–5](https://www.userneedsmapping.com/moviegoer-walkthrough).

Find the first ability the performer depends on. To-do's **Local task storage** retains accepted changes between sessions. **Task service** is the component that provides that ability. In the capability file, link to the desired outcome under **Supports** and to the known provider under **Provided by**. If the ability lacks a specification, create it in the owning system or component directory. [Specify a capability](specify-a-capability.md) · [Worked capability](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).

Follow the dependency as far as it helps answer your question. Where one ability relies on another, link from the first capability to the second under **Depends on**. If a human practice helps make an ability possible, describe it under **Enablers**. Keep each relative link in the file of the subject making the claim. For example, if a requirement contributes to the desired outcome, put the **Supports** link in the requirement. If a desired outcome applies to two job steps, link to both under **Applies to** in the outcome. [To-do restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md#supports) · [Where the outcome applies](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#applies-to).

## Generate a needs-to-specification view

With your job and capability saved, run `dotspec model needs` from your project root. To reproduce the To-do view, run `dotspec model needs --project todo` from the dotspec repository root. In **Relationships**, follow the names and read each row as a sentence. The example gives three routes to check.

- **Where does the criterion apply?** Less reconstruction applies to both *Recall outstanding work* and *Check what remains*.
- **What contributes to it?** Local task storage and the requirement *Preserve tasks across restart* each support the same outcome.
- **What provides the ability?** Local task storage names Task service as its provider.

If an expected row is missing, open the source under **Definitions and evidence**. A link in **Description** can lead to the right file without recording its role. Put the link under its relationship heading, check the path and fragment, and rerun the view. If you chose a system goal under **Supports**, review that link in the capability’s source. The needs view shows job, step, outcome and use-case targets in its relationship table.

Review your chain with the performer. Even when all three To-do routes resolve, an entry titled “Plumber” may give the person too little to act on. To judge whether storage helped, observe what they recorded and how they found their next action after the break. [To-do participant exercise](../../../examples/todo/README.md#try-the-workflow-with-a-participant).
