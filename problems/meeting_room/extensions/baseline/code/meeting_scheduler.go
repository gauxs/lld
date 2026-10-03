package code

import "time"

type MeetingScheduler struct {
	mrhandler    *MeetingRoomsHandler
	mh           *MeetingsHandler
	notification *NotificationService
}

func (ms *MeetingScheduler) ScheduleMeeting(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []string) (string, error) {
	// 1 - Try to reserve the room's slots
	// 2 - If room reserved, create a meeting
	return "", nil
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
