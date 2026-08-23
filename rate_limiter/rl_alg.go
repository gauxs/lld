package ratelimiter

import (
	"github.com/gauxs/lld/rate_limiter/enum"
)

type RLAlgorithm interface {
	HandleResource(res *Resource, s *Storage) enum.RLStatus
}
