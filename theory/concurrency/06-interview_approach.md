---
title: Interview approach
sidebar: Interview approach
description: Classify the concern, state invariants, then pick the simplest sync or coordination tool—not a mutex by reflex.
---

# The mental model for interviews

```mermaid
graph TB
    Start([Interviewer Introduces Concurrency])

    Start -.->|instinct| Trap[/"❌ 'I will use a Mutex!'"/]
    style Trap fill:#fecaca,stroke:#b91c1c,stroke-width:2px,color:#7f1d1d,stroke-dasharray: 5 5

    Start -->|better approach| S1

    S1["<strong>Step 1: Shared state</strong><br>• Which fields can multiple goroutines write?<br>• Inventory, spots, balance, maps, queue depth"]
    --> S2["<strong>Step 2: Concurrent actions</strong><br>• Which methods can interleave?<br>• reserve(), cancel(), enqueue(), read()"]
    --> S3["<strong>Step 3: Invariants</strong><br>• inventory ≥ 0 · one seat → one holder<br>• <em>Patterns:</em> check-then-act · read-modify-write"]
    --> S4["<strong>Step 4: Classify the concern</strong><br>• <b>Correctness</b> — can updates corrupt state?<br>• <b>Coordination</b> — async work, queues, shutdown<br>• <b>Scarcity</b> — only N concurrent or N resources"]
    --> S5["<strong>Step 5: Simplest mechanism</strong><br><b>Shared state:</b> Mutex · RWMutex · atomic · sync.Cond · confinement<br><b>Message passing:</b> chan · select · WaitGroup · context<br><b>Scarcity:</b> semaphore (chan) · pool · block / reject / timeout"]
    --> Takeaway("<strong>Takeaway</strong><br>Concurrency = unpredictable interleaving on shared state.<br>Name state + invariant first, classify the concern,<br>then the smallest tool — often several combined.")

    style Start fill:#ffedd5,stroke:#c2410c,stroke-width:2px,color:#7c2d12
    style S1 fill:#ffffff,stroke:#475569,stroke-width:2px,color:#0f172a
    style S2 fill:#ffffff,stroke:#475569,stroke-width:2px,color:#0f172a
    style S3 fill:#ffffff,stroke:#475569,stroke-width:2px,color:#0f172a
    style S4 fill:#ffffff,stroke:#475569,stroke-width:2px,color:#0f172a
    style S5 fill:#ffffff,stroke:#475569,stroke-width:2px,color:#0f172a
    style Takeaway fill:#bbf7d0,stroke:#15803d,stroke-width:2px,color:#14532d
```

### Coordination follow-ups (don’t stop at “I’ll use a channel”)

When the prompt is async APIs, worker pools, or background processing, say out loud:

1. What happens when the **queue is full** — block, reject (`select` + `default`), or scale?
2. How does **graceful shutdown** work — stop accepting, drain, wait for workers?
3. Can work be **lost** — in-memory channel vs durable queue?
4. Does **ordering** or **retry on failure** matter?

### The key takeaway

For LLD, don't think:

> Concurrency = locks.

Think:

> Concurrency = unpredictable interleaving of operations on shared state.

Your job is to identify **shared state and invariants**, classify whether the hard part is **correctness**, **coordination**, or **scarcity** (often more than one), then choose the **simplest** mechanism that preserves the invariant — mutex when you must; channels when you are moving work; semaphores and pools when capacity is the limit.
