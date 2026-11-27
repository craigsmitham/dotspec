# Map needs to capabilities

## Connect needs to capabilities and specifications

Start with one performer and the need you want to investigate. A job can anchor a map of the whole service. For a narrower question, choose the desired outcome that identifies the difficulty. In To-do, **Less reconstruction** means spending less time identifying the next outstanding action after an interruption. If the need has yet to be recorded, capture it in the job so the map can link back to the performer's full criterion. This follows Richard Allen’s practice of tracing what a particular user depends on. [Write a job](write-a-job.md) · [Worked desired outcome](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#less-reconstruction) · [Allen, *User Needs Mapping* (1st ed., 2025), ch. 4, steps 3–5](https://leanpub.com/userneedsmapping).

Find the first ability the performer depends on. To-do's **Local task storage** retains accepted changes between sessions. **Task service** is the component that provides that ability. In the capability file, link to the desired outcome under **Supports** and to the known provider under **Provided by**. If the ability lacks a specification, create it in the owning system or component directory. [Specify a capability](specify-a-capability.md) · [Worked capability](../../../examples/todo/specs/capabilities/local-task-storage.capability.spec.md).

Trace what one ability needs from another. In Allen’s moviegoer walkthrough, a viewer relies on a movie listing to find a film, and the listing in turn draws on recommendations. Try the same question for **Notice a stalled errand**. When the owner reviews the list on Wednesday, **Follow-up review** needs the expected reply and follow-up day saved on Monday. Link to **Local task storage** under **Depends on** in Follow-up review. Task service provides both abilities. Naming the provider alone would leave the reason for the dependency unexplained. [Allen, “Moviegoer Experience,” step 4, “Discover a Movie”](https://www.userneedsmapping.com/moviegoer-walkthrough) · [To-do dependency](../../../examples/todo/specs/capabilities/follow-up-review.capability.spec.md#depends-on).

Read the chain from the performer's perspective, as Allen does for the moviegoer. What must someone still do when the provider supplies its part? In dotspec, describe that person's contribution under **Enablers**. The To-do owner chooses when to check back and consults their messages. If they first open Waiting on Friday, the service flags Wednesday's follow-up as due. They may already have chased the plumber by phone on Wednesday. To assess how long the errand went unattended, ask when the reply was expected and when the owner noticed it was late. [Follow-up review enablers](../../../examples/todo/specs/capabilities/follow-up-review.capability.spec.md#enablers).

Keep each relative link in the file of the subject making the claim. For example, if a requirement contributes to the desired outcome, put the **Supports** link in the requirement. If a desired outcome applies to two job steps, link to both under **Applies to** in the outcome. [To-do restart requirement](../../../examples/todo/specs/requirements/restart.requirement.spec.md#supports) · [Where the outcome applies](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#applies-to).

## Read the map in the generated specification

With your job and capability saved, run `dotspec build` from your project root and open the job’s report in `specification/jobs/`. Each desired outcome lists what contributes to it under **Supported by**, and each step lists the outcomes that apply to it under **Applies here**. Follow the names and read each link as a sentence. To read the To-do map offline, run `dotspec docs examples/todo/specification/jobs/see-personal-commitments-through.job.spec-report.md#less-reconstruction`.

- **Where does the criterion apply?** Under **Applies to**, Less reconstruction applies to both *Recall outstanding work* and *Check what remains*.
- **What contributes to it?** Under **Supported by**, Local task storage and the requirement *Preserve tasks across restart* each support the same outcome.
- **What provides the ability?** The Local task storage report names Task service under **Provided by**.
- **Which ability depends on it?** The Follow-up review report names Local task storage under **Depends on** and *Notice a stalled errand* under **Supports**.

If an expected link is missing, follow **View source** from the report that should show it. A link in **Description** can lead to the right file without recording its role. Put the link under its relationship heading, check the path and fragment, then rebuild and run `dotspec check`. If you chose a system goal under **Supports**, review that link in the capability’s source.

Review your chain with the performer. Even when the To-do links resolve, an entry titled “Plumber” may give the person too little to act on. To judge whether storage helped, observe what they recorded and how they found their next action after the break. [To-do participant exercise](../../../examples/todo/README.md#try-the-workflow-with-a-participant).
