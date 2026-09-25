---
title: Connect Four — networked multiplayer
pageClass: lld-req-page
problem: connect_four
extension: networked
next:
  text: design
  link: /problems/connect-four/extensions/networked/design
---

**Builds on:** [local two-player baseline](/problems/connect-four/extensions/baseline/requirements).

Players connect to a shared game over the network. One process owns board state; clients submit moves and receive updates.

## Functional requirements

<div class="lld-req">

### FR-1: Authoritative game server

The server holds the canonical `Game` instance. Clients never mutate board state locally except for optimistic UI (optional, out of scope here).

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Split-brain if two servers?** Single writer per match id; stick to one room owner.
- **Reconnect mid-game?** Client sends last known sequence; server returns snapshot or deltas.

</div>
</details>

</div>

<div class="lld-req">

### FR-2: Remote move

A player invokes `MakeMove(col)` over RPC/HTTP; the server validates turn, applies the baseline rules, and returns success or error.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Move on wrong turn?** Reject with error; no board change.
- **Idempotency?** Client move id or sequence number to ignore duplicates.

</div>
</details>

</div>

<div class="lld-req">

### FR-3: State broadcast

After each accepted move (and on join), subscribers receive enough state to render the board and game status.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Full snapshot vs delta?** Start with full snapshot + monotonic sequence id.
- **Transport?** WebSockets or SSE; interview choice either way if ordering is clear.

</div>
</details>

</div>

## Non-functional requirements

<div class="lld-req">

### NFR-1: Concurrency

All mutations for a given match go through one goroutine or a mutex—baseline `Game` stays logically single-threaded per room.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Many concurrent games?** One lock (or actor) per `gameID`, not one global lock.
- **Read-heavy spectators?** See spectators extension; reads can fan out from immutable snapshots.

</div>
</details>

</div>

<div class="lld-req">

### NFR-2: Failure handling

Network and validation errors are visible to the caller; server does not apply partial moves.

</div>

## Extensions from here

| Extension | What it adds |
| --- | --- |
| [Spectators](/problems/connect-four/extensions/spectators/requirements) | Read-only subscribers on the same state stream |
