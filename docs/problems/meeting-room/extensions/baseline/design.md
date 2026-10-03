---
title: Meeting room — baseline (design)
problem: meeting_room
extension: baseline
prev:
  text: requirements
  link: /problems/meeting-room/extensions/baseline/requirements
---

## Approach

Summarize the design and why it fits the requirements.

## Class design & Relationships

```go

type TimeSlot struct{
    startTime time.Time
    endTime time.Time
    isBusy bool
}

type MeetingRoom struct {
    name string
    // 15 min slots refreshed for next 30 days
    slots []*TimeSlot
}

type MeetingRoomsHandler struct{
    rooms []*MeetingRoom
}

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string) error {}

type Meeting struct{
    title string
    roomName string
    slots []*TimeSlot
    participants []string
}

type MeetingScheduler struct{
    mrhandler *MeetingRoomsHandler
}

func (ms *MeetingScheduler) ScheduleMeeting(startTime time.Time, endTime time.Time, participants []string) *Meeting {}

func (ms *MeetingScheduler) UpdateMeetingTitle(meeting *Meeting, newTitle string) error {}

func (ms *MeetingScheduler) UpdateMeetingParticipants(meeting *Meeting, newparticipants []string) error {}

func (ms *MeetingScheduler) UpdateMeetingSchedule(meeting *Meeting, newStartTime time.Time, endTime time.Time) error {}

func (ms *MeetingScheduler) CancelMeeting(meeting *Meeting) error {}

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
