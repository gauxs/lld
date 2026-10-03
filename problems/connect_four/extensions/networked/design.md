# Connect Four — networked multiplayer (design)

**Builds on:** [baseline design](/problems/connect-four/extensions/baseline/design).

Reference implementation not checked in for this extension—use the design delta in an interview or a follow-on PR.

## Delta from baseline

`Game` keeps unchanged domain logic and is wrapped by a server-side `Match` or `Room` that owns one `*Game`. The transport is a `GameService` with `Join`, `MakeMove(gameID, playerID, col)`, and `Subscribe(gameID)`. Clients stay thin: they render snapshots and call RPC for moves. Every broadcast has a monotonic `sequence` per match.

## Go design sketch

Method bodies are intentionally omitted. This shows only the design delta.

```go
// MatchRegistry maps gameID to room and owns the create/join lifecycle.
type MatchRegistry struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

// Room is the single writer queue for a match. It holds one baseline Game and
// fans out events.
type Room struct {
	game        *Game
	sequence    uint64
	commands    chan roomCommand
	subscribers map[string]chan GameEvent
}

// GameEvent is a snapshot or {sequence, state, lastMove} for clients.
type GameEvent struct {
	Snapshot Snapshot
	Sequence uint64
	State    enum.GameState
	LastMove *Move
}

// GameService is the transport for creating and joining matches, making moves,
// and subscribing to match events.
type GameService interface {
	// CreateMatch creates a new baseline game with two slots.
	CreateMatch() (gameID string, err error)

	// JoinMatch binds a player and starts the game when full.
	JoinMatch(gameID, playerID string) error

	// MakeMove validates identity and turn, then delegates to Game.MakeMove.
	MakeMove(gameID, playerID string, col int) error

	// Subscribe returns a buffered stream. Cancelling ctx removes the
	// subscription so disconnected clients do not leak.
	Subscribe(ctx context.Context, gameID string) (<-chan GameEvent, error)
}
```

Room broadcasts immutable snapshots without blocking the command loop. A slow
subscriber keeps only the latest buffered snapshot or is disconnected.
