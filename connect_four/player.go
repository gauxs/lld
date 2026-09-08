package connectfour

type Player struct {
	piece *Piece
}

func NewPlayer(p *Piece) *Player {
	return &Player{
		piece: p,
	}
}

func (p *Player) GetPiece() *Piece {
	return p.piece
}
