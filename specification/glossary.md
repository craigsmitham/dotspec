---
generated_by: dotspec
---
# Glossary

[dotspec](<dotspec.system.spec-report.md>)

## Browse by topic

- **Business requirements** · [Business rule](<#term.business-rule>), [Capability](<#term.capability>), [System](<#term.system>)
- **Software requirements** · [Capability](<#term.capability>), [Component](<#term.component>), [Concept](<#term.concept>), [Feature](<#term.feature>), [Quality attribute](<#term.quality-attribute>), [Requirement](<#term.requirement>), [System](<#term.system>), [Use case](<#term.use-case>)
- **User requirements** · [Job](<#term.job>)

<a id="term.business-rule"></a>

## Business rule

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#business-rule"} -->

A declarative business constraint independent of the software or process that enforces it.

[Model reference](<../docs/reference/model/index.md#businessrule>)

<a id="term.business-rule.topics"></a>

### Topics

- Business requirements

[View term source](<../docs/reference/model/glossary.md#business-rule>)

<a id="term.capability"></a>

## Capability

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#capability"} -->

An ability needed to deliver value, potentially combining people, procedures and technology. Distinct from a product feature or structural component; its declaration does not establish a provider.

[Model reference](<../docs/reference/model/index.md#capability>)

<a id="term.capability.topics"></a>

### Topics

- Business requirements
- Software requirements

[View term source](<../docs/reference/model/glossary.md#capability>)

<a id="term.component"></a>

## Component

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#component"} -->

A named part of a system with a described responsibility and boundary. Neither its declaration nor code location establishes structural parentage.

[Model reference](<../docs/reference/model/index.md#component>)

<a id="term.component.topics"></a>

### Topics

- Software requirements

[View term source](<../docs/reference/model/glossary.md#component>)

<a id="term.concept"></a>

## Concept

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#concept"} -->

A behavior-bearing concept described by purpose, abstract state, actions and an operational principle.

[Model reference](<../docs/reference/model/index.md#concept>)

Relationships from [Feature](<#term.feature.relationships>)

<a id="term.concept.topics"></a>

### Topics

- Software requirements

[View term source](<../docs/reference/model/glossary.md#concept>)

<a id="term.feature"></a>

## Feature

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#feature"} -->

A recognizable grouping of specified product behavior valued by stakeholders, independent of component structure. Membership does not duplicate or own requirements.

[Model reference](<../docs/reference/model/index.md#feature>)

<a id="term.feature.topics"></a>

### Topics

- Software requirements

<a id="term.feature.relationships"></a>

### Relationships

- Can group a [requirement](<#term.requirement>).
- Can group a [use case](<#term.use-case>).
- Can group a [concept](<#term.concept>).

<a id="term.feature.example"></a>

### Example

[Pick up where I left off](<../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md>) groups the [Decide what to do about a waiting errand](<../examples/todo/specs/use-cases/resume-list.use-case.spec.md>) use case.

[View term source](<../docs/reference/model/glossary.md#feature>)

<a id="term.job"></a>

## Job

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#job"} -->

What a performer seeks to accomplish or change in particular circumstances, independent of the chosen solution. Functional maps and job-story phrasing are optional complementary perspectives.

[Model reference](<../docs/reference/model/index.md#job>)

<a id="term.job.topics"></a>

### Topics

- User requirements

[View term source](<../docs/reference/model/glossary.md#job>)

<a id="term.quality-attribute"></a>

## Quality attribute

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#quality-attribute"} -->

A named quality attribute owned by exactly one system or reusable quality model. Attributes may be organized under several other attributes within that owner. This acyclic grouping supports trees and overlapping classifications; it does not inherit requirements, thresholds, or satisfaction.

[Model reference](<../docs/reference/model/index.md#qualityattribute>)

<a id="term.quality-attribute.topics"></a>

### Topics

- Software requirements

[View term source](<../docs/reference/model/glossary.md#quality-attribute>)

<a id="term.requirement"></a>

## Requirement

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#requirement"} -->

One required property or behavior, declared in its own requirement specification. Functional and quality requirements share this type. Either can reference a system quality attribute.

[Model reference](<../docs/reference/model/index.md#requirement>)

Relationships from [Feature](<#term.feature.relationships>)

<a id="term.requirement.topics"></a>

### Topics

- Software requirements

[View term source](<../docs/reference/model/glossary.md#requirement>)

<a id="term.system"></a>

## System

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#system"} -->

The system described by a system document; subordinate systems retain their own boundaries and quality attributes.

[Model reference](<../docs/reference/model/index.md#system>)

<a id="term.system.topics"></a>

### Topics

- Business requirements
- Software requirements

[View term source](<../docs/reference/model/glossary.md#system>)

<a id="term.use-case"></a>

## Use case

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#use-case"} -->

An actor goal and its interactions at a stated system boundary, with alternative paths and protections.

[Model reference](<../docs/reference/model/index.md#usecase>)

Relationships from [Feature](<#term.feature.relationships>)

<a id="term.use-case.topics"></a>

### Topics

- Software requirements

[View term source](<../docs/reference/model/glossary.md#use-case>)

[View source for “Glossary”](<../specs/dotspec.system.spec.md#glossary>)
