# Command examples

Generated from the CLI and its starter templates. `$PROJECT` stands for the absolute project directory.

## system init

```sh
dotspec system init
```

````text
Created specs/system.spec.md
Edit the system name, Purpose, and Scope.
````

## job init develop-requirements

```sh
dotspec job init develop-requirements
```

````text
Created specs/jobs/develop-requirements.job.spec.md
Edit the job statement and performer; add context, steps, and outcomes as they become known.
````

## system show

```sh
dotspec system show
```

````text
New system

Purpose
  Describe why this system exists and whom it serves.

Scope
  Describe the system's boundary and what it includes.

Parent system: None
Source: specs/system.spec.md
````

## job show develop-requirements

```sh
dotspec job show develop-requirements
```

````text
# New job

## Job

Describe what someone wants to accomplish, including where this job begins and ends.

## Performer

Describe who is trying to accomplish this job.

## Situation

Describe the circumstances in which this job arises.

## Motivation

Describe why someone wants to act in this situation.

## Intended progress

Describe the progress or result they want to achieve.
Use an optional Outcomes section for individual success criteria.

System: New system (specs/system.spec.md)
Source: specs/jobs/develop-requirements.job.spec.md
````

## job step list develop-requirements

```sh
dotspec job step list develop-requirements
```

````text
No steps defined.
````

## job outcome list develop-requirements --json

```sh
dotspec job outcome list develop-requirements --json
```

````text
{
  "schema_version": 1,
  "name": "New job",
  "path": "$PROJECT/specs/jobs/develop-requirements.job.spec.md",
  "overall": [],
  "steps": []
}
````

## job map develop-requirements --format mermaid

```sh
dotspec job map develop-requirements --format mermaid
```

````text
flowchart TD
  job["New job"]
  no_steps["No steps defined."]
  job --- no_steps
````

## job check develop-requirements --json

```sh
dotspec job check develop-requirements --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "job"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/specs/system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 1,
    "errors": 0,
    "advice": 1
  },
  "diagnostics": [
    {
      "rule": "job.boundary",
      "severity": "advice",
      "message": "Consider clarifying the execution boundary; no context marker such as 'for', 'when', or 'during' was found.",
      "location": {
        "path": "$PROJECT/specs/jobs/develop-requirements.job.spec.md",
        "line": 7,
        "column": 1,
        "element": "Job"
      }
    }
  ]
}
````

