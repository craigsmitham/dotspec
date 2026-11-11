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

Enclosing system document: None
Source: specs/system.spec.md
````

## job show develop-requirements

```sh
dotspec job show develop-requirements
```

````text
# New job

## Job

Describe what someone seeks to accomplish or change, independent of the chosen solution.

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
  "schema_version": 2,
  "name": "New job",
  "path": "$PROJECT/specs/jobs/develop-requirements.job.spec.md",
  "declared_by_job": [],
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
    "advice": 0
  },
  "diagnostics": []
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
    "goal": "Name what the primary actor wants to accomplish, using an active verb phrase.",
    "boundary_description": "Name the system whose external behavior this use case describes.",
    "primary_actor": "Identify the role pursuing the goal.",
    "goal_level": "Identify the level: summary (kite), user task (sea), or subfunction (fish).",
    "trigger": "",
    "preconditions": "",
    "success_criterion": "",
    "main_success_scenario": "Describe the typical successful interaction in a short paragraph or a few steps.\nKeep internal implementation and interface details in their own sources.",
    "extensions": "List alternative success paths and failure conditions before resolving their handling.\nMark unknown responses as unresolved; expand the branches needed for the next decision.",
    "sequence": "",
    "supporting_actors": "",
    "stakeholder_interests": "",
    "minimum_guarantees": "",
    "related_jobs": "",
    "markdown": "# New use case\n\n## Goal\n\nName what the primary actor wants to accomplish, using an active verb phrase.\n\n## System boundary\n\nName the system whose external behavior this use case describes.\n\n## Primary actor\n\nIdentify the role pursuing the goal.\n\n## Goal level\n\nIdentify the level: summary (kite), user task (sea), or subfunction (fish).\n\n## Main success scenario\n\nDescribe the typical successful interaction in a short paragraph or a few steps.\nKeep internal implementation and interface details in their own sources.\n\n## Extensions\n\nList alternative success paths and failure conditions before resolving their handling.\nMark unknown responses as unresolved; expand the branches needed for the next decision.",
    "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_attributes": null
    },
    "location": {
      "line": 5,
      "column": 1,
      "element": "New use case"
    }
  }
}
````

## use-case init review-change --template expanded

```sh
dotspec use-case init review-change --template expanded
```

````text
Created specs/use-cases/review-change.use-case.spec.md
Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.
````

## use-case list --overview

```sh
dotspec use-case list --overview
```

````text
Use cases

System: New system (specs/system.spec.md)

  New use case (specs/use-cases/check-specifications.use-case.spec.md)
    Actor: Identify the role pursuing the goal.
    Goal: Name what the primary actor wants to accomplish, using an active verb phrase.
    Level: Identify the level: summary (kite), user task (sea), or subfunction (fish).
    Stated boundary: Name the system whose external behavior this use case describes.

  New use case (specs/use-cases/review-change.use-case.spec.md)
    Actor: Describe the role pursuing the goal.
    Goal: Describe what the primary actor wants to accomplish through this interaction.
    Level: Identify the level: summary, user task, or subfunction.
    Stated boundary: Name the system whose behavior this use case describes.
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
    "advice": 2
  },
  "diagnostics": [
    {
      "rule": "use-case.sequence",
      "severity": "advice",
      "message": "State whether scenario order is representative or required, and link to the governing rule for required order.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Sequence"
      }
    },
    {
      "rule": "use-case.success-criterion",
      "severity": "advice",
      "message": "State or link to the observable condition that completes the actor's goal.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Success criterion"
      }
    }
  ]
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
    "operational_principles": "Tell a characteristic story showing how the actions and state fulfill the purpose.\nUse multiple stories when needed to explain the essential behavior.\nThe story explains the design; it does not establish exhaustive behavior or verification.",
    "markdown": "# New concept\n\n## Purpose\n\nDescribe the benefit this concept provides and why it exists.\n\n## State\n\nDescribe the abstract information needed to explain the behavior, independently of its storage representation.\n\n## Actions\n\nDescribe or reference the available actions, their inputs and outputs, conditions, and effects.\nDistinguish caller assumptions from enforced guards, and state rejection behavior where it matters.\nUse subheadings for individual actions when useful.\n\n## Operational principles\n\nTell a characteristic story showing how the actions and state fulfill the purpose.\nUse multiple stories when needed to explain the essential behavior.\nThe story explains the design; it does not establish exhaustive behavior or verification.",
    "path": "$PROJECT/specs/concepts/job-map.concept.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_attributes": null
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
    "applicability": "Identify the subject and the conditions, workload, or population covered by this requirement.",
    "statement": "State one coherent quality obligation.\nIf a target is unknown, record the open question and the evidence needed to choose it.",
    "criterion": "Record optional measures, comparisons, targets, units and conditions for assessing the statement.",
    "verification_plan": "Describe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
    "measurement_plan": "Describe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.",
    "rationale": "Explain why this quality matters and any relevant tradeoffs.",
    "quality_attributes": [],
    "classification": {
      "system_path": "",
      "resolved": [],
      "ancestors": [],
      "undefined": []
    },
    "markdown": "# New quality requirement\n\n## Applies to\n\nIdentify the subject and the conditions, workload, or population covered by this requirement.\n\n## Statement\n\nState one coherent quality obligation.\nIf a target is unknown, record the open question and the evidence needed to choose it.\n\n## Verification plan\n\nDescribe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.\n\n## Criterion\n\nRecord optional measures, comparisons, targets, units and conditions for assessing the statement.\n\n## Measurement plan\n\nDescribe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.\n\n## Rationale\n\nExplain why this quality matters and any relevant tradeoffs.",
    "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_attributes": null
    },
    "location": {
      "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
      "line": 5,
      "column": 1,
      "element": "New quality requirement"
    }
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
    "applicability": "Identify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.",
    "statement": "State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.",
    "verification_plan": "Describe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
    "rationale": "Explain why this obligation is needed and reference its source or governing goal.",
    "quality_attributes": [],
    "classification": {
      "system_path": "",
      "resolved": [],
      "ancestors": [],
      "undefined": []
    },
    "markdown": "# New requirement\n\n## Applies to\n\nIdentify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.\n\n## Statement\n\nState one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.\n\n## Rationale\n\nExplain why this obligation is needed and reference its source or governing goal.\n\n## Verification plan\n\nDescribe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
    "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
    "system": {
      "name": "New system",
      "purpose": "Describe why this system exists and whom it serves.",
      "scope": "Describe the system's boundary and what it includes.",
      "path": "$PROJECT/specs/system.spec.md",
      "quality_attributes": null
    },
    "location": {
      "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
      "line": 5,
      "column": 1,
      "element": "New requirement"
    }
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
        "quality_attributes": null
      },
      "requirements": [
        {
          "type": "quality",
          "slug": "diagnostic-severity",
          "name": "New quality requirement",
          "applicability": "Identify the subject and the conditions, workload, or population covered by this requirement.",
          "statement": "State one coherent quality obligation.\nIf a target is unknown, record the open question and the evidence needed to choose it.",
          "criterion": "Record optional measures, comparisons, targets, units and conditions for assessing the statement.",
          "verification_plan": "Describe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
          "measurement_plan": "Describe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.",
          "rationale": "Explain why this quality matters and any relevant tradeoffs.",
          "quality_attributes": [],
          "classification": {
            "system_path": "",
            "resolved": [],
            "ancestors": [],
            "undefined": []
          },
          "markdown": "# New quality requirement\n\n## Applies to\n\nIdentify the subject and the conditions, workload, or population covered by this requirement.\n\n## Statement\n\nState one coherent quality obligation.\nIf a target is unknown, record the open question and the evidence needed to choose it.\n\n## Verification plan\n\nDescribe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.\n\n## Criterion\n\nRecord optional measures, comparisons, targets, units and conditions for assessing the statement.\n\n## Measurement plan\n\nDescribe how to assess the criterion, including sampling, measurement conditions, and evidence to retain.\n\n## Rationale\n\nExplain why this quality matters and any relevant tradeoffs.",
          "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
          "system": {
            "name": "New system",
            "purpose": "Describe why this system exists and whom it serves.",
            "scope": "Describe the system's boundary and what it includes.",
            "path": "$PROJECT/specs/system.spec.md",
            "quality_attributes": null
          },
          "location": {
            "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
            "line": 5,
            "column": 1,
            "element": "New quality requirement"
          }
        },
        {
          "type": "requirement",
          "slug": "reject-existing-target",
          "name": "New requirement",
          "applicability": "Identify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.",
          "statement": "State one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.",
          "verification_plan": "Describe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
          "rationale": "Explain why this obligation is needed and reference its source or governing goal.",
          "quality_attributes": [],
          "classification": {
            "system_path": "",
            "resolved": [],
            "ancestors": [],
            "undefined": []
          },
          "markdown": "# New requirement\n\n## Applies to\n\nIdentify the subject and applicable conditions. Distinguish assumptions from conditions the system must enforce.\n\n## Statement\n\nState one required result, rule, property, or restriction. Use general rules or contracts where they fit; prescribe ordering only when it is required.\n\n## Rationale\n\nExplain why this obligation is needed and reference its source or governing goal.\n\n## Verification plan\n\nDescribe how to check a system or process realization against this obligation, including conditions, success criteria and evidence to retain. This plan is not a completed verification result.",
          "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
          "system": {
            "name": "New system",
            "purpose": "Describe why this system exists and whom it serves.",
            "scope": "Describe the system's boundary and what it includes.",
            "path": "$PROJECT/specs/system.spec.md",
            "quality_attributes": null
          },
          "location": {
            "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
            "line": 5,
            "column": 1,
            "element": "New requirement"
          }
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

## sources --json

```sh
dotspec sources --json
```

````text
{
  "schema_version": 1,
  "scope": {
    "path": "$PROJECT",
    "types": [
      "quality-model",
      "business-rule",
      "system",
      "job",
      "use-case",
      "concept",
      "quality",
      "requirement"
    ],
    "discovery": "git"
  },
  "complete": true,
  "files": [
    {
      "path": "$PROJECT/specs/concepts/job-map.concept.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/jobs/develop-requirements.job.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/quality/diagnostic-severity.quality.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/system.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
      "referenced_by": []
    },
    {
      "path": "$PROJECT/specs/use-cases/review-change.use-case.spec.md",
      "referenced_by": []
    }
  ],
  "diagnostics": [
    {
      "rule": "use-case.sequence",
      "severity": "advice",
      "message": "State whether scenario order is representative or required, and link to the governing rule for required order.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Sequence"
      }
    },
    {
      "rule": "use-case.success-criterion",
      "severity": "advice",
      "message": "State or link to the observable condition that completes the actor's goal.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Success criterion"
      }
    }
  ]
}
````

