package code

import (
	"sync"
	"time"

	"github.com/gauxs/lld/problems/meeting_room/extensions/baseline/code/enum"
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

func (ms *MeetingsHandler) getMeeting(meetingID string) (*Meeting, error) {
	val, ok := ms.meetings.Load(meetingID)
	if !ok {
		return nil, nil
	}

	meeting, ok := val.(*Meeting)
	if !ok {
		return nil, nil
	}

	return meeting, nil
}

func (ms *MeetingsHandler) GetMeetingParticipants(meetingID string) ([]*User, error) {
	m, err := ms.getMeeting(meetingID)
	if err != nil {
		return nil, err
	}

	m.rwMu.Lock()
	defer m.rwMu.Unlock()

	// Create a copy so the caller gets a safe snapshot
	copied := make([]*User, len(m.participants))
	copy(copied, m.participants)

	return copied, nil
}

func (ms *MeetingsHandler) GetMeetingTitle(meetingID string) (string, error) {
	m, err := ms.getMeeting(meetingID)
	if err != nil {
		return "", err
	}

	m.rwMu.Lock()
	defer m.rwMu.Unlock()

	return m.title, nil
}

func (ms *MeetingsHandler) GetMeetingRoomName(meetingID string) (string, error) {
	m, err := ms.getMeeting(meetingID)
	if err != nil {
		return "", err
	}

	m.rwMu.Lock()
	defer m.rwMu.Unlock()

	return m.roomName, nil
}

func (ms *MeetingsHandler) UpdateTitle(meetingID string, newTitle string) error {
	meeting, err := ms.getMeeting(meetingID)
	if err != nil {
		return err
	}

	meeting.rwMu.Lock()
	defer meeting.rwMu.Unlock()

	meeting.title = newTitle
	return nil
}

func (ms *MeetingsHandler) UpdateParticipants(meetingID string, newparticipants []*User) error {
	meeting, err := ms.getMeeting(meetingID)
	if err != nil {
		return err
	}

	meeting.rwMu.Lock()
	defer meeting.rwMu.Unlock()

	meeting.participants = newparticipants
	return nil
}

func (ms *MeetingsHandler) UpdateSchedule(meetingID string, newStartTime time.Time, newEndTime time.Time) error {
	meeting, err := ms.getMeeting(meetingID)
	if err != nil {
		return err
	}

	meeting.rwMu.Lock()
	defer meeting.rwMu.Unlock()

	meeting.startTime = newStartTime
	meeting.endTime = newEndTime

	return nil
}

func (ms *MeetingsHandler) Cancel(meetingID string) error {
	meeting, err := ms.getMeeting(meetingID)
	if err != nil {
		return err
	}

	meeting.rwMu.Lock()
	defer meeting.rwMu.Unlock()

	meeting.state = enum.MEETINGSTATE_CANCELLED
	return nil
}
