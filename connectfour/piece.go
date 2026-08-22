package connectfour

type Color int

const (
	INVALID Color = iota
	BLUE
	RED
)

type Piece struct {
	color Color
}

func NewPiece(c Color) *Piece {
	return nil
}
