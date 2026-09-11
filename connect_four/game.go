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
	winRule       *Rule
}

func NewGame(boardX int, boardY int) *Game {
	return &Game{
		players:       make([]*Player, 0),
		board:         NewBoard(boardX, boardY),
		state:         enum.GAMESTATE_NOT_PLAYING,
		currentPlayer: nil,
		winner:        nil,
		winRule:       NewRule(),
	}
}

func (g *Game) AddPlayer(name string, d enum.Disc) error {
	for _, p := range g.players {
		if p.GetDisk() == d {
			return ErrDiskAlreadyTaken
		}
	}

	g.players = append(g.players, NewPlayer(name, d))
	g.currentPlayer = g.players[0]
	return nil
}

func (g *Game) StartGame() {
	g.state = enum.GAMESTATE_PLAYING
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
	if err := g.board.PlaceDisk(g.currentPlayer.GetDisk(), col); err != nil {
		return err
	}

	if g.winRule.Satisfied(g.board, g.board.getNextFreeSlot(col)+1, col) {
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

type Rule struct {
	directions       []enum.Direction
	consecutiveCount int
}

func NewRule() *Rule {
	return &Rule{
		directions: []enum.Direction{enum.DIRECTION_HORIZONTAL, enum.DIRECTION_VERTICAL,
			enum.DIRECTION_DIAGONAL_BACK, enum.DIRECTION_DIAGONAL_FRONT},
		consecutiveCount: 4,
	}
}

func (r *Rule) Satisfied(b *Board, row int, col int) bool {
	for _, direction := range r.directions {
		if b.CountInDirection(row, col, direction, r.consecutiveCount) >= r.consecutiveCount {
			return true
		}
	}

	return false
}
