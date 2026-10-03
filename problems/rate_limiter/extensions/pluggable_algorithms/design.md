# Rate limiter — pluggable algorithms (design)

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

- **RateLimiter** — Already holds `RLAlgorithm`; this extension specifies registration, swap semantics, and testing strategy for multiple implementations.
- **RLAlgorithm** — Document expected side effects on `Storage` and key naming conventions per algorithm.

## Go design sketch

Method bodies are intentionally omitted. Comments describe the design delta.

```go
// New types—e.g. TokenBucketAlgorithm and SlidingWindowLogAlgorithm—each own a
// config map shape analogous to FixedWindowAlgorithm.
type TokenBucketAlgorithm struct {
    config sync.Map
}

func (a *TokenBucketAlgorithm) HandleResource(
    resource *Resource,
    storage *Storage,
) enum.RLStatus

// SlidingWindowLogAlgorithm owns a config map shape analogous to
// FixedWindowAlgorithm.
type SlidingWindowLogAlgorithm struct {
    config sync.Map
}

func (a *SlidingWindowLogAlgorithm) HandleResource(
    resource *Resource,
    storage *Storage,
) enum.RLStatus

// UpdateRLAlgorithm is the baseline atomic pointer swap under a write lock;
// document happens-before for readers.
func (rl *RateLimiter) UpdateRLAlgorithm(newAlg RLAlgorithm)
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
