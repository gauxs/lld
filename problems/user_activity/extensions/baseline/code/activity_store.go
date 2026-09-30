package code

import (
	"cmp"
	"slices"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

const (
	DefaultNumberOfMinutesInWindow = 24 * 60
)

type ActivityStore struct {
	window []*Record
}

func NewActivityStore(windowLength uint) *ActivityStore {
	if windowLength == 0 {
		windowLength = DefaultNumberOfMinutesInWindow
	}

	w := make([]*Record, windowLength)

	for i := 0; i < int(windowLength); i++ {
		w[i] = NewRecord()
	}

	return &ActivityStore{
		window: w,
	}
}

func (as *ActivityStore) Store(timeInMin time.Time, userID uint, activity enum.Activity) {
	windowMinute := as.convertTimeTowindowMinute(UnixMinute(timeInMin))

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
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			continue
		}

		if as.window[curWindowMin].byActivity[activity] == nil {
			continue
		}

		totalCount += uint(as.window[curWindowMin].byActivity[activity][userID])
	}

	return totalCount
}

func (as *ActivityStore) GetRateInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) float64 {
	totalCount := as.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity)
	startUnixMin := UnixMinute(startTimeInMin)
	endUnixMin := UnixMinute(endTimeInMin)
	return float64(totalCount) / float64(endUnixMin-startUnixMin+1)
}

func (as *ActivityStore) GetDistinctCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) uint {
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
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			continue
		}
		if as.window[curWindowMin].byActivity[activity] == nil {
			continue
		}

		for userID, _ := range as.window[curWindowMin].byActivity[activity] {
			distinctUserMap[userID] = struct{}{}
		}

	}

	return uint(len(distinctUserMap))
}

func (as *ActivityStore) GetInTimerangeByActivityID(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) []uint {
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
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			continue
		}
		if as.window[curWindowMin].byActivity[activity] == nil {
			continue
		}

		for userID, _ := range as.window[curWindowMin].byActivity[activity] {
			distinctUserMap[userID] = struct{}{}
		}
	}

	distinctUserIDs := make([]uint, 0)
	for userID, _ := range distinctUserMap {
		distinctUserIDs = append(distinctUserIDs, userID)
	}

	return distinctUserIDs
}

func (as *ActivityStore) GetTopKInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity, k uint) []uint {
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
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			continue
		}

		if as.window[curWindowMin].byActivity[activity] == nil {
			continue
		}

		for userID, count := range as.window[curTimeInMin].byActivity[activity] {
			distinctUserCount[userID] += count
		}
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
		if as.window[curWindowMin].unixMinute != curTimeInMin {
			continue
		}
		if as.window[curWindowMin].byUserID[userID] == nil {
			continue
		}

		for activity, _ := range as.window[curWindowMin].byUserID[userID] {
			distinctActivityMap[activity] = struct{}{}
		}
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
