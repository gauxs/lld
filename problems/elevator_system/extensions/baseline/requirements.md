# Elevator System

Design and implement a normal passenger elevator system like you'd see in a gated community/residential apartment complex.

## Functional requirements

### FR-1: Number of elevators & Floors

- Number of elevators and floors: configurable at system initialization
- There is a hall/lobby on every floor containing all elevator entrances.

### FR-2: User interaction

- Outside: each floor has Up and Down buttons.
- Inside: the elevator has a button for each floor.

### FR-3: Elevator Selection

- Nearest non-moving elevator as the initial selection strategy. We won't optimize/batch multiple requests in the initial implementation.
- Selection logic should be replaceable/extensible in the future

### FR-4: Elevator movement
- An elevator moves one floor at a time toward its destination.
- While handling a request, it does not accept another request.

### FR-4: Elevator Details

- Each elevator has its own external display, showing that elevator's current floor.
- The floor displays outside show the current floor of the elevator(s) serving that floor.
- Each elevator has a fixed capacity.
- The system must prevent boarding beyond that capacity.
- Capacity is considered when handling passengers/requests.

### FR-4: Multiple Elevators

- Multiple elevators can be moving simultaneously; their movements are independent.

## Out of scope (this extension)

- Deferred capability.

## Non-functional requirements

### NFR-1: Concurrent operation

- Multiple elevators can operate independently/concurrently.
- The simulation should be capable of representing their independent movement.
