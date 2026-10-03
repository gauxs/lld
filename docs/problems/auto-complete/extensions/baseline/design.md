---
title: Search autocomplete — baseline (design)
problem: auto_complete
extension: baseline
prev:
  text: requirements
  link: /problems/auto-complete/extensions/baseline/requirements
next:
  text: codebase
  link: /problems/auto-complete/extensions/baseline/codebase
---

## Approach

`Autocomplete` is the public API. It delegates word storage and prefix lookup to an in-memory trie owned by `SearchDB`, while `SearchOrder` selects the ranking strategy used for suggestions.

## Go design sketch

Method bodies are intentionally omitted. Comments describe ownership,
responsibilities, and invariants.

```go
// Autocomplete coordinates additions and searches. It owns the result limit
// and active ranking order.
type Autocomplete struct {
    rwMu            *sync.RWMutex
    similarMaxCount int
    searchDB        *SearchDB
    order           enum.SearchOrder
}

// SearchDB owns the in-memory trie. One lock protects structural changes and
// traversals so readers never observe a partially inserted word.
type SearchDB struct {
    mu   *sync.Mutex
    root *SearchNode
}

// SearchNode represents one character in the trie. Terminal nodes retain the
// complete normalized word and its explicit-addition frequency. Each node
// caches its highest-ranked descendants for low-latency prefix reads.
type SearchNode struct {
    letter         string
    isLastLetter   bool
    entryFrequency int
    completeEntry  string
    children       []*SearchNode
    topSuggestions []RankedWord
}

type RankedWord struct {
    word      string
    frequency int
}

func NewAutocomplete() *Autocomplete

// SetSearchOrder replaces the ranking order used by subsequent searches.
func (a *Autocomplete) SetSearchOrder(order enum.SearchOrder) error

// SetMaxSuggestions configures the maximum number returned by a search.
func (a *Autocomplete) SetMaxSuggestions(limit int) error

// AddWord inserts or increments a normalized word.
func (a *Autocomplete) AddWord(word string) error

// SearchSimilar returns at most similarMaxCount words matching the prefix.
func (a *Autocomplete) SearchSimilar(prefix string) ([]string, error)

func NewSearchDB() *SearchDB

func (db *SearchDB) AddEntry(word string) error

func (db *SearchDB) SearchNSimilarByFrequency(
    prefix string,
    limit int,
) ([]string, error)
```

## Concurrency

- `Autocomplete` protects ranking configuration with `sync.RWMutex`.
- `SearchDB` serializes trie insertion and traversal with one mutex.

## Complexity and trade-offs

- Lookup is `O(prefix length + returned suggestions)` using the cached
  `topSuggestions` at the prefix node.
- Adding or incrementing a word updates cached rankings along its trie path,
  trading additional write work and memory for low-latency reads.
- The current reference implementation traverses and sorts the matching
  subtree; it must adopt this cache to meet the stated query target.
