---
generated_by: dotspec
---
[dotspec](<dotspec.system.spec-report.md>)

# Companions

## Illustrated by

- [companions.txtar](<../internal/integration/testdata/companions.txtar>) — Declare, check, report and export companions · [List the companions that carry a specification's intent](<requirements/list-the-companions-that-carry-a-specification-s-intent.requirement.spec-report.md>)
- [docs.txtar](<../internal/integration/testdata/docs.txtar>) — Offline guidance · [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>)
- [help.txtar](<../internal/integration/testdata/help.txtar>) — Help, redirected output, completions and manual · [Find authoring guidance in the installed CLI](<requirements/find-authoring-guidance-in-the-installed-cli.requirement.spec-report.md>)
- [quality-convention.txtar](<../internal/integration/testdata/quality-convention.txtar>) — Adoption into the nearest system, refusal of a declared attribute and parent reuse · [Adopt a quality attribute for a system](<use-cases/adopt-quality-attribute.use-case.spec-report.md>)
- [quality-models.txtar](<../internal/integration/testdata/quality-models.txtar>) — Preview, bundled and local models, and a term with two parents · [Adopt a quality attribute for a system](<use-cases/adopt-quality-attribute.use-case.spec-report.md>)
- [init.txtar](<../internal/integration/testdata/init.txtar>) — A build without a project, a first build and the check that confirms it · [Build the generated specification](<use-cases/build-generated-specification.use-case.spec-report.md>)
- [check.py](<../scripts/check-example/check.py>) — Rebuilding the To-do generated specification and following its links · [Build the generated specification](<use-cases/build-generated-specification.use-case.spec-report.md>)
- [discovery.txtar](<../internal/integration/testdata/discovery.txtar>) — Repository discovery · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [diagnostics.txtar](<../internal/integration/testdata/diagnostics.txtar>) — Diagnostics and exit statuses · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [use-cases.txtar](<../internal/integration/testdata/use-cases.txtar>) — Focused and project-wide checks of draft and malformed use cases · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [boundary-advice.txtar](<../internal/integration/testdata/boundary-advice.txtar>) — Advice on assumptions without Basis or Applicability and interfaces without a contract link · [Check specifications before committing](<use-cases/check-specifications.use-case.spec-report.md>)
- [workflow.txtar](<../internal/integration/testdata/workflow.txtar>) — Create and inspect · [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>)
- [use-cases.txtar](<../internal/integration/testdata/use-cases.txtar>) — Use-case template selection, non-overwrite and inspection · [Create a specification draft](<use-cases/create-specification.use-case.spec-report.md>)
- [sources.txtar](<../internal/integration/testdata/sources.txtar>) — Source discovery through ordinary references and companions · [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>)
- [kinds_test.go](<../internal/cli/kinds_test.go>) — Lists of every specification kind, grouped by system · [Discover specification sources](<use-cases/discover-specification-sources.use-case.spec-report.md>)
- [snapshot.txtar](<../internal/integration/testdata/snapshot.txtar>) — Export and offline validation · [Export a specification snapshot](<use-cases/export-specification-model.use-case.spec-report.md>)
- [init.txtar](<../internal/integration/testdata/init.txtar>) — Starting a project, refusing to overwrite, and starting one inside another · [Start a specification project](<use-cases/start-specification-project.use-case.spec-report.md>)

## Verified by

- [quality_contracts_test.go › TestInspectionQualityContracts](<../internal/cli/quality_contracts_test.go#L38>) — CLI quality-contract tests for inspection · [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>)
- [quality_contracts_test.go › TestInspectionQualityContracts](<../internal/cli/quality_contracts_test.go#L38>) — CLI quality-contract tests for inspection · [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>)
- [quality_contracts_test.go › TestDiagnosticSeverityQualityContract](<../internal/cli/quality_contracts_test.go#L104>) — CLI quality-contract test for severity without styling · [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>)
