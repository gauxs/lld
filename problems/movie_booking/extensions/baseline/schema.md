```sql
CREATE TABLE cities (
    id VARCHAR(36) PRIMARY KEY,
    name NOT NULL
)

CREATE TABLE movies (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(36) NOT NULL
)

CREATE TABLE movie_theaters (
    id VARCHAR(36) PRIMARY KEY,
    city_id INT NOT NULL,
    FOREIGN KEY (city_id) REFERENCES cities(id)
)

CREATE TABLE screens (
    id VARCHAR(36) PRIMARY KEY,
    movie_theater_id INT NOT NULL,
    FOREIGN KEY (movie_theater_id) REFERENCES movie_theaters(id)
)

CREATE TABLE seats (
    id VARCHAR(36) PRIMARY KEY,
    row_no INT NOT NULL,
    seat_no INT NOT NULL,
    screen_id VARCHAR(36),
    seat_type VARCHAR(20) NOT NULL, -- 'SILVER', 'GOLD', 'PLATINUM'
    FOREIGN KEY (screen_id) REFERENCES screens(id)
)

CREATE TABLE shows (
    id VARCHAR(36) PRIMARY KEY,
    movie_id VARCHAR(36) NOT NULL,
    screen_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id),
    FOREIGN KEY (screen_id) REFERENCES screens(id)
)

CREATE TABLE show_seats (
    id VARCHAR(36) PRIMARY KEY,
    show_id VARCHAR(36) NOT NULL,
    seat_id VARCHAR(36) NOT NULL,
    seat_status VARCHAR(20) NOT NULL, -- 'AVAILABLE', 'LOCKED', 'BOOKED'
    locked_by_user_id VARCHAR(36) NULL,
    locked_at TIMESTAMP NULL,
    lock_version INT NOT NULL DEFAULT 0,
    FOREIGN KEY (show_id) REFERENCES shows(id),
    FOREIGN KEY (seat_id) REFERENCES seats(id),
    UNIQUE(show_id, seat_id)
)

CREATE TABLE bookings (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    show_id VARCHAR(36) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'PENDING', 'CONFIRMED', 'EXPIRED', 'CANCELLED'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (show_id) REFERENCES shows(id)
)
```