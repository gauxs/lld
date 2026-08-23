package ratelimiter

import "sync"

type Storage struct {
	s sync.Map
}

func (s *Storage) Get(key string) int {
	return 0
}

func (s *Storage) CompareAndIncrement(key string, lessThan int) bool {
	return false
}

func (s *Storage) Delete(key string) {
}
