---
type: concept
id: reservation
---

# Reservation

## Purpose

Hold a limited resource for one claimant while a decision is pending.

## Operational principles

A claimant holds available stock, then commits or releases its hold. A release makes the stock available again.


## State

Available quantities and active holds, each identifying a claimant and quantity.

## Actions

### Hold

When sufficient stock is available, create a hold and reduce availability.

### Commit

Finalize an existing hold for its claimant.

### Release

Remove an existing hold and restore its quantity to availability.
