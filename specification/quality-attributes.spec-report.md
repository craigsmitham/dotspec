---
generated_by: dotspec
---
# Quality attributes

[dotspec](<dotspec.system.spec-report.md>)

These are dotspec's declared quality attributes. The definitions identify concerns;
individual quality requirements establish the obligations.

## Usability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#usability"} -->

**ID:** usability

How effectively contributors can author, inspect, and maintain specifications.

**Classifies** · [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>)

## Diagnosability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#diagnosability"} -->

**ID:** diagnosability\
**Sub-attribute of:** usability

How effectively findings help contributors understand and locate specification
problems.

**Classifies**

- [Continue checking after a structural error](<requirements/continue-after-invalid-spec.requirement.spec-report.md>)
- [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>)

## Accessibility

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#accessibility"} -->

**ID:** accessibility\
**Sub-attribute of:** usability

How effectively people with differing abilities and ways of accessing output
can use dotspec.

**Classifies** · [Diagnostic severity is available without styling](<requirements/diagnostic-severity.requirement.spec-report.md>)

## Reliability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#reliability"} -->

**ID:** reliability

How consistently dotspec performs its intended operations under stated
conditions.

## Repeatability

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#repeatability"} -->

**ID:** repeatability\
**Sub-attribute of:** reliability

Whether equivalent inspections under unchanged conditions produce equivalent
reports.

**Classifies**

- [Read and verify the generated requirements specification](<requirements/read-and-verify-the-generated-requirements-specification.requirement.spec-report.md>)
- [Repeated inspections produce equivalent reports](<requirements/repeatable-inspection.requirement.spec-report.md>)

## Source preservation

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#source-preservation"} -->

**ID:** source-preservation\
**Sub-attribute of:** reliability

Whether inspection preserves the source material it examines.

**Classifies**

- [Inspection preserves source material](<requirements/inspection-preserves-source.requirement.spec-report.md>)
- [Preserve source and represented content in model exports](<requirements/model-source-fidelity.requirement.spec-report.md>)
- [Reject initialization at an existing target](<requirements/reject-existing-target.requirement.spec-report.md>)

## Performance

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#performance"} -->

**ID:** performance

The time and computing resources required to perform operations under stated
workloads.

## Responsiveness

<!-- dotspec-doc {"type":"quality-attribute-spec-report","documents":"../specs/dotspec.system.spec.md#responsiveness"} -->

**ID:** responsiveness\
**Sub-attribute of:** performance

The elapsed time before an operation delivers its result.

No responsiveness target has been adopted. Selecting one requires representative
workloads, defined measurement conditions, and baseline measurements.

[Edit source](<../specs/dotspec.system.spec.md#quality-attributes>)
