// package main

// import connectfour "github.com/gauxs/lld/problems/connect_four/extensions/baseline/code"

// func main() {
// 	connectfour.Execute()
// }

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
