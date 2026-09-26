---
title: Interview approach
description: Shared state, invariants, then pick a synchronization primitive—before reaching for a mutex.
prev:
  text: scarcity
  link: /learn/concurrency/02-problems/03-scarcity
---


# Interview approach

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
