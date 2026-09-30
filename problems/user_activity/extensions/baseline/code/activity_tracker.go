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

func (as *ActivityTracker) StoreUserActivity(timeInMin time.Time, userID uint, activity enum.Activity) error {
	if timeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return ErrInvalidTimeRange
	}

	as.store.Store(timeInMin, userID, activity)
	return nil
}

func (as *ActivityTracker) GetActivityCountForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) (uint, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}

	return as.store.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity), nil
}

func (as *ActivityTracker) GetActivityRateForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) (float64, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}
	return as.store.GetRateInTimerange(startTimeInMin, endTimeInMin, userID, activity), nil
}

func (as *ActivityTracker) GetDistinctUserCountByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) (uint, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}
	return as.store.GetDistinctCountInTimerange(startTimeInMin, endTimeInMin, activity), nil
}

func (as *ActivityTracker) GetUsersByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) ([]uint, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return []uint{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}
	return as.store.GetInTimerangeByActivityID(startTimeInMin, endTimeInMin, activity), nil
}

func (as *ActivityTracker) GetTopKUsersForActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity, k uint) ([]uint, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return []uint{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}
	return as.store.GetTopKInTimerange(startTimeInMin, endTimeInMin, activity, k), nil
}

func (as *ActivityTracker) GetUsersActivitySummaryInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) ([]enum.Activity, error) {
	if endTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		return []enum.Activity{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-24 * time.Hour)) {
		startTimeInMin = time.Now().Add(-24 * time.Hour).Add(time.Minute)
	}
	return as.store.GetInTimerangeByUserID(startTimeInMin, endTimeInMin, userID), nil
}
