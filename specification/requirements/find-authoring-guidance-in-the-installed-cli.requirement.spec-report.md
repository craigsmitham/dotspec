---
generated_by: dotspec
type: requirement-spec-report
source: "../../specs/requirements/offline-authoring-guidance.requirement.spec.md#find-authoring-guidance-in-the-installed-cli"
---
[dotspec](<../dotspec.system.spec-report.md>) › [Requirements](<../requirements.md>)

# Find authoring guidance in the installed CLI

## Statement

A contributor must be able to reach the relevant versioned guidance from a specification type, a question in their own words, or a mapped diagnostic, then read the selected page, section, or worked example offline. The guidance must identify what the contributor needs to carry out its instructions with that installation.

## Applies to

Contributors using an installed version of dotspec, including in an empty directory without a project or network connection.

## Criterion

Root help identifies the handbook. Type help opens a route to the topic. Creating a draft reports its path and a direct reading command. A topic and its synonym show the same ordered destinations, each with its title, documentation form and reading command. Searching by words gives title and topic matches priority over body matches, with repeatable ordering. An exact address opens its destination. An invalid address reports the failed lookup with a route back to browsing.

A section includes its subordinate content. Offline reading commands for linked sections preserve their anchors. Redirecting a whole example reproduces the bundled source. The CLI handbook index and knowledge index lead readers to the same catalog destinations as the installed CLI.

Only the Markdown in `docs/cli-handbook/` and the selected worked examples and model or specification references are available for offline reading. Search and shell completion follow the same publication boundary. A nested handbook page becomes addressable when the CLI is rebuilt. Repository-only citations identify their repository paths. Maintainer pages and plans are excluded from the executable, knowledge bundle and platform archives.

Required local reading in handbook pages and linked walkthroughs opens in the installed CLI, including the requested heading. A supplemental checkout citation identifies its purpose and checkout prerequisite at the link. Such a citation can be skipped while completing the stated authoring task.

Before an instruction to run a command, the reader can identify the working directory, selected project and files it assumes. When an example uses the dotspec checkout, the guidance identifies how to do the same work in the reader’s own project where applicable. A bundled specification can be read from an empty directory. Inspecting or editing a project uses files in the selected project.

## Evidence

The [publication check](<../../internal/integration/handbook_test.go>) follows authored local links from outside the checkout and checks supplemental citations against their repository targets. The [offline discovery scenarios](<../../internal/integration/testdata/docs.txtar>) exercise help, topic selection, search, diagnostic lookup, exact destinations, completion and redirection. The [package check](<../../scripts/package-knowledge/test_package.py>) compares index links with the topic catalog.

### Reader exercise

The [handbook reader exercise](<../../docs/maintainers/handbook.md#try-the-guidance-from-an-installation>) asks a contributor to revise a thin job statement and identify what still needs evidence.

## Rationale

A structurally complete draft can still leave the reader guessing about the performer’s circumstances. Authors need guidance while they are writing.

## Related specifications

**Quality attribute**

- [Usability](<../quality/usability/usability.quality-attribute.spec-report.md>)

---

[View source](<../../specs/requirements/offline-authoring-guidance.requirement.spec.md#find-authoring-guidance-in-the-installed-cli>)
