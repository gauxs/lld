package ratelimiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/gauxs/lld/rate_limiter/enum"
)

type RLAlgorithm interface {
	HandleResource(res *Resource, s *Storage) enum.RLStatus
}

type resourceConfig struct {
	l                sync.RWMutex
	maxResourceCount int
	windowDuration   time.Duration
}

type FixedWindowAlgorithm struct {
	config map[string]*resourceConfig
}

func (fwa *FixedWindowAlgorithm) generateKey(resID string, windowNo int) string {
	return resID + "#" + strconv.Itoa(windowNo)
}

func (fwa *FixedWindowAlgorithm) currentWindow(rc *resourceConfig) int {
	return int(float64(time.Now().UnixNano()) / float64(rc.windowDuration.Nanoseconds()))
}

func (fwa *FixedWindowAlgorithm) HandleResource(res *Resource, s *Storage) enum.RLStatus {
	resConfig := fwa.config[res.GetID()]
	resConfig.l.RLock()
	defer resConfig.l.RUnlock()

	if !s.CompareAndIncrement(fwa.generateKey(res.GetID(), fwa.currentWindow(resConfig)), resConfig.maxResourceCount) {
		return enum.RLStatus_REJECTED
	}

	return enum.RLStatus_ACCEPTED
}

func (fwa *FixedWindowAlgorithm) UpdateWindowDuration(res *Resource, newWD time.Duration) {
	resConfig := fwa.config[res.GetID()]
	resConfig.l.Lock()
	defer resConfig.l.Unlock()

	resConfig.windowDuration = newWD
}

func (fwa *FixedWindowAlgorithm) UpdateMaxResourceCount(res *Resource, newmrc int) {
	resConfig := fwa.config[res.GetID()]
	resConfig.l.Lock()
	defer resConfig.l.Unlock()

	resConfig.maxResourceCount = newmrc
}
