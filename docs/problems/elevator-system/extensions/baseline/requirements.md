---
title: Elevator System
pageClass: lld-req-page
problem: elevator_system
extension: baseline
next:
  text: design
  link: /problems/elevator-system/extensions/baseline/design
---

Design and implement a normal passenger elevator system like you'd see in a gated community/residential apartment complex.

## Functional requirements

### FR-1: Number of elevators & Floors

- Number of elevators: configurable at system initialization
- Number of floors: configurable at system initialization
- Both remain fixed after initialization
- There is a hall/lobby on every floor containing all elevator entrances.

> **Interview prompts**
>
> - **Question?** Agreed answer.

### FR-2: User interaction

- Outside: each floor has Up and Down buttons.
- Inside: the elevator has a button for each floor.

### FR-3: Elevator Selection

- Nearest non-moving elevator as the initial selection strategy. We won't optimize/batch multiple requests in the initial implementation.
- Selection logic should be replaceable/extensible in the future

### FR-3: Elevator Display

- Each elevator has its own external display, showing that elevator's current floor.
- The floor displays outside show the current floor of the elevator(s) serving that floor.

### FR-4: Multiple Elevators

- Multiple elevators can be moving simultaneously; their movements are independent.

## Out of scope (this extension)

- Deferred capability.

## Non-functional requirements

### NFR-1: Short title

Constraint or quality requirement.

## Extensions from here

- [Child extension](/problems/problem-name/extensions/child/requirements) — What it adds
