---
title: Rate limiter — rejection metadata (design)
problem: rate_limiter
extension: rejection_metadata
prev:
  text: requirements
  link: /problems/rate-limiter/extensions/rejection-metadata/requirements
---

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

- **Return type** — Replace or wrap `enum.RLStatus` with `LimitResult` (names illustrative).
- **FixedWindowAlgorithm** — Reuse `currentWindowExpiry` logic when returning reject.
- **RateLimiter.Handle** — Populate metadata on reject path only.

## Go design sketch

Method bodies are intentionally omitted. Comments describe the design delta.

```go
// LimitResult replaces or wraps enum.RLStatus with accept/reject plus an
// optional retry hint. Names are illustrative.
type LimitResult struct {
    Status     enum.RLStatus
    RetryAfter time.Duration
}

// Handle returns accept/reject plus an optional retry hint and populates
// metadata on the reject path only.
func (rl *RateLimiter) Handle(req *pkg.Request) LimitResult
```

Baseline `FixedWindowAlgorithm` already computes window expiry internally—surface that value on reject instead of discarding it.

Implement in `extensions/rejection_metadata/code/` during practice.
