---
title: Connect Four — spectators (design)
problem: connect_four
extension: spectators
prev:
  text: requirements
  link: /problems/connect-four/extensions/spectators/requirements
---

**Builds on:** [networked design](/problems/connect-four/extensions/networked/design).

## Delta from networked

| Area | Change |
| --- | --- |
| **Subscribe** | Split into `SubscribePlayer` vs `SubscribeSpectator`, or one stream with role in session |
| **Authorization** | `MakeMove` checks `RolePlayer`; spectators get `403` |
| **Broadcast** | Same `GameEvent` payload; optional redaction (hidden until start) is out of scope |

## API (sketch)

| Method | Description |
| --- | --- |
| `WatchMatch(gameID, spectatorID) (<-chan GameEvent, error)` | Read-only stream |
| `MakeMove(...)` | Unchanged; rejects non-players |
