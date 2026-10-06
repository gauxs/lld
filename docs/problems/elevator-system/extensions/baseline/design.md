---
title: Elevator System
problem: elevator_system
extension: baseline
prev:
  text: requirements
  link: /problems/elevator-system/extensions/baseline/requirements
next:
  text: codebase
  link: /problems/elevator-system/extensions/baseline/codebase
---

## Approach

Summarize the design and why it fits the requirements.

## Class design & Relationships

```go
type ElevatorMovementState int

const (
    ELEVATORMOVEMENTSTATE_INVALID iota = 0
    ELEVATORMOVEMENTSTATE_IDLE
)

type Elevator struct {
    capacity int
    passengers []*User
    currentFloor int
    movementState ElevatorMovementState
    requests chan *ElevatorRequest
}

type ElevatorRequest struct {
    sourceFloor int
    destinationFloor int
}

type ElevatorManagementSystem struct {
    elevators []*Elevator
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
