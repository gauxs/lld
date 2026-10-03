---
title: Rate limiter — distributed limits (design)
problem: rate_limiter
extension: distributed
prev:
  text: requirements
  link: /problems/rate-limiter/extensions/distributed/requirements
---

**Builds on:** [baseline design](/problems/rate-limiter/extensions/baseline/design).

## Delta from baseline

- **Storage** — Interface extracted from in-memory map; implementations: `MemoryStorage` (baseline), `RedisStorage`, etc.
- **FixedWindowAlgorithm** — Unchanged policy; depends on `Storage` port
- **Key design** — Global key namespace: env prefix + resource id + window + version

## Storage port (sketch)

```text
CompareAndIncrement(ctx, key, limit, ttl) (allowed bool, err error)
  Atomic or best-effort increment with expiry
```

## Deployment sketch

```text
[Instance A] ──┐
               ├──► Redis (counters + TTL)
[Instance B] ──┘
```

Implement `DistributedStorage` and wiring in `extensions/distributed/code/` during practice.
