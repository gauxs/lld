package code

import (
	"time"

	"github.com/gauxs/lld/problems/meeting_room/extensions/baseline/code/enum"
)

type MeetingScheduler struct {
	mrhandler    *MeetingRoomsHandler
	mhandler     *MeetingsHandler
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

	return meeting.id, ms.mhandler.AddMeeting(meeting)
}

func (ms *MeetingScheduler) UpdateMeetingTitle(meetingID string, newTitle string) error {
	err := ms.mhandler.UpdateTitle(meetingID, newTitle)
	if err != nil {
		return err
	}

	participants, err := ms.mhandler.GetMeetingParticipants(meetingID)
	if err != nil {
		return nil
	}

	title, err := ms.mhandler.GetMeetingTitle(meetingID)
	if err != nil {
		return nil
	}

	for _, participant := range participants {
		ms.notification.NotifyUser(participant, title)
	}

	return nil
}

func (ms *MeetingScheduler) UpdateMeetingParticipants(meetingID string, newparticipants []*User) error {
	meetingRoomName, err := ms.mhandler.GetMeetingRoomName(meetingID)
	if err != nil {
		return err
	}

	// update meeting room participant count
	err = ms.mrhandler.UpdateMeetingRoomParticipantCount(meetingRoomName, meetingID, len(newparticipants))
	if err != nil {
		return err
	}

	// update meeting participants
	err = ms.mhandler.UpdateParticipants(meetingID, newparticipants)
	if err != nil {
		return err
	}

	// notify particiapnts
	participants, err := ms.mhandler.GetMeetingParticipants(meetingID)
	if err != nil {
		return nil
	}

	title, err := ms.mhandler.GetMeetingTitle(meetingID)
	if err != nil {
		return nil
	}

	for _, participant := range participants {
		ms.notification.NotifyUser(participant, title)
	}

	return nil
}

func (ms *MeetingScheduler) UpdateMeetingSchedule(meetingID string, newStartTime time.Time, newEndTime time.Time) error {
	meetingRoomName, err := ms.mhandler.GetMeetingRoomName(meetingID)
	if err != nil {
		return err
	}

	// update meeting room participant count
	err = ms.mrhandler.UpdateMeetingSchedule(meetingRoomName, meetingID, newStartTime, newEndTime)
	if err != nil {
		return err
	}

	// notify particiapnts
	participants, err := ms.mhandler.GetMeetingParticipants(meetingID)
	if err != nil {
		return nil
	}

	title, err := ms.mhandler.GetMeetingTitle(meetingID)
	if err != nil {
		return nil
	}

	for _, participant := range participants {
		ms.notification.NotifyUser(participant, title)
	}

	return nil
}

func (ms *MeetingScheduler) CancelMeeting(meetingID string) error {
	meetingRoomName, err := ms.mhandler.GetMeetingRoomName(meetingID)
	if err != nil {
		return err
	}

	// update meeting room participant count
	err = ms.mrhandler.RemoveMeeting(meetingRoomName, meetingID)
	if err != nil {
		return err
	}

	err = ms.mhandler.Cancel(meetingID)
	if err != nil {
		return err
	}

	// notify particiapnts
	participants, err := ms.mhandler.GetMeetingParticipants(meetingID)
	if err != nil {
		return nil
	}

	title, err := ms.mhandler.GetMeetingTitle(meetingID)
	if err != nil {
		return nil
	}

	for _, participant := range participants {
		ms.notification.NotifyUser(participant, title)
	}

	return nil
}
