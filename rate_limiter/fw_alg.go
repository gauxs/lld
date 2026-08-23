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

func (fwa *FixedWindowAlgorithm) generateKey(res *Resource) string {
	// resourceID + lastWindowStartTime
	return ""
}

func (fwa *FixedWindowAlgorithm) isNewWindow(rc *resourceConfig) string {
	return ""
}

func (fwa *FixedWindowAlgorithm) generateNewWindow(rc *resourceConfig) string {
	// cleanup last window here
	return ""
}

func (fwa *FixedWindowAlgorithm) HandleResource(res *Resource, s *Storage) enum.RLStatus {
	return enum.RLStatus_INVALID
}

func (fwa *FixedWindowAlgorithm) UpdateWindowDuration(res *Resource, newWD time.Duration) {

}

func (fwa *FixedWindowAlgorithm) UpdateMaxResourceCount(res *Resource, newmrc int) {

}

func (fwa *FixedWindowAlgorithm) UpdateLastWindowStartTime(res *Resource, newlwst time.Time) {

}
