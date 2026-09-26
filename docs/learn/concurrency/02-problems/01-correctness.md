---
title: Correctness
description: Shared state corruption, check-then-act, and mutual exclusion.
prev:
  text: overview
  link: /learn/concurrency/02-problems/00-overview
next:
  text: coordination
  link: /learn/concurrency/02-problems/02-coordination
---


# Correctness

<p class="lead">Correctness problems happen when two threads observe or update shared state in an order that violates your invariants—double booking, lost increments, or stale reads.</p>

## Failure modes

**Check-then-act:** Thread A and B both read “seat available,” both proceed, one update is lost.

**Read-modify-write:** `counter++` is load, add, store. Two interleaved increments can both read the same value.

**Invariant violations:** A structure is valid only if several fields agree; without a lock, another thread sees a half-updated object.

## What to reach for

1. **Mutex** around the smallest critical section that must appear atomic to other threads.
2. **Atomics** when a single word is the whole story (metrics, feature flags).
3. **Confinement**—one goroutine owns the data; others send messages (no shared mutable state).

## In interviews

State the invariant aloud (“at most one owner per seat”), then show where two threads can break it. Prefer one clear lock over a clever lock-free structure unless the prompt demands throughput.

## Example sketch (Go)

```go
type Inventory struct {
    mu    sync.Mutex
    stock map[string]int
}

func (i *Inventory) Reserve(sku string) error {
    i.mu.Lock()
    defer i.mu.Unlock()
    if i.stock[sku] <= 0 {
        return ErrOutOfStock
    }
    i.stock[sku]--
    return nil
}
```

Connect Four **baseline** deliberately avoids this layer (single-threaded CLI). The **[rate limiter](/problems/rate-limiter/extensions/baseline/requirements)** trail is where thread safety matters (`problems/rate_limiter/extensions/baseline/code/`).
