---
title: Connect Four — local two-player
pageClass: lld-req-page
problem: connect_four
extension: baseline
prev:
  text: scarcity
  link: /learn/concurrency/02-problems/scarcity
next:
  text: design
  link: /problems/connect-four/extensions/baseline/design
---

Build the object-oriented design for a two-player Connect Four game. Players take turns dropping discs into a 7-column, 6-row board. The first to align four of their own discs vertically, horizontally, or diagonally wins.

## Functional requirements

<div class="lld-req">

### FR-1: Column drop

Discs are placed by choosing a column; the piece occupies the lowest empty row in that column.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Can the player choose a row?** No—column only; gravity fills the bottom slot.
- **What if the column is full?** Reject the move; the active player tries again.

</div>
</details>

</div>

<div class="lld-req">

### FR-2: Win detection

Four consecutive discs of the same player in a horizontal, vertical, or diagonal line wins the game.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Both diagonal directions?** Yes.
- **Who wins on the connecting move?** The player who placed the fourth disc.

</div>
</details>

</div>

<div class="lld-req">

### FR-3: Draw

If the board is full and no player has won, the game is a draw.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Is a draw a terminal state?** Yes—no further moves.

</div>
</details>

</div>

<div class="lld-req">

### FR-4: Turn order

Exactly two players alternate turns while the game is in progress.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Can the same player move twice?** No.
- **Moves after the game ends?** Rejected.

</div>
</details>

</div>

<div class="lld-req">

### FR-5: Invalid input

Moves out of range, into a full column, or while the game is not in `PLAYING` state are rejected with an error.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Wrong player in API?** This extension uses `MakeMove(col)` for the current player only.
- **Silent ignore vs error?** Return an error; CLI prompts again.

</div>
</details>

</div>

## Out of scope (this extension)

<div class="lld-req lld-req--scope">

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Undo?** No.
- **More than two players?** Not in this extension.
- **Remote play?** See [networked extension](/problems/connect-four/extensions/networked/requirements).
- **Configurable board or win rule?** See [configurable rules extension](/problems/connect-four/extensions/configurable-rules/requirements).

</div>
</details>

</div>

## Non-functional requirements

<div class="lld-req">

### NFR-1: Concurrency model

Single-threaded, in-process execution—no locking in this extension.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Multiple threads?** Not here; the networked extension adds server-side serialization.
- **Reentrancy?** Not required for CLI-driven play.

</div>
</details>

</div>

<div class="lld-req">

### NFR-2: Persistence

All state lives in memory for the lifetime of the game instance.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Save / resume game?** Out of scope.
- **Crash recovery?** Not required.

</div>
</details>

</div>

<div class="lld-req">

### NFR-3: Failure handling

Invalid operations surface as errors instead of silent no-ops.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Return codes vs exceptions?** Go errors returned to caller.
- **User-facing messages?** CLI prints error string; API stays domain-focused.

</div>
</details>

</div>

## Extensions from here

| Extension | What it adds |
| --- | --- |
| [Networked multiplayer](/problems/connect-four/extensions/networked/requirements) | Authoritative server, remote `MakeMove`, live state sync |
| [Configurable rules & board](/problems/connect-four/extensions/configurable-rules/requirements) | Pluggable win rule and non-default board sizes |
