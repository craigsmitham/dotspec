# Getting the sink fixed

```mermaid
stateDiagram-v2
    Open --> Waiting: Ask the plumber for a visit time
    Waiting --> Open: Plumber asks for a photo
    Open --> Completed: Repair is finished
    Completed --> Open: Leak returns
```
