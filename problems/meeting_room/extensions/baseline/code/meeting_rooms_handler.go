package code

import (
	"sync"
	"time"
)

type MeetingRoomsHandler struct {
	rooms sync.Map
}

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string) error {
	return nil
}

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) error {
	return nil
}

func (mrh *MeetingRoomsHandler) ReserveMeetingRoom(name string, meetingID string, startTime time.Time, endTime time.Time) error {
	return nil
}
