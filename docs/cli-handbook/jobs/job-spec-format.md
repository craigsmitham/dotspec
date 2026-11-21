# Job specification format

## Job commands

In your project directory, choose a name for the new job in place of `personal-commitments`. For an existing job, use its filename slug and its own step IDs. To inspect the worked example, run the read commands with `--project todo` from the dotspec checkout. To read its specification from any directory, use `dotspec docs job` and open the worked job.

| Command | Result |
| --- | --- |
| `dotspec job init personal-commitments --preset odi` | Creates an editable starter under `specs/jobs/`. An existing file is preserved and the command reports an error. |
| `dotspec job show personal-commitments` | Displays the authored Markdown, source path and associated system. |
| `dotspec job list` | Lists jobs in the selected subtree, grouped by associated system. Jobs without an enclosing system appear under `System: None`. |
| `dotspec job step list personal-commitments` | Lists the steps in document order with their descriptions, stages and explicit IDs where supplied. |
| `dotspec job outcome list personal-commitments` | Lists desired outcomes grouped by declaration owner. |
| `dotspec job outcome list personal-commitments --step recall` | Selects the outcomes declared under the step with that explicit ID. In the To-do example the selection is empty. An unknown step ID is an error. |
| `dotspec job map personal-commitments --outcomes` | Shows the job and authored step order, with outcomes attached to their declaration owner. Omit `--outcomes` for the steps alone. Use `--format mermaid` for flowchart source. |
| `dotspec model needs` | Shows explicit outcome applicability and contributions from other specifications. |
| `dotspec job check personal-commitments --preset odi` | Checks structure and references and suggests improvements. Use `--json` for structured findings. |

The filename slug uses lowercase ASCII letters and digits with single hyphens between words. Changing an H1 leaves the selector unchanged. Jobs in different directories may share a slug. For `show` and `check`, a source path selects a specific file. The nearest owning subject determines the default system association. [Project layout](../../reference/layout.md#project-layout).

## Job spec format

A job file requires `type: job-spec` in its frontmatter and one nonempty H1. The only other supported frontmatter key is an optional `id`. A minimal draft can consist of the frontmatter and a title.

| Conventional H2 | Content |
| --- | --- |
| `Job` | The accomplishment or change the performer seeks. |
| `Performer` | Who carries out the job. |
| `Situation`, `Motivation`, `Intended progress` | The circumstances, reason to act and change sought. Useful together for a job story. |
| `Steps` | H3 headings naming accomplishments in authored map order. For a partial map, describe the limits under `Evidence`. |
| `Outcomes` | H3 headings naming desired results or success criteria declared by the job. |
| `Evidence` | Source and limits of an account, study links or hypotheses still to investigate. |

These sections are optional, case-sensitive and may appear once each. Add a map when it helps describe the work. Arrows in the map view show document order. Real work can loop or branch. A job without steps can still declare outcomes. Use explicit [relationships](../guides/map-needs-to-capabilities.md) such as `Supports` to connect specifications.

### Steps, stages and desired outcomes

This partial example shows the heading levels for a step-owned outcome and a job-owned outcome that applies to two steps.

```markdown
---
type: job-spec
---
# See personal commitments through

## Job

See personal commitments through from accepting responsibility for an errand to settling it.

## Performer

The person responsible for the errand.

## Steps

### Recall outstanding work

**ID:** recall
**Stage:** locate

Recover the commitments that may need attention.

### Check what remains

**Stage:** monitor

Compare what has happened with what was promised.

#### Outcomes

##### Notice a stalled errand

Minimize the time a commitment goes unattended after an expected reply is overdue.

## Outcomes

### Less reconstruction

**ID:** less-reconstruction

Minimize the time it takes to identify the next outstanding action when returning to an errand after an interruption.

#### Applies to

[Recall outstanding work](#recall-outstanding-work) and [Check what remains](#check-what-remains).

## Evidence

This partial map covers follow-up after an interruption. Both criteria are hypotheses.
```

- **Steps.** Place optional `**ID:**` and `**Stage:**` in the introductory prose beneath the H3. Stage accepts a free-text label. You may omit or repeat a stage. The [universal job steps](../explanations/jobs-to-be-done.md#the-universal-job-map) offer conventional labels for exploring the work.
- **Outcomes.** Give each outcome a short heading and write its full statement beneath. Each outcome has one declaration owner. Under a job's H2 `Outcomes`, use an H3 for each name and an H4 for its fields. Under a step's H4 `Outcomes`, use H5 names and H6 fields. `Applies to` links to the jobs or steps for which it is a criterion. Ownership alone establishes no applicability. Other supported outcome fields include `Supporting content`, `Supports` and `Derived from`.
- **Identifiers.** An explicit job, step or outcome ID starts with an ASCII letter or digit and may also contain hyphens and underscores after the first character. IDs are case-sensitive and share a namespace within the file. Use a step ID for `--step`. Local fragments can resolve to explicit IDs or heading anchors. Heading anchors also work in ordinary Markdown viewers. [Reference checking](../../reference/inspection.md#references).

In the example, selecting `--step recall` returns no declared outcomes even though “Less reconstruction” explicitly applies to that step. Use `model needs` to inspect that link. [Worked applicability links](../../../examples/todo/specs/jobs/personal-commitments.job.spec.md#applies-to) · [JSON result fields](../../reference/inspection.md#output-and-automation).

### Checks and incomplete work

`init` and `check` accept the same presets. The choice affects the starter text or writing advice for that invocation. All three use the same file format.

| Preset | Writing guidance |
| --- | --- |
| `general` (default) | Prompts for a job statement, a performer and text for declared outcomes. Qualitative results are valid. |
| `odi` | Adds prompts about the functional boundary, an optional map and outcome statements with a direction and measure. |
| `job-story` | Adds prompts for situation, motivation and intended progress. |

Structural errors such as duplicate conventional sections, malformed metadata, conflicting explicit IDs and unresolved local references fail the check. Advisory findings allow a successful exit. Wording checks can recognize a measure word. Establishing that the proposed measure matters to performers requires research. [Write a job from a concrete situation](../guides/write-a-job.md) · [JTBD perspectives](../explanations/jobs-to-be-done.md).
