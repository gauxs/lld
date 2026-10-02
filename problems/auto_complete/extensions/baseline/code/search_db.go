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
