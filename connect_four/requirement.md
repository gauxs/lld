# Problem Statement
> Build the object-oriented design for a two-player Connect Four game. Players take turns dropping discs into a 7-column, 6-row board. The first to align four of their own discs vertically, horizontally, or diagonally wins.

# [Step: I] Requirement Gathering
### Functional requirement
> Player can choose only column?
Yes and the last empty slot will be filled

> What happens if the board fills w/o a winner?
Its a draw

> Both diagonals?
Yes

### What we won’t we do (this version) and what we might extend later (YAGNI - Design with extension in mind, but only implement what's needed now) ?
> Can the player take-back its last move?
No undo

> Can the numerb of players increase?
Yes, but not in this version

> Can the board size change?
Yes, but not in this version

> Can the win rule change? Escpecially the number of disks?
Yes, configurable consecutive-piece count and winning directions.

### Non-functional requirement: Assumptions & constraints (NFR)
> Single threaded? 
Yes

> In memory?
Yes

### What fails and how do we handle it?
> What happens if player moves out of turn or after game is in terminal state?
Do no allow the player to make the move

> What happens if player keeps the disc on a full column or non-existent column?
Do no allow the player to make the move

## Requirements
1. 7 × 6 board
2. Two players
3. Players take turns dropping pieces into columns.
4. Pieces occupy the lowest available position.
5. Four consecutive pieces wins.
6. Horizontal, vertical, and diagonal wins.
7. Full board without winner = draw.
8. No moves after game completion.
9. Future support for more players.
10. Future configurable consecutive-piece count and winning directions.

# [Step: II] Entities and relationships
## Entities
1. Game: Holds the board and players. Manages game's state. Enforces players turns. 
2. Board: Represents a board. The grid where the discs will be placed. Checks the validity of moves. Knows if the column / board is full. Can check if 4 discs are connected.
3. Player: Represents a player. Holds the name and disc color.

# [Step: III] Class design

### Player
```go
type Disc int

const(
    DISK_RED = iota
    DISK_BLUE
)

type Player struct{
    name string
    color Disc
}

func (p *Player) GetName() string{}
func (p *Player) GetColor() string{}
```

### Board
```go
type Direction int

const(
    Direction_Horizontal = iota
    Direction_Vertical
    Direction_Diagonal
)

type Board struct{
    row int
    col int
    grid [][]Disc
}

func (b *Board) PlaceDisk(d Disc, col int) error{}
func (b *Board) CountConsecutiveInDirection(col int, d Direction) int{}
func (b *Board) IsFull() bool{}
```

### Game
```go
type GameState int

const(
    GameState_Not_Playing = iota
    GameState_Playing
    GameState_Won
    GameState_Draw
)

type Game struct{
    players []*Player
    board *Board
    state GameState
    currentPlayer *Player
    winner *Player
}

func (g *Game) WinningPlayer() *Player{}
func (g *Game) GetGameState() GameState{}
func (g *Game) GetCurrentPlayer() *Player{}
func (g *Game) GetWinner() *Player{}
func (g *Game) MakeMove(*Player, col int) error{}
```