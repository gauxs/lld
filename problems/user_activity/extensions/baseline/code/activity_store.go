package code

import (
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
	windowMinute := as.convertTimeTowindowMinute(timeInMin)

	if as.window[windowMinute].byActivity[activity] == nil {
		as.window[windowMinute].byActivity[activity] = make(map[uint]uint)
	}

	as.window[windowMinute].byActivity[activity][userID]++
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

	startTimeWindowMin := as.convertTimeTowindowMinute(startTimeInMin)
	endTimeWindowMin := as.convertTimeTowindowMinute(endTimeInMin)

	totalCount := uint(0)
	for curTimeInMin := startTimeWindowMin; curTimeInMin <= endTimeWindowMin; curTimeInMin++ {
		if as.window[curTimeInMin].byActivity[activity] == nil {
			continue
		}

		totalCount += uint(as.window[curTimeInMin].byActivity[activity][userID])
	}

	return totalCount
}

func (as *ActivityStore) GetRateInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) float64 {
	totalCount := as.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity)
	startTimeWindowMin := as.convertTimeTowindowMinute(startTimeInMin)
	endTimeWindowMin := as.convertTimeTowindowMinute(endTimeInMin)
	return float64(totalCount) / float64(endTimeWindowMin-startTimeWindowMin+1)
}

func (as *ActivityStore) GetDistinctCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetInTimerangeByActivityID(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetTopKInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetInTimerangeByUserID(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) {

}

func (as *ActivityStore) convertTimeTowindowMinute(t time.Time) uint {
	return uint(t.Unix() / 60 % int64(len(as.window)))
}

type Record struct {
	// activity X user - count
	byActivity map[enum.Activity]map[uint]uint
}

func NewRecord() *Record {
	return &Record{
		byActivity: make(map[enum.Activity]map[uint]uint),
	}
}
