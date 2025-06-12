-- +migrate Up
ALTER TABLE movie_casts ADD CONSTRAINT unique_movie_casts UNIQUE (movie_id, person_id);
