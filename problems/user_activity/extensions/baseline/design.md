# User activity tracker — baseline (design)

## Approach

`ActivityTracker` is the public facade over a rolling array of minute buckets.
Each bucket indexes counts both by activity and by user, while independent
locks keep concurrent reads and writes scoped to the affected minute.

## Go design sketch

Method bodies are intentionally omitted. Comments describe ownership,
responsibilities, and invariants.

```go
import (
    "sync"
    "time"

    "github.com/gauxs/lld/problems/user_activity/extensions/baseline/code/enum"
)

// Activitystore — Store all the activities with time as the first class dimensions.
// Activitystore stores Record.
type ActivityStore struct {
    windowMux []*sync.RWMutex
    window    []*Record
}

// Record is one minute bucket indexed both by activity and by user.
type Record struct {
    unixMinute int64
    byActivity map[enum.Activity]map[uint]uint
    byUserID   map[uint]map[enum.Activity]uint
}

// User — An user.
// User represents a user.
type User struct {
    ID uint
}

// ActivityTracker — Holds ActivityStore and provides activity tracking functionality.
// ActivityTracker holds ActivityStore.
type ActivityTracker struct {
    windowLenInMin time.Duration
    store          *ActivityStore
}

// Store record for timestamp, user and activity.
func (s *ActivityStore) Store(timestamp time.Time, userID uint, activityID enum.Activity) error

// Get the count of record in time range by userID and activityID.
func (s *ActivityStore) GetCountInTimerange(startTimestamp, endTimestamp time.Time, userID uint, activityID enum.Activity) uint

// Get the rate in time range by userID and activityID.
func (s *ActivityStore) GetRateInTimerange(startTimestamp, endTimestamp time.Time, userID uint, activityID enum.Activity) float64

// Get the distinct count of record in time range by activityID.
func (s *ActivityStore) GetDistinctCountInTimerange(startTimestamp, endTimestamp time.Time, activityID enum.Activity) uint

// Get user IDs that performed activityID in the time range.
func (s *ActivityStore) GetInTimerangeByActivityID(startTimestamp, endTimestamp time.Time, activityID enum.Activity) []uint

// Get top K frequency records by userID in time range for activityID.
func (s *ActivityStore) GetTopKInTimerange(startTimestamp, endTimestamp time.Time, activityID enum.Activity, k uint) []uint

// Get activity counts for userID in the time range.
func (s *ActivityStore) GetInTimerangeByUserID(startTimestamp, endTimestamp time.Time, userID uint) map[enum.Activity]uint

// Store users activity for timestamp.
func (t *ActivityTracker) StoreUserActivity(timestamp time.Time, userID uint, activityID enum.Activity) error

// Get activity count for user in time range.
func (t *ActivityTracker) GetActivityCountForUserInTimerange(startTimestamp, endTimestamp time.Time, userID uint, activityID enum.Activity) (uint, error)

// Get activity rate for user in time range.
func (t *ActivityTracker) GetActivityRateForUserInTimerange(startTimestamp, endTimestamp time.Time, userID uint, activityID enum.Activity) (float64, error)

// Get distinct user count by activity in time range.
func (t *ActivityTracker) GetDistinctUserCountByActivityInTimerange(startTimestamp, endTimestamp time.Time, activityID enum.Activity) (uint, error)

// Get users by activity in time range.
func (t *ActivityTracker) GetUsersByActivityInTimerange(startTimestamp, endTimestamp time.Time, activityID enum.Activity) ([]uint, error)

// Get topK users for activity in time range.
func (t *ActivityTracker) GetTopKUsersForActivityInTimerange(startTimestamp, endTimestamp time.Time, activityID enum.Activity, k uint) ([]uint, error)

// Get users activity counts grouped by activity type in time range.
func (t *ActivityTracker) GetUsersActivitySummaryInTimerange(startTimestamp, endTimestamp time.Time, userID uint) (map[enum.Activity]uint, error)
```

## Implementation note

The design returns counts grouped by activity, as required. The current
reference implementation returns only `[]enum.Activity`; it must be updated to
match this contract.
