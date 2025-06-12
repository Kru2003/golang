-- +migrate Up
CREATE TABLE IF NOT EXISTS movie_crew (
    credit_id VARCHAR(50) primary key,
    movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
    person_id INTEGER REFERENCES credits (id) ON DELETE CASCADE,
    department VARCHAR(100),
    job VARCHAR(100)
);
