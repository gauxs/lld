package enum

type Direction int

const (
	DIRECTION_INVALID Direction = iota
	DIRECTION_HORIZONTAL
	DIRECTION_VERTICAL
	DIRECTION_DIAGONAL
)
