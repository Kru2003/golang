-- +migrate Down
ALTER TABLE movie_crew
DROP CONSTRAINT unique_movie_credit;
