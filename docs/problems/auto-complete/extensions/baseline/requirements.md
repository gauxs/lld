---
title: User activity tracker
pageClass: lld-req-page
problem: auto_complete
extension: baseline
next:
  text: design
  link: /problems/auto-complete/extensions/baseline/design
---

Design and implement a Search Autocomplete System


## Functional requirements

### FR-1: Return suggestions when a user types a prefix. The number of suggestions should be configurable.

### FR-3: Words are explicitly added to the system. Search queries do not automatically add new words.

### FR-4: Suggestions are ranked by frequency of explicit additions: a word that has been added more frequently ranks higher for a given prefix.

### FR-5: The ranking strategy should be replaceable/configurable

### FR-6: Adding the same word multiple times increments its frequency rather than creating duplicate entries.

### FR-7: Empty query should be supported; return the top-N suggestions globally.

### FR-8: Non-alphabetic input: reject/ignore characters outside a-z after normalization.

## Out of scope (this extension)

## Non-functional requirements

### NFR-1: concurrent additions and retrievals are supported.

### NFR-2: Query latency: autocomplete retrieval should be low-latency; target is O(query length + number of returned suggestions)

### NFR-3: When a word is added or its frequency changes, subsequent retrievals should reflect the updated ranking without requiring a full rebuild.