---
title: Coordination
sidebar: Coordination
description: Ordering work between threads—producer–consumer and backpressure.
---

# Coordination

<p class="lead">Coordination problems are about who runs when: handing tasks between threads, waiting for readiness, and avoiding busy loops.</p>

## Failure modes

- Consumers spin while the queue is empty (wasted CPU).
- Producers overwrite or drop work when consumers lag.
- Shutdown leaves goroutines blocked forever.

## What to reach for

- **Blocking queue** or **buffered channel** between producers and workers.
- **Condition variables** (or channel select) to park threads until data arrives.
- Clear **lifecycle**: close the channel or signal done when shutting down.

## In interviews

Draw the flow: request arrives → enqueue → worker pool → response. Name where threads block and what unblocks them.

```go
jobs := make(chan Job, 100)
for w := 0; w < workerCount; w++ {
    go func() {
        for job := range jobs {
            process(job)
        }
    }()
}
```

Follow-ups on Connect Four (networked play) are coordination-heavy: authoritative server, push state to clients.
