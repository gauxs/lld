package code

import "time"

type SeatType int

const (
	SEATTYPE_INVALID SeatType = iota
	SEATTYPE_RECLINER
)

type MovieTheatre struct {
	id      string
	screens []Screen
	cityID  string
}

type Screen struct {
	slots []Slot
}

type Slot struct {
	movieID   string
	startTime time.Time
	endTime   time.Time
	seating   *ScreenSeatingArrangement
}

type ScreenSeatingArrangement struct {
	seating [][]Seat
}

type Seat struct {
	seatType SeatType
}
