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
└── user.go
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
```

## ` meeting_room.go `

```go
package code

import (
	"sync"
	"time"
)

type TimeSlot struct {
	startTime time.Time
	endTime   time.Time
	meetingID string
}

type MeetingRoom struct {
	rwMu *sync.RWMutex

	name        string
	capacity    int
	bookedslots []*TimeSlot // NOTE: deadlock is not possible
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

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string) error {
	return nil
}

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) error {
	return nil
}

func (mrh *MeetingRoomsHandler) ReserveMeetingRoom(name string, meetingID string, startTime time.Time, endTime time.Time) error {
	return nil
}
```

## ` meeting_scheduler.go `

```go
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
```

## ` meetings_handler.go `

```go
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
```

## ` notification_service.go `

```go
package code

type NotificationService struct{}

func (ns *NotificationService) NotifyUser(userName string, msg string) {}
```

## ` user.go `

```go
package code

type User struct {
	name string
}
```
