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
	grid := make([][]*Piece, rows)

	for i := uint(0); i < rows; i++ {
		grid[i] = make([]*Piece, col)
	}

	return &Board{
		gridX: rows,
		gridY: col,
		grid:  grid,
	}
}

func (b *Board) Place(col uint, p *Piece) bool {
	if !b.isValidMove(col) {
		return false
	}

	er := b.emptyRow(col)
	b.grid[er][col] = p
	return true
}

func (b *Board) emptyRow(col uint) int {
	if col >= b.gridY {
		return -1
	}

	curRow := int(b.gridX - 1)
	for curRow >= 0 && b.grid[curRow][col] != nil {
		curRow--
	}

	return curRow
}

func (b *Board) isValidMove(col uint) bool {
	row := b.emptyRow(col)
	return row >= 0 && row < int(b.gridX) && col < b.gridY
}

func (b *Board) TopPieceHasNConsecutive(col uint, n uint, directions []enum.Direction) bool {
	row := b.emptyRow(col) + 1
	for _, d := range directions {
		won := true
		switch d {
		case enum.DIRECTION_VERTICAL:
			maxRow := min(int(b.gridX-1), row+int(n-1))
			r := row + 1
			for ; r <= maxRow; r++ {
				if b.grid[r][col] == nil {
					won = false
					break
				}

				if b.grid[r][col].Color() != b.grid[r-1][col].Color() {
					won = false
					break
				}
			}

			if r > maxRow && won {
				return won
			}
		case enum.DIRECTION_HORIZONTAL:
			// TODO
		case enum.DIRECTION_DIAGONAL:
			// TODO
		}
	}

	return false
}

func (b *Board) IsBoardFull() bool {
	for i := 0; i < int(b.gridY); i++ {
		if b.grid[0][i] == nil {
			return false
		}
	}

	return true
}
