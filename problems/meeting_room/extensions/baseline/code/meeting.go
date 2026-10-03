package code

import (
	"sync"

	"github.com/gauxs/lld/problems/meeting_room/extensions/baseline/code/enum"
)

type Meeting struct {
	rwMu *sync.RWMutex

	id           string
	title        string
	roomName     string
	slot         *TimeSlot
	participants []*User
	state        enum.MeetingState
}
