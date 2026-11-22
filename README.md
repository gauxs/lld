# LLDZen

Low-level design theory, interview problems, and reference implementations.
Problems progress through a baseline and optional extensions, each with
requirements, design, and optional code.

## Run locally

```bash
npm install
npm run docs:dev
```

Open `http://localhost:5173/lld/`.

## Repository

- `theory/` — source Markdown for theory trails
- `problems/` — problem requirements, designs, extensions, and code
- `practice_with_ai/` — prompts for chatbot-guided interview practice
- `_template/problem/` — starting point for a new problem
- `docs/` — generated VitePress site and theme

Theory, problems, and practice guides are discovered automatically from their
directory structure. Do not edit generated pages under `docs/learn/`,
`docs/problems/`, or `docs/practice-with-ai/`.

See **[Add study content](_template/setup.md)** for the authoring workflow.

## Build

```bash
npm run docs:build
npm run docs:preview
```
# Staff-Level LLD & Machine Coding Prep Checklist (Go)

At **Staff level**, interviewers treat clean interfaces and modularity as table stakes. The bar shifts to **concurrency correctness, lock granularity, failure handling (idempotency, retries, rollbacks), memory/latency trade-offs, and API ergonomics**.

---

## Tier 1: Core Infrastructure & Concurrency (Highest Staff Yield)

- [ ] **1. Movie Ticket Booking (BookMyShow)**
  * **Look out for:** Lock granularity (show-level vs. seat-level locks), temporary seat holds with TTL expiry, preventing deadlocks when booking multiple seats simultaneously (sort seat IDs before locking), and idempotent payment callbacks.
- [ ] **2. Rate Limiter**
  * **Look out for:** Race conditions on read-modify-write counters, memory leaks from inactive user keys (TTL cleanup), clock skew (`time.Now()` vs. monotonic clock), and pluggable algorithms (`TokenBucket`, `SlidingWindowLog`).
- [ ] **3. LRU / LFU Cache**
  * **Look out for:** Lock contention under high concurrency (sharded/striped locks vs. single global `sync.RWMutex`), $O(1)$ eviction without memory leaks in linked list pointers, and optional TTL expiration.
- [ ] **4. In-Memory Pub-Sub / Message Queue (Kafka-like)**
  * **Look out for:** Slow consumer backpressure (blocking vs. buffered channels vs. drop), per-consumer offset tracking, graceful shutdown (`context.Context`), and goroutine leak prevention.
- [ ] **5. In-Memory Key-Value Store with Transactions**
  * **Look out for:** Nested transactions (`BEGIN`/`COMMIT`/`ROLLBACK`), isolation semantics, snapshot vs. write-ahead undo logs, and thread-safe concurrent readers/writers.
- [ ] **6. Task Scheduler (Cron / Delayed Job Executor)**
  * **Look out for:** Min-heap (`container/heap`) wakeup precision without busy-waiting, worker pool bounding, retry with exponential backoff, and cancellation of pending tasks.
- [ ] **7. Logging Framework (Log4j)**
  * **Look out for:** Non-blocking async log ingestion (ring buffer / channel), backpressure policy when sinks are slow, structured formatting, and safe sink rotation without dropping logs.

---

## Tier 2: High-Contention Domain & Financial Systems

- [ ] **8. Payment Wallet / Ledger System (PhonePe / Paytm)**
  * **Look out for:** Double-entry bookkeeping (never just mutate a balance field), deadlock prevention on bidirectional transfers (lock accounts in deterministic ID order), and idempotency keys to prevent duplicate charges.
- [ ] **9. Parking Lot**
  * **Look out for:** Concurrent entry/exit gates competing for the last spot, scaling spot lookup from $O(N)$ scan to $O(1)$ per vehicle type (min-heap by distance to gate), and pluggable pricing/allocation strategies.
- [ ] **10. Splitwise (Expense Sharing)**
  * **Look out for:** Floating-point precision bugs (use integer paise/cents, handle remainder distribution), concurrent group updates, and debt simplification graph trade-offs (NP-complete exact vs. greedy heuristic).
- [ ] **11. Ride-Sharing App (Uber / Ola)**
  * **Look out for:** Preventing two drivers from accepting the same ride (atomic state transition / CAS), spatial indexing abstraction for nearby drivers, and handling driver/rider cancellation race conditions.
- [ ] **12. Food Delivery App (Swiggy / Zomato)**
  * **Look out for:** Strict order state machine transitions (rejecting invalid jumps like `Delivered -> Cancelled`), atomic inventory reservation at restaurant, and async event hooks for notifications.
- [ ] **13. Shopping Cart + Coupon/Discount Engine**
  * **Look out for:** Deterministic rule evaluation order (priority/stacking conflicts), immutability of applied snapshots at checkout, and extensibility for new rule predicates without touching core cart logic.
- [ ] **14. Meeting Room Booking System**
  * **Look out for:** Interval overlap queries (interval tree / sorted slice vs. brute force), atomic check-and-book per room, and handling recurring bookings cleanly.

---

## Tier 3: Hardware, Resource & State-Heavy Systems

- [ ] **15. Elevator System**
  * **Look out for:** Dispatch algorithm pluggability (SCAN/LOOK vs. nearest-car), decoupling hall/cabin requests from movement ticks, and safe concurrent button presses.
- [ ] **16. Notification Service**
  * **Look out for:** Per-channel rate limiting, retry queues with dead-letter handling, user preference/deduplication filtering, and template rendering isolation.
- [ ] **17. Car Rental System (ZoomCar)**
  * **Look out for:** Time-window availability locking, grace-period/late-fee calculation edge cases, and vehicle maintenance state blocking reservations.
- [ ] **18. Library Management System**
  * **Look out for:** Separating `Book` (metadata) from `BookItem` (physical copy), FIFO reservation queue when a copy is returned, and concurrent checkout of the last available copy.

---

## Tier 4: Games & State Machines (Clean Abstraction & Testing)

- [ ] **19. Vending Machine**
  * **Look out for:** Explicit `State` interface transitions (`Idle -> HasMoney -> Dispensing`), exact change calculation failure (rolling back when machine lacks change), and inventory consistency.
- [ ] **20. ATM Machine**
  * **Look out for:** `Chain of Responsibility` for note dispensation with rollback if exact amount cannot be formed, PIN retry lockout state, and atomic account debit.
- [ ] **21. Tic-Tac-Toe / Chess Validator**
  * **Look out for:** Pure domain model completely decoupled from I/O (easy to unit-test board states), $O(1)$ win detection on $N \times N$ boards, and `Command` pattern for move history/undo.
- [ ] **22. Snake & Ladder**
  * **Look out for:** Deterministic board initialization (preventing infinite snake-ladder loops), injectable `Dice` interface for deterministic unit tests, and clean game-loop termination.

---

## Go-Specific Staff Signals to Showcase

- [ ] **Clean package & interface boundaries:** Define interfaces on the *consumer* side (small, 1–3 methods) rather than massive Java-style service interfaces.
- [ ] **Explicit concurrency control:** Always explain *why* you chose `sync.Mutex` vs. `sync.RWMutex` vs. channels, and run `go test -race ./...` during machine coding demos.
- [ ] **Deterministic testability:** Inject `Clock` (`func() time.Time`) and `IDGenerator` interfaces so TTL expirations, rate limits, and bookings can be tested deterministically without `time.Sleep()`.