# Problem name

## Approach

Summarize the design and why it fits the requirements.

## Class design & Relationships

```go
// Example owns the state and coordinates the main workflow.
type Example struct {
    repository *Repository
}

// Repository owns storage and the access patterns over it.
type Repository struct {
    records map[string]Record
}

// NewExample creates an initialized Example.
func NewExample(repository *Repository) *Example

// Create validates and stores a record.
func (e *Example) Create(input Input) (Record, error)

// Get returns a record by id.
func (e *Example) Get(id string) (Record, error)
```

## Main flow

1. Validate the request.
2. Read or update the relevant state.
3. Return the result.

## Invariants and concurrency

Explain shared state, invariants, and synchronization.

## Complexity

- Write: `O(?)`
- Read: `O(?)`
- Space: `O(?)`

## Trade-offs

- Decision and benefit
- Limitation and possible future improvement
