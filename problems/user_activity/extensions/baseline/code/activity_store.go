package code

import (
	"cmp"
	"slices"
	"sync"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

const (
	DefaultNumberOfMinutesInWindow = 24 * 60
)

type ActivityStore struct {
	windowMux []*sync.RWMutex
	window    []*Record
}

func NewActivityStore(windowLength uint) *ActivityStore {
	wm := make([]*sync.RWMutex, windowLength)
	w := make([]*Record, windowLength)

	for i := 0; i < int(windowLength); i++ {
		wm[i] = &sync.RWMutex{}
		w[i] = NewRecord()
	}

	return &ActivityStore{
		windowMux: wm,
		window:    w,
	}
}

func (as *ActivityStore) Store(timeInMin time.Time, userID uint, activity enum.Activity) {
	windowMinute := as.convertTimeTowindowMinute(UnixMinute(timeInMin))
	as.windowMux[windowMinute].Lock()
	defer as.windowMux[windowMinute].Unlock()

	if as.window[windowMinute].unixMinute != UnixMinute(timeInMin) {
		r := NewRecord()
		r.unixMinute = UnixMinute(timeInMin)
		as.window[windowMinute] = r
	}

	if as.window[windowMinute].byActivity[activity] == nil {
		as.window[windowMinute].byActivity[activity] = make(map[uint]uint)
	}

	as.window[windowMinute].byActivity[activity][userID]++

	if as.window[windowMinute].byUserID[userID] == nil {
		as.window[windowMinute].byUserID[userID] = make(map[enum.Activity]uint)
	}

	as.window[windowMinute].byUserID[userID][activity]++
}

func (as *ActivityStore) GetCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) uint {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	if startTimeInMin.After(endTimeInMin) {
		return 0
	}

	if userID == 0 {
		return 0
	}

	if activity == enum.UNDEFINED_ACTIVITY {
		return 0
	}

	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)

	totalCount := uint(0)
	for curTimeInMin := startUnixMin; curTimeInMin <= endUnixMin; curTimeInMin++ {
		curWindowMin := as.convertTimeTowindowMinute(curTimeInMin)
		as.windowMux[curWindowMin].RLock()

		if as.window[curWindowMin].unixMinute != curTimeInMin {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		if as.window[curWindowMin].byActivity[activity] == nil {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		totalCount += uint(as.window[curWindowMin].byActivity[activity][userID])
		as.windowMux[curWindowMin].RUnlock()
	}

	return totalCount
}

func (as *ActivityStore) GetRateInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) float64 {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	totalCount := as.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity)
	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)
	return float64(totalCount) / float64(endUnixMin-startUnixMin+1)
}

func (as *ActivityStore) GetDistinctCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) uint {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	if startTimeInMin.After(endTimeInMin) {
		return 0
	}

	if activity == enum.UNDEFINED_ACTIVITY {
		return 0
	}

	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)

	distinctUserMap := make(map[uint]struct{})
	for curTimeInMin := startUnixMin; curTimeInMin <= endUnixMin; curTimeInMin++ {
		curWindowMin := as.convertTimeTowindowMinute(curTimeInMin)
		as.windowMux[curWindowMin].RLock()
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}
		if as.window[curWindowMin].byActivity[activity] == nil {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		for userID, _ := range as.window[curWindowMin].byActivity[activity] {
			distinctUserMap[userID] = struct{}{}
		}
		as.windowMux[curWindowMin].RUnlock()
	}

	return uint(len(distinctUserMap))
}

func (as *ActivityStore) GetInTimerangeByActivityID(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) []uint {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	if startTimeInMin.After(endTimeInMin) {
		return []uint{}
	}

	if activity == enum.UNDEFINED_ACTIVITY {
		return []uint{}
	}

	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)

	distinctUserMap := make(map[uint]struct{})
	for curTimeInMin := startUnixMin; curTimeInMin <= endUnixMin; curTimeInMin++ {
		curWindowMin := as.convertTimeTowindowMinute(curTimeInMin)
		as.windowMux[curWindowMin].RLock()
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}
		if as.window[curWindowMin].byActivity[activity] == nil {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		for userID, _ := range as.window[curWindowMin].byActivity[activity] {
			distinctUserMap[userID] = struct{}{}
		}
		as.windowMux[curWindowMin].RUnlock()
	}

	distinctUserIDs := make([]uint, 0)
	for userID, _ := range distinctUserMap {
		distinctUserIDs = append(distinctUserIDs, userID)
	}

	return distinctUserIDs
}

func (as *ActivityStore) GetTopKInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity, k uint) []uint {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	if startTimeInMin.After(endTimeInMin) {
		return []uint{}
	}

	if activity == enum.UNDEFINED_ACTIVITY {
		return []uint{}
	}

	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)

	distinctUserCount := make(map[uint]uint)
	for curTimeInMin := startUnixMin; curTimeInMin <= endUnixMin; curTimeInMin++ {
		curWindowMin := as.convertTimeTowindowMinute(curTimeInMin)
		as.windowMux[curWindowMin].RLock()
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		if as.window[curWindowMin].byActivity[activity] == nil {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		for userID, count := range as.window[curWindowMin].byActivity[activity] {
			distinctUserCount[userID] += count
		}
		as.windowMux[curWindowMin].RUnlock()
	}

	type userAndCount struct {
		userID uint
		count  uint
	}

	distinctUserIDs := make([]userAndCount, 0, len(distinctUserCount))
	for userID, count := range distinctUserCount {
		distinctUserIDs = append(distinctUserIDs, userAndCount{
			userID: userID,
			count:  count,
		})
	}

	slices.SortFunc(distinctUserIDs, func(a, b userAndCount) int {
		return cmp.Compare(b.count, a.count)
	})

	limit := int(k)
	if limit > len(distinctUserIDs) {
		limit = len(distinctUserIDs)
	}

	// 4. Extract just the userIDs for the top K results
	topKUsers := make([]uint, 0, limit)
	for i := 0; i < limit; i++ {
		topKUsers = append(topKUsers, distinctUserIDs[i].userID)
	}

	return topKUsers
}

func (as *ActivityStore) GetInTimerangeByUserID(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) []enum.Activity {
	if endTimeInMin.After(time.Now()) {
		endTimeInMin = time.Now()
	}

	if startTimeInMin.After(endTimeInMin) {
		return []enum.Activity{}
	}

	if userID == 0 {
		return []enum.Activity{}
	}

	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)

	distinctActivityMap := make(map[enum.Activity]struct{})
	for curTimeInMin := startUnixMin; curTimeInMin <= endUnixMin; curTimeInMin++ {
		curWindowMin := as.convertTimeTowindowMinute(curTimeInMin)
		as.windowMux[curWindowMin].RLock()
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}
		if as.window[curWindowMin].byUserID[userID] == nil {
			as.windowMux[curWindowMin].RUnlock()
			continue
		}

		for activity, _ := range as.window[curWindowMin].byUserID[userID] {
			distinctActivityMap[activity] = struct{}{}
		}
		as.windowMux[curWindowMin].RUnlock()
	}

	distinctActivities := make([]enum.Activity, 0)
	for activity, _ := range distinctActivityMap {
		distinctActivities = append(distinctActivities, activity)
	}

	return distinctActivities
}

func (as *ActivityStore) convertTimeTowindowMinute(unixMin int64) uint {
	return uint(unixMin % int64(len(as.window)))
}

func UnixMinute(t time.Time) int64 {
	return t.Unix() / 60
}

type Record struct {
	unixMinute int64
	byActivity map[enum.Activity]map[uint]uint
	byUserID   map[uint]map[enum.Activity]uint
}

func NewRecord() *Record {
	return &Record{
		unixMinute: 0,
		byActivity: make(map[enum.Activity]map[uint]uint),
		byUserID:   make(map[uint]map[enum.Activity]uint),
	}
}
