# To-do

A personal task list you can close today and pick up tomorrow. One local Go process serves HTML forms and keeps the list in SQLite.

## Run the app

From a fresh checkout, install the project's pinned tools using mise. Run the following from the repository root.

```sh
mise trust
mise install
cd examples/todo/app
mise exec -- go run .
```

Open the loopback URL printed in the terminal. Add “Call the plumber” and “Collect keys”. Rename the first task, mark it complete and look in **Completed**. Stop the process with Ctrl+C and run the same command again in the same directory. Both lists should look as you left them. Reopen the finished errand and find it in **Open**. Expand **Edit task** to rename or delete an entry.

The default file is `todo.sqlite` in the app's working directory. Use `-db PATH` to choose a file or `-port 0` to choose a free port. The first build downloads the Go dependencies. The app runs locally after that.

## Follow one decision

Start with the [job of resuming interrupted work](specs/jobs/resume.job.spec.md). The [business rule](specs/business-rules/completion.business-rule.spec.md) explains why finishing an errand keeps it available. The [return scenario](specs/use-cases/resume-list.use-case.spec.md) connects that policy to the [checklist actions](specs/concepts/checklist.concept.spec.md).

The [restart requirement](specs/requirements/restart.requirement.spec.md) refines the scenario's promise and classifies it against the system's [durability concern](specs/quality-models/personal-data.quality-model.spec.md#durability). Follow its verification link to the [acceptance test](tests/acceptance.py). Compare the fields it reads through the [forms](app/list.html) with the updates the [server](app/main.go) commits before redirecting.

## Verify the example

From the repository root, run the acceptance test alone, inspect the authored relationships, or run the whole example gate.

```sh
mise exec -- python -I examples/todo/tests/acceptance.py
mise exec -- go run . model trace --path examples/todo
mise run example-check
```

The acceptance test builds the app in a temporary directory, submits the rendered forms and restarts the server on a temporary database. For a session with a participant, use the [validation exercise](specs/system.spec.md#validation-plan). The [coverage inventory](COVERAGE.md) identifies what the gate exercises and which fields still need a worked example.
