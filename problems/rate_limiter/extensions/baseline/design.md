# Rate limiter — fixed window (design)

Reference implementation: [`problems/rate_limiter/extensions/baseline/code/`](https://github.com/gauxs/lld/tree/main/problems/rate_limiter/extensions/baseline/code) (package `ratelimiter`).

## Approach

`RateLimiter` is a facade that resolves a resource from each request and delegates to the active `RLAlgorithm`. The fixed-window policy keeps versioned per-resource configuration, while in-memory storage atomically compares and increments counters and removes keys after window expiry.

## Go design sketch

Method bodies are intentionally omitted. Comments describe ownership,
responsibilities, and invariants.

```go
// RateLimiter resolves the resource from a request, delegates to the active
// RLAlgorithm, and supports an optional algorithm swap under lock.
type RateLimiter struct {
    mu        sync.RWMutex
    storage   *Storage
    algorithm RLAlgorithm
    generator ResourceIDGenerator
}

// ResourceIDGenerator maps a Request (client id + API) to a Resource id string.
type ResourceIDGenerator interface {
    GetResource(req *pkg.Request) *Resource
}

// Resource is an opaque limit key identity.
type Resource struct {
    id string
}

// RLAlgorithm is the policy that returns accept or reject for a resource and
// storage.
type RLAlgorithm interface {
    HandleResource(res *Resource, storage *Storage) enum.RLStatus
}

// FixedWindowAlgorithm owns per-resource configuration: maximum count, window
// duration, and version. It derives clock-aligned window indexes and storage keys.
type FixedWindowAlgorithm struct {
    config sync.Map
}

// Storage owns an in-memory map of counter entries. Compare-and-increment runs
// under the limit and key deletion is scheduled after window expiry.
type Storage struct {
    counters sync.Map
}

// Handle accepts or rejects this request.
func (rl *RateLimiter) Handle(req *pkg.Request) enum.RLStatus

// UpdateRLAlgorithm points new traffic at a different algorithm; old buckets age out.
func (rl *RateLimiter) UpdateRLAlgorithm(newAlg RLAlgorithm)

// HandleResource increments the counter for the current window or rejects.
func (fwa *FixedWindowAlgorithm) HandleResource(res *Resource, storage *Storage) enum.RLStatus

// UpdateWindowDuration bumps the config version; new windows use the new duration.
func (fwa *FixedWindowAlgorithm) UpdateWindowDuration(res *Resource, duration time.Duration)

// UpdateMaxResourceCount bumps the config version; new windows use the new cap.
func (fwa *FixedWindowAlgorithm) UpdateMaxResourceCount(res *Resource, count int)

// CompareAndIncrement creates or increments a counter if below the cap and
// schedules cleanup.
func (s *Storage) CompareAndIncrement(key string, lessThan int, expiry time.Duration) bool

// NewStorage returns storage with an empty backing map.
func NewStorage() *Storage
```

## Concurrency notes

- **RateLimiter** uses `RWMutex` around algorithm pointer reads and exclusive lock on swap.
- **FixedWindowAlgorithm** holds per-resource config in `sync.Map`; each `resourceConfig` has its own lock for version and limits.
- **Storage** uses `sync.Map` of per-key mutexes for increment; first insert schedules `DeleteAfterDuration`.

## Configuration versioning

Storage keys embed algorithm id, **config version**, resource id, and **window number** so runtime config changes do not require migrating old counters.
