---
title: Scarcity
description: Limited resources—pools, semaphores, and rate limits.
prev:
  text: coordination
  link: /learn/concurrency/coordination
next:
  text: baseline
  link: /problems/connect-four/extensions/baseline/requirements
---

# Scarcity

<p class="lead">Scarcity problems appear when only N concurrent operations are allowed—DB connections, API quota, memory-heavy workers.</p>

## Failure modes

- Unbounded goroutines exhaust memory or file descriptors.
- Fairness: one client starves others without per-tenant limits.
- Leaked permits after panics (semaphore never released).

## What to reach for

- **Semaphore** (or weighted semaphore in Go) to cap concurrency.
- **Object pool** with acquire/release and timeouts.
- **Rate limiter** keyed by client + route (see `rate_limiter/` in this repo).

Always release permits in `defer` or `finally` semantics.

```go
sem := make(chan struct{}, maxConcurrent)
sem <- struct{}{} // acquire
defer func() { <-sem }() // release
```

## What's next

Apply the requirement-first trail on a full problem: [Connect Four baseline](/problems/connect-four/extensions/baseline/requirements), then open child extensions from there.
