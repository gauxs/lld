package connectfour

import (
	"github.com/gauxs/lld/connect_four/enum"
)

type Game struct {
	players       []*Player
	board         *Board
	state         enum.GameState
	currentPlayer *Player
	winner        *Player
	winRule       Rule
}

func (g *Game) GetGameState() enum.GameState {
	return g.state
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *Game) GetWinner() *Player {
	return g.winner
}

func (g *Game) MakeMove(p *Player, col int) error {
	if err := g.board.PlaceDisk(p.GetDisk(), col); err != nil {
		return err
	}

	return nil
}

type Rule struct {
}
