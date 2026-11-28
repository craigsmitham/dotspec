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

A recognizable grouping of specified product behavior valued by stakeholders, independent of component structure. Membership does not duplicate or own requirements.

[Model reference](index.md#feature)

### Topics

- Software requirements

### Relationships

- Can group a [requirement](#requirement).
- Can group a [use case](#use-case).
- Can group a [concept](#concept).

### Example

[Pick up where I left off](../../../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md) groups the [See what is left to do](../../../examples/todo/specs/use-cases/see-remaining.use-case.spec.md) use case.

## Job

What a performer seeks to accomplish or change in particular circumstances, independent of the chosen solution. Functional maps and job-story phrasing are optional complementary perspectives.

[Model reference](index.md#job)

### Topics

- User requirements

## Quality attribute

A named quality definition with acyclic grouping within its declaring model or effective system.

[Model reference](index.md#qualityattribute)

### Topics

- Software requirements

## Requirement

One required property or behavior, declared in its own requirement specification. Functional and quality requirements share this type. Either can reference a system quality attribute.

[Model reference](index.md#requirement)

### Topics

- Software requirements

## System

The system described by a system document; subordinate systems retain their own boundaries and quality attributes.

[Model reference](index.md#system)

### Topics

- Business requirements
- Software requirements

## System stakeholder role

A named capacity in which people or organizations have an interest in a particular system.

[Model reference](index.md#systemstakeholderrole)

### Topics

- Business requirements

## Use case

An actor goal and its interactions at a stated system boundary, with alternative paths and protections.

[Model reference](index.md#usecase)

### Topics

- Software requirements
