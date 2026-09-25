# Follow-up requirements

Future variations (not implemented in the baseline).

## FU-1: Networked multiplayer

<details>
<summary>Questions to ask</summary>

- **Who owns state?** An authoritative server; clients call `MakeMove` remotely.
- **How is state pushed?** WebSockets or SSE with full or delta snapshots.

</details>

**Requirement:** Serialize moves through a single writer; clients are read-only except their move RPC.

## FU-2: Spectators

<details>
<summary>Questions to ask</summary>

- **Can spectators move?** No—subscribe to read-only game state.
- **Late join?** Send current snapshot + sequence id for ordering.

</details>

**Requirement:** Publish the same state stream as players see, without mutation APIs.

## FU-3: Configurable rules & board

**Requirement:** Reuse `NewGame(rows, cols)` and swap `Rule` implementations; add configuration and tests without changing turn flow.
