package ratelimiter

import "github.com/gauxs/lld/rate_limiter/pkg"

type ResourceIDGenerator interface {
	GetResourceID(r *pkg.Request) string
}

type ClientAPIResourceGenerator struct {
}

func (carg *ClientAPIResourceGenerator) GetResourceID(r *pkg.Request) string {
	return ""
}
