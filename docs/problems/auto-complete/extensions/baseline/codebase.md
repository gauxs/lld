---
title: Auto complete — codebase
problem: auto_complete
extension: baseline
prev:
  text: design
  link: /problems/auto-complete/extensions/baseline/design
---

# Codebase

Source: [GitHub](https://github.com/gauxs/lld/tree/main/problems/auto_complete/extensions/baseline/code)

Path: ` problems/auto_complete/extensions/baseline/code `

## Directory structure

```text
code/
├── enum/
│   └── search_order.go
├── auto_complete.go
├── error.go
├── orchestrator.go
├── search_db.go
└── search_db_test.go
```

## ` auto_complete.go `

```go
package code

import (
	"sync"

	"github.com/gauxs/lld/problems/auto_complete/extensions/baseline/code/enum"
)

// Autocomplete provides the functionality of searching similar words given the search word
type Autocomplete struct {
	rwMu            *sync.RWMutex
	similarMaxCount int
	sdb             *SearchDB
	order           enum.SearchOrder
}

func NewAutocomplete() *Autocomplete {
	return &Autocomplete{
		rwMu:            &sync.RWMutex{},
		similarMaxCount: 3,
		sdb:             NewSearchDB(),
		order:           enum.SEARCHORDER_FREQUENCY,
	}
}

// SetSearchOrder updates the search order
func (ac *Autocomplete) SetSearchOrder(newSearchOrder enum.SearchOrder) error {
	ac.rwMu.Lock()
	defer ac.rwMu.Unlock()

	ac.order = newSearchOrder
	return nil
}

// AddWord adds a word to the underlying datastore
func (ac *Autocomplete) AddWord(word string) error {
	return ac.sdb.AddEntry(word)
}

// SearchSimilar provides the seach functionality
func (ac *Autocomplete) SearchSimilar(searchWord string) ([]string, error) {
	ac.rwMu.RLock()
	defer ac.rwMu.RUnlock()

	switch ac.order {
	case enum.SEARCHORDER_FREQUENCY:
		return ac.sdb.SearchNSimilarByFrequency(searchWord, ac.similarMaxCount)
	}

	return nil, nil
}
```

## ` enum/search_order.go `

```go
package enum

// SearchOrder define how the results of a similar search given order.
type SearchOrder int

const (
	SEARCHORDER_INVALID SearchOrder = iota
	SEARCHORDER_FREQUENCY
)
```

## ` error.go `

```go
package code

import "errors"

var ErrInvalidEntryLength = errors.New("invalid entry length")
```

## ` orchestrator.go `

```go
package code

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gauxs/lld/problems/auto_complete/extensions/baseline/code/enum"
)

func Execute() {
	ac := NewAutocomplete()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Autocomplete System CLI ===")
	fmt.Println("Commands:")
	fmt.Println("  add <word>          - Add a new word to the system")
	fmt.Println("  search <prefix>     - Search similar words matching the prefix")
	fmt.Println("  mode <order_type>   - Set search order (e.g., frequency)")
	fmt.Println("  exit                - Exit the program")
	fmt.Println("===============================")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Split command action from arguments
		parts := strings.Fields(line)
		command := strings.ToLower(parts[0])
		args := parts[1:]

		switch command {
		case "exit", "quit":
			fmt.Println("Exiting Autocomplete CLI. Goodbye!")
			return

		case "add":
			if len(args) < 1 {
				fmt.Println("❌ Error: Missing word. Usage: add <word>")
				continue
			}
			word := args[0]
			if err := ac.AddWord(word); err != nil {
				fmt.Printf("❌ Error adding word: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully added word: %q\n", word)
			}

		case "search":
			// If no argument is provided, default to an empty string to return top N overall
			prefix := ""
			if len(args) >= 1 {
				prefix = args[0]
			}

			results, err := ac.SearchSimilar(prefix)
			if err != nil {
				fmt.Printf("❌ Search failed: %v\n", err)
				continue
			}

			if len(results) == 0 {
				fmt.Println("🔍 No matching similarities found.")
			} else {
				fmt.Printf("🔍 Top Results: %s\n", strings.Join(results, ", "))
			}

		case "mode":
			if len(args) < 1 {
				fmt.Println("❌ Error: Specify a mode. Usage: mode frequency")
				continue
			}

			var targetOrder enum.SearchOrder
			modeStr := strings.ToLower(args[0])

			switch modeStr {
			case "frequency":
				targetOrder = enum.SEARCHORDER_FREQUENCY
			default:
				fmt.Printf("❌ Unknown mode %q. Defaulting or failing.\n", modeStr)
				continue
			}

			if err := ac.SetSearchOrder(targetOrder); err != nil {
				fmt.Printf("❌ Failed to change search order: %v\n", err)
			} else {
				fmt.Printf("⚙️  Search mode updated to: %s\n", modeStr)
			}

		default:
			fmt.Printf("❌ Unknown command: %q. Try 'add', 'search', 'mode', or 'exit'.\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
	}
}
```

## ` search_db.go `

```go
package code

import (
	"cmp"
	"slices"
	"strings"
	"sync"
)

const (
	LetterAASCIIValue = int('a')
)

// SearchDB is the in memeory datastore and retrieval entity
type SearchDB struct {
	mu   *sync.Mutex // we cant use per SearchNode mutex because it can cause deadlock
	root *SearchNode
}

func NewSearchDB() *SearchDB {
	return &SearchDB{
		mu:   &sync.Mutex{},
		root: NewSearchNode("", false),
	}
}

// AddEntry adds an entry in the DB
func (sdb *SearchDB) AddEntry(entry string) error {
	if len(entry) == 0 {
		return ErrInvalidEntryLength
	}

	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	entry = strings.ToLower(entry)
	curSearchNode := sdb.root
	for wIdx := 0; wIdx < len(entry); wIdx++ {
		charIdx := int(entry[wIdx]) - LetterAASCIIValue
		if curSearchNode.childrens[charIdx] == nil {
			curSearchNode.childrens[charIdx] = NewSearchNode(string(entry[wIdx]), false)
		}

		curSearchNode = curSearchNode.childrens[charIdx]
	}

	curSearchNode.isLastLetter = true
	curSearchNode.completeEntry = entry
	curSearchNode.entryFrequency++

	return nil
}

// SearchNSimilarByFrequency searches similar entity priortizing by frequency
func (sdb *SearchDB) SearchNSimilarByFrequency(searchEntry string, n int) ([]string, error) {
	searchEntry = strings.ToLower(searchEntry)

	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	curSearchNode := sdb.root
	for wIdx := 0; wIdx < len(searchEntry); wIdx++ {
		charIdx := int(searchEntry[wIdx]) - LetterAASCIIValue
		if curSearchNode.childrens[charIdx] == nil {
			// no word with this search prefix exists
			return []string{}, nil
		}

		curSearchNode = curSearchNode.childrens[charIdx]
	}

	// 1. Collect all matching leaf/word nodes downstream
	allNodes := sdb.listAllEntries(curSearchNode)

	// 2. Sort the slice in DECREASING frequency order
	slices.SortFunc(allNodes, func(a, b *SearchNode) int {
		return cmp.Compare(b.entryFrequency, a.entryFrequency)
	})

	limit := n
	if len(allNodes) < limit {
		limit = len(allNodes)
	}

	output := make([]string, limit)
	for i := 0; i < limit; i++ {
		output[i] = allNodes[i].completeEntry
	}

	return output, nil
}

func (sdb *SearchDB) listAllEntries(curNode *SearchNode) []*SearchNode {
	results := make([]*SearchNode, 0)

	for idx := 0; idx < len(curNode.childrens); idx++ {
		if curNode.childrens[idx] != nil {
			results = append(results, sdb.listAllEntries(curNode.childrens[idx])...)
		}
	}

	if curNode.isLastLetter {
		results = append(results, curNode)
	}

	return results
}

// SearchNode represents a single character of a word
type SearchNode struct {
	letter         string
	isLastLetter   bool
	entryFrequency int
	completeEntry  string
	childrens      []*SearchNode
}

func NewSearchNode(letter string, isWord bool) *SearchNode {
	return &SearchNode{
		letter:         letter,
		isLastLetter:   isWord,
		entryFrequency: 0,
		completeEntry:  "",
		childrens:      make([]*SearchNode, 26),
	}
}
```

## ` search_db_test.go `

```go
package code

import "testing"

func TestAddEntry(t *testing.T) {
	sdb := NewSearchDB()

	if err := sdb.AddEntry("apple"); err != nil {
		t.Errorf("%v | unable to add entry in DB: %v", t.Name(), err)
	}

	if err := sdb.AddEntry(""); err == nil {
		t.Errorf("%v | expected error to occur when adding empty entry in DB", t.Name())
	}
}

func TestSearchNSimilarByFrequency(t *testing.T) {
	sdb := NewSearchDB()

	// 1. Seed test data using a clean, readable loop
	entries := []string{"apple", "appleOne", "appleTwo", "apple", "appleOne", "apple"}
	for _, entry := range entries {
		if err := sdb.AddEntry(entry); err != nil {
			t.Fatalf("failed to add entry %q to DB: %v", entry, err)
		}
	}

	// 2. Execute search
	searchEntry := "app"
	threshold := 3
	retEntries, err := sdb.SearchNSimilarByFrequency(searchEntry, threshold)
	if err != nil {
		t.Fatalf("failed to search entry %q: %v", searchEntry, err)
	}

	// 3. Assert results match expected length and content
	// Note: Adjust the expected order slice based on how you simulated frequencies in your setup!
	expectedOrder := []string{"apple", "appleone", "appletwo"}

	if len(retEntries) != len(expectedOrder) {
		t.Errorf("expected %d results, got %d", len(expectedOrder), len(retEntries))
	}

	for i, entry := range retEntries {
		if i < len(expectedOrder) && entry != expectedOrder[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expectedOrder[i], entry)
		}
	}

	searchEntry = ""
	retEntries, err = sdb.SearchNSimilarByFrequency(searchEntry, threshold)
	if err != nil {
		t.Fatalf("failed to search entry %q: %v", searchEntry, err)
	}

	// 3. Assert results match expected length and content
	// Note: Adjust the expected order slice based on how you simulated frequencies in your setup!
	expectedOrder = []string{"apple", "appleone", "appletwo"}

	if len(retEntries) != len(expectedOrder) {
		t.Errorf("expected %d results, got %d", len(expectedOrder), len(retEntries))
	}

	for i, entry := range retEntries {
		if i < len(expectedOrder) && entry != expectedOrder[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expectedOrder[i], entry)
		}
	}

	searchEntry = ""
	threshold = 2
	retEntries, err = sdb.SearchNSimilarByFrequency(searchEntry, threshold)
	if err != nil {
		t.Fatalf("failed to search entry %q: %v", searchEntry, err)
	}

	// 3. Assert results match expected length and content
	// Note: Adjust the expected order slice based on how you simulated frequencies in your setup!
	expectedOrder = []string{"apple", "appleone"}

	if len(retEntries) != len(expectedOrder) {
		t.Errorf("expected %d results, got %d", len(expectedOrder), len(retEntries))
	}

	for i, entry := range retEntries {
		if i < len(expectedOrder) && entry != expectedOrder[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expectedOrder[i], entry)
		}
	}
}
```
