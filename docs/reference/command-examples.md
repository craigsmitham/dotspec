# Command output examples

These commands run in order in a new Git repository. `init` starts the project and the `new` commands create the specifications used by later commands. `$PROJECT` stands for that directory’s absolute path.

## init Bookings

```sh
dotspec init Bookings
```

````text
Created dotspec.yaml for project bookings
Created specs/bookings.system.spec.md
Edit Description, Purpose, and Scope. Keep the filename in sync if you rename the system.
Authoring guide: dotspec docs guides/describe-a-system
````

## new job develop-requirements

```sh
dotspec new job develop-requirements
```

````text
Created specs/jobs/develop-requirements.job.spec.md
Describe what the performer wants to accomplish and how they judge success.
Authoring guide: dotspec docs guides/write-a-job
````

## list system

```sh
dotspec list system
```

````text
Systems by directory
└── Bookings  (specs/bookings.system.spec.md)
````

## list job

```sh
dotspec list job
```

````text
Jobs

System: Bookings (specs/bookings.system.spec.md)
  Develop requirements  (specs/jobs/develop-requirements.job.spec.md)
````

## check specs/jobs/develop-requirements.job.spec.md --json

```sh
dotspec check specs/jobs/develop-requirements.job.spec.md --json
```

````text
{
  "schema_version": 1,
  "generated_spec": "out-of-scope",
  "scope": {
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT/specs/jobs",
    "types": [
      "job"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/dotspec.yaml",
    "$PROJECT/specs/bookings.system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 2,
    "errors": 0,
    "advice": 0
  },
  "diagnostics": []
}
````

## new use-case check-specifications

```sh
dotspec new use-case check-specifications
```

````text
Created specs/use-cases/check-specifications.use-case.spec.md
Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.
Authoring guide: dotspec docs guides/write-a-use-case
````

## new use-case review-change --template expanded

```sh
dotspec new use-case review-change --template expanded
```

````text
Created specs/use-cases/review-change.use-case.spec.md
Edit the actor goal and system boundary; add success criteria, scenarios, and extensions as they become known.
Authoring guide: dotspec docs guides/write-a-use-case
````

## list use-case --json

```sh
dotspec list use-case --json
```

````text
{
  "schema_version": 1,
  "entries": [
    {
      "kind": "use-case",
      "slug": "check-specifications",
      "name": "Check specifications",
      "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
      "system": {
        "name": "Bookings",
        "path": "$PROJECT/specs/bookings.system.spec.md"
      }
    },
    {
      "kind": "use-case",
      "slug": "review-change",
      "name": "Review change",
      "path": "$PROJECT/specs/use-cases/review-change.use-case.spec.md",
      "system": {
        "name": "Bookings",
        "path": "$PROJECT/specs/bookings.system.spec.md"
      }
    }
  ]
}
````

## check specs/use-cases/check-specifications.use-case.spec.md --json

```sh
dotspec check specs/use-cases/check-specifications.use-case.spec.md --json
```

````text
{
  "schema_version": 1,
  "generated_spec": "out-of-scope",
  "scope": {
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT/specs/use-cases",
    "types": [
      "use-case"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/dotspec.yaml",
    "$PROJECT/specs/bookings.system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 2,
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
      "rule": "use-case.success-guarantees",
      "severity": "advice",
      "message": "State or link to what holds when the actor's goal is achieved.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Success guarantees"
      }
    }
  ]
}
````

## new concept job-map

```sh
dotspec new concept job-map
```

````text
Created specs/concepts/job-map.concept.spec.md
Describe the purpose, state, actions, and operational principle as they become known.
Authoring guide: dotspec docs guides/specify-a-concept
````

## check specs/concepts/job-map.concept.spec.md --json

```sh
dotspec check specs/concepts/job-map.concept.spec.md --json
```

````text
{
  "schema_version": 1,
  "generated_spec": "out-of-scope",
  "scope": {
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT/specs/concepts",
    "types": [
      "concept"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/dotspec.yaml",
    "$PROJECT/specs/bookings.system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 2,
    "errors": 0,
    "advice": 0
  },
  "diagnostics": []
}
````

## new requirement diagnostic-severity

```sh
dotspec new requirement diagnostic-severity
```

````text
Created specs/requirements/diagnostic-severity.requirement.spec.md
Describe the scope, statement, any criterion, rationale, and verification as they become known.
Authoring guide: dotspec docs guides/write-an-assessable-requirement
````

## check specs/requirements/diagnostic-severity.requirement.spec.md --json

```sh
dotspec check specs/requirements/diagnostic-severity.requirement.spec.md --json
```

````text
{
  "schema_version": 1,
  "generated_spec": "out-of-scope",
  "scope": {
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT/specs/requirements",
    "types": [
      "requirement"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/dotspec.yaml",
    "$PROJECT/specs/bookings.system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 2,
    "errors": 0,
    "advice": 0
  },
  "diagnostics": []
}
````

## new requirement reject-existing-target

```sh
dotspec new requirement reject-existing-target
```

````text
Created specs/requirements/reject-existing-target.requirement.spec.md
Describe the scope, statement, any criterion, rationale, and verification as they become known.
Authoring guide: dotspec docs guides/write-an-assessable-requirement
````

## list requirement

```sh
dotspec list requirement
```

````text
Requirements

System: Bookings (specs/bookings.system.spec.md)
  Diagnostic severity  (specs/requirements/diagnostic-severity.requirement.spec.md)
  Reject existing target  (specs/requirements/reject-existing-target.requirement.spec.md)
````

## check specs/requirements/reject-existing-target.requirement.spec.md --json

```sh
dotspec check specs/requirements/reject-existing-target.requirement.spec.md --json
```

````text
{
  "schema_version": 1,
  "generated_spec": "out-of-scope",
  "scope": {
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT/specs/requirements",
    "types": [
      "requirement"
    ],
    "discovery": "direct"
  },
  "complete": true,
  "context_files": [
    "$PROJECT/dotspec.yaml",
    "$PROJECT/specs/bookings.system.spec.md"
  ],
  "summary": {
    "files_discovered": 1,
    "files_checked": 1,
    "context_files_checked": 2,
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
    "project": {
      "name": "bookings",
      "root": "$PROJECT"
    },
    "path": "$PROJECT",
    "types": [
      "quality-model",
      "business-rule",
      "system",
      "component",
      "feature",
      "capability",
      "quality-attribute",
      "job",
      "use-case",
      "concept",
      "requirement"
    ],
    "discovery": "git"
  },
  "complete": true,
  "files": [
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/dotspec.yaml",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/bookings.system.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/concepts/job-map.concept.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/jobs/develop-requirements.job.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/requirements/diagnostic-severity.requirement.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/requirements/reject-existing-target.requirement.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
      "referenced_by": [],
      "declarations": []
    },
    {
      "project": {
        "name": "bookings",
        "root": "$PROJECT"
      },
      "path": "$PROJECT/specs/use-cases/review-change.use-case.spec.md",
      "referenced_by": [],
      "declarations": []
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
      "rule": "use-case.success-guarantees",
      "severity": "advice",
      "message": "State or link to what holds when the actor's goal is achieved.",
      "location": {
        "path": "$PROJECT/specs/use-cases/check-specifications.use-case.spec.md",
        "line": 5,
        "column": 1,
        "element": "Success guarantees"
      }
    }
  ]
}
````

