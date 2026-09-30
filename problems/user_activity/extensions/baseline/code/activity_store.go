package code

import (
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

type ActivityStore struct {
}

func (as *ActivityStore) Store(timeInMin time.Time, userID uint, activity enum.Activity) {
	return
}

func (as *ActivityStore) GetCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) {

}

func (as *ActivityStore) GetRateInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) {

}

func (as *ActivityStore) GetDistinctCountInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetInTimerangeByActivityID(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetTopKInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) {

}

func (as *ActivityStore) GetInTimerangeByUserID(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) {

}
