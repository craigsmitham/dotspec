---
generated_by: dotspec
---
[dotspec](<dotspec.system.spec-report.md>)

# Glossary

**Browse by topic**

- **Business requirements** · [Business rule](<#business-rule>), [Capability](<#capability>), [System](<#system>), [System stakeholder role](<#system-stakeholder-role>)
- **Software requirements** · [Capability](<#capability>), [Component](<#component>), [Concept](<#concept>), [Feature](<#feature>), [Quality attribute](<#quality-attribute>), [Requirement](<#requirement>), [System](<#system>), [Use case](<#use-case>)
- **User requirements** · [Job](<#job>)

## Business rule

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#business-rule"} -->

A declarative business constraint independent of the software or process that enforces it.

Business requirements · [Model reference](<../docs/reference/model/index.md#businessrule>) · [Term source](<../docs/reference/model/glossary.md#business-rule>)

## Capability

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#capability"} -->

An ability needed to deliver value, potentially combining people, procedures and technology. Distinct from a product feature or structural component; its declaration does not establish a provider.

Business requirements, Software requirements · [Model reference](<../docs/reference/model/index.md#capability>) · [Term source](<../docs/reference/model/glossary.md#capability>)

## Component

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#component"} -->

A named part of a system with a described responsibility and boundary. Neither its declaration nor code location establishes structural parentage.

Software requirements · [Model reference](<../docs/reference/model/index.md#component>) · [Term source](<../docs/reference/model/glossary.md#component>)

## Concept

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#concept"} -->

A behavior-bearing concept described by purpose, abstract state, actions and an operational principle.

Software requirements · [Model reference](<../docs/reference/model/index.md#concept>) · [Term source](<../docs/reference/model/glossary.md#concept>)

From [Feature](<#feature-relationships>) · Can group a concept.

## Feature

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#feature"} -->

A recognizable grouping of specified product behavior valued by stakeholders, independent of component structure. Membership does not duplicate or own requirements.

Software requirements · [Model reference](<../docs/reference/model/index.md#feature>) · [Term source](<../docs/reference/model/glossary.md#feature>)

### Feature relationships

- Can group a [requirement](<#requirement>).
- Can group a [use case](<#use-case>).
- Can group a [concept](<#concept>).

### Feature example

[Pick up where I left off](<../examples/todo/specs/features/pick-up-where-i-left-off.feature.spec.md>) groups the [Decide what to do about a waiting errand](<../examples/todo/specs/use-cases/resume-list.use-case.spec.md>) use case.

## Job

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#job"} -->

What a performer seeks to accomplish or change in particular circumstances, independent of the chosen solution. Functional maps and job-story phrasing are optional complementary perspectives.

User requirements · [Model reference](<../docs/reference/model/index.md#job>) · [Term source](<../docs/reference/model/glossary.md#job>)

## Quality attribute

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#quality-attribute"} -->

A named quality definition with acyclic grouping within its declaring model or effective system.

Software requirements · [Model reference](<../docs/reference/model/index.md#qualityattribute>) · [Term source](<../docs/reference/model/glossary.md#quality-attribute>)

## Requirement

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#requirement"} -->

One required property or behavior, declared in its own requirement specification. Functional and quality requirements share this type. Either can reference a system quality attribute.

Software requirements · [Model reference](<../docs/reference/model/index.md#requirement>) · [Term source](<../docs/reference/model/glossary.md#requirement>)

From [Feature](<#feature-relationships>) · Can group a requirement.

## System

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#system"} -->

The system described by a system document; subordinate systems retain their own boundaries and quality attributes.

Business requirements, Software requirements · [Model reference](<../docs/reference/model/index.md#system>) · [Term source](<../docs/reference/model/glossary.md#system>)

## System stakeholder role

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#system-stakeholder-role"} -->

A named capacity in which people or organizations have an interest in a particular system.

Business requirements · [Model reference](<../docs/reference/model/index.md#systemstakeholderrole>) · [Term source](<../docs/reference/model/glossary.md#system-stakeholder-role>)

## Use case

<!-- dotspec-doc {"type":"glossary-term-spec-report","documents":"../docs/reference/model/glossary.md#use-case"} -->

An actor goal and its interactions at a stated system boundary, with alternative paths and protections.

Software requirements · [Model reference](<../docs/reference/model/index.md#usecase>) · [Term source](<../docs/reference/model/glossary.md#use-case>)

From [Feature](<#feature-relationships>) · Can group a use case.
