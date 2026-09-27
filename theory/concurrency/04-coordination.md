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

## Exercises
### Exercise 1: Worker Pool
Three workers process jobs concurrently.
Task: Make main wait until all workers have finished before printing "Done".
<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Problem code</summary>

```go
package main

import (
	"fmt"
	"time"
)

func worker(id int) {
	fmt.Println("Worker", id, "started")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Worker", id, "finished")
}

func main() {
	for i := 1; i <= 3; i++ {
		go worker(i)
	}

	fmt.Println("Done")
}
```

</details>

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Solution code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Worker", id, "started")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Worker", id, "finished")
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()

	fmt.Println("Done")
}
```
</details>

### Exercise 2: Producer-Consumer
A producer generates jobs and consumers process them.

Task: Use a channel to coordinate the producer and workers.

Requirements:
- Start 3 workers.
- Submit 10 jobs.
- Workers should process jobs as they become available.
- Workers should stop after all jobs have been processed.
- main should wait for all workers to finish.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Problem code</summary>

```go
package main

import (
	"fmt"
	"time"
)

type Job struct {
	ID int
}

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	jobs := make(chan Job)

	for i := 1; i <= 3; i++ {
		go worker(i, jobs)
	}

	for i := 1; i <= 10; i++ {
		jobs <- Job{ID: i}
	}

	fmt.Println("Done")
}
```

</details>

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Solution code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	jobs := make(chan Job)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	for i := 1; i <= 10; i++ {
		jobs <- Job{ID: i}
	}

	close(jobs)
	wg.Wait()

	fmt.Println("Done")
}

```

</details>

### Exercise 3: Bounded Queue
A queue has a fixed capacity.

Task: Implement a thread-safe queue using sync.Mutex and sync.Cond.

Requirements:
- Capacity = 3.
- Put blocks when the queue is full.
- Get blocks when the queue is empty.
- Multiple producers and consumers must be supported.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	// TODO
}

func NewQueue(capacity int) *Queue {
	// TODO
	return nil
}

func (q *Queue) Put(item int) {
	// TODO
}

func (q *Queue) Get() int {
	// TODO
	return 0
}

func main() {
	q := NewQueue(3)

	var wg sync.WaitGroup

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 10; i++ {
			fmt.Println("Producing:", i)
			q.Put(i)
		}
	}()

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			item := q.Get()
			fmt.Println("Consumed:", item)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	wg.Wait()

	fmt.Println("Done")
}
```
</details>

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Solution code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	mu         sync.Mutex
	qFullCond  *sync.Cond
	qEmptyCond *sync.Cond
	capacity   int
	arr        []int
	backIdx    int
	frontIdx   int
	len        int
}

func NewQueue(capacity int) *Queue {
	q := &Queue{
		mu:       sync.Mutex{},
		capacity: capacity,
		arr:      make([]int, capacity),
		backIdx:  0,
		frontIdx: 0,
		len:      0,
	}

	q.qFullCond = sync.NewCond(&q.mu)
	q.qEmptyCond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Put(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.len == q.capacity {
		q.qFullCond.Wait()
	}

	q.arr[q.backIdx] = item
	q.backIdx = (q.backIdx + 1) % q.capacity
	q.len++
	q.qEmptyCond.Signal()
}

func (q *Queue) Get() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	for q.len == 0 {
		q.qEmptyCond.Wait()
	}

	t := q.arr[q.frontIdx]
	q.frontIdx = (q.frontIdx + 1) % q.capacity
	q.len--
	q.qFullCond.Signal()
	return t
}

func main() {
	q := NewQueue(3)

	var wg sync.WaitGroup

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 10; i++ {
			fmt.Println("Producing:", i)
			q.Put(i)
		}
	}()

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			item := q.Get()
			fmt.Println("Consumed:", item)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	wg.Wait()

	fmt.Println("Done")
}

```

</details>

### Exercise 4: Asynchronous Email Service
Build an email service that processes emails in the background.

Requirements:
- 5 workers.
- Maximum 10 pending emails.
- SendEmail should return immediately if the queue has capacity.
- If the queue is full, SendEmail should return an error immediately.
- Shutdown should wait for accepted emails to finish.
- No email accepted by the service should be lost during graceful shutdown.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Email struct {
	To      string
	Subject string
}

type EmailService struct {
	// TODO
}

func NewEmailService() *EmailService {
	// TODO
	return nil
}

func (s *EmailService) SendEmail(email Email) error {
	// TODO
	return nil
}

func (s *EmailService) Shutdown() {
	// TODO
}

func processEmail(email Email) {
	fmt.Printf("Sending email to %s: %s\n", email.To, email.Subject)
	time.Sleep(200 * time.Millisecond)
}

func main() {
	service := NewEmailService()

	for i := 1; i <= 20; i++ {
		err := service.SendEmail(Email{
			To:      fmt.Sprintf("user%d@example.com", i),
			Subject: fmt.Sprintf("Email %d", i),
		})

		if err != nil {
			fmt.Println("Rejected:", i)
		}
	}

	fmt.Println("All requests submitted")

	service.Shutdown()

	fmt.Println("Service stopped")
}
```

</details>
<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Solution code</summary>

```go
```

</details>

## Challenges (No solutions)

### Challenge 1: Worker Pool with Results
You have 100 jobs and 5 workers.

Each job produces a result:
```go
type Result struct {
	JobID int
	Value int
}
```

Task:
- Process all jobs concurrently.
- Collect every result.
- Print the sum of all results.
- main must not exit before every result has been received.
- Avoid shared mutable state for the result collection.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Value int
}

func process(job Job) Result {
	return Result{
		JobID: job.ID,
		Value: job.ID * 10,
	}
}

func worker(
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	// TODO
}

func main() {
	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup

	// TODO: start 5 workers

	// TODO: submit 100 jobs

	// TODO: close jobs

	// TODO: wait for workers and close results

	total := 0

	// TODO: collect results

	fmt.Println("Total:", total)
}
```
</details>

### Challenge 2: Backpressure
An API receives requests faster than workers can process them.

Requirements:
- 3 workers.
- Queue capacity = 5.
- Submit must never block.
- If the queue is full, reject the request immediately.
- Every accepted request must eventually be processed.

Task: Implement the worker pool and Submit.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"time"
)

type Request struct {
	ID int
}

type Server struct {
	// TODO
}

func NewServer() *Server {
	// TODO
	return nil
}

func (s *Server) Submit(req Request) bool {
	// TODO
	return false
}

func (s *Server) worker(id int) {
	// TODO
}

func main() {
	server := NewServer()

	for i := 1; i <= 20; i++ {
		accepted := server.Submit(Request{ID: i})

		if accepted {
			fmt.Println("Accepted:", i)
		} else {
			fmt.Println("Rejected:", i)
		}
	}

	time.Sleep(2 * time.Second)
}
```
</details>

### Challenge 3: Graceful Shutdown
You have a worker pool processing requests. The server receives a shutdown signal while requests are still being processed.

Task:
- Implement Shutdown() such that:
- No new work is accepted after shutdown begins.
- Already queued work is processed.
- Workers exit after the queue is drained.
- Shutdown() returns only after all work is complete.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	// TODO
}

func NewServer() *Server {
	// TODO
	return nil
}

func (s *Server) Submit(id int) bool {
	// TODO
	return false
}

func (s *Server) Shutdown() {
	// TODO
}

func (s *Server) worker(id int) {
	// TODO
}

func main() {
	server := NewServer()

	for i := 1; i <= 10; i++ {
		server.Submit(i)
	}

	time.Sleep(300 * time.Millisecond)

	fmt.Println("Shutting down...")
	server.Shutdown()

	fmt.Println("Server stopped")
}
```
</details>