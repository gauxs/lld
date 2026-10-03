---
title: User Activity Rate Tracker — codebase
problem: user_activity
extension: baseline
prev:
  text: design
  link: /problems/user-activity/extensions/baseline/design
---

# Codebase

Source: [GitHub](https://github.com/gauxs/lld/tree/main/problems/user_activity/extensions/baseline/code)

Path: ` problems/user_activity/extensions/baseline/code `

## Directory structure

```text
code/
├── enum/
│   └── activity.go
├── activity_store.go
├── activity_store_test.go
├── activity_tracker.go
├── error.go
├── user.go
└── user_activity.go
```

## ` activity_store.go `

```go
package code

import (
	"cmp"
	"fmt"
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

func (as *ActivityStore) Store(timeInMin time.Time, userID uint, activity enum.Activity) error {
	if timeInMin.After(time.Now()) {
		return ErrInvalidTime
	}

	if userID == 0 {
		return fmt.Errorf("%w %d", ErrInvalidUser, userID)
	}

	if activity == enum.UNDEFINED_ACTIVITY {
		return fmt.Errorf("%w %d", ErrInvalidActivity, activity)
	}

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
	return nil
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
```

## ` activity_store_test.go `

```go
package code

import (
	"testing"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

var (
	OneHrBack        = time.Now().Add(-time.Hour)
	OneAndHalfHrBack = time.Now().Add(-1 * time.Hour).Add(-30 * time.Minute)
	TwoHrBack        = time.Now().Add(-2 * time.Hour)
)

func TestStore_InvalidInput(t *testing.T) {
	as := NewActivityStore(DefaultNumberOfMinutesInWindow)
	err := as.Store(time.Now(), 0, enum.LOGIN)

	if err == nil {
		t.Errorf("%s | storing activity should not have been allowed", t.Name())
	}
}

func TestStore_GetCountInTimerange(t *testing.T) {
	as := NewActivityStore(DefaultNumberOfMinutesInWindow)

	userID := uint(1)
	activity := enum.LOGIN
	if err := as.Store(OneAndHalfHrBack, userID, activity); err != nil {
		t.Errorf("%s | error in storing user %d activity %v at time %v", t.Name(), userID, activity, OneAndHalfHrBack)
	}

	if err := as.Store(OneAndHalfHrBack.Add(time.Minute), userID, activity); err != nil {
		t.Errorf("%s | error in storing user %d activity %v at time %v", t.Name(), userID, activity, OneAndHalfHrBack)
	}

	expectedActivityCount := uint(2)
	if count := as.GetCountInTimerange(TwoHrBack, OneHrBack, userID, activity); count != expectedActivityCount {
		t.Errorf("%s | unexpected activity count for user %d activity %v between time %v - %v. Expected %v, got %v",
			t.Name(), userID, activity, TwoHrBack, OneHrBack, expectedActivityCount, count)
	}
}
```

## ` activity_tracker.go `

```go
package code

import (
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

type ActivityTracker struct {
	windowLenInMin time.Duration
	store          *ActivityStore
}

func NewActivityTracker(windowLen uint) *ActivityTracker {
	if windowLen == 0 {
		windowLen = DefaultNumberOfMinutesInWindow
	}

	return &ActivityTracker{
		windowLenInMin: time.Duration(windowLen) * time.Minute,
		store:          NewActivityStore(windowLen),
	}
}

func (as *ActivityTracker) StoreUserActivity(timeInMin time.Time, userID uint, activity enum.Activity) error {
	if timeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return ErrInvalidTimeRange
	}

	as.store.Store(timeInMin, userID, activity)
	return nil
}

func (as *ActivityTracker) GetActivityCountForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) (uint, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}

	return as.store.GetCountInTimerange(startTimeInMin, endTimeInMin, userID, activity), nil
}

func (as *ActivityTracker) GetActivityRateForUserInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint, activity enum.Activity) (float64, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}
	return as.store.GetRateInTimerange(startTimeInMin, endTimeInMin, userID, activity), nil
}

func (as *ActivityTracker) GetDistinctUserCountByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) (uint, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return 0, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}
	return as.store.GetDistinctCountInTimerange(startTimeInMin, endTimeInMin, activity), nil
}

func (as *ActivityTracker) GetUsersByActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity) ([]uint, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return []uint{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}
	return as.store.GetInTimerangeByActivityID(startTimeInMin, endTimeInMin, activity), nil
}

func (as *ActivityTracker) GetTopKUsersForActivityInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, activity enum.Activity, k uint) ([]uint, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return []uint{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}
	return as.store.GetTopKInTimerange(startTimeInMin, endTimeInMin, activity, k), nil
}

func (as *ActivityTracker) GetUsersActivitySummaryInTimerange(startTimeInMin time.Time, endTimeInMin time.Time, userID uint) ([]enum.Activity, error) {
	if endTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		return []enum.Activity{}, ErrInvalidTimeRange
	}

	if startTimeInMin.Before(time.Now().Add(-as.windowLenInMin)) {
		startTimeInMin = time.Now().Add(-as.windowLenInMin).Add(time.Minute)
	}
	return as.store.GetInTimerangeByUserID(startTimeInMin, endTimeInMin, userID), nil
}
```

## ` enum/activity.go `

```go
package enum

type Activity int

const (
	UNDEFINED_ACTIVITY Activity = iota
	LOGIN
	LOGOUT
)
```

## ` error.go `

```go
package code

import "errors"

var ErrInvalidTimeRange = errors.New("invalid time range")
var ErrInvalidTime = errors.New("invalid time")
var ErrInvalidUser = errors.New("invalid user")
var ErrInvalidActivity = errors.New("invalid activity")
```

## ` user.go `

```go
package code

type User struct {
	userID uint
}

func NewUser(userID uint) *User {
	return &User{
		userID: userID,
	}
}

func (u *User) GetUserID() uint {
	if u == nil {
		// zero is invalid userID
		return 0
	}

	return u.userID
}
```

## ` user_activity.go `

```go
package code

import (
	"fmt"
	"time"

	"github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

func Execute() {
	activityTracker := NewActivityTracker(24 * 60)
	for {
		fmt.Println("\n--- Select Operation ---")
		fmt.Println("1. Store User Activity")
		fmt.Println("2. Get Activity Count For User In Timerange")
		fmt.Println("3. Get Activity Rate For User In Timerange")
		fmt.Println("4. Get Distinct User Count By Activity In Timerange")
		fmt.Println("5. Get Users By Activity In Timerange")
		fmt.Println("6. Get Top K Users For Activity In Timerange")
		fmt.Println("7. Get Users Activity Summary In Timerange")
		fmt.Println("8. EXIT")
		fmt.Print("Enter option (1-8): ")

		var op int
		_, err := fmt.Scan(&op)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch op {
		case 1: // StoreUserActivity
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <userID> and <activityID> (1=LOGIN, 2=LOGOUT):")
			_, err := fmt.Scan(&userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			currentTime := time.Now()
			activityTracker.StoreUserActivity(currentTime, userID, enum.Activity(activityID))
			fmt.Printf("Stored activity %d for user %d at %s\n", activityID, userID, currentTime.Format("15:04"))

		case 2: // GetActivityCountForUserInTimerange
			var startOffset, endOffset int
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID> <activityID>:")
			fmt.Println("(e.g., '-10 0 1 1' checks from 10 mins ago until now)")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if count, err := activityTracker.GetActivityCountForUserInTimerange(startTime, endTime, userID, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Resulting Count: %d\n", count)
			}

		case 3: // GetActivityRateForUserInTimerange
			var startOffset, endOffset int
			var userID uint
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if rate, err := activityTracker.GetActivityRateForUserInTimerange(startTime, endTime, userID, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Resulting Rate: %.4f\n", rate)
			}

		case 4: // GetDistinctUserCountByActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if distinctCount, err := activityTracker.GetDistinctUserCountByActivityInTimerange(startTime, endTime, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Distinct User Count: %d\n", distinctCount)
			}

		case 5: // GetUsersByActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if users, err := activityTracker.GetUsersByActivityInTimerange(startTime, endTime, enum.Activity(activityID)); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Matching User IDs: %v\n", users)
			}

		case 6: // GetTopKUsersForActivityInTimerange
			var startOffset, endOffset int
			var activityID int
			var k uint
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <activityID> <k>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &activityID, &k)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if topUsers, err := activityTracker.GetTopKUsersForActivityInTimerange(startTime, endTime, enum.Activity(activityID), k); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Top %d Users: %v\n", k, topUsers)
			}

		case 7: // GetUsersActivitySummaryInTimerange
			var startOffset, endOffset int
			var userID uint
			fmt.Println("Enter space-separated <start_min_offset> <end_min_offset> <userID>:")
			_, err := fmt.Scan(&startOffset, &endOffset, &userID)
			if err != nil {
				fmt.Println("Error reading parameters.")
				continue
			}
			startTime := time.Now().Add(time.Duration(startOffset) * time.Minute)
			endTime := time.Now().Add(time.Duration(endOffset) * time.Minute)

			if summary, err := activityTracker.GetUsersActivitySummaryInTimerange(startTime, endTime, userID); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Activity History Summary: %v\n", summary)
			}

		case 8: // Exit
			fmt.Println("Exiting application...")
			return

		default:
			fmt.Println("Unknown operation number. Please select between 1 and 8.")
		}
	}
}
```
