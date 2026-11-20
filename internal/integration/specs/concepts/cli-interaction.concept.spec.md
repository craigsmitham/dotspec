---
type: concept-spec
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

The [command reference](../../../../docs/reference/cli.md) describes the available inputs.
The [output contract](../../../../docs/reference/inspection.md#output-and-automation) defines process
outputs, formats and exit statuses.

## Actions

### Invoke

Submit a supported operation and its arguments. Invalid invocation syntax produces
an error result rather than performing the requested operation.

### Observe the result

Read the output and error streams and the exit status when the process completes.
A successful invocation reports exit status 0; specification findings and invocation
or operational errors use the distinct outcomes defined by the output contract.

## Operational principles

Invoke an operation, read its result, and use the exit status to distinguish success
from findings or an operational failure. A program can make its next decision from
that result without scraping presentation prose for an implied success claim.

The workflow of creating, editing and inspecting specifications spans several
concepts. The linked executable cases below describe those interactions and their
expected process outcomes; they do not make every command a separate action of
this concept.

## Related specifications

Each script supplies its inputs, commands and expected observations. Read the governing specification alongside the cases it exercises.

| Governing source | Executable cases |
| --- | --- |
| [Create a specification](../../../../specs/use-cases/create-specification.use-case.spec.md), including the linked non-overwrite requirement | [Create and inspect](../../testdata/workflow.txtar), [use-case presets and inspection](../../testdata/use-cases.txtar) |
| [Inspection interaction](../../../../specs/use-cases/check-specifications.use-case.spec.md) and [source selection](../../../../docs/reference/layout.md) | [Repository discovery](../../testdata/discovery.txtar), [diagnostics and exit statuses](../../testdata/diagnostics.txtar) |
| [Command arguments](../../../../docs/reference/cli.md) and [process output](../../../../docs/reference/inspection.md#output-and-automation) | [Offline guidance](../../testdata/guide.txtar), [help, redirected output, completions and manual](../../testdata/help.txtar) |
| [Concept format](../../../../docs/reference/formats/concept.md) and [business-rule authoring guide](../../../../docs/handbook/business-rules.md) | [Concept commands](../../testdata/concepts.txtar), [business-rule commands](../../testdata/business-rules.txtar) |
| [System quality attributes](../../../../docs/reference/formats/system.md#system-quality-attributes), [quality-model population](../../../../docs/handbook/quality-models.md) and [requirement format](../../../../docs/reference/formats/requirement.md) | [Quality models and system population](../../testdata/quality-models.txtar), [classified requirements](../../testdata/classified-requirements.txtar), [general requirements](../../testdata/requirements.txtar) |
| [Discover specification sources](../../../../specs/use-cases/discover-specification-sources.use-case.spec.md) and [export a specification model](../../../../specs/use-cases/export-specification-model.use-case.spec.md) | [Source discovery through ordinary references](../../testdata/sources.txtar), [export and offline validation](../../testdata/model.txtar) |

The Go test harness executes these scripts; it does not define additional
product requirements. The repository's [specification-test convention](../../../../CONTRIBUTING.md#specification-tests-and-implementation-tests)
identifies the sources that its coverage script requires.
