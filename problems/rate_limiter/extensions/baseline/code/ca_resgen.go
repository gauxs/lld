package ratelimiter

import "github.com/gauxs/lld/problems/rate_limiter/extensions/baseline/code/pkg"

type ResourceIDGenerator interface {
	GetResource(r *pkg.Request) *Resource
}

type ClientAPIResourceGenerator struct {
}

func (carg *ClientAPIResourceGenerator) GetResource(r *pkg.Request) *Resource {
	return NewResource(r.GetClientID() + r.GetAPI())
}
