package connectfour

import (
	"fmt"

	"github.com/gauxs/lld/connect_four/enum"
)

type Board struct {
	row    int
	column int
	grid   [][]enum.Disc
}

func (b *Board) getNextFreeSlot(col int) int {
	for r := b.row - 1; r >= 0; r-- {
		if b.grid[r][col] == enum.DISK_INVALID {
			return r
		}
	}

	return -1
}

func (b *Board) PlaceDisk(d enum.Disc, col int) error {
	if col < 0 || col >= b.column {
		return fmt.Errorf("column %d: %w", col, ErrInvalidMove)
	}

	r := b.getNextFreeSlot(col)
	if r == -1 {
		return fmt.Errorf("column %d if full: %w", col, ErrInvalidMove)
	}

	b.grid[r][col] = d
	return nil
}

func (b *Board) countInDirection(r int, c int, limit int, refDisk enum.Disc, dr int, dc int) int {
	if c < 0 || c >= b.column || r < 0 || r >= b.row {
		return 0
	}

	if b.grid[r][c] != refDisk || b.grid[r][c] == enum.DISK_INVALID {
		return 0
	}

	if limit <= 0 {
		return 0 // early cut
	}

	return 1 + b.countInDirection(r+dr, c+dc, limit-1, refDisk, dr, dc)
}

var directionRowColSteps = map[enum.Direction][][]int{
	enum.DIRECTION_HORIZONTAL:     {{0, 1}, {0, -1}},
	enum.DIRECTION_VERTICAL:       {{1, 0}, {-1, 0}},
	enum.DIRECTION_DIAGONAL_BACK:  {{-1, -1}, {1, 1}},
	enum.DIRECTION_DIAGONAL_FRONT: {{-1, 1}, {1, -1}},
}

func (b *Board) CountInDirection(r int, c int, d enum.Direction, limit int) int {
	if c < 0 || c >= b.column || r < 0 || r >= b.row {
		return 0
	}

	count := 0
	for _, rowColStep := range directionRowColSteps[d] {
		count += b.countInDirection(r+rowColStep[0], c+rowColStep[1], limit, b.grid[r][c], rowColStep[0], rowColStep[1])
	}

	return 1 + count
}

func (b *Board) IsFull() bool {
	for c := 0; c < b.column; c++ {
		if b.getNextFreeSlot(c) != -1 {
			return false
		}
	}

	return true
}
