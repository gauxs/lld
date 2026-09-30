package code

import (
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

type ActivityTracker struct {
	store *ActivityStore
}

func NewActivityTracker(windowLen uint) *ActivityTracker {
	return &ActivityTracker{
		store: NewActivityStore(windowLen),
	}
}

func (as *ActivityTracker) StoreUserActivity(timeInMin time.Time, userID uint, activity enum.Activity) {
	as.store.Store(timeInMin, userID, activity)
}

func (as *ActivityTracker) GetActivityCountForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) uint {
	return as.store.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity)
}

func (as *ActivityTracker) GetActivityRateForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) float64 {
	return as.store.GetRateInTimerange(startTimeInMin, endTimeInMin, userID, activity)
}

func (as *ActivityTracker) GetDistinctUserCountByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) uint {
	return as.store.GetDistinctCountInTimerange(startTimeInMin, endTimeInMin, activity)
}

func (as *ActivityTracker) GetUsersByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) []uint {
	return as.store.GetInTimerangeByActivityID(startTimeInMin, endTimeInMin, activity)
}

func (as *ActivityTracker) GetTopKUsersForActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity, k uint) []uint {
	return as.store.GetTopKInTimerange(startTimeInMin, endTimeInMin, activity, k)
}

func (as *ActivityTracker) GetUsersActivitySummaryInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) []enum.Activity {
	return as.store.GetInTimerangeByUserID(startTimeInMin, endTimeInMin, userID)
}
