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

## use-case init check-specifications

```sh
dotspec use-case init check-specifications
```

````text
Created specs/use-cases/check-specifications.use-case.spec.md
Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.
````

## use-case show check-specifications --json

```sh
dotspec use-case show check-specifications --json
```

````text
{
  "schema_version": 1,
  "use_case": {
    "slug": "check-specifications",
    "name": "New use case",
    "goal": "Describe what the primary actor wants to accomplish through this interaction.",
    "system_boundary": "Name the system whose behavior this use case describes.",
    "primary_actor": "Describe the role pursuing the goal.",
    "goal_level": "Identify the level: summary, user task, or subfunction.",
    "trigger": "Describe the action or time event that starts the interaction.",
    "preconditions": "Describe the state expected to be true before the interaction, or state that none is assumed.",
    "success_criterion": "State or link to the observable condition that completes the actor's goal.",
    "main_success_scenario": "1. Describe the actor's action.\n2. Describe the system's response and how the actor reaches the goal.",
    "extensions": "Describe exceptional conditions and their responses, linking to governing rules.",
    "sequence": "State whether the scenario's order is representative or required; link to the governing rule for required order.",
    "related_jobs": "Link to the jobs, steps, or desired outcomes this use case supports.",
    "markdown": "# New use case\n\n## Goal\n\nDescribe what the primary actor wants to accomplish through this interaction.\n\n## System boundary\n\nName the system whose behavior this use case describes.\n\n## Primary actor\n\nDescribe the role pursuing the goal.\n\n## Goal level\n\nIdentify the level: summary, user task, or subfunction.\n\n## Trigger\n\nDescribe the action or time event that starts the interaction.\n\n## Preconditions\n\nDescribe the state expected to be true before the interaction, or state that none is assumed.\n\n## Success criterion\n\nState or link to the observable condition that completes the actor's goal.\n\n## Main success scenario\n\n1. Describe the actor's action.\n2. Describe the system's response and how the actor reaches the goal.\n\n## Extensions\n\nDescribe exceptional conditions and their responses, linking to governing rules.\n\n## Sequence\n\nState whether the scenario's order is representative or required; link to the governing rule for required order.\n\n## Related jobs\n\nLink to the jobs, steps, or desired outcomes this use case supports.",
    "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_model": null
    },
    "location": {
      "line": 5,
      "column": 1,
      "element": "New use case"
    }
  }
}
````

## use-case check check-specifications --json

```sh
dotspec use-case check check-specifications --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "use-case"
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
    "advice": 0
  },
  "diagnostics": []
}
````

## concept init job-map

```sh
dotspec concept init job-map
```

````text
Created specs/concepts/job-map.concept.spec.md
Describe the purpose, state, actions, and operational principle as they become known.
````

## concept show job-map --json

```sh
dotspec concept show job-map --json
```

````text
{
  "schema_version": 1,
  "concept": {
    "slug": "job-map",
    "name": "New concept",
    "purpose": "Describe the benefit this concept provides and why it exists.",
    "state": "Describe the abstract information needed to explain the behavior, independently of its storage representation.",
    "actions": "Describe or reference the available actions, their inputs and outputs, conditions, and effects.\nDistinguish caller assumptions from enforced guards, and state rejection behavior where it matters.\nUse subheadings for individual actions when useful.",
    "operational_principle": "Tell a characteristic story showing how the actions and state fulfill the purpose.\nThe story explains the design; it does not establish exhaustive behavior or verification.",
    "markdown": "# New concept\n\n## Purpose\n\nDescribe the benefit this concept provides and why it exists.\n\n## State\n\nDescribe the abstract information needed to explain the behavior, independently of its storage representation.\n\n## Actions\n\nDescribe or reference the available actions, their inputs and outputs, conditions, and effects.\nDistinguish caller assumptions from enforced guards, and state rejection behavior where it matters.\nUse subheadings for individual actions when useful.\n\n## Operational principle\n\nTell a characteristic story showing how the actions and state fulfill the purpose.\nThe story explains the design; it does not establish exhaustive behavior or verification.",
    "path": "$PROJECT/specs/concepts/job-map.concept.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_model": null
    },
    "location": {
      "line": 5,
      "column": 1,
      "element": "New concept"
    }
  }
}
````

## concept check job-map --json

```sh
dotspec concept check job-map --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "concept"
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
    "advice": 0
  },
  "diagnostics": []
}
````

## quality init diagnostic-severity

```sh
dotspec quality init diagnostic-severity
```

````text
Created specs/quality/diagnostic-severity.quality.spec.md
Describe the scope, criterion, measurement, and rationale as they become known.
````

## quality show diagnostic-severity --json

```sh
dotspec quality show diagnostic-severity --json
```

````text
{
  "schema_version": 1,
  "quality": {
    "slug": "diagnostic-severity",
    "name": "New quality requirement",
    "applies_to": "Identify the subject and the conditions, workload, or population covered by this requirement.",
    "rationale": "Explain why this quality matters and any relevant tradeoffs.",
    "attributes": [],
    "classification": {
      "model_path": "",
      "resolved": [],
      "broader": [],
      "undefined": []
    },
    "markdown": "# New quality requirement\n\n## Applies to\n\nIdentify the subject and the conditions, workload, or population covered by this requirement.\n\n## Criterion\n\nState one coherent obligation, including its measure, comparison, and target where applicable.\nIf a target is unknown, record the open decision and the evidence needed to choose it.\n\n## Measurement\n\nDescribe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.\n\n## Rationale\n\nExplain why this quality matters and any relevant tradeoffs.",
    "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_model": null
    },
    "location": {
      "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
      "line": 5,
      "column": 1,
      "element": "New quality requirement"
    },
    "criterion": "State one coherent obligation, including its measure, comparison, and target where applicable.\nIf a target is unknown, record the open decision and the evidence needed to choose it.",
    "measurement": "Describe how to assess the criterion, including sampling, measurement conditions, and evidence to retain."
  }
}
````

## quality check diagnostic-severity --json

```sh
dotspec quality check diagnostic-severity --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "quality"
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
    "advice": 0
  },
  "diagnostics": []
}
````

## requirement init reject-existing-target

```sh
dotspec requirement init reject-existing-target
```

````text
Created specs/requirements/reject-existing-target.requirement.spec.md
Describe the scope, statement, rationale, and verification as they become known.
````

## requirement show reject-existing-target --json

```sh
dotspec requirement show reject-existing-target --json
```

````text
{
  "schema_version": 1,
  "requirement": {
    "type": "requirement",
    "slug": "reject-existing-target",
    "name": "New requirement",
    "applies_to": "Identify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.",
    "rationale": "Explain why this obligation is needed and reference its source or governing goal.",
    "attributes": [],
    "classification": {
      "model_path": "",
      "resolved": [],
      "broader": [],
      "undefined": []
    },
    "markdown": "# New requirement\n\n## Applies to\n\nIdentify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.\n\n## Statement\n\nState one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.\n\n## Rationale\n\nExplain why this obligation is needed and reference its source or governing goal.\n\n## Verification\n\nDescribe how to establish that the obligation holds, including success criteria and evidence to retain.",
    "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_model": null
    },
    "location": {
      "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
      "line": 5,
      "column": 1,
      "element": "New requirement"
    },
    "statement": "State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.",
    "verification": "Describe how to establish that the obligation holds, including success criteria and evidence to retain."
  }
}
````

## requirement list --json

```sh
dotspec requirement list --json
```

````text
{
  "schema_version": 1,
  "groups": [
    {
      "system": {
        "name": "New system",
        "purpose": "Describe why this system exists and whom it serves.",
        "scope": "Describe the system's boundary and what it includes.",
        "path": "$PROJECT/specs/system.spec.md",
        "quality_model": null
      },
      "requirements": [
        {
          "type": "quality",
          "slug": "diagnostic-severity",
          "name": "New quality requirement",
          "applies_to": "Identify the subject and the conditions, workload, or population covered by this requirement.",
          "rationale": "Explain why this quality matters and any relevant tradeoffs.",
          "attributes": [],
          "classification": {
            "model_path": "",
            "resolved": [],
            "broader": [],
            "undefined": []
          },
          "markdown": "# New quality requirement\n\n## Applies to\n\nIdentify the subject and the conditions, workload, or population covered by this requirement.\n\n## Criterion\n\nState one coherent obligation, including its measure, comparison, and target where applicable.\nIf a target is unknown, record the open decision and the evidence needed to choose it.\n\n## Measurement\n\nDescribe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.\n\n## Rationale\n\nExplain why this quality matters and any relevant tradeoffs.",
          "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
          "system": {
            "name": "New system",
            "purpose": "Describe why this system exists and whom it serves.",
            "scope": "Describe the system's boundary and what it includes.",
            "path": "$PROJECT/specs/system.spec.md",
            "quality_model": null
          },
          "location": {
            "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
            "line": 5,
            "column": 1,
            "element": "New quality requirement"
          },
          "statement": "State one coherent obligation, including its measure, comparison, and target where applicable.\nIf a target is unknown, record the open decision and the evidence needed to choose it.",
          "verification": "Describe how to assess the criterion, including sampling, measurement conditions, and evidence to retain."
        },
        {
          "type": "requirement",
          "slug": "reject-existing-target",
          "name": "New requirement",
          "applies_to": "Identify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.",
          "rationale": "Explain why this obligation is needed and reference its source or governing goal.",
          "attributes": [],
          "classification": {
            "model_path": "",
            "resolved": [],
            "broader": [],
            "undefined": []
          },
          "markdown": "# New requirement\n\n## Applies to\n\nIdentify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.\n\n## Statement\n\nState one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.\n\n## Rationale\n\nExplain why this obligation is needed and reference its source or governing goal.\n\n## Verification\n\nDescribe how to establish that the obligation holds, including success criteria and evidence to retain.",
          "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
          "system": {
            "name": "New system",
            "purpose": "Describe why this system exists and whom it serves.",
            "scope": "Describe the system's boundary and what it includes.",
            "path": "$PROJECT/specs/system.spec.md",
            "quality_model": null
          },
          "location": {
            "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
            "line": 5,
            "column": 1,
            "element": "New requirement"
          },
          "statement": "State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.",
          "verification": "Describe how to establish that the obligation holds, including success criteria and evidence to retain."
        }
      ]
    }
  ]
}
````

## requirement check reject-existing-target --json

```sh
dotspec requirement check reject-existing-target --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "requirement"
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
    "advice": 0
  },
  "diagnostics": []
}
````

