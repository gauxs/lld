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

Split `Subscribe` into `SubscribePlayer` versus `SubscribeSpectator`, or use one stream with the role in the session. For authorization, `MakeMove` checks `RolePlayer`; spectators get `403`. Broadcasts use the same `GameEvent` payload; optional redaction (hidden until start) is out of scope.

## Go design sketch

Method bodies are intentionally omitted. This shows only the design delta.

```go
// SpectatorGameService preserves the networked operations while making every
// subscription role-aware. There is no unauthenticated base Subscribe method.
type SpectatorGameService interface {
	CreateMatch() (gameID string, err error)
	JoinMatch(gameID, playerID string) error
	MakeMove(gameID, playerID string, col int) error
	SubscribePlayer(
		ctx context.Context,
		gameID, playerID string,
	) (<-chan GameEvent, error)
	SubscribeSpectator(
		ctx context.Context,
		gameID, spectatorID string,
	) (<-chan GameEvent, error)
}
```

Spectator channels use the same bounded, nonblocking snapshot policy as player
subscriptions, so fan-out cannot stall the room's single-writer command loop.
