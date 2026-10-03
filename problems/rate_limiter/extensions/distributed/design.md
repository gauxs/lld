# Rate limiter — distributed limits (design)

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

- **Storage** — Interface extracted from in-memory map.
- **FixedWindowAlgorithm** — Unchanged policy; depends on `Storage` port.
- **Key design** — Global key namespace: env prefix + resource id + window + version.

## Go design sketch

Method bodies are intentionally omitted. Comments describe the design delta.

```go
// Storage implementations provide an atomic per-key compare-and-increment
// with expiry. RedisStorage uses one server-side transaction/script.
type Storage interface {
    CompareAndIncrement(
        ctx context.Context,
        key string,
        limit int,
        ttl time.Duration,
    ) (allowed bool, err error)
}

// RateLimiter now depends on the Storage interface instead of *Storage.
type RateLimiter struct {
    mu        sync.RWMutex
    storage   Storage
    algorithm RLAlgorithm
    generator ResourceIDGenerator
    failClosed bool // zero value preserves the default fail-open behavior
}

// RLAlgorithm receives the Storage interface and propagates backend failures
// so RateLimiter can apply the documented fail-open or fail-closed policy.
type RLAlgorithm interface {
    HandleResource(
        ctx context.Context,
        resource *Resource,
        storage Storage,
    ) (enum.RLStatus, error)
}

// Handle preserves the baseline allow/reject contract. Storage errors fail
// open unless failClosed is explicitly enabled.
func (rl *RateLimiter) Handle(
    ctx context.Context,
    request *pkg.Request,
) enum.RLStatus

// MemoryStorage is the baseline process-local Storage implementation.
type MemoryStorage struct {
    counters sync.Map
}

func (s *MemoryStorage) CompareAndIncrement(
    ctx context.Context,
    key string,
    limit int,
    ttl time.Duration,
) (bool, error)

// RedisStorage is a shared Storage implementation.
type RedisStorage struct {
    client *redis.Client
}

func (s *RedisStorage) CompareAndIncrement(
    ctx context.Context,
    key string,
    limit int,
    ttl time.Duration,
) (bool, error)
```

## Deployment sketch

```text
[Instance A] ──┐
               ├──► Redis (counters + TTL)
[Instance B] ──┘
```

Implement `DistributedStorage` and wiring in `extensions/distributed/code/` during practice.

Per-key updates are atomic across application instances. The design does not
permit race-based over-limit increments; a backend outage follows the
documented fail-open policy and emits an operational signal.
