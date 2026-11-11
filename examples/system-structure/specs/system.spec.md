---
type: system
id: store
---

# Online Store

## Purpose

Let customers buy available products.

## Scope

Order placement, stock reservation, and confirmation. Payment processing is external.

## Components

### Checkout

**ID:** checkout

#### Purpose

Accept and price customer orders.

#### Part of

[Online Store](#store)

#### Responsibilities

Calculate totals and coordinate order submission.

#### Boundary

Owns order submission; stock availability remains Inventory's responsibility.

#### Contracts

[Order submission contract](../order-contract.md)

#### Architecture

[Order coordination decision](../order-coordination.md)

#### Implementation

Implementation locations can be linked here when they exist.

### Inventory

**ID:** inventory

#### Purpose

Track available stock.

#### Part of

[Online Store](#store)

### Reservations

**ID:** reservations

#### Purpose

Hold stock for orders awaiting confirmation.

#### Part of

[Inventory](#inventory)

## Features

### Place an order

**ID:** place-order

#### Purpose

Let a customer purchase products with confidence that stock is reserved.

#### Groups

- [Reserve purchased stock](requirements/reserve-stock.requirement.spec.md)

### Prevent overselling

**ID:** prevent-overselling

#### Purpose

Keep accepted orders within available stock.

#### Groups

- [Reserve purchased stock](requirements/reserve-stock.requirement.spec.md)


## Concept uses

### Stock holds

#### Definition

[Reservation](concepts/reservation.concept.spec.md)
