package enum

type MeetingState int

const (
	MEETINGSTATE_INVALID MeetingState = iota
	MEETINGSTATE_BOOKED
	MEETINGSTATE_CANCELLED
)
