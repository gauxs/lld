package connectfour

import (
	"github.com/gauxs/lld/connectfour/enum"
)

type Board struct {
	gridX uint
	gridY uint
	grid  [][]*Piece
}

func NewBoard(rows uint, col uint) *Board {
	return nil
}

func (b *Board) Place(col uint, p *Piece) bool {
	return false
}

func (b *Board) isValidMove(row uint, col uint) bool {
	return true
}

func (b *Board) TopPieceHasNConsecutive(col uint, n uint, directions []enum.Direction) bool {
	return true
}

func (b *Board) IsBoardFull() bool {
	return false
}
