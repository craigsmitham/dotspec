# CLI reference

Run project commands in the directory containing your specifications. Substitute your own names, paths and project selectors in the examples. Handbook reading commands work from any directory.

## dotspec

Manage software requirements alongside your code

```sh
dotspec [flags]
```

Manage human-readable software requirements in *.spec.md files.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec docs job
dotspec system init
dotspec job init develop-requirements
dotspec check
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for dotspec |
| `--version, -v` | `false` | version for dotspec |

## dotspec business-rule

Describe business constraints independently of their enforcement

```sh
dotspec business-rule [flags]
```

Authoring handbook: dotspec docs business-rule

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for business-rule |

## dotspec business-rule check

Check structure and suggest missing business-rule content

```sh
dotspec business-rule check <slug-or-source-path> [flags]
```

Check structure and suggest missing business-rule content.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs business-rule

```sh
dotspec business-rule check purchase-approval --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec business-rule init

Create a business-rule specification without overwriting an existing file

```sh
dotspec business-rule init <slug> [flags]
```

Authoring handbook: dotspec docs business-rule

```sh
dotspec business-rule init purchase-approval --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec business-rule list

Discover business rules, grouped by associated system

```sh
dotspec business-rule list [flags]
```

```sh
dotspec business-rule list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec business-rule show

Show authored Markdown and the associated system

```sh
dotspec business-rule show <slug-or-source-path> [flags]
```

```sh
dotspec business-rule show purchase-approval --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec capability

Describe an ability and its providers

```sh
dotspec capability [flags]
```

Authoring handbook: dotspec docs capability

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for capability |

## dotspec capability check

Check the specification and its references

```sh
dotspec capability check <slug-or-source-path> [flags]
```

Check the specification and its references.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs capability

```sh
dotspec capability check local-storage --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec capability init

Create a standalone specification

```sh
dotspec capability init <slug> [flags]
```

Authoring handbook: dotspec docs capability

```sh
dotspec capability init local-storage --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec capability list

List specifications in the selected subtree

```sh
dotspec capability list [flags]
```

```sh
dotspec capability list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec capability show

Show the authored specification

```sh
dotspec capability show <slug-or-source-path> [flags]
```

```sh
dotspec capability show local-storage --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec check

Check specifications and reference documents in a project subtree

```sh
dotspec check [source-path] [flags]
```

Check authored specifications and reference documents using the same discovery rules as list. A whole-project check also inspects generated Markdown for projects declared with dotspec.yaml unless report.enabled is false. These automated authoring and selected requirements-quality checks support requirements verification; passing them does not establish full requirements verification, requirements validation, or system verification/validation.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec check
dotspec check --path internal --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec completion

Generate a shell completion script

```sh
dotspec completion [flags]
```

Generate a completion script, then load it using your shell's completion setup.

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for completion |

## dotspec completion bash

Generate the bash completion script

```sh
dotspec completion bash [flags]
```

```sh
dotspec completion bash
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for bash |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion fish

Generate the fish completion script

```sh
dotspec completion fish [flags]
```

```sh
dotspec completion fish
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for fish |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion powershell

Generate the powershell completion script

```sh
dotspec completion powershell [flags]
```

```sh
dotspec completion powershell
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for powershell |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec completion zsh

Generate the zsh completion script

```sh
dotspec completion zsh [flags]
```

```sh
dotspec completion zsh
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for zsh |
| `--no-descriptions` | `false` | omit completion descriptions |

## dotspec component

Describe a component in its own specification

```sh
dotspec component [flags]
```

Authoring handbook: dotspec docs component

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for component |

## dotspec component init

Create a named component specification

```sh
dotspec component init <name> [flags]
```

Authoring handbook: dotspec docs component

```sh
dotspec component init Intake --path internal --part-of specs/new-system.system.spec.md
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--part-of` | `` | Parent system or component source file, relative to the working directory |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec concept

Describe purposeful product behavior through state and actions

```sh
dotspec concept [flags]
```

Authoring handbook: dotspec docs concept

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for concept |

## dotspec concept check

Check structure and suggest missing concept content

```sh
dotspec concept check <slug-or-source-path> [flags]
```

Check structure and suggest missing concept content.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs concept

```sh
dotspec concept check job-map --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec concept init

Create a concept specification without overwriting an existing file

```sh
dotspec concept init <slug> [flags]
```

Authoring handbook: dotspec docs concept

```sh
dotspec concept init job-map --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec concept list

Discover concepts, grouped by associated system

```sh
dotspec concept list [flags]
```

```sh
dotspec concept list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec concept show

Show authored Markdown and the associated system

```sh
dotspec concept show <slug-or-source-path> [flags]
```

```sh
dotspec concept show job-map --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec docs

Browse the offline authoring handbook

```sh
dotspec docs [TOPIC | PAGE[#SECTION] | RULE | WORDS...] [flags]
```

Browse topics, open a page or section, or search this version's handbook and bundled examples. Topic names and exact addresses open directly. Enter words to find matching pages.

```sh
dotspec docs
dotspec docs job
dotspec docs success criteria
dotspec docs guides/write-a-job
dotspec docs job.performer
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for docs |

## dotspec feature

Group stakeholder-visible behavior

```sh
dotspec feature [flags]
```

Authoring handbook: dotspec docs feature

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for feature |

## dotspec feature check

Check the specification and its references

```sh
dotspec feature check <slug-or-source-path> [flags]
```

Check the specification and its references.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs feature

```sh
dotspec feature check resume-work --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec feature init

Create a standalone specification

```sh
dotspec feature init <slug> [flags]
```

Authoring handbook: dotspec docs feature

```sh
dotspec feature init resume-work --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec feature list

List specifications in the selected subtree

```sh
dotspec feature list [flags]
```

```sh
dotspec feature list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec feature show

Show the authored specification

```sh
dotspec feature show <slug-or-source-path> [flags]
```

```sh
dotspec feature show resume-work --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec help

Help about any command

```sh
dotspec help [command] [flags]
```

Help provides help for any command in the application.
Simply type dotspec help [path to command] for full details.

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for help |

## dotspec job

Define and discover jobs to be done (JTBDs)

```sh
dotspec job [flags]
```

Authoring handbook: dotspec docs job

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for job |

## dotspec job check

Check structure and suggest writing improvements

```sh
dotspec job check <slug-or-source-path> [flags]
```

Check structure and suggest writing improvements.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs job

```sh
dotspec job check develop-requirements --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--preset` | `general` | JTBD authoring guidance: general, odi, or job-story |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec job init

Create a job specification without overwriting an existing file

```sh
dotspec job init <slug> [flags]
```

Authoring handbook: dotspec docs job

```sh
dotspec job init develop-requirements --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--preset` | `general` | JTBD authoring guidance: general, odi, or job-story |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec job list

Discover jobs, grouped by associated system

```sh
dotspec job list [flags]
```

```sh
dotspec job list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec job map

Map authored step order (text, JSON, or Mermaid)

```sh
dotspec job map <slug> [flags]
```

```sh
dotspec job map develop-requirements --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json or mermaid |
| `--help, -h` | `false` | help for map |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--outcomes` | `false` | include job and step outcomes |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec job outcome

List outcomes by declaration owner

```sh
dotspec job outcome [flags]
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for outcome |

## dotspec job outcome list

List outcomes by declaration owner

```sh
dotspec job outcome list <slug> [flags]
```

```sh
dotspec job outcome list develop-requirements --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--step` | `` | select the declaring step by its explicit ID |

## dotspec job show

Show authored Markdown and the associated system

```sh
dotspec job show <slug-or-source-path> [flags]
```

```sh
dotspec job show develop-requirements --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec man

Write the dotspec manual in roff format

```sh
dotspec man [flags]
```

```sh
dotspec man
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for man |

## dotspec model

Export and validate complete specification snapshots

```sh
dotspec model [flags]
```

Export source text, named domain fields, located relationships and inspection findings in a versioned specification model.

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for model |

## dotspec model export

Export a source-backed specification snapshot as JSON

```sh
dotspec model export [flags]
```

Export selected specifications, reusable quality models, enclosing systems, and direct reference context as JSON. Retain readable Markdown and source locations, including diagnostics for invalid documents. Filesystem discovery and authoring findings follow check; external references are not fetched.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec model export
dotspec model export --path internal
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for export |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec model needs

Generate a needs-to-specification map with capability dependencies

```sh
dotspec model needs [flags]
```

Show jobs, shared outcomes, features, requirements and optional capabilities. Short labels link to full statements and evidence. Only explicit relationships are projected; contribution does not prove achievement.

```sh
dotspec model needs --path . --relative-to views
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for needs |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--relative-to` | `.` | Directory containing the generated Markdown; source links are relative to it |

## dotspec model overview

Generate system structure, features, and applicable requirements as Markdown

```sh
dotspec model overview [flags]
```

Show explicit system parts, non-owning feature groups, applicability, concept uses, native references and requirement criteria. Directory context remains separate; no implementation or requirement inheritance is inferred.

```sh
dotspec model overview --path . --relative-to views
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for overview |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--relative-to` | `.` | Directory containing the generated Markdown; source links are relative to it |

## dotspec model pegs

Generate a Project, Environment, Goals, System reading perspective

```sh
dotspec model pegs [flags]
```

Generate a Project, Environment, Goals, System reading perspective. Requirements verification concerns specification quality; requirements and system validation concern needs and intended use. Only acquired content and explicit relationships are shown; missing information produces advisory review prompts, not completeness or endorsement claims.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec model pegs --path . --relative-to views
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for pegs |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--relative-to` | `.` | Directory containing the generated Markdown; source links are relative to it |

## dotspec model product

Show the specification roots, owned elements and references for one system

```sh
dotspec model product [SYSTEM] [flags]
```

```sh
dotspec model product
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for product |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec model review

Generate advisory requirements review and open-question inventory

```sh
dotspec model review [flags]
```

Generate advisory requirements review and open-question inventory. Requirements verification concerns specification quality; requirements and system validation concern needs and intended use. Only acquired content and explicit relationships are shown; missing information produces advisory review prompts, not completeness or endorsement claims.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec model review --path . --relative-to views
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for review |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--relative-to` | `.` | Directory containing the generated Markdown; source links are relative to it |

## dotspec model schema

Print the embedded snapshot JSON Schema

```sh
dotspec model schema [flags]
```

```sh
dotspec model schema
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for schema |
| `--workspace` | `false` | print the schema for a complete model export --all envelope |

## dotspec model trace

Generate layered requirements traceability and change-review views as Markdown

```sh
dotspec model trace [flags]
```

Show motivation, derivation, refinement, assumptions, use-case participation, reverse change-review links and requirement criteria. Resolved explicit relationships only; paths do not establish satisfaction or evidence.

```sh
dotspec model trace --path . --relative-to views
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--help, -h` | `false` | help for trace |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--relative-to` | `.` | Directory containing the generated Markdown; source links are relative to it |

## dotspec model validate

Check a snapshot's shape, source ranges, and relationships

```sh
dotspec model validate SNAPSHOT [flags]
```

Validate a snapshot file, or use - to read standard input. Snapshot validation checks representation conformance and integrity, not requirements fitness or system behavior; faithfully represented authoring errors do not make a snapshot invalid. No source paths or external targets are fetched.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for validate |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |

## dotspec quality-attribute

Browse a quality model and adopt system quality attributes

```sh
dotspec quality-attribute [flags]
```

Authoring handbook: dotspec docs quality-attribute

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for quality-attribute |

## dotspec quality-attribute add

Adopt a quality attribute and its missing parents into one system

```sh
dotspec quality-attribute add <attribute-id> [flags]
```

```sh
dotspec quality-attribute add recoverability --dry-run
```

| Flag | Default | Description |
| --- | --- | --- |
| `--dry-run` | `false` | Preview the decisions and new file contents |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for add |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--model` | `` | Use a bundled model, local slug, or source path for this invocation |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--system` | `` | Target system source path, relative to --path |

## dotspec quality-attribute browse

Browse top-level terms or inspect a term and its neighbors

```sh
dotspec quality-attribute browse [attribute-id] [flags]
```

```sh
dotspec quality-attribute browse reliability
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for browse |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--model` | `` | Use a bundled model, local slug, or source path for this invocation |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec quality-model

Define reusable quality attributes and copy them into systems

```sh
dotspec quality-model [flags]
```

Authoring handbook: dotspec docs quality-model

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for quality-model |

## dotspec quality-model apply

Adopt selected quality attributes into one system

```sh
dotspec quality-model apply [model] [flags]
```

Select IDs with --attributes or request the whole model with --all. Existing system declarations keep their wording and grouping. To explore the model first, use quality-attribute browse.

```sh
dotspec quality-model apply --attributes recoverability --dry-run
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | Adopt every attribute from the model into one system |
| `--attributes` | `[]` | Comma-separated attribute IDs to adopt |
| `--dry-run` | `false` | Preview the decisions and new file contents |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for apply |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--system` | `` | Target system source path, relative to --path |

## dotspec quality-model check

Check a local quality model's structure and suggest missing content

```sh
dotspec quality-model check <slug-or-source-path> [flags]
```

Check a local quality model's structure and suggest missing content.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs quality-model

```sh
dotspec quality-model check organization-quality --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec quality-model init

Create an editable project quality model

```sh
dotspec quality-model init <slug> [flags]
```

Authoring handbook: dotspec docs quality-model

```sh
dotspec quality-model init organization-quality --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--from` | `` | Copy a bundled starter into the new project model |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec quality-model list

List bundled starters and project quality models

```sh
dotspec quality-model list [flags]
```

```sh
dotspec quality-model list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec quality-model show

Show a bundled starter or a local quality model

```sh
dotspec quality-model show <slug-or-source-path> [flags]
```

```sh
dotspec quality-model show furps --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec requirement

Specify required results, rules, properties, and restrictions

```sh
dotspec requirement [flags]
```

Authoring handbook: dotspec docs requirement

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for requirement |

## dotspec requirement check

Check requirement specifications and attribute references

```sh
dotspec requirement check <slug-or-source-path> [flags]
```

Check requirement specifications and attribute references.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs requirement

```sh
dotspec requirement check reject-existing-target --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec requirement init

Create a requirement specification without overwriting an existing file

```sh
dotspec requirement init <slug> [flags]
```

Authoring handbook: dotspec docs requirement

```sh
dotspec requirement init reject-existing-target --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--attribute` | `` | Write this governing system attribute ID in quality_attributes. |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--preset` | `general` | Authoring template: general, or quality to prompt for an explicit criterion |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec requirement list

Discover requirements, grouped by associated system

```sh
dotspec requirement list [flags]
```

```sh
dotspec requirement list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--attribute` | `` | Include requirements with this attribute or a sub-attribute in their system's model |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec requirement show

Show authored Markdown and the associated system

```sh
dotspec requirement show <slug-or-source-path> [flags]
```

```sh
dotspec requirement show reject-existing-target --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec sources

List specifications, reference documents, system context, and directly referenced files

```sh
dotspec sources [flags]
```

Inventory sources using the same scan and reference checks as check. File references do not assign authority or classify their contents.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

```sh
dotspec sources
dotspec sources --path internal --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for sources |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec spec

Generate the project's readable requirements specification

```sh
dotspec spec [flags]
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for spec |

## dotspec spec build

Build SPEC.md and the supporting Markdown tree

```sh
dotspec spec build [flags]
```

Build the requirements specification for one whole project. Reports are enabled by default in projects declared with dotspec.yaml. Validate the candidate before replacing generator-owned output. Run dotspec check to verify the resulting documents.

```sh
dotspec spec build
```

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for build |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec system

Describe systems and discover their directory context

```sh
dotspec system [flags]
```

Authoring handbook: dotspec docs system

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for system |

## dotspec system init

Create a named system specification

```sh
dotspec system init [name] [flags]
```

Authoring handbook: dotspec docs system

```sh
dotspec system init --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec system list

Display systems in the selected subtree

```sh
dotspec system list [flags]
```

```sh
dotspec system list --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec system show

Describe the nearest system and its enclosing system document

```sh
dotspec system show [source-path] [flags]
```

```sh
dotspec system show --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec use-case

Describe actor goals and interactions with a system

```sh
dotspec use-case [flags]
```

Authoring handbook: dotspec docs use-case

| Flag | Default | Description |
| --- | --- | --- |
| `--help, -h` | `false` | help for use-case |

## dotspec use-case check

Check structure and suggest missing use-case content

```sh
dotspec use-case check <slug-or-source-path> [flags]
```

Check structure and suggest missing use-case content.

Exit codes: 0 success (advice is nonfatal), 1 specification findings, 2 invalid invocation or operational failure.

Authoring handbook: dotspec docs use-case

```sh
dotspec use-case check check-specifications --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for check |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec use-case init

Create a use-case specification without overwriting an existing file

```sh
dotspec use-case init <slug> [flags]
```

Authoring handbook: dotspec docs use-case

```sh
dotspec use-case init check-specifications --path .
dotspec use-case init review-change --template expanded --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for init |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |
| `--template` | `brief` | Writing prompts: brief or expanded (same document format) |

## dotspec use-case list

Discover use cases, grouped by associated system

```sh
dotspec use-case list [flags]
```

```sh
dotspec use-case list --path .
dotspec use-case list --overview --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--all` | `false` | inspect all discovered projects separately |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for list |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--overview` | `false` | Show actor, goal, level, and stated boundary in text; JSON retains the full records |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

## dotspec use-case show

Show authored Markdown and the associated system

```sh
dotspec use-case show <slug-or-source-path> [flags]
```

```sh
dotspec use-case show check-specifications --path .
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `text` | output format: text or json |
| `--help, -h` | `false` | help for show |
| `--json` | `false` | emit structured JSON (equivalent to --format json) |
| `--path` | `.` | project directory to inspect or initialize |
| `--project` | `` | select a project by its dotspec.yaml name (--path is relative to its root) |

