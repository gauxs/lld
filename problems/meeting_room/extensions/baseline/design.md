# Meeting room — baseline (design)

## Approach

## Class design & Relationships

```go

type TimeSlot struct{
    startTime time.Time
    endTime time.Time
    meetingID string
}

type MeetingRoom struct {
    rwMu *sync.RWMutex

    name string
    capacity int
    bookedslots []*TimeSlot // NOTE: deadlock is not possible
}

type MeetingRoomsHandler struct{
    rooms sync.Map[string]*MeetingRoom
}

func (mrh *MeetingRoomsHandler) AddMeetingRoom(name string) error {}

func (mrh *MeetingRoomsHandler) GetMeetingRoom(name string) error {}

func (mrh *MeetingRoomsHandler) ReserveMeetingRoom(name string, startTime time.Time, endTime time.Time) error {}

type Meeting struct{
    rwMu *sync.RWMutex

    id string
    title string
    roomName string
    slot *TimeSlot
    participants []MeetingParticipant
    state enum.MeetingState
}

type MeetingsHandler struct{
    meetings sync.Map[string]*Meeting
}

func (ms *MeetingsHandler) Schedule(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []string) (string, error) {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, create a meeting 
}

func (ms *MeetingsHandler) UpdateTitle(meetingID string, newTitle string) error {}

func (ms *MeetingsHandler) UpdateParticipants(meetingID string, newparticipants []string) error {}

func (ms *MeetingsHandler) UpdateSchedule(meetingID string, newStartTime time.Time, endTime time.Time) error {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, update the meeting 
}

func (ms *MeetingsHandler) Cancel(meetingID string) error {
    // 1 - Free up the room's slots
    // 2 - Update the meeting with CANCELLED state
}

type MeetingScheduler struct{
    mrhandler *MeetingRoomsHandler
    mh *MeetingsHandler
    notification *NotificationService
}

func (ms *MeetingScheduler) ScheduleMeeting(meetingTitle string, roomName string, startTime time.Time, endTime time.Time, participants []string) (string, error) {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, create a meeting 
}

func (ms *MeetingScheduler) UpdateMeetingTitle(meetingID string, newTitle string) error {}

func (ms *MeetingScheduler) UpdateMeetingParticipants(meetingID string, newparticipants []string) error {}

func (ms *MeetingScheduler) UpdateMeetingSchedule(meetingID string, newStartTime time.Time, endTime time.Time) error {
    // 1 - Try to reserve the room's slots
    // 2 - If room reserved, update the meeting 
}

func (ms *MeetingScheduler) CancelMeeting(meetingID string) error {
    // 1 - Free up the room's slots
    // 2 - Update the meeting with CANCELLED state
}

type User struct {
    name string
}

type NotificationService struct{}

func (ns *NotificationService) NotifyUser(userName string, msg string){}

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
