---
title: Connect Four — spectators
pageClass: lld-req-page
problem: connect_four
extension: spectators
next:
  text: design
  link: /problems/connect-four/extensions/spectators/design
---

**Builds on:** [networked multiplayer](/problems/connect-four/extensions/networked/requirements).

Observers can watch a live match without playing.

## Functional requirements

### FR-1: Read-only subscription

Spectators receive the same state stream as players but have no `MakeMove` API (or calls are rejected).

> **Interview prompts**
>
> - **Can spectators chat?** Out of scope unless asked.
> - **Late join?** Send latest snapshot + sequence so UI catches up.

### FR-2: Identity

Spectators authenticate or receive a guest token so rate limits and room caps can apply.

> **Interview prompts**
>
> - **Anonymous ok?** Yes for interviews; mention abuse/capacity if probed.

## Non-functional requirements

### NFR-1: Fan-out

Many spectators must not block the single writer path—publish from an immutable snapshot copy.

## Extensions from here

No further extensions in this branch.
