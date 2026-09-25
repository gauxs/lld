---
title: Requirements — Variation 1
description: Functional scope, NFRs, and failure handling for Connect Four.
prev: /problems/connect-four/variation-1/
next: /problems/connect-four/variation-1/design
---

# Requirements gathering

## Functional

| Question | Answer |
| --- | --- |
| How does a player move? | Choose a column only; the disc falls to the lowest empty row. |
| Win condition? | Four consecutive own discs in a row, column, or either diagonal. |
| Board full with no winner? | Draw. |
| Moves after the game ends? | Rejected. |
| Invalid column (full or out of range)? | Rejected; player retries. |
| Wrong player or wrong game state? | Rejected. |

## YAGNI (later, not v1)

| Idea | v1 | Future |
| --- | --- | --- |
| Undo last move | No | Possible extension |
| More than two players | No | Supported in design discussions |
| Board size | Fixed 7×6 for the story | Constructor already takes rows/cols |
| Win rule | Four in a line | `Rule` interface for count + directions |

## Non-functional

- **Single-threaded** — no locking in this variation.  
- **In-memory** — no persistence between runs.

## Consolidated requirements

1. 7 × 6 board (6 rows, 7 columns in code: `NewGame(6, 7)`).  
2. Two players with distinct disc colors.  
3. Alternating turns until win or draw.  
4. Gravity: lowest available cell in the column.  
5. Win: four in a line (horizontal, vertical, both diagonals).  
6. No moves after terminal state (`WON` or `DRAW`).

## Source notes

Authoring notes also live in the repo at [`connect_four/docs/requirement.md`](https://github.com/gauxs/lld/blob/main/connect_four/docs/requirement.md).

Next: [Design](/problems/connect-four/variation-1/design).
