package ratelimiter

type Storage struct {
}

func (s *Storage) GetResourceCount(res *Resource) int {
	return 0
}

func (s *Storage) IncrementResourceCount(res *Resource) int {
	return 0
}

func (s *Storage) DecrementResourceCount(res *Resource) int {
	return 0
}
