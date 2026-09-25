# Rate limiter — distributed limits (design)

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

| Area | Change |
| --- | --- |
| **Storage** | Interface extracted from in-memory map; implementations: `MemoryStorage` (baseline), `RedisStorage`, etc. |
| **FixedWindowAlgorithm** | Unchanged policy; depends on `Storage` port |
| **Key design** | Global key namespace: env prefix + resource id + window + version |

## Storage port (sketch)

| Method | Description |
| --- | --- |
| `CompareAndIncrement(ctx, key, limit, ttl) (allowed bool, err error)` | Atomic or best-effort increment with expiry |

## Deployment sketch

```text
[Instance A] ──┐
               ├──► Redis (counters + TTL)
[Instance B] ──┘
```

Implement `DistributedStorage` and wiring in `extensions/distributed/code/` during practice.
