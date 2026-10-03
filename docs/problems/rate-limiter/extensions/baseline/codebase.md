---
title: Fixed window (in-process) — codebase
problem: rate_limiter
extension: baseline
prev:
  text: design
  link: /problems/rate-limiter/extensions/baseline/design
---

# Codebase

Source: [GitHub](https://github.com/gauxs/lld/tree/main/problems/rate_limiter/extensions/baseline/code)

Path: ` problems/rate_limiter/extensions/baseline/code `

## Directory structure

```text
code/
├── enum/
│   └── rl_status.go
├── pkg/
│   └── req.go
├── ca_resgen.go
├── fw_alg.go
├── res.go
├── rl.go
└── storage.go
```

## ` ca_resgen.go `

```go
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
```

## ` enum/rl_status.go `

```go
package enum

type RLStatus int

const (
	RLStatus_INVALID RLStatus = iota
	RLStatus_ACCEPTED
	RLStatus_REJECTED
)
```

## ` fw_alg.go `

```go
package ratelimiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/gauxs/lld/problems/rate_limiter/extensions/baseline/code/enum"
)

type RLAlgorithm interface {
	HandleResource(res *Resource, s *Storage) enum.RLStatus
}

type resourceConfig struct {
	l                sync.RWMutex
	versNo           int
	maxResourceCount int
	windowDuration   time.Duration
}

type FixedWindowAlgorithm struct {
	config sync.Map
}

func (fwa *FixedWindowAlgorithm) generateKey(versionNo int, resID string, windowNo int64) string {
	return "FWA" + strconv.Itoa(versionNo) + "#" + resID + "#" + strconv.Itoa(int(windowNo))
}

func (fwa *FixedWindowAlgorithm) currentWindow(rc *resourceConfig) int64 {
	// epoch time divided into duration sized chuncks. Each chunck is termed as
	// window number
	return int64(time.Now().UnixNano() / rc.windowDuration.Nanoseconds())
}

func (fwa *FixedWindowAlgorithm) currentWindowExpiry(rc *resourceConfig) time.Duration {
	curWindow := fwa.currentWindow(rc)
	return time.Duration((curWindow+1)*(int64(rc.windowDuration.Nanoseconds()))-time.Now().UnixNano()) - 1
}

func (fwa *FixedWindowAlgorithm) HandleResource(res *Resource, s *Storage) enum.RLStatus {
	v, ok := fwa.config.Load(res.GetID())
	if !ok {
		// allow incase of configuration issues
		return enum.RLStatus_ACCEPTED
	}

	rs, ok := v.(*resourceConfig)
	if !ok {
		// allow incase of configuration issues
		return enum.RLStatus_ACCEPTED
	}

	rs.l.RLock()
	defer rs.l.RUnlock()
	cw := fwa.currentWindow(rs)
	if !s.CompareAndIncrement(fwa.generateKey(rs.versNo, res.GetID(), cw),
		rs.maxResourceCount, fwa.currentWindowExpiry(rs)) {
		return enum.RLStatus_REJECTED
	}

	return enum.RLStatus_ACCEPTED
}

func (fwa *FixedWindowAlgorithm) UpdateWindowDuration(res *Resource, newWD time.Duration) {
	v, ok := fwa.config.Load(res.GetID())
	if !ok {
		return
	}

	rs, ok := v.(*resourceConfig)
	if !ok {
		return
	}

	rs.l.Lock()
	defer rs.l.Unlock()

	rs.versNo += 1
	rs.windowDuration = newWD
}

func (fwa *FixedWindowAlgorithm) UpdateMaxResourceCount(res *Resource, newmrc int) {
	v, ok := fwa.config.Load(res.GetID())
	if !ok {
		return
	}

	rs, ok := v.(*resourceConfig)
	if !ok {
		return
	}

	rs.l.Lock()
	defer rs.l.Unlock()

	rs.versNo += 1
	rs.maxResourceCount = newmrc
}
```

## ` pkg/req.go `

```go
package pkg

type Request struct {
	clientID string
	api      string
}

func NewRequest(cID string, a string) *Request {
	return &Request{
		clientID: cID,
		api:      a,
	}
}

func (r *Request) GetClientID() string {
	return r.clientID
}

func (r *Request) GetAPI() string {
	return r.api
}
```

## ` res.go `

```go
package ratelimiter

type Resource struct {
	id string
}

func NewResource(id string) *Resource {
	return &Resource{
		id: id,
	}
}

func (r *Resource) GetID() string {
	return r.id
}
```

## ` rl.go `

```go
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
```

## ` storage.go `

```go
package ratelimiter

import (
	"sync"
	"time"
)

type ds struct {
	l     sync.RWMutex
	key   string
	count int
}

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{
		m: sync.Map{},
	}
}

func (s *Storage) CompareAndIncrement(key string, lessThan int, expiry time.Duration) bool {
	if actual, loaded := s.m.LoadOrStore(key, &ds{
		sync.RWMutex{},
		key,
		1,
	}); loaded {
		assertedV, _ := actual.(*ds)
		assertedV.l.Lock()
		defer assertedV.l.Unlock()

		if assertedV.count < lessThan {
			assertedV.count++
			return true
		}

		return false
	} else {
		s.DeleteAfterDuration(key, expiry)
	}

	return true
}

func (s *Storage) DeleteAfterDuration(key string, expiry time.Duration) {
	// Automatically runs the function in its own goroutine after the duration
	time.AfterFunc(expiry, func() {
		s.m.Delete(key)
	})
}
```
