-- +migrate Up
CREATE TABLE IF NOT EXISTS languages (iso_code VARCHAR(2) PRIMARY KEY, name VARCHAR(50));
