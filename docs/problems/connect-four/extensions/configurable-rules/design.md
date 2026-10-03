---
title: Connect Four — configurable rules & board (design)
problem: connect_four
extension: configurable_rules
prev:
  text: requirements
  link: /problems/connect-four/extensions/configurable-rules/requirements
---

**Builds on:** [baseline design](/problems/connect-four/extensions/baseline/design).

## Delta from baseline

The baseline code already accepts `(rows, cols)` and a `Rule` interface—this extension documents making that explicit in interviews and tests. `NewGame` accepts an optional `Rule` parameter or functional options such as `NewGame(rows, cols, WithRule(r))`. `FourRule` remains the default; add `ConnectFiveRule`, etc. Validation rejects `rows` and `cols` smaller than the win streak length.

## Go design sketch

Method bodies are intentionally omitted. This shows only the design delta.

```go
// FourRule is the default Rule.
type FourRule struct{}

func (r *FourRule) Satisfied(board *Board, row, col int) bool
func (r *FourRule) WinLength() int

// ConnectFiveRule is an additional Rule implementation.
type ConnectFiveRule struct{}

func (r *ConnectFiveRule) Satisfied(board *Board, row, col int) bool
func (r *ConnectFiveRule) WinLength() int

// SizedRule extends the baseline Rule with the metadata needed to validate
// board dimensions during construction.
type SizedRule interface {
	Rule
	WinLength() int
}

// GameOption configures construction and may reject an invalid option.
type GameOption func(config *gameConfig) error

func WithRule(rule SizedRule) GameOption

// NewGame wires the chosen Rule and rejects rows or cols smaller than the
// win streak length. Rule.Satisfied has an unchanged contract.
func NewGame(rows, cols int, opts ...GameOption) (*Game, error)
```
