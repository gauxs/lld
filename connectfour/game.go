package connectfour

import "github.com/gauxs/lld/connectfour/enum"

type GameRule struct {
	consecutiveCount uint
	directions       []enum.Direction
}
type Game struct {
	currentTurn   uint
	board         *Board
	players       []*Player
	winningPlayer *Player
}

func NewGame(numOfPlayers uint, boardDimensionX uint, boardDimensionY uint) *Game {
	return nil
}

// Start coordinates among the players and act as the game orchestrator
func (g *Game) Start() error {
	// while playable()
	// 	p = nextTurn()
	//	g.board.Place(col, p.GetPiece())
	// if g.boardTopPieceHasNConsecutive(col, consecutiveCount, directions)
	//	g.winner = p
	return nil
}

func (g *Game) Winner() *Player {
	return nil
}

func (g *Game) nextTurn() *Player {
	return nil
}

func (g *Game) playable() bool {
	return false
}
