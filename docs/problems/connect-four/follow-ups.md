---
title: Connect Four — Follow-ups
description: Extensions for future variations (networking, spectators, concurrency).
prev: /problems/connect-four/variation-1/codebase
---

# Follow-ups

Ideas for **variation 2+**—not implemented in variation 1.

## Networked multiplayer

- Run game logic on an **authoritative server**; clients call `MakeMove` over RPC or WebSocket.  
- Server serializes moves → naturally solves correctness for shared board state.  
- Clients receive pushed game state snapshots.

Touches [Coordination](/learn/concurrency/coordination) more than correctness in a single mutex.

## Spectators

- Publish read-only game state events (same channel as players, minus mutation APIs).  
- Consider versioning state or sequence numbers so late joiners catch up.

## Configurable rules and board

Already partially modeled:

- `NewGame(rows, cols)` for board size.  
- `Rule` interface for win detection.  

A new variation trail would add requirements + design for admin configuration and tests.

## Source

Notes from [`connect_four/docs/followups.md`](https://github.com/gauxs/lld/blob/main/connect_four/docs/followups.md).
