---
type: system
---

# dotspec

## Purpose

Help teams manage software product requirements alongside the code
they describe.

## Scope

The system includes the `dotspec` CLI and the `*.spec.md` files it manages.

## Quality attributes

This is dotspec's local vocabulary. The definitions identify concerns;
individual quality requirements establish the obligations.

### Usability

**ID:** usability

How effectively contributors can author, inspect, and maintain specifications.

### Diagnosability

**ID:** diagnosability
**Broader:** usability

How effectively findings help contributors understand and locate specification
problems.

### Accessibility

**ID:** accessibility
**Broader:** usability

How effectively people with differing abilities and ways of accessing output
can use dotspec.

### Reliability

**ID:** reliability

How consistently dotspec performs its intended operations under stated
conditions.

### Repeatability

**ID:** repeatability
**Broader:** reliability

Whether equivalent inspections under unchanged conditions produce equivalent
reports.

### Source preservation

**ID:** source-preservation
**Broader:** reliability

Whether inspection preserves the source material it examines.

### Performance

**ID:** performance

The time and computing resources required to perform operations under stated
workloads.

### Responsiveness

**ID:** responsiveness
**Broader:** performance

The elapsed time before an operation delivers its result.

No responsiveness target has been adopted. Selecting one requires representative
workloads, defined measurement conditions, and baseline measurements.
