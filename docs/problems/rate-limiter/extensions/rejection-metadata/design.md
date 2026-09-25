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

| Area | Change |
| --- | --- |
| **Return type** | Replace or wrap `enum.RLStatus` with `LimitResult { Status, RetryAfter time.Duration }` (names illustrative) |
| **FixedWindowAlgorithm** | Reuse `currentWindowExpiry` logic when returning reject |
| **RateLimiter.Handle** | Populate metadata on reject path only |

## API (additions)

| Method | Description |
| --- | --- |
| `Handle(req *Request) LimitResult` | Accept/reject plus optional retry hint |

Baseline `FixedWindowAlgorithm` already computes window expiry internally—surface that value on reject instead of discarding it.

Implement in `extensions/rejection_metadata/code/` during practice.
