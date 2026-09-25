---
title: Connect Four — networked multiplayer (design)
problem: connect_four
extension: networked
prev:
  text: requirements
  link: /problems/connect-four/extensions/networked/requirements
---

**Builds on:** [baseline design](/problems/connect-four/extensions/baseline/design).

Reference implementation not checked in for this extension—use the design delta in an interview or a follow-on PR.

## Delta from baseline

| Area | Change |
| --- | --- |
| **Game** | Unchanged domain logic; wrapped by a server-side `Match` or `Room` that owns one `*Game` |
| **Transport** | `GameService` with `Join`, `MakeMove(gameID, playerID, col)`, `Subscribe(gameID)` |
| **Clients** | Thin: render snapshots; call RPC for moves |
| **Ordering** | Monotonic `sequence` per match on every broadcast |

## New entities

| Entity | Responsibilities |
| --- | --- |
| **MatchRegistry** | Map `gameID` → room; create/join lifecycle |
| **Room** | Single writer queue, holds `*Game`, fans out events |
| **GameEvent** | Snapshot or `{ sequence, state, lastMove }` for clients |

## API (sketch)

| Method | Description |
| --- | --- |
| `CreateMatch() (gameID, error)` | New baseline game, two slots |
| `JoinMatch(gameID, playerID) error` | Bind player; start when full |
| `MakeMove(gameID, playerID, col) error` | Validate identity + turn; delegate to `Game.MakeMove` |
| `Subscribe(gameID) <-chan GameEvent` | Blocking stream for players (spectators extension narrows permissions) |
