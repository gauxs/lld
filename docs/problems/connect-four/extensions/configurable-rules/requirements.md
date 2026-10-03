---
title: Connect Four — configurable rules & board
pageClass: lld-req-page
problem: connect_four
extension: configurable_rules
next:
  text: design
  link: /problems/connect-four/extensions/configurable-rules/design
---

**Builds on:** [local two-player baseline](/problems/connect-four/extensions/baseline/requirements).

Same turn flow and column-drop mechanics, but board dimensions and win condition are configurable.

## Functional requirements

### FR-1: Board dimensions

`NewGame(rows, cols)` defines the grid; defaults remain 6×7 for backward compatibility.

> **Interview prompts**
>
> - **Minimum size?** At least `winLength` columns/rows for a meaningful game—validate at construction.
> - **Non-rectangular?** No.

### FR-2: Pluggable win rule

Win detection is delegated to a `Rule` implementation (default: four in a line).

> **Interview prompts**
>
> - **Connect five?** New `Rule` type; `Game` unchanged except injected dependency.
> - **Custom patterns?** Same interface; swap implementation.

## Out of scope (this extension)

- None stated.

## Non-functional requirements

### NFR-1: Testability

Rules and boards are covered by table-driven tests without network or CLI.

## Extensions from here

No further extensions in this branch.
