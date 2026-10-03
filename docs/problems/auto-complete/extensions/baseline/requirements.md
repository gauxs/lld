---
title: Search autocomplete
pageClass: lld-req-page
problem: auto_complete
extension: baseline
next:
  text: design
  link: /problems/auto-complete/extensions/baseline/design
---

Design and implement a Search Autocomplete System

## Functional requirements

### FR-1: Return suggestions

Return suggestions when a user types a prefix. The number of suggestions should be configurable.

### FR-2: Explicit word additions

Words are explicitly added to the system. Search queries do not automatically add new words.

### FR-3: Frequency ranking

Suggestions are ranked by frequency of explicit additions: a word that has been added more frequently ranks higher for a given prefix.

### FR-4: Replaceable ranking strategy

The ranking strategy should be replaceable/configurable

### FR-5: Repeated additions

Adding the same word multiple times increments its frequency rather than creating duplicate entries.

### FR-6: Empty query

Empty query should be supported; return the top-N suggestions globally.

### FR-7: Non-alphabetic input

Non-alphabetic input: reject/ignore characters outside a-z after normalization.

## Out of scope (this extension)

- None stated.

## Non-functional requirements

### NFR-1: Concurrent additions and retrievals

concurrent additions and retrievals are supported.

### NFR-2: Query latency

Query latency: autocomplete retrieval should be low-latency; target is O(query length + number of returned suggestions)

### NFR-3: Updated ranking

When a word is added or its frequency changes, subsequent retrievals should reflect the updated ranking without requiring a full rebuild.

## Extensions from here

- _(none yet)_