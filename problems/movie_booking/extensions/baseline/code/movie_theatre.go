package code

import "time"

type MovieTheatre struct {
	id      string
	screens []Screen
	cityID  string
}

type Screen struct {
	slots   []Slot
	seating [][]Seat
}

type Slot struct {
	startTime time.Time
	endTime   time.Time
}

type Seat struct {
}
