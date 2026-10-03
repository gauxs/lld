package code

import (
	"sync"
	"time"
)

type MeetingsHandler struct {
	meetings sync.Map
}

func (ms *MeetingsHandler) AddMeeting(meeting *Meeting) error {
	if _, ok := ms.meetings.Load(meeting.id); ok {
		return nil
	}

	ms.meetings.Store(meeting.id, meeting)
	return nil
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
