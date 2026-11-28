---
generated_by: dotspec
type: component-spec-report
source: "../../../app/specs/task-service.component.spec.md#task-service"
---
[Personal task list](<../../personal-task-list.system.spec-report.md>) › [Components](<../../personal-task-list.system.spec-report.md#components>)

# Task service

Checks whether a submitted action still applies to the selected entry and records accepted changes.

<details>
<summary>On this page</summary>

- [Responsibilities](<#responsibilities>)
- [Architecture](<#architecture>)
- [Contracts](<#contracts>)
- [Implementation](<#implementation>)

</details>

<details>
<summary>Browse specifications</summary>

- [Contracts](<contracts.md>)
- [Requirements](<requirements.md>)

</details>

## Responsibilities

A browser tab can still show Open after the owner has moved an errand to Waiting in another tab. When the service receives the old form, it compares the submitting view with the entry’s current state. A mismatch leaves the entry as it was and the browser receives a conflict. [Checklist action effects](<../../concepts/checklist.concept.spec-report.md#effects-table>).

The service supplies the current local day when it prepares the Waiting view. [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>) governs when it may report a successful change.

## Architecture

### Part of

[Personal task list](<../../personal-task-list.system.spec-report.md>).

### Capabilities

- [Follow-up review](<../../capabilities/follow-up-review.capability.spec-report.md>) — When the owner reviews waiting errands, they can see which replies they planned to chase by today and what to ask for.
- [Remember task progress between visits](<../../capabilities/remember-task-progress-between-visits.capability.spec-report.md>) — Remember my tasks and which ones I have checked off between visits. For a waiting task, keep the expected reply and follow-up day too.

## Contracts

- The browser submits changes through [Task list forms](<../../interfaces/task-list-forms.interface.spec-report.md>), whose [form template](<#specification.task-form-actions-and-fields>) supplies the actions and fields.

## Implementation

The [Go server](<../../../app/main.go>) owns form handling and storage.

## Specification sources

<a id="specification.task-form-actions-and-fields"></a>

### [Task form actions and fields](<../../../app/list.html>)

Defines contracts.

---

[View source](<../../../app/specs/task-service.component.spec.md#task-service>)
