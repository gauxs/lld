# Rate limiter — fixed window (in-process)

Design an in-process rate limiter for a single application. Incoming API requests are counted per **client** and **API endpoint**. Each endpoint can have its own limit and time window. Start with a **fixed window** algorithm; storage is in-memory and must stay correct under concurrent requests.

## Functional requirements

<div class="lld-req">

### FR-1: Allow or reject

For each request, the limiter returns whether the request is **accepted** or **rejected** based on the current count in the active window.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Idempotent retries?** Each call consumes quota unless you add a separate idempotency extension.
- **Unknown client?** Treat as a normal client id; no special case required here.

</div>
</details>

</div>

<div class="lld-req">

### FR-2: Per client + API resource

Limits apply to a composite key (client identifier + API route). Different APIs on the same client have independent counters.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Global client cap?** Out of scope—only per-endpoint limits in this extension.
- **Key format?** Implementation detail; keep generation behind a small abstraction.

</div>
</details>

</div>

<div class="lld-req">

### FR-3: Per-API configuration

Each API can define its own **maximum request count** and **window duration** (between one second and one hour).

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Default for unconfigured API?** Define explicitly (fail closed vs allow)—document in design.
- **Invalid duration?** Reject configuration; do not silently clamp without stating rules.

</div>
</details>

</div>

<div class="lld-req">

### FR-4: Fixed window on clock boundaries

Use a fixed-window counter algorithm first. Windows align to clock time (epoch divided into window-sized chunks), not a rolling timer per first request.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Burst at window edge?** Known fixed-window artifact; sliding window is a [pluggable algorithms extension](/problems/rate-limiter/extensions/pluggable-algorithms/requirements).
- **Timezone?** Use UTC / monotonic epoch math consistently.

</div>
</details>

</div>

<div class="lld-req">

### FR-5: Runtime configuration changes

Operators may change limits or window duration at runtime. Existing counter entries do **not** need migration—new config applies via a version or key scheme so old buckets expire naturally.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Mid-window tighten?** New limit applies to new window keys; old keys TTL out.
- **Loosen limit?** Same; no backfill of rejected requests.

</div>
</details>

</div>

<div class="lld-req">

### FR-6: Request count only

Rate limiting is based on **number of requests**, not payload size or compute cost.

</div>

## Out of scope (this extension)

<div class="lld-req lld-req--scope">

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Token bucket / sliding window?** [Pluggable algorithms](/problems/rate-limiter/extensions/pluggable-algorithms/requirements).
- **Retry-after on reject?** [Rejection metadata](/problems/rate-limiter/extensions/rejection-metadata/requirements).
- **Shared counters across machines?** [Distributed limits](/problems/rate-limiter/extensions/distributed/requirements).
- **Persistence across restarts?** In-memory only here.

</div>
</details>

</div>

## Non-functional requirements

<div class="lld-req">

### NFR-1: Thread safety

Concurrent callers must not corrupt counters or observe torn updates. Correctness beats micro-optimizations.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Hot keys?** Single mutex per key or sharded locks; discuss tradeoffs.
- **Read-heavy metrics?** Not required in baseline; avoid premature caching.

</div>
</details>

</div>

<div class="lld-req">

### NFR-2: Memory and cleanup

Storage stays in-process. Stale client/API counter state should **eventually** be removed (TTL or lazy expiry) so maps do not grow without bound.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **TTL accuracy?** Eventual cleanup is enough; exact timing is best-effort.
- **Leaked goroutines?** If using timers for expiry, bound work per key.

</div>
</details>

</div>

<div class="lld-req">

### NFR-3: Extensible algorithm slot

The core `RateLimiter` should delegate counting rules to a replaceable algorithm implementation even though baseline ships fixed window only.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Swap at runtime?** Supported in reference sketch via algorithm update; full story in [pluggable algorithms](/problems/rate-limiter/extensions/pluggable-algorithms/requirements).

</div>
</details>

</div>

## Extensions from here

| Extension | What it adds |
| --- | --- |
| [Pluggable algorithms](/problems/rate-limiter/extensions/pluggable-algorithms/requirements) | Token bucket, sliding window, hot-swap without rewriting storage |
| [Rejection metadata](/problems/rate-limiter/extensions/rejection-metadata/requirements) | Retry-after / window reset time on rejected requests |
| [Distributed limits](/problems/rate-limiter/extensions/distributed/requirements) | Shared counters across processes or regions |
