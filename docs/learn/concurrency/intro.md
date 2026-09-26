---
title: Introduction to concurrency
description: Shared memory, interleaving, and the three problem types you see in LLD interviews.
next:
  text: correctness
  link: /learn/concurrency/correctness
---

# Introduction

Concurrency means multiple threads executing against the same in-memory state at overlapping times. The core problem is that the order of execution is unpredictable.

## Why do we care?
Consider a parking lot with one empty spot:
```mermaid
sequenceDiagram
    autonumber
    actor A as Thread A
    actor B as Thread B
    participant S as Spot

    A->>S: check spot
    S-->>A: available
    B->>S: check spot
    S-->>B: available
    A->>S: occupy spot
    B->>S: occupy spot
```
Both threads saw the same state. The resulting system state is invalid: two cars believe they own one spot. The important interview question is therefore:
> What happens if two operations execute concurrently?

That's the lens you should apply to almost every LLD problem.

## Concurrency vs Parallelism
These are related but different.
Concurrency = multiple tasks can make progress independently and their operations can interleave.
Parallelism = multiple tasks are literally executing at the same time, typically on different CPU cores.

For LLD interviews, concurrency is the important concept. The operations can interleave unpredictably. On multiple cores, they can actually execute simultaneously. Either way, the problem is the same:
> Multiple execution paths are accessing shared state.

```mermaid
sequenceDiagram
    box #fafafa CONCURRENCY (Single Core)
    participant C1 as CPU Core 1
    end
    
    box #e6f5ff PARALLELISM (Multi-Core)
    participant C2 as CPU Core 2
    participant C3 as CPU Core 3
    end

    %% Concurrency Flow
    Note over C1: 1 Core Context-Switching
    C1->>C1: RUNNING: Task A <br> (Task B is PAUSED)
    C1->>C1: RUNNING: Task B <br> (Task A is PAUSED)
    
    %% Parallelism Flow
    Note over C2, C3: 2 Cores Running Simultaneously
    par True Parallelism
        C2->>C2: RUNNING: Task A
        C3->>C3: RUNNING: Task B
    end

```

## The real source of concurrency bugs
Most concurrency bugs come from shared mutable state.
```mermaid
sequenceDiagram
    autonumber
    actor A as Thread A
    actor B as Thread B
    participant State as Shared State: [inventory = 1]

    Note over A, B: "Check-Then-Act" Concurrency Bug

    %% Phase 1: Interleaved Reads
    A->>State: READ inventory (returns 1)
    B->>State: READ inventory (returns 1)
    
    %% Phase 2: Independent Checks
    Note over A: CHECK: 1 > 0 ? (True)
    Note over B: CHECK: 1 > 0 ? (True)

    %% Phase 3: Double Overwrite
    rect #ffdee2
        Note over A, B: Race Condition: Overwriting without synchronization
        A->>State: WRITE inventory - 1 (sets to 0)
        B->>State: WRITE inventory - 1 (sets to 0)
    end

    Note over State: Final State: inventory = 0 <br> [BUG: 1 item was consumed twice!]

```
One item was consumed twice. This pattern is called check-then-act and will show up constantly in LLD interviews.

## What concurrency adds to an LLD problem
Normally you might design:
```text
ParkingLot
    └── ParkingSpot

park(car)
unpark(car)
```

Then the interviewer asks:
> What if two cars try to take the same spot simultaneously?

Now you need to answer:
1. What state is shared?
2. What operations can overlap?
3. What invariant must never be violated?
4. How do we protect that invariant?

For example:
```text
Invariant: A parking spot can belong to at most one car.
```
Then:
```text
find available spot + claim spot  --> must effectively behave as one atomic operation.
```
That's the essence of concurrency in LLD.

## The 3 concurrency problems you should recognize
Most LLD concurrency questions fall into three buckets:

| Problem | Core question | Typical solution |
| :--- | :--- | :--- |
| **Correctness** | Can concurrent operations corrupt state? | Locks, atomics |
| **Coordination** | How do threads communicate/wait? | Queues, conditions, channels |
| **Scarcity** | How do we limit access to a finite resource? | Semaphores, pools |



#### Correctness
```text
Two threads → same seat
```
Need to protect shared state.

#### Coordination
```text
Producer → Queue → Consumer
```
Consumer may need to wait until work exists.

#### Scarcity
```text
100 requests
     ↓
10 DB connections
```
Only 10 operations can use the resource concurrently. These categories are more useful in interviews than memorizing individual concurrency APIs.

## The mental model for interviews
When the interviewer introduces concurrency, don't immediately say:
> I'll use a mutex.

First ask:
#### Step 1: What is shared?
- inventory
- parking spots
- account balance
- queue
- connection pool

#### Step 2: What can happen concurrently?
- reserve()
- cancel()
- update()
- read()

#### Step 3: What must remain true?
These are your invariants. Example:
- inventory >= 0
- one seat → at most one reservation

#### Step 4: What synchronization is required?
Only now choose:
- atomic
- mutex
- RW lock
- semaphore
- condition variable
- blocking queue

### The key takeaway
For LLD, don't think:
> Concurrency = locks.

Think:
> Concurrency = unpredictable interleaving of operations on shared state.

Your job is to identify the shared state + invariant, then choose the simplest mechanism that preserves that invariant.

That's the foundation. The next topic should be Correctness, where we go deep into race conditions, atomicity, check-then-act, locks, atomics, and how to reason about whether a piece of code is actually thread-safe.