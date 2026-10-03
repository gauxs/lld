package code

import (
	"sync"
	"time"

	"github.com/gauxs/lld/problems/meeting_room/extensions/baseline/code/enum"
)

type Meeting struct {
	rwMu *sync.RWMutex

	id           string
	title        string
	roomName     string
	startTime    time.Time
	endTime      time.Time
	participants []*User
	state        enum.MeetingState
}

func NewMeeting(title string, roomName string, startTime time.Time, endTime time.Time, p []*User) *Meeting {
	return &Meeting{
		rwMu: &sync.RWMutex{},

		id:           GenerateID(),
		title:        title,
		roomName:     roomName,
		startTime:    startTime,
		endTime:      endTime,
		participants: p,
		state:        enum.MEETINGSTATE_INVALID,
	}
}
