-- +migrate Up
CREATE TABLE IF NOT EXISTS movie_genres (
    movieid INTEGER REFERENCES movies (id) on delete cascade,
    genreid INTEGER REFERENCES genres (id) on delete cascade,
    PRIMARY KEY (movieid, genreid)
);
