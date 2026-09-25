---
title: Rate limiter — fixed window (design)
problem: rate_limiter
extension: baseline
prev:
  text: requirements
  link: /problems/rate-limiter/extensions/baseline/requirements
next:
  text: codebase
  link: /problems/rate-limiter/extensions/baseline/codebase
---

Reference implementation: [`problems/rate_limiter/extensions/baseline/code/`](https://github.com/gauxs/lld/tree/main/problems/rate_limiter/extensions/baseline/code) (package `ratelimiter`).

## Entities & responsibilities

| Entity | Responsibilities |
| --- | --- |
| **RateLimiter** | Facade: resolve resource from request, delegate to active `RLAlgorithm`, optional algorithm swap under lock |
| **ResourceIDGenerator** | Map `Request` (client id + API) → `Resource` id string |
| **Resource** | Opaque limit key identity |
| **RLAlgorithm** | Policy: given resource + storage, return accept/reject |
| **FixedWindowAlgorithm** | Per-resource config (max count, window duration, version), clock-aligned window index, storage keys |
| **Storage** | In-memory map of counter entries; compare-and-increment under limit; schedule key deletion after window expiry |

## API

### RateLimiter

| Method | Description |
| --- | --- |
| `Handle(req *Request) RLStatus` | Accept or reject for this request |
| `UpdateRLAlgorithm(newAlg RLAlgorithm)` | Point new traffic at a different algorithm; old buckets age out |

### Request (`pkg`)

| Field / accessor | Description |
| --- | --- |
| `clientID`, `api` | Inputs to resource key generation |

### FixedWindowAlgorithm

| Method | Description |
| --- | --- |
| `HandleResource(res *Resource, s *Storage) RLStatus` | Increment counter for current window or reject |
| `UpdateWindowDuration(res *Resource, d time.Duration)` | Bump config version; new windows use new duration |
| `UpdateMaxResourceCount(res *Resource, n int)` | Bump config version; new windows use new cap |

### Storage

| Method | Description |
| --- | --- |
| `CompareAndIncrement(key string, lessThan int, expiry time.Duration) bool` | Create or increment counter if below cap; schedule cleanup |
| `NewStorage() *Storage` | Empty backing map |

## Concurrency notes

- **RateLimiter** uses `RWMutex` around algorithm pointer reads and exclusive lock on swap.
- **FixedWindowAlgorithm** holds per-resource config in `sync.Map`; each `resourceConfig` has its own lock for version and limits.
- **Storage** uses `sync.Map` of per-key mutexes for increment; first insert schedules `DeleteAfterDuration`.

## Configuration versioning

Storage keys embed algorithm id, **config version**, resource id, and **window number** so runtime config changes do not require migrating old counters.
