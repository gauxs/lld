package iface

import (
	"github.com/gauxs/lld/rate_limiter/pkg"
)

type ResourceIDGenerator interface {
	GetResourceID(r *pkg.Request) string
}
