---
title: Scarcity
sidebar: Scarcity
description: Limited resources—pools, semaphores, and rate limits.
next:
  text: baseline
  link: /problems/connect-four/extensions/baseline/requirements
---

# Scarcity

Scarcity is about managing a limited resource when demand can exceed supply.Examples:
- Only 10 requests should call a downstream API concurrently.
- Only 100 MB of memory should be consumed by active jobs.
- Only 5 DB connections are available.
- A resource is expensive to create, so we should reuse it.

The core question is:
> What happens when all available capacity is already being used?

Unlike correctness, where the concern is data corruption, scarcity is about capacity and resource management.

## The Problem
Consider a service with only 3 database connections.
```text
Request A ──┐
Request B ──┤
Request C ──┼──> Connection Pool (3)
Request D ──┤
Request E ──┘
```

A, B and C can execute immediately. D and E must wait until a connection is returned.

```text
Acquire → Use → Release
            ↓
        connection
        available
```

The important invariant is:
> At most 3 connections can be in use at any time.

The design must answer:
1. How do we limit usage?
2. What happens when capacity is exhausted?
3. How does a waiting operation get capacity?
4. What happens if a caller never releases it?
5. Should callers wait forever, timeout, or fail immediately?

## Solutions
There are two common solutions.
| Mechanism | What it manages |
|---|---|
| **Semaphore** | Permission to use a limited resource |
| **Resource Pool** | The actual reusable resources |

#### Rule of thumb
```text
Need to limit concurrency? → Semaphore
Need an actual reusable object? → Resource Pool
```

### Semaphores
A semaphore maintains a fixed number of permits. For capacity 3:
```text
Semaphore
┌─────────────┐
│ ● ● ●       │  3 permits
└─────────────┘

Acquire → take a permit
Release → return a permit
```
If all permits are taken, Acquire() blocks.

In golang, a buffered channel can act as a counting semaphore. Go's documentation explicitly describes this pattern.

```go
type Semaphore struct {
	permits chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: make(chan struct{}, n),
	}
}

func (s *Semaphore) Acquire() {
	s.permits <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.permits
}
```

Usage
```go
sem := NewSemaphore(3)

sem.Acquire()
defer sem.Release()

doWork()
```
Only 3 goroutines can be inside doWork() concurrently.

#### Why not a mutex?
A mutex says:
> Only one goroutine can enter.

A semaphore says:
> Up to N goroutines can enter.

#### Challenge
What happens if this code panics?
```go
sem.Acquire()

doWork() // panic

sem.Release()
```
The permit is never returned. Use:
```go
sem.Acquire()
defer sem.Release()

doWork()
```
**Lesson: resource acquisition should be paired with guaranteed release.**

### Resource Pooling
A semaphore only gives you permission. A resource pool gives you an actual resource.
```text
Pool
┌─────────────────────┐
│ Conn1 Conn2 Conn3   │
└─────────────────────┘
       ↓
     Get()

Request
       ↓
   use Conn2
       ↓
     Put()
       ↓
      Pool
```

The pool typically contains:
- available resources
- a blocking queue for waiters
- a maximum capacity

#### Simple Go implementation
```go
type Connection struct {
	ID int
}

type Pool struct {
	resources chan *Connection
}

func NewPool(size int) *Pool {
	p := &Pool{
		resources: make(chan *Connection, size),
	}

	for i := 0; i < size; i++ {
		p.resources <- &Connection{ID: i}
	}

	return p
}

func (p *Pool) Get() *Connection {
	return <-p.resources
}

func (p *Pool) Put(conn *Connection) {
	p.resources <- conn
}
```

Usage:
```go
conn := pool.Get()
defer pool.Put(conn)

conn.Query(...)
```

When all connections are occupied, Get() blocks until one becomes available.

### Pool vs Semaphore
```text
Semaphore: Request → Acquire permit → Work → Release permit
Pool: Request → Get connection → Query → Put connection
```

A pool can internally use a semaphore-like mechanism, but the distinction is important:
> Semaphore controls capacity. Pool manages reusable resources.

### Challenges
1. What if a resource is never returned? --> The pool slowly becomes exhausted.

2. What if a resource is broken? --> Returning it blindly can poison the pool.

3. What if requests wait forever? --> Use a timeout/context

```go
select {
case conn := <-p.resources:
	return conn
case <-ctx.Done():
	return nil, ctx.Err()
}
```

Go's concurrency primitives commonly use context to make blocking operations cancellable.

## Common Problems
### Limit Concurrent Operations
A resource can handle only a fixed number of concurrent operations. Use a semaphore:
```go
sem := make(chan struct{}, 3)

sem <- struct{}{}        // acquire
defer func() { <-sem }() // release

process()
```
The semaphore controls how many operations may proceed at once. The key invariant is:
> At most N operations are active at the same time.

### Limit aggregate consumption
Sometimes the scarce resource is not the number of operations, but the amount consumed by those operations. For example, suppose the system has a 100 MB memory budget:
```text
Request A → 20 MB
Request B → 30 MB
Request C → 40 MB

Used = 90 MB

Request D → 20 MB
          ↓
        WAIT
```
A normal semaphore cannot model this because every operation consumes a different amount. Use a weighted semaphore or equivalent accounting mechanism. Conceptually:
```go
limiter.Acquire(ctx, 20)
// use 20 units
limiter.Release(20)
```

### Reuse expensive resources
Some resources are expensive to create. Examples:
- Database connections
- HTTP connections
- Large buffers
- Thread/worker objects
- GPU contexts

Instead of repeatedly creating and destroying them, reuse a bounded pool. The pool controls both:
- resource reuse
- resource capacity

A typical API looks like:
```go
resource := pool.Get()
defer pool.Put(resource)

use(resource)
```
When no resources are available, Get() can block until one is returned.

### Maximize utilization
Having a capacity limit is not enough. You also want to keep available capacity useful. For example:
```text
Capacity = 10 workers

Worker 1  ██████████
Worker 2  ██████████
Worker 3  ██
Worker 4
Worker 5
Worker 6
...
```
There is capacity available, but work may not be distributed efficiently. Common techniques include:

#### Worker Pool
Keep a fixed number of workers consuming from a queue, this avoids creating a new worker for every task.

#### Batching
Instead of repeatedly consuming a resource, combine work. This can reduce per-operation overhead and improve resource utilization.

#### Work Sharing
If one worker has work while another is idle, the idle worker can take work from A.

The general goal is:
> Keep scarce resources busy without exceeding their capacity.

## Exercises
### Exercise 1: Concurrency Limit
A service calls a downstream API that allows at most 3 concurrent requests.

Task: Use a semaphore to enforce the limit.
Requirements:
1. Process 10 requests concurrently.
2. Never have more than 3 requests in flight.
3. Wait for all requests to finish.
4. Always release the permit, even if processing fails.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Problem code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func callDownstream(id int) {
	fmt.Println("Starting request:", id)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Finished request:", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			callDownstream(id)
		}(i)
	}

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

type Semaphore struct {
	permit chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permit: make(chan struct{}, n),
	}
}

func (s *Semaphore) Acquire() {
	s.permit <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.permit
}

func callDownstream(id int) {
	fmt.Println("Starting request:", id)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Finished request:", id)
}

func main() {
	sem := NewSemaphore(3)
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()

			callDownstream(id)
		}(i)
	}

	wg.Wait()
	fmt.Println("Done")
}
```
</details>

### Exercise 2: Resource Pool
A service has 3 expensive database connections that should be reused.

Task: Implement a thread-safe connection pool.

Requirements:
1. The pool starts with 3 connections.
2. Get() returns an available connection.
3. Get() blocks when all connections are in use.
4. Put() returns the connection to the pool.
5. Multiple goroutines can use the pool concurrently.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	ID int
}

type ConnectionPool struct {
	// TODO
}

func NewConnectionPool(size int) *ConnectionPool {
	// TODO
	return nil
}

func (p *ConnectionPool) Get() *Connection {
	// TODO
	return nil
}

func (p *ConnectionPool) Put(conn *Connection) {
	// TODO
}

func main() {
	pool := NewConnectionPool(3)

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			conn := pool.Get()
			fmt.Printf("Request %d using connection %d\n", id, conn.ID)

			time.Sleep(200 * time.Millisecond)

			pool.Put(conn)
		}(i)
	}

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

type Connection struct {
	ID int
}

type ConnectionPool struct {
	pool chan *Connection
}

func NewConnectionPool(size int) *ConnectionPool {
	cp := &ConnectionPool{
		pool: make(chan *Connection, size),
	}

	for i := 0; i < size; i++ {
		cp.pool <- &Connection{ID: i}
	}

	return cp
}

func (p *ConnectionPool) Get() *Connection {
	return <-p.pool
}

func (p *ConnectionPool) Put(conn *Connection) {
	p.pool <- conn
}

func main() {
	pool := NewConnectionPool(3)

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			conn := pool.Get()
			fmt.Printf("Request %d using connection %d\n", id, conn.ID)

			time.Sleep(200 * time.Millisecond)

			pool.Put(conn)
		}(i)
	}

	wg.Wait()
	fmt.Println("Done")
}

```
</details>

### Exercise 3: Weighted Capacity
A service has a maximum 100 MB memory budget for concurrent jobs. Each job requires a different amount of memory.

Task: Implement a weighted limiter.

Requirements:
1. Total acquired memory must never exceed 100 MB.
2. A job blocks when insufficient memory is available.
3. Releasing memory allows waiting jobs to proceed.
4. Multiple goroutines must be supported.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Starter code</summary>

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type MemoryLimiter struct {
	// TODO
}

func NewMemoryLimiter(capacity int) *MemoryLimiter {
	// TODO
	return nil
}

func (l *MemoryLimiter) Acquire(amount int) {
	// TODO
}

func (l *MemoryLimiter) Release(amount int) {
	// TODO
}

func processJob(id, memory int) {
	fmt.Printf("Job %d using %d MB\n", id, memory)
	time.Sleep(200 * time.Millisecond)
}

func main() {
	limiter := NewMemoryLimiter(100)

	jobs := []struct {
		id     int
		memory int
	}{
		{1, 40},
		{2, 30},
		{3, 50},
		{4, 20},
	}

	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)

		go func(id, memory int) {
			defer wg.Done()

			limiter.Acquire(memory)
			defer limiter.Release(memory)

			processJob(id, memory)
		}(job.id, job.memory)
	}

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

type MemoryLimiter struct {
	mu       sync.Mutex
	capCond  *sync.Cond
	capacity int
	current  int
}

func NewMemoryLimiter(capacity int) *MemoryLimiter {
	ml := &MemoryLimiter{
		mu:       sync.Mutex{},
		capacity: capacity,
		current:  0,
	}

	ml.capCond = sync.NewCond(&ml.mu)
	return ml
}

func (l *MemoryLimiter) Acquire(amount int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for (l.current + amount) > l.capacity {
		l.capCond.Wait()
	}

	l.current += amount
}

func (l *MemoryLimiter) Release(amount int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.current -= amount
	l.capCond.Broadcast()
}

func processJob(id, memory int) {
	fmt.Printf("Job %d using %d MB\n", id, memory)
	time.Sleep(200 * time.Millisecond)
}

func main() {
	limiter := NewMemoryLimiter(100)

	jobs := []struct {
		id     int
		memory int
	}{
		{1, 40},
		{2, 30},
		{3, 50},
		{4, 20},
	}

	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)

		go func(id, memory int) {
			defer wg.Done()

			limiter.Acquire(memory)
			defer limiter.Release(memory)

			processJob(id, memory)
		}(job.id, job.memory)
	}

	wg.Wait()
	fmt.Println("Done")
}


```
</details>