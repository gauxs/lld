# Rate limiter — pluggable algorithms

**Builds on:** [fixed window baseline](/problems/rate-limiter/extensions/baseline/requirements).

Same request model and per client + API keys, but counting policy is swappable: fixed window remains the default, with room for **token bucket**, **sliding window log**, or **leaky bucket** without rewriting the facade or storage contract.

## Functional requirements

### FR-1: Algorithm interface

All policies implement a shared interface invoked by `RateLimiter` with the resolved `Resource` and shared `Storage` (or a narrowed storage port if an algorithm needs different primitives).

> **Interview prompts**
>
> - **Storage leakage?** Prefer algorithm-specific keys or a small `Storage` interface per family.
> - **Unit test algorithms?** Mock storage; no network.

### FR-2: Runtime swap

Operators or code can replace the active algorithm instance. Requests after the swap use the new policy; old algorithm state expires via existing TTL/version keys—no migration job.

> **Interview prompts**
>
> - **In-flight request?** Define whether a single `Handle` sees old or new algorithm for its full call—typically snapshot pointer at start.
> - **Per-API algorithm?** Optional stretch; baseline uses one active algorithm globally.

### FR-3: Additional algorithms (candidates)

Document at least two alternatives beyond fixed window and when each is appropriate (burst tolerance, smoothness, memory).

## Out of scope (this extension)

- None stated.

## Non-functional requirements

### NFR-1: Backward compatibility

Fixed-window behavior from baseline must remain available as one registered implementation.

## Extensions from here

- No further siblings defined from this extension; see baseline for [**rejection metadata**](/problems/rate-limiter/extensions/rejection-metadata/requirements) and [**distributed limits**](/problems/rate-limiter/extensions/distributed/requirements).
