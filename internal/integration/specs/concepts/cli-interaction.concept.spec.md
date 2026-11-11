---
type: concept
id: cli-interaction
---

# Command Invocation

## Purpose

Let a person or program submit an operation and determine its outcome through a
consistent process interface.

## State

One invocation has arguments, an input environment, output and error streams, and
an eventual exit status. The command reference defines available operations;
individual concept and requirement specifications define their effects on files.
The invocation concept does not own those operations' domain behavior.

The [command reference](../../../../docs/cli.md) describes the available inputs.
The [output contract](../../../../README.md#output-and-automation) defines process
outputs, formats and exit statuses.

## Actions

### Invoke

Submit a supported operation and its arguments. Invalid invocation syntax produces
an error result rather than performing the requested operation.

### Observe the result

Read the output and error streams and the exit status when the process completes.
A successful invocation reports exit status 0; specification findings and invocation
or operational errors use the distinct outcomes defined by the output contract.

## Operational principle

Invoke an operation, read its result, and use the exit status to distinguish success
from findings or an operational failure. A program can make its next decision from
that result without scraping presentation prose for an implied success claim.

The workflow of creating, editing and inspecting specifications spans several
concepts. The linked executable cases below describe those interactions and their
expected process outcomes; they do not make every command a separate action of
this concept.

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
- [Quality models and system population](../../testdata/quality-models.txtar)
- [Quality requirements and attributes](../../testdata/quality.txtar)
- [Business-rule commands](../../testdata/business-rules.txtar)
- [General requirements](../../testdata/requirements.txtar)
- [Source discovery through ordinary file references](../../testdata/sources.txtar)
- [Export and validate source-backed specification models](../../testdata/model.txtar)

The Go test harness executes these scripts; it does not define additional
product obligations. The repository's [specification-test convention](../../../../CONTRIBUTING.md#specification-tests-and-implementation-tests)
identifies the sources that its coverage script requires.
