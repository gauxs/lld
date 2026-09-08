package ratelimiter

import "github.com/gauxs/lld/rate_limiter/pkg"

type ResourceIDGenerator interface {
	GetResource(r *pkg.Request) *Resource
}

type ClientAPIResourceGenerator struct {
}

func (carg *ClientAPIResourceGenerator) GetResourceID(r *pkg.Request) *Resource {
	return NewResource(r.GetClientID() + r.GetAPI())
}
