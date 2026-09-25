package connectfour

import "github.com/gauxs/lld/problems/connect_four/code/enum"

type Rule interface {
	Satisfied(b *Board, row int, col int) bool
}

type FourRule struct {
	directions       []enum.Direction
	consecutiveCount int
}

func NewFourRule() Rule {
	return &FourRule{
		directions: []enum.Direction{enum.DIRECTION_HORIZONTAL, enum.DIRECTION_VERTICAL,
			enum.DIRECTION_DIAGONAL_BACK, enum.DIRECTION_DIAGONAL_FRONT},
		consecutiveCount: 4,
	}
}

func (r *FourRule) Satisfied(b *Board, row int, col int) bool {
	for _, direction := range r.directions {
		if b.CountInDirection(row, col, direction, r.consecutiveCount) >= r.consecutiveCount {
			return true
		}
	}

	return false
}
