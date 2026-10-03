package code

import (
	"time"

	"github.com/gauxs/lld/problems/meeting_room/extensions/baseline/code/enum"
)

type MeetingScheduler struct {
	mrhandler    *MeetingRoomsHandler
	mh           *MeetingsHandler
	notification *NotificationService
}

func (ms *MeetingScheduler) ScheduleMeeting(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []*User) (string, error) {
	meeting := NewMeeting(meetingTitle, roomName, startTime, endTime, participants)
	// 1 - Try to reserve the room's slots
	if err := ms.mrhandler.ReserveMeetingRoom(roomName, meeting.id, startTime, endTime, len(participants)); err != nil {
		return "", nil
	}

	// 2 - If room reserved, update the meeting
	meeting.state = enum.MEETINGSTATE_BOOKED

	return meeting.id, ms.mh.AddMeeting(meeting)
}

func (ms *MeetingScheduler) UpdateMeetingTitle(meetingID string, newTitle string) error {
	return nil
}

func (ms *MeetingScheduler) UpdateMeetingParticipants(meetingID string, newparticipants []string) error {
	return nil
}

func (ms *MeetingScheduler) UpdateMeetingSchedule(meetingID string, newStartTime time.Time, endTime time.Time) error {
	// 1 - Try to reserve the room's slots
	// 2 - If room reserved, update the meeting
	return nil
}

func (ms *MeetingScheduler) CancelMeeting(meetingID string) error {
	// 1 - Free up the room's slots
	// 2 - Update the meeting with CANCELLED state
	return nil
}
