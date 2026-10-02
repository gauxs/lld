package code

import (
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
	curSearchNode.entryFrequency++

	return nil
}

// SearchNSimilarByFrequency searches similar entity priortizing by frequency
func (sdb *SearchDB) SearchNSimilarByFrequency(searchWord string, n int) ([]string, error) {
	return nil, nil
}

// SearchNode represents a single character of a word
type SearchNode struct {
	letter         string
	isLastLetter   bool
	entryFrequency int
	childrens      []*SearchNode
}

func NewSearchNode(letter string, isWord bool) *SearchNode {
	return &SearchNode{
		letter:         letter,
		isLastLetter:   isWord,
		entryFrequency: 0,
		childrens:      make([]*SearchNode, 26),
	}
}
