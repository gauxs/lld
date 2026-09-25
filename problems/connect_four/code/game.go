package connectfour

import (
	"github.com/gauxs/lld/problems/connect_four/code/enum"
)

type Game struct {
	players       []*Player
	board         *Board
	state         enum.GameState
	currentPlayer *Player
	winner        *Player
	winRule       Rule
}

func NewGame(boardX int, boardY int) *Game {
	return &Game{
		players:       make([]*Player, 0),
		board:         NewBoard(boardX, boardY),
		state:         enum.GAMESTATE_NOT_PLAYING,
		currentPlayer: nil,
		winner:        nil,
		winRule:       NewFourRule(),
	}
}

func (g *Game) AddPlayer(name string, d enum.Disc) error {
	if g.state != enum.GAMESTATE_NOT_PLAYING {
		return ErrGameNotInCorrectState
	}
	for _, p := range g.players {
		if p.GetDisk() == d {
			return ErrDiskAlreadyTaken
		}
	}

	g.players = append(g.players, NewPlayer(name, d))
	g.currentPlayer = g.players[0]
	return nil
}

func (g *Game) StartGame() error {
	if g.state != enum.GAMESTATE_NOT_PLAYING {
		return ErrGameNotInCorrectState
	}

	if len(g.players) < 2 {
		return ErrInsufficientPlayers
	}

	g.state = enum.GAMESTATE_PLAYING
	return nil
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

func (g *Game) nextPlayer() *Player {
	found := false
	for _, p := range g.players {
		if found {
			return p
		}

		if p.GetName() == g.currentPlayer.name {
			found = true
		}
	}

	return g.players[0]
}

func (g *Game) MakeMove(col int) error {
	if g.state != enum.GAMESTATE_PLAYING {
		return ErrGameNotInCorrectState
	}

	row, err := g.board.PlaceDisk(g.currentPlayer.GetDisk(), col)
	if err != nil {
		return err
	}

	if g.winRule.Satisfied(g.board, row, col) {
		g.winner = g.currentPlayer
		g.state = enum.GAMESTATE_WON
		return nil
	}

	if g.board.IsFull() {
		g.state = enum.GAMESTATE_DRAW
		return nil
	}

	g.currentPlayer = g.nextPlayer()
	return nil
}
