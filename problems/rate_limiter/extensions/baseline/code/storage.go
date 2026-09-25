package ratelimiter

import (
	"sync"
	"time"
)

type ds struct {
	l     sync.RWMutex
	key   string
	count int
}

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{
		m: sync.Map{},
	}
}

func (s *Storage) CompareAndIncrement(key string, lessThan int, expiry time.Duration) bool {
	if actual, loaded := s.m.LoadOrStore(key, &ds{
		sync.RWMutex{},
		key,
		1,
	}); loaded {
		assertedV, _ := actual.(*ds)
		assertedV.l.Lock()
		defer assertedV.l.Unlock()

		if assertedV.count < lessThan {
			assertedV.count++
			return true
		}

		return false
	} else {
		s.DeleteAfterDuration(key, expiry)
	}

	return true
}

func (s *Storage) DeleteAfterDuration(key string, expiry time.Duration) {
	// Automatically runs the function in its own goroutine after the duration
	time.AfterFunc(expiry, func() {
		s.m.Delete(key)
	})
}
