---
title: Variation 1 — In-memory OOP
description: Single-process, single-threaded Connect Four with extensible win rules.
prev: /problems/connect-four/
next: /problems/connect-four/variation-1/requirements
---

# Variation 1 — In-memory OOP

<div class="trail-meta">
  <span>Single process</span>
  <span>Single-threaded</span>
  <span>In-memory</span>
  <span>Go package: connect_four</span>
</div>

<p class="lead">Baseline LLD: model the game in one process, enforce rules in code, and keep win detection pluggable for future rule changes.</p>

## Trail

| Step | Page | Outcome |
| --- | --- | --- |
| 1 | [Requirements](/problems/connect-four/variation-1/requirements) | Agreed scope and error behavior |
| 2 | [Design](/problems/connect-four/variation-1/design) | Entities and public APIs |
| 3 | [Codebase](/problems/connect-four/variation-1/codebase) | Mapping to `github.com/gauxs/lld/connect_four` |

## Out of scope (this variation)

- Undo / take-back  
- More than two players  
- Configurable board size (constructor accepts dimensions, but interview scope is 7×6)  
- Network or spectator clients  

These are captured under [Follow-ups](/problems/connect-four/follow-ups) for a future variation.
