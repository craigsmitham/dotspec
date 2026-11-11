---
type: system
---
# Document workspace

## Purpose

Help teams preserve and recover valuable work.

## Scope

Document deletion and restoration within a workspace.

## Goals

### Reduce lost work

Reduce the amount of valuable work lost through accidental deletion.

## Assumptions

### Trusted time

The platform supplies a trustworthy UTC clock for retention calculations.

## Concept uses

### DocumentTrash

#### Definition

[Trash](concepts/trash.concept.spec.md)

#### Supports

[Recover lost work](jobs/recover.job.spec.md).

#### Relies on

[Trusted time](#trusted-time) bounds retention intervals.
