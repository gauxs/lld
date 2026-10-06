package code

import "time"

type SeatType int

const (
	SEATTYPE_INVALID SeatType = iota
	SEATTYPE_REGULAR
	SEATTYPE_PREMIUM
	SEATTYPE_RECLINER
)

type MovieTheatre struct {
	id      string
	screens []Screen
	cityID  string
}

type Screen struct {
	id    string
	slots []Slot
}

type Slot struct {
	id        string
	movieID   string
	startTime time.Time
	endTime   time.Time
	seating   *ScreenSeatingArrangement
}

type ScreenSeatingArrangement struct {
	seating []*Seat
}

type Seat struct {
	id       string
	seatType SeatType
	row      int
	seatNo   int
}
