---
type: system
---
# Expiring links

## Purpose

Provide temporary public and internal links with separately configured lifetimes.

## Scope

An illustrative concept composition, not an implemented dotspec service. Public and
internal registrations are independent uses of the same reusable definitions.

## Concept uses

### PublicLinks

Named shorthand registrations for public links.

#### Definition

[Shorthand](concepts/shorthand.concept.spec.md).

#### Bindings

##### Public targets

**Parameter:** [Target](concepts/shorthand.concept.spec.md#target)
**Primitive:** URL

### PublicExpiry

Tracks lifetimes for codes registered by PublicLinks.

#### Definition

[Expiry](concepts/expiry.concept.spec.md).

#### Bindings

##### Public resources

**Parameter:** [Resource](concepts/expiry.concept.spec.md#resource)
**Type:** [Code](concepts/shorthand.concept.spec.md#code)
**Use:** [PublicLinks](#publiclinks)

#### Dependencies

[PublicLinks](#publiclinks) supplies the identities whose lifetimes justify
this use. This is an application-scoping dependency, not a call from Expiry.

#### Configuration

The application allocates each registration a lifetime of 1 hour.

### InternalLinks

Named shorthand registrations for internal links.

#### Definition

[Shorthand](concepts/shorthand.concept.spec.md).

#### Bindings

##### Internal targets

**Parameter:** [Target](concepts/shorthand.concept.spec.md#target)
**Primitive:** URL

### InternalExpiry

Tracks lifetimes for codes registered by InternalLinks.

#### Definition

[Expiry](concepts/expiry.concept.spec.md).

#### Bindings

##### Internal resources

**Parameter:** [Resource](concepts/expiry.concept.spec.md#resource)
**Type:** [Code](concepts/shorthand.concept.spec.md#code)
**Use:** [InternalLinks](#internallinks)

#### Dependencies

[InternalLinks](#internallinks) supplies the identities whose lifetimes justify
this use. This is an application-scoping dependency, not a call from Expiry.

#### Configuration

The application allocates each registration a lifetime of 1 day.

## Synchronizations

### Start public lifetime

#### When

##### Completed action

**Use:** [PublicLinks](#publiclinks)
**Action:** [Register](concepts/shorthand.concept.spec.md#register)

###### Result

A successful completion and its returned identity, correlated to this named use.

#### Where

Only successful completions participate. An error is not a successful result.

#### Then

##### Invoke action

**Use:** [PublicExpiry](#publicexpiry)
**Action:** [Allocate](concepts/expiry.concept.spec.md#allocate)

###### Arguments

Pass the returned Code as Resource and 1 hour as Lifetime.

#### Semantics

Event-based coordination, without an atomic transaction across concepts. Retain
successful completion events until the resulting invocation succeeds; retry failures
with the same arguments. Delivery may repeat, so the receiving action is idempotent.
Do not treat emission of an event as completion of the receiving action. Links may
remain accessible while expiration delivery is delayed; no instantaneous revocation
promise is made. Each use has independent state even when definitions or code values
coincide. These are requirements for an implementation, not guarantees from dotspec.

### Remove expired public link

#### When

##### Completed action

**Use:** [PublicExpiry](#publicexpiry)
**Action:** [Expire](concepts/expiry.concept.spec.md#expire)

###### Result

A successful completion and its returned identity, correlated to this named use.

#### Where

Only successful completions participate. An error is not a successful result.

#### Then

##### Invoke action

**Use:** [PublicLinks](#publiclinks)
**Action:** [Unregister](concepts/shorthand.concept.spec.md#unregister)

###### Arguments

Pass the expired Resource identity as the Code to unregister.

#### Semantics

Event-based coordination, without an atomic transaction across concepts. Retain
successful completion events until the resulting invocation succeeds; retry failures
with the same arguments. Delivery may repeat, so the receiving action is idempotent.
Do not treat emission of an event as completion of the receiving action. Links may
remain accessible while expiration delivery is delayed; no instantaneous revocation
promise is made. Each use has independent state even when definitions or code values
coincide. These are requirements for an implementation, not guarantees from dotspec.

### Start internal lifetime

#### When

##### Completed action

**Use:** [InternalLinks](#internallinks)
**Action:** [Register](concepts/shorthand.concept.spec.md#register)

###### Result

A successful completion and its returned identity, correlated to this named use.

#### Where

Only successful completions participate. An error is not a successful result.

#### Then

##### Invoke action

**Use:** [InternalExpiry](#internalexpiry)
**Action:** [Allocate](concepts/expiry.concept.spec.md#allocate)

###### Arguments

Pass the returned Code as Resource and 1 day as Lifetime.

#### Semantics

Event-based coordination, without an atomic transaction across concepts. Retain
successful completion events until the resulting invocation succeeds; retry failures
with the same arguments. Delivery may repeat, so the receiving action is idempotent.
Do not treat emission of an event as completion of the receiving action. Links may
remain accessible while expiration delivery is delayed; no instantaneous revocation
promise is made. Each use has independent state even when definitions or code values
coincide. These are requirements for an implementation, not guarantees from dotspec.

### Remove expired internal link

#### When

##### Completed action

**Use:** [InternalExpiry](#internalexpiry)
**Action:** [Expire](concepts/expiry.concept.spec.md#expire)

###### Result

A successful completion and its returned identity, correlated to this named use.

#### Where

Only successful completions participate. An error is not a successful result.

#### Then

##### Invoke action

**Use:** [InternalLinks](#internallinks)
**Action:** [Unregister](concepts/shorthand.concept.spec.md#unregister)

###### Arguments

Pass the expired Resource identity as the Code to unregister.

#### Semantics

Event-based coordination, without an atomic transaction across concepts. Retain
successful completion events until the resulting invocation succeeds; retry failures
with the same arguments. Delivery may repeat, so the receiving action is idempotent.
Do not treat emission of an event as completion of the receiving action. Links may
remain accessible while expiration delivery is delayed; no instantaneous revocation
promise is made. Each use has independent state even when definitions or code values
coincide. These are requirements for an implementation, not guarantees from dotspec.
