package code

// SearchDB is the in memeory datastore and retrieval entity
type SearchDB struct {
	root []*SearchNode
}

// AddEntry adds an entry in the DB
func (sdb *SearchDB) AddEntry(word string) error {
	return nil
}

// SearchNSimilarByFrequency searches similar entity priortizing by frequency
func (sdb *SearchDB) SearchNSimilarByFrequency(searchWord string, n int) ([]string, error) {
	return nil, nil
}

// SearchNSimilarByLex searches similar entity priortizing by lex order
func (sdb *SearchDB) SearchNSimilarByLex(searchWord string, n int) ([]string, error) {
	return nil, nil
}

// SearchNode represents a single character of a word
type SearchNode struct {
	letter        string
	isWord        bool
	wordFrequency int
	childrens     []*SearchNode
}
