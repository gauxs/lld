package code

import (
	"sync"
	"time"
)

type TimeSlot struct {
	startTime      time.Time
	endTime        time.Time
	meetingID      string
	bookedCapacity int
}

func NewTimeSlot(startTime time.Time, endTime time.Time, meetingID string, participantsCapacity int) *TimeSlot {
	return &TimeSlot{
		startTime:      startTime,
		endTime:        endTime,
		meetingID:      meetingID,
		bookedCapacity: participantsCapacity,
	}
}

type MeetingRoom struct {
	rwMu *sync.RWMutex

	name        string
	capacity    int
	bookedslots []*TimeSlot // NOTE: deadlock is not possible
}

func NewMeetingRoom(name string, cap int) *MeetingRoom {
	return &MeetingRoom{
		rwMu:        &sync.RWMutex{},
		name:        name,
		capacity:    cap,
		bookedslots: make([]*TimeSlot, 0),
	}
}

func (mr *MeetingRoom) IsAvailaible(startTime time.Time, endTime time.Time, participantsCount int) bool {
	if mr.capacity < participantsCount {
		return false
	}

	// check slot availability
	// NOTE: use binary search for optimization
	for _, slot := range mr.bookedslots {
		if slot.startTime.Before(endTime) && startTime.Before(slot.endTime) {
			return false
		}
	}

	return true
}

func (mr *MeetingRoom) BookSlot(meetingID string, startTime time.Time, endTime time.Time, participantsCount int) error {
	mr.bookedslots = append(mr.bookedslots, NewTimeSlot(startTime, endTime, meetingID, participantsCount))
	return nil
}
