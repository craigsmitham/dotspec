---
type: feature-spec
specifications:
  - name: Calendar dates accepted for follow-up
    defines: requirements
    source: ../../app/follow_up_date_test.go
---

# Follow up when someone owes me a reply

## Description

Find an errand after a break, see whether a reply needs chasing and take up the work again when an answer arrives.

The [calendar date suite](#specification.calendar-dates-accepted-for-follow-up) defines which dates a person can enter.

## Purpose

Someone can finish making calls for the day while still owing several results. When they return, they can see which calls need chasing and what answer they were waiting for. [Follow-up review](../capabilities/follow-up-review.capability.spec.md) helps them choose whom to contact first.

## Groups

- [Decide what to do about a waiting errand](../use-cases/resume-list.use-case.spec.md) shows how the owner decides whether to take up the work, keep waiting or mark the errand settled.
- [Show which replies need follow-up](../requirements/follow-up.requirement.spec.md) specifies the cue for deciding where to act.
- [Keep recorded work available](../requirements/availability.requirement.spec.md) lets the owner find each errand in the view for its recorded state when they come back.
- [Preserve tasks across restart](../requirements/restart.requirement.spec.md) covers returning after the service itself has stopped.
