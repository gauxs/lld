---
title: Connect Four — local two-player (design)
problem: connect_four
extension: baseline
prev:
  text: requirements
  link: /problems/connect-four/extensions/baseline/requirements
next:
  text: codebase
  link: /problems/connect-four/extensions/baseline/codebase
---

## Approach

`Game` orchestrates player registration, the game lifecycle, turn order, and `MakeMove`. `Board` owns grid storage, column validity, gravity placement, and directional line counting from the last move. A pluggable `Rule` decides whether the last placement satisfies the win condition.

## Go design sketch

Method bodies are intentionally omitted. Comments describe responsibilities and invariants.

```go
// Game owns player registration, turn order, and the game lifecycle:
// NOT_PLAYING → PLAYING → WON | DRAW. It orchestrates MakeMove.
type Game struct {
	board         *Board
	players       [2]*Player
	currentPlayer int
	state         enum.GameState
	winner        *Player
	rule          Rule
}

// Board owns grid storage, column validity, gravity placement, and
// directional line counting from the last move.
type Board struct {
	grid [][]enum.Disc
}

// Player has a display name and disc color.
type Player struct {
	name string
	disc enum.Disc
}

// Rule determines whether the last placement satisfies the win condition.
type Rule interface {
	Satisfied(b *Board, row, col int) bool
}

// NewGame creates an empty board with the default FourRule.
func NewGame(rows, cols int) *Game

// AddPlayer registers a player before start; disc colors must be unique.
func (g *Game) AddPlayer(name string, d enum.Disc) error

// StartGame requires two players.
func (g *Game) StartGame() error

// MakeMove drops the current player's disc in col, then updates state and turn.
func (g *Game) MakeMove(col int) error

func (g *Game) GetGameState() enum.GameState

func (g *Game) GetCurrentPlayer() *Player

// GetWinner is set when the state is WON.
func (g *Game) GetWinner() *Player

// PlaceDisk places the disc in the lowest free row in the column.
func (b *Board) PlaceDisk(d enum.Disc, col int) (row int, err error)

func (b *Board) IsFull() bool

// CountInDirection is used by rules.
func (b *Board) CountInDirection(row, col int, dir enum.Direction, count int) int

func (p *Player) GetName() string

func (p *Player) GetDisk() enum.Disc
```
