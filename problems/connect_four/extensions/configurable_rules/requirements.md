# Connect Four — configurable rules & board

**Builds on:** [local two-player baseline](/problems/connect-four/extensions/baseline/requirements).

Same turn flow and column-drop mechanics, but board dimensions and win condition are configurable.

## Functional requirements

<div class="lld-req">

### FR-1: Board dimensions

`NewGame(rows, cols)` defines the grid; defaults remain 6×7 for backward compatibility.

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Minimum size?** At least `winLength` columns/rows for a meaningful game—validate at construction.
- **Non-rectangular?** No.

</div>
</details>

</div>

<div class="lld-req">

### FR-2: Pluggable win rule

Win detection is delegated to a `Rule` implementation (default: four in a line).

<details class="lld-reveal">
<summary><span class="lld-reveal-icon" aria-hidden="true"></span>Interview prompts</summary>

<div class="lld-reveal-inner">

- **Connect five?** New `Rule` type; `Game` unchanged except injected dependency.
- **Custom patterns?** Same interface; swap implementation.

</div>
</details>

</div>

## Non-functional requirements

<div class="lld-req">

### NFR-1: Testability

Rules and boards are covered by table-driven tests without network or CLI.

</div>

## Extensions from here

No further extensions in this branch.
