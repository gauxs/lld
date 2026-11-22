# Machine Coding: [System / Problem Name]

> **Interview Format:** Machine Coding (90–120 mins)  
> **Focus:** Working end-to-end implementation, clean OOP/SOLID boundaries, thread-safe in-memory state, and verifiable edge-case flows.

---

## 1. Problem Scope & P0 / P1 Breakdown

### P0: Must-Have Core Flows (Build First)
- **[P0-1] Entity & State Setup:** [e.g., Register users/resources, toggle availability/status, validate core inputs.]
- **[P0-2] Primary End-to-End Flow:** [e.g., Create request -> match/allocate -> progress state -> complete.]
- **[P0-3] Basic Query / Status View:** [e.g., Inspect current state, active allocations, or final result/fare.]

### P1: Stretch / Bonus Flows (Build Only After P0 Works)
- **[P1-1] Cancellation & Re-allocation:** [e.g., Cancel flow, fallback matching, or penalty rules.]
- **[P1-2] Pluggable Rule Variants:** [e.g., Alternative matching/pricing/ranking strategies.]
- **[P1-3] Concurrency / Stress Scenario:** [e.g., Multi-threaded contention on shared resources.]

### Explicit Non-Goals (For Machine Coding)
- No external databases, message brokers, HTTP frameworks, or over-engineered boilerplate unless requested.

---

## 2. Package & File Structure

```text
src/
├── models/          # Entities, Value Objects, Status Enums
├── exceptions/      # Domain-specific exceptions (invalid state, conflict, not found)
├── repositories/    # Thread-safe in-memory stores & lookup helpers
├── strategies/      # Interfaces & concrete implementations (matching, pricing, rules)
├── services/        # Core orchestration, state guards, concurrency control
└── Main / Driver    # Runnable scenario harness (happy path + edge cases)
```

---

## 3. Core Models, Enums & State Transitions

### Enums & States
- **`PrimaryStatus`:** `[INITIAL, IN_PROGRESS_STATE_1, IN_PROGRESS_STATE_2, TERMINAL_1, TERMINAL_2]`
- **`ResourceStatus`:** `[AVAILABLE, RESERVED, BUSY, OFFLINE]`

### Entities & Key Fields
- **`PrimaryEntity`:** `[id, status, referenceId, timestamps, metadata]`
- **`ResourceEntity`:** `[id, status, attributes/location, activeReferenceId]`
- **`Value Objects`:** `[e.g., Location(x, y), Money(amount), TimeWindow]`

### Allowed State Transitions & Guards
- `INITIAL` $\rightarrow$ `IN_PROGRESS_STATE_1`: [Guard condition]
- `IN_PROGRESS_STATE_1` $\rightarrow$ `TERMINAL_1`: [Guard condition]
- Any `TERMINAL_*` state $\rightarrow$ *None* (throw `InvalidStateTransitionException` or return idempotently).

---

## 4. Interfaces & Extensibility Points

- **`[Matching / Selection]Strategy`:**
  - *Responsibility:* Select or rank eligible candidates for a request.
  - *Contract:* `Candidate select(Request req, List<Candidate> candidates)`
- **`[Pricing / Scoring / Policy]Strategy`:**
  - *Responsibility:* Compute cost, score, or policy outcome cleanly without `if/else` sprawl.
  - *Contract:* `Result compute(Context ctx)`
- **`[Entity]Repository`:**
  - *Responsibility:* Thread-safe CRUD and filtered lookups over in-memory collections.

---

## 5. Concurrency & Invariant Guards

- **Shared Mutable State:** [e.g., `ConcurrentHashMap` for entity maps; synchronized blocks / `ReentrantLock` / atomic state transitions when updating both `PrimaryEntity` and `ResourceEntity`.]
- **Critical Invariants:**
  1. [e.g., A `BUSY` resource cannot be assigned to another active request.]
  2. [e.g., State transitions must be validated before mutating fields.]

---

## 6. Execution Plan & Timebox (90–120 Mins)

- **00–15 mins:** Clarify requirements, lock P0 scope, sketch models, enums, and exceptions.
- **15–35 mins:** Implement thread-safe repositories and core strategy interfaces/defaults.
- **35–65 mins:** Implement core services and state transition guards; get P0 happy path compiling and running.
- **65–90 mins:** Add edge-case handling, concurrency guards, and P1 bonus features.
- **90–120 mins:** Build comprehensive `Main` demo/test runner, verify edge cases, and clean up code.

---

## 7. Verification & Edge-Case Checklist

| # | Scenario / Edge Case | Expected Behavior |
| :- | :--- | :--- |
| **1** | **Happy Path** | End-to-end flow succeeds and outputs expected final state/result |
| **2** | **Duplicate / Conflicting Action** | Rejects or handles safely with clear domain exception/message |
| **3** | **Invalid State Transition** | Blocks illegal transitions (e.g., completing a cancelled item) |
| **4** | **Resource Contention** | Concurrent requests for the last/single resource do not double-book |
| **5** | **Empty / Boundary Input** | Handles no available candidates, missing IDs, or zero/edge values gracefully |