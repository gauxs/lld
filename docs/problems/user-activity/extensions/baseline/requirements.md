---
title: User activity tracker
pageClass: lld-req-page
problem: user_activity
extension: baseline
next:
  text: design
  link: /problems/user-activity/extensions/baseline/design
---

Design and implement a User Activity Rate Tracker. At a high level, the system should track user activity and determine activity rates.

## Functional requirements

### FR-1: Track activity of user

Track activity of user where activities can be: LOGIN / PURCHASE / API_CALL / POST

> **What are the activties which we need to track?** Start with LOGIN but assume the type of activities can be added in future.

### FR-2: Querying activity of a user

The activity of a user will be queried as:

- Count occurrences for (user, activity, time window).
- Calculate the rate for (user, activity, time window)

> **define rate and exactly how the stored activity historical record of user be queried?** number of occurrences of an activity within a time window.

### FR-3: Other access pattern

Given an activity and a time window, the system should support:

- Count of distinct users who performed the activity in the window.
- List of users who performed the activity in the window.
- Top K users by activity count

> **will there be reverse lookup? i.e. given the activity, tell how many users or which users performed this activity in last X minute?** Yes.

### FR-4: Activity summary for a user

Given a user and time window, return the activity counts across all activity types.

## Out of scope

- **<>** <>

## Constraints

### NFR-1: Max users and activity numbers

Up to 10 million users. Up to 100 predefined activity types.

> **What are the max number of users we need to track?**

### NFR-2: Solution to be in-memory

For this LLD round, assume the tracker is a single-process, in-memory component.

- No Redis, database, Kafka, or other external dependency.
- Persistence across process restarts is not required.

> **Can we assume that solution has to be in-memeory or we need to user external dependencies like redis cache etc?**

### NFR-3: Concurrency

Yes. Assume the tracker is accessed concurrently:

- Multiple users can record activities at the same time.
- Multiple queries can execute concurrently with activity recording.
- The component must be thread-safe.

> **I am assuming there can be concurrent activities by multliple users?**

### NFR-4: Window size and dynamic behaviour

The supported historical window is up to 24 hours. The query window is dynamic per request. A caller can query 1 minute, 10 minutes, 1 hour, or any interval up to 24 hours.

> **How far back can the time window be?**
>
> **can the window size change dynamically? And what to do with the data if the window size decreases, should we discard it?**
