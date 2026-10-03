---
title: Rate limiter — rejection metadata
pageClass: lld-req-page
problem: rate_limiter
extension: rejection_metadata
next:
  text: design
  link: /problems/rate-limiter/extensions/rejection-metadata/design
---

**Builds on:** [fixed window baseline](/problems/rate-limiter/extensions/baseline/requirements).

When a request is rejected, callers need enough information to back off: time until the **current window resets** or a standard **Retry-After** hint. Accepted requests may omit timing fields.

## Functional requirements

### FR-1: Rejection payload

Reject responses include **reset time** or **retry-after duration** derived from the active window boundary (fixed window: end of current clock-aligned window).

> **Interview prompts**
>
> - **HTTP mapping?** `Retry-After` seconds or HTTP-date—caller’s problem; core returns structured duration/time.
> - **Accepted response?** No retry metadata required.

### FR-2: Stable contract

Extend status beyond bare enum where needed (e.g. result struct with `Status`, `RetryAfter`) without breaking baseline callers that only check allow/reject.

> **Interview prompts**
>
> - **Wrapper vs new method?** `HandleWithMetadata` vs enriched return type—pick one and document migration.

## Out of scope (this extension)

- None stated.

## Non-functional requirements

### NFR-1: Clock consistency

Reset time uses the same clock source as window alignment in baseline (wall clock / UTC epoch chunks).

## Extensions from here

- See [baseline](/problems/rate-limiter/extensions/baseline/requirements) for other follow-ups.
