---
title: Connect Four
description: Object-oriented design for a two-player Connect Four game.
next: /problems/connect-four/variation-1/
---

# Connect Four

<p class="lead">Build the object-oriented design for a two-player Connect Four game. Players take turns dropping discs into a 7-column, 6-row board. The first to align four of their own discs vertically, horizontally, or diagonally wins.</p>

## How this problem is organized

Each **variation** is a self-contained trail:

1. **Requirements** — interview-style Q&A and consolidated scope  
2. **Design** — entities, relationships, and class sketches  
3. **Codebase** — how the Go package implements the design  

Later variations can add concurrency, networking, or rule plug-ins without rewriting the earlier trail.

| Variation | Scope | Status |
| --- | --- | --- |
| [Variation 1 — In-memory OOP](/problems/connect-four/variation-1/) | 7×6 board, two players, CLI, single-threaded | Implemented in [`connect_four/`](https://github.com/gauxs/lld/tree/main/connect_four) |
| Variation 2 (planned) | Networked multiplayer, authoritative server | [Follow-ups](/problems/connect-four/follow-ups) |

## Package map

```text
connect_four/
├── game.go          # turn order, state machine, MakeMove
├── board.go         # grid, gravity, win counting
├── player.go
├── rule.go          # Rule interface, FourRule
├── enum/            # Disc, GameState, Direction
└── connectfour.go   # CLI demo (Execute)
```

Start the trail: [Variation 1](/problems/connect-four/variation-1/).
