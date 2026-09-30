package code

import (
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

type ActivityTracker struct {
}

func (as *ActivityTracker) StoreUserActivity(timeInMin time.Time, userID uint, activity enum.Activity) {
	return
}

func (as *ActivityTracker) GetActivityCountForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) {

}

func (as *ActivityTracker) GetActivityRateForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) {

}

func (as *ActivityTracker) GetDistinctUserCountByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityTracker) GetUsersByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityTracker) GetTopKUsersForActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityTracker) GetUsersActivitySummaryInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) {

}
