---
title: Rate limiter — pluggable algorithms (design)
problem: rate_limiter
extension: pluggable_algorithms
prev:
  text: requirements
  link: /problems/rate-limiter/extensions/pluggable-algorithms/requirements
---

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

- **RateLimiter** — Already holds `RLAlgorithm`; this extension specifies registration, swap semantics, and testing strategy for multiple implementations
- **RLAlgorithm** — Document expected side effects on `Storage` and key naming conventions per algorithm
- **New types** — e.g. `TokenBucketAlgorithm`, `SlidingWindowLogAlgorithm`—each owns config map shape analogous to `FixedWindowAlgorithm`

## API (additions)

```text
UpdateRLAlgorithm(newAlg RLAlgorithm)
  (baseline) Atomic pointer swap under write lock; document happens-before for readers
```

## Algorithm comparison (sketch)

- **Fixed window**
  - Burst at boundary: High at rollovers
  - Memory: O(keys)
  - Notes: Baseline
- **Sliding window log**
  - Burst at boundary: Lower
  - Memory: O(limit × keys)
  - Notes: Store timestamps per accept
- **Token bucket**
  - Burst at boundary: Controlled burst
  - Memory: O(keys)
  - Notes: Refill rate + bucket size

Implement chosen algorithms under `problems/rate_limiter/extensions/pluggable_algorithms/code/` when practicing—no reference code in repo yet.
