---
generated_by: dotspec
type: concept-spec-report
source: "../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#command-invocation"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Concept definitions](<../dotspec.system.spec-report.md#concept.command-invocation>)

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

[View section source](<../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#invoke>)

Submit a supported operation and its arguments. Invalid invocation syntax produces
an error result rather than performing the requested operation.

### Observe the result

<!-- dotspec-doc {"type":"concept-action-spec-report","documents":"../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#observe-the-result"} -->

[View section source](<../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#observe-the-result>)

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
| [Create a specification](<../use-cases/create-specification.use-case.spec-report.md>), including the linked non-overwrite requirement | [Create and inspect](<../../internal/integration/testdata/workflow.txtar>), [use-case presets and inspection](<../../internal/integration/testdata/use-cases.txtar>) |
| [Inspection interaction](<../use-cases/check-specifications.use-case.spec-report.md>) and [source selection](<../../docs/reference/layout.md>) | [Repository discovery](<../../internal/integration/testdata/discovery.txtar>), [diagnostics and exit statuses](<../../internal/integration/testdata/diagnostics.txtar>) |
| [Offline authoring guidance](<../requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>), [command arguments](<../../docs/reference/cli.md>) and [process output](<../../docs/reference/inspection.md#output-and-automation>) | [Offline guidance](<../../internal/integration/testdata/docs.txtar>), [help, redirected output, completions and manual](<../../internal/integration/testdata/help.txtar>) |
| [Concept format](<../../model/markdown-mapping.md#concept>) and [business-rule mapping](<../../model/markdown-mapping.md#business-rule>) | [Concept commands](<../../internal/integration/testdata/concepts.txtar>), [business-rule commands](<../../internal/integration/testdata/business-rules.txtar>) |
| [System quality attributes](<../../docs/reference/model/index.md#systemqualityattribute>), [quality-model application](<../../docs/reference/cli.md#dotspec-quality-model-apply>) and [requirement format](<../../model/markdown-mapping.md#requirement>) | [Project defaults and incremental adoption](<../../internal/integration/testdata/quality-convention.txtar>), [quality models and system population](<../../internal/integration/testdata/quality-models.txtar>), [classified requirements](<../../internal/integration/testdata/classified-requirements.txtar>), [general requirements](<../../internal/integration/testdata/requirements.txtar>) |
| [Discover specification sources](<../use-cases/discover-specification-sources.use-case.spec-report.md>) and [export a specification model](<../use-cases/export-specification-model.use-case.spec-report.md>) | [Source discovery through ordinary references](<../../internal/integration/testdata/sources.txtar>), [export and offline validation](<../../internal/integration/testdata/model.txtar>) |

The Go test harness executes these scripts; it does not define additional
product requirements. The repository's [specification-test convention](<../../CONTRIBUTING.md#specification-tests-and-implementation-tests>)
identifies the sources that its coverage script requires.

---

[View source](<../../internal/integration/specs/concepts/cli-interaction.concept.spec.md#command-invocation>)
