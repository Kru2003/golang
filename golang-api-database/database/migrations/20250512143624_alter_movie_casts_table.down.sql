-- +migrate Down
ALTER TABLE movie_casts
DROP CONSTRAINT unique_movie_casts;
