# To-do

A personal task list you can close today and pick up tomorrow. One local Go process serves HTML forms and keeps the list in SQLite.

From any directory, use `dotspec docs job` to find the worked job’s reading command. For example, `dotspec docs examples/todo/specs/jobs/personal-commitments.job.spec.md#job` opens the job statement. Running the app, editing its specifications and running its tests require the dotspec checkout.

## Run the app

From a fresh dotspec checkout, install the project's pinned tools using mise. Run the following from the repository root.

```sh
mise trust
mise install
cd examples/todo/app
mise exec -- go run .
```

Open the loopback URL printed in the terminal. Add “Call the plumber” and “Collect keys”. Rename the first task, mark it complete and look in **Completed**. Stop the process with Ctrl+C and run the same command again in the same directory. Both lists should look as you left them. Reopen the finished errand and find it in **Open**. Expand **Edit task** to rename or delete an entry.

The default file is `todo.sqlite` in the app's working directory. Use `-db PATH` to choose a file or `-port 0` to choose a free port. The first build downloads the Go dependencies. The app runs locally after that.

## Follow one decision

Start with the [job of seeing personal commitments through](specs/jobs/personal-commitments.job.spec.md). The [business rule](specs/business-rules/completion.business-rule.spec.md) explains why finishing an errand keeps it available. The [return scenario](specs/use-cases/resume-list.use-case.spec.md) connects that policy to the [checklist actions](specs/concepts/checklist.concept.spec.md).

The [Task service specification](app/specs/task-service.component.spec.md) links back to the system and links to the [acknowledgement requirement](app/specs/requirements/acknowledge-after-commit.requirement.spec.md). Its [component report](specification/components/task-service/task-service.component.spec-report.md) brings those sources together.

The [restart requirement](specs/requirements/restart.requirement.spec.md) refines [Keep recorded work available](specs/requirements/availability.requirement.spec.md) for a service interruption and is classified under the system's [durability concern](specs/quality/durability/durability.quality-attribute.spec.md). Read the [acceptance test](tests/acceptance.py) to see what it compares. For implementation detail, inspect the [HTML forms (dotspec checkout)](app/list.html) and the [server’s update handler (dotspec checkout)](app/main.go).

## Read the requirements specification

Start with the [system overview](specification/personal-task-list.system.spec-report.md) or its [requirements index](specification/requirements.md). The pages are generated from the [To-do system](specs/personal-task-list.system.spec.md) and its related specifications. To rebuild after an edit, run `dotspec spec build` and `dotspec check` in your project. For this example, run those commands in `examples/todo/` in the dotspec checkout.

| Question | Start here |
| --- | --- |
| Which features, interfaces and behaviors does the system expose? | [System overview](specification/personal-task-list.system.spec-report.md) |
| How do the user's goals connect to the specifications? | [User job](specification/jobs/see-personal-commitments-through.job.spec-report.md) |
| What changes when someone reopens a task? | [Checklist](specification/concepts/checklist.concept.spec-report.md) |
| Which promise does the restart requirement refine? | [Restart requirement](specification/requirements/preserve-tasks-across-restart.requirement.spec-report.md) |
| Which scenarios and decisions need review? | Run `dotspec model review` in the checkout’s `examples/todo/` directory |
| How do project, environment, goals and system fit together? | [Purpose, environment and scope](specification/personal-task-list.system.spec-report.md) |

## Verify the example

From the root of the dotspec checkout, run the acceptance test alone, inspect the authored relationships, or run the whole example gate.

```sh
mise exec -- python -I examples/todo/tests/acceptance.py
mise exec -- go run . model trace --path examples/todo
mise run example-check
```

The acceptance test builds the app in a temporary directory, submits the rendered forms and restarts the server on a temporary database. The [example coverage inventory (dotspec checkout)](COVERAGE.md) shows which fields the example exercises.

## Try the workflow with a participant

With the To-do app running from the dotspec checkout, invite someone to record two real errands and mark one finished. After a break,
restart the app on the same database. Ask which errand still needs attention,
then ask them to resume the finished errand for a follow-up. Let them choose the
controls. With their consent, record what they recognized and where they needed
help. Compare those observations with the [user’s job](specs/jobs/personal-commitments.job.spec.md).
