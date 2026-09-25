---
title: Design
prev: /problems/connect-four/functional-requirement
next: /problems/connect-four/codebase
---

## Entities & responsibilities

| Entity | Responsibilities |
| --- | --- |
| **Game** | Player registration, game lifecycle (`NOT_PLAYING` → `PLAYING` → `WON` \| `DRAW`), turn order, orchestrates `MakeMove` |
| **Board** | Grid storage, column validity, gravity placement, directional line counting from last move |
| **Player** | Display name and disc color |
| **Rule** | Whether the last placement satisfies the win condition (pluggable) |

## API

### Game

| Method | Description |
| --- | --- |
| `NewGame(rows, cols int) *Game` | Empty board, default `FourRule` |
| `AddPlayer(name string, d enum.Disc) error` | Register before start; unique disc colors |
| `StartGame() error` | Requires two players |
| `MakeMove(col int) error` | Current player drops in `col`; updates state / turn |
| `GetGameState() enum.GameState` | |
| `GetCurrentPlayer() *Player` | |
| `GetWinner() *Player` | Set when state is `WON` |

### Board

| Method | Description |
| --- | --- |
| `PlaceDisk(d enum.Disc, col int) (row int, error)` | Lowest free row in column |
| `IsFull() bool` | |
| `CountInDirection(row, col, dir, count int) int` | Used by rules |

### Rule

| Method | Description |
| --- | --- |
| `Satisfied(b *Board, row, col int) bool` | Win check after a placement |

### Player

| Method | Description |
| --- | --- |
| `GetName() string` | |
| `GetDisk() enum.Disc` | |
