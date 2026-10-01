package enum

type Activity int

const (
	UNDEFINED_ACTIVITY Activity = iota
	LOGIN
	LOGOUT
)
