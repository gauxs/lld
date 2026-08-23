package ratelimiter

import "github.com/gauxs/lld/rate_limiter/pkg"

type ClientAPIResourceGenerator struct {
}

func (carg *ClientAPIResourceGenerator) GetResourceID(r *pkg.Request) string {
	return ""
}
