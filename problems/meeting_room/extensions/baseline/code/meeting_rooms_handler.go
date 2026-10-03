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

func (mrh *MeetingRoomsHandler) getMeetingRoom(name string) (*MeetingRoom, error) {
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
	meetingRoom, err := mrh.getMeetingRoom(name)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	if !meetingRoom.IsAvailaible(startTime, endTime, participantsCount) {
		return nil
	}

	return meetingRoom.BookSlot(meetingID, startTime, endTime, participantsCount)
}

func (mrh *MeetingRoomsHandler) UpdateMeetingRoomParticipantCount(roomName string, meetingID string, newParticipantsCount int) error {
	meetingRoom, err := mrh.getMeetingRoom(roomName)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	if meetingRoom.capacity < newParticipantsCount {
		return nil
	}

	var currentSlot *TimeSlot
	for _, slot := range meetingRoom.bookedslots {
		if slot.meetingID == meetingID {
			currentSlot = slot
		}
	}

	currentSlot.bookedCapacity = newParticipantsCount
	return nil
}

func (mrh *MeetingRoomsHandler) UpdateMeetingSchedule(roomName string, meetingID string, newStartTime time.Time, newEndTime time.Time) error {
	meetingRoom, err := mrh.getMeetingRoom(roomName)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	var targetIndex = -1
	var currentSlot *TimeSlot
	for idx, slot := range meetingRoom.bookedslots {
		if slot.meetingID == meetingID {
			currentSlot = slot
			targetIndex = idx
		}
	}

	if targetIndex != -1 {
		meetingRoom.bookedslots = append(meetingRoom.bookedslots[:targetIndex], meetingRoom.bookedslots[targetIndex+1:]...)
	}

	if !meetingRoom.IsAvailaible(newStartTime, newEndTime, currentSlot.bookedCapacity) {
		meetingRoom.BookSlot(meetingID, currentSlot.startTime, currentSlot.endTime, currentSlot.bookedCapacity)
		return nil
	}

	meetingRoom.BookSlot(meetingID, newEndTime, newEndTime, currentSlot.bookedCapacity)

	return nil
}

func (mrh *MeetingRoomsHandler) RemoveMeeting(roomName string, meetingID string) error {
	meetingRoom, err := mrh.getMeetingRoom(roomName)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	var targetIndex = -1
	for idx, slot := range meetingRoom.bookedslots {
		if slot.meetingID == meetingID {
			targetIndex = idx
		}
	}

	if targetIndex != -1 {
		meetingRoom.bookedslots = append(meetingRoom.bookedslots[:targetIndex], meetingRoom.bookedslots[targetIndex+1:]...)
	}

	return nil
}
