-- +migrate Up
CREATE TABLE IF NOT EXISTS movie_casts (
    credit_id VARCHAR(50) primary key,
    movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
    person_id INTEGER REFERENCES credits (id) ON DELETE CASCADE,
    cast_id INTEGER,
    character VARCHAR(350),
    cast_order INTEGER
);
