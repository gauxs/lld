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

### FR-1: Shared storage backend

Replace or adapt `Storage` so increments are visible across processes (e.g. Redis, DynamoDB, or a dedicated rate-limit service).

> **Interview prompts**
>
> - **Strong vs eventual consistency?** Over-limit by one under race may be acceptable—state assumption.
> - **Hot keys?** Shard or local coalescing with global sync—advanced topic.

### FR-2: Same external behavior

Per client + API limits, configurable windows, and allow/reject semantics match baseline from the caller’s perspective.

### FR-3: Failure modes

When the shared store is unavailable, default to **fail open** so an
infrastructure outage does not block all application traffic. Emit an
operational signal so the degraded enforcement is visible.

## Out of scope (this extension)

- None stated.

## Non-functional requirements

### NFR-1: Latency budget

Each `Handle` adds at most one round trip to shared storage on the happy path unless batching is explicitly designed.

## Extensions from here

- See [baseline](/problems/rate-limiter/extensions/baseline/requirements) for algorithm and metadata extensions; combine with [**pluggable algorithms**](/problems/rate-limiter/extensions/pluggable-algorithms/requirements) when policies differ per route.
