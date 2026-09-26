---
title: Sync in Go
sidebar: Sync in Go
description: The stdlib pieces interviewers expect you to name—and when each one earns its place.
---

# Sync in Go

Interviews are not a quiz on obscure APIs. They check whether you can **see shared state**, state an **invariant**, and reach for something from `sync`, `sync/atomic`, or a channel without turning the design into a mutex soup.

## Goroutines: many cooks, one kitchen

A goroutine is just a function running concurrently. Cheap to start, scheduled by the runtime, and here's the catch, it shares the **heap** with every other goroutine in the process.

```go
go handle(req)
```

That one line means overlap. Anything on the heap both goroutines can reach is fair game for interleaving unless you deliberately isolate it (one owner, or a lock, or a message passed on a channel).

Interview habit worth building: 
> when you draw a struct on the whiteboard, circle the fields that **more than one goroutine can write**. Those fields are where the story gets interesting.

## Mutex: one person edits the ledger

`sync.Mutex` is a turnstile. One goroutine passes through the critical section; everyone else waits. Use it when several steps must look like one step to the rest of the program, classic check-then-act on shared fields.

```go
type Ledger struct {
    mu sync.Mutex
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

`sync.RWMutex` is the same idea with a reading room: many readers **or** one writer. Mention it when the prompt is read-heavy (metrics snapshot, config cache) and writes are rare.

**Granularity trade-off:** 
> one big lock is easy to defend in an interview; many small locks buy concurrency but invite deadlock if ordering is sloppy. You'll practice both extremes in [Correctness](/learn/concurrency/03-correctness).

## Atomics: one counter, one instruction

`sync/atomic` is for when the *entire* update is "bump this number" or "flip this flag", nothing else has to stay in sync with it.

```go
var completed atomic.Uint64
completed.Add(1)
```

Hardware makes that increment indivisible. It does **not** make your *design* indivisible if you still need a check before the bump:

```go
if completed.Load() > 0 {
    completed.Add(-1) // still a race: load and add are two steps
}
```

When two fields must move together, balance and pending holds, seats and waitlist, atomics alone won't save you; a mutex or a single owner goroutine will. More examples in [Correctness](/learn/concurrency/03-correctness).

## Channels: pass notes, don't share the notebook

Channels are how Go prefers you **hand off** work or results instead of letting every goroutine scribble on the same map.

- **Unbuffered** - sender and receiver meet at the handoff (synchronization built in).
- **Buffered** - a waiting line of fixed length; a full buffer slows producers down (backpressure).

```go
tasks := make(chan Task, 64)

tasks <- t    // blocks if the line is full
job := <-tasks // blocks if the line is empty
```

Closing tells receivers "no more coming":

```go
close(tasks)
for job := range tasks {
    run(job)
}
```

Think conveyor belt, not shared whiteboard. Producer–consumer sketches and worker pools live in [Coordination](/learn/concurrency/04-coordination).

## WaitGroup: roll call before you leave

`sync.WaitGroup` answers "are the helper goroutines done yet?" It does **not** protect shared data, it only waits.

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

We often slap a WaitGroup on a broken counter and wonder why the count is still wrong. The group waits for **completion**; you still need mutex, atomic, or confinement for **correctness**.

## sync.Cond: the custom waiting room

`sync.Cond` pairs with a mutex: "I can't proceed until *this* boolean becomes true." Lower level than channels; you'll see it inside home-grown queues.

```go
mu.Lock()
for !queueNonEmpty {
    cond.Wait()
}
item := dequeue()
mu.Unlock()
```

In many LLD answers, a buffered channel or `select` reads cleaner. Knowing `Cond` exists signals you understand what blocking queues are doing under the hood, covered more in [Coordination](/learn/concurrency/04-coordination).

## Semaphore: only N guests inside

Go doesn't ship `Semaphore` in the stdlib, but the idiom is a buffered channel of empty structs **N tokens**, take one before work, put it back after.

```go
slots := make(chan struct{}, maxOpen)

slots <- struct{}{}
defer func() { <-slots }()

queryDB()
```

Weighted limits (memory-sized permits) show up in `golang.org/x/sync/semaphore`. 

> Always release in `defer` so a panic doesn't permanently shrink capacity. Pools, caps, and rate limits are the [Scarcity](/learn/concurrency/05-scarcity) chapter.

## sync.Map and context (the supporting cast)

**sync.Map** - a concurrent map tuned for caches that mostly read stable keys. In interviews, `map` + `Mutex` is usually easier to explain unless the prompt is explicitly a hot shared cache.

**context.Context** - cancellation and deadlines flowing down the call stack. It won't fix a race on your counter, but it stops goroutines from working forever after the client hung up, worth mentioning when you describe HTTP handlers or graceful shutdown.

## Follow-up reading

| If the interview worry sounds like… | Start here |
| --- | --- |
| "Two requests corrupt the same record" | [Correctness](/learn/concurrency/03-correctness) - mutex, atomic, confinement |
| "Workers need a queue and clean shutdown" | [Coordination](/learn/concurrency/04-coordination) - channels, `select`, lifecycle |
| "We only have ten DB connections" | [Scarcity](/learn/concurrency/05-scarcity) - semaphores, pools, limits |

Real systems blend all three; separating them in your explanation is a clarity move, not a claim that the world comes in neat boxes.
