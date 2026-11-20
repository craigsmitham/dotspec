---
generated_by: dotspec
type: component-spec-report
source: "../../../app/specs/task-service.component.spec.md#task-service"
---
# Task service

The local service that handles browser requests and saves changes to the task list.

[Personal task list](<../../personal-task-list.system.spec-report.md>)

| Section | Explore |
| --- | --- |
| [Responsibilities](<#responsibilities>) | Accept browser changes and persist them locally. See Acknowledge after commit for the confirmation guarantee. |
| [System topology](<#system-topology>) | Part of |
| Capabilities | [All capabilities](<capabilities.md>) |
| Requirements | [All requirements](<requirements.md>) |
| [Contracts](<#contracts>) | The browser uses the HTML forms to submit changes. |
| [Implementation](<#implementation>) | The Go server owns form handling and storage. |

## Responsibilities

Accept browser changes and persist them locally. See [Acknowledge after commit](<requirements/acknowledge-after-commit.requirement.spec-report.md>) for the confirmation guarantee.

## System topology

### Part of

[Personal task list](<../../personal-task-list.system.spec-report.md>).

```mermaid
flowchart LR
    n0["Task service"]
    n1["Personal task list"]
    n0 -->|part of| n1
    n2["Local task storage"]
    n0 -->|provides| n2
```

- **Parent system** · [Personal task list](<../../personal-task-list.system.spec-report.md>) — A single-user task list that stores open and completed tasks on one computer.
- **Provided capability** · [Local task storage](<../../capabilities/local-task-storage.capability.spec-report.md>) — Retains accepted task changes between sessions on the same computer.

## Contracts

- The browser uses the [HTML forms](<../../../app/list.html>) to submit changes.

## Implementation

The [Go server](<../../../app/main.go>) owns form handling and storage.
