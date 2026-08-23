package ratelimiter

import (
	"sync"
	"time"

	"github.com/gauxs/lld/rate_limiter/enum"
)

type RLAlgorithm interface {
	HandleResource(res *Resource, s *Storage) enum.RLStatus
}

type resourceConfig struct {
	l                   sync.RWMutex
	maxResourceCount    int
	windowDuration      time.Duration
	lastWindowStartTime time.Time
}

type FixedWindowAlgorithm struct {
	config map[string]*resourceConfig
}

func (fwa *FixedWindowAlgorithm) generateKey(resID string, lastWindowStartTime time.Time) string {
	return resID + "#" + lastWindowStartTime.String()
}

func (fwa *FixedWindowAlgorithm) isNewWindow(rc *resourceConfig) bool {
	return time.Now().After(rc.lastWindowStartTime.Add(rc.windowDuration))
}

func (fwa *FixedWindowAlgorithm) getOrGenerateKey(resID string, rc *resourceConfig, s *Storage) string {
	if fwa.isNewWindow(rc) {
		rc.l.RUnlock()
		defer rc.l.Lock()
		rc.l.Lock()
		if rc.lastWindowStartTime.Add(rc.windowDuration).Before(time.Now()) {
			s.Delete(fwa.generateKey(resID, rc.lastWindowStartTime))
			rc.lastWindowStartTime = rc.lastWindowStartTime.Add(rc.windowDuration)
		}
		rc.l.Unlock()
	}

	return fwa.generateKey(resID, rc.lastWindowStartTime)
}

func (fwa *FixedWindowAlgorithm) HandleResource(res *Resource, s *Storage) enum.RLStatus {
	resConfig := fwa.config[res.GetID()]
	resConfig.l.RLock()
	defer resConfig.l.RUnlock()

	if !s.CompareAndIncrement(fwa.getOrGenerateKey(res.GetID(), resConfig, s), resConfig.maxResourceCount) {
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
