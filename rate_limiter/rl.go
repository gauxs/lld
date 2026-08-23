package ratelimiter

import (
	"github.com/gauxs/lld/rate_limiter/enum"
	"github.com/gauxs/lld/rate_limiter/pkg"
)

type RateLimiter struct {
	s      *Storage
	alg    RLAlgorithm
	resgen ResourceIDGenerator
}

func (rl *RateLimiter) Handle(req *pkg.Request) enum.RLStatus {
	return enum.RLStatus_INVALID
}

func (rl *RateLimiter) UpdateRLAlgorithm(newAlg RLAlgorithm) {

}
