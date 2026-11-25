---
generated_by: dotspec
type: quality-model-report
source: "../../specs/quality-models/personal-data.quality-model.md#personal-data-quality"
---
[Personal task list](<../personal-task-list.system.spec-report.md>) › [Quality models](<../quality-models.md>)

# Personal data quality

## Description

A task list and an expense journal both hold decisions a person expects to find again. Use durability to ask which acknowledged records and revisions survive an interruption. Each tool chooses the records and interruptions its own requirements cover.

## Sources

Adapted from recoverability in [ISO/IEC 25010:2023, §3.5.4](<https://www.iso.org/standard/78176.html>), which covers recovering the data directly affected by an interruption and re-establishing the system's desired state, and from [Q42 durability](<https://quality.arc42.org/qualities/durability>), under which a committed transaction survives system failures such as a power outage. This model keeps only the retention of changes a person sees acknowledged, so tools can discuss it apart from restoring a service or committing a transaction.

## Quality attributes

### Durability

<!-- dotspec-doc {"type":"quality-model-attribute-report","documents":"../../specs/quality-models/personal-data.quality-model.md#durability"} -->

[View section source](<../../specs/quality-models/personal-data.quality-model.md#durability>)

**ID:** durability

The ability to retain acknowledged changes when the application stops and starts again.

---

[View source](<../../specs/quality-models/personal-data.quality-model.md#personal-data-quality>)
