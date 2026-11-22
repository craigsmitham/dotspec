---
generated_by: dotspec
type: quality-model-spec-report
source: "../../specs/quality-models/personal-data.quality-model.spec.md#personal-data-quality"
---
# Personal data quality

## Description

A task list and an expense journal both hold decisions a person expects to find again. Use durability to ask which acknowledged records and revisions survive an interruption. Each tool chooses the records and interruptions its own requirements cover.

[Personal task list](<../personal-task-list.system.spec-report.md>)

## Quality attributes

### Durability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../../specs/quality-models/personal-data.quality-model.spec.md#durability"} -->

**ID:** durability

The ability to retain acknowledged changes when the application stops and starts again.
