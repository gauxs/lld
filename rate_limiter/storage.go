package ratelimiter

import "sync"

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

func (s *Storage) CompareAndIncrement(key string, lessThan int) bool {
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
	}

	return true
}

func (s *Storage) Delete(key string) {
	s.m.Delete(key)
}

func (s *Storage) GC() {
	// run GC every duration and clear
}
