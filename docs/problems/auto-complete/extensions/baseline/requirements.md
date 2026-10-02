---
title: User activity tracker
pageClass: lld-req-page
problem: auto_complete
extension: baseline
prev:
  text: scarcity
  link: /learn/concurrency/02-problems/scarcity
next:
  text: design
  link: /problems/auto-complete/extensions/baseline/design
---

Design and implement a Search Autocomplete System


## Functional requirements
<div class="lld-req">

### FR-1: Support inserting words into an internal dictionary.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

<!-- - **What are the activties which we need to track?** Start with LOGIN but assume the type of activities can be added in future. -->

</div>
</details>

### FR-2: Return suggestions when a user types a prefix.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **define rate and exactly how the stored activity historical record of user be queried?** number of occurrences of an activity within a time window.

</div>
</details>


### FR-3: Other access pattern

Given an activity and a time window, the system should support:
- Count of distinct users who performed the activity in the window.
- List of users who performed the activity in the window.
- Top K users by activity count

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **will there be reverse lookup? i.e. given the activity, tell how many users or which users performed this activity in last X minute?** Yes.

</div>
</details>

### FR-4: Activity summary for a user

Given a user and time window, return the activity counts across all activity types.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

</div>
</details>

</div>

## Out of scope (this extension)
<div class="lld-req lld-req--scope">

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **<>** <>

</div>
</details>

</div>

## Non-functional requirements
<div class="lld-req">

### NFR-1: Max users and activity numbers

Up to 10 million users. Up to 100 predefined activity types.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **What are the max number of users we need to track?** 

</div>
</details>

### NFR-2: Solution to be in-memory

For this LLD round, assume the tracker is a single-process, in-memory component.
 - No Redis, database, Kafka, or other external dependency.
 - Persistence across process restarts is not required.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Can we assume that solution has to be in-memeory or we need to user external dependencies like redis cache etc?** 

</div>
</details>

### NFR-3: Concurrency

Yes. Assume the tracker is accessed concurrently:
- Multiple users can record activities at the same time.
- Multiple queries can execute concurrently with activity recording.
- The component must be thread-safe.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **I am assuming there can be concurrent activities by multliple users?** 

</div>
</details>

### NFR-4: Window size and dynamic behaviour

The supported historical window is up to 24 hours. The query window is dynamic per request. A caller can query 1 minute, 10 minutes, 1 hour, or any interval up to 24 hours.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **How far back can the time window be?** 
- **can the window size change dynamically? And what to do with the data if the window size decreases, should we discard it?**

</div>
</details>

</div>