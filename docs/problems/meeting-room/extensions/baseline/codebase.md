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
```

## ` meeting_room.go `

```go
package code
```

## ` meeting_rooms_handler.go `

```go
package code
```

## ` meeting_scheduler.go `

```go
package code
```

## ` meetings_handler.go `

```go
package code
```

## ` notification_service.go `

```go
package code
```

## ` user.go `

```go
package code
```
