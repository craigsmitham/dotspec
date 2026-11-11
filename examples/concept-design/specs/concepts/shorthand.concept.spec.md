---
type: concept
---
# Shorthand

## Purpose

Let people refer to a target with a short, memorable code.

## Type parameters

### Target

An opaque target identity. The concept makes no assumptions about its structure.

## Introduced types

### Code

An identity for a registered shorthand. Each use maintains its own codes.

## State

### Initial state

No codes are registered and no codes have been issued.

### Issued codes

All identities previously issued by this use, retained even after unregistration.

### Targets

A partial mapping from Code to Target. Registered codes are its domain, a subset
of Issued codes.

#### Type

Code → Target.

### Registered count

The number of currently registered codes.

#### Derivation

The cardinality of the domain of Targets; no separate counter is required.

### Invariants

- Each registered code identifies exactly one target.

## Actions

### Register

Create a code not in Issued codes for the supplied target. Do not replace an earlier registration.

#### Inputs

##### Target input

**Type:** [Target](#target)

The target to make accessible.

#### Outputs

##### Code result

**Type:** [Code](#code)

The newly registered code.

#### Effects

Add the code to Issued codes and add the code-to-target association.

### Lookup

For a registered code, return its target. Otherwise return an unknown-code result.

#### Inputs

##### Code input

**Type:** [Code](#code)

The code to look up.

#### Outputs

##### Target result

**Type:** [Target](#target)

Returned on a successful lookup. An unknown code returns no target.

### Unregister

Remove a code's association. Repeated removal succeeds without further changes.
Never reassign that code to another target in this use.

## Operational principles

### Access through a shorthand

Register a target, retain the returned code, and look it up later to recover the
same target while the registration remains active.

### End access

Unregister the code. Subsequent lookups report that it is unknown.
