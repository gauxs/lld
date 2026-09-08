package enum

type BoardStatus int

const (
	INVALID BoardStatus = iota
	DRAW
	WON
	PLAYABLE
)
