# The dotspec handbook

## What am I trying to specify?

Choose the subject you want to write about. Open its explanation to understand which decisions belong there, or its how-to guide to draft and assess a specification.

Run `dotspec docs` from any directory to browse the topics offline.

<!-- Generated topic navigation -->

### Systems

Establish which promises the system makes and which conditions it relies on. Run `dotspec docs system`.

- [Describe a system and its boundary](guides/describe-a-system.md) (How-to)
- [Systems and their boundaries](explanations/systems.md) (Explanation)

### Components

Identify the part responsible for a behavior, its parent and its contracts. Run `dotspec docs component`.

- [Specify a component](guides/specify-a-component.md) (How-to)
- [Components and their responsibilities](explanations/components.md) (Explanation)

### Jobs to be done

Map a functional job and its success criteria, investigate switching, and connect user needs to capabilities. Run `dotspec docs job`.

- [Write a job from a concrete situation](guides/write-a-job.md) (How-to)
- [What makes something a job to be done?](explanations/jobs-to-be-done.md) (Explanation)
- [See personal commitments through](../../examples/todo/specs/jobs/personal-commitments.job.spec.md) (Example)
- [Connect needs to capabilities and specifications](guides/map-needs-to-capabilities.md#connect-needs-to-capabilities-and-specifications) (How-to)

### Capabilities

Connect an ability to the need it supports, the abilities it depends on and the parts that provide it. Run `dotspec docs capability`.

- [Specify a capability](guides/specify-a-capability.md) (How-to)
- [How capabilities connect needs to providers](explanations/capabilities.md) (Explanation)
- [Map needs to capabilities](guides/map-needs-to-capabilities.md) (How-to)

### Features

Group the behaviors that together deliver a benefit a stakeholder can recognize. Run `dotspec docs feature`.

- [Define the scope of a feature](guides/define-the-scope-of-a-feature.md) (How-to)
- [What holds a feature together?](explanations/features.md) (Explanation)

### Business rules

State what governs a business decision, identify its authority, and test its scope with concrete cases. Run `dotspec docs business-rule`.

- [Write a business rule](guides/write-a-business-rule.md) (How-to)
- [Business rules](explanations/business-rules.md) (Explanation)

### Use cases

Describe the paths an actor can take toward a goal and what the system promises at each ending. Run `dotspec docs use-case`.

- [Write a use case](guides/write-a-use-case.md) (How-to)
- [Use cases](explanations/use-cases.md) (Explanation)

### Concepts

Explain how a concept’s actions serve its purpose and how a system uses that behavior. Run `dotspec docs concept`.

- [Specify a concept](guides/specify-a-concept.md) (How-to)
- [Concepts in software design](explanations/concepts.md) (Explanation)

### Requirements

Write a promise a reader can assess and review whether related requirements can serve a stakeholder goal together. Run `dotspec docs requirement`.

- [Write an assessable requirement](guides/write-an-assessable-requirement.md) (How-to)
- [What makes a requirement assessable?](explanations/requirements.md) (Explanation)
- [Review a set of requirements](guides/review-a-set-of-requirements.md) (How-to)

### Quality models

Choose a shared vocabulary for the qualities that matter across systems. Run `dotspec docs quality-model`.

- [Define a quality model](guides/define-a-quality-model.md) (How-to)
- [What does a quality model help a team decide?](explanations/quality-models.md) (Explanation)

### System quality attributes

Define a quality for one system and connect it to requirements with assessable criteria. Run `dotspec docs quality-attribute`.

- [Define a quality attribute for a system](guides/define-a-system-quality-attribute.md) (How-to)
- [What does a quality mean for one system?](explanations/system-quality-attributes.md) (Explanation)
