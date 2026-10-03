---
title: Baseline — codebase
problem: meeting_room
extension: baseline
prev:
  text: design
  link: /problems/meeting-room/extensions/baseline/design
---

# Codebase

Source: [GitHub](https://github.com/gauxs/lld/tree/main/problems/meeting_room/extensions/baseline/code)

Path: ` problems/meeting_room/extensions/baseline/code `

## Directory structure

```text
code/
├── enum/
│   └── meeting_state.go
├── error.go
├── meeting.go
├── meeting_room.go
├── meeting_rooms_handler.go
├── meeting_scheduler.go
├── meetings_handler.go
├── notification_service.go
├── user.go
└── util.go
```

## ` enum/meeting_state.go `

```go
package enum

type MeetingState int

const (
	MEETINGSTATE_INVALID MeetingState = iota
	MEETINGSTATE_BOOKED
	MEETINGSTATE_CANCELLED
)
```

## ` error.go `

```go
package code
```

## ` meeting.go `

```go
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
```

## ` meeting_room.go `

```go
package code

import (
	"sync"
	"time"
)

type TimeSlot struct {
	startTime      time.Time
	endTime        time.Time
	meetingID      string
	bookedCapacity int
}

func NewTimeSlot(startTime time.Time, endTime time.Time, meetingID string, participantsCapacity int) *TimeSlot {
	return &TimeSlot{
		startTime:      startTime,
		endTime:        endTime,
		meetingID:      meetingID,
		bookedCapacity: participantsCapacity,
	}
}

type MeetingRoom struct {
	rwMu *sync.RWMutex

	name        string
	capacity    int
	bookedslots []*TimeSlot // NOTE: deadlock is not possible
}

func NewMeetingRoom(name string, cap int) *MeetingRoom {
	return &MeetingRoom{
		rwMu:        &sync.RWMutex{},
		name:        name,
		capacity:    cap,
		bookedslots: make([]*TimeSlot, 0),
	}
}

func (mr *MeetingRoom) IsAvailaible(startTime time.Time, endTime time.Time, participantsCount int) bool {
	if mr.capacity < participantsCount {
		return false
	}

	// check slot availability
	// NOTE: use binary search for optimization
	for _, slot := range mr.bookedslots {
		if (slot.startTime.After(startTime) && slot.endTime.Before(startTime)) || (slot.startTime.After(endTime) && slot.endTime.Before(endTime)) {
			return false
		}
	}

	return true
}

func (mr *MeetingRoom) BookSlot(meetingID string, startTime time.Time, endTime time.Time, participantsCount int) error {
	mr.bookedslots = append(mr.bookedslots, NewTimeSlot(startTime, endTime, meetingID, participantsCount))
	return nil
}
```

## ` meeting_rooms_handler.go `

```go
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

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) (*MeetingRoom, error) {
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
	meetingRoom, err := mrh.GetMeetingRoom(name)
	if err != nil {

	}

	meetingRoom.rwMu.Lock()
	defer meetingRoom.rwMu.Unlock()

	if !meetingRoom.IsAvailaible(startTime, endTime, participantsCount) {
		return nil
	}

	return meetingRoom.BookSlot(meetingID, startTime, endTime, participantsCount)
}
```

## ` meeting_scheduler.go `

```go
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
	err := ms.mh.UpdateTitle(meetingID, newTitle)
	if err != nil {
		return err
	}

	participants, err := ms.mh.GetMeetingParticipants(meetingID)
	if err != nil {
		return nil
	}

	title, err := ms.mh.GetMeetingTitle(meetingID)
	if err != nil {
		return nil
	}

	for _, participant := range participants {
		ms.notification.NotifyUser(participant, title)
	}

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
```

## ` meetings_handler.go `

```go
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
	if ok {
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
```

## ` notification_service.go `

```go
package code

import "fmt"

type NotificationService struct{}

func (ns *NotificationService) NotifyUser(user *User, meetingTitle string) {
	fmt.Printf("%v, meeting %v is updated", user.name, meetingTitle)
}
```

## ` user.go `

```go
package code

type User struct {
	name string
}

func NewUser(name string) *User {
	return &User{
		name: name,
	}
}
```

## ` util.go `

```go
package code

import "github.com/google/uuid"

func GenerateID() string {
	return uuid.New().String()
}
```
