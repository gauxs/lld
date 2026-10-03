package code

import (
	"sync"
	"time"
)

type MeetingsHandler struct {
	meetings sync.Map
}

func (ms *MeetingsHandler) Schedule(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []string) (string, error) {
	// 1 - Try to reserve the room's slots
	// 2 - If room reserved, create a meeting
	return "", nil
}

func (ms *MeetingsHandler) UpdateTitle(meetingID string, newTitle string) error {
	return nil
}

func (ms *MeetingsHandler) UpdateParticipants(meetingID string, newparticipants []string) error {
	return nil
}

func (ms *MeetingsHandler) UpdateSchedule(meetingID string, newStartTime time.Time, endTime time.Time) error {
	return nil
}

func (ms *MeetingsHandler) Cancel(meetingID string) error {
	return nil
}
