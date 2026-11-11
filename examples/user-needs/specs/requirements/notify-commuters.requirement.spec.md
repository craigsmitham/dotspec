---
type: requirement
---
# Notify affected commuters

## Statement

The notification service must submit a disruption notification for each affected subscribed commuter within 30 seconds of receiving the disruption event.

## Supports

[Timely awareness](../jobs/reach-work.job.spec.md#timely-awareness)

## Applies to

**Targets:** [Notification service](../system.spec.md#notification-service)

## Verification plan

In a controlled test, record event receipt and notification submission timestamps for each affected subscription. Check that each elapsed duration is at most 30 seconds. This tests submission, not delivery, reading, or improved user confidence.
