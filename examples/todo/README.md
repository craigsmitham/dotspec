# To-do

A personal task list you can close today and pick up tomorrow. One local Go process serves HTML forms and keeps the list in SQLite.

With the dotspec CLI on your PATH, use `dotspec docs job` from any directory to find the worked job’s reading command. For example, `dotspec docs examples/todo/specs/jobs/personal-commitments.job.spec.md#job` opens the job statement. Running the app, editing its specifications and running its tests require the dotspec checkout.

## Run the app

Start in the dotspec checkout with mise installed. The first build needs access to the Go dependencies. From the repository root, install the pinned tools and start the app.

```sh
mise trust
mise install
cd examples/todo/app
mise exec -- go run .
```

Open the loopback URL printed in the terminal. Add “Call the plumber” and “Collect keys”. Rename the first task, mark it complete and look in **Completed**. Stop the process with Ctrl+C and run the same command again in the same directory. Both lists should look as you left them. Reopen the finished errand and find it in **Open**. Expand **Edit task** to rename or delete an entry.

The default database is `todo.sqlite` in the app’s working directory. If port 8080 is busy, run `mise exec -- go run . -port 0` from the same directory and open the new address it prints. For a separate session, use `-db PATH` to choose a file. If both lists are unexpectedly empty after a restart, check the working directory and database path before adding a task. Stop the app and restart it with the original file.

## Follow one decision

Start with the [job of seeing personal commitments through](specs/jobs/personal-commitments.job.spec.md). The [business rule](specs/business-rules/completion.business-rule.spec.md) explains why finishing an errand keeps it available. The [return scenario](specs/use-cases/resume-list.use-case.spec.md) connects that policy to the [checklist actions](specs/concepts/checklist.concept.spec.md).

The [Task service specification](app/specs/task-service.component.spec.md) links back to the system and links to the [acknowledgement requirement](app/specs/requirements/acknowledge-after-commit.requirement.spec.md). Its [component report](specification/components/task-service/task-service.component.spec-report.md) brings those sources together.

The [restart requirement](specs/requirements/restart.requirement.spec.md) refines [Keep recorded work available](specs/requirements/availability.requirement.spec.md) for a service interruption and is classified under the system's [durability concern](specs/quality/durability/durability.quality-attribute.spec.md). Read the [acceptance test](tests/acceptance.py) to see what it compares. For implementation detail, inspect the [HTML forms (dotspec checkout)](app/list.html) and the [server’s update handler (dotspec checkout)](app/main.go).

## Read the requirements specification

Start with the [system overview](specification/personal-task-list.system.spec-report.md) or its [requirements index](specification/requirements.md). The pages are generated from the [To-do system](specs/personal-task-list.system.spec.md) and its related specifications. With the CLI on your PATH, run `dotspec spec build` and `dotspec check` in your project after editing an authored source. For To-do, use the checkout’s `examples/todo/` directory. If you have only the checkout, run `mise run build` from its root and use `../../dist/dotspec` in place of `dotspec` from `examples/todo/`.

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

The acceptance test builds the app in a temporary directory, submits the rendered forms and restarts the server on a temporary database. It compares entry identities and titles in both views. If it reports a changed list, compare the last acknowledged entries with the ones served after restart. The [example coverage inventory (dotspec checkout)](COVERAGE.md) shows which fields the example exercises.

## Try the workflow with a participant

Use a separate database for the session. Start the app from the checkout’s `examples/todo/app/` directory, adding `-db session.sqlite` with an unused filename to the run command above. Invite the participant to record two errands they are comfortable sharing and mark one finished. After a break, restart on the same file and open the address printed by the new process.

Ask which errand still needs attention and how they know. Give the finished errand a follow-up and ask how they would resume it. Let them choose the controls. With their consent, record the cue they recognized, the action they took and any help they needed. If the entry is present but they cannot tell what to do next, capture the words that left them guessing. Compare the observation with the [desired outcome for returning after an interruption](specs/jobs/personal-commitments.job.spec.md#less-reconstruction). The acceptance test can check persistence with its own saved entries. To investigate a missing errand from this session, check which database each process opened and what the participant saw after submitting it.
