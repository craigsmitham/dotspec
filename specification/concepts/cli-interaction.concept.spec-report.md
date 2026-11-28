---
generated_by: dotspec
type: concept-spec-report
source: "../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#command-invocation"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Concepts](<../dotspec.system.spec-report.md#concepts>)

# Command Invocation

## Purpose

Let a person or program submit an operation and determine its outcome through a
consistent process interface.

## State

One invocation has arguments, an input environment, output and error streams, and
an eventual exit status. The command reference defines available operations;
individual concept and requirement specifications define their effects on files.
The invocation concept does not own those operations' domain behavior.

The [command reference](<../../docs/reference/cli.md>) describes the available inputs.
The [output contract](<../../docs/reference/inspection.md#output-and-automation>) defines process
outputs, formats and exit statuses.

## Actions

[Invoke](<#invoke>) · [Observe the result](<#observe-the-result>)

### Invoke

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#invoke"} -->

Submit a supported operation and its arguments. Invalid invocation syntax produces
an error result rather than performing the requested operation.

### Observe the result

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#observe-the-result"} -->

Read the output and error streams and the exit status when the process completes.
A successful invocation reports exit status 0; specification findings and invocation
or operational errors use the distinct outcomes defined by the output contract.

## Operational principles

Invoke an operation, read its result, and use the exit status to distinguish success
from findings or an operational failure. A program can make its next decision from
that result without scraping presentation prose for an implied success claim.

The workflow of creating, editing and inspecting specifications spans several
concepts. Its executable cases are named sources on the
[create a specification draft](<../use-cases/create-specification.use-case.spec-report.md>)
and [check specifications before committing](<../use-cases/check-specifications.use-case.spec-report.md>)
use cases. Scripts that no specification names yet are linked below
beside the sources they exercise; they do not make every command a separate action
of this concept.

## Related specifications

Each script supplies its inputs, commands and expected observations. Read the governing specification alongside the cases it exercises.

| Governing source | Executable cases |
| --- | --- |
| [Concept format](<../../model/markdown-mapping.md#concept>) and [business-rule mapping](<../../model/markdown-mapping.md#business-rule>) | [Concept commands](<../../internal/integration/testdata/concepts.txtar>), [business-rule commands](<../../internal/integration/testdata/business-rules.txtar>) |
| [System quality attributes](<../../docs/reference/model/index.md#systemqualityattribute>), [quality attribute adoption](<../../docs/reference/cli.md#dotspec-new-quality-attribute>) and [requirement format](<../../model/markdown-mapping.md#requirement>) | [Explicit model selection and incremental adoption](<../../internal/integration/testdata/quality-convention.txtar>), [quality models and system population](<../../internal/integration/testdata/quality-models.txtar>), [classified requirements](<../../internal/integration/testdata/classified-requirements.txtar>), [general requirements](<../../internal/integration/testdata/requirements.txtar>) |

The Go test harness executes these scripts; it does not define additional
product requirements. The repository's [specification-test convention](<../../CONTRIBUTING.md#specification-tests-and-implementation-tests>)
identifies the sources that its coverage script requires.

---

[View source](<../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#command-invocation>)
