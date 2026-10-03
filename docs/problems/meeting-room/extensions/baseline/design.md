---
title: Meeting room — baseline (design)
problem: meeting_room
extension: baseline
prev:
  text: requirements
  link: /problems/meeting-room/extensions/baseline/requirements
---

## Approach

## Class design & Relationships

```go

type TimeSlot struct{
    startTime time.Time
    endTime time.Time
    isBusy bool
}

type MeetingRoom struct {
    rwMu *sync.RWMutex
    name string
    // 15 min slots refreshed for next 30 days
    slots []*TimeSlot
}

type MeetingRoomsHandler struct{
    rooms sync.Map[string]*MeetingRoom
}

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string) error {}

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) error {}

func (mrh *MeetingRoomsHandler) ReserveMeetingRoom(name string, startTime time.Time, endTime time.Time) error {}

type Meeting struct{
    title string
    roomName string
    slots []*TimeSlot
    participants []string
    state enum.MeetingState
}

type MeetingsHandler struct{
    mrhandler *MeetingRoomsHandler
    meetings sync.Map[string]*Meeting
}

func (ms *MeetingScheduler) ScheduleMeeting(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []string) (*Meeting, error) {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, create a meeting 
}

func (ms *MeetingScheduler) UpdateMeetingTitle(currentTitle string, newTitle string) error {}

func (ms *MeetingScheduler) UpdateMeetingParticipants(meetingTitle string, newparticipants []string) error {}

func (ms *MeetingScheduler) UpdateMeetingSchedule(meetingTitle string, newStartTime time.Time, endTime time.Time) error {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, update the meeting 
}

func (ms *MeetingScheduler) CancelMeeting(meetingTitle string) error {
    // 1 - Free up the room's slots
    // 2 - Update the meeting with CANCELLED state
}

func ()
```

## Main flow

1. Validate the request.
2. Read or update the relevant state.
3. Return the result.

## Invariants and concurrency

Explain shared state, invariants, and synchronization.

## Complexity

- Write: `O(?)`
- Read: `O(?)`
- Space: `O(?)`

## Trade-offs

- Decision and benefit
- Limitation and possible future improvement
