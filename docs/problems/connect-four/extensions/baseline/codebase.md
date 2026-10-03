---
title: Local two-player — codebase
problem: connect_four
extension: baseline
prev:
  text: design
  link: /problems/connect-four/extensions/baseline/design
---

# Codebase

Source: [GitHub](https://github.com/gauxs/lld/tree/main/problems/connect_four/extensions/baseline/code)

Path: ` problems/connect_four/extensions/baseline/code `

## Directory structure

```text
code/
├── enum/
│   ├── direction.go
│   ├── disk.go
│   └── gamestate.go
├── board.go
├── connectfour.go
├── error.go
├── game.go
├── player.go
└── rule.go
```

## ` board.go `

```go
package connectfour

import (
	"fmt"

	"github.com/gauxs/lld/problems/connect_four/extensions/baseline/code/enum"
)

var directionRowColSteps = map[enum.Direction][][]int{
	enum.DIRECTION_HORIZONTAL:     {{0, 1}, {0, -1}},
	enum.DIRECTION_VERTICAL:       {{1, 0}, {-1, 0}},
	enum.DIRECTION_DIAGONAL_BACK:  {{-1, -1}, {1, 1}},
	enum.DIRECTION_DIAGONAL_FRONT: {{-1, 1}, {1, -1}},
}

type Board struct {
	row    int
	column int
	grid   [][]enum.Disc
}

func NewBoard(r int, c int) *Board {
	g := make([][]enum.Disc, r)

	for i := 0; i < r; i++ {
		g[i] = make([]enum.Disc, c)
	}

	return &Board{
		row:    r,
		column: c,
		grid:   g,
	}
}

func (b *Board) getNextFreeSlot(col int) int {
	for r := b.row - 1; r >= 0; r-- {
		if b.grid[r][col] == enum.DISK_INVALID {
			return r
		}
	}

	return -1
}

func (b *Board) PlaceDisk(d enum.Disc, col int) (int, error) {
	if col < 0 || col >= b.column {
		return -1, fmt.Errorf("column %d: %w", col, ErrInvalidMove)
	}

	r := b.getNextFreeSlot(col)
	if r == -1 {
		return -1, fmt.Errorf("column %d if full: %w", col, ErrInvalidMove)
	}

	b.grid[r][col] = d
	return r, nil
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

func (b *Board) CountInDirection(r int, c int, d enum.Direction, limit int) int {
	if c < 0 || c >= b.column || r < 0 || r >= b.row || b.grid[r][c] == enum.DISK_INVALID {
		return 0
	}

	count := 0
	for _, rowColStep := range directionRowColSteps[d] {
		count += b.countInDirection(r+rowColStep[0], c+rowColStep[1], limit,
			b.grid[r][c], rowColStep[0], rowColStep[1])
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
```

## ` connectfour.go `

```go
package connectfour

import (
	"fmt"

	"github.com/gauxs/lld/problems/connect_four/extensions/baseline/code/enum"
)

func Execute() {
	fmt.Println("Starting the game ConnectFour")

	game := NewGame(6, 7)
	if err := game.AddPlayer("Player-A", enum.DISK_RED); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := game.AddPlayer("Player-B", enum.DISK_BLUE); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := game.StartGame(); err != nil {
		fmt.Println(err.Error())
		return
	}

	for game.GetGameState() == enum.GAMESTATE_PLAYING {
		col := 0
		curPlayer := game.GetCurrentPlayer()
		fmt.Printf("Enter column for player %v: ", curPlayer.GetName())
		fmt.Scanln(&col)

		if err := game.MakeMove(col); err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Printf("Game ended in %v \n", game.GetGameState().String())
	if game.GetGameState() == enum.GAMESTATE_WON {
		fmt.Printf("The winner of te game is: %v \n", game.GetWinner().GetName())
	}
}
```

## ` enum/direction.go `

```go
package enum

type Direction int

const (
	DIRECTION_INVALID Direction = iota
	DIRECTION_HORIZONTAL
	DIRECTION_VERTICAL
	DIRECTION_DIAGONAL_FRONT
	DIRECTION_DIAGONAL_BACK
)
```

## ` enum/disk.go `

```go
package enum

type Disc int

const (
	DISK_INVALID = iota
	DISK_RED
	DISK_BLUE
)
```

## ` enum/gamestate.go `

```go
package enum

type GameState int

const (
	GAMESTATE_INVALID = iota
	GAMESTATE_NOT_PLAYING
	GAMESTATE_PLAYING
	GAMESTATE_WON
	GAMESTATE_DRAW
)

func (gs GameState) String() string {
	switch gs {
	case GAMESTATE_NOT_PLAYING:
		return "Not Playing"
	case GAMESTATE_PLAYING:
		return "Playing"
	case GAMESTATE_WON:
		return "Won"
	case GAMESTATE_DRAW:
		return "Draw"
	}

	return "Invalid"
}
```

## ` error.go `

```go
package connectfour

import "errors"

var ErrInvalidMove = errors.New("invalid move")
var ErrDiskAlreadyTaken = errors.New("disk with this color is already taken")
var ErrGameNotInCorrectState = errors.New("game is not in correct state")
var ErrInsufficientPlayers = errors.New("insufficient players")
```

## ` game.go `

```go
package connectfour

import (
	"github.com/gauxs/lld/problems/connect_four/extensions/baseline/code/enum"
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
```

## ` player.go `

```go
package connectfour

import (
	"github.com/gauxs/lld/problems/connect_four/extensions/baseline/code/enum"
)

type Player struct {
	name string
	disk enum.Disc
}

func NewPlayer(name string, d enum.Disc) *Player {
	return &Player{
		name: name,
		disk: d,
	}
}

func (p *Player) GetName() string {
	return p.name
}
func (p *Player) GetDisk() enum.Disc {
	return p.disk
}
```

## ` rule.go `

```go
package connectfour

import "github.com/gauxs/lld/problems/connect_four/extensions/baseline/code/enum"

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
```
