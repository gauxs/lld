# Movie Booking

Design a Movie Ticket Booking System that enables users to search for movies, view showtimes, select seats, and book tickets at cinemas or multiplexes

## Functional requirements

### FR-1: Users can search for shows based on a movie title and a city.

### FR-2: The system should support multiple cities, cinemas, screens, and shows.

### FR-3: Each screen has a defined layout of seats with different types (e.g., REGULAR, PREMIUM, RECLINER).

### FR-4: A user can book one or more available seats for a specific show. Double booking should be prevented.

### FR-5: The ticket price should be calculated dynamically based on configurable rules (e.g., seat types)

### FR-6: The system must be flexible to support different payment methods.

## Non-functional requirements

### NFR-1: Concurrency: The system must be designed to handle concurrent booking requests gracefully, ensuring data integrity and preventing race conditions like double-booking.

### NFR-2: Extensibility: The design should be modular. It should be easy to add new pricing strategies (e.g., holiday pricing) or new payment methods without significant changes to the core system.

### NFR-3: Modularity: The system should follow good object-oriented principles with a clear separation of concerns.

### NFR-4: Simplified Interface: The system should expose a simple API for clients to interact with, hiding the underlying complexity of the booking, locking, and payment processes

## Out of scope (this extension)