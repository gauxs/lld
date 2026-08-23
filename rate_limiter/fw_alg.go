package ratelimiter

import (
	"github.com/gauxs/lld/rate_limiter/enum"
)

type FixedWindowAlgorithm struct {
}

func (fwa *FixedWindowAlgorithm) HandleRequest(res *Resource, s *Storage) enum.RLStatus {
	return enum.RLStatus_INVALID
}
