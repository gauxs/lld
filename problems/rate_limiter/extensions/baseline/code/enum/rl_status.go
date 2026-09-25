package enum

type RLStatus int

const (
	RLStatus_INVALID RLStatus = iota
	RLStatus_ACCEPTED
	RLStatus_REJECTED
)
