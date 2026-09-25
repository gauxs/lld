# Connect Four — follow-ups

Future variations (not implemented in the baseline).

<div class="lld-req">

### FU-1: Networked multiplayer

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Who owns state?** An authoritative server; clients call `MakeMove` remotely.
- **How is state pushed?** WebSockets or SSE with full or delta snapshots.

</div>
</details>

**Outcome:** Serialize moves through a single writer; clients are read-only except their move RPC.

</div>

<div class="lld-req">

### FU-2: Spectators

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Can spectators move?** No—subscribe to read-only game state.
- **Late join?** Send current snapshot + sequence id for ordering.

</div>
</details>

**Outcome:** Publish the same state stream as players see, without mutation APIs.

</div>

<div class="lld-req">

### FU-3: Configurable rules & board

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **What varies?** Board dimensions and win rule (`Rule` implementation).
- **What stays stable?** Turn flow and `MakeMove` orchestration.

</div>
</details>

**Outcome:** Reuse `NewGame(rows, cols)` and swap `Rule` implementations with tests and configuration hooks.

</div>
