---
title: Coordination
sidebar: Coordination
description: Ordering work between threads—producer–consumer and backpressure.
---

# Coordination

Concurrency creates a second problem after correctness:
> How do concurrent workers coordinate their work?

A useful mental model is:
- Shared state → coordinate by accessing the same state safely.
- Message passing → coordinate by sending work/data between workers.

In LLD interviews, the key is choosing the simplest mechanism that matches the problem.

## Coordination via shared state
Multiple goroutines access the same state.
```go
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value++
}
```
Use this when workers genuinely need to observe or modify common state. Common tools:
| Tool | Use |
|---|---|
| `Mutex` | Protect shared mutable state |
| `RWMutex` | Many readers, few writers |
| `atomic` | Simple atomic state such as counters |
| `sync.Cond` | Wait for a state condition |

The important question is:
> What state are the goroutines coordinating around?

## Coordination via message passing
Instead of sharing state, goroutines communicate by sending messages. In Go, channels are the primary mechanism.
```go
jobs := make(chan Job)

go worker(jobs)

jobs <- job
```
The worker receives the job:
```go
func worker(jobs <-chan Job) {
	for job := range jobs {
		process(job)
	}
}
```

Now the workers don't need to directly share the job queue.

## Shared state vs message passing
Neither is universally better. A useful rule:
```text
If workers need to coordinate state, use synchronization.
If workers need to coordinate work, consider message passing.
```

## Common Coordination Problems
### Process Requests Asynchronously
A common LLD requirement is:
> The API should accept requests quickly, while expensive processing happens in the background.

Don't make the request handler perform the expensive work. Instead:
```text
                    ┌── Worker
                    │
Request → Queue ────┼── Worker
                    │
                    └── Worker
```
A simple Go implementation:
```go
type Job struct {
	ID int
}

jobs := make(chan Job, 100)

func worker(jobs <-chan Job) {
	for job := range jobs {
		process(job)
	}
}

func handler(job Job) {
	jobs <- job
	// Return immediately.
}
```
The request path only submits the work. Workers process it independently.

#### What should you consider?
In an interview, don't stop at the channel.
Ask:
1. What happens if the queue is full?
2. How many workers should there be?
3. Can jobs be lost?
4. Can jobs be retried?
5. Does ordering matter?
6. How does the system shut down?
7. What happens if processing fails?

An in-memory channel is not durable. If requirement is that jobs should not be lost:

> API → Durable Queue → Workers

Examples include Kafka, SQS, and RabbitMQ.

### Handle Bursty Traffic
Suppose the system normally receives: 100 requests/sec but occasionally receives: 10,000 requests/sec. If every request immediately consumes a worker, the downstream system can become overloaded. 

A queue provides a buffer:
```text
             ┌── Worker
             │
Requests → Queue ── Worker
             │
             └── Worker
```

The queue absorbs short-term bursts while workers process at a controlled rate. A bounded Go channel:
```go
jobs := make(chan Job, 1000)
```

When the queue is full, you need an explicit policy. The policies can be:
#### Block
Wait until capacity becomes available.
```go
jobs <- job
```
This creates **backpressure**.

#### Reject
Fail immediately when the queue is full.
```go
select {
case jobs <- job:
	// accepted
default:
	// queue full
	// reject
}
```
For an HTTP API, this might result in a 429 or another appropriate overload response.

#### Scale
Increase the number of workers:
```text
10 workers → 50 workers
```
But only if downstream dependencies can handle the additional load.

#### Important

A queue does not solve sustained overload. If:
```text
arrival rate > processing rate
```
for long enough, the queue will eventually fill.

You need backpressure, rejection, load shedding, or additional processing capacity.

## Conclusion
Coordination is about how concurrent workers cooperate.

Remember these two models:
```text
Shared state    ->  Mutex / RWMutex / Atomic / Cond

Message passing ->  Channels / Queues
```

For LLD interviews, focus less on memorizing primitives and more on the behavior:

> Who produces work, who consumes it, what happens when consumers are busy, and what happens when the system is overloaded?