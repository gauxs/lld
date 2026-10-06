package code

type BookingStatus int

const (
	BOOKINGSTATUS_INVALID BookingStatus = iota
	BOOKINGSTATUS_PAYMENT_PENDING
	BOOKINGSTATUS_CONFIRMED
)

type Booking struct {
	id             string
	movieTheatreID string
	screenID       string
	slotID         string
	seatIDs        []string
	status         BookingStatus
}
