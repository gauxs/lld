package code

import (
	"sync"
	"time"
)

type TimeSlot struct {
	startTime time.Time
	endTime   time.Time
	meetingID string
}

type MeetingRoom struct {
	rwMu *sync.RWMutex

	name        string
	capacity    int
	bookedslots []*TimeSlot // NOTE: deadlock is not possible
}
