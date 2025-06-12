-- +migrate Up
CREATE TABLE IF NOT EXISTS credits (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255),
    gender INTEGER,
    profile_path VARCHAR(255)
);
