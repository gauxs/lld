Authoring moved to [`problems/rate_limiter/extensions/baseline/requirements.md`](../problems/rate_limiter/extensions/baseline/requirements.md). This file is kept as a pointer only.

## Consolidated Requirements (archive)

Single-application/component scope.
API requests are rate limited.
Multiple APIs/endpoints exist.
Different APIs can have different limits.
Fixed-window algorithm initially.
Algorithm should be replaceable/extensible.
Concurrent requests are expected.
Correctness is more important than memory optimization, but memory usage should remain reasonable.
Rate limiting is based on request count.
In-memory storage is sufficient.
Window duration is configurable from 1 second to 1 hour.
Rate limiting is per client + API endpoint.
Each API can have its own request limit and window duration.
Configuration can change at runtime.
Existing counters do not need migration after configuration changes.
The API should return allow/reject and reset timing for rejected requests.
Fixed windows are aligned to clock boundaries.
Thread safety is required.
Stale client/API state should eventually be cleaned up.
