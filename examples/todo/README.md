# To-do

Suppose you promised a guest a working sink by Friday. On Monday you ask the plumber for a visit time and decide to chase the reply on Wednesday. What should the list tell you on Wednesday? If the plumber asks for a photo, whose turn is it to act? This example follows that decision from a personal commitment to a working browser form and a restart test.

## Run the app

From the root of the dotspec checkout, install the pinned tools and start the local Go service. The first build needs access to the Go dependencies. Use a temporary directory so you can repeat the exercise with an empty list.

```sh
mise trust
mise install
cd examples/todo/app
todo_dir=$(mktemp -d)
mise exec -- go run . -port 0 -db "$todo_dir/list.sqlite"
```

Open the address printed in the terminal. The person keeping the list is the only one who uses it. Calls and messages in the walkthrough are fictional.

1. In **Open**, add **Fix the guest sink** and **Collect keys**. You have called the plumber. Under the sink entry, expand **Wait for reply**. Enter **Plumber’s visit time** and choose today as the **Follow-up date**. Submit the form. Open now contains just the keys.
2. Open **Waiting**. The sink entry shows your expected reply, today’s date and **Follow up now**. Suppose you chase the plumber and agree to check back tomorrow. Expand **Update follow-up**, change just the date to tomorrow and submit. The same entry now reads **Upcoming**. That day is when to chase the plumber. Friday is still when you owe your guest a working sink.
3. Stop the service with Ctrl+C. In the same shell, run the last command again and open the new address. Check Waiting for the visit time and the day you agreed. If the list is empty, check that you restarted with the same database path. With the default invocation the service uses `todo.sqlite` in its working directory.
4. Try two replies. “Send a photo” means you have work to do. Choose **Resume** and look in Open for the sink. **Wait for reply** now starts with empty fields. After sending the photo, put the errand back in Waiting for **Confirmation the repair is finished**, using yesterday as the follow-up date. It immediately reads **Follow up now** and stays in Waiting.
5. This time the plumber says the repair is done and you confirm the sink works. Choose **Complete** and open **Completed**. If the leak recurs, **Reopen** returns that entry to Open. Under **Edit task**, **Delete task** removes the entry entirely.

## Follow one decision

Start at [SPEC.md](SPEC.md) and open the system specification. Under **Goals**, open **Resume unfinished work**. Under **Supported by**, choose the job **See personal commitments through**. In the job’s **Evidence** section, follow the **return scenario**. At step 2, open **Show which replies need follow-up**. Its **Derived from** link takes you to the rule that explains why an unanswered errand is still yours after the day passes. Its **Evidence** link takes you to the test that reads the cue on the rendered form.

For the distinction between facts the owner saved and a cue that can change overnight, start with the restart requirement. The feature groups both promises for the returning owner. [Pick up where I left off](specs/features/pick-up-where-i-left-off.feature.spec.md) · [Preserve tasks across restart](specs/requirements/restart.requirement.spec.md).

## Verify the example

From the root of the dotspec checkout, run the acceptance test by itself or run the full example gate.

```sh
mise exec -- python -I examples/todo/tests/acceptance.py
mise run example-check
```

The acceptance test builds the app, submits its HTML forms and starts a second process on the same disposable database. It checks the entry identities and visible details in each view, including the saved reply and day. The [coverage inventory (dotspec checkout)](COVERAGE.md) gives cases to try against each authored specification and identifies concerns left for a larger system.

To follow the specifications with an installed CLI, run `dotspec docs job` from any directory and open the worked job listed there. To inspect your edits in the checkout, run `mise run generate` from its root, then `mise exec -- go run . model trace --path examples/todo` to inspect the relationships.

## Try the workflow with a participant

Use a fresh database. Invite someone to recall a recent errand where they needed an answer from another person. What had they promised, what did they do while waiting, and how did they tell it was finished? Let them record an example they are comfortable sharing, take a break, and then return to the list. Ask what needs attention today and what they would do if a reply requested more work. Let them choose the controls.

With their consent, record the cue they used, the action they took and anything they had to reconstruct elsewhere. The job and its outcomes are hypotheses. A passing restart test establishes what the software preserved. Whether those details help someone resume a real commitment is still open. [Less reconstruction](specs/jobs/personal-commitments.job.spec.md#less-reconstruction).
