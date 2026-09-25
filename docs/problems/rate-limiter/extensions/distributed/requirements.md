---
title: Rate limiter — distributed limits
pageClass: lld-req-page
problem: rate_limiter
extension: distributed
next:
  text: design
  link: /problems/rate-limiter/extensions/distributed/design
---

**Builds on:** [fixed window baseline](/problems/rate-limiter/extensions/baseline/requirements).

Multiple application instances must enforce the **same** limits for a given client + API. In-memory `sync.Map` is insufficient; counters live in a shared store with acceptable latency and correctness tradeoffs.

## Functional requirements

<div class="lld-req">

### FR-1: Shared storage backend

Replace or adapt `Storage` so increments are visible across processes (e.g. Redis, DynamoDB, or a dedicated rate-limit service).

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Strong vs eventual consistency?** Over-limit by one under race may be acceptable—state assumption.
- **Hot keys?** Shard or local coalescing with global sync—advanced topic.

</div>
</details>

</div>

<div class="lld-req">

### FR-2: Same external behavior

Per client + API limits, configurable windows, and allow/reject semantics match baseline from the caller’s perspective.

</div>

<div class="lld-req">

### FR-3: Failure modes

Define behavior when the store is unavailable (fail open vs fail closed) and document operational choice.

</div>

## Non-functional requirements

<div class="lld-req">

### NFR-1: Latency budget

Each `Handle` adds at most one round trip to shared storage on the happy path unless batching is explicitly designed.

</div>

## Extensions from here

See baseline for algorithm and metadata extensions; combine with **pluggable algorithms** when policies differ per route.
