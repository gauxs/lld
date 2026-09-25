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
