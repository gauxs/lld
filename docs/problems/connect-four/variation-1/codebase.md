---
title: Codebase — Variation 1
description: How the connect_four package maps to the design.
prev: /problems/connect-four/variation-1/design
next: /problems/connect-four/follow-ups
---

# Codebase walkthrough

Package import path: `github.com/gauxs/lld/connect_four` (package name `connectfour`).

## Entry point

`main.go` delegates to the demo driver:

```go
func main() {
    connectfour.Execute()
}
```

`connectfour.Execute()` in `connect_four/connectfour.go` builds a 6×7 game, registers two players, runs a read loop on stdin, and prints the terminal outcome.

## Game (`game.go`)

- `NewGame(boardX, boardY)` creates board + default `FourRule`.  
- `AddPlayer` enforces `NOT_PLAYING`, unique disc colors, sets first player as current.  
- `StartGame` requires two players.  
- `MakeMove(col)` checks state, delegates placement to `Board`, runs `winRule.Satisfied`, advances turn or sets `DRAW` when full.

## Board (`board.go`)

- `grid` is `[][]enum.Disc` with `DISK_INVALID` as empty.  
- `getNextFreeSlot` scans from bottom row upward.  
- `CountInDirection` walks both rays from `(row, col)` using a direction → step map (horizontal, vertical, both diagonals).

## Win rule (`rule.go`)

```go
type Rule interface {
    Satisfied(b *Board, row int, col int) bool
}
```

`FourRule` checks four directions with `consecutiveCount: 4`. Swap the `Rule` implementation to vary win conditions without editing board logic.

## Supporting types

| Path | Role |
| --- | --- |
| `player.go` | `Player` struct |
| `error.go` | Domain errors (`ErrInvalidMove`, game state errors) |
| `enum/disk.go`, `enum/gamestate.go`, `enum/direction.go` | Shared enums |

## Run locally

```bash
go run .   # from repo root; follow column prompts
```

## Related docs in repo

- [`connect_four/docs/requirement.md`](https://github.com/gauxs/lld/blob/main/connect_four/docs/requirement.md) — original interview notes  
- [`connect_four/docs/followups.md`](https://github.com/gauxs/lld/blob/main/connect_four/docs/followups.md) — networking / spectators

Next: [Follow-ups](/problems/connect-four/follow-ups) for variation ideas.
