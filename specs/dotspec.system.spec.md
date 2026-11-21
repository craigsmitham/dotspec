---
type: system-spec
glossary_sources:
  - ../docs/cli-handbook/reference/model/glossary.md
---

# dotspec

## Description

A requirements authoring framework and command-line toolkit for Markdown specifications stored alongside source code.

## Framing

The idea beheind dotspec is that having a distinctly seperated functional specification for a system and its parts is beneficial for both humans and agents

* Provides shared model and language for describing the intent and shape of a system
* Easy to infer (agents) and comprhened (humans)

dotspec consists of *.spec.md files that are co-located with the system or part/component of the system they specify. These files comprise the foundational "specification" layer for the system. Specs for a system can take many forms (formal methods, executable specifications, etc.), and dotspec is not intended to be the sole source of truth for all specifications, just those it chooses to specify. It is intended to support "progressive specifcitaion" with "just enough" specs to the extent they are helpful.

dotspec provides a rich but simple model for specifying requirements. Not all need to be used, but they can all be used and composed together if desired.

- identify/trace to other sources of truth for specs
- provide affordances for ensuring verification coverage/traceability
- non goal: perform verification of requirements or implementation


## Purpose

Help teams develop and maintain software requirements alongside the
code they describe.

## Goals


## Scope

dotspec covers the authoring framework, the `*.spec.md` files teams maintain
alongside their code, and the CLI that discovers and checks those files and
builds a requirements specification from them.

## Environment

Contributors edit Markdown in ordinary text editors. Git repositories provide
version history and discovery rules. The CLI communicates through process
arguments, standard streams, and exit statuses.

## Assumptions


## Constraints

## Glossary
