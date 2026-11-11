---
type: system
id: journey-service
---
# Journey service

## Purpose

Help commuters respond to changing travel conditions.

## Scope

Journey information and disruption notifications. Transport operations remain external.

## Components

### Notification service

**ID:** notification-service

#### Part of

[Journey service](#journey-service)

#### Responsibilities

Deliver relevant disruption information to affected commuters.

## Features

### Disruption alerts

**ID:** disruption-alerts

#### Purpose

Help a commuter discover that their travel plan needs attention.

#### Supports

[Timely awareness](jobs/reach-work.job.spec.md#timely-awareness)

#### Groups

[Notify affected commuters](requirements/notify-commuters.requirement.spec.md)

## Capabilities

### Disruption communication

**ID:** disruption-communication

#### Purpose

Make relevant changes understandable and available to affected travelers.

#### Supports

[Timely awareness](jobs/reach-work.job.spec.md#timely-awareness)
[Confidence in arrival](jobs/reach-work.job.spec.md#arrival-confidence)

#### Depends on

[Service-change detection](#change-detection)

#### Provided by

[Notification service](#notification-service)

#### Enablers

Transport staff interpret ambiguous reports and provide explanations. The software provider link is non-exhaustive.

### Service-change detection

**ID:** change-detection

#### Purpose

Recognize changes that could affect a planned journey.

#### Enablers

Transport operations staff and external information feeds. No software provider is specified here.
