package enum

// SearchOrder define how the results of a similar search given order.
type SearchOrder int

const (
	SEARCHORDER_INVALID SearchOrder = iota
	SEARCHORDER_FREQUENCY
)
