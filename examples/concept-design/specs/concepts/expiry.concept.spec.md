---
type: concept
---
# Expiry

## Purpose

Identify when a resource's explicitly allocated lifetime has ended.

## Type parameters

### Resource

An opaque resource identity; expiry does not inspect or modify the resource itself.

## State

### Initial state

No resources have an allocated deadline or a completed expiration.

### Deadlines

A partial mapping from Resource to deadline. Pending expiration retains the identity
and deadline until the expiration action completes.

#### Type

Resource → instant.

### Completed resources

The identities whose expiration completed, retained to make replayed allocations
idempotent even after their deadlines have been removed.

## Actions

### Allocate

Associate the resource with a deadline after the supplied positive lifetime.
A repeated allocation for a pending resource retains its original deadline.
An allocation for a completed resource succeeds without adding another deadline.

#### Inputs

##### Resource input

**Type:** [Resource](#resource)

The identity whose lifetime is tracked.

##### Lifetime

**Primitive:** Duration

A positive elapsed duration.

#### Rejections

A nonpositive lifetime returns an error and leaves state unchanged.

### Expire

Complete expiration for a resource whose deadline has passed, remove its deadline,
add it to Completed resources, and return its identity. An already expired or unknown resource produces no new
expiration result.

#### Initiation

System-initiated when a deadline has passed. The implementation must arrange timely
execution; this specification does not assume an external user invokes expiration.

#### Guards

- The resource has an allocated deadline that has passed.

#### Outputs

##### Expired resource

**Type:** [Resource](#resource)

The identity whose lifetime ended.

## Operational principle

Allocate a resource with a lifetime. Once that time has passed, expiration reports
the resource identity so the surrounding application can react.
