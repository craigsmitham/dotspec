# Dotspec model glossary

Generated from the [metamodel](../../../model/metamodel.yaml).

## Business rule

A declarative business constraint independent of the software or process that enforces it.

[Model reference](index.md#businessrule)

### Topics

- Business requirements

## Capability

An ability needed to deliver value, potentially combining people, procedures and technology. Distinct from a product feature or structural component; its declaration does not establish a provider.

[Model reference](index.md#capability)

### Topics

- Business requirements
- Software requirements

## Component

A named part of a system with a described responsibility and boundary. Neither its declaration nor code location establishes structural parentage.

[Model reference](index.md#component)

### Topics

- Software requirements

## Concept

A behavior-bearing concept described by purpose, abstract state, actions and an operational principle.

[Model reference](index.md#concept)

### Topics

- Software requirements

## Feature

A recognizable grouping of specified product behavior valued by stakeholders, independent of component structure. Membership does not duplicate or own obligations.

[Model reference](index.md#feature)

### Topics

- Software requirements

### Relationships

- Can group a [requirement](#requirement).
- Can group an [obligation clause](#obligation-clause).
- Can group a [use case](#use-case).
- Can group a [concept](#concept).

### Example

[Pick up where I left off](../../../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md) groups the [Return to my list](../../../examples/todo/specs/use-cases/resume-list.use-case.spec.md) use case.

## Job

What a performer seeks to accomplish or change in particular circumstances, independent of the chosen solution. Functional maps, progress cases and job-story phrasing are optional complementary perspectives.

[Model reference](index.md#job)

### Topics

- User requirements

## Obligation clause

An explicitly distinguished obligation inside another specification or behavioral model. Its content is the authored expression, not a duplicate requirement document. Project authority determines acceptance; declaration alone does not.

[Model reference](index.md#obligationclause)

### Topics

- Business requirements
- User requirements
- Software requirements

## Quality attribute

A named quality attribute owned by exactly one system or reusable quality model. Attributes may be organized under several other attributes within that owner. This acyclic grouping supports trees and overlapping classifications; it does not inherit obligations, thresholds, or satisfaction.

[Model reference](index.md#qualityattribute)

### Topics

- Software requirements

## Requirement

A specification that gives one obligation its own document. An obligation is what ISO/IEC/IEEE 29148 calls a requirement, wherever it is written; this class names the specification type, not every normative expression. Functional and quality requirements share this type. Either can reference a quality attribute.

[Model reference](index.md#requirement)

### Topics

- Software requirements

## System

The system described by a system document; subordinate systems retain their own boundaries and quality attributes.

[Model reference](index.md#system)

### Topics

- Business requirements
- Software requirements

## Use case

An actor goal and its interactions at a stated system boundary, with alternative paths and protections.

[Model reference](index.md#usecase)

### Topics

- Software requirements
