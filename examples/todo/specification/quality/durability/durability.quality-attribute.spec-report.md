---
generated_by: dotspec
type: system-quality-attribute-spec-report
source: "../../../specs/quality/durability/durability.quality-attribute.spec.md#durability"
---
[Personal task list](<../../personal-task-list.system.spec-report.md>) › [Quality attributes](<../../personal-task-list.system.spec-report.md#quality-attributes>)

# Durability

## Definition

Preservation of the owner’s acknowledged decisions between local sessions, including the reply they expect for a waiting errand and the day they chose to follow up. The follow-up cue can change overnight even though those decisions survive intact.

The [time-zone cases](<#specification.saved-follow-up-dates-across-time-zones>) govern what the owner sees when the same saved list is opened on different local days. The suite exercises those cases on Unix.

[Source definition](<../../quality-models/personal-data-quality.quality-model-report.md#durability>).

[How SQLite keeps a committed transaction across a crash or power loss](<https://www.sqlite.org/atomiccommit.html>).

## Specification sources

<a id="specification.saved-follow-up-dates-across-time-zones"></a>

### [Saved follow-up dates across time zones](<../../../tests/acceptance.py#L278>)

Defines requirements.

#### Verification

- [Run the acceptance suite on Unix](<../../../tests/acceptance.py#L324>)

No matching result recorded.

## Verification

- [Recover acknowledged changes after restart](<../../../tests/acceptance.py#L177>)

No matching result recorded.

## Related specifications

**Classifies**

- [Acknowledge after commit](<../../components/task-service/requirements/acknowledge-after-commit.requirement.spec-report.md>)
- [Preserve tasks across restart](<../../requirements/preserve-tasks-across-restart.requirement.spec-report.md>)

---

[View source](<../../../specs/quality/durability/durability.quality-attribute.spec.md#durability>)
