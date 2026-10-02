package code

import "github.com/gauxs/lld/problems/auto_complete/extensions/baseline/code/enum"

// Autocomplete provides the functionality of searching similar words given the search word
type Autocomplete struct {
	similarMaxCount int
	sdb             *SearchDB
	order           enum.SearchOrder
}

// SetSearchOrder updates the search order
func (ac *Autocomplete) SetSearchOrder(newSearchOrder enum.SearchOrder) error {
	return nil
}

// AddWord adds a word to the underlying datastore
func (ac *Autocomplete) AddWord(word string) error {
	return nil
}

// SearchSimilar provides the seach functionality
func (ac *Autocomplete) SearchSimilar(searchWord string) ([]string, error) {
	return nil, nil
}
