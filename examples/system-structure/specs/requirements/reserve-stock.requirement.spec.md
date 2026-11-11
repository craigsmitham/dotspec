---
type: requirement
id: reserve-stock
---

# Reserve purchased stock

## Specification system

[Online Store](../system.spec.md#store)

## Statement

An accepted order must reserve its purchased quantities before confirmation.

## Applies to

**Targets:** [Reservations](../system.spec.md#reservations)

## Rationale

Prevent confirmation of orders that cannot be fulfilled from available stock.

## Derived from

[Order submission contract](../../order-contract.md#order-submission-contract)

## Verification plan

Submit concurrent orders against limited stock and verify that accepted quantities never exceed stock.
