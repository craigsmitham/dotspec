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
concepts. Its executable cases are companions of the
[create a specification draft](../../../../specs/use-cases/create-specification.use-case.spec.md)
and [check specifications before committing](../../../../specs/use-cases/check-specifications.use-case.spec.md)
use cases. Scripts that no specification lists as a companion yet are linked below
beside the sources they exercise; they do not make every command a separate action
of this concept.

## Related specifications

Each script supplies its inputs, commands and expected observations. Read the governing specification alongside the cases it exercises.

| Governing source | Executable cases |
| --- | --- |
| [Concept format](../../../../model/markdown-mapping.md#concept) and [business-rule mapping](../../../../model/markdown-mapping.md#business-rule) | [Concept commands](../../testdata/concepts.txtar), [business-rule commands](../../testdata/business-rules.txtar) |
| [System quality attributes](../../../../docs/reference/model/index.md#systemqualityattribute), [quality attribute adoption](../../../../docs/reference/cli.md#dotspec-new-quality-attribute) and [requirement format](../../../../model/markdown-mapping.md#requirement) | [Explicit model selection and incremental adoption](../../testdata/quality-convention.txtar), [quality models and system population](../../testdata/quality-models.txtar), [classified requirements](../../testdata/classified-requirements.txtar), [general requirements](../../testdata/requirements.txtar) |

The Go test harness executes these scripts; it does not define additional
product requirements. The repository's [specification-test convention](../../../../CONTRIBUTING.md#specification-tests-and-implementation-tests)
identifies the sources that its coverage script requires.
