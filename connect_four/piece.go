package connectfour

import "github.com/gauxs/lld/connect_four/enum"

type Piece struct {
	color enum.Color
}

func NewPiece(c enum.Color) *Piece {
	return &Piece{
		color: c,
	}
}

func (p *Piece) Color() enum.Color {
	return p.color
}
