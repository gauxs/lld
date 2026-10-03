package code

import (
	"sync"
	"time"
)

type MeetingRoomsHandler struct {
	rooms sync.Map
}

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string, capacity int) error {
	mrh.rooms.Store(name, NewMeetingRoom(name, capacity))
	return nil
}

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) (*MeetingRoom, error) {
	if val, ok := mrh.rooms.Load(name); !ok {
		return nil, nil
	} else {
		meetingRoom, ok := val.(*MeetingRoom)
		if !ok {
			return nil, nil
		}

		return meetingRoom, nil
	}
}

func (mrh *MeetingRoomsHandler) ReserveMeetingRoom(name string, meetingID string, startTime time.Time, endTime time.Time, participantsCount int) error {
	meetingRoom, err := mrh.GetMeetingRoom(name)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	if !meetingRoom.IsAvailaible(startTime, endTime, participantsCount) {
		return nil
	}

	return meetingRoom.BookSlot(meetingID, startTime, endTime, participantsCount)
}
