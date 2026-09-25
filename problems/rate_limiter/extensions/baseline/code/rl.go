package ratelimiter

import (
	"sync"

	"github.com/gauxs/lld/problems/rate_limiter/extensions/baseline/code/enum"
	"github.com/gauxs/lld/problems/rate_limiter/extensions/baseline/code/pkg"
)

type RateLimiter struct {
	l      sync.RWMutex
	s      *Storage
	alg    RLAlgorithm
	resgen ResourceIDGenerator
}

func (rl *RateLimiter) Handle(req *pkg.Request) enum.RLStatus {
	rl.l.RLock()
	alg := rl.alg
	rl.l.RUnlock()

	return alg.HandleResource(rl.resgen.GetResource(req), rl.s)
}

func (rl *RateLimiter) UpdateRLAlgorithm(newAlg RLAlgorithm) {
	// new request will follow this algorithm instantaneously
	// old algorithms bucket will auto cleanup
	rl.l.Lock()
	rl.alg = newAlg
	rl.l.Unlock()
}
