# Low-Level Design (LLD): [System / Problem Name]

> **Interview Level:** Staff / Principal Engineer  
> **Focus:** Domain modeling, concurrency & state integrity, extensibility, and edge-case resilience (Non-Machine Coding).

---

## 1. Functional Requirements (FR)

### Core Capabilities (In-Scope)
* **[FR-1] Actor & Resource Lifecycle:** [e.g., Onboarding, online/offline state transitions, metadata updates.]
* **[FR-2] Primary Workflow / Execution:** [e.g., Request creation, matching/allocation, state progression.]
* **[FR-3] Termination, Cancellation & Settlement:** [e.g., Completion, user/system cancellation, post-processing/billing hooks.]
* **[FR-4] Query & Observability:** [e.g., Fetching current status, active sessions, or audit history.]

### Explicit Out-of-Scope
* [e.g., External payment gateway settlement mechanics, ML ranking model training, authentication/IAM.]

---

## 2. Non-Functional Requirements (NFR) & System Constraints

* **Concurrency & Consistency:** [e.g., Zero double-allocation under high contention; atomic state transitions.]
* **Idempotency & Fault Tolerance:** [e.g., Safe handling of duplicate client retries via idempotency keys; graceful recovery from partial failures.]
* **Latency & Performance:** [e.g., Sub-50ms in-memory candidate lookup; non-blocking I/O during external notifications.]
* **Extensibility (Open/Closed):** [e.g., Pluggable strategies for matching, pricing, and policy evaluation without modifying core orchestration.]
* **Scale Assumptions (Single-Node vs. Distributed Context):** [e.g., Clarify whether state lives in-memory on a partitioned shard or backed by an optimistic/pessimistic datastore.]

---

## 3. Core Domain Entities & Relationships

### Entities & Value Objects
* **`EntityA` (Aggregate Root):** `[id, version, state, metadata, timestamps]`
* **`EntityB`:** `[id, status, currentLocation/attributes, activeReferenceId]`
* **`Offer / AllocationToken` (Ephemeral Lease):** `[offerId, entityAId, entityBId, expiresAt, status]`
* **`Value Objects`:** `Location(lat, lng, timestamp)`, `Money(amount, currency)`, `TimeWindow`

### Key Relationships
* `EntityA` `1 --- 0..*` `Offer` (Sequential or broadcast dispatch leases)
* `EntityA` `1 --- 0..1` `EntityB` (Active assignment invariant)

---

## 4. State Machines & Invariants

> *Staff Signal: Explicitly define valid transitions, terminal states, and cross-entity invariants.*

### Primary Entity State Machine (`EntityA`)
```text
[CREATED] ──(dispatch)──> [MATCHING] ──(accepted)──> [ASSIGNED] ──(start)──> [IN_PROGRESS] ──(finish)──> [COMPLETED]
                              │                          │                        │
                              └────────(cancel/timeout)──┴────────(cancel)────────┴──> [CANCELLED]
```

### Cross-Entity Invariants
1. **Mutual Exclusion:** If `EntityB.state == BUSY`, `EntityB.activeEntityAId` must point to a non-terminal `EntityA`.
2. **Terminal Immutability:** Once an entity reaches `[COMPLETED, CANCELLED, EXPIRED]`, all subsequent state-mutating calls must fail fast or return idempotently.

---

## 5. Key Interfaces & Design Patterns (Extensibility)

* **Strategy Pattern — `[PolicyName]Strategy`:**
  * *Why:* Decouples volatile business rules (e.g., ranking, pricing, cancellation penalties) from the core workflow engine.
  * *Signature:* `Result evaluate(Context ctx)`
* **State Pattern / Explicit Transition Table:**
  * *Why:* Prevents scattered `if/else` status checks and enforces legal state transitions centrally.
* **Observer / Domain Events (`EventPublisher`):**
  * *Why:* Decouples core state commits from side effects (push notifications, analytics, ledger entries).
* **Repository / Store Abstraction:**
  * *Why:* Separates domain logic from storage concurrency mechanisms (e.g., `compareAndSet`, row locks, or spatial indexes).

---

## 6. Concurrency, Locking & Lease Strategy

* **Granularity of Locking:** [e.g., Fine-grained per-entity lock / striped lock / optimistic CAS (`version`) rather than coarse global locks.]
* **Lock Ordering (Deadlock Prevention):** [e.g., Always acquire lock on `Trip` before `Driver`, or use single-entity CAS + compensating transitions.]
* **Handling Slow External Calls:** [e.g., Never hold in-memory locks while awaiting network/push notifications; reserve via time-bounded lease (`expiresAt`) first, release lock, then notify.]
* **Idempotency Mechanism:** [e.g., Deduplication store keyed by `(actorId, idempotencyKey)` storing request hash + serialized outcome.]

---

## 7. Deep Dive: Hard Core & Edge Cases

| # | Edge Case / Failure Mode | Impact if Unhandled | Resolution Strategy |
| :--- | :--- | :--- | :--- |
| **1** | **Concurrent Claims (Race Condition)** | Double-booking / split-brain assignment | Atomic CAS on resource state + lease token validation |
| **2** | **Lease Expiry vs. Late Client Ack** | Ghost assignment or stranded resource | Fencing tokens / epoch check; reject ack if `now > expiresAt` |
| **3** | **Simultaneous Cancel & Accept** | Trip cancelled while resource marked busy | Serialize on Aggregate Root lock/version; deterministic tie-breaker |
| **4** | **Stale / Out-of-Order Telemetry** | Matching based on outdated state/location | Monotonic timestamps / sequence numbers; TTL eviction on index |
| **5** | **Cascading Rejections / Partial Failure** | Request starvation or leaked reservations | Bounded retry loop with backoff/radius expansion + background sweeper |

---

## 8. Sequence Walkthrough (Happy Path + Race Path)

### A. Happy Path (`Request -> Reserve -> Accept`)
1. Client sends `createRequest(idempotencyKey, payload)`.
2. Orchestrator validates & persists `EntityA` in `MATCHING`.
3. Spatial/Candidate Index returns top-$K$ eligible candidates.
4. Dispatcher atomically transitions Candidate #1 to `RESERVED` with `offerId` & `expiresAt`.
5. Dispatcher emits async `OfferSentEvent` (outside lock).
6. Candidate #1 calls `acceptOffer(offerId)` before `expiresAt` $\rightarrow$ atomic commit of `ASSIGNED`.

### B. Contention / Race Path (`Cancel` vs. `Accept`)
1. Thread 1 (`cancel`) acquires lock/CAS on `EntityA`, transitions to `CANCELLED`, releases `EntityB` back to `AVAILABLE`.
2. Thread 2 (`accept`) subsequently checks `EntityA.state` / `offer.status`, sees `CANCELLED`, and returns a clean `409 Conflict / OfferExpired` response.

---

## 9. Trade-offs & Future Extensions

* **Trade-off 1:** [e.g., Sequential offer dispatch (lower contention, higher latency) vs. Broadcast offer dispatch (lower latency, requires multi-winner arbitration).]
* **Trade-off 2:** [e.g., Pessimistic fine-grained locks vs. Optimistic Concurrency Control (CAS) under high contention spikes.]
* **Future Extension:** [e.g., How adding pooled/shared rides or multi-leg workflows impacts the `1:1` assignment invariant.]