---
generated_by: dotspec
---
[dotspec](<dotspec.system.spec-report.md>)

# Companions

## Illustrated by

- [companions.txtar](<../internal/integration/testdata/companions.txtar>) — Declare, check, report and export companions · [List the companions that carry a specification's intent](<requirements/list-the-companions-that-carry-a-specification-s-intent.requirement.spec-report.md>)
- [docs.txtar](<../internal/integration/testdata/docs.txtar>) — Offline guidance · [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>)
- [help.txtar](<../internal/integration/testdata/help.txtar>) — Help, redirected output, completions and manual · [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>)
- [discovery.txtar](<../internal/integration/testdata/discovery.txtar>) — Repository discovery · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [diagnostics.txtar](<../internal/integration/testdata/diagnostics.txtar>) — Diagnostics and exit statuses · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [workflow.txtar](<../internal/integration/testdata/workflow.txtar>) — Create and inspect · [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>)
- [use-cases.txtar](<../internal/integration/testdata/use-cases.txtar>) — Use-case presets and inspection · [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>)
- [sources.txtar](<../internal/integration/testdata/sources.txtar>) — Source discovery through ordinary references and companions · [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>)
- [model.txtar](<../internal/integration/testdata/model.txtar>) — Export and offline validation · [Export a specification model for review or integration](<use-cases/export-specification-model.use-case.spec-report.md>)

## Verified by

- [quality_contracts_test.go › TestInspectionQualityContracts](<../internal/cli/quality_contracts_test.go#L38>) — CLI quality-contract tests for inspection · [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>)
- [quality_contracts_test.go › TestInspectionQualityContracts](<../internal/cli/quality_contracts_test.go#L38>) — CLI quality-contract tests for inspection · [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>)
- [quality_contracts_test.go › TestDiagnosticSeverityQualityContract](<../internal/cli/quality_contracts_test.go#L104>) — CLI quality-contract test for severity without styling · [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>)
