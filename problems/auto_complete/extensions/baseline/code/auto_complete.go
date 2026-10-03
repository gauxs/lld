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
