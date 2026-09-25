# Connect Four — configurable rules & board (design)

**Builds on:** [baseline design](/problems/connect-four/extensions/baseline/design).

The baseline code already accepts `(rows, cols)` and a `Rule` interface—this extension documents making that explicit in interviews and tests.

## Delta from baseline

| Area | Change |
| --- | --- |
| **NewGame** | Accept optional `Rule` parameter or functional options: `NewGame(rows, cols, WithRule(r))` |
| **FourRule** | Default; add `ConnectFiveRule`, etc. |
| **Validation** | Reject `rows`, `cols` smaller than win streak length |

## API (sketch)

| Method | Description |
| --- | --- |
| `NewGame(rows, cols int, opts ...GameOption) *Game` | Wire chosen `Rule` |
| `Rule.Satisfied(b *Board, row, col int) bool` | Unchanged contract |
