// package main

// import connectfour "github.com/gauxs/lld/problems/connect_four/extensions/baseline/code"

//	func main() {
//		connectfour.Execute()
//	}
package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	balance int64
}

func transfer(from, to *Account, amount int64, mu *sync.Mutex) bool {
	mu.Lock()
	defer mu.Unlock()

	if from.balance < amount {
		return false
	}
	from.balance -= amount

	// to bring out the race issue
	time.Sleep(1 * time.Microsecond)

	to.balance += amount
	return true
}

func main() {
	a := Account{balance: 1000}
	b := Account{balance: 1000}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 500; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			transfer(&a, &b, 1, &mu)
		}()
		go func() {
			defer wg.Done()
			transfer(&b, &a, 1, &mu)
		}()
	}

	wg.Wait()
	total := a.balance + b.balance
	fmt.Println("Expected total:", 2000)
	fmt.Println("Actual total:", total)
}
