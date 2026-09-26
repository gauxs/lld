---
title: Sync in Go
description: Standard library synchronization in Go, and when to use each mechanism in LLD interviews.
prev:
  text: introduction
  link: /learn/concurrency/01-introduction
next:
  text: correctness
  link: /learn/concurrency/03-correctness
---


# Sync in Go

In a low-level design interview, concurrency questions usually test one skill: can you identify **shared mutable state**, state the **invariant** it must satisfy, and choose a standard library mechanism that preserves that invariant without over-engineering the design?

Go gives you goroutines plus a focused set of tools in `sync`, `sync/atomic`, channels, and the `select` statement. This page summarizes what each tool does and when it belongs in your answer. For worked patterns and exercises, continue with [Correctness](/learn/concurrency/03-correctness), [Coordination](/learn/concurrency/04-coordination), and [Scarcity](/learn/concurrency/05-scarcity).

## Goroutines and shared memory

A goroutine is a function executing concurrently with the rest of the program. The runtime schedules many goroutines on fewer OS threads; starting one is inexpensive compared to a new thread.

```go
go handle(req)
```

All goroutines in a process share the same address space: heap allocations, package-level variables, and fields reachable through shared pointers. That sharing is what creates races when two goroutines read and write the same memory without coordination.

Useful interview discipline: when you sketch a type on the board, mark which fields can be written by more than one goroutine. Those fields need a deliberate strategy (lock, atomic, confinement, or message passing).

## Mutex and RWMutex

`sync.Mutex` enforces **mutual exclusion**: at most one goroutine runs the guarded code at a time. Hold the lock across every step that must appear atomic to callers, especially check-then-act sequences on shared fields.

```go
type Ledger struct {
    mu      sync.Mutex
    balance int64
}

func (l *Ledger) Withdraw(amount int64) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.balance < amount {
        return false
    }
    l.balance -= amount
    return true
}
```

`sync.RWMutex` allows many concurrent readers or a single writer. Consider it when reads dominate (for example, serving a cached configuration) and writes are infrequent.

**Lock granularity:** a single mutex over a service is simple to explain and often enough in an interview. Finer locks (per shard or per resource) can improve throughput but require consistent lock ordering to avoid deadlock. [Correctness](/learn/concurrency/03-correctness) walks through both approaches.

## Package sync/atomic

The `atomic` package provides indivisible operations on individual values: add, load, store, compare-and-swap, and related helpers on types such as `atomic.Uint64`.

```go
var completed atomic.Uint64
completed.Add(1)
```

An atomic increment is safe as a single operation. A **sequence** of operations is not automatically safe:

```go
if completed.Load() > 0 {
    completed.Add(-1) // load and add are still two steps; another goroutine can interleave
}
```

Reach for atomics when updating one counter or flag is the full story. When correctness depends on multiple fields staying consistent with each other, prefer a mutex or thread confinement. See [Correctness](/learn/concurrency/03-correctness).

## Channels

Channels combine a queue with blocking send and receive. They are the idiomatic way to **transfer** data or work between goroutines instead of sharing a mutable structure every worker updates.

- **Unbuffered:** the sender blocks until a receiver is ready; the handoff itself synchronizes the two goroutines.
- **Buffered:** sends proceed until the buffer fills, then senders block (backpressure on producers).

```go
tasks := make(chan Task, 64)

tasks <- t     // blocks when the buffer is full
job := <-tasks // blocks when the buffer is empty
```

Closing a channel signals that no more values will be sent. Receivers can drain remaining values with `range`:

```go
close(tasks)
for job := range tasks {
    run(job)
}
```

Producer-consumer layouts, worker pools, and shutdown ordering are developed in [Coordination](/learn/concurrency/04-coordination).

## The select statement

`select` waits on **multiple channel operations** and runs the branch for whichever is ready first. If several branches are ready, the runtime chooses one pseudo-randomly (do not rely on a fixed priority unless you structure the code that way).

Use it when a goroutine must react to more than one source of events: incoming work, results, timeouts, or shutdown.

```go
select {
case job := <-tasks:
    process(job)
case err := <-errc:
    return err
case <-ctx.Done():
    return ctx.Err()
}
```

A `default` branch makes the `select` non-blocking: if nothing is ready, execution continues immediately. That pattern suits polling or try-send/try-receive without parking the goroutine.

```go
select {
case tasks <- job:
    // enqueued
default:
    // queue full; apply backpressure or drop policy
}
```

In LLD answers, `select` is how you tie channels to **timeouts** (`context`) and **graceful shutdown** (a dedicated `done` channel or `ctx.Done()`). It does not replace mutexes for shared structs; it coordinates **communication** between goroutines. More patterns in [Coordination](/learn/concurrency/04-coordination).

## sync.WaitGroup

A `WaitGroup` waits until a fixed set of goroutines finishes. It coordinates **completion**, not access to shared memory.

```go
var wg sync.WaitGroup
for i := 0; i < workers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        processBatch()
    }()
}
wg.Wait()
```

A common mistake is to use `WaitGroup` alone when goroutines also update a shared counter or map. Waiting for goroutines to exit does not make those updates race-free; you still need a mutex, atomics where appropriate, or confinement.

## sync.Cond

A condition variable lets a goroutine wait until a predicate becomes true, typically while associated with a mutex. It appears in custom queues and similar structures when you already hold a lock and need efficient waiting.

```go
mu.Lock()
for !queueNonEmpty {
    cond.Wait()
}
item := dequeue()
mu.Unlock()
```

Many interview designs are clearer with a buffered channel or `select`. Mentioning `sync.Cond` shows you understand how blocking queues synchronize internally. [Coordination](/learn/concurrency/04-coordination) expands on these patterns.

## Semaphore pattern

The standard library does not define a `Semaphore` type. The usual Go pattern is a buffered channel of empty structs: capacity `N` means at most `N` goroutines hold a permit at once.

```go
slots := make(chan struct{}, maxOpen)

slots <- struct{}{}
defer func() { <-slots }()

queryDB()
```

For weighted or dynamic limits, `golang.org/x/sync/semaphore` provides `Acquire` and `Release`. Release permits in a `defer` (or equivalent) so panics do not leak capacity. Connection pools, concurrency caps, and rate limiting are covered in [Scarcity](/learn/concurrency/05-scarcity).

## sync.Map and context.Context

**sync.Map** is a concurrent map intended for caches with many reads and stable keys. In interviews, a `map` protected by `sync.Mutex` is often easier to justify unless the problem explicitly calls for a heavily contended shared cache.

**context.Context** carries cancellation and deadlines through a call tree. It does not replace synchronization for shared fields, but it helps you stop work when a client disconnects or a timeout expires. Mention it when describing request handlers, pipelines, or graceful shutdown.

## Where to go next

| Interview concern | Chapter |
| --- | --- |
| Concurrent updates corrupt shared state | [Correctness](/learn/concurrency/03-correctness) |
| Goroutines must queue work, wait, or shut down cleanly | [Coordination](/learn/concurrency/04-coordination) |
| Only a bounded number of operations or resources may run at once | [Scarcity](/learn/concurrency/05-scarcity) |

Production designs often combine all three concerns. Separating them in your explanation helps the interviewer follow your reasoning even when the final design uses several mechanisms together.
