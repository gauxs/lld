---
title: Connect Four
---

# Connect Four

Extensions are **siblings** in the sidebar (baseline first). Use **Builds on** or the graph to see dependencies.

```mermaid
flowchart TD
  baseline["Local two-player"]
  networked["Networked multiplayer"]
  configurable_rules["Configurable rules & board"]
  spectators["Spectators"]
  baseline --> networked
  baseline --> configurable_rules
  networked --> spectators
```

| Extension | Builds on |
| --- | --- |
| [Local two-player](/problems/connect-four/extensions/baseline/requirements) | — |
| [Networked multiplayer](/problems/connect-four/extensions/networked/requirements) | [Local two-player](/problems/connect-four/extensions/baseline/requirements) |
| [Configurable rules & board](/problems/connect-four/extensions/configurable-rules/requirements) | [Local two-player](/problems/connect-four/extensions/baseline/requirements) |
| [Spectators](/problems/connect-four/extensions/spectators/requirements) | [Networked multiplayer](/problems/connect-four/extensions/networked/requirements) |
