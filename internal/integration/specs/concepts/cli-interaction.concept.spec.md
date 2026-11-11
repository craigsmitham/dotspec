---
type: concept
id: cli-interaction
---

# CLI Interaction

## Purpose

Let contributors and programs author, inspect, and discover specification
material through ordinary command invocations.

## State

An invocation receives arguments and the selected project files. It produces
standard output, standard error, and an exit status. Creation commands can add
specification files; read commands leave source material unchanged.

The [command reference](../../../../docs/cli.md) describes the available inputs.
The [output contract](../../../../README.md#output-and-automation) describes
formats and process behavior.

## Actions

Create a specification, read its authored content, inspect a selected scope,
or inventory its source files. The executable examples below specify required
outcomes for their illustrated inputs and conditions. Their sequences are
representative workflows, not a universal ordering constraint on commands.

## Operational principle

Create and edit specifications, then inspect them before sharing a change.
When connecting other specification sources, add ordinary Markdown links and
use the source inventory to check what can be found. Read the sources and their
surrounding prose to understand what they define or verify.

## Related specifications

These scripts are the authoritative executable examples for the cases they
exercise. They can supplement general rules defined elsewhere without
restating those rules here. A passing example does not settle untested cases.

- [Create and inspect specifications](../../testdata/workflow.txtar)
- [Repository discovery](../../testdata/discovery.txtar)
- [Diagnostics and exit statuses](../../testdata/diagnostics.txtar)
- [Help and redirected output](../../testdata/help.txtar)
- [Use-case commands](../../testdata/use-cases.txtar)
- [Concept commands](../../testdata/concepts.txtar)
- [Quality requirements and attributes](../../testdata/quality.txtar)
- [General requirements](../../testdata/requirements.txtar)
- [Source discovery through ordinary file references](../../testdata/sources.txtar)
- [Export and validate source-backed specification models](../../testdata/model.txtar)

The Go test harness executes these scripts; it does not define additional
product obligations. The repository's [specification-test convention](../../../../CONTRIBUTING.md#specification-tests-and-implementation-tests)
identifies the sources that its coverage script requires.
