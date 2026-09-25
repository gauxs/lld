# Functional & non-functional requirements

## Functional requirements

### FR-1: Column-only moves

Discs are placed by choosing a column; the piece occupies the lowest empty row in that column.

<details>
<summary>Questions to ask</summary>

- **Can the player choose a row?** No—column only; gravity fills the bottom slot.
- **What if the column is full?** Reject the move; the active player tries again.

</details>

### FR-2: Win detection

Four consecutive discs of the same player in a horizontal, vertical, or diagonal line wins the game.

<details>
<summary>Questions to ask</summary>

- **Both diagonal directions?** Yes.
- **Who wins on the connecting move?** The player who placed the fourth disc.

</details>

### FR-3: Draw

If the board is full and no player has won, the game is a draw.

<details>
<summary>Questions to ask</summary>

- **Is a draw a terminal state?** Yes—no further moves.

</details>

### FR-4: Turn order

Exactly two players alternate turns while the game is in progress.

<details>
<summary>Questions to ask</summary>

- **Can the same player move twice?** No.
- **Moves after the game ends?** Rejected.

</details>

### FR-5: Invalid input

Moves out of range, into a full column, or while the game is not in `PLAYING` state are rejected with an error.

<details>
<summary>Questions to ask</summary>

- **Wrong player object in API?** This variation uses `MakeMove(col)` on behalf of the current player only.
- **Silent ignore vs error?** Return an error; CLI prompts again.

</details>

## Out of scope (this version)

<details>
<summary>Questions to ask (YAGNI)</summary>

- **Undo?** No.
- **More than two players?** Not now; registration API can be extended later.
- **Configurable board size?** Constructor accepts dimensions; interview story uses 6×7.
- **Configurable win rule?** Use a `Rule` strategy; default is four in a line.

</details>

## Non-functional requirements

| ID | Requirement |
| --- | --- |
| NFR-1 | Single-threaded, in-process |
| NFR-2 | In-memory state only |
| NFR-3 | Errors over silent failure for invalid moves |

## Consolidated checklist

1. 7 columns × 6 rows (`NewGame(6, 7)` — rows × cols in code).
2. Two players with distinct disc colors.
3. Alternating turns until win or draw.
4. Gravity placement.
5. Four in a line wins (all directions).
6. No moves after terminal state.
