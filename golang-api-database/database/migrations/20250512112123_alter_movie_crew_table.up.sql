-- +migrate Up
ALTER TABLE movie_crew ADD CONSTRAINT unique_movie_credit UNIQUE (movie_id, person_id);
