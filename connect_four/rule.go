package connectfour

import "github.com/gauxs/lld/connect_four/enum"

type Rule interface {
	HasWon(board *Board, c int) bool
}

type StandardRules struct {
	consecutiveCount uint
	directions       []enum.Direction
}

func NewStandardRules(consecutiveCount int, directions []enum.Direction) *StandardRules {
	return &StandardRules{
		consecutiveCount: uint(consecutiveCount),
		directions:       directions,
	}
}

func (sr *StandardRules) HasWon(board *Board, c int) bool {
	return true
}
