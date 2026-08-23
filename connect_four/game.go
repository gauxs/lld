package connectfour

import (
	"fmt"

	"github.com/gauxs/lld/connectfour/enum"
)

type Game struct {
	turnNo        int
	board         *Board
	players       []*Player
	winningPlayer *Player
	rule          Rule
}

func NewGame(numOfPlayers uint, boardDimensionX uint, boardDimensionY uint) *Game {
	players := make([]*Player, numOfPlayers)

	for i := uint(0); i < numOfPlayers; i++ {
		players[i] = NewPlayer(NewPiece(nextAvailaibleColor(players)))
	}

	return &Game{
		turnNo:        -1,
		board:         NewBoard(boardDimensionX, boardDimensionY),
		players:       players,
		winningPlayer: nil,
		rule: NewStandardRules(4, []enum.Direction{
			enum.DIRECTION_DIAGONAL,
			enum.DIRECTION_HORIZONTAL,
			enum.DIRECTION_VERTICAL,
		}),
	}
}

// Start coordinates among the players and act as the game orchestrator
func (g *Game) Start() error {
	for g.playable() {
		col := uint(0)
		curentPlayer := g.nextTurn()

		fmt.Println("Enter valid column for player %s: ", curentPlayer.GetPiece().Color())
		_, _ = fmt.Scan(&col)

		for !g.board.isValidMove(col) {
			fmt.Println("Enter valid column for player %s: ", curentPlayer.GetPiece().Color())
			_, _ = fmt.Scan(&col)
		}

		g.board.Place(col, curentPlayer.GetPiece())

		// if g.board.TopPieceHasNConsecutive(col, g.rules.consecutiveCount, g.rules.directions) {
		if g.rule.HasWon(g.board, int(col)) {
			g.winningPlayer = curentPlayer
			break
		}
	}

	if g.Winner() == nil {
		fmt.Println("The game ended with draw")
	} else {
		fmt.Println("The game ended with winning color: ", g.Winner().GetPiece().Color())
	}

	return nil
}

func (g *Game) Winner() *Player {
	return g.winningPlayer
}

func (g *Game) nextTurn() *Player {
	g.turnNo = (g.turnNo + 1) % len(g.players)
	return g.players[g.turnNo]
}

func (g *Game) playable() bool {
	return !g.board.IsBoardFull() && g.winningPlayer == nil
}
