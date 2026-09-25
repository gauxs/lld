# Rate limiter — pluggable algorithms (design)

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

| Area | Change |
| --- | --- |
| **RateLimiter** | Already holds `RLAlgorithm`; this extension specifies registration, swap semantics, and testing strategy for multiple implementations |
| **RLAlgorithm** | Document expected side effects on `Storage` and key naming conventions per algorithm |
| **New types** | e.g. `TokenBucketAlgorithm`, `SlidingWindowLogAlgorithm`—each owns config map shape analogous to `FixedWindowAlgorithm` |

## API (additions)

| Method | Description |
| --- | --- |
| `UpdateRLAlgorithm(newAlg RLAlgorithm)` | (baseline) Atomic pointer swap under write lock; document happens-before for readers |

## Algorithm comparison (sketch)

| Algorithm | Burst at boundary | Memory | Notes |
| --- | --- | --- | --- |
| Fixed window | High at rollovers | O(keys) | Baseline |
| Sliding window log | Lower | O(limit × keys) | Store timestamps per accept |
| Token bucket | Controlled burst | O(keys) | Refill rate + bucket size |

Implement chosen algorithms under `problems/rate_limiter/extensions/pluggable_algorithms/code/` when practicing—no reference code in repo yet.
